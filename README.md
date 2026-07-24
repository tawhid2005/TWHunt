# Subfnder 🚀

Subfnder is a blazingly fast, concurrent, and modular passive subdomain enumeration tool written in Go. It integrates **24 highly reliable, completely FREE APIs** without requiring any API keys. 

Designed for Bug Bounty Hunters and Pentesters, it grabs thousands of subdomains in just a few seconds and provides a beautiful, animated terminal output.

## Features ✨

- **24 Built-in Free Sources**: Uses Regex & JSON scraping engines to fetch data from crt.sh, HackerTarget, URLScan, RapidDNS, ShodanCT, ThreatCrowd, and many more!
- **Extremely Fast**: Uses Go routines (`goroutines`) and channels for concurrent fetching.
- **Zero Configuration**: No API keys required. Plug and play!
- **Live Host Verification**: Optionally resolve DNS to verify which subdomains are alive.
- **Beautiful UI**: Colorful ASCII banners, live spinner animations, and clean summary tables.
- **Cross-Platform**: Works natively on Linux, Windows, and macOS.

## Installation 🛠️

Ensure you have [Go](https://golang.org/doc/install) installed on your system.

### Install from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/subfnder.git

# Navigate to the directory
cd subfnder

# Build the binary
go build -o subfnder main.go

# (Optional) Move to bin path for global usage (Linux/macOS)
sudo mv subfnder /usr/local/bin/
```

## Usage 🎯

It's incredibly simple to use! 

**1. Basic Subdomain Enumeration**
```bash
./subfnder -d target.com
```

**2. Enumerate & Verify Live Hosts**
By passing the `-v` flag, Subfnder will perform DNS resolution to find out which subdomains are actively resolving.
```bash
./subfnder -d target.com -v
```

## Supported Sources 📡

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

## Contributing 🤝
Pull requests are welcome! If you find a new free API source, feel free to add it to the `sources/` directory.

## License 📄
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Happy Hunting! 🐞
