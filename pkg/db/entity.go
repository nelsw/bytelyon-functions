package db

type Entity interface {
	Dir() string
	Key() string
}
