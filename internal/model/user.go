package model

import (
	"bytelyon-functions/internal/client/s3"
	"errors"
	"fmt"
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

func FindAllUsers(db s3.Client) (users Users, err error) {
	var keys []string
	var after string
	for {
		kk, e := db.Keys(UserPath, after, 1000)
		if e != nil {
			err = errors.Join(err, e)
			continue
		}
		for _, k := range kk {
			keys = append(keys, k)
		}
		if len(keys) == 1000 {
			after = keys[999]
			continue
		}
		break
	}

	m := map[string]ulid.ULID{}
	for _, k := range keys {
		v := strings.Split(k, "/")[1]
		ID, e := ulid.Parse(v)
		if e != nil {
			continue
		}
		m[v] = ID
	}

	for _, id := range m {
		users = append(users, User{ID: id})
	}

	log.Err(err).Int("users", len(users)).Msg("find all users")
	return
}

func (u User) FindAllJobs(db s3.Client) (jobs Jobs, err error) {
	var after string
	prefix := Job{UserID: u.ID}.Path()
	var paths []string
	for {
		keys, e := db.Keys(prefix, after, 1000) // todo - define job limit
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
