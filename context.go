package main

import "container/list"

type OperatorTree struct {
	branches      map[rune]*OperatorTree
	childOpCount  int
	operatorToken int
}

func NewOperatorTree() *OperatorTree {
	m := make(map[rune]*OperatorTree)
	o := new(OperatorTree)
	o.branches = m
	o.childOpCount = 0
	o.operatorToken = -1
	return o
}

func (tree *OperatorTree) AddOperator(ra []rune, token int) bool {
	if len(ra) == 0 {
		if tree.operatorToken == -1 {
			tree.operatorToken = token
			return true
		} else {
			return false
		}
	}
	c := ra[0]
	branch, ok := tree.branches[c]
	if !ok {
		branch = NewOperatorTree()
		tree.branches[c] = branch
	}
	success := branch.AddOperator(ra[1:], token)
	if success {
		tree.childOpCount++
	}
	return success
}

func (tree *OperatorTree) PossibleCount(ra []rune) (int, *OperatorTree) {
	count, subtree := tree.PossibleChildCount(ra)
	if subtree.operatorToken != -1 {
		count++
	}
	return count, subtree
}

func (tree *OperatorTree) PossibleChildCount(ra []rune) (int, *OperatorTree) {
	if len(ra) == 0 {
		if tree.operatorToken != -1 {
			return 1, tree
		}
		return tree.childOpCount, tree
	}
	c := ra[0]
	branch, ok := tree.branches[c]
	if ok {
		return branch.PossibleChildCount(ra[1:])
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
	levels  *list.List
	opToken int
}

func NewOpContext() *OpContext {
	l := list.New()
	t := 10
	r := new(OpContext)
	r.levels = l
	r.opToken = t
	return r
}

func (ctx *OpContext) AddOperator(op []rune, precedence int) bool {
	level := ctx.levels.Front()
	for i := 0; i < precedence; i++ {
		if level == nil {
			level = ctx.levels.PushBack(NewOperatorTree())
		}
		level = level.Next()
	}
	if level == nil {
		level = ctx.levels.PushBack(NewOperatorTree())
	}
	tree := level.Value.(*OperatorTree)
	success := tree.AddOperator(op, ctx.opToken)
	if success {
		ctx.opToken++
	}
	return success
}

func (ctx *OpContext) PossibleCount(ra []rune) int {
	element := ctx.levels.Front()
	tree := element.Value.(*OperatorTree)
	total, _ := tree.PossibleCount(ra)
	for i := 1; i < ctx.levels.Len(); i++ {
		element = element.Next()
		tree := element.Value.(*OperatorTree)
		inc, _ := tree.PossibleCount(ra)
		total += inc
	}
	return total
}

func (ctx *OpContext) OperatorExists(ra []rune) bool {
	element := ctx.levels.Front()
	tree := element.Value.(*OperatorTree)
	exists := tree.OperatorExists(ra)
	for !exists {
		element = element.Next()
		tree := element.Value.(*OperatorTree)
		exists = exists || tree.OperatorExists(ra)
	}
	return exists
}