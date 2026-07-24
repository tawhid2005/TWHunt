package utils

import (
	"fmt"
	"os"
	"time"
	"twhunt/core"
)

// GenerateHTMLReport creates a responsive HTML dashboard for the results
func GenerateHTMLReport(filename string, subdomains []string, portResults map[string][]string, silent bool) {
	if !silent {
		fmt.Printf(" %s[*] GENERATING HTML REPORT...%s\n", core.Gold, core.EndC)
	}

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TWHunt Reconnaissance Report</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #0d1117; color: #c9d1d9; margin: 0; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { text-align: center; padding: 40px 0; border-bottom: 2px solid #30363d; margin-bottom: 30px; }
        h1 { color: #58a6ff; margin: 0 0 10px 0; font-size: 2.5em; }
        .stats { display: flex; justify-content: center; gap: 30px; margin-bottom: 30px; }
        .stat-box { background-color: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 20px; text-align: center; width: 200px; }
        .stat-value { font-size: 2em; color: #3fb950; font-weight: bold; }
        table { width: 100%; border-collapse: collapse; background-color: #161b22; border-radius: 8px; overflow: hidden; }
        th, td { padding: 15px; text-align: left; border-bottom: 1px solid #30363d; }
        th { background-color: #21262d; color: #8b949e; font-weight: 600; }
        tr:hover { background-color: #30363d; }
        .subdomain { color: #58a6ff; font-weight: bold; text-decoration: none; }
        .subdomain:hover { text-decoration: underline; }
        .badge { background-color: #238636; color: white; padding: 3px 8px; border-radius: 12px; font-size: 0.8em; margin-right: 5px; }
        .footer { text-align: center; margin-top: 50px; color: #8b949e; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>TWHunt Pro Report</h1>
            <p>Advanced Passive & Active Subdomain Enumeration</p>
            <p>Generated on: ` + time.Now().Format(time.RFC1123) + `</p>
        </div>

        <div class="stats">
            <div class="stat-box">
                <div>Total Discovered</div>
                <div class="stat-value">` + fmt.Sprintf("%d", len(subdomains)) + `</div>
            </div>
            <div class="stat-box">
                <div>With Open Ports</div>
                <div class="stat-value">` + fmt.Sprintf("%d", len(portResults)) + `</div>
            </div>
        </div>

        <table>
            <thead>
                <tr>
                    <th>#</th>
                    <th>Subdomain</th>
                    <th>Open Ports</th>
                </tr>
            </thead>
            <tbody>`

	for i, sub := range subdomains {
		html += `<tr><td>` + fmt.Sprintf("%d", i+1) + `</td>`
		html += `<td><a href="http://` + sub + `" target="_blank" class="subdomain">` + sub + `</a></td>`
		
		portsHtml := ""
		if ports, ok := portResults[sub]; ok && len(ports) > 0 {
			for _, p := range ports {
				portsHtml += `<span class="badge">` + p + `</span>`
			}
		} else {
			portsHtml = `<span style="color:#8b949e">-</span>`
		}
		
		html += `<td>` + portsHtml + `</td></tr>`
	}

	html += `
            </tbody>
        </table>
        
        <div class="footer">
            Developed by MD TALHA HUSSAIN TAWHID | Open Source Intelligence
        </div>
    </div>
</body>
</html>`

	err := os.WriteFile(filename, []byte(html), 0644)
	if err != nil {
		if !silent {
			fmt.Printf(" %s[!] ERROR SAVING HTML REPORT: %v%s\n", core.Coral, err, core.EndC)
		}
		return
	}

	if !silent {
		fmt.Printf(" %s[✓] HTML REPORT SAVED TO: %s%s\n", core.Mint, filename, core.EndC)
	}
}
