//
// config  - application specification and CLI parsing
//
// @author darryl.west <darryl.west@ebay.com>
// @created 2017-11-27 17:56:46

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
	cfg.MaxConcurrent = 8
	cfg.DbFilename = "data/hub.db"
	cfg.StaticFolder = "public-html"
	cfg.Site = "US"
	cfg.Pool = "www.sellexsvc-stg.qa.ebay.com" // "sellexsvc.stratus.qa.ebay.com"
	cfg.MaxJobList = 8
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
	concurrent := flag.Int("max-concurrent", dflt.MaxConcurrent, "set the maximum number of concurrent tests")
	dbfilename := flag.String("db-filename", dflt.DbFilename, "set the databse file")
	staticFolder := flag.String("static", dflt.StaticFolder, "set the dashboard's static html folder")
	site := flag.String("site", dflt.Site, "set the default site")
	pool := flag.String("pool", dflt.Pool, "set the default pool")
	maxJobList := flag.Int("max-job-list", dflt.MaxJobList, "the maximum number of jobs to display in detailed view")
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
		MaxConcurrent: *concurrent,
		DbFilename:    *dbfilename,
		StaticFolder:  *staticFolder,
		Site:          *site,
		Pool:          *pool,
		MaxJobList:    *maxJobList,
		Timeout:       *timeout,
	}

	log.SetLevel(cfg.LogLevel)

	return &cfg
}
