package sewan_go_sdk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
	req.Header.Add(httpAuthorization, httpTokenHeader+api.Token)
	logger.Println("req = ", req)
	logger.Println("reqError = ", reqError)
	if reqError == nil {
		resp, respError = clientTooler.Client.Do(api, req)
		logger.Println("resp = ", resp)
		if respError == nil {
			templateList, handlerRespError = clientTooler.Client.HandleResponse(resp,
				http.StatusOK,
				httpJsonContentType)
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
	if resp == nil {
		return "", errEmptyResp
	}
	if resp.Body == nil {
		return "", errEmptyRespBody
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get(httpRespContentType)
	if contentType != expectedBodyFormat {
		return nil, errRespStatusCodeBuilder(resp, expectedCode,
			"Wrong response content type,\n"+
				"expected :"+expectedBodyFormat+"\ngot :"+contentType)
	}
	switch contentType {
	case httpJsonContentType:
		return handleJsonContentType(resp, expectedCode)
	case httpHtmlTextContentType:
		return handleHtmlContentType(resp, expectedCode)
	case "":
		return nil, errRespStatusCodeBuilder(resp, expectedCode, "")
	default:
		return nil, errRespStatusCodeBuilder(resp, expectedCode,
			errApiUnhandledRespType+
				resp.Header.Get(httpRespContentType)+
				errValidateApiUrl)
	}
}

func handleJsonContentType(resp *http.Response,
	expectedCode int) (interface{}, error) {
	var (
		respBodyReader interface{}
	)
	bodyBytes, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, errRespStatusCodeBuilder(resp, expectedCode,
			"\nRead of response body error "+err1.Error())
	}
	if string(bodyBytes) == "" {
		return nil, errRespStatusCodeBuilder(resp, expectedCode, "")
	}
	err2 := json.Unmarshal(bodyBytes, &respBodyReader)
	if err2 != nil {
		return nil, errRespStatusCodeBuilder(resp, expectedCode,
			errJsonFormat+err2.Error()+
				"\nJson :"+string(bodyBytes))
	}
	err3 := errRespStatusCodeBuilder(resp, expectedCode, "")
	if err3 != nil {
		return nil, errors.New(err3.Error() +
			"\nResponse body error :" + string(bodyBytes))
	}
	return respBodyReader.(interface{}), nil
}

func handleHtmlContentType(resp *http.Response,
	expectedCode int) (interface{}, error) {
	bodyBytes, err4 := ioutil.ReadAll(resp.Body)
	if err4 != nil {
		return nil, errRespStatusCodeBuilder(resp, expectedCode, err4.Error())
	}
	err5 := errRespStatusCodeBuilder(resp, expectedCode, "")
	if err5 != nil {
		return nil, errors.New(err5.Error() +
			"\nResponse body error :" + string(bodyBytes))
	}
	return string(bodyBytes), nil
}
