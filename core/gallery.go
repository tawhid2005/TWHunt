package core

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// GenerateDesktopGallery creates a standalone HTML gallery of the screenshots on the user's Desktop
func GenerateDesktopGallery(domain string, screenResults map[string]string, tkResults map[string]string, probeResults map[string]ProbeResult, vulnResults map[string][]string, corsResults map[string]string, wafResults map[string]string, paramResults map[string]string, silent bool) {
	if len(screenResults) == 0 {
		return // Don't generate if no screenshots
	}

	if !silent {
		fmt.Printf(" %s[*] GENERATING DESKTOP SCREENSHOT GALLERY...%s\n", Gold, EndC)
	}

	// 1. Find Desktop Path
	usr, err := user.Current()
	var desktopPath string
	if err == nil {
		// Works for Windows and standard Linux Desktop setups
		desktopPath = filepath.Join(usr.HomeDir, "Desktop")
	} else {
		// Fallback to current directory if Desktop isn't found
		desktopPath = "."
	}

	// Double check if Desktop exists, if not, use current directory
	if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
		desktopPath = "."
	}

	filename := filepath.Join(desktopPath, fmt.Sprintf("TWHunt_Gallery_%s.html", strings.ReplaceAll(domain, ".", "_")))

	// 2. Generate HTML
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TWHunt Screenshot Gallery - ` + domain + `</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #0d1117; color: #c9d1d9; margin: 0; padding: 20px; }
        h1 { text-align: center; color: #38b2ac; border-bottom: 2px solid #30363d; padding-bottom: 10px; }
        .gallery-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 20px; margin-top: 30px; }
        .card { background-color: #161b22; border: 1px solid #30363d; border-radius: 8px; overflow: hidden; transition: transform 0.2s; }
        .card:hover { transform: translateY(-5px); box-shadow: 0 10px 20px rgba(0,0,0,0.5); border-color: #38b2ac; }
        .card-img-container { width: 100%; height: 250px; overflow: hidden; background: #000; display: flex; align-items: center; justify-content: center; }
        .card-img { width: 100%; height: auto; transition: transform 0.3s; cursor: pointer; }
        .card-img:hover { transform: scale(1.05); }
        .card-body { padding: 15px; }
        .card-title { margin: 0 0 10px 0; font-size: 1.1em; font-weight: bold; }
        .card-title a { color: #58a6ff; text-decoration: none; word-break: break-all; }
        .card-title a:hover { text-decoration: underline; }
        .badge { display: inline-block; padding: 3px 8px; border-radius: 12px; font-size: 0.75em; font-weight: bold; margin: 2px 2px 2px 0; }
        .badge-red { background-color: #f85149; color: #fff; }
        .badge-yellow { background-color: #d29922; color: #000; }
        .badge-blue { background-color: #1f6feb; color: #fff; }
        .badge-purple { background-color: #8e44ad; color: #fff; }
        .badge-green { background-color: #238636; color: #fff; }
        .footer { text-align: center; margin-top: 50px; color: #8b949e; font-size: 0.9em; }
        .vuln-list { margin: 5px 0 0 0; padding-left: 20px; font-size: 0.85em; color: #ff7b72; }
    </style>
</head>
<body>
    <h1>📸 TWHunt Screenshot Gallery for ` + domain + `</h1>
    <div style="text-align: center; margin-bottom: 20px;">
        <span class="badge badge-blue">Total Screenshots: ` + fmt.Sprintf("%d", len(screenResults)) + `</span>
        <span class="badge badge-purple">Generated: ` + time.Now().Format("2006-01-02 15:04:05") + `</span>
    </div>
    
    <div class="gallery-grid">
`

	// 3. Add each screenshot card
	for sub, imgPath := range screenResults {
		// Need absolute path for the image so it loads from the Desktop HTML
		absImgPath, _ := filepath.Abs(imgPath)
		// For Windows paths in file:// URI
		fileUri := "file:///" + strings.ReplaceAll(absImgPath, "\\", "/")
		
		targetUrl := "https://" + sub
		if probe, ok := probeResults[sub]; ok && probe.StatusCode != 443 {
			// This is a rough estimation since ProbeResult doesn't store the scheme directly, 
			// but we know it probed successfully
		}

		cardHtml := `
        <div class="card">
            <div class="card-img-container">
                <a href="` + targetUrl + `" target="_blank" title="Open in browser">
                    <img src="` + fileUri + `" class="card-img" alt="` + sub + `">
                </a>
            </div>
            <div class="card-body">
                <h3 class="card-title"><a href="` + targetUrl + `" target="_blank">` + sub + `</a></h3>
                <div class="badges">`

		// Add logic for badges (Bugs)
		hasBugs := false
		
		if probe, ok := probeResults[sub]; ok {
			cardHtml += `<span class="badge badge-green">` + fmt.Sprintf("%d", probe.StatusCode) + `</span>`
		}
		
		if tk, ok := tkResults[sub]; ok {
			cardHtml += `<span class="badge badge-red">TAKEOVER: ` + tk + `</span>`
			hasBugs = true
		}
		
		if w, ok := wafResults[sub]; ok {
			cardHtml += `<span class="badge badge-yellow">WAF: ` + w + `</span>`
		}
		
		if p, ok := paramResults[sub]; ok {
			cardHtml += `<span class="badge badge-blue">PARAM: ` + p + `</span>`
			hasBugs = true
		}
		
		if c, ok := corsResults[sub]; ok {
			cardHtml += `<span class="badge badge-red">CORS: ` + c + `</span>`
			hasBugs = true
		}

		if v, ok := vulnResults[sub]; ok && len(v) > 0 {
			cardHtml += `<span class="badge badge-red">EXPOSED FILES</span>`
			cardHtml += `<ul class="vuln-list">`
			for _, s := range v {
				cardHtml += `<li>` + s + `</li>`
			}
			cardHtml += `</ul>`
			hasBugs = true
		}
		
		if !hasBugs {
			cardHtml += `<span class="badge badge-green" style="background:#2ea043; color:white;">Clean</span>`
		}

		cardHtml += `
                </div>
            </div>
        </div>`

		html += cardHtml
	}

	html += `
    </div>
    <div class="footer">
        Generated by <a href="https://github.com/tawhid2005/TWHunt" style="color: #58a6ff;">TWHunt</a> - The Ultimate Bug Bounty Framework
    </div>
</body>
</html>`

	// 4. Write to file
	err = os.WriteFile(filename, []byte(html), 0644)
	if err == nil {
		if !silent {
			fmt.Printf(" %s[!] DESKTOP GALLERY SAVED AT: %s%s\n", Coral, filename, EndC)
		}
	} else {
		if !silent {
			fmt.Printf(" %s[-] Failed to save Desktop Gallery: %v%s\n", Silver, err, EndC)
		}
	}
}
