package main

import (
	"fmt"
	"hash/crc32"
	"time"
)

type HashFunc func([]byte) uint32

type HashMaper interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

type HashMap struct {
	hm   HashMaper
	hash HashFunc
}

func (m *HashMap) Set(key string, value interface{}) {
	m.hm.Set(key, value)
}

func (m *HashMap) Get(key string) (interface{}, bool) {
	return m.hm.Get(key)
}

type ListMap struct {
	hash HashFunc
	data []*Node
}

type Node struct {
	key   string
	value interface{}
	next  *Node
}

func (l ListMap) CalcHash(key string) uint32 {
	hash := l.hash([]byte(key))
	return hash % uint32(len(l.data))
}

func (l ListMap) Set(key string, value interface{}) {
	if len(l.data) == 0 {
		return
	}
	hashVal := l.CalcHash(key)

	cell := l.data[hashVal]

	if cell == nil {
		l.data[hashVal] = &Node{key, value, nil}
		return
	}
	if cell.key == key {
		cell.value = value
		return
	}
	for cell.next != nil {
		if cell.next.key == key {
			cell.next.value = value
			return
		}
		cell = cell.next
	}
	cell.next = &Node{key, value, nil}
}

func (l ListMap) Get(key string) (interface{}, bool) {
	if len(l.data) == 0 {
		return nil, false
	}

	hash := l.CalcHash(key)
	cell := l.data[hash]

	for cell != nil {
		if cell.key == key {
			return cell.value, true
		}
		cell = cell.next
	}
	return nil, false
}

type SliceMap struct {
	hash HashFunc
	data [][]*SliceNode
}

type SliceNode struct {
	key   string
	value interface{}
}

func (m SliceMap) Set(key string, value interface{}) {
	if len(m.data) == 0 {
		return
	}

	hashVal := m.CalcHash(key)

	cell := m.data[hashVal]

	if len(cell) == 0 {
		m.data[hashVal] = append(m.data[hashVal], &SliceNode{key, value})
		return
	}

	for i := 0; i < len(cell); i++ {
		if cell[i].key == key {
			cell[i].value = value
			return
		}
	}
	m.data[hashVal] = append(m.data[hashVal], &SliceNode{key, value})
}

func (m SliceMap) Get(key string) (interface{}, bool) {
	if len(m.data) == 0 {
		return nil, false
	}

	hashVal := m.CalcHash(key)

	cell := m.data[hashVal]

	for i := 0; i < len(cell); i++ {
		if cell[i].key == key {
			return cell[i].value, true
		}
	}
	return nil, false
}

func (m SliceMap) CalcHash(key string) uint32 {
	hash := m.hash([]byte(key))
	return hash % uint32(len(m.data))
}

func WithHashFunc(hash func([]byte) uint32) func(listMap *HashMap) {
	return func(m *HashMap) {
		m.hash = hash
	}
}

func NewHashMapSlice(count int, options ...func(*HashMap)) HashMaper {
	m := &HashMap{}
	for _, option := range options {
		option(m)
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	m.hm = SliceMap{
		hash: m.hash,
		data: make([][]*SliceNode, count),
	}
	return m
}
func NewHashMapList(count int, options ...func(*HashMap)) HashMaper {
	m := &HashMap{}
	for _, option := range options {
		option(m)
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	m.hm = ListMap{
		hash: m.hash,
		data: make([]*Node, count),
	}
	return m
}
func MeassureTime(f func()) time.Duration {
	start := time.Now()
	f()
	since := time.Since(start)
	return since
}
func main() {
	time := MeassureTime(TestSlice16)
	fmt.Println(time)
	time = MeassureTime(TestSlice1000)
	fmt.Println(time)
	time = MeassureTime(TestList16)
	fmt.Println(time)
	time = MeassureTime(TestList1000)
	fmt.Println(time)
}
func TestList16() {
	m := NewHashMapList(16)
	for i := 0; i < 16; i++ {
		m.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	for i := 0; i < 16; i++ {
		value, ok := m.Get(fmt.Sprintf("key%d", i))
		if !ok {
			fmt.Printf("Expected key to exist in the HashMap")
		}
		if value != fmt.Sprintf("value%d", i) {
			fmt.Printf("Expected value to be 'value%d', got '%v'", i, value)
		}
	}
}
func TestList1000() {
	m := NewHashMapList(1000)
	for i := 0; i < 1000; i++ {
		m.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	for i := 0; i < 1000; i++ {
		value, ok := m.Get(fmt.Sprintf("key%d", i))
		if !ok {
			fmt.Printf("Expected value to be 'value%d', got '%v'\n", i, value)
		}
	}
}

func TestSlice16() {
	m := NewHashMapSlice(16)
	for i := 0; i < 16; i++ {
		m.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	for i := 0; i < 16; i++ {
		value, ok := m.Get(fmt.Sprintf("key%d", i))
		if !ok {
			fmt.Printf("Expected key to exist in the HashMap\n")
		}
		if value != fmt.Sprintf("value%d", i) {
			fmt.Printf("Expected value to be 'value%d', got '%v'\n", i, value)
		}
	}
}
func TestSlice1000() {
	m := NewHashMapSlice(1000)
	for i := 0; i < 1000; i++ {
		m.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	for i := 0; i < 1000; i++ {
		value, ok := m.Get(fmt.Sprintf("key%d", i))
		if !ok {
			fmt.Printf("Expected key to exist in the HashMap\n")
		}
		if value != fmt.Sprintf("value%d", i) {
			fmt.Printf("Expected value to be 'value%d', got '%v'\n", i, value)
		}
	}
}
