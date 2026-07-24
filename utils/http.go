package utils

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPClient হচ্ছে একটি কাস্টম ক্লায়েন্ট যা SSL এরর ইগনোর করে এবং ৩ সেকেন্ডের টাইমআউট রাখে
var HTTPClient = &http.Client{
	Timeout: 35 * time.Second, // পাইথন স্ক্রিপ্টের মতো ৩৫ সেকেন্ড
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

// RequestAPI একটি জেনেরিক ফাংশন যা যেকোনো URL এ GET রিকোয়েস্ট করে রেসপন্স বডি রিটার্ন করে
func RequestAPI(url string) ([]byte, error) {
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		// পাইথন স্ক্রিপ্টের মতো User-Agent সেট করা
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

		resp, err := HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second) // ফেইল করলে ২ সেকেন্ড ওয়েট করে আবার ট্রাই করবে
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = http.ErrServerClosed
			time.Sleep(2 * time.Second)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		
		if err != nil {
			lastErr = err
			continue
		}
		
		return body, nil
	}

	return nil, lastErr
}
