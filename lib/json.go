package lib

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func PrettyJSON(b []byte) []byte {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, b, "", "    ")
	if err != nil {
		logrus.Warnf("prettyJSON parse error, value: %v err: %v", string(b), err)
		return b
	}
	return prettyJSON.Bytes()
}
