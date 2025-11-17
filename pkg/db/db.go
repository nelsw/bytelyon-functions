package db

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
)

type Entity interface {
	Key() string
}

var db s3.Service

func init() {
	db = s3.New()
}

func Find[T Entity](e T) error {
	if b, err := db.Get(e.Key()); err != nil {
		return err
	} else if err = json.Unmarshal(b, &e); err != nil {
		return err
	}
	return nil
}

func Save[T Entity](e T) (b []byte, err error) {
	if b, err = json.MarshalIndent(&e, "", "\t"); err != nil {
		return nil, err
	} else if err = db.Put(e.Key(), b); err != nil {
		return nil, err
	}
	return b, nil
}
