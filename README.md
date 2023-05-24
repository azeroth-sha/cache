# cache

## 说明

- 轻量级的缓存字典

- 特色：线程安全、高性能、过期回调、健全的方法

- PS: 当前仅实现memcache, 后续考虑增加 Leveldb、SSDB、Redis等支持

## 版本

- go version >= 1.17

## 方法说明

```go
// Reply 缓存器响应内容
type Reply interface {
	Has() bool          // 是否存在
	Val() interface{}   // 返回的值
	Err() error         // 返回的错误
	Dur() time.Duration // 剩余时长
}

// Cache 字典缓存器
type Cache interface {
	Has(k string) Reply                                     // 是否存在
	Set(k string, v interface{}) Reply                      // 设置键值
	SetX(k string, v interface{}, d time.Duration) Reply    // 设置键值并设置过期时长
	SetN(k string, v interface{}) Reply                     // 设置键值存在时返回错误
	SetNX(k string, v interface{}, d time.Duration) Reply   // 设置键值并设置过期时长，如果存在则返回错误
	Del(k string) Reply                                     // 删除键值(不存在为成功)
	DelExpired(k string) Reply                              // 删除过期键值(不存在为成功)
	Get(k string) Reply                                     // 获取键值
	GetDel(k string) Reply                                  // 获取并删除键值
	GetSet(k string, v interface{}) Reply                   // 获取并设置新的键值
	GetSetX(k string, v interface{}, d time.Duration) Reply // 获取并设置新的键值和过期时长
	Expire(k string, d time.Duration) Reply                 // 设置新的过期时间
	Dur(k string) Reply                                     // 获取过期时间
	Len(f RangeFunc) Reply                                  // 遍历键
	Range(f RangeFunc) Reply                                // 遍历键值
}

// RangeFunc 遍历回调方法
type RangeFunc func(k string, v interface{}) bool
```

## 基准参考

```
goos: windows
goarch: amd64
pkg: github.com/azeroth-sha/cache/memory
cpu: Intel(R) Xeon(R) CPU E3-1231 v3 @ 3.40GHz
BenchmarkSet-8                   5631349               613.8 ns/op            95 B/op          4 allocs/op
BenchmarkSet-8                   5577706               617.4 ns/op            95 B/op          4 allocs/op
BenchmarkSet-8                   5821218               627.6 ns/op            95 B/op          4 allocs/op
BenchmarkGet-8                   6949197               506.9 ns/op            63 B/op          2 allocs/op
BenchmarkGet-8                   7159614               500.1 ns/op            63 B/op          2 allocs/op
BenchmarkGet-8                   7267500               501.8 ns/op            63 B/op          2 allocs/op
BenchmarkDel-8                  10551502               322.0 ns/op            64 B/op          2 allocs/op
BenchmarkDel-8                  11155561               320.9 ns/op            64 B/op          2 allocs/op
BenchmarkDel-8                  11227999               327.3 ns/op            64 B/op          2 allocs/op
BenchmarkSetX-8                  5064267               714.8 ns/op           119 B/op          5 allocs/op
BenchmarkSetX-8                  5028674               755.2 ns/op           119 B/op          5 allocs/op
BenchmarkSetX-8                  4281289               735.4 ns/op           119 B/op          5 allocs/op
BenchmarkSetWithParallel-8      13636311               255.0 ns/op            95 B/op          4 allocs/op
BenchmarkSetWithParallel-8      13856289               253.1 ns/op            95 B/op          4 allocs/op
BenchmarkSetWithParallel-8      15273561               243.6 ns/op            95 B/op          4 allocs/op
BenchmarkGetWithParallel-8      23599113               161.9 ns/op            63 B/op          2 allocs/op
BenchmarkGetWithParallel-8      23006031               162.7 ns/op            63 B/op          2 allocs/op
BenchmarkGetWithParallel-8      23599701               164.0 ns/op            63 B/op          2 allocs/op
BenchmarkDelWithParallel-8      25255748               129.6 ns/op            64 B/op          2 allocs/op
BenchmarkDelWithParallel-8      24468538               124.8 ns/op            64 B/op          2 allocs/op
BenchmarkDelWithParallel-8      28035482               128.7 ns/op            64 B/op          2 allocs/op
BenchmarkSetXWithParallel-8     12816495               294.1 ns/op           119 B/op          5 allocs/op
BenchmarkSetXWithParallel-8     12107102               291.6 ns/op           119 B/op          5 allocs/op
BenchmarkSetXWithParallel-8     12172800               287.5 ns/op           119 B/op          5 allocs/op
PASS
ok      github.com/azeroth-sha/cache/memory     98.954s
```

## 灵感来源

- 灵感来源于[go-cache](https://github.com/patrickmn/go-cache)并新增了一些方法
