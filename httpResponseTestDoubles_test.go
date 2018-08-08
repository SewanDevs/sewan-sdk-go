package sewan_go_sdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpResponseFake_OKJson() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpJsonContentType)
	Resp_BodyJson, _ := json.Marshal(Resp_Body{Detail: "a simple json"})
	response.Body = ioutil.NopCloser(bytes.NewBuffer(Resp_BodyJson))
	return &response
}

func HttpResponseFake_OKTemplateListJson() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpJsonContentType)
	Resp_BodyJson, _ := json.Marshal(templatesList)
	response.Body = ioutil.NopCloser(bytes.NewBuffer(Resp_BodyJson))
	return &response
}

func HttpResponseFake_500_texthtml() *http.Response {
	response := http.Response{}
	response.Status = "500 Internal Server Error"
	response.StatusCode = http.StatusInternalServerError
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpHtmlTextContentType)
	response.Body = ioutil.NopCloser(bytes.NewBufferString("<h1>Server Error (500)</h1>"))
	return &response
}

func HttpResponseFake_500Json() *http.Response {
	response := http.Response{}
	response.Status = "500 Internal Server Error"
	response.StatusCode = http.StatusInternalServerError
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpJsonContentType)
	Resp_BodyJson, _ := json.Marshal(Resp_Body{Detail: "a json response Resp_Body"})
	response.Body = ioutil.NopCloser(bytes.NewBuffer(Resp_BodyJson))
	return &response
}

func HttpResponseFake_OK_txthtml() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpHtmlTextContentType)
	response.Body = ioutil.NopCloser(bytes.NewBufferString("<h1>An html text</h1>"))
	return &response
}

func HttpResponseFakeOkNilBody() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	return &response
}

func HttpResponseFake_OK_wrongjson() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, httpJsonContentType)
	response.Body = ioutil.NopCloser(bytes.NewBufferString("a bad formated json"))
	return &response
}

func HttpResponseFake_OK_image() *http.Response {
	response := http.Response{}
	response.Status = "200 OK"
	response.StatusCode = http.StatusOK
	response.Header = map[string][]string{}
	response.Header.Add(httpRespContentType, "image")
	response.Body = ioutil.NopCloser(bytes.NewBufferString("a response non empty body"))
	return &response
}
