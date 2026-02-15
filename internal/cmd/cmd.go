// Package cmd is responsible for the program's command-line interface.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AdguardTeam/golibs/log"
	goFlags "github.com/jessevdk/go-flags"
	"github.com/zamibd/gorao/internal/dnsproxy"
	gorao "github.com/zamibd/gorao/internal/sniproxy"
	"github.com/zamibd/gorao/internal/version"
	"gopkg.in/yaml.v3"
)

// Main is the entry point of the program.
func Main() {
	for _, arg := range os.Args {
		if arg == "--version" {
			fmt.Printf("gorao version: %s\n", version.VersionString)
			os.Exit(0)
		}
	}

	options := DefaultOptions()

	// Parse config file if exists.
	configFile := "config.yaml"
	for i, arg := range os.Args {
		if arg == "--config" && i+1 < len(os.Args) {
			configFile = os.Args[i+1]
		} else if strings.HasPrefix(arg, "--config=") {
			configFile = strings.TrimPrefix(arg, "--config=")
		}
	}

	if content, err := os.ReadFile(configFile); err == nil {
		content = []byte(os.ExpandEnv(string(content)))
		if err = yaml.Unmarshal(content, options); err != nil {
			log.Fatalf("cmd: cannot parse config file: %s", err)
		}
		log.Info("cmd: loaded configuration from %s", configFile)
	}

	parser := goFlags.NewParser(options, goFlags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*goFlags.Error); ok && flagsErr.Type == goFlags.ErrHelp {
			os.Exit(0)
		}

		os.Exit(1)
	}

	if options.Verbose {
		log.SetLevel(log.DEBUG)
	}
	if options.LogOutput != "" {
		var file *os.File
		file, err = os.OpenFile(options.LogOutput, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil {
			log.Fatalf("cannot create a log file: %s", err)
		}
		defer log.OnCloserError(file, log.INFO)
		log.SetOutput(file)
	}

	// Load rules from CSV files if specified
	if options.ForwardRulesFile != "" {
		fileRules, err := loadRulesFromFile(options.ForwardRulesFile)
		if err != nil {
			log.Fatalf("cmd: failed to load forward rules from %s: %v", options.ForwardRulesFile, err)
		}
		options.ForwardRules = append(options.ForwardRules, fileRules...)
	}

	if options.DNSRedirectRulesFile != "" {
		fileRules, err := loadRulesFromFile(options.DNSRedirectRulesFile)
		if err != nil {
			log.Fatalf("cmd: failed to load DNS redirect rules from %s: %v", options.DNSRedirectRulesFile, err)
		}
		options.DNSRedirectRules = append(options.DNSRedirectRules, fileRules...)
	}

	if options.BlockRulesFile != "" {
		fileRules, err := loadRulesFromFile(options.BlockRulesFile)
		if err != nil {
			log.Fatalf("cmd: failed to load block rules from %s: %v", options.BlockRulesFile, err)
		}
		options.BlockRules = append(options.BlockRules, fileRules...)
	}

	if options.DropRulesFile != "" {
		fileRules, err := loadRulesFromFile(options.DropRulesFile)
		if err != nil {
			log.Fatalf("cmd: failed to load drop rules from %s: %v", options.DropRulesFile, err)
		}
		options.DropRules = append(options.DropRules, fileRules...)
	}

	run(options)
}

// Proxy is an interface that combines Start and Close methods.
type Proxy interface {
	Start() error
	io.Closer
}

// run starts reads the configuration options and starts the gorao.
func run(options *Options) {
	log.Info("cmd: run gorao with the following configuration:\n%s", options)

	dnsProxy := newDNSProxy(options)
	err := dnsProxy.Start()
	check(err)

	gorao := newgorao(options)
	err = gorao.Start()
	check(err)

	// Subscribe to the OS events.
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	log.Info("cmd: stopping gorao")
	log.OnCloserError(dnsProxy, log.INFO)
	log.OnCloserError(gorao, log.INFO)
}

// newDNSProxy creates a new instance of [*dnsproxy.DNSProxy] or panics if any
// error happens.
func newDNSProxy(options *Options) (d *dnsproxy.DNSProxy) {
	cfg := toDNSProxyConfig(options)

	d, err := dnsproxy.New(cfg)
	check(err)

	return d
}

// newgorao creates a new instance of the gorao proxy or panics if any
// error happens.
func newgorao(options *Options) (p Proxy) {
	cfg := togoraoConfig(options)

	p, err := gorao.New(cfg)
	check(err)

	return p
}

// check log.Fatalf if err is not nil.
func check(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
