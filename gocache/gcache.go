package gocache

import (
	"MiniArch/gocache/gcachepb"
	"MiniArch/gocache/singleflight"
	"context"
	"errors"
	"sync"
)

type Getter interface {
	Get(ctx context.Context, key string) ([]byte, error)
}

type GetFunc func(ctx context.Context, key string) ([]byte, error)

func (f GetFunc) Get(ctx context.Context, key string) ([]byte, error) {
	return f(ctx, key)
}

func NewGroup(name string, callback Getter, cacheBytes int64) *Group {
	mu.Lock()
	defer mu.Unlock()
	if callback == nil {
		panic("nil CallBack function")
	}
	if _, ok := groups[name]; ok {
		panic("duplicate register cache")
	}
	group := newGroup(name, callback, cacheBytes)
	groups[name] = group
	return group
}

func newGroup(name string, callback Getter, cacheBytes int64) *Group {
	return &Group{
		name:       name,
		callback:   callback,
		cacheBytes: cacheBytes,
		loadOnce:   &singleflight.SingleCaller{},
		mainCache:  new(cache),
	}
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func GetGroup(GroupName string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[GroupName]
}

type Group struct {
	name string

	// callback when you cannot get a key from this machine or promote machines,
	// this Getter function will be automatically called
	callback Getter

	mainCache *cache

	loadOnce *singleflight.SingleCaller

	picker PromoteGetterPicker

	cacheBytes int64
}

// Get will try to get value from local group firstly
func (g *Group) Get(ctx context.Context, key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key cannot be empty")
	}
	view, ok := g.mainCache.get(key)
	if ok {
		return view, nil
	}
	return g.load(ctx, key)
}

func (g *Group) load(ctx context.Context, key string) (ByteView, error) {
	value, err := g.loadOnce.Do(key, func() (interface{}, error) {
		if value, cacheHit := g.mainCache.get(key); cacheHit {
			return value, nil
		}
		if g.picker == nil {
			panic("group picker must be registered")
		}
		getter, ok := g.picker.PickPromoteGetter(key)
		if ok {
			value, err := g.getFromPeer(ctx, getter, key)
			if err == nil {
				return value, nil
			}
		}
		value, err := g.getLocally(ctx, key)
		if err != nil {
			return nil, err
		}
		g.populateCache(key, value)
		return value, nil
	})
	if err == nil {
		return value.(ByteView), nil
	}
	return ByteView{}, err
}

// getLocally will call g.callback function to get data from source
func (g *Group) getLocally(ctx context.Context, key string) (ByteView, error) {
	if g.callback == nil {
		return ByteView{}, errors.New("")
	}
	bytes, err := g.callback.Get(ctx, key)
	if err != nil {
		return ByteView{}, errors.New("")
	}
	value := ByteView{b: cloneBytes(bytes)}
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	if g.cacheBytes <= 0 {
		return
	}
	g.mainCache.add(key, value)
}
func (g *Group) getFromPeer(ctx context.Context, getter PromoteGetter, key string) (ByteView, error) {
	req := &gcachepb.GetRequest{
		Group: g.name,
		Key:   key,
	}
	rep := &gcachepb.GetResponse{}

	err := getter.GetFromPromote(ctx, req, rep)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{s: rep.Value}, nil
}
func (g *Group) RegisterPicker(picker PromoteGetterPicker) {
	if g.picker != nil {
		panic("picker has been register")
	}
	g.picker = picker
}

type cache struct {
	mu     sync.RWMutex
	lru    *Cache // lazily init
	nBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = New(0)
		c.lru.SetEvictedFunc(func(key Key, value Any) {
			val := value.(ByteView)
			c.nBytes -= int64(len(key.(string)) + val.Len())
		})
	}
	IsSwap := c.lru.Add(key, value)
	if IsSwap {
		view := c.lru.cache[key].value.(*entry).value.(*ByteView)
		c.nBytes += int64(view.Len() - value.Len())
	} else {
		c.nBytes += int64(len(key) + value.Len())
	}
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.lru == nil {
		c.lru = New(0)
		c.lru.SetEvictedFunc(func(key Key, value Any) {
			val := value.(ByteView)
			c.nBytes -= int64(len(key.(string)) + val.Len())
		})
	}
	value, ok := c.lru.Get(key)
	if !ok {
		return ByteView{}, false
	}
	return value.(ByteView), true
}

func (c *cache) removeOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = New(0)
		c.lru.SetEvictedFunc(func(key Key, value Any) {
			val := value.(ByteView)
			c.nBytes -= int64(len(key.(string)) + val.Len())
		})
	}
	c.lru.DeleteOldest()
}
func (c *cache) bytes() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nBytes
}

func (c *cache) Items() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.lru == nil {
		return 0
	}
	return c.lru.Len()
}
