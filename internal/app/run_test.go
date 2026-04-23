package app

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/jbradley/dns-discovery/internal/discovery"
)

type noopLogger struct{}

func (noopLogger) Infof(string, ...any) {}
func (noopLogger) Errorf(string, ...any) {}
func (noopLogger) Close() error { return nil }

func withRunDiscoveryHooks(t *testing.T) {
	t.Helper()

	originalScanner := testHookScanner
	originalWriter := testHookReportWriter
	originalLoggerFactory := testHookLoggerFactory

	t.Cleanup(func() {
		testHookScanner = originalScanner
		testHookReportWriter = originalWriter
		testHookLoggerFactory = originalLoggerFactory
	})
}

func TestValidateOutputFormatAcceptsSupportedValues(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  OutputFormat
	}{
		{name: "markdown", input: "markdown", want: OutputMarkdown},
		{name: "json", input: " JSON ", want: OutputJSON},
		{name: "text", input: "text", want: OutputText},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ValidateOutputFormat(tc.input)
			if err != nil {
				t.Fatalf("validate output format: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestValidateOutputFormatRejectsInvalidValues(t *testing.T) {
	if _, err := ValidateOutputFormat("xml"); err == nil {
		t.Fatal("expected error for unsupported output format")
	}
}

func TestRunDiscoveryCollectsOrderedResults(t *testing.T) {
	withRunDiscoveryHooks(t)

	testHookLoggerFactory = func(logLocation string, verbose bool) (runLoggerSink, error) {
		return noopLogger{}, nil
	}
	testHookScanner = func(domain string) (*discovery.DiscoveryResult, error) {
		if domain == "bad.example" {
			return nil, errors.New("lookup failed")
		}
		return &discovery.DiscoveryResult{Domain: domain}, nil
	}
	testHookReportWriter = func(baseDir string, res *discovery.DiscoveryResult, output OutputFormat) (string, error) {
		return filepath.Join(baseDir, res.Domain, "report.md"), nil
	}

	summary, err := RunDiscovery(
		[]string{"good.example", "bad.example", "great.example"},
		RunOptions{OutputDir: "output", Output: OutputMarkdown, LogLocation: "/tmp/dns-discovery.log"},
	)
	if err != nil {
		t.Fatalf("run discovery: %v", err)
	}

	if len(summary.Succeeded) != 2 {
		t.Fatalf("expected 2 successes, got %d", len(summary.Succeeded))
	}
	if summary.Succeeded[0].Domain != "good.example" || summary.Succeeded[1].Domain != "great.example" {
		t.Fatalf("success order mismatch: %#v", summary.Succeeded)
	}
	if len(summary.Failed) != 1 {
		t.Fatalf("expected 1 failure, got %d", len(summary.Failed))
	}
	if summary.Failed[0].Domain != "bad.example" || summary.Failed[0].Err == nil || summary.Failed[0].Err.Error() != "lookup failed" {
		t.Fatalf("unexpected failure entry: %#v", summary.Failed[0])
	}
}

func TestRunDiscoverySkipsBlankDomains(t *testing.T) {
	withRunDiscoveryHooks(t)

	calls := 0
	testHookLoggerFactory = func(logLocation string, verbose bool) (runLoggerSink, error) {
		return noopLogger{}, nil
	}
	testHookScanner = func(domain string) (*discovery.DiscoveryResult, error) {
		calls++
		return &discovery.DiscoveryResult{Domain: domain}, nil
	}
	testHookReportWriter = func(baseDir string, res *discovery.DiscoveryResult, output OutputFormat) (string, error) {
		return filepath.Join(baseDir, res.Domain, "report.md"), nil
	}

	summary, err := RunDiscovery(
		[]string{" github.com ", "", "   ", "cloudflare.com"},
		RunOptions{OutputDir: "output", Output: OutputMarkdown, LogLocation: "/tmp/dns-discovery.log"},
	)
	if err != nil {
		t.Fatalf("run discovery: %v", err)
	}

	if calls != 2 {
		t.Fatalf("expected 2 scanner calls, got %d", calls)
	}
	if summary.Total() != 2 {
		t.Fatalf("expected total 2, got %d", summary.Total())
	}
}

func TestRunDiscoveryCapturesWriterErrorsAsFailures(t *testing.T) {
	withRunDiscoveryHooks(t)

	testHookLoggerFactory = func(logLocation string, verbose bool) (runLoggerSink, error) {
		return noopLogger{}, nil
	}
	testHookScanner = func(domain string) (*discovery.DiscoveryResult, error) {
		return &discovery.DiscoveryResult{Domain: domain}, nil
	}
	testHookReportWriter = func(baseDir string, res *discovery.DiscoveryResult, output OutputFormat) (string, error) {
		if res.Domain == "bad.example" {
			return "", errors.New("write failed")
		}
		return filepath.Join(baseDir, res.Domain, "report.md"), nil
	}

	summary, err := RunDiscovery(
		[]string{"good.example", "bad.example"},
		RunOptions{OutputDir: "output", Output: OutputMarkdown, LogLocation: "/tmp/dns-discovery.log"},
	)
	if err != nil {
		t.Fatalf("run discovery: %v", err)
	}

	if len(summary.Succeeded) != 1 || summary.Succeeded[0].Domain != "good.example" {
		t.Fatalf("unexpected successes: %#v", summary.Succeeded)
	}
	if len(summary.Failed) != 1 || summary.Failed[0].Domain != "bad.example" {
		t.Fatalf("unexpected failures: %#v", summary.Failed)
	}
}
