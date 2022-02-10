package main

type AST interface {
}

type CodeBlock struct {
	lines []AST
}

type Statement struct {
	terms      []AST
	properties *OpProp
}

type Identifier struct {
	name string
}

type IntLiteral struct {
	value int64
}

type FloatLiteral struct {
	value float64
}

type CharLiteral struct {
	value rune
}

type StringLiteral struct {
	value string
}

func NewCodeBlock(lines []AST) *CodeBlock {
	result := new(CodeBlock)
	result.lines = lines
	return result
}

func NewStatement(terms []AST, properties *OpProp) *Statement {
	result := new(Statement)
	result.terms = terms
	result.properties = properties
	return result
}
