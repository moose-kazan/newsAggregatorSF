package logger

import (
	"net/http"
	"testing"
)

func TestLogger_formatMessage(t *testing.T) {
	type fields struct {
		prefix string
	}
	type args struct {
		req *http.Request
		msg string
	}
	var r http.Request = http.Request{
		RemoteAddr: "127.0.0.1",
		Header:     http.Header{"X-Request-Id": []string{"xxx"}},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Main",
			fields: fields{prefix: "test1"},
			args:   args{req: &r, msg: "message"},
			want:   "test1: [127.0.0.1] [xxx] message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				prefix: tt.fields.prefix,
			}
			if got := l.formatMessage(tt.args.req, tt.args.msg); got != tt.want {
				t.Errorf("Logger.formatMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
