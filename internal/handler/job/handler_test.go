package job

import (
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
	req := test.NewRequest(t).Header("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImlkIjoiMDFLNDhQQzBCSzEzQldWMkNHV0ZQOFFRSDAifSwiaXNzIjoiQnl0ZUx5b24iLCJleHAiOjE3NTcwMDA4NDYsIm5iZiI6MTc1Njk5OTA0NiwiaWF0IjoxNzU2OTk5MDQ2LCJqdGkiOiIwZmNmZWZhNi0zMTE3LTRjNWItOTlhMy04OGM0ODhhNTc1N2UifQ.ZGk73d0aWuC08l-zeDzclVEmcEyYcYiWfHO57dZ1swY").Post(job)

	res, err := Handler(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
}
