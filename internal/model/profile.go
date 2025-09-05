package model

import "fmt"

type Profile struct {
	User  User   `json:"-"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (p Profile) Path() string {
	return fmt.Sprintf("%s/profile/%s", p.User.Key(), p.User.ID)
}
