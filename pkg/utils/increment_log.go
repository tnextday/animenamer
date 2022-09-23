package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

func MakeIncrementLog(dir, name string) string {
	existMaxLogIndex := 0
	iRe := regexp.MustCompile(fmt.Sprintf(`%s\.(\d+)\.log$`, name))
	if files, err := filepath.Glob(fmt.Sprintf("%s/%s.*.log", dir, name)); err == nil {
		for _, f := range files {
			_, fn := filepath.Split(f)
			matchs := iRe.FindStringSubmatch(fn)
			if len(matchs) != 2 {
				continue
			}
			if i, err := strconv.Atoi(matchs[1]); err == nil {
				if i > existMaxLogIndex {
					existMaxLogIndex = i
				}
			}
		}
	}
	return path.Join(dir, fmt.Sprintf("%s.%d.log", name, existMaxLogIndex+1))
}
