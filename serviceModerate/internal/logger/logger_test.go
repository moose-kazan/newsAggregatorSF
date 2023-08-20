package logger

import "testing"

func TestLogger_formatMessage(t *testing.T) {
	type fields struct {
		prefix string
	}
	type args struct {
		reqId string
		msg   string
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
			args:   args{reqId: "xxx", msg: "message"},
			want:   "test1: [xxx] message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				prefix: tt.fields.prefix,
			}
			if got := l.formatMessage(tt.args.reqId, tt.args.msg); got != tt.want {
				t.Errorf("Logger.formatMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
