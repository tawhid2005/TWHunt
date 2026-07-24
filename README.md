# TWHunt 🚀

**TWHunt** is an advanced, blazingly fast, and modular passive subdomain enumeration tool written in Go. It integrates **24 highly reliable, completely FREE APIs** without requiring any API keys. 

Designed for Bug Bounty Hunters and Pentesters, it grabs thousands of subdomains in just a few seconds and provides a beautiful, colorful, and animated terminal output.

## Features ✨

- **24 Built-in Free Sources**: Uses Regex & JSON scraping engines to fetch data from crt.sh, HackerTarget, URLScan, RapidDNS, ShodanCT, ThreatCrowd, and many more!
- **Extremely Fast**: Uses Go routines (`goroutines`) and channels for concurrent fetching.
- **Zero Configuration**: No API keys required. Just plug and play!
- **Live Host Verification**: Optionally resolve DNS to verify which subdomains are actively alive.
- **Save to File**: Save all discovered subdomains to a `.txt` file automatically.
- **Beautiful UI**: Colorful ASCII banners, live spinner animations, and clean summary tables.
- **Cross-Platform**: Works natively on Linux (Kali/Parrot), Windows, and macOS.

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

## 🎯 Usage

It's incredibly simple to use! 

**1. Basic Subdomain Enumeration**
Just pass the target domain using the `-d` flag. It will quickly find all subdomains.
```bash
twhunt -d target.com
```

**2. Enumerate & Verify Live Hosts**
By passing the `-v` flag, TWHunt will also perform DNS resolution to find out which of the discovered subdomains are actively alive and resolving.
```bash
twhunt -d target.com -v
```

**3. Save Results to a Text File**
You can use the `-o` flag to save the final discovered subdomains into a `.txt` file for later use.
```bash
twhunt -d target.com -o subdomains.txt
twhunt -d target.com -v -o live_subdomains.txt
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
