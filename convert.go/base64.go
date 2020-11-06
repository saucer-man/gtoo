package convert

import (
	"encoding/base64"
)

func Base64encode(input []byte) string {
	str := base64.StdEncoding.EncodeToString(input)
	return str
}

func Base64decode(input string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return data, nil
}
