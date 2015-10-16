package main

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/bsiegert/ranges"
)

var ErrParsePorts = fmt.Errorf("cannot parse ports argument")

// ParseVars returns parsed and validated slice of strings with
// variables names that will be used for monitoring.
func ParseVars(vars string) ([]VarName, error) {
	if vars == "" {
		return nil, errors.New("no vars specified")
	}

	ss := strings.FieldsFunc(vars, func(r rune) bool { return r == ',' })
	var ret []VarName
	for _, s := range ss {
		ret = append(ret, VarName(s))
	}
	return ret, nil
}

// BaseCommand returns cleaned command name from Cmdline array.
//
// I.e. "./some.service/binary.name -arg 1 -arg" will be "binary.name".
func BaseCommand(cmdline []string) string {
	if len(cmdline) == 0 {
		return ""
	}
	return filepath.Base(cmdline[0])
}

// flattenURLs returns URLs for the given addr and set of ports.
//
// Note, rawurl shouldn't contain port, as port will be appended.
func flattenURLs(rawurl string, ports []string) ([]url.URL, error) {
	var urls []url.URL

	// Add http by default
	if !strings.HasPrefix(rawurl, "http") {
		rawurl = fmt.Sprintf("http://%s", rawurl)
	}

	// Make URL from rawurl
	baseURL, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if baseURL.Path == "" {
		baseURL.Path = DefaultEndpoint
	}

	// Create new URL for each port
	for _, port := range ports {
		u := *baseURL
		u.Host = fmt.Sprintf("%s:%s", u.Host, port)
		urls = append(urls, u)
	}
	return urls, nil
}

// ParsePorts parses and flattens comma-separated ports/urls into URLs slice
func ParsePorts(s string) ([]url.URL, error) {
	var urls []url.URL
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ',' })
	for _, field := range fields {
		rawurl, portsRange := extractURLAndPorts(field)

		ports, err := parseRange(portsRange)
		if err != nil {
			return nil, ErrParsePorts
		}

		purls, err := flattenURLs(rawurl, ports)
		if err != nil {
			return nil, ErrParsePorts
		}

		urls = append(urls, purls...)
	}

	return urls, nil
}

// extractUrlAndPorts attempts to split url and extract raw url
// for the single port and range of ports to parse.
//
// i.e. "http://name:1234-1236/_endpoint" would return "http://name/_endpoint" and
// "1234-1236"
func extractURLAndPorts(s string) (string, string) {
	var rawurl, ports string
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 1:
		// "1234-234"
		rawurl = "http://localhost"
		ports = parts[0]
	case 2:
		// "localhost:1234"
		rawurl, ports = parts[0], parts[1]
	default:
		// "https://user:pass@remote.name:1234" or "http://name:1234-1236/_endpoint"

		// construct endpoint from the first part of URI, before ports appera
		rawurl = strings.Join(parts[:len(parts)-1], ":")

		// get either "1234-1235" or "1234-1235/_endpoint"
		lastPart := parts[len(parts)-1]

		// try to find endpoint and attach it to rawurl
		fields := strings.SplitN(lastPart, "/", 2)
		ports = fields[0]
		if len(fields) > 1 {
			rawurl = fmt.Sprintf("%s/%s", rawurl, fields[1])
		}
	}

	return rawurl, ports
}

// parseRange flattens port ranges, such as "1234-1240,1333"
func parseRange(s string) ([]string, error) {
	portsInt, err := ranges.Parse(s)
	if err != nil {
		return nil, err
	}

	var ports []string
	for _, port := range portsInt {
		ports = append(ports, fmt.Sprintf("%d", port))
	}
	return ports, nil
}

// NewURL returns net.URL for the given port, with expvarmon defaults set.
func NewURL(port string) url.URL {
	return url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port),
		Path:   "/debug/vars",
	}
}
