/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/29
 * @Desc    : gmi engine for karmada
 */

package pkg

import (
	"context"

	api "github.com/karmada-io/karmada/pkg/gmi/agent/logic/apis"
	"github.com/karmada-io/karmada/pkg/gmi/agent/logic/config"
	"github.com/karmada-io/karmada/pkg/gmi/agent/logic/jobs"
	"github.com/karmada-io/karmada/pkg/gmi/agent/logic/models"
	"github.com/karmada-io/karmada/pkg/gmi/agent/pkg/core"
	"github.com/piaobeizu/titan"
	"github.com/piaobeizu/titan/service"
	"github.com/piaobeizu/titan/storage"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	name  string
	ctx   context.Context
	titan *titan.Titan
	agent *core.GmiAgent
}

func NewEngine(ctx context.Context, name string) *Engine {
	if err := config.InitCfg(); err != nil {
		panic(err)
	}
	cfg := config.GetConfig()
	// log.InitLog(name, cfg.LogMode)

	// start titan engine
	agent, err := core.NewGmiAgent(ctx)
	if err != nil {
		panic(err)
	}
	titan := titan.NewTitan(ctx, name, cfg.LogMode).
		ApiServer(cfg.Http.ApiAddr, cfg.Http.Version).
		Scheduler().
		Handler(&api.Handler{}).
		WSHandler(&api.WSHandler{GMIAgent: agent}).
		Middleware(&service.ApiMiddleware{})
	// add middlewares and routers
	titan = titan.Middlewares(cfg.Http.Middlewares)
	for group, route := range cfg.Http.Routes {
		titan = titan.Routers(group, route.Middlewares, route.Routers, route.Sses, route.Websockets)
	}
	// add jobs
	for _, scheduler := range cfg.Schedulers {
		if scheduler.Enabled {
			subCtx, cancel := context.WithCancel(ctx)
			titan = titan.Job(&service.Job{
				Name:   scheduler.Name,
				Cron:   scheduler.Cron,
				Method: scheduler.Method,
				Runner: &jobs.Job{
					Ctx:    subCtx,
					Cancel: cancel,
					Name:   scheduler.Name,
					Detail: scheduler.Detail,
					Args:   scheduler.Args,
					Agent:  agent,
				},
			})
		}
	}
	e := &Engine{
		name:  name,
		ctx:   ctx,
		titan: titan,
		agent: agent,
	}

	// init db
	if len(cfg.Mysql) > 8 {
		for _, model := range models.Models() {
			storage.RegisterModel(model)
		}
		storage.Migrate(cfg.Mysql)
		logrus.Infof("init mysql succeed...")
	}
	// init redis
	if len(cfg.Redis) > 8 {
		storage.InitRedis(cfg.Redis)
		logrus.Infof("init redis succeed...")
	}
	// init oss
	if len(cfg.Oss) > 8 {
		logrus.Infof("init oss succeed...")
	}
	return e
}

func (e *Engine) Start() {
	logrus.Printf("%s starting...", e.name)
	go e.agent.Start()
	e.titan.Start()
}

func (e *Engine) Stop() {
	e.titan.Stop()
	logrus.Printf("%s stopped, byebye!", e.name)
}
