package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/userinternal/alert"
	{
	var (
		portRange = flag.String("ports", "1-1024", "Port range to scan (e.g. 1-1024)")
		interval  = flag.Duration("interval", 30*time.Second, "How often to scan ports")
		stateFile = flag.String("state", "/tmp/portwatch.state", "Path to persist port state")
		workers   = flag.Int("workers", 100, "Number of concurrent scan workers")
	)
	flag.Parse()

	notifier := alert.New(os.Stdout)
	sc, err := scanner.New(*portRange, *workers)
	if err != nil {
		log.Fatalf("invalid port range: %v", err)
	}

	// Load previously saved state if it exists.
	prev, err := state.Load(*stateFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("warning: could not load state file: %v", err)
	}

	log.Printf("portwatch started — scanning %s every %s", *portRange, *interval)

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Run an immediate first scan before waiting for the ticker.
	prev = runScan(sc, notifier, prev, *stateFile, prev == nil)

	for {
		select {
		case <-ticker.C:
			prev = runScan(sc, notifier, prev, *stateFile, false)
		case sig := <-sigs:
			fmt.Fprintf(os.Stderr, "\nreceived %s, shutting down\n", sig)
			return
		}
	}
}

// runScan performs a single scan cycle: scan ports, compare with previous
// state, notify on changes, persist new state, and return the new snapshot.
func runScan(
	sc *scanner.Scanner,
	notifier *alert.Notifier,
	prev []uint16,
	stateFile string,
	baseline bool,
) []uint16 {
	current, err := sc.Scan()
	if err != nil {
		log.Printf("scan error: %v", err)
		return prev
	}

	if baseline {
		log.Printf("baseline established: %d open port(s)", len(current))
	} else {
		opened, closed := state.Compare(prev, current)
		notifier.Notify(opened, closed)
	}

	if err := state.Save(stateFile, current); err != nil {
		log.Printf("warning: could not save state: %v", err)
	}

	return current
}
