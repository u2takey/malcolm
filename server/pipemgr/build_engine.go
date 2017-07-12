package pipemgr

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"k8s.io/client-go/rest"

	"github.com/arlert/malcolm/model"
)

const (
	defaultQPS   = 1e6
	defaultBurst = 1e6

	defaultNameSpace = "malcolm" // #todo ->param
	defaultTimeout   = 60 * time.Minute
)

type Engine struct {
	client *kubernetes.Clientset
}

func NewEngine(client *kubernetes.Clientset) *Engine {
	return &Engine{client: client}
}

func CreateK8sClientByConfig(cfg *rest.Config) (*kubernetes.Clientset, error) {
	if cfg.QPS == 0 {
		cfg.QPS = defaultQPS
	}
	if cfg.Burst == 0 {
		cfg.Burst = defaultBurst
	}
	if cfg.ContentType == "" {
		cfg.ContentType = "application/vnd.kubernetes.protobuf"
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getlabel(job *batchv1.Job) string {
	labels := []string{}
	for k, v := range job.ObjectMeta.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(labels, ",")
}

func (e *Engine) RunSyncJobch(work *model.WorkStep, msg chan<- *meassge) {
	logrus.Debug("RunSyncJobch")
	job, err := e.client.Jobs(defaultNameSpace).Create(work.K8sjob)
	_ = job
	logrus.Debugf("create:%+v", work.K8sjob.Spec.Template.Spec.Containers)
	if err != nil {
		msg <- &meassge{
			err: err,
		}
		return
	}
	timeout := int64(defaultTimeout)
	opts := meta_v1.ListOptions{
		Watch:          true,
		LabelSelector:  getlabel(work.K8sjob),
		TimeoutSeconds: &timeout,
	}
	watcher, err := e.client.Jobs(defaultNameSpace).Watch(opts)
	if err != nil {
		msg <- &meassge{
			err: err,
		}
		return
	}

	for {
		select {
		case watchEvent, open := <-watcher.ResultChan():
			if !open {
				return
			}

			switch object := watchEvent.Object.(type) {
			case *batchv1.Job:
				logrus.Debug(fmt.Sprintf("type : %s, active : %d, success : %d, fail : %d",
					watchEvent.Type, object.Status.Active, object.Status.Succeeded, object.Status.Failed))
				for _, cond := range object.Status.Conditions {
					if cond.Type == batchv1.JobComplete {
						msg <- &meassge{
							data: "job complete",
						}
						return
					}
				}

				threshold := int32(4)
				if object.Status.Failed > threshold {
					msg <- &meassge{
						data: "job watch abort after fail 4 times",
					}
					return
				}
			case *meta_v1.Status:
				logrus.Debug(object.Status)
			}
		case <-time.After(defaultTimeout):
			msg <- &meassge{
				err: errors.New("watch time out"),
			}
			return
		}
	}
}