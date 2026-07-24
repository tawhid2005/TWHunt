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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// পাইথন স্ক্রিপ্টের মতো User-Agent সেট করা
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, http.ErrServerClosed // যেকোনো এররের জন্য ডামি এরর
	}

	return ioutil.ReadAll(resp.Body)
}
