package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/docker/docker/pkg/archive"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var VERSION = "0.4.0"

func usage() {
	example_url := "https://example.org/download.tar.bz2"
	fmt.Printf("Usage: %s %s\n", os.Args[0], example_url)
	fmt.Printf("Alternate usage: DUCKTAPE_URL=%s %s\n", example_url, os.Args[0])
	os.Exit(1)
}

func version() {
	fmt.Println(VERSION)
	os.Exit(0)
}

func get_path() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Failed to find current dir\n")
		os.exit(1)
	}
	return dir
}

func get_file(path string) string {
	return filepath.Join(get_path(), path)
}

func tls_config() *tls.Config {
	pool := x509.NewCertPool()
	cert, err := ioutil.ReadFile(get_file("cert"))
	if err != nil {
		fmt.Printf("Failed to load certificate -- %s\n", err)
		os.Exit(1)
	}
	pool.AppendCertsFromPEM(cert)
	return &tls.Config{RootCAs: pool}
}

func download(url string) io.Reader {
	transport := &http.Transport{TLSClientConfig: tls_config()}
	client := &http.Client{Transport: transport}
	response, err := client.Get(url)
	if err != nil {
		fmt.Printf("Failed to download %s -- %s\n", url, err)
		os.Exit(1)
	}
	fmt.Printf("Downloaded %s\n", url)
	return response.Body
}

func main() {
	url := os.Getenv("DUCKTAPE_URL")
	if len(os.Args) > 1 {
		if os.Args[1] == "-v" {
			version()
		}
		url = os.Args[1]
	}
	if len(url) == 0 {
		usage()
	}
	err := archive.Untar(download(url), "/", nil)
	if err != nil {
		fmt.Printf("Failed to extract -- %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully extracted archive\n")
}
