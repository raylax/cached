package redis

import (
	"io"
	"reflect"
	"testing"
)

func TestReadResp(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "null bulk string",
			args: args{
				bytes: []byte("$-1\r\n"),
			},
			want: &RespString{resp: resp{Type: RespTypeBulkString}, Data: nil},
		},
		{
			name: "bulk string size error",
			args: args{
				bytes: []byte("$1\r\n"),
			},
			wantErr: true,
		},
		{
			name: "bulk string size error",
			args: args{
				bytes: []byte("$x1\r\n"),
			},
			wantErr: true,
		},
		{
			name: "bulk string skip crlf error",
			args: args{
				bytes: []byte("$1\r\nx"),
			},
			wantErr: true,
		},
		{
			name: "array skip CRLF error",
			args: args{
				bytes: []byte("*1\r\n:123,"),
			},
			wantErr: true,
		},
		{
			name: "array error",
			args: args{
				bytes: []byte("*1\r\n"),
			},
			wantErr: true,
		},
		{
			name: "skip CRLF error",
			args: args{
				bytes: []byte(":123,"),
			},
			wantErr: true,
		},
		{
			name: "invalid CRLF",
			args: args{
				bytes: []byte("+OK\n"),
			},
			wantErr: true,
		},
		{
			name: "invalid CRLF",
			args: args{
				bytes: []byte("+OK\r"),
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: args{
				bytes: []byte("123"),
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: args{
				bytes: []byte(""),
			},
			wantErr: true,
		},
		{
			name: "integer",
			args: args{
				bytes: []byte(":123\r\n"),
			},
			want: &RespInteger{resp: resp{Type: RespTypeInteger}, Value: 123},
		},
		{
			name: "simple string",
			args: args{
				bytes: []byte("+OK\r\n"),
			},
			want: &RespString{resp: resp{Type: RespTypeSimpleString}, Data: []byte("OK")},
		},
		{
			name: "bulk string",
			args: args{
				bytes: []byte("$6\r\nfoobar\r\n"),
			},
			want: &RespString{resp: resp{Type: RespTypeBulkString}, Data: []byte("foobar")},
		},
		{
			name: "error",
			args: args{
				bytes: []byte("-Error message\r\n"),
			},
			want: &RespError{resp: resp{Type: RespTypeError}, Message: "Error message"},
		},
		{
			name: "null array",
			args: args{
				bytes: []byte("*-1\r\n"),
			},
			want: &RespArray{resp: resp{Type: RespTypeArray}, Elements: nil},
		},
		{
			name: "integer array",
			args: args{
				bytes: []byte("*2\r\n:1\r\n:2\r\n"),
			},
			want: &RespArray{resp: resp{Type: RespTypeArray}, Elements: []any{
				&RespInteger{resp: resp{Type: RespTypeInteger}, Value: 1},
				&RespInteger{resp: resp{Type: RespTypeInteger}, Value: 2},
			}},
		},
		{
			name: "simple string array",
			args: args{
				bytes: []byte("*2\r\n+foo\r\n+bar\r\n"),
			},
			want: &RespArray{resp: resp{Type: RespTypeArray}, Elements: []any{
				&RespString{resp: resp{Type: RespTypeSimpleString}, Data: []byte("foo")},
				&RespString{resp: resp{Type: RespTypeSimpleString}, Data: []byte("bar")},
			}},
		},
		{
			name: "bulk string array",
			args: args{
				bytes: []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"),
			},
			want: &RespArray{resp: resp{Type: RespTypeArray}, Elements: []any{
				&RespString{resp: resp{Type: RespTypeBulkString}, Data: []byte("foo")},
				&RespString{resp: resp{Type: RespTypeBulkString}, Data: []byte("bar")},
			}},
		},
		{
			name: "array size error",
			args: args{
				bytes: []byte("*a2\r\n"),
			},
			wantErr: true,
		},
		{
			name: "error crlf",
			args: args{
				bytes: []byte("-Error message\n"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadResp(&fakeReader{b: tt.args.bytes})
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadResp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadResp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type fakeReader struct {
	b []byte
	i int
}

func (f *fakeReader) Next(n int) (p []byte, err error) {
	if f.i+n > len(f.b) {
		return nil, io.EOF
	}
	p = f.b[f.i : f.i+n]
	f.i += n
	return
}

func (f *fakeReader) Peek(n int) (buf []byte, err error) {
	if f.i+n > len(f.b) {
		return nil, io.EOF
	}
	buf = f.b[f.i : f.i+n]
	return
}

func (f *fakeReader) Skip(n int) (err error) {
	if f.i+n > len(f.b) {
		return io.EOF
	}
	f.i += n
	return
}
