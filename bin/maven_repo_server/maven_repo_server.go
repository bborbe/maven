package main

import (
	"fmt"
	"net/http"

	debug_handler "github.com/bborbe/http_handler/debug"

	"runtime"

	flag "github.com/bborbe/flagenv"
	io_util "github.com/bborbe/io/util"
	"github.com/bborbe/maven_repo/model"
	"github.com/bborbe/maven_repo/upload_file"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

const (
	PARAMETER_ROOT = "root"
	PARAMETER_PORT = "port"
)

var (
	portPtr         = flag.Int(PARAMETER_PORT, 8080, "Port")
	documentRootPtr = flag.String(PARAMETER_ROOT, "", "Document root directory")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	server, err := createServer()
	if err != nil {
		return err
	}

	glog.V(2).Infof("start server")
	return gracehttp.Serve(server)
}

func createServer() (*http.Server, error) {
	port := model.Port(*portPtr)
	root := *documentRootPtr
	if port <= 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_PORT)
	}
	if len(root) == 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_ROOT)
	}
	root, err := io_util.NormalizePath(root)
	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("root dir: %s", root)
	handler := createHandler(root)

	if glog.V(4) {
		handler = debug_handler.New(handler)
	}

	glog.V(2).Infof("create http server on %s", port.Address())
	return &http.Server{Addr: port.Address(), Handler: handler}, nil
}

func createHandler(root string) http.Handler {

	handler := mux.NewRouter()
	handler.Methods("GET").Handler(http.FileServer(http.Dir(root)))
	handler.Methods("PUT").Handler(upload_file.New(root))
	return handler
}
