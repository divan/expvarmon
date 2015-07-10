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
	baseUrl, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	baseUrl.Path = ExpvarsPath

	// Create new URL for each port
	for _, port := range ports {
		u := *baseUrl
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
		// Try simple 'ports range' mode, ports only ("1234-1235,80")
		// Defaults to "localhost" will be used.
		ports, err := parseRange(field)
		if err == nil {
			furls, err := flattenURLs("http://localhost", ports)
			if err != nil {
				return nil, err
			}
			urls = append(urls, furls...)
			continue
		}

		// then, try host:ports notation ("localhost:1234-1235,https://remote:2000,2345")
		var rawurl, portsRange string
		parts := strings.FieldsFunc(field, func(r rune) bool { return r == ':' })
		switch len(parts) {
		case 1:
			// "1234-234"
			rawurl = "http://localhost"
		case 2:
			// "localhost:1234"
			rawurl, portsRange = parts[0], parts[1]
		default:
			// "https://user:pass@remote.name:1234"
			rawurl = strings.Join(parts[:len(parts)-1], ":")
			portsRange = parts[len(parts)-1]
		}

		ports, err = parseRange(portsRange)
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
		Host:   fmt.Sprintf("localhost:%s"),
		Path:   "/debug/vars",
	}
}
