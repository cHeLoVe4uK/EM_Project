package app

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Deps struct {
	Router *httprouter.Router
}

type Server struct {
	// TODO logger, configuration
	r *httprouter.Router
}

// NewServer
// return new application muxer with panic recover and logger middleware
func NewServer(d *Deps) *Server {
	return &Server{
		r: d.Router,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rcv := recover(); rcv != nil {
			s.PanicHandler(w, r, rcv)
		}
	}()
	// body limit
	r.Body = http.MaxBytesReader(w, r.Body, 10<<10)
	s.LogMiddleware(s.r).ServeHTTP(w, r)
}

// PanicHandler
// write 500 status
func (s *Server) PanicHandler(w http.ResponseWriter, r *http.Request, rcv any) {
	w.WriteHeader(http.StatusInternalServerError)
	//s.l.Error("recover from panic: ", rcv)
}

// LogMiddleware
// logging all requests
func (s *Server) LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sTime := time.Now().UTC()
		handler.ServeHTTP(w, r)
		_ = time.Since(sTime)
		// TODO logging
		//s.l.Debugf("[%s]\t%s\t%s\ttime:%s", r.Method, r.URL.Path, r.RemoteAddr, eTime.String())
	})
}
