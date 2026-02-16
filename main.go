package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	host := flag.String("host", "https://idetrust.com", "host URL")
	daid := flag.String("daid", "QC-DEMO", "distributed authority identifier")
	cid := flag.Int("cid", 3, "certificate identifier")
	skipRoot := flag.Bool("skip-root", false, "skips the verification of the root node")
	flag.Parse()

	u, err := url.Parse(*host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	u = u.JoinPath("daid", *daid, "cid", strconv.Itoa(*cid))
	fmt.Println("URL:", u.String())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(u.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "unexpected status: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
		os.Exit(1)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/pem-certificate-chain") {
		fmt.Fprintf(os.Stderr, "unexpected content type: %s\n", contentType)
		os.Exit(1)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	certs, err := parseCertificateChain(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = verifyCertificateChain(certs, *skipRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	printCertificates(certs)
}
