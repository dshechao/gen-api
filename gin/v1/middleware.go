package gin

import (
	"log"
	"strings"

	"github.com/dshechao/gen-api/gen"
	"github.com/dshechao/gen-api/gen/models"
	"github.com/dshechao/gen-api/middleware"
)

func Document() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !gen.IsOn() {
			return
		}
		apiCall := models.ApiCall{}
		middleware.Before(&apiCall, c.Request)
		c.Next()
		if gen.IsStatusCodeValid(c.Writer.Status()) {
			apiCall.MethodType = c.Request.Method
			apiCall.CurrentPath = strings.Split(c.Request.RequestURI, "?")[0]
			apiCall.ResponseBody = ""
			apiCall.ResponseCode = c.Writer.Status()
			headers := map[string]string{}
			for k, v := range c.Writer.Header() {
				log.Println(k, v)
				headers[k] = strings.Join(v, " ")
			}
			apiCall.ResponseHeader = headers
			go gen.GenerateHtml(&apiCall)
		}
	}
}
