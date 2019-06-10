/**
	HLC FOUNDATION
	james
 */

package tool

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"log"
	"Nox-DAG-test/script/tool/go-flags"
)

const (
	defaultConfigFilename = "halalchainminer.conf"
	defaultLogLevel       = "info"
	defaultLogDirname     = "logs"
)

var (
	minerHomeDir          = GetCurrentDir()
	noxHomeDir           = AppDataDir("nox", false)
	defaultConfigFile     = filepath.Join(minerHomeDir, defaultConfigFilename)
	defaultRPCServer      = "127.0.0.1"
	defaultRPCCertFile    = filepath.Join(noxHomeDir, "rpc.cert")
	defaultLogDir         = filepath.Join(minerHomeDir, defaultLogDirname)
)

type Config struct {
	// Config / log options
	Experimental bool   `long:"experimental" description:"enable EXPERIMENTAL features such as setting a temperature target with (-t/--temptarget) which may DAMAGE YOUR DEVICE(S)."`
	ConfigFile   string `short:"C" long:"configfile" description:"Path to configuration file"`
	Temp   string `short:"t" long:"temp" description:"temp"`
	Network   string `short:"o" long:"network" description:"test"`
	LogDir       string `long:"logdir" description:"Directory to log output."`
	DebugLevel   string `short:"d" long:"debuglevel" description:"Logging level for all subsystems {trace, debug, info, warn, error, critical} -- You may also specify <subsystem>=<level>,<subsystem2>=<level>,... to set the log level for individual subsystems -- Use show to list available subsystems"`

	// RPC connection options
	RPCUser     string `short:"u" long:"rpcuser" description:"RPC username"`
	//Dag     	bool `short:"dag" long:"dag" description:"dag mining"`
	RPCPassword string `short:"P" long:"rpcpass" default-mask:"-" description:"RPC password"`
	RPCServer   string `short:"s" long:"rpcserver" description:"RPC server to connect to"`
	RPCCert     string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	NoTLS       bool   `long:"notls" description:"Disable TLS"`
	Send       bool   `long:"send" description:"send tx"`
	Height       int   `long:"height" description:"send height"`
	Proxy       string `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser   string `long:"proxyuser" description:"Username for proxy server"`
	ProxyPass   string `long:"proxypass" default-mask:"-" description:"Password for proxy server"`
	FromPrivateKey   string `short:"f" long:"privkey" default-mask:"-" description:"private key"`
	FromAddress   string `short:"x" long:"faddress" default-mask:"-" description:"from address"`
	AddressFile   string `short:"q" long:"addressfile" default-mask:"-" description:"from address"`
	TXFile   string `short:"r" long:"txfile" default-mask:"-" description:"from address"`
	FromTransactionHash   string `short:"y" long:"fromtransaction" default-mask:"-" description:"from transaction"`

	Benchmark bool `short:"B" long:"benchmark" description:"Run in benchmark mode."`
	Action string `long:"action" description:"Run action."`

	TestNet       bool `long:"testnet" description:"Connect to testnet"`
}

// removeDuplicateAddresses returns a new slice with all duplicate entries in
// addrs removed.
func removeDuplicateAddresses(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	seen := map[string]struct{}{}
	for _, val := range addrs {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = struct{}{}
		}
	}
	return result
}

// normalizeAddress returns addr with the passed default port appended if
// there is not already a port specified.
func normalizeAddress(addr string, defaultPort string) string {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return net.JoinHostPort(addr, defaultPort)
	}
	return addr
}

// normalizeAddresses returns a new slice with all the passed peer addresses
// normalized with the given default port, and all duplicates removed.
func normalizeAddresses(addrs []string, defaultPort string) []string {
	for i, addr := range addrs {
		addrs[i] = normalizeAddress(addr, defaultPort)
	}

	return removeDuplicateAddresses(addrs)
}

// cleanAndExpandPath expands environement variables and leading ~ in the
// passed path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// Expand initial ~ to OS specific home directory.
	if strings.HasPrefix(path, "~") {
		homeDir := filepath.Dir(minerHomeDir)
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%,
	// but they variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}

// loadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
//
// The above results in btcd functioning properly without any config settings
// while still allowing the user to override settings with config files and
// command line options.  Command line options always take precedence.
func LoadConfig() (*Config, []string, error) {
	// Default config.
	cfg := Config{
		ConfigFile: defaultConfigFile,
		DebugLevel: defaultLogLevel,
		LogDir:     defaultLogDir,
		RPCServer:  defaultRPCServer,
		RPCCert:    defaultRPCCertFile,
	}

	// Create the home directory if it doesn't already exist.
	err := os.MkdirAll(minerHomeDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}

	// Pre-parse the command line options to see if an alternative config
	// file or the version flag was specified.
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.Default)
	_, err = preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			preParser.WriteHelp(os.Stderr)
		}
		return nil, nil, err
	}

	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))

	// Load additional config from file.
	var configFileError error
	parser := flags.NewParser(&cfg, flags.Default)
	err = flags.NewIniParser(parser).ParseFile(preCfg.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			fmt.Fprintln(os.Stderr, err)
			parser.WriteHelp(os.Stderr)
			return nil, nil, err
		}
		configFileError = err
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, nil, err
	}


	if cfg.Experimental {
		fmt.Fprintln(os.Stderr, "enabling EXPERIMENTAL features "+
			"that may possibly DAMAGE YOUR DEVICE(S)")
		time.Sleep(time.Second * 3)
	}

	// Handle environment variable expansion in the RPC certificate path.
	cfg.RPCCert = cleanAndExpandPath(cfg.RPCCert)

	var defaultRPCPort string

	// Add default port to RPC server based on --testnet flag
	// if needed.
	cfg.RPCServer = normalizeAddress(cfg.RPCServer, defaultRPCPort)

	// Warn about missing config file only after all other configuration is
	// done.  This prevents the warning on help messages and invalid
	// options.  Note this should go directly before the return.
	if configFileError != nil {
		log.Printf("%v", configFileError)
	}

	return &cfg, remainingArgs, nil
}
