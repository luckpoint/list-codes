package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"myutilsport/utils"
)

var (
	folder     string
	outputFile string
	readmeOnly bool
	maxDepth   int
	debugMode  bool
	includes   []string
	excludes   []string
	prompt     string
	langFlag   string
)

func init() {
	// Parse --lang flag early before Cobra processes it
	langSet := parseLanguageFlag()
	
	// Initialize i18n with debug mode check
	debugMode := false
	for _, arg := range os.Args {
		if arg == "--debug" {
			debugMode = true
			break
		}
	}
	
	if !langSet {
		utils.InitI18n(debugMode)
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
	rootCmd.PersistentFlags().StringSliceVarP(&includes, "include", "i", []string{}, "Folder path to include (repeatable)")
	rootCmd.PersistentFlags().StringSliceVarP(&excludes, "exclude", "e", []string{}, "Folder path to exclude (repeatable)")
	rootCmd.PersistentFlags().Int64Var(&utils.MaxFileSizeBytes, "max-file-size", utils.MaxFileSizeBytesDefault, "Maximum file size in bytes to include in the summary")
	rootCmd.PersistentFlags().StringVarP(&prompt, "prompt", "p", "", "Prompt template name, file path, or custom prompt text to prepend to output")
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "Force language (ja|en) instead of auto-detection")

	rootCmd.AddCommand(completionCmd)
}

var rootCmd = &cobra.Command{
	Use:   "list-codes",
	// Short and Long will be set dynamically based on locale
	Run: func(cmd *cobra.Command, args []string) {
		excludeNames := make(map[string]struct{})
		for k := range utils.DefaultExcludeNames {
			excludeNames[k] = struct{}{}
		}
		utils.PrintDebug("Default exclude names: "+joinSet(excludeNames), debugMode)

		excludePaths := make(map[string]struct{})
		for _, p := range excludes {
			abs := p
			if !filepath.IsAbs(abs) {
				abs = filepath.Join(folder, p)
			}
			if resolved, err := filepath.Abs(abs); err == nil {
				excludePaths[resolved] = struct{}{}
			} else {
				utils.PrintWarning(fmt.Sprintf("Could not resolve absolute path for exclude '%s': %v", p, err), debugMode)
			}
		}
		utils.PrintDebug("User excluded absolute paths: "+joinSet(excludePaths), debugMode)

		includePaths := make(map[string]struct{})
		for _, p := range includes {
			abs := p
			if !filepath.IsAbs(abs) {
				abs = filepath.Join(folder, p)
			}
			if resolved, err := filepath.Abs(abs); err == nil {
				includePaths[resolved] = struct{}{}
			} else {
				utils.PrintWarning(fmt.Sprintf("Could not resolve absolute path for include '%s': %v", p, err), debugMode)
			}
		}
		utils.PrintDebug("User included absolute paths: "+joinSet(includePaths), debugMode)

		folderAbs, err := filepath.Abs(folder)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Could not resolve absolute path for folder '%s': %v", folder, err))
			os.Exit(1)
		}
		utils.PrintDebug("Scanning folder: "+folderAbs, debugMode)

		var outputMD string
		if readmeOnly {
			utils.PrintDebug("Mode: Collecting README.md files only.", debugMode)
			outputMD = utils.CollectReadmeFiles(folderAbs, includePaths, excludeNames, excludePaths, debugMode)
		} else {
			utils.PrintDebug("Mode: Summarizing project.", debugMode)
			outputMD = utils.ProcessSourceFiles(folderAbs, maxDepth, includePaths, excludeNames, excludePaths, debugMode)
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}