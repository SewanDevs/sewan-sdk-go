package sewan_go_sdk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ClientTooler struct {
	Client Clienter
}
type Clienter interface {
	Do(api *API, req *http.Request) (*http.Response, error)
	GetTemplatesList(clientTooler *ClientTooler,
		enterpriseSlug string, api *API) ([]interface{}, error)
	HandleResponse(resp *http.Response,
		expectedCode int,
		expectedBodyFormat string) (interface{}, error)
}
type HttpClienter struct{}

func (client HttpClienter) Do(api *API, req *http.Request) (*http.Response, error) {
	resp, err := api.Client.Do(req)
	return resp, err
}

func (client HttpClienter) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {

	var (
		reqError           error
		respError          error
		handlerRespError   error
		returnError        error         = nil
		templateList       interface{}   = nil
		returnTemplateList []interface{} = nil
		templatesListUrl   strings.Builder
	)
	resp := &http.Response{}
	req := &http.Request{}
	templatesListUrl.WriteString(api.URL)
	templatesListUrl.WriteString("template/?enterprise__slug=")
	templatesListUrl.WriteString(enterpriseSlug)
	logger := loggerCreate("getTemplatesList.log")

	logger.Println("templatesListUrl.String() = ", templatesListUrl.String())
	req, reqError = http.NewRequest("GET",
		templatesListUrl.String(),
		nil)
	req.Header.Add("authorization", "Token "+api.Token)
	logger.Println("req = ", req)
	logger.Println("reqError = ", reqError)
	if reqError == nil {
		resp, respError = clientTooler.Client.Do(api, req)
		logger.Println("resp = ", resp)
		if respError == nil {
			templateList, handlerRespError = clientTooler.Client.HandleResponse(resp,
				http.StatusOK,
				"application/json")
			if templateList != nil {
				returnTemplateList = templateList.([]interface{})
			}
			if handlerRespError != nil {
				returnError = handlerRespError
			}
		} else {
			returnError = respError
		}
	} else {
		returnError = reqError
	}
	logger.Println("returnTemplateList = ", returnTemplateList)
	return returnTemplateList, returnError
}

func (client HttpClienter) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {

	var (
		respError      error       = nil
		responseBody   interface{} = nil
		contentType    string
		bodyBytes      []byte
		respBodyReader interface{}
		readBodyError  error = nil
		readJsonError  error = nil
	)

	if resp.StatusCode == expectedCode {
		contentType = resp.Header.Get("Content-Type")

		if contentType == expectedBodyFormat {
			switch contentType {
			case "application/json":
				bodyBytes, readBodyError = ioutil.ReadAll(resp.Body)
				readJsonError = json.Unmarshal(bodyBytes, &respBodyReader)
				switch {
				case readBodyError != nil:
					respError = errors.New("Read of response body error " +
						readBodyError.Error())
				case readJsonError != nil:
					respError = errors.New("Response body is not a properly formated json :" +
						readJsonError.Error())
				default:
					responseBody = respBodyReader.(interface{})
				}
			case "text/html":
				bodyBytes, readBodyError = ioutil.ReadAll(resp.Body)
				responseBody = string(bodyBytes)
			case "":
				responseBody = nil
			default:
				respError = errors.New("Unhandled api response type : " +
					resp.Header.Get("Content-Type") +
					"\nPlease validate the configuration api url.")
			}
		} else {
			respError = errors.New("Wrong content type, \n\r expected :" +
				expectedBodyFormat + "\n\r got :" + contentType)
		}
	} else {
		respError = errors.New("Wrong response status code, \n\r expected :" +
			strconv.Itoa(expectedCode) + "\n\r got :" + strconv.Itoa(resp.StatusCode) +
			"\n\rFull response status : " + resp.Status)
	}
	return responseBody, respError
}
