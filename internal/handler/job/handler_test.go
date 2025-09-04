package job

import (
	"bytelyon-functions/internal/app"
	"bytelyon-functions/internal/client/s3"
	"bytelyon-functions/internal/model"
	"bytelyon-functions/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Save(t *testing.T) {

	job := model.Job{
		Type: model.NewsJobType,
		Frequency: model.Frequency{
			Unit:  "h",
			Value: 12,
		},
		Name:     "Test Job Name",
		Desc:     "Test Job Description",
		Keywords: []string{"Ford", "Bronco"},
	}
	ctx := context.Background()
	user := model.User{ID: app.NewUlid()}
	token := model.CreateJWTString(ctx, user)
	req := test.NewRequest(t).Header("authorization", "Bearer "+token).Post(job)

	res, err := Handler(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)

	_ = s3.NewWithContext(ctx).Delete(user.Path())
}
