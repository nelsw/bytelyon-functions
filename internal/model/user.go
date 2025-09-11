package model

import (
	"bytelyon-functions/internal/client/s3"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

const UserPath = "user"

type Users []User
type User struct {
	ID ulid.ULID `json:"id"`
}

func (u User) Key() string {
	return fmt.Sprintf("%s/%s", UserPath, u.ID)
}

func FindAllUsers(db s3.Client) (Users, error) {

	m := map[string]User{}
	var after string
	var err error
	for {

		var keys []string
		if keys, err = db.Keys(UserPath, after, "", 1000); err != nil {
			log.Err(err).Msg("FindAllUsers")
			return nil, err
		}

		for _, k := range keys {

			if strings.HasPrefix(k, "user/") && strings.HasSuffix(k, "/_.json") {

				k = strings.TrimPrefix(k, "user/")
				k = strings.TrimSuffix(k, "/_.json")

				if len(k) > 26 {
					continue
				}

				if id, e := ulid.ParseStrict(k); e != nil {
					err = errors.Join(err, e)
				} else if _, ok := m[k]; !ok {
					m[k] = User{ID: id}
				}
			}
		}

		if len(keys) == 1000 {
			after = keys[999]
			continue
		}

		break
	}

	users := slices.Collect(maps.Values(m))

	log.Err(err).Int("count", len(users)).Msg("FindAllUsers")

	return users, err
}

func (u User) FindAllJobs(db s3.Client) (jobs Jobs, err error) {
	var after string
	prefix := Job{UserID: u.ID}.Path()
	var paths []string
	for {
		keys, e := db.Keys(prefix, after, "", 1000) // todo - define job limit
		if e != nil {
			err = errors.Join(err, e)
			continue
		}
		fmt.Println(keys)
		for _, key := range keys {
			paths = append(paths, key)
		}
		if len(keys) == 1000 {
			after = keys[999]
			continue
		}
		break
	}

	m := map[string]ulid.ULID{}
	for _, path := range paths {
		v := strings.Split(path, "/")[3]
		id, e := ulid.Parse(v)
		if e != nil {
			continue
		}
		m[v] = id
	}

	for _, id := range m {
		jobs = append(jobs, Job{ID: id})
	}

	log.Err(err).Int("jobs", len(jobs)).Msg("find all jobs")
	return
}

func UserKey(ID ulid.ULID) string {
	return fmt.Sprintf("%s/%s", UserPath, ID)
}
