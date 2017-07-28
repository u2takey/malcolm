package mgoutil

import (
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// ------------------------------------------------------------------------

type Collection struct {
	*mgo.Collection
}

// ensure indexes for a collection
//
// eg. c.EnsureIndexes(
//			"uid :unique", "email :unique",
//			"serial_num", "uid,status,delete :sparse,background")
//
func (c Collection) EnsureIndexes(indexes ...string) {

	for _, colIndex := range indexes {
		var index mgo.Index
		pos := strings.Index(colIndex, ":")
		if pos >= 0 {
			parseIndexOptions(&index, colIndex[pos+1:])
			colIndex = colIndex[:pos]
		}
		index.Key = strings.Split(strings.TrimRight(colIndex, " "), ",")
		err := c.EnsureIndex(index)
		if err != nil {
			log.Fatal("<Mongo.C>:", c.Name, "Index:", index.Key, " error:", err)
			break
		}
	}
}

func parseIndexOptions(index *mgo.Index, options string) {

	for {
		var option string
		pos := strings.Index(options, ",")
		if pos < 0 {
			option = options
		} else {
			option = options[:pos]
			options = options[pos+1:]
		}
		switch option {
		case "unique":
			index.Unique = true
		case "sparse":
			index.Sparse = true
		case "background":
			index.Background = true
		default:
			log.Fatal("Unknown option:", option)
		}
		if pos < 0 {
			return
		}
	}
}

func getNameOf(tag string) string {

	pos := strings.Index(tag, ",")
	if pos < 0 {
		return tag
	}
	return tag[:pos]
}

// ensure indexes for a collection
//
// eg.
//	  v := new(struct{
//		  Id        string `bson:"_id"`
//		  Pid       string `bson:"pid" index:"+name+time,unique"`
//		  Name      string `bson:"name"`
//		  Email     string `bson:"email" index:"unique"`
//		  SerialNum int    `bson:"serial_num" index:"index"`
//	  })
//	  c.EnsureIndexesByType(reflect.TypeOf(v))
//
func (c Collection) EnsureIndexesByType(t reflect.Type) {

	n := t.NumField()
	for i := 0; i < n; i++ {
		tag := t.Field(i).Tag
		if options := tag.Get("index"); options != "" {
			index := mgo.Index{Key: []string{getNameOf(tag.Get("bson"))}}
			parseIndexOptionsByType(&index, options)
			err := c.EnsureIndex(index)
			if err != nil {
				log.Fatal("<Mongo.C>:", c.Name, "Index:", index.Key, " error:", err)
				break
			}
		}
	}
}

func parseIndexOptionsByType(index *mgo.Index, options string) {

	for {
		var option string
		pos := strings.Index(options, ",")
		if pos < 0 {
			option = options
		} else {
			option = options[:pos]
			options = options[pos+1:]
		}
		switch option {
		case "index":
		case "unique":
			index.Unique = true
		case "sparse":
			index.Sparse = true
		case "background":
			index.Background = true
		default:
			if !strings.HasPrefix(option, "+") {
				log.Fatal("Unknown option:", option)
			}
			index.Key = append(index.Key, strings.Split(option[1:], "+")...)
		}
		if pos < 0 {
			return
		}
	}
}

func (c Collection) CopySession() Collection {

	db := c.Database
	return Collection{db.Session.Copy().DB(db.Name).C(c.Name)}
}

func (c Collection) CloseSession() (err error) {

	c.Database.Session.Close()
	return nil
}

// ------------------------------------------------------------------------
