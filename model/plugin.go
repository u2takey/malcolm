package model

import (
	"errors"
	"time"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/mgo.v2/bson"
	"k8s.io/kubernetes/pkg/util/sets"
)

//Plugin type
var (
	PluginTypeScm     = "scm"
	PluginTypeBuild   = "build"
	PluginTypeArchive = "archive"
	PluginTypeNotify  = "notify"

	AvailablePluginTypes = sets.NewString(PluginTypeScm, PluginTypeBuild, PluginTypeArchive, PluginTypeNotify)
)

// Plugin for build process
type Plugin struct {
	ID          bson.ObjectId `bson:"_id"`
	Type        string        `bson:"type" index:"+name,unique"`
	Name        string        `bson:"name"  index:"unique"`
	Description string        `bson:"description,omitempty"`
	Version     string        `bson:"version"`
	Path        string        `bson:"url"`
	Schema      interface{}   `bson:"schema,omitempty"`      // for pluginConfig setting
	Credentials []string      `bson:"credentials,omitempty"` // credential needed
	Commands    bool          `bson:"commands,omitempty"`    // allow arbitrary commands
	Args        bool          `bson:"args,omitempty"`        // allow arbitrary args
	Updated     time.Time     `bson:"updated,omitempty"`
}

// Valid Valid plugin input
func (p *Plugin) Valid() (err error) {
	if !AvailablePluginTypes.HasAny(p.Type) {
		err = WarpErrors(err, errors.New("plugin type not valid"))
	}
	if p.Description == "" || p.Name == "" {
		err = WarpErrors(err, errors.New("Description and Name cannot be empty"))
	}
	if p.Version == "" {
		err = WarpErrors(err, errors.New("Version cannot be empty"))
	}
	if p.Path == "" {
		err = WarpErrors(err, errors.New("Path cannot be empty"))
	}
	schemaLoader := gojsonschema.NewGoLoader(p.Schema)
	_, err1 := schemaLoader.LoadJSON()
	err = WarpErrors(err, err1)
	return
}

// ValidDoc valid input document json with schema
func (p *Plugin) ValidDoc(document interface{}) (ret bool, errDetails []string) {
	schemaLoader := gojsonschema.NewGoLoader(p.Schema)
	documentLoader := gojsonschema.NewGoLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	ret = false
	if err != nil {
		errDetails = append(errDetails, err.Error())
		return
	}
	if result.Valid() {
		ret = true
		return
	}
	for _, err := range result.Errors() {
		errDetails = append(errDetails, err.String())
	}
	return
}
