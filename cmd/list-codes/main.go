package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/luckpoint/list-codes/utils"
	"github.com/spf13/cobra"
)

var (
	folder          string
	outputFile      string
	readmeOnly      bool
	maxDepth        int
	debugMode       bool
	includes        []string
	excludes        []string
	prompt          string
	langFlag        string
	version         = "dev" // Will be overridden by build flags
	includeTests    bool
	maxFileSizeStr  string
	maxTotalSizeStr string
	noGitignore     bool
)

var doubleStarRegexpCache sync.Map

func init() {
	// Parse --lang flag early before Cobra processes it
	langSet := parseLanguageFlag()

	// Initialize i18n with debug mode check
	earlyDebugMode := false
	for _, arg := range os.Args {
		if arg == "--debug" {
			earlyDebugMode = true
			break
		}
	}

	if !langSet {
		utils.InitI18n(earlyDebugMode)
	}

	// Set help messages based on current language
	short, long, completionHelp := utils.GetHelpMessages()
	rootCmd.Short = short
	rootCmd.Long = long
	completionCmd.Long = completionHelp

	// Flag descriptions (always in English for technical consistency)
	rootCmd.PersistentFlags().StringVarP(&folder, "folder", "f", ".", "Folder to scan")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Output Markdown file path")
	rootCmd.PersistentFlags().BoolVar(&readmeOnly, "readme-only", false, "Only collect README.md files")
	rootCmd.PersistentFlags().IntVar(&maxDepth, "max-depth", utils.MaxStructureDepthDefault, "Max depth for directory structure")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringSliceVarP(&includes, "include", "i", []string{}, "Additional folder/file paths to include beyond defaults (repeatable)")
	rootCmd.PersistentFlags().StringSliceVarP(&excludes, "exclude", "e", []string{}, "Folder path to exclude (repeatable)")
	rootCmd.PersistentFlags().StringVar(&maxFileSizeStr, "max-file-size", "1m", "Maximum file size to include (e.g., 1m, 500k, 2g)")
	rootCmd.PersistentFlags().StringVar(&maxTotalSizeStr, "max-total-size", "", "Maximum total file size to collect (e.g., 10m, 1g) - empty means no limit")
	rootCmd.PersistentFlags().StringVarP(&prompt, "prompt", "p", "", "Prompt text or template name to prepend to output (accepts both predefined templates and custom text)")
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "Force language (ja|en) instead of auto-detection")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show version information")
	rootCmd.PersistentFlags().BoolVar(&includeTests, "include-tests", false, "Include test files in the output")
	rootCmd.PersistentFlags().BoolVar(&noGitignore, "no-gitignore", false, "Disable .gitignore file processing")

	// Register custom completion for --prompt flag
	rootCmd.RegisterFlagCompletionFunc("prompt", promptCompletion)

	rootCmd.AddCommand(completionCmd)
}

var rootCmd = &cobra.Command{
	Use: "list-codes",
	// Short and Long will be set dynamically based on locale
	Run: func(cmd *cobra.Command, args []string) {
		// Handle version flag
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Println("list-codes version", version)
			return
		}

		// Parse size strings to bytes
		var err error
		utils.MaxFileSizeBytes, err = utils.ParseSize(maxFileSizeStr)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Invalid --max-file-size: %v", err))
			os.Exit(1)
		}

		utils.TotalMaxFileSizeBytes, err = utils.ParseSize(maxTotalSizeStr)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Invalid --max-total-size: %v", err))
			os.Exit(1)
		}

		excludeNames := make(map[string]struct{})
		for k := range utils.DefaultExcludeNames {
			excludeNames[k] = struct{}{}
		}
		utils.PrintDebug("Default exclude names: "+joinSet(excludeNames), debugMode)

		folderAbs, err := filepath.Abs(folder)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Could not resolve absolute path for folder '%s': %v", folder, err))
			os.Exit(1)
		}
		utils.PrintDebug("Scanning folder: "+folderAbs, debugMode)

		// Process excludes
		var excludePatterns []string
		for _, p := range excludes {
			// If it contains glob characters, handle anchoring
			if strings.ContainsAny(p, "*?[]") {
				pattern := filepath.ToSlash(p)
				// If it doesn't contain "**" and doesn't start with "/", anchor it to root
				// This follows the requirement: "*.md" matches root only, "**/*.md" matches recursively.
				if !strings.Contains(pattern, "**") && !strings.HasPrefix(pattern, "/") {
					pattern = "/" + pattern
				}
				excludePatterns = append(excludePatterns, pattern)
				continue
			}

			// Resolve to absolute path first to handle "relative to current dir" vs "relative to folder"
			abs := p
			if !filepath.IsAbs(abs) {
				abs = filepath.Join(folder, p)
			}
			abs, err = filepath.Abs(abs)
			if err == nil {
				rel, err := filepath.Rel(folderAbs, abs)
				if err == nil {
					// Windows backslash to slash
					rel = filepath.ToSlash(rel)
					// Anchor to root to preserve "path relative to folder" semantics
					// and prevent accidental matching of deeply nested files with same name
					if !strings.HasPrefix(rel, "/") {
						rel = "/" + rel
					}
					excludePatterns = append(excludePatterns, rel)
				}
			} else {
				utils.PrintWarning(fmt.Sprintf("Could not resolve path for exclude '%s': %v", p, err), debugMode)
			}
		}

		excludeMatcher, err := utils.NewSimpleMatcher(folderAbs, excludePatterns)
		if err != nil {
			utils.PrintWarning(fmt.Sprintf("Failed to create exclude matcher: %v", err), debugMode)
		}
		if len(excludePatterns) > 0 {
			utils.PrintDebug("User exclude patterns: "+strings.Join(excludePatterns, ", "), debugMode)
		}

		// Process includes
		includePaths := make(map[string]struct{})
		var includePatterns []string

		for _, p := range includes {
			// If it contains glob characters, handle anchoring
			if strings.ContainsAny(p, "*?[]") {
				pattern := filepath.ToSlash(p)
				// If it doesn't contain "**" and doesn't start with "/", anchor it to root
				if !strings.Contains(pattern, "**") && !strings.HasPrefix(pattern, "/") {
					pattern = "/" + pattern
				}
				includePatterns = append(includePatterns, pattern)
				continue
			}

			abs := p
			if !filepath.IsAbs(abs) {
				abs = filepath.Join(folder, p)
			}
			resolved, err := filepath.Abs(abs)
			if err == nil {
				// Add to map for "parent of" traversal logic
				includePaths[resolved] = struct{}{}

				// Add to patterns
				rel, err := filepath.Rel(folderAbs, resolved)
				if err == nil {
					rel = filepath.ToSlash(rel)
					if !strings.HasPrefix(rel, "/") {
						rel = "/" + rel
					}
					includePatterns = append(includePatterns, rel)
				}
			} else {
				utils.PrintWarning(fmt.Sprintf("Could not resolve absolute path for include '%s': %v", p, err), debugMode)
			}
		}

		includeMatcher, err := utils.NewSimpleMatcher(folderAbs, includePatterns)
		if err != nil {
			utils.PrintWarning(fmt.Sprintf("Failed to create include matcher: %v", err), debugMode)
		}
		utils.PrintDebug("User included absolute paths: "+joinSet(includePaths), debugMode)

		// Create GitIgnoreMatcher if --no-gitignore is not set
		var gitIgnoreMatcher *utils.GitIgnoreMatcher
		if !noGitignore {
			matcher, err := utils.NewGitIgnoreMatcher(folderAbs)
			if err != nil {
				utils.PrintWarning(fmt.Sprintf("Could not create gitignore matcher: %v", err), debugMode)
			} else {
				gitIgnoreMatcher = matcher
				utils.PrintDebug("GitIgnore matcher created successfully", debugMode)
			}
		} else {
			utils.PrintDebug("GitIgnore processing disabled via --no-gitignore flag", debugMode)
		}

		var outputMD string
		if readmeOnly {
			utils.PrintDebug("Mode: Collecting README.md files only.", debugMode)
			outputMD = utils.CollectReadmeFiles(folderAbs, includePaths, includeMatcher, excludeNames, excludeMatcher, debugMode, gitIgnoreMatcher)
		} else {
			utils.PrintDebug("Mode: Summarizing project.", debugMode)
			outputMD = utils.ProcessSourceFiles(folderAbs, maxDepth, includePaths, includeMatcher, excludeNames, excludeMatcher, debugMode, includeTests, gitIgnoreMatcher)
		}

		// Apply prompt if specified
		if prompt != "" {
			promptText, err := utils.GetPrompt(prompt, debugMode)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Could not process prompt '%s': %v", prompt, err))
				os.Exit(1)
			}
			outputMD = utils.FormatWithPrompt(promptText, outputMD)
			utils.PrintDebug("Applied prompt to output", debugMode)
		}

		if err := utils.SaveToMarkdown(outputMD, outputFile); err != nil {
			utils.PrintError(fmt.Sprintf("Could not save output to '%s': %v", outputFile, err))
			os.Exit(1)
		}

		utils.PrintDebug("Processing complete.", debugMode)
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: utils.T("Generate completion script", "補完スクリプトを生成"),
	// Long will be set in init() based on locale
	Args: cobra.ExactValidArgs(1), // Only one arg, which is the name of the app
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func joinSet(m map[string]struct{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func resolvePathPatterns(baseFolder string, patterns []string, flagName string, debug bool) map[string]struct{} {
	resolvedPaths := make(map[string]struct{})

	for _, rawPattern := range patterns {
		pattern := strings.TrimSpace(rawPattern)
		if pattern == "" {
			continue
		}

		fullPattern := pattern
		if !filepath.IsAbs(fullPattern) {
			fullPattern = filepath.Join(baseFolder, fullPattern)
		}

		expanded, err := expandPathPattern(fullPattern)
		if err != nil {
			utils.PrintWarning(fmt.Sprintf("Could not expand %s pattern '%s': %v", flagName, pattern, err), debug)
			continue
		}

		for _, p := range expanded {
			absPath, err := toAbsCleanPath(p)
			if err != nil {
				utils.PrintWarning(fmt.Sprintf("Could not resolve absolute path for %s '%s': %v", flagName, p, err), debug)
				continue
			}
			resolvedPaths[absPath] = struct{}{}
		}
	}

	return resolvedPaths
}

func expandPathPattern(pattern string) ([]string, error) {
	cleanPattern := filepath.Clean(pattern)
	if !hasGlobMeta(cleanPattern) {
		return []string{cleanPattern}, nil
	}

	if !strings.Contains(cleanPattern, "**") {
		matches, err := filepath.Glob(cleanPattern)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			// Keep the literal pattern to ensure include-only mode still takes effect.
			return []string{cleanPattern}, nil
		}
		return matches, nil
	}

	searchRoot := patternSearchRoot(cleanPattern)
	searchRootAbs, err := toAbsCleanPath(searchRoot)
	if err != nil {
		return nil, err
	}

	re, err := getCachedDoubleStarRegexp(cleanPattern)
	if err != nil {
		return nil, err
	}

	var matches []string
	walkErr := filepath.WalkDir(searchRootAbs, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		normalizedPath := filepath.ToSlash(filepath.Clean(path))
		if re.MatchString(normalizedPath) {
			matches = append(matches, path)
		}
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	if len(matches) == 0 {
		// Keep the literal pattern to ensure include-only mode still takes effect.
		return []string{cleanPattern}, nil
	}
	return matches, nil
}

func toAbsCleanPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Abs(path)
}

func getCachedDoubleStarRegexp(pattern string) (*regexp.Regexp, error) {
	normalized := filepath.ToSlash(filepath.Clean(pattern))
	if cached, ok := doubleStarRegexpCache.Load(normalized); ok {
		return cached.(*regexp.Regexp), nil
	}

	re, err := regexp.Compile(doubleStarPatternToRegexp(normalized))
	if err != nil {
		return nil, err
	}
	if cached, loaded := doubleStarRegexpCache.LoadOrStore(normalized, re); loaded {
		return cached.(*regexp.Regexp), nil
	}
	return re, nil
}

func hasGlobMeta(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[")
}

func patternSearchRoot(pattern string) string {
	idx := strings.IndexAny(pattern, "*?[")
	if idx < 0 {
		return pattern
	}

	prefix := pattern[:idx]
	root := filepath.Dir(prefix)
	if root == "." || root == "" {
		return string(filepath.Separator)
	}
	return root
}

func doubleStarPatternToRegexp(pattern string) string {
	var b strings.Builder
	b.WriteString("^")

	for i := 0; i < len(pattern); {
		ch := pattern[i]
		switch ch {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				if i+2 < len(pattern) && pattern[i+2] == '/' {
					// `**/` means zero or more nested directories.
					b.WriteString("(?:[^/]+/)*")
					i += 3
					continue
				}
				b.WriteString(".*")
				i += 2
				continue
			} else {
				b.WriteString("[^/]*")
			}
		case '?':
			b.WriteString("[^/]")
		default:
			b.WriteString(regexp.QuoteMeta(string(ch)))
		}
		i++
	}

	b.WriteString("$")
	return b.String()
}

// promptCompletion provides completion for --prompt flag
func promptCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	availablePrompts := utils.GetAvailablePrompts()
	sort.Strings(availablePrompts)
	return availablePrompts, cobra.ShellCompDirectiveNoFileComp
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
