package utils

import "testing"

func TestNamedFormat(t *testing.T) {
	type args struct {
		format string
		params map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"default",
			args{
				format: "{name}.{absolute.2}.{ext}",
				params: map[string]interface{}{
					"name":     "test",
					"ext":      "mkv",
					"absolute": 1,
				},
			},
			"test.01.mkv",
		},
		{
			"padding string",
			args{
				format: "{absolute.2}.{ext}",
				params: map[string]interface{}{
					"ext":      "mkv",
					"absolute": "3",
				},
			},
			"03.mkv",
		},
		{
			"multi key format",
			args{
				format: "{absolute.2}.{absolute.3}.{ext}",
				params: map[string]interface{}{
					"ext":      "mkv",
					"absolute": 3,
				},
			},
			"03.003.mkv",
		},
		{
			"repeat key format",
			args{
				format: "{absolute.2}.{absolute.2}.{ext}",
				params: map[string]interface{}{
					"ext":      "mkv",
					"absolute": 3,
				},
			},
			"03.03.mkv",
		},
		{
			"white space",
			args{
				format: "{name }.{ absolute.2 }.{ext  }",
				params: map[string]interface{}{
					"name":     "test",
					"ext":      "mkv",
					"absolute": 1,
				},
			},
			"test.01.mkv",
		},
		{
			"no value",
			args{
				format: "{name }.{ absolute.2 }.{ext  }",
				params: map[string]interface{}{
					"name": "test",
					"ext":  "mkv",
				},
			},
			"test..mkv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NamedFormat(tt.args.format, tt.args.params); got != tt.want {
				t.Errorf("NamedFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
