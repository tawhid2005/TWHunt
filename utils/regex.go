package utils

import (
	"regexp"
	"strings"
)

// ExtractAll একটি টেক্সট থেকে ডোমেইন এর সাথে ম্যাচ করা সব সাবডোমেইন খুঁজে বের করে
func ExtractAll(text string, domain string) []string {
	// পাইথন স্ক্রিপ্টের প্যাটার্ন: r'(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+' + re.escape(domain)
	pattern := `(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+` + regexp.QuoteMeta(domain)
	re := regexp.MustCompile(pattern)
	
	matches := re.FindAllString(strings.ToLower(text), -1)
	
	// ডুপ্লিকেট রিমুভ করা
	unique := make(map[string]bool)
	var result []string
	
	for _, match := range matches {
		match = strings.TrimSpace(match)
		if match != "" && !unique[match] {
			unique[match] = true
			result = append(result, match)
		}
	}
	
	return result
}
