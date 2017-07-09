package model

import (
	"time"

	"labix.org/v2/mgo/bson"
)

// save into database as json
type Plugin struct {
	ID          bson.ObjectId `bson:"_id"`
	Type        string        `bson:"type" index:"+name,unique"`
	Name        string        `bson:"name"`
	Version     string        `bson:"version"`
	Identifier  string        `bson:"indentifier" index:"unique"`
	Description string        `bson:"description,omitempty"`
	Options     []Option      `bson:"options,omitempty"`
	Updated     time.Time     `bson:"updated,omitempty"`
}

type Option struct {
	Must        bool     `bson:"must,omitempty"`
	Help        string   `bson:"help,omitempty"`
	Key         string   `bson:"key,omitempty"`
	Default     string   `bson:"default,omitempty"`
	Choose      []string `bson:"choose,omitempty"`
	MultiChoose bool     `bson:"multichoose,omitempty"`
}

func (p *Plugin) Valid() error {
	return nil
}

func (o *Option) Valid() error {
	return nil
}
