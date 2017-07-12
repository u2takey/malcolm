package pipemgr

import (
	"time"

	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	"labix.org/v2/mgo/bson"

	"github.com/arlert/malcolm/model"
)

type buildtemplate struct {
}

func (t *buildtemplate) ConfigToBuild(pipe *model.PipeConfig) *model.Build {
	build := &model.Build{
		ID:              bson.NewObjectId(),
		PipeID:          pipe.ID,
		PipeName:        pipe.Name,
		Title:           pipe.Name,
		Description:     pipe.Description,
		ConcurrentBuild: pipe.ConcurrentBuild,
		Created:         time.Now(),
	}
	work := model.Work{
		Title:       pipe.Name + "-0",
		Description: "",
	}

	for _, cfg := range pipe.Tasks {
		k8sjob := &batchv1.Job{
			Spec: batchv1.JobSpec{
				// Template: v1.PodTemplate{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						RestartPolicy: v1.RestartPolicyNever,
						Containers: []v1.Container{
							v1.Container{
								Name:    build.ID.Hex(),
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

	// single work for now
	build.Works = append(build.Works, work)
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
