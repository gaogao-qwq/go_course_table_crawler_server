package configs

import "flag"

var (
	Address  string
	Port     string
	LoginUrl string
	HomeUrl  string
)

func init() {
	flag.StringVar(&Address, "address", "localhost", "server listen address")
	flag.StringVar(&Port, "port", "56789", "server listen port")
	flag.StringVar(&LoginUrl, "loginurl", "http://targeturl/login", "login page url")
	flag.StringVar(&HomeUrl, "homeurl", "http://targeturl/home.action", "home page url")
	flag.Parse()
}
