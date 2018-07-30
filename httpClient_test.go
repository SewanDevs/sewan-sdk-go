package sewan_go_sdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
)

func TestDo(t *testing.T) {
	//Not tested, ref=TD-35489-UT-35737-1
}

//------------------------------------------------------------------------------
func TestGetTemplatesList(t *testing.T) {
	testCases := []struct {
		Id              int
		TC_clienter     Clienter
		Enterprise_slug string
		TemplateList    []interface{}
		Error           error
	}{
		{
			1,
			GetTemplatesList_Success_HttpClienterFake{},
			"unit test enterprise",
			TEMPLATES_LIST,
			nil,
		},
		{
			2,
			GetTemplatesList_Failure_HttpClienterFake{},
			"unit test enterprise",
			nil,
			errors.New("HandleResponse() error"),
		},
		{
			3,
			ErrorResponse_HttpClienterFake{},
			"unit test enterprise",
			nil,
			errors.New(REQ_ERR),
		},
	}
	var (
		templatesList []interface{}
		err           error
		diffs         string
	)
	client_tooler := ClientTooler{}
	client_tooler.Client = HttpClienter{}
	fake_client_tooler := ClientTooler{}
	apiTooler := APITooler{}
	api := apiTooler.New(TOKEN_FIELD, "url")

	for _, testCase := range testCases {
		fake_client_tooler.Client = testCase.TC_clienter
		templatesList, err = client_tooler.Client.GetTemplatesList(&fake_client_tooler,
			testCase.Enterprise_slug, api)
		switch {
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : GetTemplatesList() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.Id, err, testCase.Error)
			} else {
				diffs = cmp.Diff(testCase.TemplateList, templatesList)
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong template list (-got +want) :\n%s",
						testCase.Id, diffs)
				}
			}
		case err != nil && testCase.Error != nil:
			if templatesList != nil {
				t.Errorf("\n\nTC %d : Wrong response read element,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.Id, templatesList, testCase.TemplateList)
			}
			if err.Error() != testCase.Error.Error() {
				t.Errorf("\n\nTC %d : Wrong response handle error,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Error.Error())
			}
		}
	}
}

//------------------------------------------------------------------------------
func TestHandleResponse(t *testing.T) {
	testCases := []struct {
		Id                 int
		Response           *http.Response
		ExpectedCode       int
		ExpectedBodyFormat string
		ResponseBody       interface{}
		Error              error
	}{
		{
			1,
			HttpResponseFake_OKJson(),
			http.StatusOK,
			"application/json",
			JsonStub(),
			nil,
		},
		{
			2,
			HttpResponseFake_OKTemplateListJson(),
			http.StatusOK,
			"application/json",
			JsonTemplateListFake(),
			nil,
		},
		{
			3,
			HttpResponseFake_500_texthtml(),
			http.StatusInternalServerError,
			"text/html",
			"<h1>Server Error (500)</h1>",
			nil,
		},
		{
			4,
			HttpResponseFake_500Json(),
			http.StatusInternalServerError,
			"text/html",
			nil,
			errors.New("Wrong response content type, \n\r expected :text/html\n\r got :application/json"),
		},
		{
			5,
			HttpResponseFake_OKJson(),
			http.StatusInternalServerError,
			"text/html",
			nil,
			errors.New("Wrong response status code, \n\r expected :500\n\r got :200" +
				"\n\rFull response status : 200 OK"),
		},
		{
			6,
			HttpResponseFake_OKJson(),
			http.StatusInternalServerError,
			"text/html",
			nil,
			errors.New("Wrong response status code, \n\r expected :500\n\r got :200" +
				"\n\rFull response status : 200 OK"),
		},
		{
			7,
			HttpResponseFake_OK_no_content(),
			http.StatusOK,
			"",
			nil,
			nil,
		},
		{
			8,
			HttpResponseFake_OK_wrongjson(),
			http.StatusOK,
			"application/json",
			nil,
			errors.New("Response body is not a properly formated json :" +
				"invalid character 'a' looking for beginning of value"),
		},
		{
			9,
			HttpResponseFake_OK_image(),
			http.StatusOK,
			"image",
			nil,
			errors.New("Unhandled api response type : image" +
				"\nPlease validate the configuration api url."),
		},
	}

	var (
		responseBody interface{}
		err          error
		diffs        string
	)
	clienter := HttpClienter{}

	for _, testCase := range testCases {
		responseBody, err = clienter.HandleResponse(testCase.Response,
			testCase.ExpectedCode, testCase.ExpectedBodyFormat)

		switch {
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : HandleResponse() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.Id, err, testCase.Error)
			} else {
				diffs = cmp.Diff(testCase.ResponseBody, responseBody)
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong response read element (-got +want) :\n%s",
						testCase.Id, diffs)
				}
			}
		case err != nil && testCase.Error != nil:
			if responseBody != nil {
				t.Errorf("\n\nTC %d : Wrong response read element,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.Id, responseBody, testCase.ResponseBody)
			}
			if err.Error() != testCase.Error.Error() {
				t.Errorf("\n\nTC %d : Wrong response handle error,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.Id, err.Error(), testCase.Error.Error())
			}
		}
	}
}
