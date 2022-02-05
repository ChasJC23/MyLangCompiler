package main

type OpContext struct {
	tree    *OperatorTree
	opToken Token
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
