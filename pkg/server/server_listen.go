package server

func (srv *Server) initSidecar() error {
	scc, err := srv.cfg.export(srv.tls)
	if err != nil {
		return err
	}

	sc, err := scc.New()
	if err != nil {
		return err
	}

	srv.sc = sc
	return nil
}

// Listen listens all ports.
func (srv *Server) Listen() error {
	return srv.sc.Listen()
}
