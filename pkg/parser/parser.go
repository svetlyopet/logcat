package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reRepository     = regexp.MustCompile(`/([^/]+)/[^/]+/([^/]+)`)
	rePathTechnology = regexp.MustCompile(`/[^/]+/[^/]+/[^/]+/[^/]+/([^/].*)`)
	rePathGeneric    = regexp.MustCompile(`/[^/]+/([^/].*)`)
)

// Parse takes a log line containing data separated by a delimiter
// and returns a string with the parsed and modified data
func Parse(line string, delimiter string) (string, error) {
	split := strings.Split(line, delimiter)

	r := RequestLogs{
		timestamp: split[0],
		ip:        split[2],
		user:      split[3],
		method:    split[4],
		path:      split[5],
		status:    split[6],
		size:      split[8],
	}

	// add checks if the request contains information suitable for billing log
	if r.status != "200" || r.method != "GET" || r.size == "0" || r.user == "non_authenticated_user" || r.user == "anonymous" {
		return "", nil
	}

	// naming pattern suggested by Artifactory is to have all remote repositories have a suffix "-remote"
	repo := reRepository.FindStringSubmatch(r.path)
	if repo == nil || (!strings.Contains(repo[1], "-remote") && !strings.Contains(repo[2], "-remote")) {
		return "", nil
	}

	// start parsing the valuable information
	var repository string
	var artifactoryPath string

	// check if it's a generic remote repository
	if strings.Contains(repo[1], "-remote") {
		repository = repo[1]
		path := rePathGeneric.FindStringSubmatch(r.path)
		if path == nil {
			return "", nil
		}
		artifactoryPath = path[1]
	}

	// check if it`s a technology remote repository
	if strings.Contains(repo[2], "-remote") {
		repository = repo[2]
		path := rePathTechnology.FindStringSubmatch(r.path)
		// check if the requests is getting a token
		if path == nil || strings.Compare(path[1], "token") == 0 {
			return "", nil
		}
		artifactoryPath = strings.Replace(strings.Replace(path[1], ":", "__", 1), "/blobs/", "/", 1)
	}

	time := strings.Split(r.timestamp, "T")
	if len(time) != 2 {
		return "", fmt.Errorf("could not parse timestamp from request log")
	}

	splitTimestampTime := strings.Split(time[1], ":")
	if len(time) != 2 {
		return "", fmt.Errorf("could not parse timestamp from request log")
	}

	timestamp := time[0] + " " + splitTimestampTime[0] + ":00:00.000"

	quantity, err := strconv.ParseInt(r.size, 10, 64)
	if err != nil {
		return "", fmt.Errorf("cound not parse response size from request log: %v", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not get node hostname: %v", err)
	}

	billingLog := BillingLogs{
		Timestamp:       timestamp,
		ServerName:      hostname,
		Service:         "artifactory",
		Action:          "download",
		RemoteIP:        r.ip,
		Repository:      repository,
		Project:         "default",
		ArtifactoryPath: artifactoryPath,
		User:            r.user,
		ConsumptionUnit: "bytes",
		Quantity:        quantity,
	}

	logEntry, err := json.Marshal(billingLog)
	if err != nil {
		return "", fmt.Errorf("could not log billing entry: %v", err)
	}

	return string(logEntry), nil
}
