# cache

## 说明

- 轻量级的缓存字典

- 特色：线程安全、高性能、过期回调、健全的方法

- PS: 当前仅实现memcache, 后续考虑增加 Leveldb、SSDB、Redis等支持

## 版本

- go version >= 1.17

## 灵感

- 灵感来源于[go-cache](https://github.com/patrickmn/go-cache)并新增了一些方法
