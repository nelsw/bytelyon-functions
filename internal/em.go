package db

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Entity interface {
	String() string
}

func Find(e Entity) error {
	b, err := New().Get(e.String() + `/_.json`)
	if err == nil {
		err = json.Unmarshal(b, e)
	}
	return err
}

func Save(e Entity) error {
	return New().Put(e.String()+"/_.json", e)
}

func Delete(e Entity) error {
	return New().Delete(e.String() + "/_.json")
}

func Keys(e Entity) ([]string, error) {

	var keys []string
	var after string

	regex := regexp.MustCompile(fmt.Sprintf(`.*%s/([A-Za-z0-9]{26}/_.json)$`, e))
	for {

		arr, err := New().Keys(e.String(), after, 1_000)
		if err != nil {
			return nil, err
		}

		for _, k := range arr {
			if regex.MatchString(e.String()) {
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

	keys, err := Keys(e)
	if err != nil {
		return nil, err
	}

	var entities []E
	if len(keys) == 0 {
		return entities, nil
	}

	var strs []string
	var wg sync.WaitGroup

	for i, k := range keys {
		if i > 0 && i%1_500 == 0 {
			time.Sleep(1 * time.Second)
		}
		wg.Go(func() {
			var b []byte
			if b, err = New().Get(k); err == nil {
				strs = append(strs, string(b))
			}
		})
	}

	wg.Wait()

	return entities, json.Unmarshal([]byte("["+strings.Join(strs, ",")+"]"), &entities)
}

func FindLast(e Entity) error {

	var key string
	for {
		arr, err := New().Keys(e.String(), key, 1000)
		if err != nil {
			return err
		}

		if key = arr[len(arr)-1]; len(arr) < 1000 {
			break
		}
	}

	b, err := New().Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, e)
}
