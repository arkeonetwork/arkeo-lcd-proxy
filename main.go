package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	logFilePath := getenv("LOG_FILE", "~/.lcd-proxy/lcd-proxy.log")
	logFile, err := setupLogging(logFilePath)
	if err != nil {
		log.Printf("log file disabled (%s): %v", logFilePath, err)
	}
	if logFile != nil {
		defer func() {
			_ = logFile.Close()
		}()
	}

	target := getenv("BACKEND_LCD_URL", "http://127.0.0.1:1317")
	listen := getenv("LISTEN", ":1318")

	u, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	orig := proxy.Director
	proxy.Director = func(r *http.Request) {
		orig(r)
		r.Host = u.Host
		if r.URL.Path == "/cosmos/tx/v1beta1/txs" {
			q := r.URL.Query()
			if q.Get("query") == "" {
				if evs, ok := q["events"]; ok && len(evs) > 0 {
					q.Set("query", strings.Join(evs, " AND "))
					q.Del("events")
					r.URL.RawQuery = q.Encode()
				}
			}
		}
	}

	srv := &http.Server{
		Addr:              listen,
		Handler:           proxy,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("proxy %s -> %s", listen, target)
	log.Fatal(srv.ListenAndServe())
}

func setupLogging(path string) (*os.File, error) {
	log.SetOutput(os.Stdout)
	if path == "" || path == "-" {
		return nil, nil
	}
	path = expandHome(path)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	return file, nil
}

func expandHome(path string) string {
	if path == "" || path == "-" {
		return path
	}
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
		return path
	}
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
