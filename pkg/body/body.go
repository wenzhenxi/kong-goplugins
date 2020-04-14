package body

import "errors"

type Body interface {
	CheckContentType(string) bool
	GetStringByKey(string, string, string) (string, error)
}

func GetDataByKey(contentType, body, key string) (string, error) {

	if FormUrlHeader.CheckContentType(contentType) {
		s, err := FormUrlHeader.GetStringByKey(body, key, contentType)
		if err != nil {
			return "", err
		}
		return s, nil
	} else if RawHeader.CheckContentType(contentType) {
		s, err := RawHeader.GetStringByKey(body, key, contentType)
		if err != nil {
			return "", err
		}
		return s, nil
	} else if FormDataHeader.CheckContentType(contentType) {
		s, err := FormDataHeader.GetStringByKey(body, key, contentType)
		if err != nil {
			return "", err
		}
		return s, nil
	}
	return "", errors.New("Error ContentType")
}
