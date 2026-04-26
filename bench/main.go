// bench/main.go is an entry for running the benchmark
// - checks .env for reading the LBStrategy that is running

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type IterationConfig struct {
	targeter    vegeta.Targeter
	loadPattern string
	tokenSize   int
	promptType  string
	promAddr    string
	strategy    string
	outputDir   string
}

func main() {
	env := loadEnv("../.env")

	strategy := env["LB_STRATEGY"]
	if strategy == "" {
		strategy = os.Getenv("LB_STRATEGY")
	}
	if strategy == "" {
		strategy = "unknown"
	}

	mode := env["MODE"]
	if mode == "" {
		mode = os.Getenv("MODE")
	}

	promAddr := "http://localhost:7779"
	if mode == "local" {
		promAddr = "http://localhost:7998"
	}
	if addr := os.Getenv("PROMETHEUS_ADDR"); addr != "" {
		promAddr = addr
	}

	baseDir := "bench-results/" + strategy
	experimentID := env["EXPERIMENT_ID"]
	if experimentID == "" {
		experimentID = os.Getenv("EXPERIMENT_ID")
	}
	if experimentID == "" {
		log.Fatal("EXPERIMENT_ID must be set in .env or environment")
	}
	outputDir := fmt.Sprintf("%s/experiment_%s", baseDir, experimentID)
	log.Printf("Experiment #%s — results will be saved to %s/", experimentID, outputDir)

	loadPatterns := []string{"uniform", "bursty", "rampup"}
	tokenSizes := []int{100, 500, 2000}
	promptTypes := []string{"short", "long"}

	total := len(loadPatterns) * len(tokenSizes) * len(promptTypes)
	run := 0

	for _, lp := range loadPatterns {
		for _, ts := range tokenSizes {
			for _, pt := range promptTypes {
				run++
				log.Printf("[%d/%d] strategy=%s pattern=%s tokens=%d prompt=%s",
					run, total, strategy, lp, ts, pt)
				iterationInstance := IterationConfig{
					targeter:    MakeTargeter([]int{ts}, []string{pt}),
					loadPattern: lp,
					tokenSize:   ts,
					promptType:  pt,
					promAddr:    promAddr,
					strategy:    strategy,
					outputDir:   outputDir,
				}
				iterationInstance.startIteration()
				if run < total {
					log.Printf("  cooldown 60s before next run...")
					time.Sleep(60 * time.Second)
				}
			}
		}
	}

	run = 0
	total = len(loadPatterns)
	for _, lp := range loadPatterns {
		run++
		iterationInstance := IterationConfig{
			targeter:    MakeTargeter(tokenSizes, promptTypes),
			loadPattern: lp,
			tokenSize: -1,
			promptType: "heterogenous",
			promAddr:    promAddr,
			strategy:    strategy,
			outputDir:   outputDir,
		}
		iterationInstance.startIteration()
		if run < total {
			log.Printf("  cooldown 60s before next run...")
			time.Sleep(60 * time.Second)
		}

	}

	log.Printf("Done. Results written to %s/", outputDir)
}

func (iConfig *IterationConfig) startIteration() {
	preRM, err := ScrapeRouterMetrics(iConfig.promAddr)
	if err != nil {
		log.Printf("warn: pre-scrape router metrics: %v", err)
	}
	preBM, _ := ScrapeBackendMetrics(iConfig.promAddr)

	startTime := time.Now()

	var vm vegeta.Metrics
	var rawReqs []RawRequest
	switch iConfig.loadPattern {
	case "uniform":
		vm, rawReqs = RunUniform(iConfig.targeter)
	case "bursty":
		vm, rawReqs = RunBursty(iConfig.targeter)
	case "rampup":
		vm, rawReqs = RunRampUp(iConfig.targeter)
	}

	endTime := time.Now()

	gauges := ScrapeGaugeRanges(iConfig.promAddr, startTime, endTime)

	waitForQueueDrain(iConfig.promAddr)
	log.Printf("  waiting %s for Prometheus scrape alignment...", scrapeInterval+scrapeBuffer)
	time.Sleep(scrapeInterval + scrapeBuffer)

	postRM, err := ScrapeRouterMetrics(iConfig.promAddr)
	if err != nil {
		log.Printf("warn: post-scrape router metrics: %v", err)
	}
	postBM, _ := ScrapeBackendMetrics(iConfig.promAddr)

	result := ExperimentResult{
		Metadata: ExperimentMetadata{
			Strategy:    iConfig.strategy,
			LoadPattern: iConfig.loadPattern,
			TokenSize:   iConfig.tokenSize,
			PromptType:  iConfig.promptType,
			StartTime:   startTime,
			EndTime:     endTime,
		},
		VegetaMetrics:  vm,
		RouterMetrics:  DeltaRouterMetrics(preRM, postRM),
		BackendMetrics: DeltaBackendMetrics(preBM, postBM, gauges),
		RawRequests:    rawReqs,
	}

	if err := WriteResult(result, iConfig.outputDir); err != nil {
		log.Printf("error: write result for %s_%d_%s: %v", iConfig.loadPattern, iConfig.tokenSize, iConfig.promptType, err)
	} else {
		log.Printf("Note: saved  p99=%s throughput=%.2f req/s success=%.1f%%",
			vm.Latencies.P99,
			vm.Throughput,
			vm.Success*100,
		)
	}
}

// loadEnv reads a .env file and returns a map of key=value pairs.
// Lines starting with # are ignored.
func loadEnv(path string) map[string]string {
	env := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return env
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return env
}
