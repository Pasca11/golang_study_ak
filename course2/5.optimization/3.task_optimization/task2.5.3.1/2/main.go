package main

import "sort"

//https://leetcode.com/problems/sort-the-students-by-their-kth-score/

func main() {
	sco
}

func sortTheStudents(score [][]int, k int) [][]int {
	sort.Slice(score, func(i, j int) bool {
		return score[k][i] < score[k][j]
	})
	return score
}
