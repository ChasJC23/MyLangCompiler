package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestParser_ParseSource(t *testing.T) {
	testContext := NewOpContext()
	testContext.opTree.AddOperatorRune(';', STATEMENT_ENDING_TOKEN, 0)
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_RIGHT_ASSOCIATIVE,
		operators:  make(map[int]*OpProp),
	})
	assignProperties := testContext.AddOperatorToHighest([]string{":="}, 0, 2)
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_RIGHT_ASSOCIATIVE,
		operators:  make(map[int]*OpProp),
	})
	ternProperties := testContext.AddOperatorToHighest([]string{"?", ":"}, 0, 3)
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_LEFT_ASSOCIATIVE | IMPLIED_OPERATION,
		operators:  make(map[int]*OpProp),
	})
	eqProperties := testContext.AddOperatorToHighest([]string{"="}, 0, 2)
	ltProperties := testContext.AddOperatorToHighest([]string{"<"}, 0, 2)
	geProperties := testContext.AddOperatorToHighest([]string{">="}, 0, 2)
	conjProperties := testContext.AddOperatorToHighest([]string{}, 0, 2)
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_LEFT_ASSOCIATIVE,
		operators:  make(map[int]*OpProp),
	})
	addProperties := testContext.AddOperatorToHighest([]string{"+"}, 0, 2)
	subProperties := testContext.AddOperatorToHighest([]string{"-"}, 0, 2)
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: PREFIX,
		operators:  make(map[int]*OpProp),
	})
	//posProperties := testContext.AddOperatorToHighest([]string{"+"}, 0, 1)
	negProperties := testContext.AddOperatorToHighest([]string{"-"}, 0, 1)
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_LEFT_ASSOCIATIVE,
		operators:  make(map[int]*OpProp),
	})
	prodProperties := testContext.AddOperatorToHighest([]string{"*"}, 0, 2)
	divProperties := testContext.AddOperatorToHighest([]string{"/"}, 0, 2)
	/*juxProperties :=*/ testContext.AddOperatorToHighest([]string{}, 0, 2)
	tests := []struct {
		name       string
		expression *bufio.Reader
		context    *OpContext
		want       AST
	}{
		{"what's 9 + 10?", bufio.NewReader(strings.NewReader("9 + 10 = 21")), testContext,
			&CodeBlock{
				[]AST{
					&Statement{
						terms: []AST{
							&Statement{
								terms:      []AST{IntLiteral{9}, IntLiteral{10}},
								properties: addProperties,
							},
							IntLiteral{21},
						},
						properties: eqProperties,
					},
				},
			},
		},
		{"Inequalities", bufio.NewReader(strings.NewReader("0/1 < 1/3 < 1/2 < 2/3 < 1/1")), testContext,
			&CodeBlock{
				[]AST{
					&Statement{
						terms: []AST{
							&Statement{
								terms: []AST{
									&Statement{
										terms: []AST{
											&Statement{
												terms: []AST{
													&Statement{
														terms:      []AST{IntLiteral{0}, IntLiteral{1}},
														properties: divProperties,
													},
													&Statement{
														terms:      []AST{IntLiteral{1}, IntLiteral{3}},
														properties: divProperties,
													},
												},
												properties: ltProperties,
											},
											&Statement{
												terms: []AST{
													&Statement{
														terms:      []AST{IntLiteral{1}, IntLiteral{3}},
														properties: divProperties,
													},
													&Statement{
														terms:      []AST{IntLiteral{1}, IntLiteral{2}},
														properties: divProperties,
													},
												},
												properties: ltProperties,
											},
										},
										properties: conjProperties,
									},
									&Statement{
										terms: []AST{
											&Statement{
												terms:      []AST{IntLiteral{1}, IntLiteral{2}},
												properties: divProperties,
											},
											&Statement{
												terms:      []AST{IntLiteral{2}, IntLiteral{3}},
												properties: divProperties,
											},
										},
										properties: ltProperties,
									},
								},
								properties: conjProperties,
							},
							&Statement{
								terms: []AST{
									&Statement{
										terms:      []AST{IntLiteral{2}, IntLiteral{3}},
										properties: divProperties,
									},
									&Statement{
										terms:      []AST{IntLiteral{1}, IntLiteral{1}},
										properties: divProperties,
									},
								},
								properties: ltProperties,
							},
						},
						properties: conjProperties,
					},
				},
			},
		},
		{"Ternary time", bufio.NewReader(strings.NewReader("x >= 0 ? 1 : -1")), testContext,
			&CodeBlock{
				[]AST{
					&Statement{
						terms: []AST{
							&Statement{
								terms:      []AST{Identifier{"x"}, IntLiteral{0}},
								properties: geProperties,
							},
							IntLiteral{1},
							&Statement{
								terms:      []AST{IntLiteral{1}},
								properties: negProperties,
							},
						},
						properties: ternProperties,
					},
				},
			},
		},
		{"Multiple lines", bufio.NewReader(strings.NewReader("x := 9 - 4;y := x * 3")), testContext,
			&CodeBlock{
				[]AST{
					&Statement{
						terms: []AST{
							Identifier{"x"},
							&Statement{
								terms:      []AST{IntLiteral{9}, IntLiteral{4}},
								properties: subProperties,
							},
						},
						properties: assignProperties,
					},
					&Statement{
						terms: []AST{
							Identifier{"y"},
							&Statement{
								terms:      []AST{Identifier{"x"}, IntLiteral{3}},
								properties: prodProperties,
							},
						},
						properties: assignProperties,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				tokeniser: NewTokeniser(tt.expression, tt.context),
				opctx:     tt.context,
			}
			if got := p.ParseSource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
