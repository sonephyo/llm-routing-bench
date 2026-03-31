package scraper

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

func GetFilteredMetrics(url string, keep []string) (map[string]float64, error) {

	resp, err := http.Get(url + "/metrics")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d from %s/metrics: %s", resp.StatusCode, url, string(body))
	}

	parser := expfmt.NewTextParser(model.UTF8Validation)
	mf, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil && len(mf) == 0 {
		return nil, err
	}

	keepSet := make(map[string]bool)
	for _, name := range keep {
		keepSet[name] = true
	}

	result := make(map[string]float64)
	for name, family := range mf {
		if !keepSet[name] {
			continue
		}

		for _, m := range family.GetMetric() {
			result[name] = m.GetGauge().GetValue()
		}
	}

	return result, nil
}
