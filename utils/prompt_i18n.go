package utils

import (
	"fmt"
	"os"
	"strings"
)

// PromptTemplatesJA contains Japanese versions of prompt templates
var PromptTemplatesJA = map[string]string{
	"explain": `以下のプロジェクトのコードベースを分析し、以下の項目について説明してください：

1. **プロジェクトの目的と主要機能**
2. **アーキテクチャと設計パターン**
3. **主要なコンポーネントとその役割**
4. **技術スタックと依存関係**
5. **コードの構造と組織化**

技術的でない人にも理解できるよう、簡潔で分かりやすい説明をお願いします。`,

	"find-bugs": `以下のコードベースを詳細に分析し、潜在的なバグ、エラー、問題を特定してください：

1. **論理エラーや実装上の問題**
2. **メモリリークや性能問題**
3. **エラーハンドリングの不備**
4. **型安全性の問題**
5. **境界値やエッジケースの処理不備**
6. **競合状態やスレッドセーフティの問題**

各問題について、該当箇所、問題の詳細、修正案を提示してください。`,

	"refactor": `以下のコードベースを分析し、リファクタリングの機会を特定してください：

1. **コードの重複排除**
2. **関数やクラスの責任分離**
3. **命名の改善**
4. **複雑すぎる関数の分割**
5. **デザインパターンの適用機会**
6. **パフォーマンス改善**

各提案について、現在の問題点、改善後の利点、実装の優先度を示してください。`,

	"security": `以下のコードベースをセキュリティの観点から分析し、脆弱性やセキュリティリスクを特定してください：

1. **入力値検証の不備**
2. **SQLインジェクションやXSSの可能性**
3. **認証・認可の問題**
4. **機密情報の露出**
5. **不適切な権限設定**
6. **暗号化やハッシュ化の不備**

各問題について、リスクレベル、影響範囲、対策方法を提示してください。`,

	"optimize": `以下のコードベースを分析し、パフォーマンス最適化の機会を特定してください：

1. **計算量の改善可能性**
2. **メモリ使用量の最適化**
3. **I/O操作の効率化**
4. **キャッシュの活用機会**
5. **並列処理の導入可能性**
6. **ボトルネックとなりうる箇所**

各最適化について、現在の問題、改善案、期待される効果を示してください。`,

	"test": `以下のコードベースを分析し、テストの改善提案を行ってください：

1. **テストカバレッジの不足箇所**
2. **エッジケースのテスト不足**
3. **統合テストの必要性**
4. **テストコードの品質改善**
5. **モックやスタブの活用機会**
6. **テストの自動化改善**

具体的なテストケースの例も合わせて提示してください。`,

	"document": `以下のコードベースを分析し、ドキュメントの改善提案を行ってください：

1. **不足している仕様書や設計書**
2. **コメントが不足している複雑な処理**
3. **API仕様書の充実度**
4. **README.mdの改善点**
5. **コード内ドキュメントの品質**
6. **使用例やチュートリアルの必要性**

ユーザーと開発者の両方の視点から提案してください。`,

	"migrate": `以下のコードベースを分析し、技術スタックの移行や更新の提案を行ってください：

1. **古いライブラリやフレームワークの更新**
2. **より適切な技術選択の提案**
3. **言語バージョンのアップデート**
4. **アーキテクチャの現代化**
5. **移行に伴うリスクと利点**
6. **段階的な移行計画**

移行の複雑さとメリットを評価してください。`,

	"scale": `以下のコードベースをスケーラビリティの観点から分析してください：

1. **現在のアーキテクチャの限界**
2. **ボトルネックとなりうる箇所**
3. **水平・垂直スケーリングの対応**
4. **データベース設計の拡張性**
5. **マイクロサービス化の可能性**
6. **負荷分散の仕組み**

大規模運用に向けた改善提案を行ってください。`,

	"maintain": `以下のコードベースの保守性を分析し、改善提案を行ってください：

1. **コードの可読性向上**
2. **モジュール間の依存関係整理**
3. **設定管理の改善**
4. **ログとエラー処理の標準化**
5. **開発ワークフローの改善**
6. **技術的負債の特定と対処**

長期的な保守性向上の観点から提案してください。`,

	"api-design": `以下のコードベースのAPI設計を分析し、改善提案を行ってください：

1. **RESTful設計の妥当性**
2. **エンドポイント設計の一貫性**
3. **リクエスト・レスポンス形式**
4. **エラーレスポンスの標準化**
5. **バージョニング戦略**
6. **APIドキュメントの充実度**

使いやすく一貫性のあるAPI設計を目指した提案をしてください。`,

	"patterns": `以下のコードベースを分析し、適用可能なデザインパターンを提案してください：

1. **現在のコード構造の問題点**
2. **適用可能なGoFパターン**
3. **アーキテクチャパターンの活用**
4. **関数型プログラミングパターン**
5. **並行処理パターン**
6. **エラーハンドリングパターン**

各パターンの適用箇所と期待される効果を示してください。`,

	"review": `以下のコードベースに対して包括的なコードレビューを実施してください：

1. **コーディング規約の遵守状況**
2. **ベストプラクティスの適用度**
3. **コードの品質と一貫性**
4. **潜在的な改善点**
5. **チーム開発での課題**
6. **学習すべき技術や手法**

建設的で具体的なフィードバックを提供してください。`,

	"architecture": `以下のコードベースのアーキテクチャを分析し、評価・提案を行ってください：

1. **現在のアーキテクチャの特徴と評価**
2. **レイヤード構造の妥当性**
3. **依存関係の方向性**
4. **モジュール分割の適切性**
5. **設計原則（SOLID等）の適用状況**
6. **将来の拡張性への対応**

より良いアーキテクチャに向けた具体的な提案をしてください。`,

	"deploy": `以下のコードベースのデプロイメントと運用面を分析し、改善提案を行ってください：

1. **CI/CDパイプラインの改善**
2. **デプロイメント戦略の最適化**
3. **監視とログの充実**
4. **障害対応とロールバック**
5. **環境管理の改善**
6. **運用自動化の機会**

DevOpsの観点から実践的な改善案を提示してください。`,
}

// PromptTemplatesEN contains English versions of prompt templates
var PromptTemplatesEN = map[string]string{
	"explain": `Please analyze the following codebase and explain the following aspects:

1. **Project purpose and main functionality**
2. **Architecture and design patterns**
3. **Key components and their roles**
4. **Technology stack and dependencies**
5. **Code structure and organization**

Please provide a clear and concise explanation that can be understood by non-technical people.`,

	"find-bugs": `Please analyze the following codebase in detail and identify potential bugs, errors, and issues:

1. **Logic errors and implementation problems**
2. **Memory leaks and performance issues**
3. **Error handling deficiencies**
4. **Type safety issues**
5. **Boundary value and edge case handling problems**
6. **Race conditions and thread safety issues**

For each issue, please provide the location, detailed description, and suggested fixes.`,

	"refactor": `Please analyze the following codebase and identify refactoring opportunities:

1. **Code duplication elimination**
2. **Function and class responsibility separation**
3. **Naming improvements**
4. **Breaking down overly complex functions**
5. **Design pattern application opportunities**
6. **Performance improvements**

For each suggestion, please indicate current problems, benefits after improvement, and implementation priority.`,

	"security": `Please analyze the following codebase from a security perspective and identify vulnerabilities and security risks:

1. **Input validation deficiencies**
2. **SQL injection and XSS possibilities**
3. **Authentication and authorization issues**
4. **Sensitive information exposure**
5. **Inappropriate permission settings**
6. **Encryption and hashing deficiencies**

For each issue, please provide risk level, impact scope, and countermeasures.`,

	"optimize": `Please analyze the following codebase and identify performance optimization opportunities:

1. **Computational complexity improvement possibilities**
2. **Memory usage optimization**
3. **I/O operation efficiency**
4. **Caching utilization opportunities**
5. **Parallel processing introduction possibilities**
6. **Potential bottlenecks**

For each optimization, please indicate current problems, improvement proposals, and expected effects.`,

	"test": `Please analyze the following codebase and provide testing improvement suggestions:

1. **Areas lacking test coverage**
2. **Missing edge case tests**
3. **Need for integration tests**
4. **Test code quality improvements**
5. **Mock and stub utilization opportunities**
6. **Test automation improvements**

Please also provide specific test case examples.`,

	"document": `Please analyze the following codebase and provide documentation improvement suggestions:

1. **Missing specifications and design documents**
2. **Complex processes lacking comments**
3. **API documentation completeness**
4. **README.md improvement points**
5. **Code documentation quality**
6. **Need for usage examples and tutorials**

Please provide suggestions from both user and developer perspectives.`,

	"migrate": `Please analyze the following codebase and provide technology stack migration and update suggestions:

1. **Updating old libraries and frameworks**
2. **More appropriate technology choice suggestions**
3. **Language version updates**
4. **Architecture modernization**
5. **Migration risks and benefits**
6. **Gradual migration plans**

Please evaluate migration complexity and benefits.`,

	"scale": `Please analyze the following codebase from a scalability perspective:

1. **Current architecture limitations**
2. **Potential bottlenecks**
3. **Horizontal and vertical scaling support**
4. **Database design scalability**
5. **Microservices architecture possibilities**
6. **Load balancing mechanisms**

Please provide improvement suggestions for large-scale operations.`,

	"maintain": `Please analyze the maintainability of the following codebase and provide improvement suggestions:

1. **Code readability improvements**
2. **Inter-module dependency organization**
3. **Configuration management improvements**
4. **Log and error handling standardization**
5. **Development workflow improvements**
6. **Technical debt identification and resolution**

Please provide suggestions from a long-term maintainability perspective.`,

	"api-design": `Please analyze the API design of the following codebase and provide improvement suggestions:

1. **RESTful design validity**
2. **Endpoint design consistency**
3. **Request and response formats**
4. **Error response standardization**
5. **Versioning strategy**
6. **API documentation completeness**

Please provide suggestions aimed at user-friendly and consistent API design.`,

	"patterns": `Please analyze the following codebase and suggest applicable design patterns:

1. **Current code structure issues**
2. **Applicable GoF patterns**
3. **Architecture pattern utilization**
4. **Functional programming patterns**
5. **Concurrency patterns**
6. **Error handling patterns**

Please indicate application locations and expected effects for each pattern.`,

	"review": `Please conduct a comprehensive code review of the following codebase:

1. **Coding standard compliance**
2. **Best practice application**
3. **Code quality and consistency**
4. **Potential improvement points**
5. **Team development challenges**
6. **Technologies and methods to learn**

Please provide constructive and specific feedback.`,

	"architecture": `Please analyze the architecture of the following codebase and provide evaluation and suggestions:

1. **Current architecture characteristics and evaluation**
2. **Layered structure validity**
3. **Dependency directions**
4. **Module division appropriateness**
5. **Design principle (SOLID, etc.) application status**
6. **Future extensibility support**

Please provide specific suggestions for better architecture.`,

	"deploy": `Please analyze the deployment and operational aspects of the following codebase and provide improvement suggestions:

1. **CI/CD pipeline improvements**
2. **Deployment strategy optimization**
3. **Monitoring and logging enhancements**
4. **Incident response and rollback**
5. **Environment management improvements**
6. **Operational automation opportunities**

Please provide practical improvement proposals from a DevOps perspective.`,
}

// GetPromptI18n returns the appropriate prompt template based on current language
func GetPromptI18n(promptParam string, debugMode bool) (string, error) {
	if promptParam == "" {
		return "", nil
	}

	// Check if it's a predefined template
	var templates map[string]string
	isJa := IsJapanese()
	PrintDebug(fmt.Sprintf("Language check: IsJapanese=%v, currentLang=%v", isJa, GetCurrentLanguage()), debugMode)
	if isJa {
		templates = PromptTemplatesJA
		PrintDebug("Using Japanese prompt templates", debugMode)
	} else {
		templates = PromptTemplatesEN
		PrintDebug("Using English prompt templates", debugMode)
	}
	
	if template, exists := templates[promptParam]; exists {
		PrintDebug(fmt.Sprintf("Using predefined prompt template: %s", promptParam), debugMode)
		return template, nil
	}

	// Check if it's a file path
	if _, err := os.Stat(promptParam); err == nil {
		PrintDebug(fmt.Sprintf("Reading prompt from file: %s", promptParam), debugMode)
		content, err := os.ReadFile(promptParam)
		if err != nil {
			return "", fmt.Errorf("failed to read prompt file '%s': %v", promptParam, err)
		}
		return strings.TrimSpace(string(content)), nil
	}

	// Treat as custom prompt text
	PrintDebug("Using custom prompt text", debugMode)
	return promptParam, nil
}