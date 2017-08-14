package pipemgr

import (
	"time"

	"github.com/Sirupsen/logrus"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	api_v1 "k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"

	"github.com/arlert/malcolm/model"
)

// BuildEngineInterface is interface for running workstep
type BuildEngineInterface interface {
	RunStep(step *model.WorkStep) error
	StopStep(step *model.WorkStep) error
	SetupVolumn(volumn *api_v1.PersistentVolumeClaim) error
	CleanupVolumn(volumn *api_v1.PersistentVolumeClaim) error
	Start(event chan BuildEvent)
	Shutdown()
}

// Engine is for running job
type Engine struct {
	client    *kubernetes.Clientset
	selector  meta_v1.ListOptions
	terminate bool
	shutdown  chan bool
}

// NewEngine init new engine with client
func NewEngine(client *kubernetes.Clientset) *Engine {
	return &Engine{
		client:    client,
		terminate: false,
		shutdown:  make(chan bool),
		selector: meta_v1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labels.Set(map[string]string{
				"malcolm": "malcolm-job",
			})).String(),
		},
	}
}

// Shutdown terminates the loop
func (e *Engine) Shutdown() {
	e.shutdown <- true
}

// Start start engine
func (e *Engine) Start(event chan BuildEvent) {
	for !e.terminate {
		podWatcher, jobWatcher, err := e.createWatchers()
		if err != nil {
			logrus.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}
		for e.runOnce(podWatcher, jobWatcher, event) {
		}
	}
}

func (e *Engine) createWatchers() (watch.Interface, watch.Interface, error) {
	podWatcher, err := e.client.Core().Pods(model.DefaultNameSpace).Watch(e.selector)
	if err != nil {
		return nil, nil, err
	}
	jobWatcher, err := e.client.BatchV1().Jobs(model.DefaultNameSpace).Watch(e.selector)
	if err != nil {
		return nil, nil, err
	}
	return podWatcher, jobWatcher, nil
}

func (e *Engine) runOnce(podWatcher, jobWatcher watch.Interface, event chan BuildEvent) bool {
	select {
	case ev, ok := <-podWatcher.ResultChan():
		if !ok {
			return false
		}
		switch ev.Type {
		case watch.Added:
		case watch.Deleted:
		case watch.Modified:
			pod := ev.Object.(*api_v1.Pod)
			logrus.Debugln("POD", ev.Type, pod.Name)
		}

	case ev, ok := <-jobWatcher.ResultChan():
		if !ok {
			return false
		}
		job, ok := ev.Object.(*batchv1.Job)
		if !ok {
			logrus.Debugln("watch unexpected object:", ev.Object)
			return true
		}
		var condition *batchv1.JobCondition
		if len(job.Status.Conditions) > 0 {
			condition = &job.Status.Conditions[len(job.Status.Conditions)-1]
		}
		logrus.Debugf("%s type : %s, active : %d, success : %d, fail : %d \n condition %d : %+v",
			job.Name, ev.Type, job.Status.Active, job.Status.Succeeded,
			job.Status.Failed, len(job.Status.Conditions), condition)

		switch ev.Type {
		case watch.Added:

			// add event canbe "initevent" not status:add
			for _, cond := range job.Status.Conditions {
				if cond.Type == batchv1.JobComplete {
					event <- BuildEvent{
						eventtype: eventStatusDone,
						buildid:   job.Labels["build"],
					}
					return true
				}
			}
			event <- BuildEvent{
				eventtype: eventStatusStart,
				buildid:   job.Labels["build"],
			}

		case watch.Modified:
			// complete
			for _, cond := range job.Status.Conditions {
				if cond.Type == batchv1.JobComplete {
					event <- BuildEvent{
						eventtype: eventStatusDone,
						buildid:   job.Labels["build"],
					}
					return true
				}
			}
			// failed
			threshold := int32(4)
			if job.Status.Failed > threshold {
				event <- BuildEvent{
					eventtype: eventError,
					data:      "job watch abort after fail 4 times",
					buildid:   job.Labels["build"],
				}
				return true
			}
		case watch.Deleted:
			event <- BuildEvent{
				eventtype: eventStatusDeleted,
				buildid:   job.Labels["build"],
			}
		case watch.Error:
			event <- BuildEvent{
				eventtype: eventError,
				buildid:   job.Labels["build"],
			}
		}
	case <-e.shutdown:
		e.terminate = true
		return false
	}
	return true
}

// RunStep run a job sync with message and cancel channel
func (e *Engine) RunStep(step *model.WorkStep) error {
	logrus.Debugln("RunStep", step.Title)
	job, err := e.client.Jobs(model.DefaultNameSpace).Create(step.K8sjob)
	_ = job
	logrus.Debugf("create:%+v", step.K8sjob.Spec.Template.Spec.Containers)
	if err != nil {
		return err
	}
	//step.K8sjob = job
	return nil
}

// StopStep set Parallelism to 0 to stop a step
func (e *Engine) StopStep(step *model.WorkStep) error {
	jobService := e.client.BatchV1().Jobs(model.DefaultNameSpace)
	if step.K8sjob.Spec.Parallelism != nil && *step.K8sjob.Spec.Parallelism > 0 {
		parallelism := int32(0)
		step.K8sjob.Spec.Parallelism = &parallelism
		job, err := jobService.Update(step.K8sjob)
		if err != nil {
			return err
		}
		step.K8sjob = job
	}
	return nil
}

// SetupVolumn create volumn
func (e *Engine) SetupVolumn(volumn *api_v1.PersistentVolumeClaim) (err error) {
	_, err = e.client.CoreV1().PersistentVolumeClaims(model.DefaultNameSpace).Create(volumn)
	return
}

// CleanupVolumn delete volumn
func (e *Engine) CleanupVolumn(volumn *api_v1.PersistentVolumeClaim) (err error) {
	err = e.client.CoreV1().PersistentVolumeClaims(model.DefaultNameSpace).Delete(volumn.Name, nil)
	return
}
