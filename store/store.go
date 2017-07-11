package store

import (
	"github.com/Sirupsen/logrus"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"

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
	Job    mgoutil.Collection `coll:"job"`
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
	s.Cols.Job.EnsureIndexes()
	return
}
