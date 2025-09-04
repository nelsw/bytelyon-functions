package db

import (
	"bytelyon-functions/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveAndFindOneUser(t *testing.T) {

	exp := model.NewUser()
	err := Save(exp.Path(), exp)
	assert.NoError(t, err)

	var act model.User
	err = FindOne(exp.Path(), &act)
	assert.NoError(t, err)

	assert.Equal(t, *exp, act)

	email, _ := model.NewEmail(exp, "demo@demo.com")
	_ = Save(email.Path(), email)
	pork, _ := model.NewPassword(exp, "Demo123!")
	_ = Save(pork.Path(), pork)
}
