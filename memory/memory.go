package memory

import (
	"github.com/azeroth-sha/cache"
	"hash/fnv"
	"runtime"
	"sync"
	"time"
)

const Name = `memory`

// Handler 回调方法
type Handler func(k string, i Item) bool

func defaultCheck(_ string, i Item) bool {
	return !i.Expired()
}

type Cache interface {
	cache.Cache
	Range(Handler)
}

type barrel struct {
	checkDur      time.Duration
	shardNum      uint32
	shards        []*shard
	expireHandler Handler
	checkHandler  Handler
	closed        chan struct{}
}

func (b *barrel) getShard(k string) *shard {
	sha := fnv.New32a()
	_, _ = sha.Write(toBytes(k))
	return b.shards[sha.Sum32()&(b.shardNum-1)]
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

func (b *barrel) TTL(k string) cache.Reply {
	return b.getShard(k).TTL(k)
}

func (b *barrel) Range(f Handler) {
	r := resGet()
	defer r.Release()
	for _, s := range b.shards {
		if s.Range(f, r); !r.Has() {
			break
		}
		r.init()
	}
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
			for _, s := range b.shards {
				s.check(b.checkHandler)
			}
		}
	}
}

func New(opts ...interface{}) Cache {
	b := &barrel{
		checkDur:     time.Second,
		checkHandler: defaultCheck,
		closed:       make(chan struct{}),
	}
	for _, opt := range opts {
		if option, ok := opt.(Option); ok {
			option(b)
		}
	}
	if b.shardNum <= 0 {
		b.shardNum = uint32(runtime.NumCPU() * 2)
	}
	b.shards = make([]*shard, 0, b.shardNum)
	for i := uint32(0); i < b.shardNum; i++ {
		b.shards = append(b.shards, &shard{
			dictLock:      new(sync.RWMutex),
			expireHandler: b.expireHandler,
			dict:          make(map[string]*item),
		})
	}
	go b.check()
	runtime.SetFinalizer(b, func(c *barrel) { close(c.closed) })
	return b
}

func init() {
	cache.Register(Name, func(i ...interface{}) cache.Cache {
		return New(i...)
	})
}
