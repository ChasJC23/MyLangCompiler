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
	testContext.AddOperatorToHighest([]string{"="}, 0, 2)
	eqProperties := testContext.precedenceList.Back().Value.(*PrecedenceLevel).operators[testContext.opTree.GetToken([]rune("="))]
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_LEFT_ASSOCIATIVE,
		operators:  make(map[int]*OpProp),
	})
	testContext.AddOperatorToHighest([]string{"+"}, 0, 2)
	plusProperties := testContext.precedenceList.Back().Value.(*PrecedenceLevel).operators[testContext.opTree.GetToken([]rune("+"))]
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
		{"what's 9 + 10?", fields{testTokeniser, testContext}, &CodeBlock{[]AST{&Statement{
			terms: []AST{
				&Statement{
					terms:      []AST{IntLiteral{9}, IntLiteral{10}},
					properties: plusProperties,
				},
				IntLiteral{21},
			},
			properties: eqProperties,
		}}}},
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
