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

- **24 Built-in Free Sources**: Uses Regex & JSON scraping engines to fetch data from crt.sh, HackerTarget, URLScan, RapidDNS, ShodanCT, ThreatCrowd, and many more!
- **Extremely Fast**: Uses Go routines (`goroutines`) and channels for concurrent fetching.
- **Silent Mode**: Pipe results into other tools easily using `-silent`.
- **JSON Output**: Save or print results in valid JSON format using `-json`.
- **Multi-Domain Scanning**: Scan hundreds of domains at once using a text file with `-dL`.
- **Live Host Verification**: Optionally resolve DNS to verify which subdomains are actively alive.
- **Smart Retries**: Built-in timeout and retry logic bypasses temporary rate-limits.
- **Auto-Update**: Keep the tool updated natively using `-update`.
- **Zero Configuration**: No API keys required. Just plug and play!

## 🛠️ Installation on Kali Linux

You can easily install and run **TWHunt** on Kali Linux using one of the following methods:

### Method 1: Download the Pre-compiled Binary (Easiest)
If you don't want to install Go, you can just download the pre-compiled Linux binary.

```bash
# 1. Download the Linux binary
wget https://github.com/tawhid2005/TWHunt/raw/master/twhunt_linux_amd64 -O twhunt

# 2. Give executable permissions
chmod +x twhunt

# 3. Move it to your local bin directory to use it from anywhere
sudo mv twhunt /usr/local/bin/
```

### Method 2: Build from Source (Requires Go)
If you have Go installed on your system (`sudo apt install golang`), follow these steps:

```bash
# 1. Clone the repository
git clone https://github.com/tawhid2005/TWHunt.git

# 2. Navigate to the directory
cd TWHunt

# 3. Build the binary
go build -o twhunt main.go

# 4. Give executable permissions
chmod +x twhunt

# 5. Move it to your local bin directory
sudo mv twhunt /usr/local/bin/
```

## 🎯 Advanced Usage

**1. Basic Subdomain Enumeration**
```bash
twhunt -d target.com
```

**2. Enumerate & Verify Live Hosts**
```bash
twhunt -d target.com -v
```

**3. Save Results to a Text File**
```bash
twhunt -d target.com -o subdomains.txt
```

**4. Multi-Domain Scanning (Batch Mode)**
Have a list of domains in `domains.txt`? Scan them all at once!
```bash
twhunt -dL domains.txt
```

**5. Silent Mode & Pipelining**
Hide banners and logs. Print only raw subdomains. Perfect for pipelining into `httpx` or `nuclei`!
```bash
twhunt -d target.com -silent | httpx
```

**6. JSON Output**
Output results in a structured JSON array.
```bash
twhunt -d target.com -json
twhunt -d target.com -json -o results.json
```

**7. Auto-Update Tool**
Easily fetch and install the latest version from GitHub automatically.
```bash
twhunt -update
```

## 📡 Supported Free Sources

1. AbuseIPDB
2. AlienVault
3. Anubis
4. BeVigil
5. BufferOver
6. CertSpotter
7. CommonCrawl
8. Crtsh
9. FullHunt
10. HackerTarget
11. Netcraft
12. Omnisint
13. RapidDNS
14. Riddler
15. ShodanCT
16. ShrewdEye
17. SiteDossier
18. SubdomainCenter
19. Synapsint
20. ThreatCrowd
21. ThreatMiner
22. URLScan
23. VirusTotal
24. Wayback

## 👨‍💻 Author
**MD TALHA HUSSAIN TAWHID**
- **Location:** Sylhet, Bangladesh
- **Email:** tawhidh2005@gmail.com
- **Phone:** +8801711729858

## 🤝 Contributing
Pull requests are welcome! If you find a new free API source, feel free to add it to the `sources/` directory and submit a PR.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Happy Hunting! 🐞
