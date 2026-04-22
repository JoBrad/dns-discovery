package app

import (
	"errors"
	"testing"
)

func TestRunDomainDelegatesToRunner(t *testing.T) {
	original := domainRunner
	defer func() {
		domainRunner = original
	}()

	called := false
	domainRunner = func(domain string, outputDir string) error {
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