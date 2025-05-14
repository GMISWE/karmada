/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloua.ai
 * @Date    : 2025/05/09
 * @Desc    : 定义存储接口
 */

package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/containerd/containerd/errdefs"
	astorage "github.com/karmada-io/karmada/pkg/apis/storage/v1alpha1"
	"github.com/karmada-io/karmada/pkg/containerd"
	"github.com/karmada-io/karmada/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

var (
	cc *containerd.ContainerdClient
	// containers = map[string]*containerd.Container{}

	STORAGE_K8S_NAMESPACE = util.GetEnv("GMI_STORAGE_K8S_NAMESPACE", "gmi-storage")
	CONTAINERD_SOCKET     = util.GetEnv("CONTAINERD_SOCKET", "/run/containerd/containerd.sock")
	MOUNT_POINT           = util.GetEnv("GMI_STORAGE_MOUNT_POINT", "/mnt/juicefs")
	RETRY_COUNT           = util.GetEnv("GMI_STORAGE_RETRY_COUNT", "3")
	CONTAINER_NAMESPACE   = util.GetEnv("GMI_STORAGE_CONTAINER_NAMESPACE", "gmicloud.ai")
	STORAGE_IMAGE         = util.GetEnv("GMI_STORAGE_IMAGE", "us-west1-docker.pkg.dev/devv-404803/public/storage:v0.0.1")
)

type Storage interface {
	Mount() error
	Unmount() error
}

func Exit() {}

func Init(ctx context.Context, dynamicClientSet dynamic.Interface) error {
	// init containerd client
	if cc == nil {
		nc, err := containerd.NewContainerdClient(CONTAINERD_SOCKET)
		if err != nil {
			klog.Errorf("failed to create containerd client: %s", err.Error())
			return err
		}
		cc = nc
	}
	// TODO: more init logic
	return nil
}

func Watch(ctx context.Context, kubeClientSet kubernetes.Interface, dynamicClientSet dynamic.Interface) {
	// first check and update storages
	checkAndUpdateStorages := func() {
		storages, err := fetchAllStorages(ctx, dynamicClientSet)
		if err != nil {
			klog.Errorf("failed to fetch all storages: %s", err.Error())
			return
		}
		ctrs, err := cc.List(CONTAINER_NAMESPACE)
		if err != nil && !errdefs.IsNotFound(err) {
			klog.Errorf("failed to list containers: %s", err.Error())
			return
		}
		klog.Infof("storages: %d, containers: %d, new-mount: %v, new-unmount: %v", len(storages), len(ctrs), len(storages) > len(ctrs), len(storages) < len(ctrs))
		for _, storage := range storages {
			if err := storage.Mount(); err != nil {
				klog.Errorf("failed to mount storage: %s", err.Error())
				panic(fmt.Sprintf("failed to mount storage: %s", err.Error()))
			}
		}
		for name := range ctrs {
			if _, ok := storages[name]; !ok {
				klog.Infof("unmounting storage %s", name)
				if err := cc.Delete(ctrs[name]); err != nil {
					panic(fmt.Sprintf("failed to delete containerd container: %s", err.Error()))
				}
				ctrs[name].Cancel()
				delete(ctrs, name)
				klog.Infof("storage %s unmounted", name)
			}
		}
	}
	// check and update storages
	checkAndUpdateStorages()
	// ticker to check and update storages
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkAndUpdateStorages()
		}
	}
}

func writeScript(name, script string, config any) error {
	if _, err := os.Stat(containerd.STORAGE_PATH); os.IsNotExist(err) {
		if err := os.MkdirAll(containerd.STORAGE_PATH, 0755); err != nil {
			return err
		}
	}
	script, err := util.GenerateTemplate(script, name, config)
	if err != nil {
		return err
	}
	scriptFile := fmt.Sprintf("%s/%s.sh", containerd.STORAGE_PATH, name)
	if err := os.WriteFile(scriptFile, []byte(script), 0755); err != nil {
		return err
	}
	return nil
}

func fetchAllStorages(ctx context.Context, dynamicClientSet dynamic.Interface) (map[string]Storage, error) {
	storages := map[string]Storage{}
	// create resource request object
	gvrs := []schema.GroupVersionResource{
		{
			Group:    astorage.GroupName,             // "storage.karmada.io"
			Version:  astorage.GroupVersion.Version,  // "v1alpha1"
			Resource: astorage.ResourcePluralJuicefs, // "juicefs"
		},
		// TODO: add more resource types
	}
	for _, gvr := range gvrs {
		storageList, err := dynamicClientSet.Resource(gvr).Namespace(STORAGE_K8S_NAMESPACE).List(ctx, metav1.ListOptions{})
		if err != nil {
			klog.Errorf("failed to list storage: %s", err.Error())
			return nil, err
		}
		if storageList == nil {
			return nil, nil
		}
		for _, storage := range storageList.Items {
			if storage.GetKind() == astorage.ResourceKindJuicefs {
				sto, err := NewJuicefsFromRuntimeObject(ctx, &storage)
				if err != nil {
					klog.Errorf("failed to convert unstructured to juicefs: %s", err.Error())
					return nil, err
				}
				storages[sto.Name] = sto
			}
			// TODO: add more resource types
		}
	}
	return storages, nil
}

type StorageStatus string

const (
	StorageStatusInit       StorageStatus = "init"
	StorageStatusMounting   StorageStatus = "mounting"
	StorageStatusMounted    StorageStatus = "mounted"
	StorageStatusUnmounting StorageStatus = "unmounting"
	StorageStatusUnmounted  StorageStatus = "unmounted"
)
