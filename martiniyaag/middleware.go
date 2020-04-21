package martiniyaag

import (
	"github.com/dshechao/gen-api/gen"
	"github.com/dshechao/gen-api/gen/models"
	"github.com/dshechao/gen-api/middleware"
	"github.com/go-martini/martini"
	"net/http"
)

func Document(c martini.Context, w http.ResponseWriter, r *http.Request) {
	if !gen.IsOn() {
		c.Next()
		return
	}
	apiCall := models.ApiCall{}
	writer := middleware.NewResponseRecorder(w)
	c.MapTo(writer, (*http.ResponseWriter)(nil))
	middleware.Before(&apiCall, r)
	c.Next()
	middleware.After(&apiCall, writer, r)
}
