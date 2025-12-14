package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: direct-share <command> [arguments]")
		fmt.Println("Commands: listen, send")
		os.Exit(1)
	}

	listenCmd := flag.NewFlagSet("listen", flag.ExitOnError)
	listenPort := listenCmd.String("port", ":9000", "Port to listen on")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendAddr := sendCmd.String("addr", "localhost:9000", "Address to send to")
	sendFile := sendCmd.String("file", "", "File to send")

	switch os.Args[1] {
	case "listen":
		listenCmd.Parse(os.Args[2:])
		if err := listen(*listenPort); err != nil {
			fmt.Println("Listen error:", err)
			os.Exit(1)
		}
	case "send":
		sendCmd.Parse(os.Args[2:])
		if *sendAddr == "" {
			fmt.Println("Error: -addr argument is required for sending")
			sendCmd.PrintDefaults()
			os.Exit(1)
		}
		if *sendFile == "" {
			fmt.Println("Error: -file argument is required for sending")
			sendCmd.PrintDefaults()
			os.Exit(1)
		}
		if err := send(*sendAddr, *sendFile); err != nil {
			fmt.Println("Send error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
