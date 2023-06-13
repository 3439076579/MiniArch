package singleflight

import "sync"

type singleFight interface {
	Do(key string, fn func() (interface{}, error)) (interface{}, error)
}
type caller struct {
	wg    sync.WaitGroup
	value interface{}
	err   error
}

type SingleCaller struct {
	mu  sync.Mutex
	set map[string]*caller
}

func (s *SingleCaller) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	s.mu.Lock()
	if s.set == nil {
		s.set = make(map[string]*caller)
	}
	if call, ok := s.set[key]; ok {
		s.mu.Unlock()
		call.wg.Wait()
		return call.value, call.err
	}
	call := new(caller)
	s.set[key] = call
	call.wg.Add(1)
	s.mu.Unlock()
	call.value, call.err = fn()

	call.wg.Done()
	s.mu.Lock()
	delete(s.set, key)
	s.mu.Unlock()
	return call.value, call.err
}
