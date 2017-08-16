package atomicmap

import "qnmahjong/db/dao"

// +gen atomicmap
type Name string

// +gen atomicmap
type Age int

// +gen atomicmap
type Person struct {
	Name string
	Age  int
}

// +gen atomicmap
type Login struct {
	ID        int32
	Channel   int32
	Version   int32
	LoginType int32
	LoginTime int64
	dao.Player
}
