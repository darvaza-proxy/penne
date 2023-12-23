package server

import (
	"github.com/miekg/dns"

	"darvaza.org/penne/pkg/resolver"
)

var (
	_ dns.Handler = (*Server)(nil)
)

func defaultResolvers() []resolver.Config {
	return []resolver.Config{
		{
			Name:      "root",
			Iterative: true,
		},
	}
}

func (srv *Server) initResolvers() error {
	_, _, err := resolver.MakeResolvers(srv.cfg.Resolvers, srv.cfg.Logger)
	return err
}

// ServeDNS handles dns requests based on the IP address of the client
func (srv *Server) ServeDNS(rw dns.ResponseWriter, req *dns.Msg) {
	srv.z.ServeDNS(rw, req)
}
