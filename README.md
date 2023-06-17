# cache

## 说明

- 轻量级的缓存字典

- 特色：线程安全、高性能、过期回调、健全的方法

- PS: 当前仅实现memcache, 后续考虑增加 Leveldb、SSDB、Redis等支持

## 版本

- go version >= 1.17

## 方法说明

```go
package cache

// Reply 缓存器响应内容
type Reply interface {
    Has() bool          // 是否存在
    Val() interface{}   // 返回的值
    Err() error         // 返回的错误
    Dur() time.Duration // 剩余时长
    Release()           // 释放响应结果
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
    TTL(k string) Reply                                     // 获取过期时间
}
```

## 基准参考

```
goos: windows
goarch: amd64
pkg: github.com/azeroth-sha/cache/memory
cpu: Intel(R) Xeon(R) CPU E3-1231 v3 @ 3.40GHz
BenchmarkSet-8                   6725599               531.3 ns/op            76 B/op          3 allocs/op
BenchmarkSet-8                   7228197               516.8 ns/op            75 B/op          3 allocs/op
BenchmarkSet-8                   7281092               503.2 ns/op            75 B/op          3 allocs/op
BenchmarkGet-8                  11582373               301.1 ns/op            18 B/op          1 allocs/op
BenchmarkGet-8                  12887090               314.6 ns/op            18 B/op          1 allocs/op
BenchmarkGet-8                  10562623               298.8 ns/op            17 B/op          1 allocs/op
BenchmarkDel-8                  17233570               191.9 ns/op            15 B/op          1 allocs/op
BenchmarkDel-8                  18140698               196.7 ns/op            15 B/op          1 allocs/op
BenchmarkDel-8                  18709345               195.9 ns/op            15 B/op          1 allocs/op
BenchmarkSetX-8                  9990625               384.4 ns/op            26 B/op          2 allocs/op
BenchmarkSetX-8                  8875980               351.9 ns/op            26 B/op          2 allocs/op
BenchmarkSetX-8                 10816364               340.7 ns/op            26 B/op          2 allocs/op
BenchmarkSetWithParallel-8      22513000               146.9 ns/op            72 B/op          3 allocs/op
BenchmarkSetWithParallel-8      24671593               143.8 ns/op            72 B/op          3 allocs/op
BenchmarkSetWithParallel-8      25524890               145.6 ns/op            73 B/op          3 allocs/op
BenchmarkGetWithParallel-8      45086565                92.57 ns/op           16 B/op          1 allocs/op
BenchmarkGetWithParallel-8      45930435                92.79 ns/op           16 B/op          1 allocs/op
BenchmarkGetWithParallel-8      44680366                97.86 ns/op           16 B/op          1 allocs/op
BenchmarkDelWithParallel-8      50732734                66.31 ns/op           15 B/op          1 allocs/op
BenchmarkDelWithParallel-8      49343863                66.49 ns/op           15 B/op          1 allocs/op
BenchmarkDelWithParallel-8      53936947                67.42 ns/op           15 B/op          1 allocs/op
BenchmarkSetXWithParallel-8     34435525               100.3 ns/op            24 B/op          2 allocs/op
BenchmarkSetXWithParallel-8     36729596               101.5 ns/op            24 B/op          2 allocs/op
BenchmarkSetXWithParallel-8     33610495               106.0 ns/op            24 B/op          2 allocs/op
PASS
ok      github.com/azeroth-sha/cache/memory     93.894s
```

## 灵感来源

- 灵感来源于[go-cache](https://github.com/patrickmn/go-cache)并新增了一些方法
