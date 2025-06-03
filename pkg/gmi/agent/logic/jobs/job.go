/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/01/19 16:01:44
 Desc     :
*/

package jobs

import (
	"context"
	"time"

	"github.com/karmada-io/karmada/pkg/gmi/agent/logic/config"
	"github.com/karmada-io/karmada/pkg/gmi/agent/pkg/core"
	"github.com/piaobeizu/titan/service"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Time string `json:"time"`
}

type Job struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Name   string
	Detail string
	Status service.JobStatus
	Args   map[string]interface{}
	Agent  *core.GmiAgent
	params *config.Business
	start  time.Time
}

func (j *Job) Demo() {
	j.init()
	defer j.end()
	logrus.Infof("job %s running with args: %+v", j.Name, j.Args)
}

func (j *Job) init() {
	logrus.Printf("[job] %s: %s", j.Name, j.Detail)
	j.Status = service.JobStatusRunning
	j.params = config.GetConfig().Business
	j.start = time.Now().Local()
}

func (j *Job) end() {
	j.Status = service.JobStatusDone
	logrus.Infof("[job] %s done, cost time: %s", j.Name, time.Now().Local().Sub(j.start).String())
}
