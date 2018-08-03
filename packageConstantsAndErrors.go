package sewan_go_sdk

import (
	"errors"
	"net/http"
	"strconv"
)

const (
	NameField                      = "name"
	EnterpriseField                = "enterprise"
	DatacenterField                = "datacenter"
	VdcResourceField               = "vdc_resources"
	ResourceField                  = "resource"
	TotalField                     = "total"
	UsedField                      = "used"
	SlugField                      = "slug"
	StateField                     = "state"
	OsField                        = "os"
	RamField                       = "ram"
	CpuField                       = "cpu"
	DisksField                     = "disks"
	VDiskField                     = "v_disk"
	SizeField                      = "size"
	StorageClassField              = "storage_class"
	NicsField                      = "nics"
	VlanNameField                  = "vlan"
	MacAdressField                 = "mac_address"
	ConnectedField                 = "connected"
	VdcField                       = "vdc"
	BootField                      = "boot"
	TokenField                     = "token"
	BackupField                    = "backup"
	DiskImageField                 = "disk_image"
	PlatformNameField              = "platform_name"
	BackupSizeField                = "backup_size"
	CommentField                   = "comment"
	TemplateField                  = "template"
	IdField                        = "id"
	DynamicField                   = "dynamic_field"
	OutsourcingField               = "outsourcing"
	monoField                      = "-mono-"
	InstanceNumberField            = "instance_number"
	VmResourceType                 = "vm"
	VdcResourceType                = VdcField
	resourceNameCountSeparator     = "-"
	resourceDynamicInstanceNumber  = "${count.index + 1}"
	httpReqContentType             = "content-type"
	httpRespContentType            = "Content-Type"
	httpJsonContentType            = "application/json"
	httpHtmlTextContentType        = "text/html"
	httpAuthorization              = "authorization"
	httpTokenHeader                = "Token "
	errTestResultDiffs             = "\n\rGot: \"%s\"\n\rWant: \"%s\""
	errApiUnhandledRespType        = "Unhandled api response type : "
	errValidateApiUrl              = "\nPlease validate the configuration api url."
	errReadOf                      = "Read of \""
	errUpdateStateFailedAndRespErr = "\" state failed, response reception error : "
	errJsonRespFailedAndJsonErr    = "\" failed, response body json error :\n\r\""
	errApiDownOrWrongApiUrl        = "\", the api is down or this url is wrong."
	errEmptyResponse               = "Empty response error."
	errJsonFormat                  = "Response body is not a properly formated json :"
	creationOperation              = "Creation"
	readOperation                  = "Read"
	updateOperation                = "Update"
	deleteOperation                = "Delete"
)

var (
	ErrEmptyResp                 = errors.New("Empty API response.")
	ErrEmptyRespBody             = errors.New("Empty API response body.")
	ErrResourceNotExist          = errors.New("Resource does not exists.")
	ErrUninitializedExpectedCode = errors.New("Expected code not initialized.")
	ErrNilResponse               = errors.New("Response is nil.")
	ErrZeroStatusCode            = errors.New("Response status code is zero.")
)

func ErrRespStatusCodeBuilder(resp *http.Response,
	expectedCode int,
	additionalErrMsg string) error {
	if expectedCode == 0 {
		return ErrUninitializedExpectedCode
	}
	if resp == nil {
		return ErrNilResponse
	}
	if resp.StatusCode == 0 {
		return ErrZeroStatusCode
	}
	if expectedCode == resp.StatusCode {
		if additionalErrMsg == "" {
			return nil
		} else {
			return errors.New(additionalErrMsg)
		}
	}
	return errors.New("Wrong response status code, \n\r expected :" +
		strconv.Itoa(expectedCode) + "\n\r got :" + strconv.Itoa(resp.StatusCode) +
		"\nFull response status : " + resp.Status + "\n" + additionalErrMsg)
}

func ErrCrudRequestsBuilder(crudOperation string,
	instanceName string,
	err error) error {
	of := "of \""
	postMsg := "\" failed, POST response reception error : "
	getMsg := "\" failed, GET response reception error : "
	deleteMsg := "\" failed, DELETE response reception error : "
	if instanceName == "" {
		return errors.New("instanceName is empty string.")
	}
	if err == nil {
		return errors.New("Request execution error is nil.")
	}
	switch crudOperation {
	case creationOperation:
		return errors.New(creationOperation + of + instanceName +
			postMsg + err.Error())
	case readOperation:
		return errors.New(creationOperation + of + instanceName +
			getMsg + err.Error())
	case updateOperation:
		return errors.New(updateOperation + of + instanceName +
			postMsg + err.Error())
	case deleteOperation:
		return errors.New(creationOperation + of + instanceName +
			deleteMsg + err.Error())
	default:
		return errors.New(crudOperation + "is not a crudOperation from list :" +
			creationOperation + readOperation + updateOperation + deleteOperation)
	}
}
