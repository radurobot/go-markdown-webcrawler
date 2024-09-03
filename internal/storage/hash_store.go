package storage

type HashStore interface {
	Add(hash string)
	Exists(hash string) bool
	Close() error
}
