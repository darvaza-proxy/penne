package server

import "darvaza.org/core"

// Listen listens ports
func (*Server) Listen() error {
	return core.ErrNotImplemented
}

// ListenAndServe listens ports and runs the service
func (srv *Server) ListenAndServe() error {
	err := srv.Listen()
	if err != nil {
		return err
	}

	return srv.Serve()
}
