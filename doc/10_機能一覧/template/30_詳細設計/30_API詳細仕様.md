# API詳細仕様

## 1. 文書の目的
<!-- AI_READ_ONLY_START -->
- 本API仕様書で定義する対象（何を定義し、何を定義しないか）を記載する。
- 想定読者（アプリ開発者、バックエンド開発者、運用担当、QA）と利用目的を明記する。
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->

## 2. 適用範囲
<!-- AI_READ_ONLY_START -->
- 対象システム、対象環境（開発/検証/本番）、対象API群を記載する。
- 対象外のAPIや旧版仕様との関係（互換・非互換）を明記する。
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->

## 3. 共通仕様（認証・ヘッダ・エラー）
<!-- AI_READ_ONLY_START -->
- 認証方式、認可の前提、トークン有効期限・更新方針を記載する。
- 共通ヘッダ（`Content-Type`、`Authorization`、`X-Request-Id`等）の必須/任意を記載する。
- エラーコード体系、HTTPステータスとの対応、エラー時レスポンス形式を記載する。
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->

## 4. エンドポイント一覧
<!-- AI_READ_ONLY_START -->
- エンドポイントごとに、HTTPメソッド、パス、概要、対象機能、利用者を一覧化する。
- 各APIの公開/内部区分、レート制限、重要度（クリティカル度）を記載する。
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->

## 5. リクエスト仕様
<!-- AI_READ_ONLY_START -->
- パスパラメータ、クエリ、ヘッダ、ボディの項目定義（型・必須・制約）を記載する。
- バリデーションルール、デフォルト値、許容値、相関チェック条件を記載する。
- サンプルリクエストを正常系・異常系で必要に応じて記載する。
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->

## 6. レスポンス仕様
<!-- AI_READ_ONLY_START -->
- 正常時レスポンス（ステータス、項目、意味、単位、並び順）を記載する。
- 異常時レスポンス（エラーコード、メッセージ、リトライ可否、対処方法）を記載する。
- サンプルレスポンスを正常系・異常系で必要に応じて記載する。
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->

## 7. スキーマ定義
<!-- AI_READ_ONLY_START -->
- OpenAPI `components/schemas` に対応するデータ構造の意味と利用箇所を記載する。
- 各項目の論理名、型、桁数/最大長、必須性、業務上の意味を記載する。
- 互換性維持方針（追加・廃止・型変更時のルール）を記載する。
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->

## 8. 非機能・運用観点
<!-- AI_READ_ONLY_START -->
- 性能要件（応答時間、スループット）、可用性要件、タイムアウト値を記載する。
- 監視項目（成功率、エラー率、遅延）、ログ方針、アラート閾値を記載する。
- 運用手順（障害時対応、問い合わせ導線、リリース/ロールバック）を記載する。
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->

## 9. OpenAPI定義（YAML）
<!-- AI_READ_ONLY_START -->
- この章には実際のOpenAPI YAML本体を記載し、1〜8章の記述内容と整合させる。
- `paths`、`components`、`security`の更新時は、破壊的変更の有無をレビュー記録に残す。

```yaml
openapi: 3.0.3
info:
	title: API詳細仕様
	version: 0.1.0
	description: |
		本ファイルはAPI詳細仕様のテンプレート。
		各章相当の記載内容をOpenAPI要素に対応付けて定義する。

servers:
	- url: https://api.example.com
		description: 利用環境（本番・検証など）を記載する。

tags:
	- name: sample
		description: API群の業務分類と責務を記載する。

paths:
	/sample:
		get:
			tags: [sample]
			summary: APIの目的を1行で記載する。
			description: |
				- どの業務シナリオで呼ばれるAPIか
				- 前提条件と事後条件
				- 権限・利用制約
			operationId: getSample
			parameters:
				- name: requestId
					in: header
					required: false
					schema:
						type: string
					description: トレーシングや監査に使う識別子の説明を記載する。
			responses:
				'200':
					description: 正常応答時の意味と利用方法を記載する。
					content:
						application/json:
							schema:
								$ref: '#/components/schemas/SampleResponse'
				'400':
					description: 入力不正時の条件と対処方法を記載する。
				'401':
					description: 認証失敗時の条件と再試行方針を記載する。
				'500':
					description: サーバー障害時の扱いと問い合わせ先を記載する。

components:
	securitySchemes:
		bearerAuth:
			type: http
			scheme: bearer
			bearerFormat: JWT

	schemas:
		SampleResponse:
			type: object
			description: レスポンス構造の業務意味を記載する。
			properties:
				resultCode:
					type: string
					description: 結果コード体系と値の意味を記載する。
				message:
					type: string
					description: ユーザー向け/運用向けメッセージ方針を記載する。
			required:
				- resultCode

security:
	- bearerAuth: []
```
<!-- AI_READ_ONLY_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->
