package main

import "testing"

func TestSeparateFloat(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name         string
		args         args
		wantIml      int
		wantFml      int
		wantExl      int
		wantHasRadix bool
		wantHasExp   bool
		wantNegMan   bool
		wantNegExp   bool
		wantManSign  bool
		wantExpSign  bool
	}{
		{"all three", args{"31.574e30"}, 2, 3, 2, true, true, false, false, false, false},
		{"no exponent", args{"684.6516"}, 3, 4, 0, true, false, false, false, false, false},
		{"no fractional component", args{"468e12"}, 3, 0, 2, false, true, false, false, false, false},
		{"integer", args{"4525"}, 4, 0, 0, false, false, false, false, false, false},
		{"no integer component", args{".5467e34"}, 0, 4, 2, true, true, false, false, false, false},
		{"just fractional", args{".574654"}, 0, 6, 0, true, false, false, false, false, false},
		{"negative mantissa", args{"-643.765e3"}, 3, 3, 1, true, true, true, false, true, false},
		{"negative explicit", args{"-54.548"}, 2, 3, 0, true, false, true, false, true, false},
		{"negative exponent", args{"6e-5"}, 1, 0, 1, false, true, false, true, false, true},
		{"very positive", args{"+654.3e+345"}, 3, 1, 3, true, true, false, false, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIml, gotFml, gotExl, gotHasRadix, gotHasExp, gotNegMan, gotNegExp, gotManSign, gotExpSign := separateFloat(tt.args.s)
			if gotIml != tt.wantIml {
				t.Errorf("separateFloat() gotIml = %v, want %v", gotIml, tt.wantIml)
			}
			if gotFml != tt.wantFml {
				t.Errorf("separateFloat() gotFml = %v, want %v", gotFml, tt.wantFml)
			}
			if gotExl != tt.wantExl {
				t.Errorf("separateFloat() gotExl = %v, want %v", gotExl, tt.wantExl)
			}
			if gotHasRadix != tt.wantHasRadix {
				t.Errorf("separateFloat() gotHasRadix = %v, want %v", gotHasRadix, tt.wantHasRadix)
			}
			if gotHasExp != tt.wantHasExp {
				t.Errorf("separateFloat() gotHasExp = %v, want %v", gotHasExp, tt.wantHasExp)
			}
			if gotNegMan != tt.wantNegMan {
				t.Errorf("separateFloat() gotNegMan = %v, want %v", gotNegMan, tt.wantNegMan)
			}
			if gotNegExp != tt.wantNegExp {
				t.Errorf("separateFloat() gotNegExp = %v, want %v", gotNegExp, tt.wantNegExp)
			}
			if gotManSign != tt.wantManSign {
				t.Errorf("separateFloat() gotManSign = %v, want %v", gotManSign, tt.wantManSign)
			}
			if gotExpSign != tt.wantExpSign {
				t.Errorf("separateFloat() gotExpSign = %v, want %v", gotExpSign, tt.wantExpSign)
			}
		})
	}
}

func Test_squareBase(t *testing.T) {
	type args struct {
		s    string
		base int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"binary", args{"111100100", 2}, "13210"},
		{"binary signed", args{"+11010100101", 2}, "+122211"},
		{"ternary", args{"10210111220212200", 3}, "123456780"},
		{"quaternary", args{"-12321311132130", 4}, "-6E7579C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := squareBase(tt.args.s, tt.args.base); got != tt.want {
				t.Errorf("squareBase() = %v, want %v", got, tt.want)
			}
		})
	}
}
