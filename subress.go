package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	port       = flag.String("p", "443", "Port number")
	unique     = flag.String("d", "", "Print only results that ends with the provided domain name")
	quiet      = flag.Bool("q", false, "Don't print exceptions to screen")
	useThreads = flag.Bool("t", false, "Number of threads")
	wildcard   = flag.Bool("w", false, "Display wildcards")
	help       = flag.Bool("h", false, "Display help and usage")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Usage:\n")
		fmt.Println("Example: echo google.com | subress -d google.com -q -p 443")
		fmt.Println("Example: cat google_subdomains.txt | subress -d google.com -q -t 100")
		fmt.Println("Example: cat google_subdomains.txt | subress -q -t 100\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Read domain names from stdin
	domains := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	if *useThreads {
		var wg sync.WaitGroup
		for _, domain := range domains {
			wg.Add(1)
			go func(domain string) {
				defer wg.Done()
				getCertificateSubjectAlternativeName(domain)
			}(domain)
		}
		wg.Wait()
	} else {
		for _, domain := range domains {
			getCertificateSubjectAlternativeName(domain)
		}
	}
}

func getCertificateSubjectAlternativeName(domain string) {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", domain+":"+*port, nil)
	if err != nil {
		if !*quiet {
			fmt.Printf("Error connecting to %s: %s\n", domain, err)
		}
		return
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]

	printedNames := map[string]bool{}

	names := []string{}
	for _, san := range cert.DNSNames {
		if *unique != "" && !strings.HasSuffix(san, *unique) {
			continue
		}

		if (strings.HasPrefix(san, "*.")) && (!*wildcard) {
			san = san[2:]
		}
		names = append(names, san)
	}

	for _, name := range names {
		if !printedNames[name] {
			fmt.Println(name)
			printedNames[name] = true
		}
	}
}
