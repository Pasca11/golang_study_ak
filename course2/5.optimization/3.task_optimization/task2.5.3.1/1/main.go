package main

import "fmt"

//https://leetcode.com/problems/deepest-leaves-sum/description/

func main() {
	root := &TreeNode{Val: 1, Left: &TreeNode{Val: 2, Left: &TreeNode{Val: 4, Left: &TreeNode{Val: 7}}, Right: &TreeNode{Val: 5}}, Right: &TreeNode{Val: 3, Right: &TreeNode{Val: 6, Right: &TreeNode{Val: 8}}}}
	res := deepestLeavesSum(root)

	fmt.Println(res)
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func deepestLeavesSum(root *TreeNode) int {
	lvl := deep(root, 0)

	return recDeepestLeavesSum(root, 0, lvl)
}

func recDeepestLeavesSum(root *TreeNode, lvl, maxLvl int) int {
	if root == nil {
		return 0
	}
	if root.Left == nil && root.Right == nil {
		if lvl == maxLvl {
			return root.Val
		} else {
			return 0
		}
	}
	return recDeepestLeavesSum(root.Left, lvl+1, maxLvl) + recDeepestLeavesSum(root.Right, lvl+1, maxLvl)
}

func deep(root *TreeNode, lvl int) int {
	if root.Left == nil && root.Right == nil {
		return lvl
	}
	lvlLeft, lvlRight := 0, 0
	if root.Left != nil {
		lvlLeft = deep(root.Left, lvl+1)
	}
	if root.Right != nil {
		lvlRight = deep(root.Right, lvl+1)
	}
	if lvlLeft > lvlRight {
		return lvlLeft
	} else {
		return lvlRight
	}
}
