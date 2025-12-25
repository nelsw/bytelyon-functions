package internal

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var userKeyRegex = regexp.MustCompile(`.*user/([A-Za-z0-9]{26}/_.json)$`)

const userKey = "user"

type User struct {
	ID ulid.ULID `json:"id"`
}

func (u *User) MarshalZerologObject(evt *zerolog.Event) {
	evt.Stringer("user", u.ID)
}

func (u *User) String() string {
	return Path("user", u.ID)
}

func (u *User) FindJobs() []*Job {

	log.Trace().Msgf("find jobs")

	j := &Job{UserID: u.ID}
	keys, err := S3.Keys(j.String())
	if err != nil {
		log.Panic().Err(err).EmbedObject(u).Msg("find jobs failed")
	}

	var jobs []*Job
	for _, k := range keys {
		if k[len(k)-1] == '/' {
			continue
		}
		b, _ := S3.Get(k)
		var job Job
		_ = json.Unmarshal(b, &job)
		job.UserID = u.ID
		jobs = append(jobs, &job)
	}

	log.Trace().EmbedObject(u).Msgf("jobs found: %d", len(jobs))

	return jobs
}

func Users() []*User {

	log.Trace().Msg("getting users")

	keys, err := S3.Keys(`user/`)
	if err != nil {
		log.Panic().Err(err).Msg("failed to get users")
	}

	var users []*User
	for _, k := range keys {
		if !userKeyRegex.MatchString(k) {
			continue
		}

		b, _ := S3.Get(k)
		var user User
		err = json.Unmarshal(b, &user)
		users = append(users, &user)
	}

	log.Trace().Msgf("found %d users", len(users))

	return users
}

func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}
