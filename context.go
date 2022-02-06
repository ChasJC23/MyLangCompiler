package main

import "container/list"

type OpContext struct {
	// the tree used to find the token for a given operator
	opTree *OperatorTree
	// the list of all precedence levels used in a particular parsing session
	precList *list.List
	opToken  int
}

func NewOpContext() *OpContext {
	r := new(OpContext)
	r.opTree = NewOperatorTree()
	r.opToken = 10
	return r
}

func (ctx *OpContext) AddOperator(op []rune, preclvl *PrecedenceLevel) bool {
	token := ctx.opToken
	newToken := ctx.opTree.GetToken(op)
	if newToken != -1 {
		token = newToken
	}
	success := ctx.opTree.AddOperator(op, token)
	if success && token == ctx.opToken {
		ctx.opToken++
	}
	setbit(preclvl.operators, token)
	return success
}
