package handlers

import (
	"limecord-core/utils"

	"github.com/gin-gonic/gin"
)

func location_metadata(c *gin.Context) {
	c.JSON(200, gin.H{
		"consent_required": false,
		"country_code":     "US",
	})
}

func login(c *gin.Context) {

}

type Register struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Invite   string `json:"invite"`

	CaptchaKey    string `json:"captcha_key"`
	Fingerprint   string `json:"fingerprint"`
	Consent       bool   `json:"consent"`
	GiftCodeSkuId string `json:"gift_code_sku_id"`
}

func register(c *gin.Context) {
	var body Register
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(400, utils.GetAPIError(utils.API_INVALID_FORM_BODY, utils.API_INVALID_FORM_BODY_MESSAGE))
	}

	// validate fields

	c.JSON(200, gin.H{
		"token": body.Username,
	})
}

// register the routers for the tracking module
func RegisterAuth(router *gin.RouterGroup) {
	router = router.Group("/auth")
	router.GET("/location-metadata", location_metadata)
	router.POST("/register", register)
}
