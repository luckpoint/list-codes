package utils

const MaxStructureDepthDefault = 7
const MaxFileSizeBytesDefault = 1024 * 1024 // 1MB

var MaxFileSizeBytes int64 = MaxFileSizeBytesDefault
var TotalMaxFileSizeBytes int64 = 0 // 0 means no limit

const DependencyFilesCategory = "Dependency Files"

// DefaultExcludeNames contains default directory and file names to exclude by name.
// Dotfiles are now handled by a general rule in shouldSkipDir
var DefaultExcludeNames = map[string]struct{}{
	"node_modules": {},
	"vendor":       {},
	"target":       {},
	"build":        {},
	"dist":         {},
	"__pycache__":  {},
	"env":          {},
	"venv":         {},
}

// PROJECT_SIGNATURES contains signature files/directories for detecting project languages/frameworks.
var PROJECT_SIGNATURES = map[string][]string{
	"Go":         {"go.mod"},
	"Ruby":       {"Gemfile"},
	"Javascript": {"package.json", "node_modules"},
	"Typescript": {"tsconfig.json"},
	"Python":     {"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"},
	"Rust":       {"Cargo.toml"},
	"SQL":        {".sql"},
	"HTML":       {".html", ".erb", ".haml", ".slim"},
	"CSS":        {".css", ".scss", ".sass"},
	"Lua":        {".lua"},
	"C":          {".c"},
	"C++":        {".cpp", ".hpp", ".cxx", ".hxx", ".cc", ".hh", ".h"},
	"Java":       {"build.gradle", "settings.gradle", "pom.xml"},
	"PHP":        {"composer.json"},
	"C#":         {".csproj", ".sln", "project.json"},
	"F#":         {".fsproj"},
	"VB.NET":     {".vbproj"},
	"Scala":      {"build.sbt"},
	"Clojure":    {"project.clj", "deps.edn"},
	"Haskell":    {".cabal", "stack.yaml"},
	"Erlang":     {"rebar.config"},
	"Elixir":     {"mix.exs"},
	"Dart":       {"pubspec.yaml"},
	"R":          {"DESCRIPTION"},
	"Julia":      {"Project.toml"},
	"Swift":      {"Package.swift"},
	"Kotlin":     {"build.gradle.kts"},
	"Groovy":     {"build.gradle"},
	"CMake":      {"CMakeLists.txt"},
	"Terraform":  {"main.tf"},
	"Protobuf/Buf": {".proto", "buf.yaml"},
	"Markdown":   {"README.md"},
	"MDX":        {".mdx"},
	"Solidity":   {"hardhat.config.js", "hardhat.config.ts", "truffle-config.js", "foundry.toml"},
}

// EXTENSIONS maps file extensions to programming languages for detection.
var EXTENSIONS = map[string][]string{
	"Python":     {".py", ".pyw"},
	"Ruby":       {".rb", ".rbw"},
	"Javascript": {".js", ".jsx", ".mjs", ".cjs"},
	"Typescript": {".ts", ".tsx", ".mts", ".cts"},
	"Go":         {".go"},
	"Swift":      {".swift"},
	"Kotlin":     {".kt", ".kts"},
	"Rust":       {".rs"},
	"SQL":        {".sql"},
	"HTML":       {".html", ".htm", ".erb", ".haml", ".slim"},
	"CSS":        {".css", ".scss", ".sass", ".less", ".styl"},
	"Lua":        {".lua"},
	"C":          {".c", ".h"},
	"C++":        {".cpp", ".hpp", ".cxx", ".hxx", ".cc", ".hh", ".c++", ".h++"},
	"Java":       {".java"},
	"PHP":        {".php", ".phtml", ".php3", ".php4", ".php5"},
	"C#":         {".cs", ".csx"},
	"F#":         {".fs", ".fsi", ".fsx"},
	"VB.NET":     {".vb"},
	"Scala":      {".scala", ".sc"},
	"Clojure":    {".clj", ".cljs", ".cljc"},
	"Haskell":    {".hs", ".lhs"},
	"Erlang":     {".erl", ".hrl"},
	"Elixir":     {".ex", ".exs"},
	"Dart":       {".dart"},
	"R":          {".r", ".R"},
	"Perl":       {".pl", ".pm", ".perl"},
	"Objective-C": {".m", ".mm"},
	"Pascal":     {".pas", ".pp"},
	"Fortran":    {".f", ".f90", ".f95", ".f03", ".f08"},
	"COBOL":      {".cob", ".cbl"},
	"Ada":        {".ada", ".adb", ".ads"},
	"Lisp":       {".lisp", ".lsp", ".cl"},
	"Scheme":     {".scm", ".ss"},
	"Prolog":     {".pl", ".pro"},
	"MATLAB":     {".m"},
	"Octave":     {".m"},
	"Julia":      {".jl"},
	"Zig":        {".zig"},
	"D":          {".d"},
	"Nim":        {".nim"},
	"Crystal":    {".cr"},
	"Groovy":     {".groovy", ".gvy"},
	"PowerShell": {".ps1", ".psm1", ".psd1"},
	"Bash":       {".bash"},
	"Zsh":        {".zsh"},
	"Fish":       {".fish"},
	"Tcl":        {".tcl"},
	"Vim":        {".vim"},
	"Emacs Lisp": {".el"},
	"Protobuf/Buf": {".proto"},
	"Markdown":   {".md"},
	"MDX":        {".mdx"},
	"Solidity":   {".sol"},
	"YAML":       {".yaml", ".yml"},
	"JSON":       {".json"},
	"XML":        {".xml", ".xsd", ".xsl"},
	"Shell":      {".sh"},
	"Dockerfile": {"Dockerfile"},
	"TOML":       {".toml"},
	"Gradle":     {".gradle", ".gradle.kts"},
	"INI":        {".ini", ".cfg"},
	"Properties": {".properties"},
	"Makefile":   {"Makefile", "makefile", ".mk"},
	"CMake":      {".cmake", "CMakeLists.txt"},
	"Terraform":  {".tf", ".tfvars"},
	"HCL":        {".hcl"},
	"GraphQL":    {".graphql", ".gql"},
	"Jupyter":    {".ipynb"},
	"LaTeX":      {".tex", ".latex"},
	"Regex":      {".regex"},
	"Diff":       {".diff", ".patch"},
	"Log":        {".log"},
	"CSV":        {".csv"},
	"TSV":        {".tsv"},
	"WebAssembly": {".wat", ".wasm"},
	"Assembly":   {".asm", ".s"},
}

// FRAMEWORK_DEPENDENCY_FILES contains framework-specific dependency files.
var FRAMEWORK_DEPENDENCY_FILES = map[string][]string{
	"Ruby":       {"Gemfile", "Gemfile.lock"},
	"Javascript": {"package.json", "package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb"},
	"Typescript": {"package.json", "package-lock.json", "tsconfig.json"},
	"Python":     {"requirements.txt", "Pipfile", "Pipfile.lock", "pyproject.toml", "poetry.lock"},
	"Java":       {"build.gradle", "settings.gradle", "pom.xml"},
	"PHP":        {"composer.json", "composer.lock"},
	"C#":         {".csproj", ".sln", "project.json", "packages.config"},
	"F#":         {".fsproj"},
	"VB.NET":     {".vbproj"},
	"Scala":      {"build.sbt"},
	"Clojure":    {"project.clj", "deps.edn"},
	"Haskell":    {".cabal", "stack.yaml"},
	"Erlang":     {"rebar.config"},
	"Elixir":     {"mix.exs", "mix.lock"},
	"Dart":       {"pubspec.yaml", "pubspec.lock"},
	"R":          {"DESCRIPTION"},
	"Julia":      {"Project.toml", "Manifest.toml"},
	"Swift":      {"Package.swift"},
	"Kotlin":     {"build.gradle.kts"},
	"Groovy":     {"build.gradle"},
	"CMake":      {"CMakeLists.txt"},
	"Terraform":  {"main.tf", "variables.tf", "outputs.tf"},
	"Protobuf/Buf": {"buf.yaml"},
	"Go":         {"go.mod", "go.sum"},
	"Rust":       {"Cargo.toml", "Cargo.lock"},
	"Solidity":   {"hardhat.config.js", "hardhat.config.ts", "truffle-config.js", "foundry.toml"},
}

// EXCLUDE_TEST_KEYWORDS contains keywords for identifying test files.
// The IsTestFile function performs a case-insensitive check against the base filename.
var EXCLUDE_TEST_KEYWORDS = []string{
	"test",      // Covers "test_*.py", "*.test.js", "Test*.java", etc.
	"spec",      // Covers "*.spec.js", "*_spec.rb", etc.
	"e2e",       // For end-to-end tests
	"benchmark", // For benchmark files
	"bench",     // A common shorthand for benchmark
	"mock",      // For mock files
	"fixture",   // For test fixtures
}

// EXCLUDE_TEST_DIRS contains directory paths for identifying test files.
// The IsTestFile function checks if the file's path contains any of these strings.
var EXCLUDE_TEST_DIRS = []string{
	"/test/",      // Root-level "test" folder
	"/tests/",     // Root-level "tests" folder
	"/spec/",      // Root-level "spec" folder
	"/specs/",     // Root-level "specs" folder
	"/__tests__/", // A common convention in Jest (JavaScript)
	"/__test__/",  // Alternative Jest convention
	"/testing/",   // Generic testing folder
	"/fixtures/",  // Test fixtures
	"/mocks/",     // Mock files
	"/e2e/",       // End-to-end tests
	"/integration/", // Integration tests
	"/unit/",      // Unit tests
	"test/",       // "test" folder at any depth
	"tests/",      // "tests" folder at any depth
	"spec/",       // "spec" folder at any depth
	"specs/",      // "specs" folder at any depth
}

// EXCLUDE_TEST_PATTERNS contains specific filename patterns to identify test files.
// This is useful for conventions not easily caught by keywords alone.
var EXCLUDE_TEST_PATTERNS = []string{
	"_test.go",    // Go
	"_spec.rb",    // Ruby (RSpec)
	".test.js",    // JavaScript
	".spec.js",    // JavaScript
	".test.jsx",   // React (JSX)
	".spec.jsx",   // React (JSX)
	".test.ts",    // TypeScript
	".spec.ts",    // TypeScript
	".test.tsx",   // React TypeScript
	".spec.tsx",   // React TypeScript
	".test.py",    // Python
	".spec.py",    // Python
	"Test.java",   // Java (JUnit style)
	"Tests.java",  // Java (alternative)
	"IT.java",     // Java Integration Tests
	".test.php",   // PHP
	".spec.php",   // PHP
	"_test.rb",    // Ruby (Test::Unit style)
	".test.cs",    // C#
	".spec.cs",    // C#
	"Test.cs",     // C# (NUnit style)
	"Tests.cs",    // C# (alternative)
	".test.cpp",   // C++
	".spec.cpp",   // C++
	".test.c",     // C
	".spec.c",     // C
	".test.rs",    // Rust
	".spec.rs",    // Rust
	".test.kt",    // Kotlin
	".spec.kt",    // Kotlin
	"Test.kt",     // Kotlin (JUnit style)
	".test.swift", // Swift
	".spec.swift", // Swift
	"Test.swift",  // Swift (XCTest style)
	".t.sol",      // Solidity (Foundry)
	".test.sol",   // Solidity (alternative)
	".spec.sol",   // Solidity (alternative)
}