package main

import (
	"strconv"
)

type OperatorTree struct {
	branches      map[rune]*OperatorTree
	childOpCount  int
	operatorToken int
	root          *OperatorTree
}

func NewOperatorTree() *OperatorTree {
	m := make(map[rune]*OperatorTree)
	o := new(OperatorTree)
	o.branches = m
	o.childOpCount = 0
	o.operatorToken = -1
	return o
}

func (tree *OperatorTree) ToString() string {
	result := "[" + strconv.FormatInt(int64(tree.operatorToken), 10) + "]{"
	for i, v := range tree.branches {
		result += string(i) + ":" + v.ToString() + ","
	}
	return result + "}"
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
		branch.root = tree
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
	return tree.GetToken(ra) != -1
}

func (tree *OperatorTree) GetToken(ra []rune) int {
	if len(ra) == 0 {
		return tree.operatorToken
	}
	c := ra[0]
	branch, ok := tree.branches[c]
	if ok {
		return branch.GetToken(ra[1:])
	}
	return -1
}

type OpContext struct {
	tree    *OperatorTree
	opToken int
}

func NewOpContext() *OpContext {
	r := new(OpContext)
	r.tree = NewOperatorTree()
	r.opToken = 10
	return r
}

func (ctx *OpContext) AddOperator(op []rune) bool {
	token := ctx.opToken
	newToken := ctx.tree.GetToken(op)
	if newToken != -1 {
		token = newToken
	}
	success := ctx.tree.AddOperator(op, token)
	if success && token == ctx.opToken {
		ctx.opToken++
	}
	return success
}
