package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestParser_ParseSource(t *testing.T) {
	testContext := NewOpContext()
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_LEFT_ASSOCIATIVE | IMPLIED_OPERATION,
		operators:  make(map[int]*OpProp),
	})
	testContext.AddOperatorToHighest([]rune("="), &OpProp{
		subsequentSymbols:  nil,
		codeBlockArguments: 0,
		argumentCount:      2,
	})
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_LEFT_ASSOCIATIVE,
		operators:  make(map[int]*OpProp),
	})
	testContext.AddOperatorToHighest([]rune("+"), &OpProp{
		subsequentSymbols:  nil,
		codeBlockArguments: 0,
		argumentCount:      2,
	})
	testExpression := bufio.NewReader(strings.NewReader("9 + 10 = 21"))
	testTokeniser := NewTokeniser(testExpression, testContext)
	type fields struct {
		tokeniser *Tokeniser
		opctx     *OpContext
	}
	tests := []struct {
		name   string
		fields fields
		want   AST
	}{
		{"expression test", fields{testTokeniser, testContext}, Statement{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				tokeniser: tt.fields.tokeniser,
				opctx:     tt.fields.opctx,
			}
			if got := p.ParseSource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
