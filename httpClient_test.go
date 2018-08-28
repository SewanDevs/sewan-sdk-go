package sewansdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
)

func TestDo(t *testing.T) {
	//Not tested, ref=TD-35489-UT-35737-1
}

func TestgetTemplatesList(t *testing.T) {
	testCases := []struct {
		ID             int
		TCClienter     Clienter
		EnterpriseSlug string
		TemplateList   []interface{}
		Error          error
	}{
		{
			1,
			getTemplatesListSuccessHTTPClienterFake{},
			"unit test enterprise",
			templatesList,
			nil,
		},
		{
			2,
			getTemplatesListFailureHTTPClienterFake{},
			"unit test enterprise",
			nil,
			errors.New("handleResponse() error"),
		},
		{
			3,
			ErrorResponseHTTPClienterFake{},
			"unit test enterprise",
			nil,
			errDoRequest,
		},
		{
			4,
			handleResponseEmptyReturnTemplateListHTTPClienterFake,
			"unit test enterprise",
			nil,
			errEmptyTemplateList,
		},
	}
	clientTooler := ClientTooler{}
	clientTooler.Client = HTTPClienter{}
	fakeClientTooler := ClientTooler{}
	apiTooler := APITooler{}
	api := apiTooler.New(TokenField, "url")
	for _, testCase := range testCases {
		fakeClientTooler.Client = testCase.TCClienter
		templatesList, err := clientTooler.Client.getTemplatesList(&fakeClientTooler,
			testCase.EnterpriseSlug, api)
		switch {
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : getTemplatesList() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.Error)
			} else {
				diffs := cmp.Diff(testCase.TemplateList, templatesList)
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong template list (-got +want) :\n%s",
						testCase.ID, diffs)
				}
			}
		case err != nil && testCase.Error != nil:
			if templatesList != nil {
				t.Errorf("\n\nTC %d : Wrong response read element,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.ID, templatesList, testCase.TemplateList)
			}
			if err.Error() != testCase.Error.Error() {
				t.Errorf("\n\nTC %d : Wrong response handle error,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err.Error(), testCase.Error.Error())
			}
		}
	}
}

func TestHandleResponse(t *testing.T) {
	testCases := []struct {
		ID                 int
		Response           *http.Response
		ExpectedCode       int
		ExpectedBodyFormat string
		ResponseBody       interface{}
		Error              error
	}{
		{
			1,
			HTTPResponseFakeOKJSON(),
			http.StatusOK,
			httpJSONContentType,
			JSONStub(),
			nil,
		},
		{
			2,
			HTTPResponseFakeOKTemplateListJSON(),
			http.StatusOK,
			httpJSONContentType,
			JSONTemplateListFake(),
			nil,
		},
		{
			3,
			HTTPResponseFake500Texthtml(),
			http.StatusInternalServerError,
			httpHTMLTextContentType,
			"<h1>Server Error (500)</h1>",
			nil,
		},
		{
			4,
			HTTPResponseFake500Json(),
			http.StatusInternalServerError,
			httpHTMLTextContentType,
			nil,
			errors.New("Wrong response content type," +
				"\nexpected :text/html\ngot :application/json"),
		},
		{
			5,
			HTTPResponseFakeOKJSON(),
			http.StatusInternalServerError,
			httpHTMLTextContentType,
			nil,
			errors.New("Wrong response status code,\nexpected :500\ngot :200\n" +
				"Full response status : 200 OK" +
				"\nWrong response content type," +
				"\nexpected :text/html\ngot :application/json"),
		},
		{
			6,
			HTTPResponseFakeOKJSON(),
			http.StatusInternalServerError,
			httpHTMLTextContentType,
			nil,
			errors.New("Wrong response status code,\nexpected :500\ngot :200" +
				"\nFull response status : 200 OK" +
				"\nWrong response content type," +
				"\nexpected :text/html\ngot :application/json"),
		},
		{
			7,
			HTTPResponseFakeOkNilBody(),
			http.StatusOK,
			"",
			"",
			errEmptyRespBody,
		},
		{
			8,
			HTTPResponseFakeOKWrongjson(),
			http.StatusOK,
			httpJSONContentType,
			nil,
			errors.New(errJSONFormat +
				"invalid character 'a' looking for beginning of value" +
				"\nJson :a bad formated json"),
		},
		{
			9,
			HTTPResponseFakeOKImage(),
			http.StatusOK,
			"image",
			nil,
			errors.New(errorAPIUnhandledImageType +
				errValidateAPIURL),
		},
	}
	clienter := HTTPClienter{}
	for _, testCase := range testCases {
		responseBody, err := clienter.handleResponse(testCase.Response,
			testCase.ExpectedCode, testCase.ExpectedBodyFormat)
		switch {
		case err == nil || testCase.Error == nil:
			if !(err == nil && testCase.Error == nil) {
				t.Errorf("\n\nTC %d : handleResponse() error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.Error)
			} else {
				diffs := cmp.Diff(testCase.ResponseBody, responseBody)
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong response read element (-got +want) :\n%s",
						testCase.ID, diffs)
				}
			}
		case err != nil && testCase.Error != nil:
			if (responseBody != nil) && (responseBody != "") {
				t.Errorf("\n\nTC %d : Wrong response read element,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.ID, responseBody, testCase.ResponseBody)
			}
			if err.Error() != testCase.Error.Error() {
				t.Errorf("\n\nTC %d : Wrong response handle error,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err.Error(), testCase.Error.Error())
			}
		}
	}
}
