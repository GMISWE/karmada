/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/09
 * @Desc    : map with goroutine safe
 */

package util

import (
	"sync"
)

type GMap struct {
	mu sync.RWMutex
	m  map[any]any
}

func (g *GMap) Set(key any, value any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.m[key] = value
}

func (g *GMap) Get(key any) (any, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	value, ok := g.m[key]
	return value, ok
}

func (g *GMap) Delete(key any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.m, key)
}

func (g *GMap) Range(f func(key any, value any) bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for key, value := range g.m {
		if !f(key, value) {
			break
		}
	}
}
