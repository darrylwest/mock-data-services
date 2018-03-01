//
// config  - application specification and CLI parsing
//
// @author darryl.west <darryl.west@ebay.com>
// @created 2018-02-27 16:22:25

package proxy

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config the config structure
type Config struct {
	Port       int
	Target     string
	Bypass     bool
	LogLevel   int
	DbFilename string
	Timeout    int
	BufSize    int
}

// NewDefaultConfig default settings
func NewDefaultConfig() *Config {
	cfg := new(Config)

	cfg.Port = 3300
	cfg.Target = "127.0.0.1:9090"
	cfg.LogLevel = 2
	cfg.DbFilename = "data/proxy.db"
	cfg.Timeout = 120 // seconds
	cfg.BufSize = 64  // 1K

	return cfg
}

// ShowHelp dump out the use/command line options
func ShowHelp() {
	fmt.Printf("\n%s USE:\n\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Printf("\n%s Version %s\n", os.Args[0], Version())
}

// ParseArgs parse the command line args
func ParseArgs() *Config {
	dflt := NewDefaultConfig()

	vers := flag.Bool("version", false, "show the version and exit")
	level := flag.Int("loglevel", dflt.LogLevel, "set the server's log level 0..5, default info=2")
	port := flag.Int("port", dflt.Port, "set the server's listening port")
	target := flag.String("target", dflt.Target, "the address and port of the target machine")
	bypass := flag.Bool("bypass", dflt.Bypass, "bypass connection to target and return only mock data")
	dbfilename := flag.String("db-filename", dflt.DbFilename, "set the databse file")
	timeout := flag.Int("timeout", dflt.Timeout, "the timeout for both tests and builds in seconds")
	bufsize := flag.Int("bufsize", dflt.BufSize, "the buffer size in 1K increments")

	flag.Parse()

	if *vers == true {
		return nil
	}

	fmt.Println(logo)
	fmt.Printf("Version %s\n", Version())

	log.Info("%s Version: %s\n", filepath.Base(os.Args[0]), Version())

	cfg := Config{
		Port:       *port,
		Target:     *target,
		Bypass:     *bypass,
		LogLevel:   *level,
		DbFilename: *dbfilename,
		Timeout:    *timeout,
		BufSize:    *bufsize,
	}

	log.SetLevel(cfg.LogLevel)

	return &cfg
}
