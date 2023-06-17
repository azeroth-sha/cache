package memory

import (
	"sync"
	"unsafe"
)

var (
	itemPool = &sync.Pool{New: itemNew}
	resPool  = &sync.Pool{New: resNew}
)

func toBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func toString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func itemNew() interface{} {
	return new(item)
}

func itemGet() *item {
	return itemPool.Get().(*item)
}

func itemPut(i *item) {
	itemPool.Put(i)
}

func resNew() interface{} {
	return new(reply)
}

func resGet() *reply {
	return resPool.Get().(*reply)
}

func resPut(res *reply) {
	resPool.Put(res)
}
