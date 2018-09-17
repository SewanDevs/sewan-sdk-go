package sewansdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HTTPResponseFakeOKJSON() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpJSONContentType)
	RespBodyJSON, _ := json.Marshal(RespBody{Detail: "a simple json"})
	response.Body = ioutil.NopCloser(bytes.NewBuffer(RespBodyJSON))
	return &response
}

func HTTPResponseFakeOKTemplateListJSON() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpJSONContentType)
	RespBodyJSON, _ := json.Marshal(templateMetaDataList)
	response.Body = ioutil.NopCloser(bytes.NewBuffer(RespBodyJSON))
	return &response
}

func HTTPResponseFake500Texthtml() *http.Response {
	response := http.Response{}
	response.Status = "500 Internal Server Error"
	response.StatusCode = http.StatusInternalServerError
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpHTMLTextContentType)
	response.Body = ioutil.NopCloser(bytes.NewBufferString("<h1>Server Error (500)</h1>"))
	return &response
}

func HTTPResponseFake500Json() *http.Response {
	response := http.Response{}
	response.Status = "500 Internal Server Error"
	response.StatusCode = http.StatusInternalServerError
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpJSONContentType)
	RespBodyJSON, _ := json.Marshal(RespBody{Detail: "a json response RespBody"})
	response.Body = ioutil.NopCloser(bytes.NewBuffer(RespBodyJSON))
	return &response
}

func HTTPResponseFakeOkNilBody() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	return &response
}

func HTTPResponseFakeOKWrongjson() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpJSONContentType)
	response.Body = ioutil.NopCloser(bytes.NewBufferString("a bad formated json"))
	return &response
}

func HTTPResponseFakeOKImage() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, "image")
	response.Body = ioutil.NopCloser(bytes.NewBufferString("a response non empty body"))
	return &response
}
