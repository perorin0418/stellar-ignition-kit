# Agent Handoff Schema

このファイルは、`.github/agents/*.agent.md` で定義したカスタムエージェント同士が受け渡すデータ仕様を、人間向けに一覧化した共通リファレンスです。

## 目的

- エージェント間の受け渡しデータを YAML で明確化する
- `task` / `scope` / `context.prior_outputs` / `result` の意味を統一する
- オーケストレーション時の「何を渡すべきか」「何が返るべきか」を曖昧にしない
- `status` と `validation.verdict` の扱いを統一する

## 適用範囲

- `.github/agents/orchestrator.agent.md`
- `.github/agents/galactic-director.agent.md`
- `.github/agents/*.agent.md` の各作業エージェント

## 基本ルール

1. エージェントへの委譲入力は **YAML 1 文書** で渡す。
2. エージェントからの返却も **YAML 1 文書** のみとする。
3. YAML の前後に自然言語の前置き・補足・箇条書きを付けない。
4. 不明点は推測で埋めず、`open_issues` または各 `result` 内の未確定項目に残す。
5. `request_id` は同一ワークフロー内で原則引き継ぐ。
6. タスクが枝分かれする場合のみ、`request_id` にサフィックスを付与して分岐を識別する。
7. `context.prior_outputs` には、先行エージェントの **返却 YAML を構造のまま** 格納する。
8. `status: needs_input` / `blocked` / `failed` の返却を受けた場合は、次工程へ進めずに停止・補足・確認へ戻る。

## 共通入力スキーマ

すべてのエージェントへの入力は、原則として次の構造を基準にします。

```yaml
schema_version: "1.0"
request_id: "REQ-20260426-001"
source_agent: "Stellar Orchestrator"
target_agent: "Code Researcher"
task:
  summary: "依頼の要約"
  goal: "達成したい状態"
  background: "背景・理由"
scope:
  include:
    - ".github/agents/**"
  exclude: []
constraints:
  - "要求範囲外の変更をしない"
acceptance_criteria:
  - "候補変更ファイルと主要シンボルが明示される"
context:
  repo_root: "d:/GitWorkspace/github/stellar-ignition-kit"
  relevant_files:
    - ".github/agents/orchestrator.agent.md"
  relevant_symbols: []
  relevant_documents:
    - ".github/agents/agent-handoff-schema.md"
  prior_outputs: []
  changed_files: []
response_requirements:
  format: "yaml"
  required_sections:
    - summary
    - result
    - artifacts
    - validation
    - open_issues
```

### 共通入力フィールド定義

| フィールド | 型 | 必須 | 説明 |
| --- | --- | --- | --- |
| `schema_version` | string | 必須 | スキーマのバージョン。現時点では `"1.0"` 固定。 |
| `request_id` | string | 必須 | 同一ワークフローを追跡する識別子。 |
| `source_agent` | string | 必須 | 委譲元のエージェント名。 |
| `target_agent` | string | 必須 | 委譲先のエージェント名。 |
| `task.summary` | string | 必須 | 依頼の短い要約。 |
| `task.goal` | string | 必須 | 完了時に満たしたい状態。 |
| `task.background` | string | 推奨 | 背景、理由、前提。 |
| `scope.include` | string[] | 推奨 | 対象範囲。ファイル、フォルダ、ドキュメント分類など。 |
| `scope.exclude` | string[] | 推奨 | 非対象範囲。 |
| `constraints` | string[] | 推奨 | 禁止事項、守るべきルール、技術制約。 |
| `acceptance_criteria` | string[] | 推奨 | そのエージェントの成功条件。 |
| `context.repo_root` | string | 推奨 | リポジトリルート。 |
| `context.relevant_files` | string[] | 任意 | 直接読むべき主要ファイル。 |
| `context.relevant_symbols` | string[] | 任意 | 主にコード調査用。 |
| `context.relevant_documents` | string[] | 任意 | 主に文書調査・文書更新用。 |
| `context.prior_outputs` | object[] | 任意 | 先行エージェントの返却 YAML をそのまま格納する。 |
| `context.changed_files` | string[] | 任意 | 主に上位エージェントが変更結果を引き渡すときに使う。 |
| `response_requirements.format` | string | 必須 | 現時点では `"yaml"` を指定する。 |
| `response_requirements.required_sections` | string[] | 推奨 | 返却に必須なセクションの一覧。 |

## 共通返却スキーマ

すべてのエージェントは、原則として次の構造で返却します。

```yaml
schema_version: "1.0"
request_id: "REQ-20260426-001"
agent: "Code Researcher"
status: "ok"
summary: "返却内容の要約"
result: {}
artifacts:
  reviewed_files: []
  changed_files: []
  commands_run: []
validation:
  verdict: "not_applicable"
  checks: []
open_issues: []
next_actions: []
```

### 共通返却フィールド定義

| フィールド | 型 | 必須 | 説明 |
| --- | --- | --- | --- |
| `schema_version` | string | 必須 | スキーマのバージョン。 |
| `request_id` | string | 必須 | 委譲入力から引き継ぐ識別子。 |
| `agent` | string | 必須 | 返却エージェント名。 |
| `status` | string | 必須 | 実行状態。`ok` / `needs_input` / `blocked` / `failed`。 |
| `summary` | string | 必須 | 返却要約。 |
| `result` | object | 必須 | エージェント固有の主結果。 |
| `artifacts.reviewed_files` | string[] | 必須 | 読んだファイル。 |
| `artifacts.changed_files` | string[] | 必須 | 更新したファイル。変更がない場合は空配列。 |
| `artifacts.commands_run` | string[] | 必須 | 実行したコマンド。未実行なら空配列。 |
| `validation.verdict` | string | 必須 | 検証の総合判定。 |
| `validation.checks` | object[] | 必須 | 実施したチェック一覧。 |
| `open_issues` | array | 必須 | 未解決事項。 |
| `next_actions` | array | 必須 | 次に取るべき候補アクション。 |

## `status` の意味

| 値 | 意味 | 次のアクション |
| --- | --- | --- |
| `ok` | 依頼を処理できた | 次工程へ進めてよい |
| `needs_input` | 入力不足で判断できない | 追加調査またはユーザー確認へ戻す |
| `blocked` | 権限・前提・ゲート不足で停止 | ブロック要因を解消するまで進めない |
| `failed` | 処理自体が失敗した | エラー内容を確認し、再試行または方針変更 |

## `validation.verdict` の意味

| 値 | 意味 |
| --- | --- |
| `passed` | 検証上は問題なし |
| `warning` | 進行可能だが留意点あり |
| `failed` | 検証で不適合 |
| `not_run` | 検証未実施 |
| `not_applicable` | そのエージェントでは検証が本質でない |

## `prior_outputs` の格納ルール

`context.prior_outputs` には、先行エージェントの返却 YAML を構造ごと格納します。要約だけで済ませず、可能な限り返却全体を引き渡します。

```yaml
context:
  prior_outputs:
    - schema_version: "1.0"
      request_id: "REQ-20260426-001"
      agent: "Intent Analyzer"
      status: "ok"
      summary: "要件整理の要約"
      result:
        requirements:
          - id: "REQ-1"
            statement: "YAML で受け渡し契約を統一する"
            rationale: "エージェント間の受け渡しを安定化するため"
        constraints:
          - type: "policy"
            statement: "要求範囲外の変更をしない"
        unknowns: []
        questions: []
      artifacts:
        reviewed_files: []
        changed_files: []
        commands_run: []
      validation:
        verdict: "not_applicable"
        checks: []
      open_issues: []
      next_actions: []
```

## エージェント別 `result` 契約

以下は、各エージェントが `result` に最低限含めるべき構造です。

### 1. `Intent Analyzer`

```yaml
result:
  requirements:
    - id: "REQ-1"
      statement: "実現すべき要件"
      rationale: "その要件が必要な理由"
  constraints:
    - type: "policy"
      statement: "守るべき制約"
  unknowns:
    - item: "不足情報"
      impact: "判断にどう影響するか"
  questions:
    - id: "Q-1"
      question: "確認したい内容"
      reason: "確認が必要な理由"
```

### 2. `Document Researcher`

```yaml
result:
  references:
    - path: "doc/.../sample.md"
      focus: "確認した観点"
      evidence: "根拠となる記述の要約"
  key_points:
    - statement: "重要ポイント"
      evidence: "根拠"
      implementation_impact: "実装や判断への影響"
  prohibitions:
    - "禁止事項"
  maintainability_constraints:
    - "長期保守上の制約"
```

### 3. `Code Researcher`

```yaml
result:
  candidate_files:
    - path: "src/example.ts"
      reason: "変更候補である根拠"
  candidate_symbols:
    - file: "src/example.ts"
      symbol: "ExampleService.execute"
      reason: "主要シンボルである根拠"
  implementation_options:
    - option: "既存サービスを拡張する"
      changes:
        - "変更案の要点"
      pros:
        - "利点"
      cons:
        - "欠点"
  risks:
    confirmed:
      - "確実に言えるリスク"
    likely:
      - "強く疑われるリスク"
    possible:
      - "条件次第のリスク"
```

### 4. `Document Updater`

```yaml
result:
  changes:
    - path: "doc/.../target.md"
      summary: "更新内容の要約"
      editable_blocks_only: true
  documentation_decisions:
    - decision: "採用した記述方針"
      rationale: "採用理由"
  unresolved_tbd:
    - topic: "未確定事項"
      next_step: "次に必要な確認"
```

### 5. `Document Guideline Checker`

```yaml
result:
  target_files:
    - "doc/.../target.md"
  checks:
    - name: "テンプレート準拠"
      result: "passed"
      evidence: "根拠の要約"
  violations:
    - severity: "warning"
      message: "懸念点の説明"
      path: "doc/.../target.md"
```

### 6. `Code Updater`

```yaml
result:
  changes:
    - path: "src/example.ts"
      symbol: "ExampleService.execute"
      summary: "変更内容の要約"
  implementation_notes:
    - note: "実装上の補足"
```

### 7. `Code Guideline Checker`

```yaml
result:
  target_files:
    - "src/example.ts"
  checks:
    - name: "命名・責務分離確認"
      result: "passed"
      evidence: "根拠の要約"
  violations:
    - severity: "warning"
      message: "改善余地の説明"
      path: "src/example.ts"
```

### 8. `Maintainability Checker`

```yaml
result:
  assessment:
    - area: "responsibility_split"
      verdict: "pass"
      evidence: "根拠の要約"
  follow_ups:
    - priority: "medium"
      action: "将来対応が望ましい事項"
```

### 9. `Stellar Orchestrator`

```yaml
result:
  workflow:
    completed_steps:
      - "Document Researcher"
      - "Code Researcher"
    skipped_steps: []
  aggregated_findings:
    document_findings:
      - "文書調査の要点"
    code_findings:
      - "コード調査の要点"
    decisions:
      - "採用した判断"
  changes:
    documentation:
      - "更新したドキュメント"
    source_code:
      - "更新したコード"
```

### 10. `Galactic Director`

```yaml
result:
  task_breakdown:
    decision: "split"
    rationale: "分割した理由"
    subtasks:
      - id: "SUB-1"
        target_agent: "Stellar Orchestrator"
        scope:
          - ".github/agents/**"
  orchestration_results:
    - request_id: "REQ-20260426-001-A"
      summary: "各オーケストレーション結果の要約"
      validation_verdict: "passed"
```

## エージェント別の受け渡しマトリクス

| 委譲先 | 入力時に特に重要な項目 | `prior_outputs` の最低期待値 | 主な返却 `result` |
| --- | --- | --- | --- |
| `Intent Analyzer` | `task`, `constraints`, `acceptance_criteria` | なし | `requirements`, `constraints`, `unknowns`, `questions` |
| `Document Researcher` | `scope.include`, `context.relevant_documents` | あれば `Intent Analyzer` | `references`, `key_points`, `prohibitions` |
| `Code Researcher` | `scope.include`, `context.relevant_files`, `context.relevant_symbols` | あれば `Document Researcher` | `candidate_files`, `candidate_symbols`, `implementation_options`, `risks` |
| `Document Updater` | `task`, `scope`, `constraints` | `Intent Analyzer`, `Document Researcher` | `changes`, `documentation_decisions`, `unresolved_tbd` |
| `Document Guideline Checker` | `context.changed_files`, `scope.include` | `Document Updater` | `target_files`, `checks`, `violations` |
| `Code Updater` | `task`, `scope`, `constraints` | `Intent Analyzer`, `Code Researcher`、必要時 `Document Updater` | `changes`, `implementation_notes` |
| `Code Guideline Checker` | `context.changed_files`, `scope.include` | `Code Updater` | `target_files`, `checks`, `violations` |
| `Maintainability Checker` | `task`, `constraints`, `acceptance_criteria`, `context.changed_files` | 関連する全エージェント返却 | `assessment`, `follow_ups` |
| `Stellar Orchestrator` | `task`, `scope`, `constraints`, `acceptance_criteria` | あれば `Intent Analyzer` など | `workflow`, `aggregated_findings`, `changes` |
| `Galactic Director` | `task`, `scope`, `constraints`, `acceptance_criteria` | あれば `Intent Analyzer` / `Stellar Orchestrator` | `task_breakdown`, `orchestration_results` |

## 推奨受け渡し順序

通常の単一タスクでは、以下の順序を基準にします。

1. `Stellar Orchestrator` → `Intent Analyzer`（必要時）
2. `Stellar Orchestrator` → `Document Researcher`
3. `Stellar Orchestrator` → `Code Researcher`
4. `Stellar Orchestrator` → `Document Updater`（仕様変更がある場合）
5. `Stellar Orchestrator` → `Document Guideline Checker`（文書更新がある場合）
6. `Stellar Orchestrator` → `Code Updater`（コード変更がある場合）
7. `Stellar Orchestrator` → `Code Guideline Checker`（コード更新がある場合）
8. `Stellar Orchestrator` → `Maintainability Checker`（変更がある場合）

マルチタスク時は、`Galactic Director` が上位で以下を担います。

1. 分割要否の判断
2. 分割時の `request_id` ブランチ管理
3. 各 `Stellar Orchestrator` への YAML 委譲
4. 返却 YAML の収集と統合

## 運用ルール

- `Document Updater` の前に、少なくとも要件整理結果か文書調査結果があることが望ましい。
- `Code Updater` の前に、少なくともコード調査結果があること。
- 仕様変更を伴う場合は、原則として `Document Updater` を `Code Updater` より先に行う。
- `Maintainability Checker` には、可能な限り全先行エージェントの返却 YAML を `prior_outputs` として渡す。
- `warning` 判定を受けた場合も、その扱いを最終報告で明示する。

## 最小ハンドオフ例

### `Stellar Orchestrator` → `Document Researcher`

```yaml
schema_version: "1.0"
request_id: "REQ-20260426-001"
source_agent: "Stellar Orchestrator"
target_agent: "Document Researcher"
task:
  summary: "エージェント間受け渡し仕様を文書から確認したい"
  goal: "関連ドキュメントと制約を根拠付きで整理する"
  background: "YAML の受け渡し契約を安定化したい"
scope:
  include:
    - ".github/agents/**"
  exclude: []
constraints:
  - "要求範囲外の変更はしない"
acceptance_criteria:
  - "参照文書と重要ポイントが整理される"
context:
  repo_root: "d:/GitWorkspace/github/stellar-ignition-kit"
  relevant_files:
    - ".github/agents/orchestrator.agent.md"
  relevant_symbols: []
  relevant_documents:
    - ".github/copilot-instructions.md"
    - ".github/agents/orchestrator.agent.md"
  prior_outputs: []
  changed_files: []
response_requirements:
  format: "yaml"
  required_sections:
    - summary
    - result.references
    - result.key_points
    - validation
    - open_issues
```

### `Document Researcher` の返却

```yaml
schema_version: "1.0"
request_id: "REQ-20260426-001"
agent: "Document Researcher"
status: "ok"
summary: "関連ドキュメントと実装制約を整理した"
result:
  references:
    - path: ".github/agents/orchestrator.agent.md"
      focus: "委譲契約と返却契約"
      evidence: "YAML 1 文書で受け渡すことが明記されている"
  key_points:
    - statement: "返却は YAML のみで統一する必要がある"
      evidence: "各エージェント定義の出力契約"
      implementation_impact: "自由文ではなく構造化データを扱うべき"
  prohibitions:
    - "自然言語だけで委譲しない"
  maintainability_constraints:
    - "先行結果は `prior_outputs` で引き継ぐ"
artifacts:
  reviewed_files:
    - ".github/agents/orchestrator.agent.md"
  changed_files: []
  commands_run: []
validation:
  verdict: "not_applicable"
  checks: []
open_issues: []
next_actions:
  - "必要に応じて Code Researcher へ委譲する"
```

## 補足

- このファイルは人間向けの総覧であり、実際のエージェント実行時には `.agent.md` 側の指示も同時に満たすこと。
- もしこのファイルと各 `.agent.md` の内容が食い違った場合は、**両方を同時に更新して同期を取る**。
- 将来 `schema_version` を更新する場合は、このファイルと各 `.agent.md` のサンプルを同時に更新する。
