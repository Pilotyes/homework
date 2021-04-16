package storage

//Storage ...
type Storage interface {
	Items() ItemsRepository
}
