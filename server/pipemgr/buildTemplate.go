package pipemgr

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strconv"
	"unicode"

	"gopkg.in/mgo.v2/bson"
	"k8s.io/apimachinery/pkg/util/yaml"
	v1 "k8s.io/client-go/pkg/api/v1"
	batch_v1 "k8s.io/client-go/pkg/apis/batch/v1"

	"github.com/arlert/malcolm/model"
)

type buildtemplate struct {
}

func (t *buildtemplate) ConfigToBuild(pipe *model.Pipeline) (*model.Build, error) {
	build := &model.Build{
		ID:          bson.NewObjectId(),
		PipeID:      pipe.ID,
		Title:       pipe.Title,
		Description: pipe.Description,
		Project:     model.DefaultNameSpace,
	}

	templateData := &templateData{}
	templateData.Meta.PipeID = build.PipeID.Hex()
	templateData.Meta.BuildID = build.ID.Hex()
	templateData.Meta.Namespace = model.DefaultNameSpace
	templateData.Pipe = pipe

	if pipe.StorageClass != "" {
		volumn, err := pipe2Volumn(templateData, VolumnTemplateDefault)
		if err != nil {
			return nil, err
		}
		build.Volumn = volumn
	}

	for index, group := range pipe.TaskGroups {
		templateData.Meta.TaskID = strconv.Itoa(index)
		templateData.Meta.Type = model.PodTypeJob
		templateData.TaskGroup = &group
		k8sjob, err := taskGroup2Job(templateData, JobTemplateDefault)
		if err != nil {
			return nil, err
		}
		build.Steps = append(build.Steps, &model.WorkStep{
			StepNo: index,
			Title:  group.Title,
			Config: &group,
			K8sjob: k8sjob,
		})
	}
	return build, nil
}

type templateMeta struct {
	PipeID    string
	BuildID   string
	TaskID    string
	Namespace string
	Type      string
}

type templateData struct {
	Meta      templateMeta
	TaskGroup *model.TaskGroup
	Pipe      *model.Pipeline
}

var templateFuncs = template.FuncMap{
	"str2title":     str2title,
	"interface2str": interface2str,
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

func interface2str(in interface{}) (out string) {
	if out, ok := in.(string); ok {
		return out
	}
	buf, err := json.Marshal(in)
	if err != nil {
		return "error"
	}
	out = string(buf)
	return
}

func taskGroup2Job(data *templateData, tmpl string) (*batch_v1.Job, error) {
	var job batch_v1.Job
	if err := ExecTemplate(tmpl, data, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func pipe2Volumn(data *templateData, tmpl string) (*v1.PersistentVolumeClaim, error) {
	var volumn v1.PersistentVolumeClaim
	if err := ExecTemplate(tmpl, data, &volumn); err != nil {
		return nil, err
	}
	return &volumn, nil
}

// ExecTemplate exec template with vars out to kubernete data model
func ExecTemplate(tmplstr string, vars interface{}, out interface{}) (err error) {
	buffer := new(bytes.Buffer)
	tmpl, err := template.New("").Funcs(templateFuncs).Parse(tmplstr)
	if err != nil {
		return
	}
	err = tmpl.Execute(buffer, vars)
	if err != nil {
		return
	}
	decoder := yaml.NewYAMLOrJSONDecoder(buffer, 4096)
	err = decoder.Decode(out)
	return
}
