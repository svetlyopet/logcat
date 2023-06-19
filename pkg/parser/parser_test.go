package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		delimiter  string
		numFields  int
		serverName string
		wantResult string
		wantError  bool
	}{
		{
			name:       "ValidLogEntry",
			line:       "2023-06-15T12:34:56.789Z|abcdefgh12345678|1.2.3.4|user|GET|/api/docker/registry-docker-remote/v2/alpine/curl/manifests/latest|200|-1|1234|567|user-agent123",
			delimiter:  "|",
			numFields:  11,
			serverName: "artifactory.domain",
			wantResult: `{"billing_timestamp":"2023-06-15 12:00:00.000","server_name":"artifactory.domain","service":"artifactory","action":"download","ip":"1.2.3.4","repository":"registry-docker-remote","project":"default","artifactory_path":"alpine/curl/manifests/latest","user_name":"user","consumption_unit":"bytes","quantity":1234}`,
			wantError:  false,
		},
		{
			name:       "InvalidLogEntryNumFields",
			line:       "2023-06-15T12:34:56.789Z|GET|127.0.0.1|user|GET|/api/docker/registry-docker-remote/v2/alpine/curl/manifests/latest|200|response|12345",
			delimiter:  "|",
			numFields:  11,
			serverName: "artifactory.domain",
			wantResult: "",
			wantError:  true,
		},
		{
			name:       "InvalidLogEntryDiffDelimiter",
			line:       "2023-06-15T12:34:56.789Z|abcdefgh12345678|1.2.3.4|user|GET|/api/docker/registry-docker-remote/v2/alpine/curl/manifests/latest|200|-1|1234|567|user-agent123",
			delimiter:  ";",
			numFields:  11,
			serverName: "artifactory.domain",
			wantResult: "",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotError := Parse(tt.line, tt.delimiter, tt.numFields, tt.serverName)

			if tt.wantError {
				if gotError == nil {
					t.Errorf("Parse() error = %v, wantErr %v", gotError, tt.wantError)
				}
				return
			}

			if gotResult != tt.wantResult {
				t.Errorf("Parse() result = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
