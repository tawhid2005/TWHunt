<div align="center">
  <img src="author.png" width="150" style="border-radius:50%; border: 3px solid #38b2ac;" alt="MD TALHA HUSSAIN TAWHID">
  <h1>TWHunt - The Ultimate Reconnaissance & Bug Bounty Framework 🚀</h1>
  <p><strong>Developed by MD TALHA HUSSAIN TAWHID</strong></p>
  <p>📧 tawhidh2005@gmail.com | 📞 +8801711729858 | 🌐 <a href="https://github.com/tawhid2005">GitHub</a></p>
  <br>
  <p>
    <a href="https://github.com/tawhid2005/TWHunt/releases"><img src="https://img.shields.io/github/v/tag/tawhid2005/TWHunt?color=38b2ac&label=Version&style=for-the-badge" alt="Version"></a>
    <a href="https://github.com/tawhid2005/TWHunt/stargazers"><img src="https://img.shields.io/github/stars/tawhid2005/TWHunt?color=gold&label=Stars&style=for-the-badge" alt="Stars"></a>
    <a href="https://github.com/tawhid2005/TWHunt/network/members"><img src="https://img.shields.io/github/forks/tawhid2005/TWHunt?color=blue&label=Forks&style=for-the-badge" alt="Forks"></a>
    <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go" alt="Go">
  </p>
</div>

---

## 🔎 What is TWHunt?
**TWHunt** is a blazingly fast, highly concurrent **Subdomain Enumeration, OSINT (Open Source Intelligence), and Vulnerability Reconnaissance Framework** built in Go (Golang). 

Designed specifically for **Bug Bounty Hunters, Penetration Testers, and Red Teamers**, TWHunt completely automates your recon pipeline. It seamlessly integrates **24 Free API sources** to discover subdomains, and immediately chains them into active scans—finding **JavaScript Secrets, Subdomain Takeovers, Exposed Vulnerabilities, and Wayback Machine Endpoints** in seconds.

## 🔥 God-Tier Features

### 🌍 OSINT & Subdomain Discovery
- **24 Built-in Free Sources**: Uses Regex & JSON scraping engines (crt.sh, HackerTarget, URLScan, RapidDNS, ShodanCT, ThreatCrowd, etc.) without requiring any API keys!
- **Active Brute-Forcing**: Use custom wordlists and concurrent DNS resolvers to actively hunt down hidden subdomains.
- **Subdomain Permutation / Alteration**: Actively generates and resolves mutations (e.g., `api-dev`, `dev-staging`) to uncover unlisted endpoints.
- **Wildcard DNS Filtering**: Intelligently detect and filter out fake catch-all wildcard DNS responses to keep your results clean.

### 💣 Automated Vulnerability Probing
- **JavaScript Secrets Finder**: Actively fetches and parses `.js` files across subdomains to find hardcoded **AWS Keys, Stripe Tokens, Google APIs, and GitHub Tokens**.
- **Cloud Bucket Enumeration**: Automatically generates permutations and checks AWS S3, Google Cloud, and Azure for open, exposed buckets (`-buckets`).
- **Web Application Firewall (WAF) Detection**: Actively triggers and detects WAFs (Cloudflare, Akamai, Imperva) using XSS payloads (`-waf`).
- **CORS Misconfiguration Checker**: Actively probes live subdomains for Cross-Origin Resource Sharing vulnerabilities (`-cors`).
- **Hidden Parameter Fuzzing**: Discovers hidden developer or admin parameters (e.g., `?admin=1`, `?debug=true`) that alter responses (`-params`).
- **DNS Zone Transfer (AXFR)**: Identifies misconfigured Nameservers that allow DNS zone dumping (`-axfr`).
- **Subdomain Takeover Detection**: Automatically checks DNS CNAME records against **20+ vulnerable services** (GitHub Pages, AWS S3, Heroku) and alerts you instantly.
- **Vulnerability Prober**: Automatically probes discovered subdomains for low-hanging bugs (e.g., exposed `/.env`, `/.git/config`, `/phpinfo.php`).
- **Tech Stack Detection**: Analyzes HTTP headers and cookies to instantly identify underlying technologies (Nginx, React, PHP, ASP.NET).
- **HTTP Status Probing**: Concurrently probes live subdomains for HTTP/HTTPS status codes (e.g. 200, 403, 404) and extracts page `<title>` tags.
- **Fast TCP Port Scanning**: Automatically probe live subdomains for open web ports (`80, 443, 8080`) at lightning speeds.

### 📊 Professional Recon Pipeline
- **Wayback URLs Fetching**: Deep dive into a target's history by extracting thousands of hidden historical endpoints from the Wayback Machine.
- **Discord Notification System**: Run on a VPS. It saves scan states and alerts you on Discord via Webhooks the second a **brand new subdomain** is deployed.
- **Beautiful HTML Reports**: Generate a stunning, responsive HTML dashboard to visually analyze your findings.
- **JSON Output & Silent Mode**: Perfect for chaining with other tools like `httpx` or `nuclei`.
- **Auto-Update**: Keep the tool updated natively using `-update`.

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
twhunt -d target.com -ports 80,443 -takeover -probe -html
```

**10. Wildcard DNS Filtering**
Prevent catch-all DNS records from flooding your results.
```bash
twhunt -d target.com -nw
```

**11. Subdomain Takeover Detection**
Scan CNAME records to find high-impact Takeover vulnerabilities.
```bash
twhunt -d target.com -takeover
```

**12. HTTP Probing (Status & Titles)**
Probe live subdomains to find out their HTTP status codes and page titles.
```bash
twhunt -d target.com -probe
```

**13. Wayback Machine URLs Fetching**
Extract historical hidden endpoints and sensitive files from the Wayback Machine.
```bash
twhunt -d target.com -urls -o urls.txt
```

**14. JS Secrets Finder 🕵️‍♂️**
Actively parse `.js` files on discovered subdomains to find API Keys and leaked tokens.
```bash
twhunt -d target.com -js
```

**15. Cloud Bucket Enumeration ☁️**
Find open AWS S3, Google Cloud, and Azure buckets related to the target.
```bash
twhunt -d target.com -buckets
```

**16. WAF Detection & CORS Misconfigurations 🛡️**
Actively test subdomains for Web Application Firewalls and Cross-Origin Resource Sharing vulnerabilities.
```bash
twhunt -d target.com -waf -cors
```

**17. Tech Stack Detection ⚙️**
Analyze headers and cookies to instantly identify the technologies running on the target.
```bash
twhunt -d target.com -tech
```

**18. Vulnerability Prober 💣**
Automatically check live subdomains for low-hanging bugs (e.g. `.env`, `.git`).
```bash
twhunt -d target.com -vuln
```

**19. Hidden Parameter Fuzzing 🎛️**
Fuzz for hidden sensitive parameters (e.g. `?admin=1`).
```bash
twhunt -d target.com -params
```

**20. DNS Zone Transfer (AXFR) 🌍**
Check if the domain's Name Servers are vulnerable to a full DNS Zone Transfer.
```bash
twhunt -d target.com -axfr
```

**21. Subdomain Alteration / Permutation 🧬**
Mutate valid subdomains to find hidden ones (e.g. changes `dev.target.com` to `api-dev.target.com`).
```bash
twhunt -d target.com -alt
```

**22. Continuous Monitoring & Discord Alerts 🤖**
Run TWHunt periodically. It will save the state and alert you via Discord webhook if a brand new subdomain is deployed!
```bash
twhunt -d target.com -notify "https://discord.com/api/webhooks/your_webhook_here"
```

**23. Auto-Update Tool**
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


