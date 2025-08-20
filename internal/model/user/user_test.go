package user

import (
	"bytelyon-functions/internal/entity"
	"testing"
)

func TestNewUser(t *testing.T) {
	var user User
	_ = entity.New().Value(&user).ID("01K30NG3ZBNEFE5E549K4SRQWJ").Find()
	password := UserPassword{ID: user.ID}
	_ = entity.New().Value(&password).Find()
	err := password.Validate("Farts1234!")
	if err != nil {
		t.Error(err)
	}
}
