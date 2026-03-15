package main

import (
	"dcl/backend/aletheia"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var sandbox bool
	var demo bool
	var web bool
	var port int
	var minEmit int
	var maxEmit int

	flag.BoolVar(&sandbox, "sandbox", true, "run sandbox simulation claim stream")
	flag.BoolVar(&demo, "demo", false, "print evidence records to console")
	flag.BoolVar(&web, "web", true, "start dashboard server")
	flag.IntVar(&port, "port", 8080, "dashboard port")
	flag.IntVar(&minEmit, "emit-min", 2, "minimum simulated claim interval in seconds")
	flag.IntVar(&maxEmit, "emit-max", 5, "maximum simulated claim interval in seconds")
	flag.Parse()

	fmt.Println("==========================================")
	fmt.Println(" ALETHEIA - Sandbox Misinformation Engine ")
	fmt.Println(" Agentic verification + cryptographic log ")
	fmt.Println("==========================================")

	chain := aletheia.NewDetectionChain()
	handler := aletheia.NewHandler(chain, demo)
	watcher := aletheia.NewWatcher(handler)

	report := aletheia.NewSandboxIntelligenceReport()
	verifier := aletheia.NewSandboxVerifier(report)
	memory := aletheia.NewImmuneMemory(chain)
	aletheia.SetSandboxVerifier(verifier)
	aletheia.SetImmuneMemory(memory)

	var generator *aletheia.ClaimStreamGenerator
	if sandbox {
		generator = aletheia.NewClaimStreamGenerator(
			watcher,
			report,
			time.Duration(minEmit)*time.Second,
			time.Duration(maxEmit)*time.Second,
		)
	}

	if web {
		server := aletheia.NewServer(chain, handler, watcher, report, port)
		handler.SetServer(server)
		go server.Start()
		time.Sleep(250 * time.Millisecond)
	}

	watcher.Start()

	if generator != nil {
		generator.Start()
		fmt.Printf(
			"[sandbox] streaming enabled | segments=%d events=%d claims=%d interval=%ds-%ds\n",
			len(report.Segments),
			len(report.Events),
			report.TotalClaims(),
			minEmit,
			maxEmit,
		)
	} else {
		fmt.Println("[sandbox] streaming disabled")
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\n[shutdown] stopping ALETHEIA services...")
	if generator != nil {
		generator.Stop()
	}
	watcher.Stop()
	fmt.Println("[shutdown] complete")
}
