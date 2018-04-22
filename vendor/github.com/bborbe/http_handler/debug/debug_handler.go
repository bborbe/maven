package debug

import (
	"net/http"
	"time"

	"github.com/golang/glog"
)

type handler struct {
	subhandler http.Handler
}

func New(subhandler http.Handler) *handler {
	h := new(handler)
	h.subhandler = subhandler
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	start := time.Now()
	defer glog.V(4).Infof("%s %s takes %dms", request.Method, request.URL.Path, time.Now().Sub(start)/time.Millisecond)

	glog.V(4).Infof("request %v: ", request)
	h.subhandler.ServeHTTP(responseWriter, request)
	glog.V(4).Infof("response %v: ", responseWriter)
}
