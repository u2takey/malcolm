package mgoutil

import (
	"reflect"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func Dail(host, mode string, syncTimeoutInS int64, direct bool) (session *mgo.Session, err error) {
	addrs := getMongoHosts(host)
	timeout := time.Second * 10
	info := mgo.DialInfo{
		Addrs:   addrs,
		Direct:  direct,
		Timeout: timeout,
	}
	session, err = mgo.DialWithInfo(&info)
	if err != nil {
		return
	}
	session.SetSyncTimeout(1 * time.Minute)
	session.SetSocketTimeout(1 * time.Minute)

	if mode != "" {
		SetMode(session, mode, true)
	}
	if syncTimeoutInS != 0 {
		session.SetSyncTimeout(time.Duration(int64(time.Second) * syncTimeoutInS))
	}
	return
}

// ------------------------------------------------------------------------

type Safe struct {
	W        int    `json:"w"`
	WMode    string `json:"wmode"`
	WTimeout int    `json:"wtimeoutms"`
	FSync    bool   `json:"fsync"`
	J        bool   `json:"j"`
}

type Config struct {
	Host           string `json:"host"`
	DB             string `json:"db"`
	Mode           string `json:"mode"`
	SyncTimeoutInS int64  `json:"timeout"` // 以秒为单位
	Direct         bool   `json:"direct"`

	Safe *Safe `json:"safe"`
}

func Open(ret interface{}, cfg *Config) (session *mgo.Session, err error) {

	session, err = Dail(cfg.Host, cfg.Mode, cfg.SyncTimeoutInS, cfg.Direct)
	if err != nil {
		return
	}

	EnsureSafe(session, cfg.Safe)

	if ret != nil {
		db := session.DB(cfg.DB)
		err = InitCollections(ret, db)
		if err != nil {
			session.Close()
			session = nil
		}
	}
	return
}

func InitCollections(ret interface{}, db *mgo.Database) (err error) {

	v := reflect.ValueOf(ret)
	if v.Kind() != reflect.Ptr {
		log.Error("InitCollections: ret must be a pointer")
		return syscall.EINVAL
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		log.Error("InitCollections: ret must be a struct pointer")
		return syscall.EINVAL
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.Tag == "" {
			continue
		}
		coll := sf.Tag.Get("coll")
		if coll == "" {
			continue
		}
		switch elem := v.Field(i).Addr().Interface().(type) {
		case *Collection:
			elem.Collection = db.C(coll)
		case **mgo.Collection:
			*elem = db.C(coll)
		default:
			log.Error("InitCollections: coll must be *mgo.Collection or mongo.Collection")
			return syscall.EINVAL
		}
	}
	return
}

// ------------------------------------------------------------------------

// W 和 WMode 只在 replset 模式下生效，非replset不能配置，否则会出错
// WMode只在2.0版本以上才生效
func EnsureSafe(session *mgo.Session, safe *Safe) {
	if safe == nil {
		return
	}
	session.EnsureSafe(&mgo.Safe{
		W:        safe.W,
		WMode:    safe.WMode,
		WTimeout: safe.WTimeout,
		FSync:    safe.FSync,
		J:        safe.J,
	})
}

// ------------------------------------------------------------------------
// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
func getMongoHosts(raw string) []string {
	if strings.HasPrefix(raw, "mongodb://") {
		raw = raw[len("mongodb://"):]
	}
	if idx := strings.Index(raw, "@"); idx != -1 {
		raw = raw[idx+1:]
	}
	if idx := strings.Index(raw, "/"); idx != -1 {
		raw = raw[:idx]
	}
	if idx := strings.Index(raw, "?"); idx != -1 {
		raw = raw[:idx]
	}
	return strings.Split(raw, ",")
}

var g_modes = map[string]int{
	"eventual":  0,
	"monotonic": 1,
	"mono":      1,
	"strong":    2,
}

func SetMode(s *mgo.Session, modeFriendly string, refresh bool) {

	mode, ok := g_modes[strings.ToLower(modeFriendly)]
	if !ok {
		log.Fatal("invalid mgo mode")
	}
	switch mode {
	case 0:
		s.SetMode(mgo.Eventual, refresh)
	case 1:
		s.SetMode(mgo.Monotonic, refresh)
	case 2:
		s.SetMode(mgo.Strong, refresh)
	}
}
