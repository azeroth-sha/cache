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
	r := new(reply)
	s.mu.RLock()
	r.has = s.delExpired(k)
	s.mu.RUnlock()
	return r
}

func (s *shard) Set(k string, v interface{}) cache.Reply {
	r := new(reply)
	s.mu.Lock()
	s.delExpired(k)
	s.dict[k] = &item{v: v}
	s.mu.Unlock()
	return r
}

func (s *shard) SetX(k string, v interface{}, d time.Duration) cache.Reply {
	r := new(reply)
	e := time.Now().Add(d)
	s.mu.Lock()
	s.delExpired(k)
	s.dict[k] = &item{v: v, e: &e}
	s.mu.Unlock()
	return r
}

func (s *shard) SetN(k string, v interface{}) cache.Reply {
	r := new(reply)
	s.mu.Lock()
	if r.has = s.delExpired(k); r.has {
		return r
	}
	s.dict[k] = &item{v: v}
	s.mu.Unlock()
	return r
}

func (s *shard) SetNX(k string, v interface{}, d time.Duration) cache.Reply {
	r := new(reply)
	e := time.Now().Add(d)
	s.mu.Lock()
	if r.has = s.delExpired(k); r.has {
		return r
	}
	s.dict[k] = &item{v: v, e: &e}
	s.mu.Unlock()
	return r
}

func (s *shard) Del(k string) cache.Reply {
	r := new(reply)
	s.mu.Lock()
	if s.delExpired(k) {
		delete(s.dict, k)
		r.has = true
	}
	s.mu.Unlock()
	return r
}

func (s *shard) DelExpired(k string) cache.Reply {
	r := new(reply)
	s.mu.Lock()
	r.has = s.delExpired(k)
	s.mu.Unlock()
	return r
}

func (s *shard) Get(k string) cache.Reply {
	r := new(reply)
	s.mu.RLock()
	if s.delExpired(k) {
		r.has = true
		r.val = s.dict[k].v
	}
	s.mu.RUnlock()
	return r
}

func (s *shard) GetDel(k string) cache.Reply {
	r := new(reply)
	s.mu.Lock()
	if s.delExpired(k) {
		r.has = true
		r.val = s.dict[k].v
		delete(s.dict, k)
	}
	s.mu.Unlock()
	return r
}

func (s *shard) GetSet(k string, v interface{}) cache.Reply {
	r := new(reply)
	s.mu.Lock()
	if s.delExpired(k) {
		r.has = true
		r.val = s.dict[k].v
	}
	s.dict[k] = &item{v: v}
	s.mu.Unlock()
	return r
}

func (s *shard) GetSetX(k string, v interface{}, d time.Duration) cache.Reply {
	r := new(reply)
	e := time.Now().Add(d)
	s.mu.Lock()
	if s.delExpired(k) {
		r.has = true
		r.val = s.dict[k].v
	}
	s.dict[k] = &item{v: v, e: &e}
	s.mu.Unlock()
	return r
}

func (s *shard) Expire(k string, d time.Duration) cache.Reply {
	r := new(reply)
	e := time.Now().Add(d)
	s.mu.Lock()
	if s.delExpired(k) {
		r.has = true
		s.dict[k].e = &e
	}
	s.mu.Unlock()
	return r
}

func (s *shard) Dur(k string) cache.Reply {
	r := new(reply)
	s.mu.RLock()
	if s.delExpired(k) {
		r.has = true
		i := s.dict[k]
		if i.e != nil {
			r.dur = i.e.Sub(time.Now())
		}
	}
	s.mu.RUnlock()
	return r
}

func (s *shard) Len(f cache.RangeFunc) cache.Reply {
	r := new(reply)
	kList := make([]string, 0)
	defer func() {
		if len(kList) > 0 {
			go s.checkKeys(kList)
		}
	}()
	s.mu.RLock()
	if len(s.dict) == 0 {
		r.has = true
		goto EXIT
	}
	for k, i := range s.dict {
		if i.Expired() {
			kList = append(kList, k)
			continue
		}
		if r.has = f(k, nil); !r.has {
			break
		}
	}
EXIT:
	s.mu.RUnlock()
	return r
}

func (s *shard) Range(f cache.RangeFunc) cache.Reply {
	r := new(reply)
	kList := make([]string, 0)
	defer func() {
		if len(kList) > 0 {
			go s.checkKeys(kList)
		}
	}()
	s.mu.RLock()
	if len(s.dict) == 0 {
		r.has = true
		goto EXIT
	}
	for k, i := range s.dict {
		if i.Expired() {
			kList = append(kList, k)
			continue
		}
		if r.has = f(k, i.v); !r.has {
			break
		}
	}
EXIT:
	s.mu.RUnlock()
	return r
}

func (s *shard) delExpired(k string) bool {
	if i, ok := s.dict[k]; !ok {
		return false
	} else if i.Expired() {
		delete(s.dict, k)
		if s.call != nil {
			s.call(k, i.v)
		}
		return false
	}
	return true
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

func (s *shard) checkKeys(keys []string) {
	for _, key := range keys {
		s.DelExpired(key)
	}
}
