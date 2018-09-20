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
func (client HTTPClienterDummy) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
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
func (client ErrorResponseHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
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
func (client getJSONListFailureHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
	return nil, errEmptyResourcesList(clouddcUniTestResource)
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
func (client getListSuccessHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
	return unitTestMetaDataList, nil
}
func (client getListSuccessHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return templateMetaDataList, nil
}

type HandleRespErrHTTPClienterFake struct{}

func (client HandleRespErrHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client HandleRespErrHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
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
func (client Error404HTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
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
func (client VMCreationSuccessHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
	return nil, nil
}
func (client VMCreationSuccessHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return noTemplateVMMap, nil
}

type VDCCreationSuccessHTTPClienterFake struct{}

func (client VDCCreationSuccessHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client VDCCreationSuccessHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
	return nil, nil
}
func (client VDCCreationSuccessHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return vdcCreationResponseMap, nil
}

type VMReadSuccessHTTPClienterFake struct{}

func (client VMReadSuccessHTTPClienterFake) do(api *API,
	req *http.Request) (*http.Response, error) {
	return nil, nil
}
func (client VMReadSuccessHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
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
func (client VdcReadSuccessHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
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
func (client ResourceDeletionFailureHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
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
func (client CheckRedirectReqFailureHTTPClienterFake) getEnvResourceList(clientTooler *ClientTooler,
	api *API, resourceType string) ([]interface{}, error) {
	return nil, nil
}
func (client CheckRedirectReqFailureHTTPClienterFake) handleResponse(resp *http.Response,
	expectedCode int,
	expectedBodyFormat string) (interface{}, error) {
	return nil, nil
}
