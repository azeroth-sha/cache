package memory

import (
	"time"
)

type Option func(b *barrel)

// WithCheckDur 自定义超时检查间隔
func WithCheckDur(dur time.Duration) Option {
	return func(b *barrel) {
		b.checkDur = dur
	}
}

// WithShardNum 自定义分片数
func WithShardNum(num int) Option {
	return func(b *barrel) {
		b.shardNum = uint32(num)
	}
}

// WithCallback 自定义超时回调方法
func WithCallback(fun Handler) Option {
	return func(b *barrel) {
		b.expireHandler = fun
	}
}

// WithTimer 自定义定时检查方法
// fun 将会轮询每一个元素，返回false时将会立即触发超时回调
func WithTimer(fun Handler) Option {
	return func(b *barrel) {
		b.checkHandler = fun
	}
}
