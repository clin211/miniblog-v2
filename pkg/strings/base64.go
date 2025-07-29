package strings

import (
	"bytes"
	"encoding/base64"
	"io"
)

func DecodeBase64(i string) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(i)))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
