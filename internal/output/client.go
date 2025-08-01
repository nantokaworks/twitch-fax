package output

import (
	"git.massivebox.net/massivebox/go-catprinter"
	"github.com/nantokaworks/twitch-fax/internal/env"
	"github.com/nantokaworks/twitch-fax/internal/shared/logger"
	"github.com/nantokaworks/twitch-fax/internal/status"
	"go.uber.org/zap"
)

var latestPrinter *catprinter.Client
var opts *catprinter.PrinterOptions
var isConnected bool
var hasInitialPrintBeenDone bool

func SetupPrinter() (*catprinter.Client, error) {
	if latestPrinter != nil {
		latestPrinter.Disconnect()
		latestPrinter = nil
		isConnected = false
		status.SetPrinterConnected(false)
	}

	instance, err := catprinter.NewClient()
	if err != nil {
		return nil, err
	}
	latestPrinter = instance
	return instance, nil
}

func ConnectPrinter(c *catprinter.Client, address string) error {
	if c == nil {
		return nil
	}
	
	// Skip if already connected
	if isConnected {
		return nil
	}

	// DRY-RUNモードでも実際のプリンターに接続
	if env.Value.DryRunMode {
		logger.Info("Connecting to printer in DRY-RUN mode", zap.String("address", address))
	}

	logger.Info("Connecting to printer", zap.String("address", address))
	err := c.Connect(address)
	if err != nil {
		return err
	}
	logger.Info("Successfully connected to printer", zap.String("address", address))
	isConnected = true
	status.SetPrinterConnected(true)

	return nil
}

func SetupPrinterOptions(bestQuality, dither, autoRotate bool, blackPoint float32) error {
	// Set up the printer options
	opts = catprinter.NewOptions().
		SetBestQuality(bestQuality).
		SetDither(dither).
		SetAutoRotate(autoRotate).
		SetBlackPoint(float32(blackPoint))

	return nil
}

// Stop gracefully disconnects the printer if connected
func Stop() {
	if latestPrinter != nil && isConnected {
		latestPrinter.Disconnect()
		isConnected = false
		status.SetPrinterConnected(false)
		latestPrinter = nil
	}
}

// MarkInitialPrintDone marks that the initial print has been completed
func MarkInitialPrintDone() {
	hasInitialPrintBeenDone = true
}

// IsConnected returns whether the printer is connected
func IsConnected() bool {
	return isConnected
}

// HasInitialPrintBeenDone returns whether the initial print has been done
func HasInitialPrintBeenDone() bool {
	return hasInitialPrintBeenDone
}
