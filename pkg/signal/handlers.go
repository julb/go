package signal

import (
	"os"
	osSignal "os/signal"

	log "github.com/julb/go/pkg/logging"
)

type SignalTrappedContext struct {
	SignalHandled <-chan struct{}
}

func TrapSignal(callback func(os.Signal), signals ...os.Signal) *SignalTrappedContext {
	// Create a channel to get notification that the server has shutdown.
	signalHasBeenHandled := make(chan struct{})

	go func() {
		// Create a specific channel to receive interruption notification.
		signalReceivedChannel := make(chan os.Signal, 1)
		osSignal.Notify(signalReceivedChannel, signals...)

		// Channel triggered when signal is received.
		interruptionSignal := <-signalReceivedChannel

		// Trace signal receipt
		log.Infof("received signal '%v' - invoking callback method if any.", interruptionSignal)

		// Execute callback if provided
		if callback != nil {
			callback(interruptionSignal)
		}

		// Close chan
		close(signalHasBeenHandled)
	}()

	return &SignalTrappedContext{
		SignalHandled: signalHasBeenHandled,
	}
}
