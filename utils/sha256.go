package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/pkg/errors"
)

func Base64EncodedSha256(structured any) (string, error) {
	data, err := json.Marshal(structured)
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal")
	}

	sha256Bytes := sha256.Sum256(data)

	return base64.StdEncoding.EncodeToString(sha256Bytes[:]), nil
}
