package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jbradley/dns-discovery/internal/app"
	appconfig "github.com/jbradley/dns-discovery/internal/config"
	"github.com/spf13/cobra"
)

var outputDir string
var configPath string
var inputFile string

const defaultConfigFile = ".dns-discovery.json"

var rootCmd = &cobra.Command{
	Use:   "dns-discovery [domain]",
	Short: "DNS zone discovery tool",
	Long: `Discover DNS configuration, provider fingerprint, TLS health,
and email DNS health for any domain.`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadRuntimeConfig(configPath)
		if err != nil {
			return err
		}

		resolved := appconfig.Resolve(outputDir, cmd.Flags().Changed("output-dir"), cfg)
		domains, err := resolveDomains(args, inputFile, cfg)
		if err != nil {
			return err
		}

		if len(domains) == 1 {
			return app.RunDomain(domains[0], resolved.OutputDir)
		}

		summary := app.RunBatch(domains, resolved.OutputDir)
		printBatchSummary(summary)
		if len(summary.Failed) > 0 {
			return fmt.Errorf("batch completed with %d failed domain(s)", len(summary.Failed))
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&outputDir, "output-dir", "o", appconfig.DefaultOutputDir, "Directory to save reports")
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to JSON config file")
	rootCmd.Flags().StringVar(&inputFile, "input-file", "", "Path to newline-delimited domains file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadRuntimeConfig(explicitPath string) (appconfig.Config, error) {
	if strings.TrimSpace(explicitPath) != "" {
		return appconfig.Load(explicitPath)
	}

	if _, err := os.Stat(defaultConfigFile); err == nil {
		return appconfig.Load(defaultConfigFile)
	} else if !os.IsNotExist(err) {
		return appconfig.Config{}, fmt.Errorf("check config %q: %w", defaultConfigFile, err)
	}

	return appconfig.Config{}, nil
}

func resolveDomains(args []string, inputPath string, cfg appconfig.Config) ([]string, error) {
	if len(args) > 0 {
		domain := strings.TrimSpace(strings.ToLower(args[0]))
		if domain == "" {
			return nil, fmt.Errorf("domain argument must not be empty")
		}
		return []string{domain}, nil
	}

	if strings.TrimSpace(inputPath) != "" {
		return loadDomainsFromFile(inputPath)
	}

	if len(cfg.Domains) > 0 {
		domains := make([]string, 0, len(cfg.Domains))
		for _, domain := range cfg.Domains {
			domains = append(domains, strings.ToLower(domain))
		}
		return domains, nil
	}

	return nil, fmt.Errorf("no domains provided: pass a domain argument, --input-file, or config domains in %s", defaultConfigFile)
}

func loadDomainsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read input file %q: %w", path, err)
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		domains = append(domains, strings.ToLower(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read input file %q: %w", path, err)
	}
	if len(domains) == 0 {
		return nil, fmt.Errorf("input file %q did not contain any domains", path)
	}

	return domains, nil
}

func printBatchSummary(summary app.BatchSummary) {
	sep := strings.Repeat("=", 60)
	fmt.Printf("\n%s\n  Batch Summary\n%s\n", sep, sep)
	fmt.Printf("  Total:      %d\n", summary.Total())
	fmt.Printf("  Succeeded:  %d\n", len(summary.Succeeded))
	fmt.Printf("  Failed:     %d\n", len(summary.Failed))

	if len(summary.Failed) > 0 {
		fmt.Println("\n  Failed Domains:")
		failedDomains := summary.FailedDomains()
		for _, domain := range failedDomains {
			fmt.Printf("     • %s: %v\n", domain, summary.Failed[domain])
		}
	}

	if len(summary.Succeeded) > 0 {
		succeeded := append([]string(nil), summary.Succeeded...)
		sort.Strings(succeeded)
		fmt.Println("\n  Succeeded Domains:")
		for _, domain := range succeeded {
			fmt.Printf("     • %s\n", domain)
		}
	}

	fmt.Printf("\n%s\n", sep)
}
