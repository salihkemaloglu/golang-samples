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
