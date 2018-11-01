package main

//base repository for all end points
type BaseRepository interface {
	FindAll() ([]byte, error)
	FindById() ([]byte, error)
	Insert() error
	Update() error
	Delete() error
}

//user private repository
type UserRepository interface {
	Login() ([]byte, error)
	Register() error
	CheckUser() bool
}
