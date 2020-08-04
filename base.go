package mongox

import (
	"strings"
	"time"
)

// IsDup 判断 err 是否是 mongo E11000 错误。
func IsDup(err error) bool {
	return strings.Contains(err.Error(), "E11000")
}

// Now 以毫秒为单位返回当前时间
func Now() time.Time {
	return time.Unix(0, time.Now().UnixNano()/1e6*1e6)
}

// SplitSymbol 拆分字符串前面的 "+", "-"
// 如果是 "+"， symbol返回1. 如果是 "-"， symbol 返回 -1
func SplitSymbol(field string) (key string, symbol int32) {
	symbol = 1
	key = field

	if len(field) != 0 {
		switch field[0] {
		case '+':
			key = strings.TrimPrefix(field, "+")
			symbol = 1
		case '-':
			key = strings.TrimPrefix(field, "-")
			symbol = -1
		}
	}

	return key, symbol
}
