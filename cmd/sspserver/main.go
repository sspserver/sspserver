package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/demdxx/cloudregistry"
	"github.com/demdxx/gocast/v2"
	"github.com/demdxx/goconfig"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/context/version"
	"github.com/geniusrabbit/adcorelib/profiler"
	"github.com/geniusrabbit/adcorelib/zlogger"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/sspserver/sspserver/cmd/sspserver/appcontext"
	"github.com/sspserver/sspserver/cmd/sspserver/commands"
	"github.com/sspserver/sspserver/cmd/sspserver/jobs"
	"github.com/sspserver/sspserver/internal/cregistry"
)

var (
	config       appcontext.Config
	buildCommit  = ""
	buildVersion = "develop"
	buildDate    = "unknown"
)

// Define command list
var cmdList = commands.ICommands{
	commands.SSPServerCommand,
}

func init() {
	fmt.Println()
	fmt.Println("███████ ███████ ██████  ███████ ███████ ██████  ██    ██ ███████ ██████")
	fmt.Println("██      ██      ██   ██ ██      ██      ██   ██ ██    ██ ██      ██   ██")
	fmt.Println("███████ ███████ ██████  ███████ █████   ██████  ██    ██ █████   ██████")
	fmt.Println("     ██      ██ ██           ██ ██      ██   ██  ██  ██  ██      ██   ██")
	fmt.Println("███████ ███████ ██      ███████ ███████ ██   ██   ████   ███████ ██   ██")
	fmt.Println()
	fmt.Println("Version:", buildVersion, " (", buildCommit, ")")
	fmt.Println("Build date:", buildDate)
	fmt.Println()

	args := os.Args
	if len(args) > 1 {
		args = args[2:]
	}

	fatalError(goconfig.Load(
		&config,
		goconfig.WithDefaults(),
		goconfig.WithCustomArgs(args...),
		goconfig.WithEnv(),
	), "config loading")

	// Init new logger object
	loggerObj, err := zlogger.New(config.ServiceName, config.LogEncoder,
		config.LogLevel, config.LogAddr, zap.Fields(
			zap.String("commit", buildCommit),
			zap.String("version", buildVersion),
			zap.String("build_date", buildDate),
		))
	fatalError(err, "configure logger")

	// Replace global logger
	zap.ReplaceGlobals(loggerObj)

	// Print configuration
	if config.IsDebug() {
		fmt.Println(config.String())
	}
}

func main() {
	var (
		logger            = zap.L()
		numberOfAdServers = cloudregistry.NewSyncUInt64Value(max(1, uint64(config.Server.Datacenter.ServiceCount)))
		ctx, cancel       = signal.NotifyContext(context.Background(), os.Interrupt)
	)
	defer cancel()

	// Add logger to context
	ctx = ctxlogger.WithLogger(ctx, logger)

	// Register version information
	ctx = version.WithContext(ctx, &version.Version{
		Commit:  buildCommit,
		Version: buildVersion,
		Date:    buildDate,
	})

	if len(os.Args) < 2 {
		printCommandsUsage()
		return
	}

	// Get command name
	cmdName := os.Args[1]

	// Run command by name
	icmd := cmdList.Get(cmdName)

	// Print help if command not found
	if cmdName == "help" || icmd == nil {
		printCommandsUsage()
		return
	}

	// =========== Init cloud registry ========================
	// Cloud registry is a main entry point for service discovery
	// ========================================================

	if config.Server.Registry.Connection != "" {
		fatalError(
			initCloudRegistry(ctx, &config, numberOfAdServers),
			"cloud registry init",
		)
	}

	// Profiling server of collector
	profiler.Run(config.Server.Profile.Mode,
		config.Server.Profile.Listen, logger, true)

	// Run command with context
	fmt.Println()
	fmt.Println("░█ Run command:\x1b[31m", icmd.Cmd(), "\x1b[0m")
	fmt.Println()

	fatalError(icmd.Run(ctx, os.Args[2:], numberOfAdServers), "command execution")
}

func printCommandsUsage() {
	fmt.Println("Usage: sspserver <command> [options]")
	fmt.Println("Commands:")
	for _, cmd := range cmdList {
		fmt.Printf("  % 10s - %s\n", cmd.Cmd(), cmd.Help())
	}
	fmt.Println("  help - print this help")
}

func initCloudRegistry(ctx context.Context, config *appcontext.Config, numberOfAdServers *cloudregistry.SyncUInt64Value) error {
	// Connect to cloud registry and discover services
	registry, err := cregistry.Connect(ctx, config.Server.Registry.Connection)
	if err != nil {
		return errors.Wrap(err, "connect to cloud registry")
	}

	// Get hostname from listen address
	if config.Server.Registry.Hostname == "" {
		if config.Server.Hostname != "" {
			config.Server.Registry.Hostname = config.Server.Hostname
		} else if !strings.HasPrefix(config.Server.HTTP.Listen, ":") {
			config.Server.Registry.Hostname = config.Server.HTTP.Listen[:strings.IndexByte(config.Server.HTTP.Listen, ':')]
		}
	}

	// Get port from listen address
	if config.Server.Registry.Port == 0 {
		config.Server.Registry.Port = gocast.Int(config.Server.HTTP.Listen[strings.LastIndexByte(config.Server.HTTP.Listen, ':')+1:])
	}

	// Register service in cloud registry
	err = registry.Register(ctx, &cloudregistry.Service{
		Name:       config.ServiceName,
		InstanceID: cloudregistry.GenerateInstanceID(config.ServiceName),
		Hostname:   config.Server.Registry.Hostname,
		Port:       config.Server.Registry.Port,
		Check: cloudregistry.Check{
			ID:  "health",
			TTL: 20 * time.Second,
			HTTP: struct {
				URL     string
				Method  string
				Headers map[string][]string
			}{
				URL:    "/health",
				Method: "GET",
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "register service in cloud registry")
	}

	// Run service discovery
	go jobs.RunIntervalJob(ctx, "service-discovery", 30*time.Second, func(ctx context.Context) error {
		services, err := registry.Discover(ctx,
			&cloudregistry.ServicePrefix{Name: config.ServiceName}, 30*time.Second)
		ctxlogger.Get(ctx).Info("service discovery", zap.Int("count", len(services)), zap.Error(err))
		if err != nil {
			_ = numberOfAdServers.SetValue("service", len(services))
		}
		return nil
	})

	return nil
}

func fatalError(err error, message ...any) {
	if err != nil {
		zap.L().Fatal(fmt.Sprint(message...), zap.Error(err))
	}
}
