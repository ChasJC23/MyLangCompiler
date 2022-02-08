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
	r.opToken = INIT_TOKEN
	r.opTree.AddOperator_rune(EOF_RUNE, EOF_TOKEN)
	r.opTree.AddOperator_rune(NEWLINE_RUNE, NEWLINE_TOKEN)
	return r
}

func (ctx *OpContext) AddOperator(op []rune, preclvl *PrecedenceLevel) bool {
	token := ctx.opToken
	newToken := ctx.opTree.GetToken(op)
	if newToken != NIL_TOKEN {
		token = newToken
	}
	success := ctx.opTree.AddOperator(op, token)
	if success && token == ctx.opToken {
		ctx.opToken++
	}
	setbit(preclvl.operators, token)
	return success
}

func (ctx *OpContext) RenameOperator(oldname []rune, newname []rune) {
	oldbranch := ctx.opTree.GetBranch(oldname)
	token := oldbranch.operatorToken
	oldbranch.operatorToken = NIL_TOKEN
	newbranch := ctx.opTree.GetBranch(newname)
	newbranch.operatorToken = token
}
