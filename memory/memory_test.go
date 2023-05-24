package memory_test

import (
	"fmt"
	"github.com/azeroth-sha/cache"
	"github.com/azeroth-sha/cache/memory"
	"log"
	"sync/atomic"
	"testing"
	"time"
)

var dict cache.Cache

func init() {
	dict = cache.New(memory.Name, memory.WithCallback(func(k string, v interface{}) {
		log.Printf("k: %s v: %v", k, v)
	}))
}

//func TestNew(t *testing.T) {
//	var limit = 50
//	var keys = make([]string, 0, limit)
//	for i := 0; i < limit; i++ {
//		k := fmt.Sprintf("%08x", i)
//		dict.Set(k, i)
//		keys = append(keys, k)
//	}
//	for _, key := range keys {
//		reply := dict.Get(key)
//		fmt.Printf("has: %t val: %v dur: %v\r\n", reply.Has(), reply.Val(), reply.Dur())
//	}
//	var cnt int32
//	dict.Len(func(k string, _ interface{}) bool {
//		cnt++
//		return true
//	})
//	fmt.Printf("count: %d\r\n", cnt)
//	dict.Range(func(k string, v interface{}) bool {
//		fmt.Printf("k: %s v: %v\r\n", k, v)
//		return true
//	})
//}

func BenchmarkSet(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.Set(
			fmt.Sprintf("%08x", i%100000),
			i,
		)
	}
}

func BenchmarkGet(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.Get(
			fmt.Sprintf("%08x", i%100000),
		)
	}
}

func BenchmarkDel(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.Del(
			fmt.Sprintf("%08x", i%100000),
		)
	}
}

func BenchmarkSetX(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.SetX(
			fmt.Sprintf("%08x", i%100000),
			i,
			time.Second,
		)
	}
}

func BenchmarkSetWithParallel(b *testing.B) {
	var i int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := atomic.AddInt64(&i, 1)
			dict.Set(
				fmt.Sprintf("%08x", n%100000),
				n,
			)
		}
	})
}

func BenchmarkGetWithParallel(b *testing.B) {
	var i int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := atomic.AddInt64(&i, 1)
			dict.Get(
				fmt.Sprintf("%08x", n%100000),
			)
		}
	})
}

func BenchmarkDelWithParallel(b *testing.B) {
	var i int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := atomic.AddInt64(&i, 1)
			dict.Del(
				fmt.Sprintf("%08x", n%100000),
			)
		}
	})
}

func BenchmarkSetXWithParallel(b *testing.B) {
	var i int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := atomic.AddInt64(&i, 1)
			dict.SetX(
				fmt.Sprintf("%08x", n%100000),
				n,
				time.Second,
			)
		}
	})
}
