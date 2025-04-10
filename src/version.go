package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/Masterminds/semver/v3"
)

type GistResponse struct {
	Files map[string]struct {
		RawURL string `json:"raw_url"`
	} `json:"files"`
}

type VersionApp struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

type Response struct {
	VersionApps []VersionApp `json:"versionApps"`
}

func VersionFetcher() string {
	gistID := ""
	gistAPI := fmt.Sprintf("https://api.github.com/gists/%s", gistID)

	githubToken := ""
	if githubToken == "" {
		fmt.Println("Error: GITHUB_TOKEN tidak ditemukan")
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", gistAPI, nil)
	if err != nil {
		fmt.Println("Error membuat request:", err)
		return ""
	}
	req.Header.Set("Authorization", "token "+githubToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching Gist:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Gagal fetch Gist dengan status", resp.Status)
		return ""
	}

	var gistData GistResponse
	if err := json.NewDecoder(resp.Body).Decode(&gistData); err != nil {
		fmt.Println("Error decoding Gist JSON:", err)
		return ""
	}

	fileData, ok := gistData.Files["ccraxt.json"]
	if !ok {
		fmt.Println("File ccraxt.json not found in Gist")
		return ""
	}

	req, err = http.NewRequest("GET", fileData.RawURL, nil)
	if err != nil {
		fmt.Println("Error membuat request ke file JSON:", err)
		return ""
	}
	req.Header.Set("Authorization", "token "+githubToken)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error fetching versions JSON:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Gagal fetch file JSON dengan status", resp.Status)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		fmt.Println("Raw JSON:", string(body))
		return ""
	}

	if len(data.VersionApps) == 0 {
		fmt.Println("No versions found")
		return ""
	}

	sort.Slice(data.VersionApps, func(i, j int) bool {
		v1, _ := semver.NewVersion(data.VersionApps[i].Version)
		v2, _ := semver.NewVersion(data.VersionApps[j].Version)
		return v1.LessThan(v2)
	})

	latest := data.VersionApps[len(data.VersionApps)-1]
	fmt.Println("Latest version:", latest.Version)

	return latest.URL
}
