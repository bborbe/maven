package main

import (
	"fmt"
	"net/http"
	"os"

	debug_handler "github.com/bborbe/http_handler/debug"

	"runtime"

	flag "github.com/bborbe/flagenv"
	io_util "github.com/bborbe/io/util"
	"github.com/bborbe/log"
	"github.com/bborbe/maven_repo/upload_file"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
)

const (
	PARAMETER_ROOT     = "root"
	PARAMETER_PORT     = "port"
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_DEBUG    = "debug"
)

var (
	logger          = log.DefaultLogger
	portPtr         = flag.Int(PARAMETER_PORT, 8080, "Port")
	logLevelPtr     = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	documentRootPtr = flag.String(PARAMETER_ROOT, "", "Document root directory")
	debugPtr        = flag.Bool(PARAMETER_DEBUG, false, "debug")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	server, err := createServer(
		*portPtr,
		*debugPtr,
		*documentRootPtr,
	)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debugf("start server")
	gracehttp.Serve(server)
}

func createServer(
	port int,
	debug bool,
	root string,
) (*http.Server, error) {
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

	logger.Debugf("root dir: %s", root)

	handler := mux.NewRouter()
	handler.Methods("GET").Handler(http.FileServer(http.Dir(root)))
	handler.Methods("PUT").Handler(upload_file.New(root))

	if debug {
		handler = debug_handler.New(handler)
	}

	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
