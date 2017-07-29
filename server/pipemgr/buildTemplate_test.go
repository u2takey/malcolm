package pipemgr

import (
	"encoding/json"
	"testing"

	"github.com/arlert/malcolm/model"
	"gopkg.in/mgo.v2/bson"
)

func TestTemplate(t *testing.T) {
	testtaskgroup := []model.TaskGroup{model.TaskGroup{
		Title: "taskgroup",
		Tasks: []model.Task{model.Task{
			Title:       "tasktitle",
			Plugin:      "taksplugin",
			Environment: map[string]string{"a": "a", "b": "b"},
			Command:     []string{"1", "2", "3"},
			Args:        []string{"4", "5", "6"},
			PullPolicy:  "pullpolicy",
			Credentials: map[string]string{"x": "x", "y": "y"},
		}},
		PreTasks: []model.Task{model.Task{
			Title:       "tasktitle2",
			Plugin:      "taksplugin2",
			Environment: map[string]string{"a2": "a2", "b2": "b2"},
			Command:     []string{"1", "2", "3"},
			Args:        []string{"4", "5", "6"},
			PullPolicy:  "pullpolicy2",
			Credentials: map[string]string{"x2": "x2", "y2": "y2"},
		}},
	}}
	pipe := &model.Pipeline{
		ID:           bson.NewObjectId(),
		Title:        "a pipe",
		Description:  "a pipe",
		TaskGroups:   testtaskgroup,
		StorageClass: "ceph",
		StorageSize:  "1Gi",
		WorkSpace:    "/workspace",
	}
	tmpl := &buildtemplate{}
	job, err := tmpl.ConfigToBuild(pipe)
	b, _ := json.MarshalIndent(job, "", "  ")
	t.Log(string(b))
	if err != nil {
		t.Fatal(err)
	}
}
