package memory

import (
	"sync"
	"time"
)

type shard struct {
	mu   *sync.RWMutex
	call Handler
	dict map[string]*item
}

func (s *shard) Has(k string) *reply {
	r := resGet()
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.setReply(s.dict[k], r)
	return r
}

func (s *shard) Set(k string, v interface{}) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		i := s.dict[k]
		i.value = v
		i.expireHas = false
	} else {
		i := itemGet()
		i.value = v
		i.expireHas = false
		s.dict[k] = i
	}
	return r
}

func (s *shard) SetX(k string, v interface{}, d time.Duration) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		i := s.dict[k]
		i.value = v
		i.expireHas = true
		i.expire = time.Now().Add(d)
	} else {
		i := itemGet()
		i.value = v
		i.expireHas = true
		i.expire = time.Now().Add(d)
		s.dict[k] = i
	}
	return r
}

func (s *shard) SetN(k string, v interface{}) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		s.setReply(s.dict[k], r)
		return r
	}
	i := itemGet()
	i.value = v
	i.expireHas = false
	s.dict[k] = i
	return r
}

func (s *shard) SetNX(k string, v interface{}, d time.Duration) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		s.setReply(s.dict[k], r)
		return r
	}
	i := itemGet()
	i.value = v
	i.expireHas = true
	i.expire = time.Now().Add(d)
	s.dict[k] = i
	return r
}

func (s *shard) Del(k string) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		s.setReply(s.dict[k], r)
		itemPut(s.dict[k])
		delete(s.dict, k)
	}
	return r
}

func (s *shard) DelExpired(k string) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setReply(s.dict[k], r)
	s.delExpired(k)
	return r
}

func (s *shard) Get(k string) *reply {
	r := resGet()
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.delExpired(k) {
		s.setReply(s.dict[k], r)
	}
	return r
}

func (s *shard) GetDel(k string) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		i := s.dict[k]
		s.setReply(i, r)
		itemPut(i)
		delete(s.dict, k)
	}
	return r
}

func (s *shard) GetSet(k string, v interface{}) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		i := s.dict[k]
		s.setReply(i, r)
		i.value = v
		i.expireHas = false
	} else {
		i := itemGet()
		i.value = v
		i.expireHas = false
		s.dict[k] = i
	}
	return r
}

func (s *shard) GetSetX(k string, v interface{}, d time.Duration) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		i := s.dict[k]
		s.setReply(i, r)
		i.value = v
		i.expireHas = false
		i.expire = time.Now().Add(d)
	} else {
		i := itemGet()
		i.value = v
		i.expireHas = false
		i.expire = time.Now().Add(d)
		s.dict[k] = i
	}
	return r
}

func (s *shard) Expire(k string, d time.Duration) *reply {
	r := resGet()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delExpired(k) {
		i := s.dict[k]
		s.setReply(i, r)
		if d < 0 {
			i.expireHas = false
			return r
		}
		i.expireHas = true
		i.expire = time.Now().Add(d)
	}
	return r
}

func (s *shard) TTL(k string) *reply {
	r := resGet()
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.delExpired(k) {
		i := s.dict[k]
		s.setReply(i, r)
	}
	return r
}

func (s *shard) Range(f Handler, r *reply) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r.has = true
	if len(s.dict) == 0 {
		return
	}
	for k, i := range s.dict {
		if r.has = f(k, i); !r.has {
			break
		}
	}
}

func (s *shard) delExpired(k string) bool {
	if i, ok := s.dict[k]; !ok {
		return false
	} else if i.Expired() {
		delete(s.dict, k)
		if s.call != nil {
			s.call(k, i)
		}
		itemPut(i)
		return false
	}
	return true
}

func (s *shard) check(f Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, i := range s.dict {
		if !f(k, i) {
			if s.call != nil {
				s.call(k, i)
			}
			itemPut(i)
			delete(s.dict, k)
		}
	}
}

func (s *shard) checkKeys(keys []string) {
	for _, key := range keys {
		s.DelExpired(key)
	}
}

func (s *shard) setReply(i *item, r *reply) {
	if i == nil {
		r.init()
	} else {
		r.err = nil
		r.has = true
		r.val = i.value
		r.dur = i.Expire()
	}
}
