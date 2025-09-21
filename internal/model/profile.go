package model

import "fmt"

type Profile struct {
	User   User   `json:"-"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Bio    string `json:"bio"`
}

func (p Profile) Path() string {
	return fmt.Sprintf("%s/profile/%s", p.User.Key(), p.User.ID)
}
