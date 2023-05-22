package memory

import "time"

type reply struct {
	has bool
	val interface{}
	err error
	dur time.Duration
}

func (r *reply) Has() bool {
	return r.has
}

func (r *reply) Val() interface{} {
	return r.val
}

func (r *reply) Err() error {
	return r.err
}

func (r *reply) Dur() time.Duration {
	return r.dur
}
