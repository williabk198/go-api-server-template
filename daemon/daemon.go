package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/williabk198/go-api-server-template/controller"
	"github.com/williabk198/go-api-server-template/db/dummydb"
	"github.com/williabk198/go-api-server-template/router"
)

// Start is a blocking function that will initialize and startup the API server.
// This function only returns if and error occured when starting the server or if a
// terminate or interrupt signal was recieved from the OS
func Start() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	database := dummydb.NewSession() // Update
	controls := controller.NewController(logger, database)
	routes := router.NewRouter(controls)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: routes,
	}

	// Create an error channel so that fatal server errors can be logged appropriately
	errChan := make(chan error)
	defer close(errChan)
	go func() { // Have the server start on a new thread, and send any errors it may encounter to the error channel
		logger.Info("starting server", "address", server.Addr)
		errChan <- server.ListenAndServe()
	}()

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)

	// Indicate that we want to listen for the interupt signal (SIGINT) and terminate signal (SIGTERM) from the OS
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for something to come in on one of the channels
	select {
	case sig := <-sigChan: // If we recieve a signal from the OS, then shutdown gracefully
		logger.Info("Recieved signal from the OS and is shutting down.", "signal", sig)
	case err := <-errChan: // If an error occured in server.ListenAndServe(), the just return
		logger.Error("The server encountered an error", "error", err)
		return
	}

	// Create a context with a deadline so that server.Shutdown doesn't
	// wait indefinitely for connections to complete their requests
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	// Try to let active requests to finish up and then close.
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Warn("encountered error during server shutdown", "error", err)
	}
}
