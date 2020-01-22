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

var VERSION = "0.5.0"

func usage() {
	example_url := "https://example.org/download.tar.bz2"
	fmt.Printf("Usage: %s %s\n", os.Args[0], example_url)
	fmt.Printf("Alternate usage: DUCKTAPE_URL=%s %s\n", example_url, os.Args[0])
}

func version() {
	fmt.Println(VERSION)
}

func get_dir_path() (string, error) {
	return os.Executable()
}

func get_file_path(name string) (string, error) {
	dir_path, err := get_dir_path()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir_path, name), nil
}

func get_tmp_file() (string, error) {
	dir_path, err := get_dir_path()
	if err != nil {
		return "", err
	}
	file, err := ioutil.TempFile(dir_path, "ducktape")
	if err != nil {
		return "", err
	}
	file.Close()
	return file.Name(), nil
}

func get_tls_config() (tls.Config, error) {
	cert_file, err := get_file_path("cert")
	if err != nil {
		return tls.Config{}, err
	}
	cert, err := ioutil.ReadFile(cert_file)
	if err != nil {
		return tls.Config{}, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(cert)
	return tls.Config{RootCAs: pool}, nil
}

func get_tls_client() (http.Client, error) {
	tls_config, err := get_tls_config()
	if err != nil {
		return http.Client{}, err
	}
	transport := http.Transport{TLSClientConfig: &tls_config}
	return http.Client{Transport: &transport}, nil
}

func download(path, url string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	client, err := get_tls_client()
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
	path, err := get_tmp_file()
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
			version()
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
