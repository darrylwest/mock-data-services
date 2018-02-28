//
// config  - application specification and CLI parsing
//
// @author darryl.west <darryl.west@ebay.com>
// @created 2018-02-27 16:22:25

package hub

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config the config structure
type Config struct {
	Port           int
	LogLevel       int
	MaxConcurrent  int
	DbFilename     string
	ConfigFilename string
	StaticFolder   string
	Site           string
	Pool           string
	MaxJobList     int
	Timeout        int
}

// NewDefaultConfig default settings
func NewDefaultConfig() *Config {
	cfg := new(Config)

	cfg.Port = 3000
	cfg.LogLevel = 2
	cfg.DbFilename = "data/proxy.db"
	cfg.Timeout = 20

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
	level := flag.Int("loglevel", dflt.LogLevel, "set the server's log level 0..5 for trace..error, default info=2")
	port := flag.Int("port", dflt.Port, "set the server's listening port")
	dbfilename := flag.String("db-filename", dflt.DbFilename, "set the databse file")
	timeout := flag.Int("timeout", dflt.Timeout, "the timeout for both tests and builds in munutes")

	flag.Parse()

	if *vers == true {
		fmt.Printf("Version %s\n", Version())
		return nil
	}

	log.Info("%s Version: %s\n", filepath.Base(os.Args[0]), Version())

	cfg := Config{
		Port:          *port,
		LogLevel:      *level,
		DbFilename:    *dbfilename,
		Timeout:       *timeout,
	}

	log.SetLevel(cfg.LogLevel)

	return &cfg
}
