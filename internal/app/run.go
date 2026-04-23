package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbradley/dns-discovery/internal/discovery"
	"github.com/jbradley/dns-discovery/internal/report"
)

type OutputFormat string

const (
	OutputMarkdown OutputFormat = "markdown"
	OutputJSON     OutputFormat = "json"
	OutputText     OutputFormat = "text"
)

const (
	defaultOutputDir   = "output"
	defaultLogLocation = "logs/dns-discovery.log"
)

type RunOptions struct {
	OutputDir   string
	Output      OutputFormat
	Verbose     bool
	LogLocation string
}

type DomainSuccess struct {
	Domain     string
	ReportPath string
}

type DomainFailure struct {
	Domain string
	Err    error
}

type BatchSummary struct {
	Succeeded []DomainSuccess
	Failed    []DomainFailure
}

func (summary BatchSummary) Total() int {
	return len(summary.Succeeded) + len(summary.Failed)
}

func ValidateOutputFormat(value string) (OutputFormat, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch OutputFormat(normalized) {
	case OutputMarkdown, OutputJSON, OutputText:
		return OutputFormat(normalized), nil
	default:
		return "", fmt.Errorf("invalid output format %q: expected one of markdown, json, text", value)
	}
}

type runLoggerSink interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Close() error
}

type runLogger struct {
	file       *os.File
	fileLogger *log.Logger
	verbose    bool
}

var testHookScanner = scanDomain
var testHookReportWriter = func(baseDir string, res *discovery.DiscoveryResult, output OutputFormat) (string, error) {
	return report.SaveReportByFormat(baseDir, res, string(output))
}
var testHookLoggerFactory = func(logLocation string, verbose bool) (runLoggerSink, error) {
	return newRunLogger(logLocation, verbose)
}

func RunDiscovery(domains []string, opts RunOptions) (BatchSummary, error) {
	output, err := resolveOutputFormat(opts.Output)
	if err != nil {
		return BatchSummary{}, err
	}

	outputDir := strings.TrimSpace(opts.OutputDir)
	if outputDir == "" {
		outputDir = defaultOutputDir
	}

	logPath, err := resolveLogPath(opts.LogLocation)
	if err != nil {
		return BatchSummary{}, err
	}

	logger, err := testHookLoggerFactory(logPath, opts.Verbose)
	if err != nil {
		return BatchSummary{}, err
	}
	defer logger.Close()

	summary := BatchSummary{}
	for _, domain := range domains {
		normalized := strings.ToLower(strings.TrimSpace(domain))
		if normalized == "" {
			continue
		}

		logger.Infof("starting discovery for %s", normalized)
		result, err := testHookScanner(normalized)
		if err != nil {
			logger.Errorf("discovery failed for %s: %v", normalized, err)
			summary.Failed = append(summary.Failed, DomainFailure{Domain: normalized, Err: err})
			continue
		}

		reportPath, err := testHookReportWriter(outputDir, result, output)
		if err != nil {
			logger.Errorf("report generation failed for %s: %v", normalized, err)
			summary.Failed = append(summary.Failed, DomainFailure{Domain: normalized, Err: err})
			continue
		}

		logger.Infof("completed discovery for %s; report=%s", normalized, reportPath)
		summary.Succeeded = append(summary.Succeeded, DomainSuccess{Domain: normalized, ReportPath: reportPath})
	}

	return summary, nil
}

func scanDomain(domain string) (*discovery.DiscoveryResult, error) {
	records, err := discovery.QueryAllRecords(domain)
	if err != nil {
		return nil, err
	}

	services := discovery.DetectServices(records)
	nsHosts := records["NS"]
	providerResult := discovery.IdentifyProviders(domain, nsHosts)

	var tlsResults []discovery.TLSResult
	for range records["A"] {
		tlsResults = append(tlsResults, discovery.CheckTLS(domain))
		break
	}

	emailResult := discovery.EvaluateEmailHealth(domain)

	return &discovery.DiscoveryResult{
		Domain:   domain,
		DNS:      records,
		Services: services,
		Provider: providerResult,
		TLS:      tlsResults,
		Email:    emailResult,
	}, nil
}

func resolveOutputFormat(output OutputFormat) (OutputFormat, error) {
	normalized := strings.TrimSpace(string(output))
	if normalized == "" {
		normalized = string(OutputMarkdown)
	}
	return ValidateOutputFormat(normalized)
}

func resolveLogPath(logLocation string) (string, error) {
	location := strings.TrimSpace(logLocation)
	if location == "" {
		location = defaultLogLocation
	}

	if filepath.IsAbs(location) {
		return location, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}

	repoRoot, err := findRepoRoot(cwd)
	if err != nil {
		return "", err
	}

	return filepath.Join(repoRoot, location), nil
}

func findRepoRoot(startDir string) (string, error) {
	current := startDir
	for {
		candidate := filepath.Join(current, "go.mod")
		if _, err := os.Stat(candidate); err == nil {
			return current, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("check %q: %w", candidate, err)
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("could not find repository root containing go.mod from %q", startDir)
}

func newRunLogger(logLocation string, verbose bool) (runLoggerSink, error) {
	if err := os.MkdirAll(filepath.Dir(logLocation), 0755); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	file, err := os.OpenFile(logLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file %q: %w", logLocation, err)
	}

	return &runLogger{
		file:       file,
		fileLogger: log.New(file, "", log.LstdFlags),
		verbose:    verbose,
	}, nil
}

func (l *runLogger) Infof(format string, args ...any) {
	message := fmt.Sprintf("INFO  "+format, args...)
	l.fileLogger.Println(message)
	if l.verbose {
		fmt.Println(message)
	}
}

func (l *runLogger) Errorf(format string, args ...any) {
	message := fmt.Sprintf("ERROR "+format, args...)
	l.fileLogger.Println(message)
	fmt.Fprintln(os.Stderr, message)
}

func (l *runLogger) Close() error {
	return l.file.Close()
}
