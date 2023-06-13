package gocache

/*
**
实现LRU算法
*/

type Any interface{}
type Key interface{}

type entry struct {
	key   interface{}
	value interface{}
}

//type itemStats struct {
//}

type LinkedNode struct {
	prev *LinkedNode
	next *LinkedNode
	// use to judge whether this node belongs to the list
	list  *List
	value Any
}

type List struct {
	head   *LinkedNode
	tail   *LinkedNode
	length int
}

// NewList A List Constructor
func NewList() *List {
	var l = new(List)
	l.tail = new(LinkedNode)
	l.head = new(LinkedNode)
	l.head.next = l.tail
	l.tail.prev = l.head
	l.head.list = l
	l.tail.list = l
	return l
}

// GetTail return dummy_tail prev node
func (l *List) GetTail() *LinkedNode {
	return l.tail.prev
}

func (l *List) Len() int {
	return l.length
}

// IsExist judge whether a node exists in list
func (l *List) IsExist(node *LinkedNode) bool {
	if node.list.head == l.head && node.list.tail == l.tail {
		return true
	}
	return false
}

// Add Value into the front of List and return the node insert into list
func (l *List) Add(Value any) *LinkedNode {
	if l == nil {
		l.head = new(LinkedNode)
		l.tail = new(LinkedNode)
		l.head.next = l.tail
		l.tail.prev = l.head
	}
	node := &LinkedNode{value: Value, list: l}
	realHead := l.head.next
	l.head.next = node
	node.prev = l.head
	realHead.prev = node
	node.next = realHead
	l.length++
	return node
}

// MoveToFront will move the specific node to the front of the list
func (l *List) MoveToFront(node *LinkedNode) {
	if !l.IsExist(node) {
		panic("this node is not exist in the list")
	}
	if l.head.next == node {
		return
	}
	node.prev = node.next
	node.next.prev = node.prev
	node.next = nil
	node.prev = nil
	l.Add(node.value)
}
func (l *List) DeleteNode(node *LinkedNode) {
	if node.list != l {
		panic("delete node which is not exist in list")
	}
	node.prev.next = node.next
	node.next.prev = node.prev
	node.prev = nil
	node.next = nil
	l.length--
}

type Cache struct {
	// MaxEntries is the max entries in the list.zero means no limit
	MaxEntries int

	list *List

	cache map[Key]*LinkedNode
	// call when element be evicted from Cache
	onEvicted func(key Key, value Any)
}

func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		list:       NewList(),
		cache:      make(map[Key]*LinkedNode),
	}
}
func (c *Cache) SetEvictedFunc(f func(key Key, value Any)) {
	c.onEvicted = f
}

// Add when you add entries into Cache,First try to get node from c.cache
// if node has been existed in c.cache,just update this node.
// if node has been not existed in c.cache which represent the entry has not been inserted into c.cache
// try to add the entry into c.cache and c.list,after
// this function will return a variable type
func (c *Cache) Add(key, value Any) bool {
	if c == nil {
		c.list = NewList()
		c.cache = make(map[Key]*LinkedNode)
	}
	if node, ok := c.cache[key]; ok {
		c.list.MoveToFront(node)
		// assert value in list is a type of entry
		node.value.(*entry).value = value
		return true
	}
	// update c.list
	node := c.list.Add(&entry{key: key, value: value})
	// update cache
	c.cache[key] = node
	if c.MaxEntries != 0 && c.list.Len() > c.MaxEntries {
		c.DeleteOldest()
	}
	return false
}

func (c *Cache) DeleteOldest() {
	if c == nil {
		c.list = NewList()
		c.cache = make(map[Key]*LinkedNode)
		return
	}
	if c.list.Len() == 0 {
		return
	}
	c.onEvicted(c.list.GetTail().value.(*entry).key, c.list.GetTail().value.(*entry).value)
	c.list.DeleteNode(c.list.GetTail())
	delete(c.cache, c.list.GetTail().value.(*entry).key)
}

func (c *Cache) Get(key Key) (value Any, ok bool) {
	if c == nil {
		c.list = NewList()
		c.cache = make(map[Key]*LinkedNode)
		return nil, false
	}
	if item, ok := c.cache[key]; ok {
		// Move the node in front of the list
		c.list.MoveToFront(item)
		return item.value.(*entry).value, true
	}
	return nil, false
}

func (c *Cache) Remove(key Key) (ok bool) {
	if c == nil {
		c.list = NewList()
		c.cache = make(map[Key]*LinkedNode)
		return false
	}
	node := c.cache[key]
	c.onEvicted(key, node.value.(*entry).value)
	c.list.DeleteNode(node)
	delete(c.cache, key)
	return true
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.list.Len()
}
