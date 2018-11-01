package main

//BaseRepository repository for all end points
type BaseRepository interface {
	FindAll() ([]byte, error)
	FindById() ([]byte, error)
	Insert() error
	Update() error
	Delete() error
}

//UserRepository private repository
type UserRepository interface {
	Login() ([]byte, error)
	Register() error
	CheckUser() bool
}
