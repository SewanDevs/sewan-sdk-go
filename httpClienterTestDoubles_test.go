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
	ResourceDeletionSuccessHttpClienterFake               = HttpClienterDummy{}
	ResourceUpdateSuccessHttpClienterFake                 = HttpClienterDummy{}
	ResourceCreationFailureHttpClienterFake               = ResourceDeletionFailureHttpClienterFake{}
	ResourceUpdateFailureHttpClienterFake                 = ResourceDeletionFailureHttpClienterFake{}
	ResourceReadFailureHttpClienterFake                   = ResourceDeletionFailureHttpClienterFake{}
	HandleResponseEmptyReturnTemplateListHttpClienterFake = HttpClienterDummy{}
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
	return &resp, errCheckRedirectFailure
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
