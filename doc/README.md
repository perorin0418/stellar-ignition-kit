# doc/README

このディレクトリは、SDD（Software Design Description）を **要件定義 → 設計** の順で一貫して進めるための基点です。

## 1. まず最初に見るファイル

1. `10_requirements/10_system_overviews/system_overview.md`
2. `10_requirements/20_nonfunctional_requirements/nonfunctional_requirements.md`
3. `10_requirements/30_glossaries/glossary.md`
4. `10_requirements/40_business_codes/business_codes.md`

上記4点で「背景・制約・用語・業務コード」を固定してから、個別機能へ進んでください。

## 2. SDDの進め方（推奨導線）

### Phase A: 要件定義（What）

- 業務フロー整理: `10_requirements/50_business_flows/<業務コード>/`
- 機能要件定義: `10_requirements/60_functional_requirements/<業務コード>/`

成果物の単位は、業務コード配下の機能コード（例: `XX0001FR_*`）です。

### Phase B: 設計（How）

- 共通方針を先に確定
  - `20_design/10_coding_standards/`
  - `20_design/20_log_policies/`
  - `20_design/30_error_policies/`
- 機能詳細設計へ展開
  - `20_design/40_features/<業務コード>/<機能設計コード>/`

要件側の機能（`*FR_*`）と設計側の機能（`*DD_*`）を対応づけて管理します。

## 3. ディレクトリの役割（全体像）

- `10_requirements/`: 何を作るかを定義
  - システム概要 / 非機能 / 用語 / 業務コード / 業務フロー / 機能要件
- `20_design/`: どう作るかを定義
  - コーディング規約 / ログ方針 / エラー方針 / 機能詳細設計

## 4. 命名・配置ルール（運用ルール）

- 業務コード単位で `XX`, `YY`, `ZZ` などの配下に格納する。
- 機能コードは工程ごとに接尾辞を使い分ける。
  - 要件: `XXXXFR_*`
  - 設計: `XXXXDD_*`
  - フロー: `XXXXBF_*`
- 1機能につき1ディレクトリを基本とし、成果物を同階層に集約する。
