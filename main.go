package main

import (
    "fmt"
    "github.com/dotcloud/docker/archive"
    "io"
    "net/http"
    "os"
    "crypto/tls"
)

func usage() {
    fmt.Printf("Usage: %s https://example.org/download.tar.bz2\n", os.Args[0])
    os.Exit(1)
}

func download(url string) io.Reader {
    transport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
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
    if len(os.Args) < 2 {
        usage()
    }
    err := archive.Untar(download(os.Args[1]), "/", nil)
    if err != nil {
        fmt.Printf("Failed to extract -- %s\n", err)
    }
}
