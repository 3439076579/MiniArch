package main

import "MiniArch/giu"

//func Traversal(board [][]byte, word string, height int, pattern string) bool {
//
//	if len(word) == height {
//		if pattern == word {
//			return true
//		}
//		return false
//	}
//	for i := 0; i < len(board); i++ {
//		for j := 0; i < len(board[i]); j++ {
//			if word[height] == board[i][j] {
//
//			}
//
//		}
//	}
//
//	return false
//
//}
//
//func exist(board [][]byte, word string) bool {
//
//	return Traversal(board, word, 0, "")
//
//}

func main() {

	engine := giu.New()

	group := engine.Group("/v")

	group.Group("/c")

}
