package jobmgr

import (
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"labix.org/v2/mgo/bson"

	"github.com/u2takey/malcolm/model"
)

type buildtemplate struct {
}

func (t *buildtemplate) ConfigToBuild(job *model.JobConfig) *model.Build {
	build := &model.Build{
		ID:              bson.NewObjectId(),
		JobID:           job.ID,
		JobName:         job.Name,
		Title:           job.Name,
		Description:     job.Description,
		ConcurrentBuild: job.ConcurrentBuild,
		Created:         time.Now(),
	}
	work := model.Work{
		Title:       job.Name + "-0",
		Description: "",
	}

	workconfigs := map[string][]model.WorkConfig{
		"scm":        job.Scm,
		"prebuild":   job.PreBuild,
		"build":      job.Build,
		"afterbuild": job.AfterBuild,
		"notify":     job.Notify,
		"service":    job.Service,
	}

	for _, val := range workconfigs {
		for _, cfg := range val {
			k8sjob := &batchv1.Job{
				Spec: batchv1.JobSpec{
					// Template: v1.PodTemplate{
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:    "",
									Image:   cfg.Plugin,
									Env:     toEnvSlice(cfg.Environment),
									Command: cfg.Command,
									Args:    cfg.Args,
								},
							},
						},
						// },pbuild_engine.go:71:
					},
				},
			}
			k8sjob.Name = bson.NewObjectId().Hex()
			work.Steps = append(work.Steps, &model.WorkStep{
				Title:       cfg.Title,
				Description: cfg.Type + " : " + cfg.Plugin,
				Config:      &cfg,
				K8sjob:      k8sjob,
			})
		}
	}

	// single work for now
	build.Works = append(build.Works, work)
	b, _ := json.Marshal(build)
	logrus.Debug(string(b))
	return build
}

func toEnvSlice(envm map[string]string) (envl []v1.EnvVar) {
	for k, v := range envm {
		envl = append(envl, v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return
}
