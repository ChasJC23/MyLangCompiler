package main

import "container/list"

type OpContext struct {
	// the tree used to find the token for a given operator
	opTree *OperatorTree
	// the list of all precedence levels used in a particular parsing session
	precedenceList *list.List
	opToken        int
}

func NewOpContext() *OpContext {
	r := new(OpContext)
	r.opTree = NewOperatorTree()
	r.precedenceList = list.New()
	r.opToken = INIT_TOKEN
	r.opTree.AddOperatorRune(EOF_RUNE, EOF_TOKEN, 0)
	r.opTree.AddOperatorRune(NEWLINE_RUNE, NEWLINE_TOKEN, NEWLINE_FLAG)
	return r
}

// AddOperator adds a new operator to the current context and the given precedence level.
// It returns a boolean representing whether the addition was successful or not.
func (ctx *OpContext) AddOperator(symbols []string, precedenceLevel *PrecedenceLevel, codeBlockArguments uint, argumentCount int) bool {
	op := symbols[0]
	symbols = symbols[1:]
	token, success := ctx.addOperatorToken([]rune(op))
	subsequentSymbols := make([]int, len(symbols))
	for i, symbol := range symbols {
		newToken, newSuccess := ctx.addOperatorToken([]rune(symbol))
		success = success && newSuccess
		subsequentSymbols[i] = newToken
	}
	if success {
		properties := new(OpProp)
		properties.argumentCount = argumentCount
		properties.codeBlockArguments = codeBlockArguments
		properties.subsequentSymbols = subsequentSymbols
		precedenceLevel.operators[token] = properties
	}
	return success
}

// addOperatorToken adds a new operator symbol to the operator tree.
// Returns the generated token for this new symbol.
func (ctx *OpContext) addOperatorToken(op []rune) (int, bool) {
	token := ctx.opToken
	newToken := ctx.opTree.GetToken(op)
	if newToken != NIL_TOKEN {
		token = newToken
	}
	success := ctx.opTree.AddOperator(op, token, 0)
	if success && token == ctx.opToken {
		ctx.opToken++
	}
	return token, success
}

func (ctx *OpContext) AddOperatorAt(symbols []string, precedence int, codeBlockArguments uint, argumentCount int) bool {
	precedenceLevel := ctx.precedenceList.Front()
	for i := 0; i < precedence; i++ {
		precedenceLevel = precedenceLevel.Next()
	}
	return ctx.AddOperator(symbols, precedenceLevel.Value.(*PrecedenceLevel), codeBlockArguments, argumentCount)
}

func (ctx *OpContext) AddOperatorToLowest(symbols []string, codeBlockArguments uint, argumentCount int) bool {
	precedenceLevel := ctx.precedenceList.Front()
	return ctx.AddOperator(symbols, precedenceLevel.Value.(*PrecedenceLevel), codeBlockArguments, argumentCount)
}

func (ctx *OpContext) AddOperatorToHighest(symbols []string, codeBlockArguments uint, argumentCount int) bool {
	precedenceLevel := ctx.precedenceList.Back()
	return ctx.AddOperator(symbols, precedenceLevel.Value.(*PrecedenceLevel), codeBlockArguments, argumentCount)
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

func (ctx *OpContext) RenameOperator(oldName []rune, newName []rune) {
	oldBranch := ctx.opTree.GetBranch(oldName)
	token := oldBranch.operatorToken
	oldBranch.operatorToken = NIL_TOKEN
	newBranch := ctx.opTree.GetBranch(newName)
	newBranch.operatorToken = token
}

func (ctx *OpContext) AddLowestPrecedenceLevel(precedenceLevel *PrecedenceLevel) {
	ctx.precedenceList.PushFront(precedenceLevel)
}

func (ctx *OpContext) AddHighestPrecedenceLevel(precedenceLevel *PrecedenceLevel) {
	ctx.precedenceList.PushBack(precedenceLevel)
}

func (ctx *OpContext) AddLowerPrecedenceLevel(level *PrecedenceLevel, mark *list.Element) {
	ctx.precedenceList.InsertBefore(level, mark)
}

func (ctx *OpContext) AddHigherPrecedenceLevel(level *PrecedenceLevel, mark *list.Element) {
	ctx.precedenceList.InsertAfter(level, mark)
}
