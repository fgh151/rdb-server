package models

type Model interface {
	List(limit int, offset int, sort string, order string) []interface{}

	GetById(id string) interface{}

	Delete(id string)

	Total() *int64
}

func CreateDemo() {
	f := CreateUserForm{
		Email:    "test@example.com",
		Password: "test",
	}
	f.Save()
}
