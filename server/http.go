package server

import (
	"encoding/json"
	log "github.com/funkygao/log4go"
	"github.com/gorilla/mux"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	api *httpRestApi
)

type httpRestApi struct {
	httpListener net.Listener
	httpServer   *http.Server
	httpRouter   *mux.Router
	httpPaths    []string
}

func LaunchHttpServ(listenAddr string, debugAddr string) (err error) {
	if api != nil {
		return nil
	}

	api = new(httpRestApi)
	api.httpPaths = make([]string, 0, 10)
	api.httpRouter = mux.NewRouter()
	api.httpServer = &http.Server{Addr: listenAddr,
		Handler: api.httpRouter}

	api.httpListener, err = net.Listen("tcp", api.httpServer.Addr)
	if err != nil {
		api = nil
		return err
	}

	log.Info("HTTP serving at %s with pprof at %s", listenAddr, debugAddr)

	go api.httpServer.Serve(api.httpListener)
	if debugAddr != "" {
		go http.ListenAndServe(debugAddr, nil)
	}

	return nil
}

func StopHttpServ() {
	if api != nil && api.httpListener != nil {
		api.httpListener.Close()
		api.httpListener = nil

		log.Info("HTTP server stopped")
	}
}

func Launched() bool {
	return api != nil
}

func RegisterHttpApi(path string,
	handlerFunc func(http.ResponseWriter,
		*http.Request, map[string]interface{}) (interface{}, error)) *mux.Route {
	wrappedFunc := func(w http.ResponseWriter, req *http.Request) {
		var (
			ret interface{}
			t1  = time.Now()
		)

		params, err := api.decodeHttpParams(w, req)
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

	// path can't be duplicated
	isDup := false
	for _, p := range api.httpPaths {
		if p == path {
			log.Error("REST[%s] already registered", path)
			isDup = true
			break
		}
	}

	if !isDup {
		api.httpPaths = append(api.httpPaths, path)
	}

	return api.httpRouter.HandleFunc(path, wrappedFunc)
}

func UnregisterAllHttpApi() {
	api.httpPaths = api.httpPaths[:0]
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
