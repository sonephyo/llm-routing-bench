package scraper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/common/expfmt"
)

func GetFilteredMetrics(url string, keep []string) {
	resp, err := http.Get(url + "/metrics")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var parser expfmt.TextParser
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
