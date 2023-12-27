package server

import (
	"context"

	"github.com/miekg/dns"

	"darvaza.org/resolver/pkg/errors"
	"darvaza.org/resolver/pkg/exdns"
)

// ExchangeCHAOS handles only CHAOS requests
func (srv *Server) ExchangeCHAOS(_ context.Context, req *dns.Msg) (*dns.Msg, error) {
	var answers []dns.RR

	if srv.cfg.DisableCHAOS {
		return handleRcodeSuccess(req, []dns.RR{})
	}

	exdns.ForEachQuestionOfClass(req, dns.ClassCHAOS, func(q dns.Question) {
		rr, ok := srv.chaosAnswer(q)
		if ok {
			answers = append(answers, rr...)
		}
	})

	if len(answers) == 0 {
		// no idea what to do with this request
		return handleRcodeError(req, dns.RcodeNotImplemented)
	}

	// (partial?) success
	return handleRcodeSuccess(req, answers)
}

func (srv *Server) chaosAnswer(q dns.Question) ([]dns.RR, bool) {
	switch q.Name {
	case "authors.bind.":
		if s := srv.cfg.Authors; s != "" {
			return dnsTXTAnswer(q, s)
		}
	case "version.bind.", "version.server.":
		if s := srv.cfg.Version; s != "" {
			return dnsTXTAnswer(q, s)
		}
	case "hostname.bind.", "id.server.":
		if s := srv.cfg.Name; s != "" {
			return dnsTXTAnswer(q, s)
		}
	}

	return nil, false
}

func dnsTXTAnswer(q dns.Question, content ...string) ([]dns.RR, bool) {
	rr := &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeTXT,
			Class:  q.Qclass,
		},
		Txt: content,
	}

	return []dns.RR{rr}, true
}

func handleRcodeError(req *dns.Msg, rCode int) (*dns.Msg, error) {
	resp := new(dns.Msg)
	resp.SetRcode(req, rCode)

	return nil, errors.MsgAsError(resp)
}

func handleRcodeSuccess(req *dns.Msg, answers []dns.RR) (*dns.Msg, error) {
	resp := new(dns.Msg)
	resp.SetReply(req)
	resp.Compress = false
	resp.RecursionAvailable = true
	resp.Answer = answers

	return resp, nil
}
