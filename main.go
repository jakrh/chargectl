//go:build darwin

package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"chargectl/pkg/batt"
)

const (
	ExitWithSuccess = 0
	ExitWithFailure = 1
)

// TODO: Consider to wrap the flags in a struct for better organization
var (
	enableCharging  bool
	disableCharging bool
	maintain        bool
	status          bool
	upperLimit      int
	lowerLimit      int
	logLevel        string
)

func setupFlags() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w, `Note:
The -m option is a long-running process that will keep checking the
battery level every %g seconds.

It disables charging when the battery level is at or above the upper limit
and enables charging when the battery level is at or below the lower limit.

Please ensure that your system is not going to sleep while maintaining the
battery charging. Because if the system goes to sleep, the process will be
suspended and will not be able to check the battery level and disable or
enable charging as needed.

To prevent the system from sleeping in charging, you can turn on the option
"Prevent automatic sleeping on power adapter when the display is off."
in "System Settings" > "Battery" > "Options" section on macOS Sequoia
(version 15).

For more information, see:
https://support.apple.com/guide/mac-help/mchle41a6ccd

To stop battery charging maintenance, you can use Ctrl+C or send a SIGINT
signal to the process.
`, batt.MaintainInterval.Seconds())
	}

	flag.BoolVar(&enableCharging, "ec", false, "Enable charging")
	flag.BoolVar(&disableCharging, "dc", false, "Disable charging")
	flag.BoolVar(&maintain, "m", false, "Maintain battery charging within upper and lower limits")
	flag.BoolVar(&status, "s", false, "Show status")
	flag.IntVar(&upperLimit, "upper", 75, "Set upper limit for battery level for maintaining")
	flag.IntVar(&lowerLimit, "lower", 70, "Set lower limit for battery level for maintaining")
	flag.StringVar(&logLevel, "l", "info", "Set log level (debug, info, warn, error)")
	flag.Parse()
}

func setupLogger() {
	// Set up the logger with the specified log level
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		slog.Error("invalid log level", "level", logLevel)
		return
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}

func statusHandler(w io.Writer, smcConn *batt.AppleSMC) {
	isChargingEnabled, err := smcConn.IsChargingEnabled()
	if err != nil {
		slog.Error("failed to check if charging is enabled", "error", err)
		return
	}

	isAdapterEnabled, err := smcConn.IsAdapterEnabled()
	if err != nil {
		slog.Error("failed to check if adapter is enabled", "error", err)
		return
	}

	fmt.Fprintf(w, `Status:
  Charging Enabled: %v
  Adapter Enabled: %v
`, isChargingEnabled, isAdapterEnabled)
}

func run() int {
	setupFlags()
	setupLogger()

	// Open Apple SMC for read/writing
	smcConn := batt.New()
	if err := smcConn.Open(); err != nil {
		slog.Error("failed to open smc connection", "error", err)
		return ExitWithFailure
	}
	defer func() {
		slog.Debug("smc connection closing")
		err := smcConn.Close()
		if err != nil {
			slog.Error("failed to close smc connection", "error", err)
		}
		slog.Debug("smc connection closed")
	}()

	switch {
	case enableCharging:
		// Check if running as root
		if os.Getuid() != 0 {
			slog.Error("Not running as root, exit without changing settings")
			return ExitWithFailure
		}

		if err := smcConn.EnableCharging(); err != nil {
			slog.Error("failed to enable charging", "error", err)
			return ExitWithFailure
		}
		slog.Info("Charging enabled")
	case disableCharging:
		// Check if running as root
		if os.Getuid() != 0 {
			slog.Error("Not running as root, exit without changing settings")
			return ExitWithFailure
		}

		if err := smcConn.DisableCharging(); err != nil {
			slog.Error("failed to disable charging", "error", err)
			return ExitWithFailure
		}
		slog.Info("Charging disabled")
	case maintain:
		// Check if running as root
		if os.Getuid() != 0 {
			slog.Error("Not running as root, exit without maintaining battery")
			return ExitWithFailure
		}
		if err := batt.Maintain(upperLimit, lowerLimit); err != nil {
			slog.Error("failed to maintain battery charging", "error", err)
			return ExitWithFailure
		}
		slog.Info("Battery maintenance completed")
	case status:
		// Show status
		statusHandler(os.Stdout, smcConn)
		return ExitWithSuccess
	default:
		// If no flags are set, show status by default
		statusHandler(os.Stdout, smcConn)
		return ExitWithSuccess
	}

	return ExitWithSuccess
}

func main() {
	// Return the exit code and keep defer works
	os.Exit(run())
}
