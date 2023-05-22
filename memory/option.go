package memory

import "time"

const Name = `memory`

type Option func(b *barrel)

// WithCheckDur 选择超时检查间隔
func WithCheckDur(dur time.Duration) Option {
	return func(b *barrel) {
		b.checkDur = dur
	}
}

// WithShardNum 选择分片数
func WithShardNum(num int) Option {
	return func(b *barrel) {
		b.shardNum = uint64(num)
	}
}

// WithCallback 选择超时回调方法
func WithCallback(call Callback) Option {
	return func(b *barrel) {
		b.callback = call
	}
}
