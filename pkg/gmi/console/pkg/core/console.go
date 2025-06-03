/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/29
 * @Desc    : gmi agent
 */

package core

import (
	"context"
)

type GmiConsole struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var (
	console *GmiConsole
)

func NewGmiConsole(ctx context.Context) (*GmiConsole, error) {
	if console == nil {
		ctx, cancel := context.WithCancel(ctx)
		console = &GmiConsole{
			ctx:    ctx,
			cancel: cancel,
		}
	}
	return console, nil
}

func (a *GmiConsole) Start() {
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				return
			}
		}
	}()
}
