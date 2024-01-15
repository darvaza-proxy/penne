package server

// Serve runs the service
func (srv *Server) Serve() error {
	if err := srv.startResolvers(); err != nil {
		return err
	}
	defer srv.cancelResolvers()

	if err := srv.sc.Spawn(srv, srv.cfg.Supervision.HealthWait); err != nil {
		return err
	}

	return srv.sc.Wait()
}

// ListenAndServe listens all ports and runs the service.
func (srv *Server) ListenAndServe() error {
	if err := srv.startResolvers(); err != nil {
		return err
	}
	defer srv.cancelResolvers()

	return srv.sc.ListenAndServe(srv)
}
