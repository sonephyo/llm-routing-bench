package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	shortPrompt    = `Summarize the water cycle.`
	modelName      = "mistralai/Mistral-7B-v0.1"
	attackTimeout  = 3600 * time.Second
	scrapeInterval = 1 * time.Second
	scrapeBuffer   = 2 * time.Second
)

// longPromptBase is a ~6000-character passage used as the source for long prompts.
// makeLongPrompt trims it to ~60% of the target token size (≈5.5 chars/token),
// so each (tokenSize, "long") combination has a distinct prefill load.
const longPromptBase = `The history of computing spans nearly a century of remarkable innovation, beginning in the 1930s and 1940s when mathematicians and engineers first conceived of programmable machines. Alan Turing's theoretical work on computation, published in 1936, established the mathematical foundation for what would become computer science, introducing the concept of a universal machine capable of performing any computable function. Around the same time, Konrad Zuse built the Z3 in Germany, widely regarded as the first programmable digital computer, while in the United States, John Atanasoff and Clifford Berry developed an electronic computing device designed to solve systems of linear equations.

The 1940s saw the construction of several landmark machines. ENIAC, completed at the University of Pennsylvania in 1945, was among the first general-purpose electronic computers, weighing over 27 tons and occupying an entire room. It was programmed by physically rewiring its connections, a process that could take days. John von Neumann's influential 1945 report describing a stored-program architecture — where instructions and data share the same memory — provided the conceptual blueprint that nearly all modern computers still follow. The Manchester Baby, built in 1948, became the first machine to run a stored program.

The 1950s brought the dawn of commercial computing. IBM released the 701 in 1952, its first large-scale scientific computer, followed by the 650 which became one of the best-selling computers of the decade. Transistors, invented at Bell Labs in 1947, began replacing vacuum tubes, dramatically reducing size, power consumption, and heat while improving reliability. The programming language FORTRAN, developed by John Backus and his team at IBM and released in 1957, was the first high-level language to achieve widespread use, allowing scientists and engineers to write programs without assembling machine code by hand.

The 1960s marked a period of rapid expansion. Integrated circuits combined multiple transistors onto a single chip, further miniaturizing hardware. The development of time-sharing systems allowed multiple users to interact with a single mainframe simultaneously, laying groundwork for the interactive computing paradigm. COBOL emerged as the dominant language for business data processing. The first computer networks began to take shape, culminating in ARPANET, funded by the U.S. Department of Defense, which first connected nodes in 1969 and became the direct precursor to the internet.

By the 1970s, the microprocessor had arrived. Intel's 4004, released in 1971, placed an entire CPU on a single chip. The 8080 and Motorola 6800 followed, enabling the first hobbyist personal computers such as the Altair 8800 in 1975. Bill Gates and Paul Allen wrote a BASIC interpreter for the Altair and founded Microsoft. Steve Jobs and Steve Wozniak built the Apple I in a garage and founded Apple Computer in 1976. The Unix operating system, developed at Bell Labs throughout the 1970s, introduced hierarchical file systems, pipes, and a philosophy of small composable tools that profoundly influenced software design for decades.

The 1980s brought personal computing to the mainstream. IBM released its PC in 1981 using an open architecture, inadvertently creating an industry standard that allowed dozens of clone manufacturers to compete. Microsoft licensed MS-DOS to IBM and retained the right to sell it to other manufacturers, a strategic decision that made Microsoft the dominant software company of the era. Apple introduced the Macintosh in 1984 with a graphical user interface and mouse-driven interaction that demonstrated computing's potential beyond technical specialists. The C programming language became the lingua franca of systems programming, and the free software movement, championed by Richard Stallman, began challenging proprietary software models.

The 1990s transformed computing into a global phenomenon. Tim Berners-Lee invented the World Wide Web in 1989 and launched it publicly in 1991, creating a hyperlinked document system built on top of the internet. The release of the Mosaic web browser in 1993 made the web accessible to non-technical users and triggered explosive growth in internet adoption. Linux, created by Linus Torvalds in 1991 and developed collaboratively by thousands of contributors worldwide, demonstrated that open-source development could produce production-quality operating systems. Java, introduced by Sun Microsystems in 1995 with the promise of write-once-run-anywhere portability, became widely adopted for enterprise applications.

The 2000s saw the internet scale to billions of users and gave rise to new computing paradigms. Google, founded in 1998, built massive distributed systems to index and search the web, publishing research on MapReduce and the Google File System that inspired the Hadoop ecosystem and the broader field of big data engineering. Social networks, smartphones, and cloud computing reshaped how software was built and deployed. Amazon Web Services launched in 2006, pioneering infrastructure-as-a-service and enabling startups to build on rented compute rather than owned hardware.

The 2010s accelerated machine learning from an academic pursuit to an industrial workhorse. Deep learning, powered by large datasets, GPU computing, and architectural innovations like convolutional neural networks and later transformers, achieved superhuman performance on image recognition, speech transcription, and language translation. AlexNet's 2012 ImageNet victory marked a turning point. The transformer architecture, introduced in the 2017 paper Attention Is All You Need, became the foundation for large language models including BERT, GPT, and their successors, which demonstrated emergent reasoning and language abilities at unprecedented scale. Graphics processing units, originally designed for rendering video games, became the dominant hardware accelerator for training and serving neural networks, and specialized chips such as Google's Tensor Processing Units emerged to meet demand.`

// makeLongPrompt returns a prefix of longPromptBase sized to ~60% of tokenSize
// input tokens, keeping total context (input + output) safely within 4096 tokens.
// Approximation: 1 token ≈ 5.5 characters (measured on longPromptBase: 6000 chars / 1100 tokens).
// The cut snaps back to the nearest space to avoid splitting mid-word or mid-rune.
func makeLongPrompt(tokenSize int) string {
	targetChars := tokenSize * 60 / 100 * 11 / 2 // 60% of tokenSize * 5.5 chars/token
	if targetChars >= len(longPromptBase) {
		return longPromptBase
	}
	// Snap back to the nearest space so we don't cut mid-word or mid-rune.
	cut := targetChars
	for cut > 0 && longPromptBase[cut] != ' ' {
		cut--
	}
	return longPromptBase[:cut]
}

var routerURL = func() string {
	if u := os.Getenv("ROUTER_URL"); u != "" {
		return u
	}
	return "http://localhost:7999"
}()

func MakeTargeter(tokenSize int, promptType string) vegeta.Targeter {
	prompt := shortPrompt
	if promptType == "long" {
		prompt = makeLongPrompt(tokenSize)
	}
	body := fmt.Sprintf(
		`{"model": %q, "prompt": %q, "min_tokens": %d, "max_tokens": %d, "ignore_eos": true"}`,
		modelName, prompt, tokenSize, tokenSize,
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

