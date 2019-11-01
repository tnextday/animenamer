package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	paramRegex = regexp.MustCompile(`\{\s*(\w+?)(\.(\d+))?\s*\}`)
)

func NamedFormat(format string, params map[string]interface{}) string {
	matches := paramRegex.FindAllStringSubmatch(format, -1)
	outString := format
	for _, match := range matches {
		s := match[0]
		paramKey := match[1]
		v := ""
		if paramValue, exists := params[paramKey]; exists {
			v = strings.TrimSpace(fmt.Sprintf("%v", paramValue))
			padding, _ := strconv.Atoi(match[3])
			if padding > len(v) {
				v = strings.Repeat("0", padding-len(v)) + v
			}
		}
		outString = strings.ReplaceAll(outString, s, v)
	}
	return outString
}
