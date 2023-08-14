package env

import (
	"os"
	"testing"
)

func TestGetInt(t *testing.T) {
	type args struct {
		name string
		def  int
		val  string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "correct",
			args: args{
				name: "ENV_CORRECT",
				def:  60,
				val:  "20",
			},
			want: 20,
		},
		{
			name: "incorrect",
			args: args{
				name: "ENV_INCORRECT",
				def:  70,
				val:  "xxx",
			},
			want: 70,
		},
		{
			name: "empty",
			args: args{
				name: "ENV_EMPTY",
				def:  80,
				val:  "",
			},
			want: 80,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.args.name, tt.args.val)
			if got := GetInt(tt.args.name, tt.args.def); got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStr(t *testing.T) {
	type args struct {
		name      string
		checkName string
		def       string
		val       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "correct",
			args: args{
				name:      "ENV_CORRECT",
				checkName: "ENV_CORRECT",
				def:       "xxx",
				val:       "yyy",
			},
			want: "yyy",
		},
		{
			name: "empty",
			args: args{
				name:      "ENV_EMPTY",
				checkName: "ENV_EMPTY",
				def:       "xxx",
				val:       "",
			},
			want: "",
		},
		{
			name: "unexists",
			args: args{
				name:      "ENV_CORRECT_1",
				checkName: "ENV_CORRECT_2",
				def:       "xxx",
				val:       "yyy",
			},
			want: "xxx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.args.name, tt.args.val)
			if got := GetStr(tt.args.checkName, tt.args.def); got != tt.want {
				t.Errorf("GetStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
