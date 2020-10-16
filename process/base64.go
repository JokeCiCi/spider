package process

import "encoding/base64"

func Base64Decode(origin string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(origin)
}
