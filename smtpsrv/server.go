package smtpsrv

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/emersion/go-smtp"
)

type ServerConfig struct {
	ListenAddrSSL   []string
	ListenAddr      []string
	BannerDomain    string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	Handler         HandlerFunc
	Auther          AuthFunc
	MaxMessageBytes int
	TLSConfig       *tls.Config
}

func ListenAndServe(cfg *ServerConfig, addr string) error {
	s := smtp.NewServer(NewBackend(cfg.Auther, cfg.Handler))

	SetDefaultServerConfig(cfg)

	s.Addr = addr
	s.Domain = cfg.BannerDomain
	s.ReadTimeout = cfg.ReadTimeout
	s.WriteTimeout = cfg.WriteTimeout
	s.MaxMessageBytes = cfg.MaxMessageBytes
	s.AllowInsecureAuth = true
	s.AuthDisabled = true
	if cfg.Auther != nil {
		s.AuthDisabled = false
	}
	s.EnableSMTPUTF8 = false
	s.TLSConfig = cfg.TLSConfig

	fmt.Println("⇨ smtp server started on", s.Addr)

	return s.ListenAndServe()
}

func ListenAndServeTLS(cfg *ServerConfig, addr string) error {
	s := smtp.NewServer(NewBackend(cfg.Auther, cfg.Handler))

	SetDefaultServerConfig(cfg)

	s.Addr = addr
	s.Domain = cfg.BannerDomain
	s.ReadTimeout = cfg.ReadTimeout
	s.WriteTimeout = cfg.WriteTimeout
	s.MaxMessageBytes = cfg.MaxMessageBytes
	s.AllowInsecureAuth = true
	s.AuthDisabled = true
	if cfg.Auther != nil {
		s.AuthDisabled = false
	}
	s.EnableSMTPUTF8 = false
	s.EnableREQUIRETLS = true
	s.TLSConfig = cfg.TLSConfig

	fmt.Println("⇨ smtp server for SSL started on", s.Addr)

	return s.ListenAndServeTLS()
}
