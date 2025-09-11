package user

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/handler/jwt"
	"bytelyon-functions/internal/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"maps"
	"slices"
	"strconv"
	"strings"
)

func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	app.LogURLRequest(req)

	if app.IsOptions(req) {
		return app.OK()
	}

	user, err := jwt.Validate(ctx, req.Headers["authorization"])
	if err != nil {
		return app.Unauthorized(err)
	}

	if app.IsGet(req) {

		db := s3.New(ctx)

		after := req.QueryStringParameters["after"]
		delim := req.QueryStringParameters["delimiter"]
		size := req.QueryStringParameters["size"]
		i, _ := strconv.Atoi(size)

		var page model.Page
		if page, err = FindAll(db, user.ID, after, delim, i); err != nil {
			return app.BadRequest(err)
		}
		return app.Marshall(page)
	}

	return app.Marshall(req)
}

func FindAll(db s3.Client, userID ulid.ULID, after, delim string, size int) (page model.Page, err error) {

	log.Trace().
		Any("userID", userID).
		Str("after", after).
		Str("delimiter", delim).
		Int("size", size).
		Msg("FindAll")

	if size == 0 || size > 1000 {
		size = 10
	}

	m := map[string]any{}
	lastKey := after
	for {

		var keys []string
		if keys, err = db.Keys(model.UserPath, lastKey, "", 1000); err != nil {
			log.Err(err).Str("delimiter", delim).Msg("FindAll")
			return
		}

		for _, k := range keys {

			if delim == model.UserPath {
				if strings.HasPrefix(k, "user/") && strings.HasSuffix(k, "/_.json") {

					k = strings.TrimPrefix(k, "user/")
					k = strings.TrimSuffix(k, "/_.json")

					if len(k) > 26 {
						continue
					}

					if id, e := ulid.ParseStrict(k); e != nil {
						err = errors.Join(err, e)
					} else if _, ok := m[k]; !ok {
						m[k] = model.User{ID: id}
					}
				}
				continue
			}

			if !strings.Contains(k, userID.String()) || !strings.Contains(k, delim) {
				continue
			}

			if o, e := db.Get(k); e != nil {
				err = errors.Join(err, e)
			} else if _, ok := m[k]; !ok {
				var data map[string]any
				if e = json.Unmarshal(o, &data); e != nil {
					err = errors.Join(err, e)
				} else {
					m[k] = data
				}
			}
		}

		if len(keys) == 1000 {
			lastKey = keys[999]
			continue
		}

		break
	}

	keys := slices.Sorted(maps.Keys(m))
	page.Total = len(keys)
	afterFound := after == ""
	for _, k := range keys {
		if k == after {
			afterFound = true
		}
		if !afterFound {
			continue
		}
		if page = page.Add(m[k]); page.Size >= size {
			break
		}
	}

	log.Err(err).
		Int("total", page.Total).
		Int("size", page.Size).
		Msg("FindAll")

	return page, nil
}
