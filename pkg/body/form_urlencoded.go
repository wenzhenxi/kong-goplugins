package body

import (
	"net/url"
)

var FormUrlHeader FormUrl

type FormUrl struct {
	Body
}

func (f FormUrl) CheckContentType(contentType string) bool {
	return contentType == "application/x-www-form-urlencoded"
}

func (f FormUrl) GetStringByKey(body, key,contentType string) (string, error) {

	u, err := url.ParseQuery(body)
	if err != nil {
		return "", err
	}
	return u.Get(key), nil
}
