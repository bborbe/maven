package main

import (
	"fmt"
	"net/http"
	"os"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/log"
	"github.com/bborbe/server/handler/debug"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
	"github.com/bborbe/server/handler/static"
)

const (
	PARAMETER_ROOT = "root"
	PARAMETER_LOGLEVEL = "loglevel"
)

var (
	logger = log.DefaultLogger
	portPtr = flag.Int("port", 8080, "Port")
	logLevelPtr = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	documentRootPtr = flag.String(PARAMETER_ROOT, "", "Document root directory")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	server, err := createServer(*portPtr, *documentRootPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debugf("start server")
	gracehttp.Serve(server)
}

func createServer(port int, root string) (*http.Server, error) {

	router := mux.NewRouter()
	router.Methods("GET").Handler(http.FileServer(http.Dir(root)))
	router.Methods("PUT").Handler(static.NewHandlerStaticContent("ok"))

	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: debug.New(router)}, nil
}
