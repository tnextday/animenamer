package namer

import (
	"fmt"
	"strconv"
)

type Infos map[string]interface{}

func (infos Infos) GetString(key string) string {
	v, exists := infos[key]
	if !exists {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (infos Infos) GetInt(key string) int {
	v, exists := infos[key]
	if !exists {
		return 0
	}
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return v.(int)
	case string:
		i, _ := strconv.Atoi(v.(string))
		return i
	default:
		return 0
	}
}
