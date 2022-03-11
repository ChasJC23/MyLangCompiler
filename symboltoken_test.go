package main

import "testing"

func TestSeparateFloat(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
		want2 int
	}{
		{"all three", args{"31.574e30"}, 2, 3, 2},
		{"no exponent", args{"684.6516"}, 3, 4, 0},
		{"no fractional component", args{"468e12"}, 3, 0, 2},
		{"integer", args{"4525"}, 4, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := SeparateFloat(tt.args.s)
			if got != tt.want {
				t.Errorf("SeparateFloat() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SeparateFloat() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("SeparateFloat() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
