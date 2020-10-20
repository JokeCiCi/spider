package process

import (
	"strconv"
	"strings"
)

func UnescapeUnicode(text string) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(text), `\\u`, `\u`, -1))
	return str, err
}
