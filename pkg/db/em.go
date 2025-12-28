package db

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	db S3
)

func init() {
	godotenv.Load()
	db = NewS3()
}

type Entity interface {
	fmt.Stringer
}

func Save(e Entity) error {

	b, err := json.Marshal(e)
	if err != nil {
		log.Err(err).
			Stringer("key", e).
			Any("entity", e).
			Msg("EM - Failed to marshal")
		return err
	}

	if err = db.Put(e.String()+"/_.json", b); err != nil {
		log.Err(err).
			Str("key", e.String()).
			Bytes("body", b).
			Msg("EM - Failed to save")
	}

	return nil
}

func Find(e Entity) error {
	b, err := db.Get(e.String() + `/_.json`)
	if err == nil {
		err = json.Unmarshal(b, e)
	}
	return err
}

func Delete(e Entity) error {
	err := db.Delete(e.String() + "/_.json")
	if err != nil {
		log.Err(err).
			Stringer("key", e).
			Msg("EM - Failed to delete")
	}
	return err
}

func Keys(e Entity) ([]string, error) {

	var keys []string
	var after string

	regex := regexp.MustCompile(fmt.Sprintf(`.*%s/([A-Za-z0-9]{26}/_.json)$`, e))
	for {

		arr, err := db.Keys(e.String(), after, 1_000)
		if err != nil {
			return nil, err
		}

		for _, k := range arr {
			if regex.MatchString(k) {
				keys = append(keys, k)
			}
		}

		if len(arr) < 1_000 {
			break
		}

		after = arr[len(arr)-1]
	}

	return keys, nil
}

func List[E Entity](e E) ([]E, error) {

	var entities []E

	keys, err := Keys(e)
	if err != nil || len(keys) == 0 {
		return entities, err
	}

	var strs []string
	var wg sync.WaitGroup

	for i, k := range keys {
		if i > 0 && i%1_500 == 0 {
			time.Sleep(time.Second)
		}
		wg.Go(func() {
			var b []byte
			if b, err = db.Get(k); err == nil {
				strs = append(strs, string(b))
			}
		})
	}

	wg.Wait()

	return entities, json.Unmarshal([]byte("["+strings.Join(strs, ",")+"]"), &entities)
}
