package job

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func fakeUser() model.User {
	return model.User{ID: app.NewUlid()}
}

func fakeJob() model.Job {
	return model.Job{
		Type:      model.NewsJobType,
		Frequency: model.Frequency{Unit: "h", Value: 12},
		Name:      gofakeit.Name(),
		Desc:      gofakeit.Sentence(10),
		Keywords:  []string{"GM", "GMC", "EV", "Hummer"},
	}
}

func Test_Handler_Post(t *testing.T) {

	req := test.
		NewRequest(t).
		Bearer(model.CreateJWTString(context.Background(), fakeUser())).
		Post(fakeJob())

	res, _ := Handler(context.Background(), req)

	assert.Equal(t, res.StatusCode, 200)
}

func Test_Handler_Get(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	user := fakeUser()
	_, _ = Save(s3.New(), user, string(app.MustMarshal(fakeJob())), false)

	req := test.
		NewRequest(t).
		Bearer(model.CreateJWTString(context.Background(), user)).
		Query("size", 2).
		Get()

	res, _ := Handler(context.Background(), req)

	assert.Equal(t, res.StatusCode, 200)
}

func Test_Handler_Delete(t *testing.T) {
	t.Setenv("APP_MODE", "test")
	user := fakeUser()
	job := fakeJob()
	job.ID = app.NewUlid()
	_, _ = Save(s3.New(), user, string(app.MustMarshal(job)), false)

	req := test.
		NewRequest(t).
		Bearer(model.CreateJWTString(context.Background(), user)).
		Query("id", job.ID).
		Delete()

	res, _ := Handler(context.Background(), req)

	assert.Equal(t, res.StatusCode, 200)
}
