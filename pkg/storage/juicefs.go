/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.wang@gmicloud.com
 * @Date    : 2025/05/09
 * @Desc    : juicefs 存储类型实现
 */

package storage

import (
	"context"
	"fmt"
	"strings"

	astorage "github.com/karmada-io/karmada/pkg/apis/storage/v1alpha1"
	kcontainerd "github.com/karmada-io/karmada/pkg/containerd"
	"github.com/karmada-io/karmada/pkg/util"

	"github.com/opencontainers/runtime-spec/specs-go"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

type Juicefs struct {
	BaseStorage
	*astorage.Juicefs
}

// JuiceFSMountConfig 包含挂载 JuiceFS 所需的配置参数
type JuiceFSMountConfig struct {
	MOUNT_POINT          string
	STORAGE_PATH         string
	JUICEFS_TOKEN        string
	JUICEFS_NAME         string
	MOUNT_OPTIONS        string
	JUICEFS_CONSOLE_HOST string
	JUICEFS_META_URL     string
	JUICEFS_PATH         string
	JUICEFS_VERSION      string
	JUICEFS_CACHE_DIR    string
	JUICEFS_ACCESS_KEY   string
	JUICEFS_SECRET_KEY   string
}

func NewJuicefsFromRuntimeObject(ctx context.Context, obj runtime.Object) (*Juicefs, error) {
	subctx := context.WithoutCancel(ctx)
	unstructuredObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("failed to convert unstructured to juicefs")
	}
	juicefs := &astorage.Juicefs{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.Object, juicefs); err != nil {
		return nil, fmt.Errorf("failed to convert unstructured to juicefs: %s", err.Error())
	}
	return &Juicefs{
		Juicefs: juicefs,
		BaseStorage: BaseStorage{
			ctx: subctx,
		},
	}, nil
}

func (j *Juicefs) Mount() error {
	mountPoint := fmt.Sprintf("%s/%s", MOUNT_POINT, j.Spec.Provider.ID)
	jfsMountOptions := []string{}
	for _, opt := range j.Spec.Client.MountOptions {
		opts := strings.Split(opt, "=")
		if len(opts) == 1 {
			jfsMountOptions = append(jfsMountOptions, fmt.Sprintf("--%s", opts[0]))
		} else {
			jfsMountOptions = append(jfsMountOptions, fmt.Sprintf("--%s %s", opts[0], strings.Join(opts[1:], "=")))
		}
	}
	config := JuiceFSMountConfig{
		MOUNT_OPTIONS:     strings.Join(jfsMountOptions, " "),
		MOUNT_POINT:       mountPoint,
		JUICEFS_NAME:      j.Name,
		JUICEFS_CACHE_DIR: j.Spec.Client.CacheDir,
		JUICEFS_PATH:      fmt.Sprintf("%s/juicefs", kcontainerd.STORAGE_PATH),
		STORAGE_PATH:      kcontainerd.STORAGE_PATH,
	}
	if j.Spec.Client.EE != nil {
		config.JUICEFS_VERSION = j.Spec.Client.EE.Version
		config.JUICEFS_CONSOLE_HOST = j.Spec.Client.EE.Auth.BaseURL
		config.JUICEFS_TOKEN = j.Spec.Client.EE.Auth.Token
	} else {
		config.JUICEFS_VERSION = j.Spec.Client.CE.Version
		config.JUICEFS_META_URL = j.Spec.Client.CE.MetaURL
		config.JUICEFS_ACCESS_KEY = j.Spec.Client.CE.Backend.AccessKey
		config.JUICEFS_SECRET_KEY = j.Spec.Client.CE.Backend.SecretKey
	}
	if err := writeScript(j.Name, JUICEFS_MOUNT_SCRIPT, config); err != nil {
		klog.Errorf("failed to write script file: %s", err.Error())
		return err
	}
	if j.container == nil {
		// create storage path in host
		if err := util.RunCommand(j.ctx, "nsenter", "-t", "1", "-m", "-u", "-n", "-i", "-p", "mkdir", "-p", kcontainerd.STORAGE_PATH); err != nil {
			klog.Errorf("failed to create storage path: %s", err.Error())
			return err
		}
		j.container = kcontainerd.NewContainer(j.ctx).
			WithNamespace(CONTAINER_NAMESPACE).
			WithName(j.Name).
			WithPrivilege(true).
			WithUser("root").
			WithMounts(specs.Mount{
				Type:        "bind",
				Source:      kcontainerd.STORAGE_PATH,
				Destination: kcontainerd.STORAGE_PATH,
				Options:     []string{"bind", "rw"},
			}).
			WithImage(STORAGE_IMAGE).
			WithArgs([]string{}).
			WithAuth(&kcontainerd.Auth{
				Username:         "",
				Password:         "",
				InsecureRegistry: false,
			}).
			WithGCPCredentials(GCP_CREDENTIALS_SCRIPT).
			WithEnvs([]string{
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
				fmt.Sprintf("WATCH_PATH=%s", fmt.Sprintf("%s/%s.sh", kcontainerd.STORAGE_PATH, j.Name)),
				fmt.Sprintf("LOG_PATH=%s", fmt.Sprintf("%s/%s.log", kcontainerd.STORAGE_PATH, j.Name)),
			}).
			WithLogPath(fmt.Sprintf("%s/%s.log", kcontainerd.STORAGE_PATH, j.Name))
	}
	go func() {
		if err := cc.Run(j.container); err != nil {
			klog.Errorf("failed to run containerd container: %s", err.Error())
			cc.Delete(j.container)
			return
		}
		klog.Infof("container %s running, start log watcher", j.container.Name)
		j.container.Logs(func(line string) {
			klog.Infof("[%s] %s", j.Name, line)
		})
	}()
	return nil
}

func (j *Juicefs) Unmount() error {
	klog.Infof("unmounting storage %s", j.Name)
	if err := cc.Delete(j.container); err != nil {
		klog.Errorf("failed to delete containerd container: %s", err.Error())
		return err
	}
	j.container.Cancel()
	klog.Infof("storage %s unmounted", j.Name)
	return nil
}
