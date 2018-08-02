package sewan_go_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"net/http"
)

//------------------------------------------------------------------------------
func (apier AirDrumResources_Apier) CreateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	sewan *API) (error, map[string]interface{}) {
	var (
		resourceInstanceCreationError error = nil
		createReqError                error = nil
		createError                   error = nil
		createRespBodyError           error = nil
		createdResource               map[string]interface{}
		resourceInstance              interface{}
		responseBody                  string
		instanceName                  string = d.Get(NAME_FIELD).(string)
		resourceJson                  []byte
		respBodyReader                interface{}
		bodyBytes                     []byte
	)
	apiTools := APITooler{
		Api: apier,
	}
	req := &http.Request{}
	resp := &http.Response{}
	resourceInstanceCreationError,
		resourceInstance = apiTools.Api.ResourceInstanceCreate(d,
		clientTooler,
		templatesTooler,
		resourceType,
		sewan)
	logger := loggerCreate("create_resource_" + instanceName + ".log")

	if resourceInstanceCreationError == nil {
		logger.Println("resourceInstance = ", resourceInstance)
		resourceJson, createReqError = json.Marshal(resourceInstance)
		if createReqError == nil {
			req, createReqError = http.NewRequest("POST",
				apiTools.Api.GetResourceCreationUrl(sewan, resourceType),
				bytes.NewBuffer(resourceJson))
			logger.Println("req.Body = ", req.Body)
			if createReqError == nil {
				req.Header.Add(HTTP_AUTHORIZATION, HTTP_TOKEN_HEADER+sewan.Token)
				req.Header.Add(HTTP_REQ_CONTENT_TYPE, HTTP_JSON_CONTENT_TYPE)
				resp, createReqError = clientTooler.Client.Do(sewan, req)
			}
		}

		if resp != nil {
			if createReqError != nil {
				createError = errors.New("Creation of \"" + instanceName +
					"\" failed, response reception error : " + createReqError.Error())
			} else {
				if resp.Body != nil {
					defer resp.Body.Close()
					bodyBytes, createRespBodyError = ioutil.ReadAll(resp.Body)
					responseBody = string(bodyBytes)
				} else {
					bodyBytes = []byte{}
					createRespBodyError = nil
					responseBody = ""
				}
				switch resp.Header.Get(HTTP_RESP_CONTENT_TYPE) {
				case HTTP_JSON_CONTENT_TYPE:
					resp_bodyJsonError := json.Unmarshal(bodyBytes, &respBodyReader)
					switch {
					case createRespBodyError != nil:
						createError = errors.New("Read of " + instanceName +
							" response body error " + createRespBodyError.Error())
					case resp_bodyJsonError != nil:
						createError = errors.New("Creation of \"" + instanceName +
							"\" failed, " +
							"the response body is not a properly formated json :\n\r\"" +
							resp_bodyJsonError.Error() + "\"")
					default:
						if resp.StatusCode == http.StatusCreated {
							createdResource = respBodyReader.(map[string]interface{})
							for key, value := range createdResource {
								readValue,
									updateError := schemaTools.SchemaTools.ReadElement(key,
									value,
									logger)
								if updateError == nil {
									createdResource[key] = readValue
								}
							}
						} else {
							createError = errors.New(resp.Status + responseBody)
						}
					}
				case HTTP_HTML_TEXT_CONTENT_TYPE:
					createError = errors.New(resp.Status + responseBody)
				default:
					createError = errors.New(ERROR_API_UNHANDLED_RESP_TYPE +
						resp.Header.Get(HTTP_RESP_CONTENT_TYPE) +
						ERROR_VALIDATE_API_URL)
				}
			}
		} else {
			createError = createReqError
		}
	} else {
		createError = resourceInstanceCreationError
	}
	logger.Println("createError = ", createError,
		"\ncreatedResource = ", createdResource)
	return createError, createdResource
}

//------------------------------------------------------------------------------
func (apier AirDrumResources_Apier) ReadResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	sewan *API) (error, map[string]interface{}, bool) {

	var (
		readError                     error = nil
		readReqError                  error = nil
		resourceInstanceCreationError error = nil
		readRespBodyError             error = nil
		read_resource                 map[string]interface{}
		responseBody                  string
		respBodyReader                interface{}
		resource_exists               bool   = true
		instanceName                  string = d.Get(NAME_FIELD).(string)
		bodyBytes                     []byte
	)
	req := &http.Request{}
	resp := &http.Response{}
	logger := loggerCreate("read_resource_" + instanceName + ".log")
	apiTools := APITooler{
		Api: apier,
	}
	resourceInstanceCreationError = apiTools.Api.ValidateResourceType(resourceType)

	if resourceInstanceCreationError == nil {
		req, readReqError = http.NewRequest("GET",
			apiTools.Api.GetResourceUrl(sewan, resourceType, d.Id()), nil)
		if readReqError == nil {
			req.Header.Add(HTTP_AUTHORIZATION, HTTP_TOKEN_HEADER+sewan.Token)
			resp, readReqError = clientTooler.Client.Do(sewan, req)
		}

		if resp != nil {
			if readReqError != nil {
				readError = errors.New(ERROR_READ_OF + instanceName +
					ERROR_UPDATE_STATE_FAILED_AND_RESP_ERR + readReqError.Error())
			} else {
				if resp.Body != nil {
					defer resp.Body.Close()
					bodyBytes, readRespBodyError = ioutil.ReadAll(resp.Body)
					responseBody = string(bodyBytes)
				} else {
					bodyBytes = []byte{}
					readRespBodyError = nil
					responseBody = ""
				}
				switch resp.Header.Get(HTTP_RESP_CONTENT_TYPE) {
				case HTTP_JSON_CONTENT_TYPE:
					switch {
					case readRespBodyError != nil:
						readError = errors.New("Read of " + instanceName +
							" state response body read error " + readRespBodyError.Error())
					case resp.StatusCode == http.StatusOK:
						resp_bodyJsonError := json.Unmarshal(bodyBytes, &respBodyReader)
						if resp_bodyJsonError != nil {
							readError = errors.New(ERROR_READ_OF + instanceName +
								ERROR_JSON_RESP_FAILED_AND_JSON_ERR +
								resp_bodyJsonError.Error() + "\"")
						} else {
							read_resource = respBodyReader.(map[string]interface{})

							for key, value := range read_resource {
								readValue,
									updateError := schemaTools.SchemaTools.ReadElement(key,
									value,
									logger)
								if updateError == nil {
									read_resource[key] = readValue
								}
							}
						}
					case resp.StatusCode == http.StatusNotFound:
						resource_exists = false
					default:
						readError = errors.New(resp.Status + responseBody)
					}
				case HTTP_HTML_TEXT_CONTENT_TYPE:
					readError = errors.New(resp.Status + responseBody)
				default:
					readError = errors.New(ERROR_API_UNHANDLED_RESP_TYPE +
						resp.Header.Get(HTTP_RESP_CONTENT_TYPE) +
						ERROR_VALIDATE_API_URL)
				}
			}
		} else {
			readError = readReqError
		}
	} else {
		readError = resourceInstanceCreationError
	}

	logger.Println("readError =", readError,
		"\nread_resource =", read_resource,
		"\nresource_exists =", resource_exists,
	)
	return readError, read_resource, resource_exists
}

//------------------------------------------------------------------------------
func (apier AirDrumResources_Apier) UpdateResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	sewan *API) error {

	var (
		resourceInstanceCreationError error = nil
		updateError                   error = nil
		updateReqError                error = nil
		updateRespBodyError           error = nil
		resourceInstance              interface{}
		responseBody                  string
		instanceName                  string = d.Get(NAME_FIELD).(string)
		resourceJson                  []byte
		respBodyReader                interface{}
		bodyBytes                     []byte
	)
	req := &http.Request{}
	resp := &http.Response{}
	apiTools := APITooler{
		Api: apier,
	}
	logger := loggerCreate("update_resource_" + instanceName + ".log")
	resourceInstanceCreationError,
		resourceInstance = apiTools.Api.ResourceInstanceCreate(d,
		clientTooler,
		templatesTooler,
		resourceType,
		sewan)

	if resourceInstanceCreationError == nil {

		resourceJson, updateReqError = json.Marshal(resourceInstance)
		if updateReqError == nil {
			req, updateReqError = http.NewRequest("PUT",
				apiTools.Api.GetResourceUrl(sewan, resourceType, d.Id()),
				bytes.NewBuffer(resourceJson))
			logger.Println("req.Body = ", req.Body)
			if updateReqError == nil {
				req.Header.Add(HTTP_AUTHORIZATION, HTTP_TOKEN_HEADER+sewan.Token)
				req.Header.Add(HTTP_REQ_CONTENT_TYPE, HTTP_JSON_CONTENT_TYPE)
				resp, updateReqError = clientTooler.Client.Do(sewan, req)
			}
		}

		if resp != nil {
			if updateReqError != nil {
				updateError = errors.New("Update of \"" + instanceName +
					ERROR_UPDATE_STATE_FAILED_AND_RESP_ERR + updateReqError.Error())
			} else {
				if resp.Body != nil {
					defer resp.Body.Close()
					bodyBytes, updateRespBodyError = ioutil.ReadAll(resp.Body)
					responseBody = string(bodyBytes)
				} else {
					bodyBytes = []byte{}
					updateRespBodyError = nil
					responseBody = ""
				}
				switch resp.Header.Get(HTTP_RESP_CONTENT_TYPE) {
				case HTTP_JSON_CONTENT_TYPE:
					switch {
					case updateRespBodyError != nil:
						updateError = errors.New(ERROR_READ_OF + instanceName +
							"\" state response body read error " + updateRespBodyError.Error())
					case resp.StatusCode == http.StatusOK:
						resp_bodyJsonError := json.Unmarshal(bodyBytes, &respBodyReader)
						if resp_bodyJsonError != nil {
							updateError = errors.New(ERROR_READ_OF + instanceName +
								ERROR_JSON_RESP_FAILED_AND_JSON_ERR +
								resp_bodyJsonError.Error())
						}
					default:
						updateError = errors.New(resp.Status + responseBody)
					}
				case HTTP_HTML_TEXT_CONTENT_TYPE:
					updateError = errors.New(resp.Status + responseBody)
				default:
					updateError = errors.New(ERROR_API_UNHANDLED_RESP_TYPE +
						resp.Header.Get(HTTP_RESP_CONTENT_TYPE) +
						ERROR_VALIDATE_API_URL)
				}
			}
		} else {
			updateError = updateReqError
		}

	} else {
		updateError = resourceInstanceCreationError
	}

	logger.Println("updateError = ", updateError)
	return updateError
}

//------------------------------------------------------------------------------
func (apier AirDrumResources_Apier) DeleteResource(d *schema.ResourceData,
	clientTooler *ClientTooler,
	templatesTooler *TemplatesTooler,
	schemaTools *SchemaTooler,
	resourceType string,
	sewan *API) error {

	var (
		resourceInstanceCreationError error = nil
		deleteError                   error = nil
		deleteReqError                error = nil
		deleteRespBodyError           error = nil
		responseBody                  string
		respBodyReader                interface{}
		bodyBytes                     []byte
		instanceName                  string = d.Get(NAME_FIELD).(string)
	)
	apiTools := APITooler{
		Api: apier,
	}
	resourceInstanceCreationError = apiTools.Api.ValidateResourceType(resourceType)
	req := &http.Request{}
	resp := &http.Response{}
	logger := loggerCreate("delete_resource_" + instanceName + ".log")
	logger.Println("--------------- ", instanceName, " ( id= ", d.Id(),
		") DELETE -----------------")

	if resourceInstanceCreationError == nil {

		req, deleteReqError = http.NewRequest("DELETE",
			apiTools.Api.GetResourceUrl(sewan, resourceType, d.Id()), nil)
		if deleteReqError == nil {
			req.Header.Add(HTTP_AUTHORIZATION, HTTP_TOKEN_HEADER+sewan.Token)
			resp, deleteReqError = clientTooler.Client.Do(sewan, req)
		}

		if resp != nil {
			if deleteReqError != nil {
				deleteError = errors.New("Deletion of \"" + instanceName +
					ERROR_UPDATE_STATE_FAILED_AND_RESP_ERR + deleteReqError.Error())
			} else {
				if resp.Body != nil {
					defer resp.Body.Close()
					bodyBytes, deleteRespBodyError = ioutil.ReadAll(resp.Body)
					responseBody = string(bodyBytes)
				} else {
					bodyBytes = []byte{}
					deleteRespBodyError = nil
					responseBody = ""
				}
				if resp.StatusCode != http.StatusNoContent {
					switch resp.Header.Get(HTTP_RESP_CONTENT_TYPE) {
					case HTTP_JSON_CONTENT_TYPE:
						switch {
						case deleteRespBodyError != nil:
							deleteError = errors.New("Deletion of " + instanceName +
								" response reception error : " + deleteRespBodyError.Error())
						default:
							resp_bodyJsonError := json.Unmarshal(bodyBytes, &respBodyReader)
							switch {
							case resp_bodyJsonError != nil:
								deleteError = errors.New(ERROR_READ_OF + instanceName +
									ERROR_JSON_RESP_FAILED_AND_JSON_ERR +
									resp_bodyJsonError.Error())
							default:
								deleteError = errors.New(resp.Status + responseBody)
							}
						}
					case HTTP_HTML_TEXT_CONTENT_TYPE:
						deleteError = errors.New(resp.Status + responseBody)
					default:
						deleteError = errors.New(ERROR_API_UNHANDLED_RESP_TYPE +
							resp.Header.Get(HTTP_RESP_CONTENT_TYPE) +
							ERROR_VALIDATE_API_URL)
					}
				}
			}
		} else {
			deleteError = deleteReqError
		}
	} else {
		deleteError = resourceInstanceCreationError
	}

	logger.Println("deleteError = ", deleteError)
	return deleteError
}
