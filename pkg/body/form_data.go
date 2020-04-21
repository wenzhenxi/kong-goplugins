package body

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
)

var FormDataHeader FormData

type FormData struct {
	Body
}

func (f FormData) CheckContentType(contentType string) bool {

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return mediatype == "multipart/form-data"
}

func (f FormData) GetStringByKey(body, key,contentType string) (string, error) {

	bodyBuf := bytes.NewBuffer(nil)
	bodyBuf.WriteString(body)

	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {contentType}},
		Body: ioutil.NopCloser(bodyBuf),
	}

	req.ParseMultipartForm(128)

	return req.Form.Get(key), nil
}
