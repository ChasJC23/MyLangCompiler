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
	r.opTree.AddOperatorRune(EOF_RUNE, EOF_TOKEN, 0)
	r.opTree.AddOperatorRune(NEWLINE_RUNE, NEWLINE_TOKEN, NEWLINE_FLAG)
	return r
}

func (ctx *OpContext) AddOperator(op []rune, preclvl *PrecedenceLevel, properties *OpProp) bool {
	token := ctx.opToken
	newToken := ctx.opTree.GetToken(op)
	if newToken != NIL_TOKEN {
		token = newToken
	}
	success := ctx.opTree.AddOperator(op, token, 0)
	if success && token == ctx.opToken {
		ctx.opToken++
	}
	preclvl.operators[token] = properties
	return success
}

func (ctx *OpContext) AddControlOperator(op []rune, flags uint) bool {
	success := ctx.opTree.AddOperator(op, ctx.opToken, flags)
	if success {
		ctx.opToken++
	}
	return success
}

func (ctx *OpContext) AddFixedTokenOperator(op []rune, token int, flags uint) bool {
	return ctx.opTree.AddOperator(op, token, flags)
}

func (ctx *OpContext) RenameOperator(oldname []rune, newname []rune) {
	oldbranch := ctx.opTree.GetBranch(oldname)
	token := oldbranch.operatorToken
	oldbranch.operatorToken = NIL_TOKEN
	newbranch := ctx.opTree.GetBranch(newname)
	newbranch.operatorToken = token
}
