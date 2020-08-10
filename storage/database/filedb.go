package database

type item struct {
	key []byte
	val []byte
}

// fileDB inserts an item, which has key and value in byte slice.
// It inserts the item to somewhere and returns the location of the item.
// An item can be retrieved with the returned location, URI.
type fileDB interface {
	write(items item) (string, error)
	read(key []byte) ([]byte, error)
	delete(key []byte) error
	deleteBucket()
}
