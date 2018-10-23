package main

type BaseRepository interface {
	FindAll() ([]byte, error)
	FindById() ([]byte, error)
	Insert() error
	Update() error
	Delete() error
}
type UserRepository interface {
	Login() ([]byte, error)
	Register() error
	CheckUser() bool
}
