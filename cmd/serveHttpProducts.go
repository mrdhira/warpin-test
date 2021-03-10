package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	Routes "github.com/mrdhira/warpin-test/api/Products/deliveries/http"
	"github.com/spf13/cobra"
)

// serveHTTPProductsCmd add command
var serveHTTPProductsCmd = &cobra.Command{
	Use:   "serveHttpProducts",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Route := new(Routes.Route)

		Router := Route.Init()

		var GracefulStop = make(chan os.Signal)
		signal.Notify(GracefulStop, syscall.SIGTERM)
		signal.Notify(GracefulStop, syscall.SIGINT)

		HTTPServer := &http.Server{
			Handler:      Router,
			Addr:         "0.0.0.0:8002",
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
		}

		go func() {
			if err := HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("Error on listen and serve: %v", err)
			}
		}()

		<-GracefulStop
		if err := HTTPServer.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
		fmt.Println("Order Services Closed")
	},
}

func init() {
	rootCmd.AddCommand(serveHTTPProductsCmd)
}
