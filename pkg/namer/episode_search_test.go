package namer

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_matchEpisodeFile(t *testing.T) {
	type args struct {
		filter   *regexp.Regexp
		filename string
	}
	tests := []struct {
		name      string
		args      args
		wantInfos Infos
		wantErr   bool
	}{
		{
			name: "",
			args: args{
				filter:   regexp.MustCompile(`(?P<name>\w+)\.(?P<absolute>\d+).*\.(?P<ext>\w+)$`),
				filename: "anime.001.720P.crc_aabbccdd.mp4",
			},
			wantInfos: Infos{
				"name":     "anime",
				"absolute": 1,
				"ext":      "mp4",
			},
			wantErr: false,
		},
		{
			name: "not match",
			args: args{
				filter:   regexp.MustCompile(`(?P<name>\w+)\.(?P<absolute>\d+).*\.(?P<ext>\w+)$`),
				filename: "anime.crc_aabbccdd.mp4",
			},
			wantInfos: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfos, err := matchEpisodeFile(tt.args.filter, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchEpisodeFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfos, tt.wantInfos) {
				t.Errorf("matchEpisodeFile() gotInfos = %v, want %v", gotInfos, tt.wantInfos)
			}
		})
	}
}

func Test_searchSubtitles(t *testing.T) {
	subtitlesExt := NewFileExtFromString("ass,ssa,srt,sub", ",")
	type args struct {
		mediaName    string
		filenameList []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				mediaName: "a.mkv",
				filenameList: []string{
					"a.mkv",
					"a.ass",
					"a.en.ass",
					"a.avi",
					"b.ass",
				},
			},
			want: []string{
				"a.ass",
				"a.en.ass",
			},
		},
		{
			name: "no subtitles",
			args: args{
				mediaName: "c.mkv",
				filenameList: []string{
					"a.mkv",
					"a.en.ass",
					"b.ass",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchSubtitles(tt.args.mediaName, tt.args.filenameList, subtitlesExt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchSubtitles() = %v, want %v", got, tt.want)
			}
		})
	}
}
