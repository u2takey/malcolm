package store

import (
	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"

	mgoutil "github.com/arlert/malcolm/utils/mongo"
)

type Store struct {
	Cfg     *mgoutil.Config
	Cols    Collections
	session *mgo.Session
}

type Collections struct {
	Build  mgoutil.Collection `coll:"build"`
	Log    mgoutil.Collection `coll:"log"`
	Plugin mgoutil.Collection `coll:"plugin"`
	Pipe   mgoutil.Collection `coll:"pipe"`
}

func New(cfg *mgoutil.Config) (s *Store) {
	s = &Store{Cfg: cfg}
	session, err := mgoutil.Open(&s.Cols, s.Cfg)
	if err != nil {
		logrus.Fatalln("mgoutil.Open error -", err)
		return
	}

	err = session.Ping()
	if err != nil {
		logrus.Fatalln("session.Ping -", err)
		return
	}

	s.session = session

	s.Cols.Build.EnsureIndexes()
	s.Cols.Log.EnsureIndexes()
	s.Cols.Plugin.EnsureIndexes()
	s.Cols.Pipe.EnsureIndexes()
	return
}
