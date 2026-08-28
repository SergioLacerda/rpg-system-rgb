package main

import "testing"

// TestRunGoldenErrorText locks the exact, intentionally-public wording of
// the CLI's usage/error messages, per design.md's golden-layer guidance to
// snapshot only help/error text that is meant to stay stable. main_test.go
// covers dispatch behavior with looser Contains checks; this test is the
// stricter oracle that fails on any wording drift.
func TestRunGoldenErrorText(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "unknown top-level subcommand",
			args: []string{"bogus"},
			want: `unknown subcommand "bogus" (want validate|generate|bundle|release|docs)`,
		},
		{
			name: "missing docs subcommand",
			args: []string{"docs"},
			want: "missing docs subcommand (want library|pdf|skill|check)",
		},
		{
			name: "unknown docs subcommand",
			args: []string{"docs", "bogus"},
			want: `unknown docs subcommand "bogus" (want library|pdf|skill|check)`,
		},
		{
			name: "missing release subcommand",
			args: []string{"release"},
			want: "missing release subcommand (want manifest|check|skill-manifest|skill-check)",
		},
		{
			name: "unknown release subcommand",
			args: []string{"release", "bogus"},
			want: `unknown release subcommand "bogus" (want manifest|check|skill-manifest|skill-check)`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := run(tc.args)
			if err == nil {
				t.Fatalf("run(%v) returned no error", tc.args)
			}
			if got := err.Error(); got != tc.want {
				t.Fatalf("run(%v) error =\n  %q\nwant\n  %q", tc.args, got, tc.want)
			}
		})
	}
}
