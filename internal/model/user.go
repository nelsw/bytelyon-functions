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
	var after string
	for {
		keys, e := db.Keys(UserPath, after, 1000)
		if e != nil {
			err = errors.Join(err, e)
			continue
		}
		for _, key := range keys {
			users = append(users, User{ID: ulid.MustParse(key)})
		}
		if len(keys) == 1000 {
			after = keys[999]
			continue
		}
		break
	}
	log.Err(err).Int("users", len(users)).Msg("find all users")
	return
}

func (u User) FindAllJobs(db s3.Client) (jobs Jobs, err error) {
	var after string
	prefix := Job{UserID: u.ID}.Path()
	for {
		keys, e := db.Keys(prefix, after, 1000) // todo - define job limit
		if e != nil {
			err = errors.Join(err, e)
			continue
		}
		fmt.Println(keys)
		for _, key := range keys {
			key = strings.TrimSuffix(key, "/_.json")
			key = key[strings.LastIndex(key, "/")+1:]
			jobs = append(jobs, Job{
				ID:     ulid.MustParse(key),
				UserID: u.ID,
			})
		}
		if len(keys) == 1000 {
			after = keys[999]
			continue
		}
		break
	}
	log.Err(err).Int("jobs", len(jobs)).Msg("find all jobs")
	return
}

func UserKey(ID ulid.ULID) string {
	return fmt.Sprintf("%s/%s", UserPath, ID)
}
