package utils

import "github.com/gin-gonic/gin"

const (
	// API Error Codes
	API_NOTFOUND          = 0
	API_UNAUTHORIZED      = 0
	API_INVALID_VERSION   = 50041
	API_INVALID_FORM_BODY = 50035

	// API Error Messages
	API_NOTFOUND_MESSAGE          = "404: Not Found"
	API_UNAUTHORIZED_MESSAGE      = "404: Unauthorized"
	API_INVALID_VERSION_MESSAGE   = "Invalid API version"
	API_INVALID_FORM_BODY_MESSAGE = "Invalid Form Body"
)

func GetAPIError(code uint, message string) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
	}
}
