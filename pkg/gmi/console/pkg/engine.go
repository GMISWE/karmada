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

	api "github.com/karmada-io/karmada/pkg/gmi/console/logic/apis"
	"github.com/karmada-io/karmada/pkg/gmi/console/logic/jobs"
	"github.com/karmada-io/karmada/pkg/gmi/console/logic/models"
	"github.com/karmada-io/karmada/pkg/gmi/console/pkg/core"
	"github.com/piaobeizu/titan"
	titanConfig "github.com/piaobeizu/titan/config"
	"github.com/piaobeizu/titan/service"
	"github.com/piaobeizu/titan/storage"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type Engine struct {
	name    string
	ctx     context.Context
	titan   *titan.Titan
	console *core.GmiConsole
}

func NewEngine(ctx context.Context, name, config string, kubeClientSet *kubernetes.Clientset, dynamicClientSet *dynamic.DynamicClient) *Engine {
	if err := titanConfig.InitCfg(config); err != nil {
		panic(err)
	}
	cfg := titanConfig.GetConfig()

	// start titan engine
	console, err := core.NewGmiConsole(ctx)
	if err != nil {
		panic(err)
	}
	titan := titan.NewTitan(ctx, name, cfg.LogMode).
		ApiServer(cfg.Http.ApiAddr, cfg.Http.Version).
		Scheduler().
		Handler(&api.Handler{Console: console, KubeClientSet: kubeClientSet, DynamicClientSet: dynamicClientSet}).
		WSHandler(&api.WSHandler{Console: console, KubeClientSet: kubeClientSet, DynamicClientSet: dynamicClientSet}).
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
				},
			})
		}
	}
	e := &Engine{
		name:    name,
		ctx:     ctx,
		titan:   titan,
		console: console,
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
	go e.console.Start()
	e.titan.Start()
}

func (e *Engine) Stop() {
	e.titan.Stop()
	logrus.Printf("%s stopped, byebye!", e.name)
}
