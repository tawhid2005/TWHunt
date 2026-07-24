package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"twhunt/core"
)

func AutoUpdate() {
	fmt.Printf("\n%s[*] Checking for updates...%s\n", core.Sky, core.EndC)
	
	// Determine the binary name to download based on OS
	var downloadURL string
	var backupName string
	
	if runtime.GOOS == "windows" {
		downloadURL = "https://github.com/tawhid2005/TWHunt/raw/master/twhunt.exe"
		backupName = "twhunt.old.exe"
	} else if runtime.GOOS == "linux" {
		downloadURL = "https://github.com/tawhid2005/TWHunt/raw/master/twhunt_linux_amd64"
		backupName = "twhunt.old"
	} else {
		fmt.Printf("%s[!] Auto-update is only supported for Windows and Linux.%s\n", core.Coral, core.EndC)
		return
	}

	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("%s[!] Error determining executable path: %v%s\n", core.Coral, err, core.EndC)
		return
	}

	// Rename current executable to backup
	backupPath := filepath.Join(filepath.Dir(execPath), backupName)
	os.Remove(backupPath) // Remove old backup if exists
	
	err = os.Rename(execPath, backupPath)
	if err != nil {
		fmt.Printf("%s[!] Error renaming current executable (permission denied?): %v%s\n", core.Coral, err, core.EndC)
		return
	}

	// Download new binary
	resp, err := http.Get(downloadURL)
	if err != nil {
		os.Rename(backupPath, execPath) // Restore
		fmt.Printf("%s[!] Error downloading update: %v%s\n", core.Coral, err, core.EndC)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Rename(backupPath, execPath) // Restore
		fmt.Printf("%s[!] Update not found or error server side.%s\n", core.Coral, core.EndC)
		return
	}

	out, err := os.OpenFile(execPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		os.Rename(backupPath, execPath) // Restore
		fmt.Printf("%s[!] Error creating new executable: %v%s\n", core.Coral, err, core.EndC)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Rename(backupPath, execPath) // Restore
		fmt.Printf("%s[!] Error writing new executable: %v%s\n", core.Coral, err, core.EndC)
		return
	}

	// Give execution permissions on Linux/Mac
	if runtime.GOOS != "windows" {
		cmd := exec.Command("chmod", "+x", execPath)
		cmd.Run()
	}

	fmt.Printf("%s[✓] TWHunt has been successfully updated to the latest version!%s\n", core.Mint, core.EndC)
	fmt.Printf("%s[*] Please restart the tool.%s\n", core.Gold, core.EndC)
	os.Exit(0)
}
