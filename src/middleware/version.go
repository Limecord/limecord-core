package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"limecord-core/utils"
	"strconv"
	"strings"
)

const (
	// Minimum API version
	MIN_VERSION = 6
	// Maximum API version
	MAX_VERSION = 8
)

func APIVersionMiddleware(ctx *gin.Context) {
	var (
		version uint64 = 0
		err error = nil
	)

	// grab the version from the url
	rawVersion := ctx.Param("version")
	// make sure the string is prefixed with v
	if !strings.HasPrefix(rawVersion, "v") {
		// this isn't really required, just that error cant be nil
		err = fmt.Errorf("invalid api version")
	} else {
		// remove the v from the string and parse out the number
		// we only want to read 8-bits because we wont have an api version higher than 255
		version, err = strconv.ParseUint(rawVersion[1:], 10, 8)
	}

	// make sure we haven't ran into any errors parsing the version
	if err != nil {
		// if we couldn't parse the version just return a 404 error
		ctx.AbortWithStatusJSON(404, utils.GetAPIError(utils.API_NOTFOUND, utils.API_NOTFOUND_MESSAGE))
		// check the api version against the min and max version, and output if not matching
	} else if version > MAX_VERSION || version < MIN_VERSION {
		ctx.AbortWithStatusJSON(400, utils.GetAPIError(utils.API_INVALID_VERSION, utils.API_INVALID_VERSION_MESSAGE))
	}

	// if we made it here then we will continue onto the next handler
}
