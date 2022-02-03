package main

import "container/list"

type OperatorTree struct {
	branches      map[rune]OperatorTree
	childOpCount  int
	operatorToken int
}

func (tree *OperatorTree) PossibleCount(ra []rune) (int, *OperatorTree) {
	if len(ra) == 0 {
		if tree.operatorToken != -1 {
			return 1, tree
		}
		return tree.childOpCount, tree
	}
	c := ra[0]
	branch, ok := tree.branches[c]
	if ok {
		return branch.PossibleCount(ra[1:])
	} else {
		return 0, tree
	}
}
func (tree *OperatorTree) OperatorExists(ra []rune) bool {
	c := ra[0]
	branch, ok := tree.branches[c]
	if ok {
		if len(ra) == 1 {
			return tree.operatorToken != -1
		} else {
			return branch.OperatorExists(ra[1:])
		}
	}
	return false
}

type OpContext struct {
	levels list.List
}

func (ctx *OpContext) PossibleCount(ra []rune) int {
	element := ctx.levels.Front()
	tree := element.Value.(OperatorTree)
	total, _ := tree.PossibleCount(ra)
	for i := 1; i < ctx.levels.Len(); i++ {
		element = element.Next()
		tree := element.Value.(OperatorTree)
		inc, _ := tree.PossibleCount(ra)
		total += inc
	}
	return total
}

func (ctx *OpContext) OperatorExists(ra []rune) bool {
	element := ctx.levels.Front()
	tree := element.Value.(OperatorTree)
	exists := tree.OperatorExists(ra)
	for !exists {
		element = element.Next()
		tree := element.Value.(OperatorTree)
		exists = exists || tree.OperatorExists(ra)
	}
	return exists
}
