package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestParser_ParseSource(t *testing.T) {
	testContext := NewOpContext()
	testContext.opTree.AddOperatorRune('{', OPEN_CODE_BLOCK_TOKEN, 0)
	testContext.opTree.AddOperatorRune('}', CLOSE_CODE_BLOCK_TOKEN, 0)
	testContext.opTree.AddOperatorRune('(', OPEN_PARENS_TOKEN, 0)
	testContext.opTree.AddOperatorRune(')', CLOSE_PARENS_TOKEN, 0)
	testContext.opTree.AddOperatorRune(';', STATEMENT_ENDING_TOKEN, 0)
	testContext.AddFixedTokenOperator([]rune("true"), TRUE_LITERAL, 0)
	testContext.AddFixedTokenOperator([]rune("false"), FALSE_LITERAL, 0)
	testContext.AddFixedTokenOperator([]rune("//"), COMMENT_TOKEN, COMMENT_FLAG)
	testContext.AddControlOperator([]rune("/*"), OPEN_COMMENT_FLAG)
	testContext.AddControlOperator([]rune("*/"), CLOSE_COMMENT_FLAG)

	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: PREFIX,
		operators:  make(map[int]*OpProp),
	})
	ifProperties := testContext.AddOperatorToHighest([]string{"if"}, 0b10, 2)
	testContext.AddHighestPrecedenceLevel(&PrecedenceLevel{
		properties: INFIX_RIGHT_ASSOCIATIVE,
		operators:  make(map[int]*OpProp),
	})
	declareProperties := testContext.AddOperatorToHighest([]string{":="}, 0, 2)
	assignProperties := testContext.AddOperatorToHighest([]string{"<-"}, 0, 2)
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
	gtProperties := testContext.AddOperatorToHighest([]string{">"}, 0, 2)
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
	// juxProperties := testContext.AddOperatorToHighest([]string{}, 0, 2)
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
		{"Multiple lines", bufio.NewReader(strings.NewReader("x := 9 - 4;\ny := x * 3")), testContext,
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
						properties: declareProperties,
					},
					&Statement{
						terms: []AST{
							Identifier{"y"},
							&Statement{
								terms:      []AST{Identifier{"x"}, IntLiteral{3}},
								properties: prodProperties,
							},
						},
						properties: declareProperties,
					},
				},
			},
		},
		{"Condition", bufio.NewReader(strings.NewReader("if true { mean := (a + b) / 2 }")), testContext,
			&CodeBlock{
				[]AST{
					&Statement{
						terms: []AST{
							BoolLiteral{true},
							&CodeBlock{
								[]AST{
									&Statement{
										terms: []AST{
											Identifier{"mean"},
											&Statement{
												terms: []AST{
													&Statement{
														terms: []AST{
															Identifier{"a"},
															Identifier{"b"},
														},
														properties: addProperties,
													},
													IntLiteral{2},
												},
												properties: divProperties,
											},
										},
										properties: declareProperties,
									},
								},
							},
						},
						properties: ifProperties,
					},
				},
			},
		},
		{"Obfuscated", bufio.NewReader(strings.NewReader("x:=3;y:=x>6?3:2;if y>x y<-y-1")), testContext,
			&CodeBlock{
				[]AST{
					&Statement{
						terms:      []AST{Identifier{"x"}, IntLiteral{3}},
						properties: declareProperties,
					},
					&Statement{
						terms: []AST{
							Identifier{"y"},
							&Statement{
								terms: []AST{
									&Statement{
										terms:      []AST{Identifier{"x"}, IntLiteral{6}},
										properties: gtProperties,
									},
									IntLiteral{3},
									IntLiteral{2},
								},
								properties: ternProperties,
							},
						},
						properties: declareProperties,
					},
					&Statement{
						terms: []AST{
							&Statement{
								terms:      []AST{Identifier{"y"}, Identifier{"x"}},
								properties: gtProperties,
							},
							&Statement{
								terms: []AST{
									Identifier{"y"},
									&Statement{
										terms:      []AST{Identifier{"y"}, IntLiteral{1}},
										properties: subProperties,
									},
								},
								properties: assignProperties,
							},
						},
						properties: ifProperties,
					},
				},
			},
		},
		{"Commented Code", bufio.NewReader(strings.NewReader("x1 := -0.9 /* this is a C style comment */\nx2 := x1 * x1 // and this is a single line comment\n - 0.9\nx3 /* third iteration */ := /* square */ x2 * x2 /* seed */ - 0.9\n// and one more comment to finish it off")), testContext,
			&CodeBlock{
				[]AST{
					&Statement{
						terms: []AST{
							Identifier{"x1"},
							&Statement{
								terms:      []AST{FloatLiteral{0.5}},
								properties: negProperties,
							},
						},
						properties: declareProperties,
					},
					&Statement{
						terms: []AST{
							Identifier{"x2"},
							&Statement{
								terms: []AST{
									&Statement{
										terms:      []AST{Identifier{"x1"}, Identifier{"x1"}},
										properties: prodProperties,
									},
									FloatLiteral{0.5},
								},
								properties: subProperties,
							},
						},
						properties: declareProperties,
					},
					&Statement{
						terms: []AST{
							Identifier{"x3"},
							&Statement{
								terms: []AST{
									&Statement{
										terms:      []AST{Identifier{"x2"}, Identifier{"x2"}},
										properties: prodProperties,
									},
									FloatLiteral{0.5},
								},
								properties: subProperties,
							},
						},
						properties: declareProperties,
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
