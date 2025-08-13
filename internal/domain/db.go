package domain

type Database[T any] interface {
	Save(id string, obj T) error
	Get(id string) (T, error)
	Remove(id string) error
}
