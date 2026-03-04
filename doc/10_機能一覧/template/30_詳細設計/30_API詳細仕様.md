# API詳細仕様

## 1. 文書の目的
- APIの詳細な設計を明確化し、開発者や利用者が正確に利用できるようにする。

## 2. OpenAPI定義（YAML）
<!-- AI_CONTEXT_START --><!--
- この章には実際のOpenAPI YAML本体を記載する。
- `paths`、`components`、`security`の更新時は、破壊的変更の有無をレビュー記録に残す。

- サンプル
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
--><!-- AI_CONTEXT_END -->
<!-- AI_EXAMPLE_START --><!--
--><!-- AI_EXAMPLE_END -->
<!-- AI_EDITABLE_START -->
<!-- AI_EDITABLE_END -->
