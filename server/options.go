package server

import "crypto/tls"

type Options func(s *Server)

func WithTLS(tls *tls.Config) Options {
	return func(s *Server) {
		s.tlsConf = tls
	}
}
