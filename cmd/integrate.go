package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

const defaultRegistryURL = "https://gist.githubusercontent.com/arpitbbhayani/ab067e0ef257badf83739e636db580c2/raw/f68f2cea1a0024f2cf68114c1b7bd046952626f6/all"

type tutorial struct {
	Framework string
	URL       string
}

type registry struct {
	Tutorials []tutorial
}

type integrationStep struct {
	Name   string
	Prompt string
}

var (
	integrateLang        string
	integrateFramework   string
	integrateRegistryURL string
	integrateDryRun      bool
)

var integrateCmd = &cobra.Command{
	Use:   "integrate",
	Short: "Integrate Razorpay into the current project",
	Long: `Detect the language and framework of the project in the current directory,
fetch the appropriate Razorpay integration tutorial, and use the locally
installed Claude CLI to perform the integration step by step.

Examples:
  razorpay integrate
  razorpay integrate --language python --framework django
  razorpay integrate --dry-run`,
	RunE: runIntegrate,
}

func runIntegrate(cmd *cobra.Command, args []string) error {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude CLI not found in PATH. Install it from: https://docs.anthropic.com/en/docs/claude-code")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	language := integrateLang
	framework := integrateFramework

	// Detect language and framework using Claude
	if language == "" || framework == "" {
		printStep("Detecting language and framework")

		detectPrompt := `Look at the files in the current directory and determine:
1. What programming language this project uses
2. What framework (if any) this project uses

Reply with ONLY two lines in this exact format, nothing else:
language=<language>
framework=<framework>

Use lowercase. For language use: python, node, go, ruby, php, java, dotnet.
For framework use: django, flask, fastapi, express, nextjs, react, angular, vue, rails, laravel, spring, or generic if none detected.`

		sp := startSpinner("analyzing project...")
		output, err := runClaude(claudePath, cwd, detectPrompt)
		sp.stop()
		if err != nil {
			return fmt.Errorf("failed to detect project: %w", err)
		}

		printResult(output)

		for _, line := range strings.Split(output, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "language=") {
				language = strings.TrimPrefix(line, "language=")
			}
			if strings.HasPrefix(line, "framework=") {
				framework = strings.TrimPrefix(line, "framework=")
			}
		}

		if language == "" {
			language = "generic"
		}
		if framework == "" {
			framework = "generic"
		}
	} else {
		printStep("Using provided language and framework")
		printResult(fmt.Sprintf("Language: %s\nFramework: %s", language, framework))
	}

	// Fetch tutorial from registry
	printStep("Fetching integration tutorial")

	regURL := integrateRegistryURL
	if regURL == "" {
		regURL = defaultRegistryURL
	}

	sp := startSpinner("fetching registry...")
	reg, err := fetchRegistry(regURL)
	sp.stop()
	if err != nil {
		return fmt.Errorf("failed to fetch tutorial registry: %w", err)
	}

	tutorialEntry := matchTutorial(reg, language, framework)
	if tutorialEntry == nil {
		return fmt.Errorf("no tutorial found for language=%q framework=%q", language, framework)
	}

	printResult(fmt.Sprintf("Matched: %s\nURL: %s", tutorialEntry.Framework, tutorialEntry.URL))

	if integrateDryRun {
		fmt.Println("\n[dry-run] Would fetch tutorial and run integration steps. Exiting.")
		return nil
	}

	// Download tutorial content
	printStep("Downloading tutorial content")

	sp = startSpinner("downloading guide...")
	tutorialContent, err := fetchContent(tutorialEntry.URL)
	sp.stop()
	if err != nil {
		return fmt.Errorf("failed to fetch tutorial content: %w", err)
	}

	printResult(fmt.Sprintf("Downloaded %d bytes of integration guide.", len(tutorialContent)))

	// Step 4+: Run integration steps via Claude
	steps := []integrationStep{
		{
			Name: "Install dependencies",
			Prompt: "Based on the integration guide, install ONLY the required packages/dependencies. " +
				"Do not create any application files yet. Just install what's needed. " +
				"Show what you installed.",
		},
		{
			Name: "Create integration files",
			Prompt: "Based on the integration guide, create the necessary integration files " +
				"(configuration, helper modules, route handlers, etc). " +
				"Follow existing code patterns in the project. " +
				"Show what files you created or modified.",
		},
		{
			Name: "Add environment configuration",
			Prompt: "Based on the integration guide, add any required environment variables or " +
				"configuration entries (e.g. .env, config files). " +
				"Use placeholder values like RAZORPAY_KEY_ID and RAZORPAY_KEY_SECRET. " +
				"Show what you configured.",
		},
		{
			Name: "Verify integration",
			Prompt: "Verify the Razorpay integration is correct: " +
				"check that imports resolve, files are syntactically valid, and the app can start. " +
				"Run a quick syntax/build check if possible. " +
				"Report any issues found.",
		},
	}

	guideContext := fmt.Sprintf("Project: %s/%s\n\n%s", language, framework, tutorialContent)

	for _, step := range steps {
		printStep(step.Name)

		stepPrompt := step.Prompt + "\n\n--- INTEGRATION GUIDE ---\n" + guideContext + "\n--- END GUIDE ---"
		sp = startSpinner(strings.ToLower(step.Name) + "...")
		output, err := runClaude(claudePath, cwd, stepPrompt)
		sp.stop()
		if err != nil {
			printResult("[error] " + err.Error())
			continue
		}

		printResult(output)
	}

	printStep("Integration complete")
	fmt.Println("Set RAZORPAY_KEY_ID and RAZORPAY_KEY_SECRET environment variables to get started.")

	return nil
}

func printStep(name string) {
	fmt.Printf("\n→ %s\n", name)
}

func printResult(output string) {
	for _, line := range strings.Split(output, "\n") {
		if line != "" {
			fmt.Println(line)
		}
	}
}

type spinner struct {
	msg  string
	done chan struct{}
	wg   sync.WaitGroup
}

func startSpinner(msg string) *spinner {
	s := &spinner{msg: msg, done: make(chan struct{})}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.done:
				fmt.Printf("\r\033[K")
				return
			default:
				fmt.Printf("\r  %s %s", frames[i%len(frames)], msg)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	return s
}

func (s *spinner) stop() {
	close(s.done)
	s.wg.Wait()
}

func runClaude(claudePath, cwd, prompt string) (string, error) {
	claudeCmd := exec.Command(claudePath, "-p", "--allowedTools", "Bash,Edit,Write,Read")
	claudeCmd.Dir = cwd
	claudeCmd.Stdin = strings.NewReader(prompt)

	var stdout, stderr bytes.Buffer
	claudeCmd.Stdout = &stdout
	claudeCmd.Stderr = &stderr

	err := claudeCmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("%s", errMsg)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func fetchRegistry(registryURL string) (*registry, error) {
	resp, err := http.Get(registryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned HTTP %d", resp.StatusCode)
	}

	var reg registry
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " - ", 2)
		if len(parts) != 2 {
			continue
		}
		reg.Tutorials = append(reg.Tutorials, tutorial{
			Framework: strings.TrimSpace(parts[0]),
			URL:       strings.TrimSpace(parts[1]),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read registry: %w", err)
	}
	return &reg, nil
}

func matchTutorial(reg *registry, language, framework string) *tutorial {
	for i := range reg.Tutorials {
		t := &reg.Tutorials[i]
		if strings.EqualFold(t.Framework, framework) {
			return t
		}
	}
	for i := range reg.Tutorials {
		t := &reg.Tutorials[i]
		if strings.EqualFold(t.Framework, language) {
			return t
		}
	}
	for i := range reg.Tutorials {
		t := &reg.Tutorials[i]
		if strings.EqualFold(t.Framework, "generic") {
			return t
		}
	}
	if len(reg.Tutorials) > 0 {
		return &reg.Tutorials[0]
	}
	return nil
}

func fetchContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tutorial URL returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func init() {
	integrateCmd.Flags().StringVarP(&integrateLang, "language", "l", "", "Override detected language (e.g., python, node, go, ruby, php, java, dotnet)")
	integrateCmd.Flags().StringVarP(&integrateFramework, "framework", "f", "", "Override detected framework (e.g., django, express, nextjs, rails, laravel, spring)")
	integrateCmd.Flags().StringVar(&integrateRegistryURL, "registry-url", "", "Custom tutorial registry URL (default: "+defaultRegistryURL+")")
	integrateCmd.Flags().BoolVar(&integrateDryRun, "dry-run", false, "Show detection results without invoking Claude")
}
