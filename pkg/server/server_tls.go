package server

func (srv *Server) initTLS() error {
	s, err := srv.cfg.TLS.New(srv.cfg.Logger)
	if err != nil {
		return err
	}

	srv.tls = s
	return nil
}
