package mgo

import (
	"errors"
	"strings"
)

// a fork of Apply, fix https://bugs.launchpad.net/mgo/+bug/1232685
func (q *Query) Apply2(change Change, result interface{}) (info *ChangeInfo, err error) {
	q.m.Lock()
	session := q.session
	op := q.op // Copy.
	q.m.Unlock()

	c := strings.Index(op.collection, ".")
	if c < 0 {
		return nil, errors.New("bad collection name: " + op.collection)
	}

	dbname := op.collection[:c]
	cname := op.collection[c+1:]

	cmd := findModifyCmd{
		Collection: cname,
		Update:     change.Update,
		Upsert:     change.Upsert,
		Remove:     change.Remove,
		New:        change.ReturnNew,
		Query:      op.query,
		Sort:       op.options.OrderBy,
		Fields:     op.selector,
	}

	session = session.Clone()
	defer session.Close()
	session.SetMode(Strong, false)

	var doc valueResult
	err = session.DB(dbname).Run(&cmd, &doc)
	if err != nil {
		if qerr, ok := err.(*QueryError); ok && qerr.Message == "No matching object found" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if doc.Value.Kind == 0x0A {
		return nil, ErrNotFound
	}
	err = doc.Value.Unmarshal(result)
	if err != nil {
		return nil, err
	}
	info = &ChangeInfo{}
	lerr := &doc.LastError
	if lerr.UpdatedExisting {
		info.Updated = lerr.N
	} else if change.Remove {
		info.Removed = lerr.N
	} else if change.Upsert {
		info.UpsertedId = lerr.UpsertedId
	}
	return info, nil
}
