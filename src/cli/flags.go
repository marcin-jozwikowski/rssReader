package cli

import ("flag")

var ConfigFileName *string
var CacheFileName *string
var RunEditor *bool
var Verbose *int
var DownloadCommand *string
var DownloadParams *string
var ProxyAddr *string
var ProxyType *string
var ProxyAuth *string

func init()  {
	ConfigFileName = flag.String("configFile", "config.json", "Config file location")
	CacheFileName = flag.String("cacheFile", "cache.json", "Cache file location")
	RunEditor = flag.Bool("editConfig", false, "Run configuration editor")
	Verbose = flag.Int("verbose", DefaultVerbose, "Verbose level: 0-None ... 3-All")
	DownloadCommand = flag.String("downloadCommand", "", "Download command to run. Must return content to stdout.")
	DownloadParams = flag.String("downloadParams", "", "Download command params... See online documentation for details.")
	ProxyAddr = flag.String("proxy", "", "Proxy address to use. Valid only if proxy specified")
	ProxyType = flag.String("proxyType", "socks5", "Proxy type to use. Valid only if proxy specified")
	ProxyAuth = flag.String("proxyAuth", "", "Proxy auth to use. Valid only if proxy specified")
	flag.Parse()

	SetVerbose(*Verbose)
}
