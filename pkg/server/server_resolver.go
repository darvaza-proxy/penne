package server

import (
	"context"

	"github.com/miekg/dns"

	"darvaza.org/resolver/pkg/errors"
	"darvaza.org/resolver/pkg/exdns"
	"darvaza.org/slog"

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

func (*Server) reflectEnabled(_ string) (slog.LogLevel, bool) {
	// TODO: make this configurable
	return slog.Debug, true
}

// ServeDNS handles dns requests based on the IP address of the client
func (srv *Server) ServeDNS(rw dns.ResponseWriter, req *dns.Msg) {
	var chaos bool

	exdns.ForEachQuestionOfClass(req, dns.ClassCHAOS, func(_ dns.Question) {
		chaos = true
	})

	if chaos {
		// route CHAOS requests directly to ExchangeCHAOS
		ctx := context.Background()

		rsp, err := srv.ExchangeCHAOS(ctx, req)
		if err != nil {
			rsp = errors.ErrorAsMsg(req, err)
		}

		_ = rw.WriteMsg(rsp)
		return
	}

	srv.z.ServeDNS(rw, req)
}
