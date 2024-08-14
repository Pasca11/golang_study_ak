package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	list := &ListNode{0, &ListNode{3, &ListNode{1, &ListNode{0, &ListNode{4, &ListNode{5, &ListNode{2, &ListNode{0, nil}}}}}}}}
	root := mergeNodes(list)
	for root != nil {
		fmt.Println(root.Val)
		root = root.Next
	}
}

func mergeNodes(head *ListNode) *ListNode {
	var root *ListNode
	var rootPtr *ListNode

	ptr := head
	zPtr := head.Next
	for zPtr != nil {
		if zPtr.Val == 0 && ptr != zPtr {
			sum := 0
			for ptr != zPtr {
				sum += ptr.Val
				ptr = ptr.Next
			}
			if root != nil {
				rootPtr.Next = &ListNode{}
				rootPtr.Next.Val = sum
				rootPtr = rootPtr.Next
			} else {
				root = &ListNode{}
				root.Val = sum
				rootPtr = root
			}
		}
		zPtr = zPtr.Next
	}
	return root
}
