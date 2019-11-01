package namer

import (
	"reflect"
	"testing"
)

func TestNewFileExtFromString(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name string
		args args
		want FileExt
	}{
		{
			name: "default",
			args: args{
				s:   "Mkv, .Rm ,,",
				sep: ",",
			},
			want: FileExt{
				".mkv": struct{}{},
				".rm":  struct{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileExtFromString(tt.args.s, tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileExtFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExt_Is(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		fe   FileExt
		args args
		want bool
	}{
		{
			fe:   NewFileExtFromString("mkv,rm", ","),
			args: args{"a.rm"},
			want: true,
		},
		{
			fe:   NewFileExtFromString("mkv,rm", ","),
			args: args{"a.rmvb"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fe.Is(tt.args.filename); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
