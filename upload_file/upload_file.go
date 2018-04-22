package upload_file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
)

type handler struct {
	root string
}

func New(root string) *handler {
	h := new(handler)
	h.root = root
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if err := h.serveHTTP(responseWriter, request); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) serveHTTP(responseWriter http.ResponseWriter, request *http.Request) error {
	target, err := target(h.root, request.RequestURI)
	if err != nil {
		return err
	}
	glog.V(2).Infof("upload file to: %s", target)

	err = os.MkdirAll(path.Dir(target), 0755)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(f, request.Body)
	if err != nil {
		return err
	}
	glog.V(2).Infof("completed upload to %s", target)
	return nil
}

func target(root string, requesturi string) (string, error) {
	target := path.Join(root, requesturi)
	if strings.Index(target, root) != 0 {
		return "", fmt.Errorf("illegal target: %s", target)
	}
	return target, nil
}
