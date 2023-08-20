package filter

import "testing"

func TestBadWords(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{name: "Negative", arg: "Simple test comment", want: false},
		{name: "Positive single", arg: "qwerty", want: true},
		{name: "Positive first", arg: "qwerty/test comment", want: true},
		{name: "Positive middle", arg: "test qwerty: comment", want: true},
		{name: "Positive end", arg: "test comment, qwerty", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BadWords(tt.arg); got != tt.want {
				t.Errorf("BadWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
