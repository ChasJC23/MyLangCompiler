package main

import (
	"strconv"
	"strings"
)

type AST interface {
	String() string
}

type Statement struct {
	terms      []AST
	properties *OpProp
}

func (st Statement) String() string {
	var builder strings.Builder
	builder.WriteByte('[')
	builder.WriteString(st.properties.initSymbol)
	builder.WriteString("](")
	builder.WriteByte(' ')
	for _, term := range st.terms {
		builder.WriteString(term.String())
		builder.WriteByte(' ')
	}
	builder.WriteByte(')')
	return builder.String()
}

type Identifier struct {
	name string
}

func (id Identifier) String() string {
	var builder strings.Builder
	builder.WriteByte('<')
	builder.WriteString(id.name)
	builder.WriteByte('>')
	return builder.String()
}

type IntLiteral struct {
	value int64
}

func (il IntLiteral) String() string {
	var builder strings.Builder
	builder.WriteString(strconv.FormatInt(il.value, 10))
	return builder.String()
}

type FloatLiteral struct {
	value float64
}

func (fl FloatLiteral) String() string {
	var builder strings.Builder
	builder.WriteString(strconv.FormatFloat(fl.value, 'g', -1, 64))
	builder.WriteByte('f')
	return builder.String()
}

type CharLiteral struct {
	value rune
}

func (cl CharLiteral) String() string {
	var builder strings.Builder
	builder.WriteByte('\'')
	builder.WriteRune(cl.value)
	builder.WriteByte('\'')
	return builder.String()
}

type StringLiteral struct {
	value string
}

func (sl StringLiteral) String() string {
	var builder strings.Builder
	builder.WriteByte('"')
	builder.WriteString(sl.value)
	builder.WriteByte('"')
	return builder.String()
}

type BoolLiteral struct {
	value bool
}

func (bl BoolLiteral) String() string {
	if bl.value {
		return "true"
	} else {
		return "false"
	}
}

func NewStatement(terms []AST, properties *OpProp) *Statement {
	result := new(Statement)
	result.terms = terms
	result.properties = properties
	return result
}
