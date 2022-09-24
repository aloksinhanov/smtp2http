package main

import "flag"

var (
	flagServerName     = flag.String("name", "smtp2http", "the server name")
	flagListenAddrSSL  = flag.String("listenSSL", ":smtps", "the smtp address to listen on for SSL")
	flagListenAddr     = flag.String("listen", ":smtp", "the smtp address to listen")
	flagWebhook        = flag.String("webhook", "http://localhost:8080/my/webhook", "the webhook to send the data to")
	flagMaxMessageSize = flag.Int64("msglimit", 1024*1024*2, "maximum incoming message size")
	flagReadTimeout    = flag.Int("timeout.read", 5, "the read timeout in seconds")
	flagWriteTimeout   = flag.Int("timeout.write", 5, "the write timeout in seconds")
	flagAuthUSER       = flag.String("user", "", "user for smtp client")
	flagAuthPASS       = flag.String("pass", "", "pass for smtp client")
	flagDomain         = flag.String("domain", "", "domain for recieving mails")
	flagDBServer       = flag.String("dbserver", "", "db server")
	flagDBPort         = flag.String("dbport", "", "db server port")
	flagDBName         = flag.String("dbname", "", "db name")
	flagUserID         = flag.String("dbuser", "", "db user id")
	flagUserPwd        = flag.String("dbpwd", "", "db user paasword")
	flagPrivateKkey    = flag.String("key", "", "private key")
	flagSSLCert        = flag.String("cert", "", "ssl certificate")
)

func init() {
	flag.Parse()
}
