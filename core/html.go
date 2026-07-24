package core

import (
	"fmt"
	"os"
	"time"
)

// GenerateHTMLReport creates a responsive HTML dashboard for the results
func GenerateHTMLReport(filename string, subdomains []string, portResults map[string][]string, tkResults map[string]string, probeResults map[string]ProbeResult, waybackUrls []string, jsResults map[string][]string, techResults map[string][]string, vulnResults map[string][]string, corsResults map[string]string, wafResults map[string]string, paramResults map[string]string, silent bool) {
	if !silent {
		fmt.Printf(" %s[*] GENERATING HTML REPORT...%s\n", Gold, EndC)
	}

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TWHunt God-Tier Recon Report</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #0d1117; color: #c9d1d9; margin: 0; padding: 20px; }
        .container { max-width: 1500px; margin: 0 auto; }
        .header { text-align: center; padding: 40px 0; border-bottom: 2px solid #30363d; margin-bottom: 30px; }
        h1 { color: #58a6ff; margin: 0 0 10px 0; font-size: 2.5em; }
        .stats { display: flex; justify-content: center; gap: 20px; margin-bottom: 30px; flex-wrap: wrap; }
        .stat-box { background-color: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 20px; text-align: center; width: 180px; }
        .stat-value { font-size: 2em; color: #3fb950; font-weight: bold; }
        .stat-value.red { color: #f85149; }
        table { width: 100%; border-collapse: collapse; background-color: #161b22; border-radius: 8px; overflow: hidden; margin-bottom: 40px; }
        th, td { padding: 15px; text-align: left; border-bottom: 1px solid #30363d; }
        th { background-color: #21262d; color: #8b949e; font-weight: 600; }
        tr:hover { background-color: #30363d; }
        .subdomain { color: #58a6ff; font-weight: bold; text-decoration: none; }
        .subdomain:hover { text-decoration: underline; }
        .badge { background-color: #238636; color: white; padding: 3px 8px; border-radius: 12px; font-size: 0.85em; margin-right: 5px; display: inline-block; margin-bottom: 3px;}
        .badge-red { background-color: #f85149; }
        .badge-yellow { background-color: #d29922; color: #000; }
        .badge-blue { background-color: #1f6feb; }
        .footer { text-align: center; margin-top: 50px; color: #8b949e; font-size: 0.9em; }
        .urls-section { background-color: #161b22; padding: 20px; border-radius: 8px; max-height: 400px; overflow-y: auto; font-family: monospace; }
        .urls-section a { color: #8b949e; text-decoration: none; }
        .urls-section a:hover { color: #58a6ff; }
        .vuln-list { margin: 0; padding-left: 15px; color: #f85149; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>TWHunt Ultimate Report</h1>
            <p>God-Tier Passive & Active Subdomain Reconnaissance</p>
            <p>Generated on: ` + time.Now().Format(time.RFC1123) + `</p>
        </div>

        <div class="stats">
            <div class="stat-box">
                <div>Total Subdomains</div>
                <div class="stat-value">` + fmt.Sprintf("%d", len(subdomains)) + `</div>
            </div>
            <div class="stat-box">
                <div>Takeovers</div>
                <div class="stat-value red">` + fmt.Sprintf("%d", len(tkResults)) + `</div>
            </div>
            <div class="stat-box">
                <div>JS Secrets Leaked</div>
                <div class="stat-value red">` + fmt.Sprintf("%d", len(jsResults)) + `</div>
            </div>
            <div class="stat-box">
                <div>Vulns Exposed</div>
                <div class="stat-value red">` + fmt.Sprintf("%d", len(vulnResults)) + `</div>
            </div>
            <div class="stat-box">
                <div>Wayback URLs</div>
                <div class="stat-value">` + fmt.Sprintf("%d", len(waybackUrls)) + `</div>
            </div>
        </div>

        <h2>Subdomain Assets</h2>
        <table>
            <thead>
                <tr>
                    <th>#</th>
                    <th>Subdomain</th>
                    <th>Tech Stack</th>
                    <th>HTTP Status</th>
                    <th>Takeovers / Secrets / Vulns</th>
                </tr>
            </thead>
            <tbody>`

	for i, sub := range subdomains {
		html += `<tr><td>` + fmt.Sprintf("%d", i+1) + `</td>`
		
		subLink := `<td><a href="http://` + sub + `" target="_blank" class="subdomain">` + sub + `</a><br>`
		if ports, ok := portResults[sub]; ok && len(ports) > 0 {
			for _, p := range ports {
				subLink += `<span class="badge">Port: ` + p + `</span>`
			}
		}
		subLink += `</td>`
		html += subLink
		
		// Tech Stack
		techHtml := `<span style="color:#8b949e">-</span>`
		if tech, ok := techResults[sub]; ok && len(tech) > 0 {
			techHtml = ""
			for _, t := range tech {
				techHtml += `<span class="badge badge-blue">` + t + `</span>`
			}
		}
		html += `<td>` + techHtml + `</td>`

		// Probes (Status & Title)
		statusHtml := `<span style="color:#8b949e">-</span>`
		if pr, ok := probeResults[sub]; ok && pr.StatusCode > 0 {
			badgeClass := "badge"
			if pr.StatusCode >= 400 && pr.StatusCode < 500 {
				badgeClass = "badge badge-yellow"
			} else if pr.StatusCode >= 500 {
				badgeClass = "badge badge-red"
			}
			statusHtml = `<span class="` + badgeClass + `">` + fmt.Sprintf("%d", pr.StatusCode) + `</span><br><small>` + pr.Title + `</small>`
		}
		html += `<td>` + statusHtml + `</td>`

		// Critical findings column
		critHtml := ""
		if tk, ok := tkResults[sub]; ok {
			critHtml += `<span class="badge badge-red">TAKEOVER: ` + tk + `</span><br>`
		}
		if js, ok := jsResults[sub]; ok && len(js) > 0 {
			critHtml += `<span class="badge badge-red">SECRETS FOUND!</span>`
			critHtml += `<ul class="vuln-list">`
			for _, s := range js {
				critHtml += `<li>` + s + `</li>`
			}
			critHtml += `</ul>`
		}
		if v, ok := vulnResults[sub]; ok && len(v) > 0 {
			critHtml += `<span class="badge badge-red">EXPOSED FILES!</span>`
			critHtml += `<ul class="vuln-list">`
			for _, s := range v {
				critHtml += `<li>` + s + `</li>`
			}
			critHtml += `</ul>`
		}
		if c, ok := corsResults[sub]; ok {
			critHtml += `<span class="badge badge-red">CORS: ` + c + `</span><br>`
		}
		if w, ok := wafResults[sub]; ok {
			critHtml += `<span class="badge badge-yellow">WAF: ` + w + `</span><br>`
		}
		if p, ok := paramResults[sub]; ok {
			critHtml += `<span class="badge badge-blue">PARAMS: ` + p + `</span><br>`
		}
		
		if critHtml == "" {
			critHtml = `<span style="color:#8b949e">Safe</span>`
		}
		html += `<td>` + critHtml + `</td></tr>`
	}

	html += `
            </tbody>
        </table>`

	if len(waybackUrls) > 0 {
		html += `
		<h2>Wayback Machine Endpoints (Sampled)</h2>
		<div class="urls-section">`
		limit := len(waybackUrls)
		if limit > 1000 {
			limit = 1000 // Only show first 1000 in HTML to prevent massive lag
			html += `<p style="color:#d29922">Showing first 1000 URLs. See raw text output for full list.</p>`
		}
		for i := 0; i < limit; i++ {
			html += `<div><a href="` + waybackUrls[i] + `" target="_blank">` + waybackUrls[i] + `</a></div>`
		}
		html += `</div>`
	}

	html += `
        <div class="footer">
            Developed by MD TALHA HUSSAIN TAWHID | Open Source Intelligence
        </div>
    </div>
</body>
</html>`

	err := os.WriteFile(filename, []byte(html), 0644)
	if err != nil {
		if !silent {
			fmt.Printf(" %s[!] ERROR SAVING HTML REPORT: %v%s\n", Coral, err, EndC)
		}
		return
	}

	if !silent {
		fmt.Printf(" %s[✓] HTML REPORT SAVED TO: %s%s\n", Mint, filename, EndC)
	}
}
