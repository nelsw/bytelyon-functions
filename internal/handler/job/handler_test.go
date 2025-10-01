package job

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/handler/jwt"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func fakeJob() model.Job {
	return model.Job{
		Type:      model.NewsJobType,
		Frequency: model.Frequency{Unit: "h", Value: 12},
		Name:      gofakeit.Name(),
		Desc:      gofakeit.Sentence(10),
		Keywords:  []string{"Rivian"},
	}
}

func Test_Handler_Post(t *testing.T) {
	user := test.DemoUser()
	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, user)).
		Post(fakeJob())

	res, _ := Handler(test.CTX, req)

	assert.Equal(t, res.StatusCode, 200)
}

func Test_Handler_Get(t *testing.T) {
	test.Init(t)
	user := test.DemoUser()

	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, user)).
		Query("size", 2).
		Get()

	res, _ := Handler(test.CTX, req)

	assert.Equal(t, res.StatusCode, 200)
}

func Test_Handler_Delete(t *testing.T) {
	test.Init(t)
	user := test.DemoUser()
	job := fakeJob()
	job.ID = app.NewUlid()
	_, _ = Save(s3.New(test.CTX), user.ID, app.MustMarshal(job))

	req := test.
		NewRequest(t).
		Bearer(jwt.CreateString(test.CTX, user)).
		Query("id", job.ID).
		Delete()

	res, _ := Handler(test.CTX, req)

	assert.Equal(t, res.StatusCode, 200)
}
