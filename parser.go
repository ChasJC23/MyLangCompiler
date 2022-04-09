package main

type Parser struct {
	tokeniser *Tokeniser
	opctx     *OpContext
}

func NewParser(tokeniser *Tokeniser) *Parser {
	parser := new(Parser)
	parser.tokeniser = tokeniser
	parser.opctx = tokeniser.opctx
	return parser
}

// ParseSource is responsible for parsing an entire source file
func (p *Parser) ParseSource() AST {

	result, _ := p.ParsePrecedenceLevel(p.opctx.rootPrecedence)

	// I realise this is more of an expression thing,
	// but let's include it just in case.
	if p.tokeniser.currToken != EOF_TOKEN {
		panic("Unexpected characters at end of source file")
	}

	return result
}

func (p *Parser) ParsePrecedenceLevel(precedenceLevel *PrecedenceLevel) (tree AST, parenthesized bool) {
	if precedenceLevel == nil {
		return p.ParseLeaf()
	}
	// check bitmask in precedence_level.go
	switch precedenceLevel.properties & 0b111 {
	case PREFIX:
		return p.ParsePrefix(precedenceLevel)
	case POSTFIX:
		return p.ParsePostfix(precedenceLevel)
	case INFIX_LEFT_ASSOCIATIVE:
		return p.ParseLeftAssociative(precedenceLevel)
	case INFIX_RIGHT_ASSOCIATIVE:
		return p.ParseRightAssociative(precedenceLevel)
	case IMPLIED_OPERATION | INFIX_LEFT_ASSOCIATIVE:
		return p.ParseImpliedLeftAssociative(precedenceLevel)
	case IMPLIED_OPERATION | INFIX_RIGHT_ASSOCIATIVE:
		return p.ParseImpliedRightAssociative(precedenceLevel)
	case DELIMITER:
		return p.ParseDelimiter(precedenceLevel)
	default:
		panic("invalid configuration")
	}
}

func (p *Parser) ParseDelimiter(precedenceLevel *PrecedenceLevel) (tree AST, parenthesized bool) {
	firstArg, lp := p.ParsePrecedenceLevel(precedenceLevel.defaultNext)

	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return firstArg, lp
	}

	arguments := make([]AST, 1, 2)
	arguments[0] = firstArg
	nextOpProp := opProperties
	for nextOpProp != nil {
		p.tokeniser.ReadToken()
		nextArg, _ := p.ParsePrecedenceLevel(precedenceLevel.defaultNext)
		arguments = append(arguments, nextArg)
		nextOpProp = precedenceLevel.operators[p.tokeniser.currToken]
	}
	return NewStatement(arguments, opProperties), false
}

func (p *Parser) ParseImpliedLeftAssociative(precedenceLevel *PrecedenceLevel) (tree AST, parenthesized bool) {
	lhs, lp := p.ParsePrecedenceLevel(precedenceLevel.defaultNext)

	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return lhs, lp
	}
	p.tokeniser.ReadToken()

	args := p.getInfixArguments(precedenceLevel.defaultNext, lhs, opProperties)

	expr := NewStatement(args, opProperties)

	impliedOpProp := precedenceLevel.operators[NIL_TOKEN]

	for {
		opProperties = precedenceLevel.operators[p.tokeniser.currToken]
		if opProperties == nil {
			return expr, false
		}
		p.tokeniser.ReadToken()

		lhs = args[len(args)-1]
		args = p.getInfixArguments(precedenceLevel.defaultNext, lhs, opProperties)
		expr = NewStatement([]AST{expr, NewStatement(args, opProperties)}, impliedOpProp)
	}
}

func (p *Parser) ParseImpliedRightAssociative(precedenceLevel *PrecedenceLevel) (tree AST, parenthesized bool) {
	lhs, lp := p.ParsePrecedenceLevel(precedenceLevel.defaultNext)

	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return lhs, lp
	}
	p.tokeniser.ReadToken()

	args := p.getInfixArguments(precedenceLevel, lhs, opProperties)
	rhs := args[len(args)-1]

	binRight, ok := rhs.(*Statement)

	impliedOpProp := precedenceLevel.operators[NIL_TOKEN]

	// if rhs is an operator
	if ok {
		// if it's an operator in this precedence level
		if precedenceLevel.OperatorExists(binRight.properties) {
			binRightLeft, ok := binRight.terms[0].(*Statement)
			// if the right-hand side is the implied operation (this "ok &&" might cause some unexpected behaviour for weirdly structured trees)
			if ok && binRight.properties == impliedOpProp {
				/*
						  &
						 / \
						#   &
					   /|  / \
					  / | #  ...
					 /  |/ \ / \
					a   b  ... ...
				*/
				args[len(args)-1] = binRightLeft.terms[0]
			} else {
				/*
						&
					   / \
					  #   #
					 / \ / \
					a   b   c
				*/
				args[len(args)-1] = binRight.terms[0]
			}
			return NewStatement([]AST{NewStatement(args, opProperties), binRight}, impliedOpProp), false
		}
	}
	/*
		  #
		 / \
		a   b
	*/
	return NewStatement([]AST{lhs, rhs}, opProperties), false
}

func (p *Parser) ParseLeftAssociative(precedenceLevel *PrecedenceLevel) (tree AST, parenthesized bool) {
	lhs, lp := p.ParsePrecedenceLevel(precedenceLevel.defaultNext)
	for {
		opProperties := precedenceLevel.operators[p.tokeniser.currToken]

		// if this operator isn't defined:
		if opProperties == nil {

			// we need to check if an implied operation exists for this precedence level
			nilOpProp := precedenceLevel.operators[NIL_TOKEN]

			// we know the next symbol isn't in this precedence level,
			// but still check if it's a control token in case of higher precedence operators.
			if nilOpProp == nil || p.tokeniser.currToken > NIL_TOKEN || p.tokeniser.currToken == EOF_TOKEN {
				return lhs, lp
			} else {
				opProperties = nilOpProp
			}
		} else {
			p.tokeniser.ReadToken()
		}

		// we've parsed the first argument, now we use the operator properties to deduce subsequent symbols to expect
		args := p.getInfixArguments(precedenceLevel.defaultNext, lhs, opProperties)

		lhs = NewStatement(args, opProperties)
	}
}

func (p *Parser) ParseRightAssociative(precedenceLevel *PrecedenceLevel) (tree AST, parenthesized bool) {
	lhs, lp := p.ParsePrecedenceLevel(precedenceLevel.defaultNext)
	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	// same logic as left associative, just a bit of recursion to get the associativity right
	// private function might be useful, quite a lot of redundancy here
	if opProperties == nil {
		nilOpProp := precedenceLevel.operators[NIL_TOKEN]
		if nilOpProp == nil || p.tokeniser.currToken > NIL_TOKEN || p.tokeniser.currToken == EOF_TOKEN {
			return lhs, lp
		} else {
			opProperties = nilOpProp
		}
	} else {
		p.tokeniser.ReadToken()
	}
	args := p.getInfixArguments(precedenceLevel, lhs, opProperties)

	return NewStatement(args, opProperties), false
}

func (p *Parser) ParsePrefix(precedenceLevel *PrecedenceLevel) (tree AST, parenthesized bool) {
	opProperties := precedenceLevel.operators[p.tokeniser.currToken]
	if opProperties == nil {
		return p.ParsePrecedenceLevel(precedenceLevel.defaultNext)
	}
	argumentCount := opProperties.argumentCount
	argumentSlice := make([]AST, argumentCount)
	p.tokeniser.ReadToken()
	// uh... that works I guess
	for argumentIndex := 0; argumentIndex < argumentCount; argumentIndex = argumentIndex + 1 {
		var argument AST
		argument, _ = p.ParsePrefix(precedenceLevel)
		argumentSlice[argumentIndex] = argument
	}
	return NewStatement(argumentSlice, opProperties), false
}

func (p *Parser) ParsePostfix(precedenceLevel *PrecedenceLevel) (tree AST, parenthesized bool) {
	// stack based parsing
	stack := make(ASTStack, 0)

	opProperties := precedenceLevel.operators[p.tokeniser.currToken]

	// The stack should be the deciding factor for when we complete this precedence level
	for len(stack) != 1 || opProperties != nil {
		// for the symbols we don't recognise here, pass onto higher precedence parsing and add to the stack
		if opProperties == nil {
			arg, _ := p.ParsePrecedenceLevel(precedenceLevel.defaultNext)
			stack.Push(arg)
		} else
		// for the symbols we do recognise, replace top few elements of the stack with the parsed result
		{
			argCount := opProperties.argumentCount
			argumentSlice := make([]AST, argCount)
			for i := argCount - 1; i >= 0; i++ {
				argumentSlice[i] = stack.Pop()
			}
			stack.Push(NewStatement(argumentSlice, opProperties))
			p.tokeniser.ReadToken()
		}
		opProperties = precedenceLevel.operators[p.tokeniser.currToken]
	}
	if len(stack) != 1 {
		panic("unused arguments in postfix expression")
	}
	return stack[0], false
}

func (p *Parser) ParseLeaf() (tree AST, parenthesized bool) {
	var result AST
	parenthesized = false
	switch p.tokeniser.currToken {
	case INT_LITERAL:
		result = IntLiteral{p.tokeniser.intLiteral}
	case FLOAT_LITERAL:
		result = FloatLiteral{p.tokeniser.floatLiteral}
	case CHAR_LITERAL:
		result = CharLiteral{p.tokeniser.charLiteral}
	case STRING_LITERAL:
		result = StringLiteral{p.tokeniser.stringLiteral}
	case TRUE_LITERAL:
		result = BoolLiteral{true}
	case FALSE_LITERAL:
		result = BoolLiteral{false}
	case IDENTIFIER_TOKEN:
		result = Identifier{p.tokeniser.identifier}
	default:
		if p.tokeniser.currToken <= MAX_PARENS_TOKEN && (p.tokeniser.currToken+MAX_PARENS_TOKEN)&1 == OPEN_PARENS_MOD_2 {
			precedenceLevel := p.opctx.rootPrecedence
			for precedenceLevel.operators[p.tokeniser.currToken] == nil {
				precedenceLevel = precedenceLevel.defaultNext
			}
			p.tokeniser.ReadToken()
			result, _ = p.ParseDelimiter(precedenceLevel)
			parenthesized = true
			if (p.tokeniser.currToken+MAX_PARENS_TOKEN)&1 != CLOSE_PARENS_MOD_2 {
				panic("missing close parenthesis")
			}
		} else {
			panic("unrecognised token")
		}
	}
	p.tokeniser.ReadToken()
	return result, parenthesized
}

func (p *Parser) getInfixArguments(termPrecedence *PrecedenceLevel, firstArg AST, opProperties *OpProp) []AST {
	args := make([]AST, 2, len(opProperties.subsequentSymbols)+2)
	args[0] = firstArg
	args[1], _ = p.ParsePrecedenceLevel(termPrecedence)
	for i := 0; i < len(opProperties.subsequentSymbols); i++ {
		nextSymbol := opProperties.subsequentSymbols[i]
		if p.tokeniser.currToken == nextSymbol || p.tokeniser.currToken < NIL_TOKEN && nextSymbol == NIL_TOKEN {
			if nextSymbol != NIL_TOKEN {
				p.tokeniser.ReadToken()
			}
			nextArg, _ := p.ParsePrecedenceLevel(termPrecedence)
			args = append(args, nextArg)
		} else {
			panic("Unexpected symbol")
		}
	}
	return args
}
