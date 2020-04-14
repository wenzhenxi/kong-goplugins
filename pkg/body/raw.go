package body

import "github.com/tidwall/gjson"

var RawHeader Raw

type Raw struct {
	Body
}

func (r Raw) CheckContentType(contentType string) bool {
	return contentType == "application/json"
}

func (f Raw) GetStringByKey(body, key, contentType string) (string, error) {

	return gjson.Parse(body).Get(key).String(), nil
}
