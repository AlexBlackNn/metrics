package main

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowProjectInfo(t *testing.T) {
	type args struct {
		log *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All empty",
			args: args{
				log: slog.New(slog.NewTextHandler(io.Discard, nil)),
			},
			want: "Build version: N/A, Build date: N/A, Build commit: N/A",
		},
		{
			name: "Version only",
			args: args{
				log: slog.New(slog.NewTextHandler(io.Discard, nil)),
			},
			want: "Build version: v1.0.0, Build date: N/A, Build commit: N/A",
		},
		{
			name: "Date only",
			args: args{
				log: slog.New(slog.NewTextHandler(io.Discard, nil)),
			},
			want: "Build version: N/A, Build date: 2023-12-21, Build commit: N/A",
		},
		{
			name: "Commit only",
			args: args{
				log: slog.New(slog.NewTextHandler(io.Discard, nil)),
			},
			want: "Build version: N/A, Build date: N/A, Build commit: abc1234",
		},
		{
			name: "All filled",
			args: args{
				log: slog.New(slog.NewTextHandler(io.Discard, nil)),
			},
			want: "Build version: v1.0.0, Build date: 2023-12-21, Build commit: abc1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the global variables before each test case
			buildVersion = ""
			buildDate = ""
			buildCommit = ""

			// Set the values based on the test case
			switch tt.name {
			case "Version only":
				buildVersion = "v1.0.0"
			case "Date only":
				buildDate = "2023-12-21"
			case "Commit only":
				buildCommit = "abc1234"
			case "All filled":
				buildVersion = "v1.0.0"
				buildDate = "2023-12-21"
				buildCommit = "abc1234"
			}

			// Capture the output of the log.Info call
			var logOutput bytes.Buffer
			tt.args.log = slog.New(slog.NewTextHandler(&logOutput, nil))

			showProjectInfo(tt.args.log)
			assert.Contains(t, logOutput.String(), tt.want)
		})
	}
}
