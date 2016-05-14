package cannula

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/retailnext/cannula/expvar"
	"github.com/retailnext/cannula/internal/net/http/pprof"
)

var (
	defaultServer = newServer()
)

// Start listens on the unix socket specified by path and then calls Serve.
func Start(path string) error {
	// cleanup old instances of our domain socket
	os.Remove(path)

	ln, err := net.Listen("unix", path)
	if err != nil {
		return err
	}

	return Serve(ln)
}

// Serve starts serving debug http requests on listener.
func Serve(listener net.Listener) error {
	httpServer := http.Server{Handler: defaultServer.mux}
	return httpServer.Serve(listener)
}

// HandleFunc registers a new debug http.HandlerFunc.
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	Handle(pattern, http.HandlerFunc(handler))
}

// Handle registers a new debug handler.
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
	s.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	s.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	s.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	s.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	s.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	s.Handle("/debug/vars", http.HandlerFunc(expvar.Handler))

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
