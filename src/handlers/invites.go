package handlers

import "github.com/gin-gonic/gin"

func generic_invite(c *gin.Context) {
	invite := c.Param("invite")

	if invite == "test" {
		c.JSON(200, gin.H{
			"code": "test",
			"guild": gin.H{
				"id":                 "744930173069164639",
				"name":               "Limecord test",
				"description":        "test",
				"icon":               nil,
				"features":           []string{"COMMUNITY", "VANITY_URL"},
				"verification_level": 0,
				"vanity_url_code":    "test",
			},
			"channel": gin.H{
				"id":   "783972148175437854",
				"name": "verify",
				"type": 0,
			},
			"approximate_member_count":   0,
			"approximate_presence_count": 0,
		})
	} else {
		c.AbortWithStatus(400)
	}
}

// register the routers for the tracking module
func RegisterInvites(router *gin.RouterGroup) {
	router = router.Group("/invites")
	router.GET("/:invite", generic_invite)
}
