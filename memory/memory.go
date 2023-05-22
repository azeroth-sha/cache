package memory

import (
	"github.com/azeroth-sha/cache"
	"hash/fnv"
	"runtime"
	"sync"
	"time"
)

type Callback func(k string, v interface{})

type barrel struct {
	checkDur time.Duration
	shardNum uint64
	shardMap map[uint64]*shard
	callback Callback
	closed   chan struct{}
}

func (b *barrel) getShard(k string) *shard {
	sha := fnv.New64a()
	_, _ = sha.Write([]byte(k))
	return b.shardMap[sha.Sum64()%b.shardNum]
}

func (b *barrel) Has(k string) cache.Reply {
	return b.getShard(k).Has(k)
}

func (b *barrel) Set(k string, v interface{}) cache.Reply {
	return b.getShard(k).Set(k, v)
}

func (b *barrel) SetX(k string, v interface{}, d time.Duration) cache.Reply {
	return b.getShard(k).SetX(k, v, d)
}

func (b *barrel) SetN(k string, v interface{}) cache.Reply {
	return b.getShard(k).SetN(k, v)
}

func (b *barrel) SetNX(k string, v interface{}, d time.Duration) cache.Reply {
	return b.getShard(k).SetNX(k, v, d)
}

func (b *barrel) Del(k string) cache.Reply {
	return b.getShard(k).Del(k)
}

func (b *barrel) DelExpired(k string) cache.Reply {
	return b.getShard(k).DelExpired(k)
}

func (b *barrel) Get(k string) cache.Reply {
	return b.getShard(k).Get(k)
}

func (b *barrel) GetDel(k string) cache.Reply {
	return b.getShard(k).GetDel(k)
}

func (b *barrel) GetSet(k string, v interface{}) cache.Reply {
	return b.getShard(k).GetSet(k, v)
}

func (b *barrel) GetSetX(k string, v interface{}, d time.Duration) cache.Reply {
	return b.getShard(k).GetSetX(k, v, d)
}

func (b *barrel) Expire(k string, d time.Duration) cache.Reply {
	return b.getShard(k).Expire(k, d)
}

func (b *barrel) Dur(k string) cache.Reply {
	return b.getShard(k).Dur(k)
}

func (b *barrel) Len(f cache.RangeFunc) cache.Reply {
	var r cache.Reply
	for _, s := range b.shardMap {
		r = s.Len(f)
	}
	return r
}

func (b *barrel) Range(f cache.RangeFunc) cache.Reply {
	var r cache.Reply
	for _, s := range b.shardMap {
		r = s.Range(f)
	}
	return r
}

func (b *barrel) check() {
	if b.checkDur <= 0 {
		return
	}
	tk := time.NewTicker(b.checkDur)
	defer tk.Stop()
EXIT:
	for true {
		select {
		case <-b.closed:
			break EXIT
		case <-tk.C:
			for _, s := range b.shardMap {
				s.check()
			}
		}
	}
}

func New(opts ...interface{}) cache.Cache {
	b := &barrel{
		checkDur: time.Second,
		shardNum: uint64(runtime.NumCPU() * 2),
		shardMap: make(map[uint64]*shard),
		callback: nil,
		closed:   make(chan struct{}),
	}
	for _, opt := range opts {
		if option, ok := opt.(Option); ok {
			option(b)
		}
	}
	if b.shardNum <= 0 {
		b.shardNum = uint64(runtime.NumCPU() * 2)
	}
	for i := uint64(0); i < b.shardNum; i++ {
		b.shardMap[i] = &shard{
			mu:   new(sync.RWMutex),
			call: b.callback,
			dict: make(map[string]*item),
		}
	}
	go b.check()
	runtime.SetFinalizer(b, func(c *barrel) { close(c.closed) })
	return b
}

func init() {
	cache.Register(Name, New)
}
