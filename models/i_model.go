package models

type Model interface {
	List() []interface{}

	GetById(id string) interface{}

	Delete(id string)
}

func CreateDemo() {
	f := CreateUserForm{
		Email:    "test@example.com",
		Password: "test",
	}
	f.Save()
}
