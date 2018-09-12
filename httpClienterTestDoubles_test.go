package sewansdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HTTPClienterDummy struct{}

func (client HTTPClienterDummy) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client HTTPClienterDummy) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client HTTPClienterDummy) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client HTTPClienterDummy) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

var (
	ResourceDeletionSuccessHTTPClienterFake               = HTTPClienterDummy{}
	ResourceUpdateSuccessHTTPClienterFake                 = HTTPClienterDummy{}
	ResourceCreationFailureHTTPClienterFake               = ResourceDeletionFailureHTTPClienterFake{}
	ResourceUpdateFailureHTTPClienterFake                 = ResourceDeletionFailureHTTPClienterFake{}
	ResourceReadFailureHTTPClienterFake                   = ResourceDeletionFailureHTTPClienterFake{}
	handleResponseEmptyReturnTemplateListHTTPClienterFake = HTTPClienterDummy{}
)

type ErrorResponseHTTPClienterFake struct{}

func (client ErrorResponseHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, errDoRequest
}
func (client ErrorResponseHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client ErrorResponseHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client ErrorResponseHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

type getJSONListFailureHTTPClienterFake struct{}

func (client getJSONListFailureHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client getJSONListFailureHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, errEmptyTemplateList
}
func (client getJSONListFailureHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, errEmptyResourcesList
}
func (client getJSONListFailureHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

type getListSuccessHTTPClienterFake struct{}

func (client getListSuccessHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client getListSuccessHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return templatesList, nil
}
func (client getListSuccessHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return resourceMetaDataList, nil
}
func (client getListSuccessHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return templatesList, nil
}

type HandleRespErrHTTPClienterFake struct{}

func (client HandleRespErrHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client HandleRespErrHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client HandleRespErrHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client HandleRespErrHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, errHandleResponse
}

type Error404HTTPClienterFake struct{}

func (client Error404HTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	resp := http.Response{}
	resp.Header = map[string][]string{}
	resp.Header.Add(httpRespContentType, httpJSONContentType)
	resp.StatusCode = http.StatusNotFound
	resp.Status = notFoundRespStatus
	body := RespBody{"Not found."}
	js, _ := json.Marshal(body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(js))
	return &resp, nil
}
func (client Error404HTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client Error404HTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client Error404HTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

type VMCreationSuccessHTTPClienterFake struct{}

func (client VMCreationSuccessHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client VMCreationSuccessHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VMCreationSuccessHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VMCreationSuccessHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return noTemplateVMMap, nil
}

type VMReadSuccessHTTPClienterFake struct{}

func (client VMReadSuccessHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client VMReadSuccessHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VMReadSuccessHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VMReadSuccessHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return noTemplateVMMap, nil
}

type VdcReadSuccessHTTPClienterFake struct{}

func (client VdcReadSuccessHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client VdcReadSuccessHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VdcReadSuccessHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client VdcReadSuccessHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return vdcReadResponseMap, nil
}

type ResourceDeletionFailureHTTPClienterFake struct{}

func (client ResourceDeletionFailureHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, errEmptyResp
}
func (client ResourceDeletionFailureHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client ResourceDeletionFailureHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client ResourceDeletionFailureHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}

type CheckRedirectReqFailureHTTPClienterFake struct{}

func (client CheckRedirectReqFailureHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	resp := http.Response{}
	return &resp, errCheckRedirectFailure
}
func (client CheckRedirectReqFailureHTTPClienterFake) getTemplatesList(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client CheckRedirectReqFailureHTTPClienterFake) getPhysicalResourcesMeta(clientTooler *ClientTooler,
	api *API) ([]interface{}, error) {
	return nil, nil
}
func (client CheckRedirectReqFailureHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}
