package main

type ItemRepository interface {
	FindAll() ([]Item, error)
	FindById() (Item, error)
	Insert() error
	Update() error
	Delete() error
}

func FindAll(g ItemRepository) ([]Item, error) {
	return g.FindAll()
}
func FindById(g ItemRepository) (Item, error) {
	return g.FindById()
}
func Insert(g ItemRepository) error {
	return g.Insert()
}
func Update(g ItemRepository) error {
	return g.Update()
}
func Delete(g ItemRepository) error {
	return g.Delete()
}

type UserRepository interface {
	FindAllUser() ([]User, error)
	Login() (User, error)
	Register() bool
	UpdateUser() error
	DeleteUser() error
	CheckUser() bool
}

func FindAllUser(g UserRepository) ([]User, error) {
	return g.FindAllUser()
}
func Login(g UserRepository) (User, error) {
	return g.Login()
}
func Register(g UserRepository) bool {
	return g.Register()
}
func UpdateUser(g UserRepository) error {
	return g.UpdateUser()
}
func DeleteUser(g UserRepository) error {
	return g.DeleteUser()
}

func CheckUser(g UserRepository) bool {
	return g.CheckUser()
}
