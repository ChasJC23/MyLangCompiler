package main

import "testing"

func TestSeparateFloat(t *testing.T) {
	type args struct {
		s    string
		base int
	}
	tests := []struct {
		name        string
		args        args
		wantMl      int
		wantHasExp  bool
		wantManSign bool
		wantExpSign bool
	}{
		{"all three", args{"31.574e30", 10}, 6, true, false, false},
		{"no exponent", args{"684.6516", 10}, 8, false, false, false},
		{"no fractional component", args{"468e12", 10}, 3, true, false, false},
		{"integer", args{"4525", 10}, 4, false, false, false},
		{"no integer component", args{".5467e34", 10}, 5, true, false, false},
		{"just fractional", args{".574654", 10}, 7, false, false, false},
		{"negative mantissa", args{"-643.765e3", 10}, 7, true, true, false},
		{"negative explicit", args{"-54.548", 10}, 6, false, true, false},
		{"negative exponent", args{"6e-5", 10}, 1, true, false, true},
		{"very positive", args{"+654.3e+345", 10}, 5, true, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIml, gotHasExp, gotManSign, gotExpSign := separateFloat(tt.args.s, tt.args.base)
			if gotIml != tt.wantMl {
				t.Errorf("separateFloat() gotIml = %v, want %v", gotIml, tt.wantMl)
			}
			if gotHasExp != tt.wantHasExp {
				t.Errorf("separateFloat() gotHasExp = %v, want %v", gotHasExp, tt.wantHasExp)
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
		{"fractional component 1", args{"3.14152535", 6}, "3.ABHN"},
		{"fractional component 2", args{"1313131.3131313", 4}, "1DDD.DDDC"},
		{"fractional component 3", args{"131313.13131313", 4}, "777.7777"},
		{"empty fractional component", args{"010101.", 2}, "111."},
		{"no int component", args{".3412013", 5}, ".J71F"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := squareBase(tt.args.s, tt.args.base); got != tt.want {
				t.Errorf("squareBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	type args struct {
		s    string
		base int
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"call original library base 10", args{"3456.645136e12", 10}, 3456.645136e12, false},
		{"call original library base 16", args{"2C67.3948Ap11", 16}, 0x2C67.3948Ap11, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFloat(tt.args.s, tt.args.base)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}
