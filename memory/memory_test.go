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
	dict = cache.New(
		memory.Name,
		memory.WithCallback(callback),
	)
}

func callback(k string, i memory.Item) bool {
	log.Printf("k: %s value: %value", k, i.Value())
	return false
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
//		fmt.Printf("has: %t val: %value dur: %value\r\n", reply.Has(), reply.Val(), reply.TTL())
//	}
//	var cnt int32
//	dict.Len(func(k string, _ interface{}) bool {
//		cnt++
//		return true
//	})
//	fmt.Printf("count: %d\r\n", cnt)
//	dict.Handler(func(k string, value interface{}) bool {
//		fmt.Printf("k: %s value: %value\r\n", k, value)
//		return true
//	})
//}
//
//func BenchmarkSet(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		dict.Set(
//			fmt.Sprintf("%08x", i%100000),
//			i,
//		)
//	}
//}
//
//func BenchmarkGet(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		dict.Get(
//			fmt.Sprintf("%08x", i%100000),
//		).Release()
//	}
//}
//
//func BenchmarkDel(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		dict.Del(
//			fmt.Sprintf("%08x", i%100000),
//		).Release()
//	}
//}
//
//func BenchmarkSetX(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		dict.SetX(
//			fmt.Sprintf("%08x", i%100000),
//			i,
//			time.Second,
//		).Release()
//	}
//}

func BenchmarkSetNX(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.SetNX(
			fmt.Sprintf("%08x", i%100000),
			i,
			time.Second,
		).Release()
	}
}

//
//func BenchmarkSetWithParallel(b *testing.B) {
//	var i int64
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			n := atomic.AddInt64(&i, 1)
//			dict.Set(
//				fmt.Sprintf("%08x", n%100000),
//				n,
//			)
//		}
//	})
//}
//
//func BenchmarkGetWithParallel(b *testing.B) {
//	var i int64
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			n := atomic.AddInt64(&i, 1)
//			dict.Get(
//				fmt.Sprintf("%08x", n%100000),
//			).Release()
//		}
//	})
//}
//
//func BenchmarkDelWithParallel(b *testing.B) {
//	var i int64
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			n := atomic.AddInt64(&i, 1)
//			dict.Del(
//				fmt.Sprintf("%08x", n%100000),
//			).Release()
//		}
//	})
//}
//
//func BenchmarkSetXWithParallel(b *testing.B) {
//	var i int64
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			n := atomic.AddInt64(&i, 1)
//			dict.SetX(
//				fmt.Sprintf("%08x", n%100000),
//				n,
//				time.Second,
//			).Release()
//		}
//	})
//}

func BenchmarkSetNXWithParallel(b *testing.B) {
	var i int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := atomic.AddInt64(&i, 1)
			dict.SetNX(
				fmt.Sprintf("%08x", n%100000),
				n,
				time.Second,
			).Release()
		}
	})
}
