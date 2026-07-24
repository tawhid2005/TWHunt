<div align="center">
  <img src="author.png" width="150" style="border-radius:50%; border: 3px solid #38b2ac;" alt="MD TALHA HUSSAIN TAWHID">
  <h1>TWHunt 🚀</h1>
  <p><strong>Developed by MD TALHA HUSSAIN TAWHID</strong></p>
  <p>📧 tawhidh2005@gmail.com | 📞 +8801711729858 | 🌐 <a href="https://github.com/tawhid2005">GitHub</a></p>
</div>

---

**TWHunt** is an advanced, blazingly fast, and modular passive subdomain enumeration tool written in Go. It integrates **24 highly reliable, completely FREE APIs** without requiring any API keys. 

Designed for Bug Bounty Hunters and Pentesters, it grabs thousands of subdomains in just a few seconds and provides a beautiful, colorful, and animated terminal output.

## Features ✨

- **API Key Configuration (`config.json`)**: Expand your OSINT capabilities by seamlessly adding API keys for premium sources like Shodan and SecurityTrails.
- **Active Brute-Forcing**: Use custom wordlists to actively hunt down hidden subdomains using a highly concurrent DNS resolver engine.
- **Fast TCP Port Scanning**: Automatically probe live subdomains for open web ports (`80, 443, 8080`) at lightning speeds.
- **Beautiful HTML Reports**: Generate a stunning, responsive HTML dashboard to visualize your findings.
- **Wildcard DNS Filtering**: Intelligently detect and filter out fake catch-all wildcard DNS responses.
- **24 Built-in Free Sources**: Uses Regex & JSON scraping engines to fetch data from crt.sh, HackerTarget, URLScan, RapidDNS, ShodanCT, ThreatCrowd, and many more!
- **Extremely Fast**: Uses Go routines (`goroutines`) and channels for concurrent fetching.
- **Silent Mode**: Pipe results into other tools easily using `-silent`.
- **JSON Output**: Save or print results in valid JSON format using `-json`.
- **Multi-Domain Scanning**: Scan hundreds of domains at once using a text file with `-dL`.
- **Live Host Verification**: Optionally resolve DNS to verify which subdomains are actively alive.
- **Smart Retries**: Built-in timeout and retry logic bypasses temporary rate-limits.
- **Auto-Update**: Keep the tool updated natively using `-update`.
- **Zero Configuration out of the box**: No API keys required for the 24 built-in sources. Just plug and play!

## 🛠️ Installation Guide (Kali Linux / Parrot OS)

You can easily install and run **TWHunt** using one of the following methods:

### Method 1: The Quickest Way (Pre-compiled Binary)
If you just want to run the tool without installing Go, follow these steps:

```bash
# 1. Download the Linux binary
wget https://github.com/tawhid2005/TWHunt/raw/master/twhunt_linux_amd64 -O twhunt

# 2. Make it executable
chmod +x twhunt

# 3. Move it to your local bin directory to use it from anywhere
sudo mv twhunt /usr/local/bin/
```

### Method 2: Build from Source (Requires Go)
If you are a developer or have Go installed on your system (`sudo apt install golang`), follow these steps:

```bash
# 1. Clone the repository
git clone https://github.com/tawhid2005/TWHunt.git

# 2. Navigate to the directory
cd TWHunt

# 3. Build the binary
go build -o twhunt main.go

# 4. Make it executable and move to bin
chmod +x twhunt
sudo mv twhunt /usr/local/bin/
```

## 🎯 How to Use TWHunt (Masterclass)

TWHunt is designed to be extremely user-friendly. You can type `twhunt -h` or `twhunt --help` anytime to view the beautiful help menu right inside your terminal!

Here are the most common workflows every Bug Hunter should know:

### 1. Basic Subdomain Enumeration
Want to find subdomains for a single target quickly? Use the `-d` flag.
```bash
twhunt -d target.com
```

### 2. Verify Live Hosts
Often, many subdomains are dead or parked. By adding the `-v` flag, TWHunt will resolve the DNS and filter out only the **alive and active** subdomains!
```bash
twhunt -d target.com -v
```

### 3. Save Results to a Text File
You can use the `-o` flag to save the final discovered subdomains into a `.txt` file for later use. (Note: If you forget the `-o` flag, TWHunt will smartly ask you at the end of the scan if you want to save it!)
```bash
twhunt -d target.com -o subdomains.txt
twhunt -d target.com -v -o live_subdomains.txt
```

### 4. Batch Scanning (Multi-Domain)
Got a large scope? Save all your root domains in a file (e.g. `domains.txt`) and let TWHunt scan them all!
```bash
twhunt -dL domains.txt
```

### 5. Silent Mode & Pipeline Magic
Professional hackers love pipelining tools. The `-silent` flag removes all banners, loading animations, and logs, printing ONLY the raw subdomains. You can pipe this directly into tools like `httpx` or `nuclei`.
```bash
twhunt -d target.com -silent | httpx -title -status-code
```

### 6. JSON Output
Output results in a structured JSON array.
```bash
twhunt -d target.com -json
twhunt -d target.com -json -o results.json
```

**7. Active Subdomain Brute-Forcing**
Hunt hidden subdomains using your custom wordlist.
```bash
twhunt -d target.com -w wordlist.txt
```

**8. TCP Port Scanning**
Scan for open web ports on discovered live subdomains.
```bash
twhunt -d target.com -ports 80,443,8080,8443
```

**9. HTML Dashboard Report**
Generate an interactive, beautifully designed HTML report.
```bash
twhunt -d target.com -ports 80,443 -html
twhunt -d target.com -ports 80,443 -html -o my_report.html
```

**10. Wildcard DNS Filtering**
Prevent catch-all DNS records from flooding your results.
```bash
twhunt -d target.com -nw
```

**11. Auto-Update Tool**
Easily fetch and install the latest version from GitHub automatically.
```bash
twhunt -update
```

## 📡 Supported OSINT Sources (100% Free)

1. AbuseIPDB | 2. AlienVault | 3. Anubis | 4. BeVigil | 5. BufferOver | 6. CertSpotter | 7. CommonCrawl | 8. Crtsh | 9. FullHunt | 10. HackerTarget | 11. Netcraft | 12. Omnisint | 13. RapidDNS | 14. Riddler | 15. ShodanCT | 16. ShrewdEye | 17. SiteDossier | 18. SubdomainCenter | 19. Synapsint | 20. ThreatCrowd | 21. ThreatMiner | 22. URLScan | 23. VirusTotal | 24. Wayback

## 👨‍💻 Developed By
**MD TALHA HUSSAIN TAWHID**
- **Location:** Sylhet, Bangladesh
- **Email:** tawhidh2005@gmail.com
- **Phone:** +8801711729858

## 🤝 Contributing
Found a new API? Pull requests are always welcome! Let's build the best Subdomain Enumeration tool together.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Happy Hunting! 🐞
