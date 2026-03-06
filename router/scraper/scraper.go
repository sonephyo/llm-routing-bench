package scraper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

func GetFilteredMetrics(url string, keep []string) {

	resp, err := http.Get(url + "/metrics")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("unexpected status %d from %s/metrics: %s", resp.StatusCode, url, string(body))
	}

	parser := expfmt.NewTextParser(model.UTF8Validation)
	mf, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	keepSet := make(map[string]bool)
	for _, name := range keep {
		keepSet[name] = true
	}

	for name, family := range mf {
		if !keepSet[name] {
			continue
		}

		for _, m := range family.GetMetric() {
			fmt.Printf("%s %v\n", name, m.GetGauge().GetValue())
		}
	}
}
