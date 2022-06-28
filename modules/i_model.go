package modules

type Model interface {
	List(limit int, offset int, sort string, order string, filter map[string]interface{}) []interface{}

	GetById(id string) interface{}

	Delete(id string)

	Total() *int64
}
