package sewan_go_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type HttpClienterDummy struct{}

func (client HttpClienterDummy) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client HttpClienterDummy) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client HttpClienterDummy) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

var (
	ResourceDeletionSuccessHttpClienterFake = HttpClienterDummy{}
	ResourceUpdateSuccessHttpClienterFake   = HttpClienterDummy{}
	ResourceCreationFailureHttpClienterFake = ResourceDeletionFailureHttpClienterFake{}
	ResourceUpdateFailureHttpClienterFake   = ResourceDeletionFailureHttpClienterFake{}
	ResourceReadFailureHttpClienterFake     = ResourceDeletionFailureHttpClienterFake{}
)

type ErrorResponseHttpClienterFake struct{}

func (client ErrorResponseHttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, errDoRequest
}
func (client ErrorResponseHttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client ErrorResponseHttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

type GetTemplatesListFailureHttpClienterFake struct{}

func (client GetTemplatesListFailureHttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client GetTemplatesListFailureHttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, errors.New("GetTemplatesList() error")
}
func (client GetTemplatesListFailureHttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, errors.New("HandleResponse() error")
}

type GetTemplatesListSuccessHttpClienterFake struct{}

func (client GetTemplatesListSuccessHttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client GetTemplatesListSuccessHttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return templatesList, nil
}
func (client GetTemplatesListSuccessHttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return templatesList, nil
}

type HandleRespErrHttpClienterFake struct{}

func (client HandleRespErrHttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client HandleRespErrHttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client HandleRespErrHttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, errHandleResponse
}

type Error404HttpClienterFake struct{}

func (client Error404HttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	resp := http.Response{}
	resp.Header = map[string][]string{}
	resp.Header.Add(httpRespContentType, httpJsonContentType)
	resp.StatusCode = http.StatusNotFound
	resp.Status = notFoundRespStatus
	body := Resp_Body{"Not found."}
	js, _ := json.Marshal(body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(js))
	return &resp, nil
}
func (client Error404HttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client Error404HttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

type VmCreationSuccessHttpClienterFake struct{}

func (client VmCreationSuccessHttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client VmCreationSuccessHttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VmCreationSuccessHttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return noTemplateVmMap, nil
}

type VmReadSuccessHttpClienterFake struct{}

func (client VmReadSuccessHttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client VmReadSuccessHttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VmReadSuccessHttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return noTemplateVmMap, nil
}

type VdcReadSuccessHttpClienterFake struct{}

func (client VdcReadSuccessHttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client VdcReadSuccessHttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VdcReadSuccessHttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return vdcReadResponseMap, nil
}

type ResourceDeletionFailureHttpClienterFake struct{}

func (client ResourceDeletionFailureHttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, errEmptyResp
}
func (client ResourceDeletionFailureHttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client ResourceDeletionFailureHttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

type CheckRedirectReqFailure_HttpClienterFake struct{}

func (client CheckRedirectReqFailure_HttpClienterFake) Do(api *API,
	req *http.Request) (*http.Response, error) {
	resp := http.Response{}
	return &resp, errors.New(checkRedirectFailure)
}
func (client CheckRedirectReqFailure_HttpClienterFake) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client CheckRedirectReqFailure_HttpClienterFake) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

type FakeHttpClienter struct{}

func (client FakeHttpClienter) Do(api *API, req *http.Request) (*http.Response, error) {
	var err error
	err = nil
	type body struct {
		detail string `json:"detail"`
	}
	resp := http.Response{}
	if api.URL != noRespApiUrl {
		resp.Status = "200 OK"
		resp.StatusCode = http.StatusOK
		switch {
		case api.URL == wrongApiUrl || api.URL == noRespJsonBodyApiUrl:
			resp.Header = map[string][]string{httpRespContentType: {"text/plain; charset=utf-8"}}
			resp.Body = ioutil.NopCloser(bytes.NewBufferString("A plain text."))
		case api.URL == rightApiUrl:
			if api.Token != rightApiToken {
				resp.Status = "401 Unauthorized"
				resp.StatusCode = http.StatusUnauthorized
				resp.Body = ioutil.NopCloser(bytes.NewBufferString("{\"detail\":\"Invalid token.\"}"))
			} else {
				resp.Header = map[string][]string{httpRespContentType: {httpJsonContentType}}
				bodyJson, _ := json.Marshal(body{detail: ""})
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyJson))
			}
		}
	} else {
		err = errors.New("No response error.")
	}
	return &resp, err
}
func (client FakeHttpClienter) GetTemplatesList(clientTooler *ClientTooler,
	enterpriseSlug string, api *API) ([]interface{}, error) {
	return nil, nil
}
func (client FakeHttpClienter) HandleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}
