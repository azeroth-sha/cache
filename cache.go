package cache

import (
	"sync"
	"time"
)

var adapterMap = new(sync.Map)

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

// NewFunc 缓存器创建方法
// 参数自定义，创建缓存器时透传
type NewFunc func(...interface{}) Cache

// Register 注册缓存生成器
func Register(name string, f NewFunc) {
	adapterMap.Store(name, f)
}

// New 创建新的缓存器
func New(name string, sets ...interface{}) Cache {
	if newFun, found := adapterMap.Load(name); found {
		return newFun.(NewFunc)(sets...)
	}
	return nil
}
