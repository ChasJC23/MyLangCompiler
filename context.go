package main

type OpContext struct {
	// the tree used to find the token for a given operator
	opTree *OperatorTree
	// a slice of pointers to precedence levels each token (the index) is a member of
	precPtrs []*PrecedenceLevel
	// the list of precedence levels which may be manipulated at any time
	precList *PrecedenceList
	opToken  Token
}

func NewOpContext() *OpContext {
	r := new(OpContext)
	r.opTree = NewOperatorTree()
	r.opToken = 10
	return r
}

func (ctx *OpContext) AddOperator(op []rune) bool {
	token := ctx.opToken
	newToken := ctx.opTree.GetToken(op)
	if newToken != -1 {
		token = newToken
	}
	success := ctx.opTree.AddOperator(op, token)
	if success && token == ctx.opToken {
		ctx.opToken++
	}
	return success
}
