package pipemgr

import (
	"bytes"
	"html/template"
	"io"
	"strconv"
	"time"
	"unicode"

	"k8s.io/apimachinery/pkg/util/yaml"
	batch_v1 "k8s.io/client-go/pkg/apis/batch/v1"
	"labix.org/v2/mgo/bson"

	"github.com/arlert/malcolm/model"
)

type buildtemplate struct {
}

func (t *buildtemplate) ConfigToBuild(pipe *model.PipeConfig) (*model.Build, error) {
	build := &model.Build{
		ID:          bson.NewObjectId(),
		PipeID:      pipe.ID,
		Title:       pipe.Title,
		Description: pipe.Description,
		Created:     time.Now(),
	}

	work := model.Work{
		WorkNo:      0,
		Title:       build.ID.Hex() + "-0",
		Description: pipe.Title + pipe.Description,
		Started:     time.Now(),
	}

	for index, group := range pipe.TaskGroups {
		templateData := &TemplateData{}
		templateData.Meta.PipeID = build.PipeID.Hex()
		templateData.Meta.BuildID = build.ID.Hex()
		templateData.Meta.TaskID = strconv.Itoa(index)
		templateData.Meta.Namespace = model.DefaultNameSpace
		templateData.Meta.Type = model.PodTypeJob
		templateData.TaskGroup = &group
		k8sjob, err := TaskGroup2Job(templateData, JobTemplateDefault)
		if err != nil {
			return nil, err
		}
		work.Steps = append(work.Steps, &model.WorkStep{
			StepNo: index,
			Title:  group.Title,
			Config: &group,
			K8sjob: k8sjob,
		})
	}
	// single work for now
	build.Works = append(build.Works, work)
	return build, nil
}

type TemplateMeta struct {
	PipeID    string
	BuildID   string
	TaskID    string
	Namespace string
	Type      string
}

type TemplateData struct {
	Meta      TemplateMeta
	TaskGroup *model.TaskGroup
}

var templateFuncs = template.FuncMap{
	"str2title": str2title,
}

func str2title(in string) (out string) {
	maxsize := 10
	for index, word := range in {
		if index > maxsize {
			break
		}
		if unicode.IsLetter(word) || unicode.IsDigit(word) {
			out += string(word)
		} else {
			out += "-"
		}
	}
	return
}

func ExecTemplate(tmplstr string, writer io.Writer, vars interface{}) error {
	tmpl, err := template.New("").Funcs(templateFuncs).Parse(tmplstr)
	if err != nil {
		return err
	}
	return tmpl.Execute(writer, vars)
}

func TaskGroup2Job(taskgroup *TemplateData, tmpl string) (*batch_v1.Job, error) {
	buffer := new(bytes.Buffer)
	if err := ExecTemplate(tmpl, buffer, taskgroup); err != nil {
		return nil, err
	}
	decoder := yaml.NewYAMLOrJSONDecoder(buffer, 4096)
	var job batch_v1.Job
	if err := decoder.Decode(&job); err != nil {
		return nil, err
	}
	return &job, nil
}
