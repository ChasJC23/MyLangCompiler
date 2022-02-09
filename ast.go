package main

type AST interface {
}

type CodeBlock struct {
	lines []AST
}

type Statement struct {
	terms []AST
}

type Identifier struct {
	name string
}

type IntLiteral struct {
	value int
}

func NewCodeBlock(lines []AST) *CodeBlock {
	result := new(CodeBlock)
	result.lines = lines
	return result
}
