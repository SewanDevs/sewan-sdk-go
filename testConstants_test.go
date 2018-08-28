package sewan_go_sdk

const (
	rightAPIURL                = "https://apiURL/api/bonjour/"
	rightVMCreationAPIURL      = "https://apiURL/api/bonjour/vm/"
	rightVMURLPatate           = "https://apiURL/api/bonjour/vm/PATATE/"
	rightVMURLQuaranteDeux     = "https://apiURL/api/bonjour/vm/42/"
	wrongAPIURL                = "a wrong url"
	wrongAPIURLError           = "Wrong api url msg"
	noRespAPIURL               = "https://noRespAPIURL.fr"
	noRespBodyAPIURL           = "https://NO_BODY_API_URL.org"
	noRespJSONBodyAPIURL       = "https://tata.fr"
	rightAPIToken              = "42424242424242424242424242424242"
	wrongAPIToken              = "a wrong token"
	wrongTokenError            = "Wrong api token msg"
	errorAPIUnhandledImageType = "Unhandled api response type : image"
	notFoundRespStatus         = "404 Not Found"
	notFoundRespMsg            = "404 Not Found{\"detail\":\"Not found.\"}"
	unauthorizedStatus         = "401 Unauthorized"
	unauthorizedMsg            = "401 Unauthorized{\"detail\":\"Token non valide.\"}"
	destroyWrongMsg            = "{\"detail\":\"Destroying resource wrong body message\"}"
	vdcDestroyFailureMsg       = "Destroying the VDC now"
	vmDestroyFailureMsg        = "Destroying the VM now"
	wrongResourceType          = "a non supported ResourceType"
	enterpriseSlug             = "unit test enterprise"
)
