package models

type Model interface {
	List() []interface{}

	GetById(id string) interface{}

	Delete(id string)
}
