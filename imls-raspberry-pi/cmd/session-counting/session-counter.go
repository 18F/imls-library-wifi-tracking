package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gsa.gov/18f/config"
	"gsa.gov/18f/http"
	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/tlp"
	"gsa.gov/18f/version"
)

func run(ka *tlp.Keepalive, cfg *config.Config) {
	log.Println("Starting run")
	// Create channels for process network
	ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan []string)
	ch_macs_counted := make(chan map[string]int)
	ch_data_for_report := make(chan []map[string]string)

	// WARNING: If you get this length wrong, we have deadlock.
	// That is, every one of these needs to be used/written to/read from.
	const RESET_CHANS = 3
	// The kill channel lets us poison the network for shutdown. Really only for testing.
	var NIL_KILL_CHANNEL chan tlp.Ping = nil

	var chs_reset [RESET_CHANS]chan tlp.Ping
	for ndx := 0; ndx < RESET_CHANS; ndx++ {
		chs_reset[ndx] = make(chan tlp.Ping)
	}

	// Run the process network.
	// Driven by a 1s `tick` process.
	// Thread the keepalive through the network
	go tlp.Tick(ka, ch_sec)
	go tlp.TockEveryN(ka, 60, ch_sec, ch_nsec, NIL_KILL_CHANNEL)
	go tlp.RunWireshark(ka, cfg, ch_nsec, ch_macs, NIL_KILL_CHANNEL)
	// The reset will never be triggered in AlgoTwo unless we're rnuning in "sqlite" storage mode.
	go tlp.AlgorithmTwo(ka, cfg, ch_macs, ch_macs_counted, chs_reset[1], NIL_KILL_CHANNEL)
	go tlp.PrepareDataForStorage(ka, cfg, ch_macs_counted, ch_data_for_report, NIL_KILL_CHANNEL)
	if cfg.StorageMode == "api" {
		go tlp.StoreToCloud(ka, cfg, ch_data_for_report, chs_reset[2], NIL_KILL_CHANNEL)
	} else if cfg.StorageMode == "sqlite" {
		// At midnight, flush internal structures and restart.
		go tlp.PingAtMidnight(ka, cfg, chs_reset[0], NIL_KILL_CHANNEL)
		go tlp.StoreToSqlite(ka, cfg, ch_data_for_report, chs_reset[2], NIL_KILL_CHANNEL)
		// Fan out the ping to multiple PROCs
		go tlp.ParDelta(NIL_KILL_CHANNEL, chs_reset[:]...)
	}
}

func keepalive(ka *tlp.Keepalive, cfg *config.Config) {
	log.Println("Starting keepalive")
	var counter int64 = 0
	for {
		time.Sleep(time.Duration(cfg.Monitoring.PingInterval) * time.Second)
		ka.Publish(counter)
		counter = counter + 1
	}
}

func handleFlags() *config.Config {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	verbosePtr := flag.Bool("verbose", false, "Set log verbosity.")
	showKeyPtr := flag.Bool("show-key", false, "Tests key decryption.")
	configPathPtr := flag.String("config", "", "Path to config.yaml. REQUIRED.")
	storagePtr := flag.String("storage-mode", "api", "Either 'api', 'sqlite', or 'both'.")

	flag.Parse()

	if *configPathPtr == "" {
		log.Println("The flag --config MUST be provided.")
		os.Exit(1)
	}

	config.Verbose = *verbosePtr

	// If they just want the version, print and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	if _, err := os.Stat(*configPathPtr); os.IsNotExist(err) {
		log.Println("Looked for config at:", *configPathPtr)
		log.Fatal("Cannot find config file. Exiting.")
	} else {
		config.SetConfigPath(*configPathPtr)
	}

	cfg, err := config.ReadConfig(*configPathPtr)
	if err != nil {
		log.Fatal("session-counter: error loading config.")
	}

	if *showKeyPtr {
		fmt.Println(cfg.Auth.Token)
		os.Exit(0)
	}

	if *storagePtr == "api" || *storagePtr == "sqlite" || *storagePtr == "both" {
		cfg.StorageMode = *storagePtr
	}

	return cfg

}

func main() {
	// Read in a config
	cfg := handleFlags()

	// Set the session ID for this entire run
	if cfg.StorageMode == "sqlite" {
		t := time.Now()
		cfg.SessionId = fmt.Sprintf("%v-%v-%v", t.Year(), t.Month(), t.Day())
	} else {
		cfg.SessionId = config.CreateSessionId()
	}
	// Store this so we don't keep hitting /proc/cpuinfo
	cfg.Serial = config.GetSerial()
	// Make sure the mfg database is in place and can be loaded.
	api.CheckMfgDatabaseExists(cfg)

	el := http.NewEventLogger(cfg)
	el.Log("startup", nil)

	ka := tlp.NewKeepalive(cfg)
	go ka.Start()
	go keepalive(ka, cfg)
	go run(ka, cfg)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
