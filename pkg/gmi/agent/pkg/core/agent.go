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

	"github.com/karmada-io/karmada/pkg/gmi/base/types"
	"github.com/karmada-io/karmada/pkg/util"
	"github.com/sirupsen/logrus"
)

type GmiAgent struct {
	ctx      context.Context
	nfdwsmsg chan any
}

var (
	agent *GmiAgent
)

func NewGmiAgent(ctx context.Context) (*GmiAgent, error) {
	if agent == nil {
		agent = &GmiAgent{
			ctx:      ctx,
			nfdwsmsg: make(chan any, 1),
		}
	}
	return agent, nil
}

func (a *GmiAgent) Start() {
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				return
			case msg := <-a.nfdwsmsg:
				switch m := msg.(type) {
				case types.Node:
					resource := NewResource()
					resource.AddNode(m)
				default:
					logrus.Infof("unsupported ws msg: %s", util.FormatStruct(m))
				}
			}
		}
	}()
}

func (a *GmiAgent) SendNfdWSMsg(msg any) {
	a.nfdwsmsg <- msg
}
