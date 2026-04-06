package main

import (
	"testing"
)

func TestParseTunnelArg_OK(t *testing.T) {
	tests := []struct {
		name        string
		arg         string
		defaultHost string
		wantProto   string
		wantLocal   string
		wantReqPort int
	}{
		{
			name:        "local_only_defaults_tcp_remote_auto",
			arg:         "80",
			defaultHost: "localhost",
			wantProto:   "tcp",
			wantLocal:   "localhost:80",
			wantReqPort: 0,
		},
		{
			name:        "local_remote_defaults_tcp",
			arg:         "80:10080",
			defaultHost: "127.0.0.1",
			wantProto:   "tcp",
			wantLocal:   "127.0.0.1:80",
			wantReqPort: 10080,
		},
		{
			name:        "tcp_local_remote",
			arg:         "tcp:22:10111",
			defaultHost: "localhost",
			wantProto:   "tcp",
			wantLocal:   "localhost:22",
			wantReqPort: 10111,
		},
		{
			name:        "udp_local_remote",
			arg:         "udp:53:1053",
			defaultHost: "localhost",
			wantProto:   "udp",
			wantLocal:   "localhost:53",
			wantReqPort: 1053,
		},
		{
			name:        "http_local_remote_auto",
			arg:         "http:3000",
			defaultHost: "localhost",
			wantProto:   "http",
			wantLocal:   "localhost:3000",
			wantReqPort: 0,
		},
		{
			name:        "trim_spaces",
			arg:         "  udp : 19132  ",
			defaultHost: " localhost ",
			wantProto:   "udp",
			wantLocal:   "localhost:19132",
			wantReqPort: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseTunnelArg(tt.arg, tt.defaultHost)
			if err != nil {
				t.Fatalf("expected ok, got err: %v", err)
			}
			if cfg.Protocol != tt.wantProto {
				t.Fatalf("Protocol: want %q, got %q", tt.wantProto, cfg.Protocol)
			}
			if cfg.LocalAddr != tt.wantLocal {
				t.Fatalf("LocalAddr: want %q, got %q", tt.wantLocal, cfg.LocalAddr)
			}
			if cfg.RequestedPort != tt.wantReqPort {
				t.Fatalf("RequestedPort: want %d, got %d", tt.wantReqPort, cfg.RequestedPort)
			}
		})
	}
}

func TestParseTunnelArg_Err(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{name: "empty", arg: ""},
		{name: "too_many_parts", arg: "tcp:80:10080:extra"},
		{name: "missing_part", arg: "tcp:80:"},
		{name: "non_numeric_port", arg: "udp:abc"},
		{name: "port_out_of_range", arg: "99999"},
		{name: "unknown_proto", arg: "foo:80"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseTunnelArg(tt.arg, "localhost")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

