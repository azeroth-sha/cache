package memory

import "time"

type reply struct {
	err error
	has bool
	val interface{}
	dur time.Duration
}

func (r *reply) Err() error {
	return r.err
}

func (r *reply) Has() bool {
	return r.has
}

func (r *reply) Val() interface{} {
	return r.val
}

func (r *reply) Dur() time.Duration {
	return r.dur
}

func (r *reply) Release() {
	r.init()
	resPut(r)
}

func (r *reply) init() {
	r.err = nil
	r.has = false
	r.val = nil
	r.dur = 0
}
