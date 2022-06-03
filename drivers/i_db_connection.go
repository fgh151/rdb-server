package drivers

type DbConnection interface {
	GetDsn() string
}
