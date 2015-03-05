package server

import (
	"encoding/json"
	"errors"
	log "github.com/funkygao/log4go"
	"github.com/gorilla/mux"
	"io"
	"net"
	"net/http"
	"time"

	_ "expvar"         // localhost:xx/debug/vars
	_ "net/http/pprof" // localhost:xx/debug/pprof
)

var (
	httpApi       *httpRestApi
	httpDupLaunch = errors.New("server.LaunchHttpServer already called")
	ErrHttp404    = errors.New("Not found")
)

type httpRestApi struct {
	httpListener net.Listener
	httpServer   *http.Server
	httpRouter   *mux.Router
	httpPaths    []string
}

func LaunchHttpServer(listenAddr string, debugAddr string) (err error) {
	if httpApi != nil {
		return httpDupLaunch
	}

	httpApi = new(httpRestApi)
	httpApi.httpPaths = make([]string, 0, 10)
	httpApi.httpRouter = mux.NewRouter()
	httpApi.httpServer = &http.Server{
		Addr:    listenAddr,
		Handler: httpApi.httpRouter,
	}

	httpApi.httpListener, err = net.Listen("tcp", httpApi.httpServer.Addr)
	if err != nil {
		httpApi = nil
		return err
	}

	if debugAddr != "" {
		log.Debug("HTTP serving at %s with pprof at %s", listenAddr, debugAddr)
	} else {
		log.Debug("HTTP serving at %s", listenAddr)
	}

	go httpApi.httpServer.Serve(httpApi.httpListener)
	if debugAddr != "" {
		go http.ListenAndServe(debugAddr, nil)
	}

	return nil
}

func StopHttpServer() {
	if httpApi != nil && httpApi.httpListener != nil {
		httpApi.httpListener.Close()
		httpApi.httpListener = nil

		log.Info("HTTP server stopped")
	}
}

func RegisterHttpApi(path string,
	handlerFunc func(http.ResponseWriter,
		*http.Request, map[string]interface{}) (interface{}, error)) *mux.Route {
	wrappedFunc := func(w http.ResponseWriter, req *http.Request) {
		var (
			ret interface{}
			t1  = time.Now()
		)

		params, err := httpApi.decodeHttpParams(w, req)
		if err == nil {
			ret, err = handlerFunc(w, req, params)
		} else {
			ret = map[string]interface{}{"error": err.Error()}
		}

		w.Header().Set("Content-Type", "application/json")
		var status int
		if err == nil {
			status = http.StatusOK
		} else {
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)

		// debug request body content
		//log.Trace("req body: %+v", params)

		// access log
		log.Debug("%s \"%s %s %s\" %d %s",
			req.RemoteAddr,
			req.Method,
			req.RequestURI,
			req.Proto,
			status,
			time.Since(t1))
		if status != http.StatusOK {
			log.Error("HTTP: %v", err)
		}

		if ret != nil {
			// pretty write json result
			pretty, err := json.MarshalIndent(ret, "", "    ")
			if err != nil {
				log.Error(err)
				return
			}
			w.Write(pretty)
			w.Write([]byte("\n"))
		}
	}

	if httpApi == nil {
		panic("call server.LaunchHttpServer before server.RegisterHttpApi")
	}

	// path can't be duplicated
	isDup := false
	for _, p := range httpApi.httpPaths {
		if p == path {
			log.Error("REST[%s] already registered", path)
			isDup = true
			break
		}
	}

	if !isDup {
		httpApi.httpPaths = append(httpApi.httpPaths, path)
	}

	return httpApi.httpRouter.HandleFunc(path, wrappedFunc)
}

func UnregisterAllHttpApi() {
	httpApi.httpPaths = httpApi.httpPaths[:0]
}

func (this *httpRestApi) decodeHttpParams(w http.ResponseWriter,
	req *http.Request) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return params, nil
}
