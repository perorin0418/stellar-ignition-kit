# doc/README

このディレクトリは、SDD（Software Design Description）を **要件定義 → 設計** の順で一貫して進めるための基点です。

## 1. まず最初に見るファイル

1. `requirements/system_overviews/system_overview.md`
2. `requirements/nonfunctional_requirements/nonfunctional_requirements.md`
3. `requirements/glossaries/glossary.md`
4. `requirements/business_codes/business_codes.md`

上記4点で「背景・制約・用語・業務コード」を固定してから、個別機能へ進んでください。

## 2. SDDの進め方（推奨導線）

### Phase A: 要件定義（What）

- 業務フロー整理: `requirements/business_flows/<業務コード>/`
- 機能要件定義: `requirements/functional_requirements/<業務コード>/`

成果物の単位は、業務コード配下の機能コード（例: `XX0001FR_*`）です。

### Phase B: 設計（How）

- 共通方針を先に確定
  - `design/coding_standards/`
  - `design/log_policies/`
  - `design/error_policies/`
- 機能詳細設計へ展開
  - `design/features/<業務コード>/<機能設計コード>/`

要件側の機能（`*FR_*`）と設計側の機能（`*DD_*`）を対応づけて管理します。

## 3. ディレクトリの役割（全体像）

- `requirements/`: 何を作るかを定義
  - システム概要 / 非機能 / 用語 / 業務コード / 業務フロー / 機能要件
- `design/`: どう作るかを定義
  - コーディング規約 / ログ方針 / エラー方針 / 機能詳細設計

## 4. 命名・配置ルール（運用ルール）

- 業務コード単位で `XX`, `YY`, `ZZ` などの配下に格納する。
- 機能コードは工程ごとに接尾辞を使い分ける。
  - 要件: `XXXXFR_*`
  - 設計: `XXXXDD_*`
  - フロー: `XXXXBF_*`
- 1機能につき1ディレクトリを基本とし、成果物を同階層に集約する。
