package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v3"
)

var version = "0.5.0"

func usage() {
	exampleURL := "https://example.org/download.tar.bz2"
	fmt.Printf("Usage: %s %s\n", os.Args[0], exampleURL)
	fmt.Printf("Alternate usage: DUCKTAPE_URL=%s %s\n", exampleURL, os.Args[0])
}

func getDirPath() (string, error) {
	return os.Executable()
}

func getFilePath(name string) (string, error) {
	dirPath, err := getDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, name), nil
}

func getTmpFile() (string, error) {
	dirPath, err := getDirPath()
	if err != nil {
		return "", err
	}
	file, err := ioutil.TempFile(dirPath, "ducktape")
	if err != nil {
		return "", err
	}
	file.Close()
	return file.Name(), nil
}

func getTLSConfig() (tls.Config, error) {
	certFile, err := getFilePath("cert")
	if err != nil {
		return tls.Config{}, err
	}
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return tls.Config{}, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(cert)
	return tls.Config{RootCAs: pool}, nil
}

func getTLSClient() (http.Client, error) {
	tlsConfig, err := getTLSConfig()
	if err != nil {
		return http.Client{}, err
	}
	transport := http.Transport{TLSClientConfig: &tlsConfig}
	return http.Client{Transport: &transport}, nil
}

func download(path, url string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	client, err := getTLSClient()
	if err != nil {
		return err
	}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed http response: %s", response.Status)
	}

	_, err = io.Copy(file, response.Body)
	return err
}

func execute(url string) error {
	path, err := getTmpFile()
	if err != nil {
		return err
	}
	defer os.Remove(path)

	err = download(path, url)
	if err != nil {
		return err
	}

	return archiver.Unarchive(path, "/")
}

func main() {
	url := os.Getenv("DUCKTAPE_URL")
	if len(os.Args) > 1 {
		if os.Args[1] == "-v" {
			fmt.Println(version)
			return
		}
		url = os.Args[1]
	}
	if len(url) == 0 {
		usage()
		os.Exit(1)
	}

	err := execute(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
