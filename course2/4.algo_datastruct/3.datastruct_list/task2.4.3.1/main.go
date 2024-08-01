package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type Node struct {
	data *Commit
	prev *Node
	next *Node
}
type DoubleLinkedList struct {
	head *Node // first element in list
	tail *Node // last element in list
	curr *Node // current element. Changes when Prev or Next methods are used
	len  int   // elements quantity in list
}
type LinkedLister interface {
	LoadData(path string) error
	Init(c []Commit)
	Len() int
	SetCurrent(n int) error
	Current() *Node
	Next() *Node
	Prev() *Node
	Insert(n int, c Commit) error
	Push(c Commit) error
	Delete(n int) error
	DeleteCurrent() error
	Index() (int, error)
	GetByIndex(n int) (*Node, error)
	Pop() *Node
	Shift() *Node
	SearchUUID(uuID string) *Node
	Search(message string) *Node
	Reverse() *DoubleLinkedList
}

// LoadData loads data from a JSON file at the given path into the list.
func (d *DoubleLinkedList) LoadData(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var commits []Commit
	if err := json.Unmarshal(bytes, &commits); err != nil {
		return err
	}

	QuickSort(commits)
	d.Init(commits)

	return nil
}

func (d *DoubleLinkedList) Init(c []Commit) {
	if len(c) < 1 {
		return
	}

	d.len = len(c)
	currentNode := &Node{
		data: &c[0],
	}

	d.head = currentNode

	for i := 1; i < len(c); i++ {
		newNode := Node{
			data: &c[i],
			prev: currentNode,
		}
		currentNode.next = &newNode
		currentNode = &newNode
	}

	d.tail = currentNode
}

// Len getting list len
func (d *DoubleLinkedList) Len() int {
	return d.len
}

func (d *DoubleLinkedList) SetCurrent(n int) {
	newCurrent := d.head

	for i := 0; i < n; i++ {
		newCurrent = newCurrent.next
	}

	d.curr = newCurrent
}

// Current getting current element
func (d *DoubleLinkedList) Current() *Node {
	return d.curr
}

// Next getting next element
func (d *DoubleLinkedList) Next() *Node {
	if d.Len() == 0 {
		return nil
	}

	if d.curr == nil {
		d.curr = d.head
	} else {
		d.curr = d.curr.next
	}

	return d.curr
}

func (d *DoubleLinkedList) Prev() *Node {
	if d.Len() == 0 {
		return nil
	}

	if d.curr == nil {
		d.curr = d.tail
	} else {
		d.curr = d.curr.prev
	}

	return d.curr
}

// Insert element after n-th element
// Insert inserts a new node with commit c at position n.
func (d *DoubleLinkedList) Insert(n int, c Commit) error {
	if n < 0 || n > d.len {
		return errors.New("index out of bounds")
	}

	newNode := &Node{data: &c}
	if n == 0 {
		if d.head == nil {
			d.head = newNode
			d.tail = newNode
		} else {
			newNode.next = d.head
			d.head.prev = newNode
			d.head = newNode
		}
	} else if n == d.len {
		d.tail.next = newNode
		newNode.prev = d.tail
		d.tail = newNode
	} else {
		current := d.head

		for i := 0; i < n; i++ {
			current = current.next
		}

		newNode.next = current
		newNode.prev = current.prev
		current.prev.next = newNode
		current.prev = newNode
	}

	d.len++
	return nil
}

func (d *DoubleLinkedList) Push(c Commit) error {
	newNode := &Node{
		data: &c,
	}

	if d.Len() == 0 {
		d.head = newNode
		d.tail = newNode
	} else {
		newNode.prev = d.tail
		d.tail.next = newNode
		d.tail = newNode
	}

	d.len += 1
	return nil
}

// Delete n-th element
func (d *DoubleLinkedList) Delete(n int) error {
	if d.Len() == 0 {
		return fmt.Errorf("list is empty")
	}

	if d.Len() <= n || n < 0 {
		return fmt.Errorf("number is out of length")
	}

	if (d.Len() - 1) == n {
		d.Pop()
		return nil
	} else if n == 0 {
		d.Shift()
		return nil
	} else {
		currentNode, err := d.GetByIndex(n)

		if err != nil {
			return err
		}

		currentNode.next.prev = currentNode.prev
		currentNode.prev.next = currentNode.next
		if currentNode == d.curr {
			d.Next()
		}
	}

	d.len -= 1
	return nil
}

// DeleteCurrent удаление текущего элемента
func (d *DoubleLinkedList) DeleteCurrent() error {
	if d.curr == nil {
		return fmt.Errorf("no current element")
	}

	if d.Len() == 1 {
		d.curr = nil
		d.head = nil
		d.tail = nil
	} else if d.curr.next == nil {
		d.tail = d.curr.prev
		d.curr.prev.next = nil
		d.curr = d.Prev()
	} else if d.curr.prev == nil {
		d.head = d.curr.next
		d.curr.next.prev = nil
		d.curr = d.Next()
	} else {
		d.curr.next.prev = d.curr.prev
		d.curr.prev.next = d.curr.next
		d.curr = d.Next()
	}

	d.len -= 1
	return nil
}

// Index get current element index
func (d *DoubleLinkedList) Index() (int, error) {
	if d.curr == nil {
		return -1, fmt.Errorf("no current Node set")
	}

	currentNode := d.curr
	index := 0

	for currentNode.prev != nil {
		currentNode = currentNode.prev
		index += 1
	}

	return index, nil
}

func (d *DoubleLinkedList) GetByIndex(n int) (*Node, error) {
	if d.Len() == 0 {
		return nil, fmt.Errorf("link is empty")
	}

	if d.Len() < n || n < 0 {
		return nil, fmt.Errorf("index out of range")
	}

	currentNode := d.head

	for i := 0; i < n; i++ {
		currentNode = currentNode.next
	}

	return currentNode, nil
}

// Pop
func (d *DoubleLinkedList) Pop() *Node {
	if d.tail == nil {
		return nil
	}

	ptr := d.tail
	if d.head == d.tail {
		d.tail = nil
		d.head = nil
		d.curr = nil
	} else {
		d.tail = d.tail.prev
		d.tail.next = nil
		if ptr == d.curr {
			d.Prev()
		}
	}

	d.len -= 1
	return ptr
}

// Shift операция shift
func (d *DoubleLinkedList) Shift() *Node {
	if d.Len() == 0 {
		return nil
	}

	shiftedNode := d.head
	if d.Len() == 1 {
		d.tail = nil
		d.head = nil
		d.curr = nil
	} else {
		d.head = d.head.next
		d.head.prev = nil

		if shiftedNode == d.curr {
			d.Next()
		}
	}

	d.len -= 1
	return shiftedNode
}

// SearchUUID
func (d *DoubleLinkedList) SearchUUID(uuID string) *Node {
	if d.Len() == 0 {
		return nil
	}

	currentNode := d.head
	for currentNode != nil {
		if currentNode.data.UUID == uuID {
			return currentNode
		}

		currentNode = currentNode.next
	}

	return nil
}

func (d *DoubleLinkedList) Search(message string) *Node {
	if d.Len() == 0 {
		return nil
	}

	currentNode := d.head
	for currentNode != nil {
		if currentNode.data.Message == message {
			return currentNode
		}

		currentNode = currentNode.next
	}

	return nil
}

// Reverse the link
func (d *DoubleLinkedList) Reverse() *DoubleLinkedList {
	if d.Len() < 2 {
		return d
	}

	currentNode := d.head
	for currentNode.next != nil {
		currentNode.next, currentNode.prev = currentNode.prev, currentNode.next
		currentNode = currentNode.prev
	}

	currentNode.next, currentNode.prev = currentNode.prev, currentNode.next
	d.head, d.tail = d.tail, d.head

	return d
}

type Commit struct {
	Message string `json:"message"`
	UUID    string `json:"uuid"`
	Date    string `json:"date"`
}

func quickSort(arr []Commit, low, high int) []Commit {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func QuickSort(arr []Commit) []Commit {
	return quickSort(arr, 0, len(arr)-1)
}

func partition(arr []Commit, low, high int) ([]Commit, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		data1, _ := time.Parse("2006-01-02", arr[j].Date)
		date2, _ := time.Parse("2006-01-02", pivot.Date)
		if data1.Before(date2) {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func GenerateData(numCommits int) []Commit {
	var commits []Commit
	gofakeit.Seed(0) // Initialize the random seed
	// Define how many commits you want to generate
	for i := 0; i < numCommits; i++ {
		commit := Commit{
			Message: gofakeit.Sentence(5),
			// Generate a random sentence with 5 words
			UUID: gofakeit.UUID(),
			// Generate a random UUID
			Date: gofakeit.DateRange(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)).Format("2006-01-02"),
			// Generate a random date between 2020 and 2022
		}
		commits = append(commits, commit)
	}
	return commits
}
