package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	shortPrompt   = `Summarize the water cycle.`
	longPrompt    = `The following is a detailed history of computer science from the 1940s through the present day, covering key innovations in hardware, software, networking, and artificial intelligence.`
	modelName     = "mistralai/Mistral-7B-v0.1"
	attackTimeout = 3600 * time.Second
)

var routerURL = func() string {
	if u := os.Getenv("ROUTER_URL"); u != "" {
		return u
	}
	return "http://localhost:7999"
}()

func MakeTargeter(tokenSize int, promptType string) vegeta.Targeter {
	prompt := shortPrompt
	if promptType == "long" {
		prompt = longPrompt
	}
	body := fmt.Sprintf(
		`{"model": %q, "prompt": %q, "max_tokens": %d}`,
		modelName, prompt, tokenSize,
	)
	return vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    routerURL,
		Body:   []byte(body),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	})
}

func runPhase(attacker *vegeta.Attacker, targeter vegeta.Targeter, rate vegeta.Rate, dur time.Duration, metrics *vegeta.Metrics) {
	for res := range attacker.Attack(targeter, rate, dur, "") {
		metrics.Add(res)
	}
}

func RunUniform(targeter vegeta.Targeter) vegeta.Metrics {
	var m vegeta.Metrics
	attacker := vegeta.NewAttacker(vegeta.Timeout(attackTimeout))
	runPhase(attacker, targeter, vegeta.Rate{Freq: 10, Per: time.Second}, 10*time.Second, &m)
	m.Close()
	return m
}

func RunBursty(targeter vegeta.Targeter) vegeta.Metrics {
	var m vegeta.Metrics
	phases := []struct {
		rate vegeta.Rate
		dur  time.Duration
	}{
		{vegeta.Rate{Freq: 2, Per: time.Second}, 30 * time.Second},
		{vegeta.Rate{Freq: 20, Per: time.Second}, 10 * time.Second},
		{vegeta.Rate{Freq: 2, Per: time.Second}, 20 * time.Second},
	}
	for _, p := range phases {
		attacker := vegeta.NewAttacker(vegeta.Timeout(attackTimeout))
		runPhase(attacker, targeter, p.rate, p.dur, &m)
	}
	m.Close()
	return m
}

func RunRampUp(targeter vegeta.Targeter) vegeta.Metrics {
	var m vegeta.Metrics
	rates := []int{1, 4, 7, 10, 13, 15}
	for _, r := range rates {
		attacker := vegeta.NewAttacker(vegeta.Timeout(attackTimeout))
		runPhase(attacker, targeter, vegeta.Rate{Freq: r, Per: time.Second}, 15*time.Second, &m)
	}
	m.Close()
	return m
}

