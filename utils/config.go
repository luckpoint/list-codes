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
	// Asset directories
	"assets":  {},
	"static":  {},
	"images":  {},
	"media":   {},
	"uploads": {},
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
	"HTML": {".html", ".htm", ".rhtml", ".erb", ".haml", ".slim",
		// Common JS/HTML template engines
		".hbs", ".handlebars", ".mustache", ".hjs",
		".ejs", ".pug", ".jade",
		".njk", ".nunjucks",
		".twig", ".liquid",
		".tmpl", ".tpl", ".gohtml",
		// Django/Jinja family
		".djhtml", ".jinja", ".jinja2", ".j2",
		// ASP.NET Razor
		".cshtml", ".vbhtml",
		// JSP/Java Server Pages
		".jsp", ".jspx",
		// Phoenix (Elixir)
		".eex", ".heex", ".leex",
		// Others seen in the wild
		".dust", ".eta", ".vash"},
	"CSS":          {".css", ".scss", ".sass"},
	"Lua":          {".lua"},
	"C":            {".c"},
	"C++":          {".cpp", ".hpp", ".cxx", ".hxx", ".cc", ".hh", ".h"},
	"Java":         {"build.gradle", "settings.gradle", "pom.xml"},
	"PHP":          {"composer.json"},
	"C#":           {".csproj", ".sln", "project.json"},
	"F#":           {".fsproj"},
	"VB.NET":       {".vbproj"},
	"Scala":        {"build.sbt"},
	"Clojure":      {"project.clj", "deps.edn"},
	"Haskell":      {".cabal", "stack.yaml"},
	"Erlang":       {"rebar.config"},
	"Elixir":       {"mix.exs"},
	"Dart":         {"pubspec.yaml"},
	"R":            {"DESCRIPTION"},
	"Julia":        {"Project.toml"},
	"Swift":        {"Package.swift"},
	"Kotlin":       {"build.gradle.kts"},
	"Groovy":       {"build.gradle"},
	"CMake":        {"CMakeLists.txt"},
	"Terraform":    {"main.tf"},
	"Protobuf/Buf": {".proto", "buf.yaml"},
	"Markdown":     {"README.md"},
	"MDX":          {".mdx"},
	"Vue":          {".vue"},
	"Svelte":       {".svelte"},
	"Solidity":     {"hardhat.config.js", "hardhat.config.ts", "truffle-config.js", "foundry.toml"},
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
	"HTML": {".html", ".htm", ".rhtml", ".erb", ".haml", ".slim",
		// Common JS/HTML template engines
		".hbs", ".handlebars", ".mustache", ".hjs",
		".ejs", ".pug", ".jade",
		".njk", ".nunjucks",
		".twig", ".liquid",
		".tmpl", ".tpl", ".gohtml",
		// Django/Jinja family
		".djhtml", ".jinja", ".jinja2", ".j2",
		// ASP.NET Razor
		".cshtml", ".vbhtml",
		// JSP/Java Server Pages
		".jsp", ".jspx",
		// ASP.NET Web Forms (legacy)
		".aspx", ".ascx", ".master",
		// Phoenix (Elixir)
		".eex", ".heex", ".leex",
		// Java templating
		".ftl", ".vm",
		// Others seen in the wild
		".dust", ".eta", ".vash"},
	"CSS":          {".css", ".scss", ".sass", ".less", ".styl"},
	"Lua":          {".lua"},
	"C":            {".c", ".h"},
	"C++":          {".cpp", ".hpp", ".cxx", ".hxx", ".cc", ".hh", ".c++", ".h++"},
	"Java":         {".java"},
	"PHP":          {".php", ".phtml", ".php3", ".php4", ".php5"},
	"C#":           {".cs", ".csx"},
	"F#":           {".fs", ".fsi", ".fsx"},
	"VB.NET":       {".vb"},
	"Scala":        {".scala", ".sc"},
	"Clojure":      {".clj", ".cljs", ".cljc"},
	"Haskell":      {".hs", ".lhs"},
	"Erlang":       {".erl", ".hrl"},
	"Elixir":       {".ex", ".exs"},
	"Dart":         {".dart"},
	"R":            {".r", ".R"},
	"Perl":         {".pl", ".pm", ".perl"},
	"Objective-C":  {".m", ".mm"},
	"Pascal":       {".pas", ".pp"},
	"Fortran":      {".f", ".f90", ".f95", ".f03", ".f08"},
	"COBOL":        {".cob", ".cbl"},
	"Ada":          {".ada", ".adb", ".ads"},
	"Lisp":         {".lisp", ".lsp", ".cl"},
	"Scheme":       {".scm", ".ss"},
	"Prolog":       {".pl", ".pro"},
	"MATLAB":       {".m"},
	"Octave":       {".m"},
	"Julia":        {".jl"},
	"Zig":          {".zig"},
	"D":            {".d"},
	"Nim":          {".nim"},
	"Crystal":      {".cr"},
	"Groovy":       {".groovy", ".gvy"},
	"PowerShell":   {".ps1", ".psm1", ".psd1"},
	"Bash":         {".bash"},
	"Zsh":          {".zsh"},
	"Fish":         {".fish"},
	"Tcl":          {".tcl"},
	"Vim":          {".vim"},
	"Emacs Lisp":   {".el"},
	"Protobuf/Buf": {".proto"},
	"Markdown":     {".md"},
	"MDX":          {".mdx"},
	"Vue":          {".vue"},
	"Svelte":       {".svelte"},
	"Solidity":     {".sol"},
	"YAML":         {".yaml", ".yml"},
	"JSON":         {".json"},
	"XML":          {".xml", ".xsd", ".xsl"},
	"Shell":        {".sh"},
	"Dockerfile":   {"Dockerfile"},
	"TOML":         {".toml"},
	"Gradle":       {".gradle", ".gradle.kts"},
	"INI":          {".ini", ".cfg"},
	"Properties":   {".properties"},
	"Makefile":     {"Makefile", "makefile", ".mk"},
	"CMake":        {".cmake", "CMakeLists.txt"},
	"Terraform":    {".tf", ".tfvars"},
	"HCL":          {".hcl"},
	"GraphQL":      {".graphql", ".gql"},
	"Jupyter":      {".ipynb"},
	"LaTeX":        {".tex", ".latex"},
	"Regex":        {".regex"},
	"Diff":         {".diff", ".patch"},
	"Log":          {".log"},
	"CSV":          {".csv"},
	"TSV":          {".tsv"},
	"WebAssembly":  {".wat", ".wasm"},
	"Assembly":     {".asm", ".s"},
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
	"/test/",        // Root-level "test" folder
	"/tests/",       // Root-level "tests" folder
	"/spec/",        // Root-level "spec" folder
	"/specs/",       // Root-level "specs" folder
	"/__tests__/",   // A common convention in Jest (JavaScript)
	"/__test__/",    // Alternative Jest convention
	"/testing/",     // Generic testing folder
	"/fixtures/",    // Test fixtures
	"/mocks/",       // Mock files
	"/e2e/",         // End-to-end tests
	"/integration/", // Integration tests
	"/unit/",        // Unit tests
	"test/",         // "test" folder at any depth
	"tests/",        // "tests" folder at any depth
	"spec/",         // "spec" folder at any depth
	"specs/",        // "specs" folder at any depth
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

// ASSET_EXTENSIONS contains file extensions for asset files that should be excluded by default.
// These files are typically not useful for source code analysis.
var ASSET_EXTENSIONS = map[string]struct{}{
	// Images - Raster
	".png":  {},
	".jpg":  {},
	".jpeg": {},
	".gif":  {},
	".bmp":  {},
	".ico":  {},
	".tiff": {},
	".webp": {},
	".avif": {},
	".heic": {},
	// Images - Vector
	".svg": {},
	// Fonts
	".woff":  {},
	".woff2": {},
	".ttf":   {},
	".otf":   {},
	".eot":   {},
	// Audio
	".mp3":  {},
	".wav":  {},
	".ogg":  {},
	".aac":  {},
	".flac": {},
	".m4a":  {},
	// Video
	".mp4":  {},
	".avi":  {},
	".mov":  {},
	".webm": {},
	".mkv":  {},
	".wmv":  {},
	// Archives
	".zip": {},
	".tar": {},
	".gz":  {},
	".rar": {},
	".7z":  {},
	".bz2": {},
	// Documents
	".pdf":  {},
	".docx": {},
	".xlsx": {},
	".pptx": {},
	// Executables
	".exe": {},
	".dmg": {},
	".deb": {},
	".rpm": {},
}
