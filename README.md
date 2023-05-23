# cache

## 说明

- 轻量级的缓存字典

- 特色：线程安全、高性能、过期回调、健全的方法

- PS: 当前仅实现memcache, 后续考虑增加 Leveldb、SSDB、Redis等支持

## 版本

- go version >= 1.17

## 方法

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

## 灵感

- 灵感来源于[go-cache](https://github.com/patrickmn/go-cache)并新增了一些方法
