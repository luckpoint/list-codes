package utils

// GetHelpMessages returns localized help messages
func GetHelpMessages() (string, string, string) {
	if IsJapanese() {
		return getJapaneseHelpMessages()
	}
	return getEnglishHelpMessages()
}

func getEnglishHelpMessages() (string, string, string) {
	short := "Summarizes a project's structure and source code into a Markdown file."

	long := `list-codes is a CLI tool that scans a specified project folder and generates a Markdown summary including:

- Project directory structure
- Detected programming languages and frameworks
- Collected dependency and configuration files (in debug mode)
- Collected source code files
- Collected README files (if --readme-only is specified)

The --prompt option allows you to prepend analysis prompts for LLM processing:

Available prompt templates:
  explain       - Project overview and architecture explanation
  find-bugs     - Bug detection and code issue analysis
  refactor      - Refactoring opportunities and code improvement
  security      - Security vulnerability assessment
  optimize      - Performance optimization suggestions
  test          - Test coverage and testing strategy improvements
  document      - Documentation improvement recommendations
  migrate       - Technology stack migration proposals
  scale         - Scalability analysis and recommendations
  maintain      - Code maintainability assessment
  api-design    - API design review and improvements
  patterns      - Design pattern application suggestions
  review        - Comprehensive code review
  architecture  - Architecture analysis and proposals
  deploy        - Deployment and DevOps improvements

Example usage:
  list-codes --folder ./my-project --output summary.md --debug
  list-codes --readme-only
  list-codes --exclude node_modules,vendor
  list-codes --prompt explain
  list-codes --prompt security
  list-codes --prompt refactor`

	completion := `To load completions:

Bash:

  $ source <(list-codes completion bash)

  # To load completions for each session, add this to your ~/.bashrc
  # or ~/.profile file:
  source <(list-codes completion bash)

Zsh:

  # If shell completion is not already enabled in your environment, you will need
  # to enable it.  You can execute the following once:

  $ echo "autoload -Uz compinit" >> ~/.zshrc
  $ echo "compinit" >> ~/.zshrc

  # To load completions for each session, add this to your ~/.zshrc file:
  source <(list-codes completion zsh)

Fish:

  $ list-codes completion fish | source

  # To load completions for each session, add this to your ~/.config/fish/completions/list-codes.fish file:
  list-codes completion fish > ~/.config/fish/completions/list-codes.fish

PowerShell:

  PS> list-codes completion powershell | Out-String | Invoke-Expression

  # To load completions for each session, add this to your PowerShell profile.
  # For example, if your profile is at:
  # $HOME\Documents\PowerShell\Microsoft.PowerShell_profile.ps1
  # then add the following line to that file:
  # list-codes completion powershell | Out-String | Invoke-Expression`

	return short, long, completion
}

func getJapaneseHelpMessages() (string, string, string) {
	short := "プロジェクトの構造とソースコードをMarkdownファイルに要約します。"

	long := `list-codesは、指定されたプロジェクトフォルダをスキャンし、以下を含むMarkdownファイルを生成するCLIツールです：

- プロジェクトディレクトリ構造
- 検出されたプログラミング言語とフレームワーク
- 収集された依存関係と設定ファイル（デバッグモード時）
- 収集されたソースコードファイル
- 収集されたREADMEファイル（--readme-onlyが指定された場合）

--promptオプションを使用すると、LLM処理用の分析プロンプトを先頭に追加できます：

利用可能なプロンプトテンプレート:
  explain       - プロジェクト概要とアーキテクチャ説明
  find-bugs     - バグ検出とコード問題分析
  refactor      - リファクタリング機会とコード改善
  security      - セキュリティ脆弱性評価
  optimize      - パフォーマンス最適化提案
  test          - テストカバレッジとテスト戦略改善
  document      - ドキュメント改善推奨事項
  migrate       - 技術スタック移行提案
  scale         - スケーラビリティ分析と推奨事項
  maintain      - コード保守性評価
  api-design    - API設計レビューと改善
  patterns      - デザインパターン適用提案
  review        - 包括的コードレビュー
  architecture  - アーキテクチャ分析と提案
  deploy        - デプロイメントとDevOps改善

使用例:
  list-codes --folder ./my-project --output summary.md --debug
  list-codes --readme-only
  list-codes --exclude node_modules,vendor
  list-codes --prompt explain
  list-codes --prompt security
  list-codes --prompt refactor`

	completion := `補完を読み込むには:

Bash:

  $ source <(list-codes completion bash)

  # 各セッションで補完を読み込むには、~/.bashrcまたは~/.profileファイルに以下を追加:
  source <(list-codes completion bash)

Zsh:

  # シェル補完が環境で有効でない場合は、有効にする必要があります。
  # 以下を一度実行してください:

  $ echo "autoload -Uz compinit" >> ~/.zshrc
  $ echo "compinit" >> ~/.zshrc

  # 各セッションで補完を読み込むには、~/.zshrcファイルに以下を追加:
  source <(list-codes completion zsh)

Fish:

  $ list-codes completion fish | source

  # 各セッションで補完を読み込むには、~/.config/fish/completions/list-codes.fishファイルに以下を追加:
  list-codes completion fish > ~/.config/fish/completions/list-codes.fish

PowerShell:

  PS> list-codes completion powershell | Out-String | Invoke-Expression

  # 各セッションで補完を読み込むには、PowerShellプロファイルに追加してください。
  # 例えば、プロファイルが以下の場所にある場合:
  # $HOME\Documents\PowerShell\Microsoft.PowerShell_profile.ps1
  # そのファイルに以下の行を追加:
  # list-codes completion powershell | Out-String | Invoke-Expression`

	return short, long, completion
}

