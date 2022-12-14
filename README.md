# subress
Extract subdomains from  SSL certificates 

This code was generate by ChatGPT :)

## Install

```go install -v github.com/NightRang3r/subress@main```

## Build

```go build subress.go```

## Usage

```
Usage:

Example: echo google.com | subress -d google.com -q -p 443
Example: cat google_subdomains.txt | subress -d google.com -q -t 100
Example: cat google_subdomains.txt | subress -q -t 100

  -d string
    	Print only results that ends with the provided domain name
  -h	Display help and usage
  -p string
    	Port number (default "443")
  -q	Don't print exceptions to screen
  -t	Number of threads
  -w	Display wildcards
  
 ```
