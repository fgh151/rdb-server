package drivers

// DbConnection Database connection interface
type DbConnection interface {
	GetDsn() string
}
