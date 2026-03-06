package scraper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/prometheus/common/expfmt"
)

func GetFilteredMetrics(url string, keep []string) {

	log.Println("CP 1")
	resp, err := http.Get(url + "/metrics")
	log.Println(url + "/metrics")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("unexpected status %d from %s/metrics: %s", resp.StatusCode, url, string(body))
	}

	log.Println("CP 2")
	body, _ := io.ReadAll(resp.Body)
	log.Printf("raw body: %.500s", string(body))    // first 500 chars
	resp.Body = io.NopCloser(bytes.NewReader(body)) // restore for parser
	var parser expfmt.TextParser
	log.Println(resp.Body)
	mf, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("CP 3")
	keepSet := make(map[string]bool)
	for _, name := range keep {
		keepSet[name] = true
	}

	log.Println("CP 4")
	for name, family := range mf {
		if !keepSet[name] {
			continue
		}

		log.Println("CP 5")
		for _, m := range family.GetMetric() {
			fmt.Printf("%s %v\n", name, m.GetGauge().GetValue())
		}
	}
}
