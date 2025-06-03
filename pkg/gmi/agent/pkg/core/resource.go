/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/29
 * @Desc    : gmi gpu
 */

package core

import (
	"sync"

	"github.com/karmada-io/karmada/pkg/gmi/base/types"
	"github.com/karmada-io/karmada/pkg/util"
	"github.com/sirupsen/logrus"
)

type Resource struct {
	nodes map[string][]types.Node
	mu    sync.RWMutex
}

var resource *Resource

func NewResource() *Resource {
	if resource == nil {
		resource = &Resource{
			nodes: make(map[string][]types.Node),
			mu:    sync.RWMutex{},
		}
	}
	return resource
}

func (r *Resource) AddNode(node types.Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.nodes[node.UUID]; !ok {
		r.nodes[node.UUID] = make([]types.Node, 0)
	}
	r.nodes[node.UUID] = append(r.nodes[node.UUID], node)
}

func (r *Resource) GetNodes() map[string][]types.Node {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.nodes
}

func (r *Resource) PopNodes() map[string][]types.Node {
	r.mu.Lock()
	defer r.mu.Unlock()
	nodes := make(map[string][]types.Node)
	for uuid, node := range r.nodes {
		nodes[uuid] = node
		delete(r.nodes, uuid)
	}
	return nodes
}

func (r *Resource) CalTopo() types.Topo {
	nodes := r.PopNodes()

	logrus.Infof("nodes: %s", util.FormatStruct(nodes))

	return types.Topo{
		Cluster: "cluster",
		Region:  "region",
		Zone:    "zone",
		// Update:  time.Now(),
	}
}
