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
	expNum := nextExperimentNumber(baseDir)
	outputDir := fmt.Sprintf("%s/experiment_%d", baseDir, expNum)
	log.Printf("Experiment #%d — results will be saved to %s/", expNum, outputDir)

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

				targeter := MakeTargeter(ts, pt)

				preRM, err := ScrapeRouterMetrics(promAddr)
				if err != nil {
					log.Printf("warn: pre-scrape router metrics: %v", err)
				}
				preBM, _ := ScrapeBackendMetrics(promAddr)

				startTime := time.Now()

				var vm vegeta.Metrics
				var rawReqs []RawRequest
				switch lp {
				case "uniform":
					vm, rawReqs = RunUniform(targeter)
				case "bursty":
					vm, rawReqs = RunBursty(targeter)
				case "rampup":
					vm, rawReqs = RunRampUp(targeter)
				}

				endTime := time.Now()

				gauges := ScrapeGaugeRanges(promAddr, startTime, endTime)

				waitForQueueDrain(promAddr)
				log.Printf("  waiting %s for Prometheus scrape alignment...", scrapeInterval+scrapeBuffer)
				time.Sleep(scrapeInterval + scrapeBuffer)

				postRM, err := ScrapeRouterMetrics(promAddr)
				if err != nil {
					log.Printf("warn: post-scrape router metrics: %v", err)
				}
				postBM, _ := ScrapeBackendMetrics(promAddr)

				result := ExperimentResult{
					Metadata: ExperimentMetadata{
						Strategy:    strategy,
						LoadPattern: lp,
						TokenSize:   ts,
						PromptType:  pt,
						StartTime:   startTime,
						EndTime:     endTime,
					},
					VegetaMetrics:  vm,
					RouterMetrics:  DeltaRouterMetrics(preRM, postRM),
					BackendMetrics: DeltaBackendMetrics(preBM, postBM, gauges),
					RawRequests: rawReqs,
				}

				if err := WriteResult(result, outputDir); err != nil {
					log.Printf("error: write result for %s_%d_%s: %v", lp, ts, pt, err)
				} else {
					log.Printf("Note: saved  p99=%s throughput=%.2f req/s success=%.1f%%",
						vm.Latencies.P99,
						vm.Throughput,
						vm.Success*100,
					)
				}

				if run < total {
					log.Printf("  cooldown 60s before next run...")
					time.Sleep(60 * time.Second)
				}
			}
		}
	}

	log.Printf("Done. Results written to %s/", outputDir)
}

// nextExperimentNumber scans baseDir for subdirectories named "experiment_N"
// and returns the next available number (max existing + 1, or 1 if none).
func nextExperimentNumber(baseDir string) int {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return 1
	}
	max := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		var n int
		if _, err := fmt.Sscanf(e.Name(), "experiment_%d", &n); err == nil {
			if n > max {
				max = n
			}
		}
	}
	return max + 1
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
