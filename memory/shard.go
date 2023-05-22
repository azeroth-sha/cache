package memory

import (
	"github.com/azeroth-sha/cache"
	"sync"
	"time"
)

type shard struct {
	mu   *sync.RWMutex
	call Callback
	dict map[string]*item
}

func (s *shard) Has(k string) cache.Reply {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, found := s.dict[k]
	return &reply{has: found}
}

func (s *shard) Set(k string, v interface{}) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dict[k] = &item{v: v}
	return new(reply)
}

func (s *shard) SetX(k string, v interface{}, d time.Duration) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.Lock()
	defer s.mu.Unlock()
	e := time.Now().Add(d)
	s.dict[k] = &item{v: v, e: &e}
	return new(reply)
}

func (s *shard) SetN(k string, v interface{}) cache.Reply {
	_ = s.DelExpired(k)
	if r := s.Has(k); r.Has() {
		return r
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dict[k] = &item{v: v}
	return new(reply)
}

func (s *shard) SetNX(k string, v interface{}, d time.Duration) cache.Reply {
	_ = s.DelExpired(k)
	if r := s.Has(k); r.Has() {
		return r
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e := time.Now().Add(d)
	s.dict[k] = &item{v: v, e: &e}
	return new(reply)
}

func (s *shard) Del(k string) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.Lock()
	defer s.mu.Unlock()
	r := new(reply)
	if _, ok := s.dict[k]; ok {
		delete(s.dict, k)
		r.has = true
	}
	return r
}

func (s *shard) DelExpired(k string) cache.Reply {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := new(reply)
	if i, ok := s.dict[k]; ok && i.Expired() {
		delete(s.dict, k)
		r.has = true
		if s.call != nil {
			s.call(k, i.v)
		}
	}
	return r
}

func (s *shard) Get(k string) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := new(reply)
	if i, ok := s.dict[k]; ok {
		r.has = true
		r.val = i.v
	}
	return r
}

func (s *shard) GetDel(k string) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.Lock()
	defer s.mu.Unlock()
	r := new(reply)
	if i, ok := s.dict[k]; ok {
		delete(s.dict, k)
		r.has = true
		r.val = i.v
	}
	return r
}

func (s *shard) GetSet(k string, v interface{}) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.Lock()
	defer s.mu.Unlock()
	r := new(reply)
	if i, ok := s.dict[k]; ok {
		r.has = true
		r.val = i.v
	}
	s.dict[k] = &item{v: v}
	return r
}

func (s *shard) GetSetX(k string, v interface{}, d time.Duration) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.Lock()
	defer s.mu.Unlock()
	r := new(reply)
	if i, ok := s.dict[k]; ok {
		r.has = true
		r.val = i.v
	}
	e := time.Now().Add(d)
	s.dict[k] = &item{v: v, e: &e}
	return r
}

func (s *shard) Expire(k string, d time.Duration) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.Lock()
	defer s.mu.Unlock()
	r := new(reply)
	if i, ok := s.dict[k]; ok {
		e := time.Now().Add(d)
		i.e = &e
		r.has = true
	}
	return r
}

func (s *shard) Dur(k string) cache.Reply {
	_ = s.DelExpired(k)
	s.mu.RLock()
	defer s.mu.RUnlock()
	r := new(reply)
	if i, ok := s.dict[k]; ok {
		r.has = true
		if i.e != nil {
			r.dur = i.e.Sub(time.Now())
		}
	}
	return r
}

func (s *shard) Len(f cache.RangeFunc) cache.Reply {
	s.check()
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k := range s.dict {
		if !f(k, nil) {
			break
		}
	}
	return new(reply)
}

func (s *shard) Range(f cache.RangeFunc) cache.Reply {
	s.check()
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, i := range s.dict {
		if !f(k, i.v) {
			break
		}
	}
	return new(reply)
}

func (s *shard) check() {
	s.mu.Lock()
	defer s.mu.Unlock()
	kList := make([]string, 0)
	for k, i := range s.dict {
		if i.Expired() {
			kList = append(kList, k)
		}
	}
	for _, k := range kList {
		i := s.dict[k]
		delete(s.dict, k)
		if s.call != nil {
			s.call(k, i.v)
		}
	}
}
