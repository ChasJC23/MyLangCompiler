package main

type ASTStack []AST

func (s *ASTStack) Push(v AST) {
	*s = append(*s, v)
}

func (s *ASTStack) Pop() AST {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}
