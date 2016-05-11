package cannula

import (
	"fmt"
	"net"
	"net/http"
	httpprof "net/http/pprof"
	"sort"
	"sync"
)

func init() {
	// reset the side effect of loading net/http/pprof
	http.DefaultServeMux = http.NewServeMux()
}

var (
	defaultServer = newServer()
)

func Start(listener net.Listener) error {
	httpServer := http.Server{Handler: defaultServer.mux}
	return httpServer.Serve(listener)
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	Handle(pattern, http.HandlerFunc(handler))
}

func Handle(pattern string, handler http.Handler) {
	defaultServer.mux.Handle(pattern, handler)
}

type server struct {
	sync.Mutex
	mux   *http.ServeMux
	paths []string
}

func newServer() *server {
	s := &server{
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("/", s.index)

	// install net/http/pprof handlers also
	s.Handle("/debug/pprof/", http.HandlerFunc(httpprof.Index))
	s.Handle("/debug/pprof/cmdline", http.HandlerFunc(httpprof.Cmdline))
	s.Handle("/debug/pprof/profile", http.HandlerFunc(httpprof.Profile))
	s.Handle("/debug/pprof/symbol", http.HandlerFunc(httpprof.Symbol))
	s.Handle("/debug/pprof/trace", http.HandlerFunc(httpprof.Trace))

	return s
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "Available handlers:")

	s.Lock()
	defer s.Unlock()
	sort.Strings(s.paths)
	for _, path := range s.paths {
		fmt.Fprintln(w, path)
	}
}

func (s *server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)

	s.Lock()
	defer s.Unlock()
	s.paths = append(s.paths, pattern)
}
