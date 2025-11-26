package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	cache      map[string]string
	cacheExp   time.Time
	cacheMutex sync.Mutex
	visits     uint64
)

const counterFile = "/data/counter.txt"

func main() {
	cache = make(map[string]string)
	loadCounter()
	http.HandleFunc("/", handleIndex)
	// Winbox 4
	http.HandleFunc("/winbox4/windows", wrap("winbox4_windows"))
	http.HandleFunc("/winbox4/mac", wrap("winbox4_mac"))
	http.HandleFunc("/winbox4/linux", wrap("winbox4_linux"))
	// Winbox 3
	http.HandleFunc("/winbox3/windows", wrap("winbox3_windows"))
	http.HandleFunc("/winbox3/windows32", wrap("winbox3_windows_32"))
	// visitcounter
	http.HandleFunc("/counter", counterHandler)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	n := atomic.AddUint64(&visits, 1)
	saveCounter()

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("%d", n)))
}

func loadCounter() {
	data, err := os.ReadFile(counterFile)
	if err != nil {
		return // no file yet, starts at 0
	}

	n, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err == nil {
		atomic.StoreUint64(&visits, n)
	}
}
func saveCounter() {
	n := atomic.LoadUint64(&visits)
	os.WriteFile(counterFile, []byte(fmt.Sprintf("%d", n)), 0644)
}
func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// preload cache for version extraction
	getURL("winbox4_windows")

	html := `
<!DOCTYPE html>
<html>
<head>
  <title>WinboxGet</title>
  <meta charset="UTF-8">
</head>
<style>
  body {
    font-family: system-ui, sans-serif;
    background: #fafafa;
    padding: 2rem;
    color: #333;
    max-width: 600px;
    margin: auto;
  }
  h1 {
    text-align: center;
    color: #222;
    margin-bottom: 1.5rem;
  }
  ul { list-style: none; padding: 0; margin: 0; }
  ul li { margin: 0.6rem 0; }
  a {
    display: block;
    padding: 0.8rem 1rem;
    background: #fff;
    border: 1px solid #ddd;
    border-radius: 8px;
    text-decoration: none;
    color: #0074d9;
    transition: all 0.15s ease-in-out;
  }
  a:hover { background: #e9f3ff; border-color: #bcdcff; }
</style>
<h1>Winboxget</h1>
<ul>
  <li><a href="/winbox4/windows">Winbox 4 (Windows) - v%v</a></li>
  <li><a href="/winbox4/mac">Winbox 4 (Mac) - v%v</a></li>
  <li><a href="/winbox4/linux">Winbox 4 (Linux) - v%v</a></li>
  <li><a href="/winbox3/windows">Winbox 3 x64 (Windows) - v%v</a></li>
  <li><a href="/winbox3/windows32">Winbox 3 32-bit (Windows) - v%v</a></li>
</ul>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(html,
		extractVersion(cache["winbox4_windows"]),
		extractVersion(cache["winbox4_mac"]),
		extractVersion(cache["winbox4_linux"]),
		extractVersion(cache["winbox3_windows"]),
		extractVersion(cache["winbox3_windows_32"]),
	)))
	w.Write([]byte(`
<div style="font-size: 0.8rem; color: #555; margin-top: 2rem;">
  <hr>
  <p><strong>Legendary Legal Disclaimer</strong></p>
  <p>
    This site is an unofficial convenience tool for downloading Winbox.  
    It is <strong>not</strong> created by, endorsed by, sponsored by, or in any way
    affiliated with <strong>SIA Mikrotīkls</strong>. They probably don't even know this page exists,
    and if they do, <i> hello there</i>
  </p>
  <p>
    Winbox, MikroTik, and any related names or logos are trademarks or registered
    trademarks of <strong>SIA Mikrotīkls</strong>. All rights belong to their respective owners.
  </p>
  <p>
    This site comes with absolutely <strong>no warranty</strong> of any kind whatsoever.  
    No warranty of correctness, uptime, usefulness, fitness for any purpose,
    or resilience against solar flares, router goblins, or spontaneous combustion.
  </p>
  <p>
    By using this site, you agree that:
    <ul>
    <li> You receive <strong>no license</strong> to any MikroTik software from here.</li>  
    <li> All downloads remain subject to MikroTik's original terms.</li>  
    <li> If something breaks, crashes, explodes, or starts speaking in tongues, that is 100% on you.</li>  
    <li> <strong>nugge cannot be blamed for lost limbs, missing fingers, damaged routers,  
      corrupted configs, emotional distress, or unexpected portal openings.</strong></li>
    </ul>
  </p>
  <p>
    For official software, support, the best routers and the real deal, please visit:<br>
    <a href="https://mikrotik.com" target="_blank" rel="noopener noreferrer">
      https://mikrotik.com
    </a>
  </p>
<p><i>P.S. The official WinBox download is a 4-click adventure. 5 including your click to mikrotik.com.
6 if you miss-click. (You will.)</i></p>
</div>
<p id="counter" style="text-align:center;color:#777;font-size:0.9rem;">
  Loading visitors…
</p>

<script>
fetch("/counter")
  .then(r => r.text())
  .then(n => {
    document.getElementById("counter").textContent =
      "You are visitor number:  " + n;
  })
  .catch(() => {
    document.getElementById("counter").textContent =
      "Visitors: unavailable";
  });
</script>
<p style="
    text-align:center;
    margin-top:2rem;
    font-size:0.85rem;
    color:#888;
    letter-spacing:0.5px;
">
&#169; <span style="font-weight:600;">nugge boman (nugge (a) nugge.fi)</span> - All vibes reserved.
</p>
<p style="text-align:center; margin-top:1.5rem;">
  <a href="https://github.com/nuggex/winboxget"
     target="_blank"
     rel="noopener noreferrer"
     style="text-decoration:none; font-size:0.9rem;">
    <span style="font-size:1.2rem; vertical-align:middle;">🐙</span>
    <span style="vertical-align:middle;">View on GitHub</span>
  </a>
</p>
</html>
	`))
}

func wrap(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := getURL(key)
		if err != nil {
			http.Error(w, "Failed to fetch Winbox"+err.Error(), 500)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func extractVersion(url string) string {
	parts := strings.Split(url, "/")
	return parts[5]
}

func getURL(key string) (string, error) {
	cacheMutex.Lock()
	if time.Now().Before(cacheExp) && cache[key] != "" {
		val := cache[key]
		cacheMutex.Unlock()
		return val, nil
	}
	cacheMutex.Unlock()

	urls, err := scrape()
	if err != nil {
		return "", err
	}

	cacheMutex.Lock()
	for k, v := range urls {
		cache[k] = v
	}
	cacheExp = time.Now().Add(1 * time.Hour)
	cacheMutex.Unlock()

	return cache[key], nil
}

func scrape() (map[string]string, error) {
	req, _ := http.NewRequest("GET", "https://mikrotik.com/download/winbox", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	results := make(map[string]string)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if href == "" {
			return
		}

		l := strings.ToLower(href)

		if strings.Contains(l, "winbox") || strings.Contains(l, "WinBox") {
			full := href
			switch {
			case strings.Contains(l, "winbox_windows"):
				results["winbox4_windows"] = full
			case strings.Contains(l, "winbox") && strings.Contains(l, ".dmg"):
				results["winbox4_mac"] = full
			case strings.Contains(l, "winbox_linux"):
				results["winbox4_linux"] = full
			case strings.Contains(l, "winbox64") && strings.Contains(l, ".exe"):
				results["winbox3_windows"] = full
			case strings.Contains(l, "winbox.exe"):
				results["winbox3_windows_32"] = full
			}
		}
	})

	return results, nil
}
