package em

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var db = s3.New()

type Entity interface {
	// Path defines the location where all entities of a given type can be found
	Path() string
	// Dir defines the location to a single entity its children
	Dir() string
	// Key defines the location of a single entity only
	Key() string
}

func Save(e Entity) error {
	if b, err := json.Marshal(&e); err != nil {
		return err
	} else if err = db.Put(e.Key(), b); err != nil {
		return err
	}
	return nil
}

func Find(e Entity) error {
	if b, err := db.Get(e.Key()); err != nil {
		return err
	} else if err = json.Unmarshal(b, e); err != nil {
		return err
	}
	return nil
}

func Delete(e Entity) error {
	return db.Delete(e.Key())
}

func Keys(e Entity, r ...*regexp.Regexp) ([]string, error) {

	var keys []string
	var after string

	for {

		arr, err := db.Keys(e.Path(), after, 1000)
		if err != nil {
			return nil, err
		}

		if len(r) == 0 {
			keys = append(keys, arr...)
		} else {
			for _, k := range arr {
				if r[0].MatchString(k) {
					keys = append(keys, k)
				}
			}
		}

		if len(arr) < 1000 {
			break
		}

		after = arr[len(arr)-1]
	}

	log.Info().Int("count", len(keys)).Msg("Keys found")

	return keys, nil
}

func FindAll[T Entity](e T, r ...*regexp.Regexp) ([]T, error) {

	keys, err := Keys(e, r...)
	if err != nil {
		return nil, err
	}

	var strs []string
	var wg sync.WaitGroup

	for i, k := range keys {
		if i%1_000 == 0 {
			time.Sleep(1 * time.Second)
		}
		wg.Go(func() {
			var b []byte
			if b, err = db.Get(k); err != nil {
				log.Warn().Err(err).Str("key", k).Msg("failed to get entity")
			} else {
				strs = append(strs, string(b))
			}
		})
	}

	log.Info().Int("count", len(keys)).Msg("waiting to get all entities")

	wg.Wait()

	log.Info().Int("count", len(strs)).Msg("all entities retrieved")

	var entities []T
	if err = json.Unmarshal([]byte("["+strings.Join(strs, ",")+"]"), &entities); err != nil {
		return nil, err
	}

	return entities, nil
}

func FindLast(e Entity) error {
	var key string
	for {

		arr, err := db.Keys(e.Path(), key, 1000)
		if err != nil {
			return err
		}

		if key = arr[len(arr)-1]; len(arr) < 1000 {
			break
		}
	}

	b, _ := db.Get(key)
	return json.Unmarshal(b, e)
}
