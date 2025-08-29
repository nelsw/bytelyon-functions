package user

import (
	"bytelyon-functions/internal/entity"
	"bytelyon-functions/internal/model/id"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {

	u := User{
		ID:    id.NewULID(),
		Email: gofakeit.Email(),
	}
	e := Email{
		ID:     u.Email,
		UserID: u.ID,
	}
	plainTextPassword := gofakeit.Password(true, true, true, true, true, 8)
	b := []byte(plainTextPassword)
	v, _ := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	p := Password{
		ID:    u.ID,
		Value: v,
	}
	_ = entity.New().Value(&u).Save()
	_ = entity.New().Value(&e).Save()
	_ = entity.New().Value(&p).Save()

	var user User
	_ = entity.New().Value(&user).ID(u.ID.String()).Find()
	password := Password{ID: user.ID}
	_ = entity.New().Value(&password).ID(u.ID.String()).Find()
	err := password.Validate(plainTextPassword)
	if err != nil {
		t.Error(err)
	}
}
