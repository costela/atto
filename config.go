package main

import "time"

var conf = struct {
	LogLevel string `default:"info" short:"l" desc:"level of logging output (one of debug/info/warn/error)"`
	Port     int    `default:"8080" desc:"port to listen on"`
	Path     string `default:"." desc:"path to serve"`
	Path404  string `default:"404.html" desc:"path to a file whose content will be returned on 404 (relative to --path)"`
	Prefix   string `default:"" desc:"prefix under which atto will be accessed (this will be stripped before accessing 'path')"`
	ShowList bool   `default:"false" desc:"whether to display folder contents"`
	Compress bool   `default:"true" desc:"whether to transparently compress served files"`
	Timeout  struct {
		ReadHeader *duration `default:"5s" desc:"time to wait for request headers"`
		Shutdown   *duration `default:"30s" desc:"time to wait for ungoing requests to finish before shutting down"`
	}
}{}

// work around time.Duration's lack of UnmarshalText
type duration time.Duration

func (d *duration) UnmarshalText(data []byte) error {
	dd, err := time.ParseDuration(string(data))
	*d = duration(dd)
	return err
}
