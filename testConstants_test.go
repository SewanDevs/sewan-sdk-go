package sewan_go_sdk

const (
	rightApiUrl                = "https://api_url/api/bonjour/"
	rightVmCreationApiUrl      = "https://api_url/api/bonjour/vm/"
	rightVmUrlPatate           = "https://api_url/api/bonjour/vm/PATATE/"
	rightVmUrlQuaranteDeux     = "https://api_url/api/bonjour/vm/42/"
	wrongApiUrl                = "a wrong url"
	wrongApiUrlError           = "Wrong api url msg"
	noRespApiUrl               = "https://noRespApiUrl.fr"
	noRespBodyApiUrl           = "https://NO_BODY_API_URL.org"
	noRespJsonBodyApiUrl       = "https://tata.fr"
	rightApiToken              = "42424242424242424242424242424242"
	wrongApiToken              = "a wrong token"
	wrongTokenError            = "Wrong api token msg"
	errorApiUnhandledImageType = "Unhandled api response type : image"
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
