package app

import (
	"errors"
	"testing"
)

// Verifies that RunDomain delegates to the injected runner seam instead of
// directly calling executeDomain. It also asserts that domain/output arguments
// are forwarded unchanged and that a successful runner returns nil.
func TestRunDomainDelegatesToRunner(t *testing.T) {
	original := testHookDomainRunner
	defer func() {
		testHookDomainRunner = original
	}()

	called := false
	testHookDomainRunner = func(domain string, outputDir string) error {
		called = true
		if domain != "github.com" {
			t.Fatalf("unexpected domain: %q", domain)
		}
		if outputDir != "reports" {
			t.Fatalf("unexpected output dir: %q", outputDir)
		}
		return nil
	}

	if err := RunDomain("github.com", "reports"); err != nil {
		t.Fatalf("run domain: %v", err)
	}
	if !called {
		t.Fatal("expected domain runner to be called")
	}
}

// Verifies batch aggregation behavior when domains have mixed outcomes.
// Successful domains are collected in Succeeded, while runner errors are
// recorded in Failed keyed by normalized domain name.
func TestRunBatchCollectsMixedResults(t *testing.T) {
	summary := runBatch([]string{"good.example", "bad.example"}, "output", func(domain string, outputDir string) error {
		if domain == "bad.example" {
			return errors.New("lookup failed")
		}
		return nil
	})

	if len(summary.Succeeded) != 1 || summary.Succeeded[0] != "good.example" {
		t.Fatalf("unexpected successes: %#v", summary.Succeeded)
	}
	if len(summary.Failed) != 1 {
		t.Fatalf("unexpected failures: %#v", summary.Failed)
	}
	if err := summary.Failed["bad.example"]; err == nil || err.Error() != "lookup failed" {
		t.Fatalf("unexpected batch failure for bad.example: %v", err)
	}
}

// Verifies that runBatch trims/normalizes input domains and skips blank entries.
// The runner should be invoked only for non-empty domains, and the summary total
// should match the number of processed domains.
func TestRunBatchSkipsBlankDomains(t *testing.T) {
	called := 0
	summary := runBatch([]string{" github.com ", "", "   ", "cloudflare.com"}, "output", func(domain string, outputDir string) error {
		called++
		return nil
	})

	if called != 2 {
		t.Fatalf("expected 2 runner calls, got %d", called)
	}
	if summary.Total() != 2 {
		t.Fatalf("expected total 2, got %d", summary.Total())
	}
}
