package main

import (
	"fmt"
	"math/rand"
	"time"
)

type User struct {
	ID   int
	Name string
	Age  int
}
type Node struct {
	index int
	data  *User
	left  *Node
	right *Node
}
type BinaryTree struct {
	root *Node
}

func (t *BinaryTree) insert(user *User) *BinaryTree {
	if t.root == nil {
		t.root = &Node{
			data: user,
		}
		return t
	}
	t.root.insert(user)

	return t
}
func (n *Node) insert(user *User) {
	ptr := n

	for ptr != nil {
		if user.ID > ptr.data.ID {
			if ptr.right == nil {
				ptr.right = &Node{
					data: user,
				}
				return
			} else {
				ptr = ptr.right
			}
		} else {
			if ptr.left == nil {
				ptr.left = &Node{
					data: user,
				}
				return
			} else {
				ptr = ptr.left
			}
		}
	}
}
func (t *BinaryTree) search(key int) *User {
	if t.root == nil {
		return nil
	}
	return t.root.search(key)
}
func (n *Node) search(key int) *User {
	ptr := n
	for ptr != nil {
		if key < ptr.data.ID {
			ptr = ptr.left
		} else {
			ptr = ptr.right
		}
		if ptr != nil {
			if ptr.data.ID == key {
				return ptr.data
			}
		}
	}
	return nil
}
func generateData(n int) *BinaryTree {
	rand.Seed(time.Now().UnixNano())
	bt := &BinaryTree{}
	for i := 0; i < n; i++ {
		val := rand.Intn(100)
		bt.insert(&User{
			ID:   val,
			Name: fmt.Sprintf("User%d", val),
			Age:  rand.Intn(50) + 20,
		})
	}
	return bt
}
func main() {
	bt := generateData(50)
	user := bt.search(30)
	if user != nil {
		fmt.Printf("Найден пользователь: %+v\n", user)
	} else {
		fmt.Println("Пользователь не найден")
	}
}
