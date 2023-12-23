package server

import "net/http"

var (
	_ http.Handler = (*Server)(nil)
)

// ServeHTTP handles HTTP requests based on the IP address of the client
func (srv *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	srv.z.ServeHTTP(rw, req)
}
