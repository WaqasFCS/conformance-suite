package main

import (
	"fmt"
	"os"

	os2 "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/os"
	"github.com/spf13/cobra"
)

const (
	defaultHostServer          = "https://localhost:8443"
	defaultWebsocketHostServer = "wss://localhost:8443"
)

func main() {
	fmt.Println("Functional Conformance Suite CLI")

	insecureConn, err := NewConnection()
	if err == errInsecure {
		fmt.Println("server's certificate chain and host name not verified")
	} else if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	service := newService(
		os2.GetEnvOrDefault("FCS_HOST", defaultHostServer),
		os2.GetEnvOrDefault("FCS_WEBSOCKET_HOST", defaultWebsocketHostServer),
		insecureConn,
	)

	rootCmd := newRootCommand(service)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCommand(service Service) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fcs",
		Short: "Functional Conformance Suite CLI",
		Long:  `To use with pipelines and reproducible test runs`,
	}
	rootCmd.AddCommand(runCmd(service))
	rootCmd.AddCommand(versionCmd(service))
	return rootCmd
}
