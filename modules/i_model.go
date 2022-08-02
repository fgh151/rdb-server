package modules

type Model interface {
	List(limit int, offset int, sort string, order string, filter map[string]string) ([]interface{}, error)

	GetById(id string) (interface{}, error)

	Delete(id string)

	Total() *int64
}
