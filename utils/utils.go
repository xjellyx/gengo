package utils

import (
	"os"
	"strings"
)

// SQLColumnToHumpStyle sql转换成驼峰模式
func SQLColumnToHumpStyle(in string) (ret string) {
	for i := 0; i < len(in); i++ {
		if i > 0 && in[i-1] == '_' && in[i] != '_' {
			s := strings.ToUpper(string(in[i]))
			ret += s
		} else if in[i] == '_' {
			continue
		} else {
			ret += string(in[i])
		}
	}
	return
}

// SQLColumn2PkgStyle sql转-格式
func SQLColumn2PkgStyle(in string) (ret string) {
	var (
		arr  []string
		data []string
	)
	arr = strings.Split(in, "_")
	for _, v := range arr {
		if len(v) > 0 {
			data = append(data, v)
		}
	}
	return strings.Join(data, "-")
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
