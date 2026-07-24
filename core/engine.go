package core

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// ANSI Colors
const (
	Mint     = "\033[38;5;121m"
	Sky      = "\033[38;5;117m"
	Gold     = "\033[38;5;222m"
	Coral    = "\033[38;5;210m"
	Lavender = "\033[38;5;183m"
	Silver   = "\033[38;5;249m"
	Slate    = "\033[38;5;241m"
	Bold     = "\033[1m"
	EndC     = "\033[0m"
)

type Result struct {
	SourceName string
	Data       []string
	Error      error
}

type Engine struct {
	Sources []Source
}

func NewEngine(sources []Source) *Engine {
	return &Engine{Sources: sources}
}

func (e *Engine) Run(domain string) []string {
	var wg sync.WaitGroup
	resultsChan := make(chan Result, len(e.Sources))
	
	fmt.Printf("\n%s[*] TARGETING : %s%s\n", Sky, strings.ToUpper(domain), EndC)
	fmt.Printf("%s[*] SOURCES   : %d PASSIVE ENGINES ACTIVATED%s\n", Gold, len(e.Sources), EndC)
	fmt.Printf("%s%s%s\n", Slate, strings.Repeat("-", 60), EndC)

	startTime := time.Now()

	// স্পিনার অ্যানিমেশনের জন্য
	stopAnim := make(chan bool)
	go func() {
		frames := []string{"⠏", "⠛", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠍"}
		i := 0
		for {
			select {
			case <-stopAnim:
				fmt.Print("\r" + strings.Repeat(" ", 50) + "\r")
				return
			default:
				fmt.Printf("\r %s%s%s %sEXTRACTING DEEP DATA...%s", Mint, frames[i%len(frames)], EndC, Gold, EndC)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	for _, source := range e.Sources {
		wg.Add(1)
		go func(s Source) {
			defer wg.Done()
			data, err := s.Run(domain)
			resultsChan <- Result{
				SourceName: s.Name(),
				Data:       data,
				Error:      err,
			}
		}(source)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	uniqueSubdomains := make(map[string]bool)
	var allSubdomains []string
	
	// সোর্স অনুযায়ী রেজাল্ট রাখার ম্যাপ
	sourceSummary := make(map[string]int)

	for result := range resultsChan {
		if result.Error != nil {
			sourceSummary[result.SourceName] = -1 // -1 মানে এরর
			continue
		}
		
		for _, sub := range result.Data {
			sub = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(sub, ".")))
			if sub != "" && sub != domain && strings.HasSuffix(sub, domain) && !uniqueSubdomains[sub] {
				uniqueSubdomains[sub] = true
				allSubdomains = append(allSubdomains, sub)
			}
		}
		sourceSummary[result.SourceName] = len(result.Data)
	}

	stopAnim <- true // স্পিনার বন্ধ করা

	// সামারি প্রিন্ট করা
	for _, s := range e.Sources {
		name := s.Name()
		count := sourceSummary[name]
		nameUpper := strings.ToUpper(name)
		
		if count > 0 {
			fmt.Printf(" %s[✓]%s %-18s : %s%d%s FOUND\n", Mint, EndC, nameUpper, Bold, count, EndC)
		} else {
			fmt.Printf(" %s[✗]%s %-18s : 0 FOUND\n", Coral, EndC, nameUpper)
		}
	}

	fmt.Printf("%s%s%s\n", Slate, strings.Repeat("-", 60), EndC)
	fmt.Printf(" %s[★] TOTAL UNIQUE DISCOVERED: %d%s\n", Lavender, len(allSubdomains), EndC)
	fmt.Printf(" %s[*] TIME TAKEN: %v%s\n", Slate, time.Since(startTime), EndC)

	return allSubdomains
}
