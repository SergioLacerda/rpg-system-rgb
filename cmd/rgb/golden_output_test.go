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
			want: "missing release subcommand (want manifest|check|pdf-regression-check|skill-manifest|skill-check)",
		},
		{
			name: "unknown release subcommand",
			args: []string{"release", "bogus"},
			want: `unknown release subcommand "bogus" (want manifest|check|pdf-regression-check|skill-manifest|skill-check)`,
		},
		{
			name: "unknown flag on docs library",
			args: []string{"docs", "library", "--bogus"},
			want: "flag provided but not defined: -bogus",
		},
		{
			name: "unknown flag on docs pdf",
			args: []string{"docs", "pdf", "--bogus"},
			want: "flag provided but not defined: -bogus",
		},
		{
			name: "unknown flag on docs skill",
			args: []string{"docs", "skill", "--bogus"},
			want: "flag provided but not defined: -bogus",
		},
		{
			name: "unknown flag on docs check",
			args: []string{"docs", "check", "--bogus"},
			want: "flag provided but not defined: -bogus",
		},
		{
			name: "unknown flag on release manifest",
			args: []string{"release", "manifest", "--bogus"},
			want: "flag provided but not defined: -bogus",
		},
		{
			name: "unknown flag on release check",
			args: []string{"release", "check", "--bogus"},
			want: "flag provided but not defined: -bogus",
		},
		{
			name: "unknown flag on release pdf-regression-check",
			args: []string{"release", "pdf-regression-check", "--bogus"},
			want: "flag provided but not defined: -bogus",
		},
		{
			name: "unknown flag on release skill-manifest",
			args: []string{"release", "skill-manifest", "--bogus"},
			want: "flag provided but not defined: -bogus",
		},
		{
			name: "unknown flag on release skill-check",
			args: []string{"release", "skill-check", "--bogus"},
			want: "flag provided but not defined: -bogus",
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
