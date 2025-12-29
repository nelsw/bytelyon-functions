package db

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

func Save(e Entity) error {

	b, err := json.Marshal(e)
	if err != nil {
		log.Err(err).Str("key", e.Key()).Msg("EM - Failed to marshal")
		return err
	}

	if err = NewS3().Put(e.Key(), b); err != nil {
		log.Err(err).Str("key", e.Key()).Msg("EM - Failed to save")
	}

	return nil
}

func Find(e Entity) error {
	b, err := NewS3().Get(e.Key())
	if err == nil {
		err = json.Unmarshal(b, e)
	}
	return err
}

func Delete(e Entity) error {
	err := NewS3().Delete(e.Key())
	if err != nil {
		log.Err(err).Str("key", e.Key()).Msg("EM - Failed to delete")
	}
	return err
}

func Keys(DB S3, e Entity) ([]string, error) {

	var keys []string
	var after string

	regex := regexp.MustCompile(fmt.Sprintf(`.*%s/([A-Za-z0-9]{26}.json)$`, e))
	for {

		arr, err := DB.Keys(e.Dir(), after)
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

	DB := NewS3()

	var entities []E

	keys, err := Keys(DB, e)
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
			if b, err = DB.Get(k); err == nil {
				strs = append(strs, string(b))
			}
		})
	}

	wg.Wait()

	return entities, json.Unmarshal([]byte("["+strings.Join(strs, ",")+"]"), &entities)
}
