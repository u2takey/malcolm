package service

import (
	"io/ioutil"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/arlert/malcolm/model"
	"github.com/arlert/malcolm/server/pipemgr"
	"github.com/arlert/malcolm/store"
	req "github.com/arlert/malcolm/utils/reqlog"
)

var (
	bearer_token_file = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

type Service struct {
	config *model.Config
	store  *store.Store
	pipem  *pipemgr.PipeMgr
	engine *pipemgr.Engine
}

func New(cfg *model.Config) *Service {
	token := ""
	if bearer_token_file != "" {
		bf, err := ioutil.ReadFile(bearer_token_file)
		if err != nil {
			logrus.Error("read bearer_token err ", err)
		}
		token = string(bf)
	}
	if !strings.HasPrefix(cfg.KubeAddr, "http") {
		cfg.KubeAddr = "http://" + cfg.KubeAddr
	}
	resconfig := &rest.Config{
		Host:        cfg.KubeAddr,
		BearerToken: token,
	}
	resconfig.Insecure = true
	client, err := pipemgr.CreateK8sClientByConfig(resconfig)
	if err != nil {
		logrus.Fatalln("CreateK8sClientByConfig fail", err)
	}
	svc := &Service{
		config: cfg,
		store:  store.New(&cfg.MgoCfg),
		engine: pipemgr.NewEngine(client),
	}
	svc.pipem = pipemgr.NewPipeMgr(svc.store, svc.engine)
	return svc
}

func (s *Service) GetPing(c *gin.Context) {
	req.Entry(c).Debug("GetPing")
	c.String(200, "pong")
}
