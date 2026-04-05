# Python AWS CDK ディレクトリ構成・設計規約（準デファクト）

本規約は、Python を用いた AWS CDK (v2) プロジェクトにおいて、
**理解しやすさ・保守性・監査耐性・AI支援との親和性**を最大化することを目的とする。

対象読者：
- AWS CDK を業務で利用する開発者
- CloudFormation を前提知識としていること

---

## 1. 基本方針（Principles）

### P1. CloudFormation を最終成果物とみなす
- CDK は **CloudFormation を生成するための DSL** である
- 常に `cdk synth` の結果を説明できる構成にする

### P2. Stack は「デプロイ単位」、Construct は「再利用単位」
- Stack を肥大化させない
- 再利用可能なロジックは必ず Construct に切り出す

### P3. 暗黙の自動生成を減らし、明示指定を優先する
- VPC / IAM / SecurityGroup / LogGroup は可能な限り明示する
- L3 Construct の乱用を避ける

### P4. app.py は配線専用
- app.py にリソース定義を書かない
- 環境切替・Stack の組み立てのみを行う

### P5. 依存管理は pyproject.toml / uv.lock を正とする
- Python 依存関係は `pyproject.toml` で管理する
- ロック情報は `uv.lock` を正とし、`requirements.txt` は採用しない

---

## 2. 現行ディレクトリ構成（本リポジトリ準拠）

```text
infra/
├── app.py
├── cdk.context.json
├── cdk.json
├── pyproject.toml
├── README.md
├── uv.lock
│
├── cdk_stacks/                         # CloudFormation Stack 定義（デプロイ単位）
│   ├── （機能名）_stack.py              # 機能ごとにスタックを作成する
│   ├── common_stack1.py                # 共通的なリソースは別スタックに切り出す（例: VPC, IAM）
│   └── common_stack2.py                # 別リージョンや依存関係上、統合できないリソースはさらに分割する（例: HTTPS証明書）
│
└── cdk_constructs/                     # 再利用可能な Construct
    ├── （機能名）/
    │   └── （サービス名）_construct.py  # サービス単位で Construct を定義する（例: API Gateway, Lambda）
    └── common/                         # Stack 間で共有される Construct を配置する（例: VPC, IAM Role）

```

注記：
- `cdk.out/`、`.venv/`、`__pycache__/` は生成物として扱い、設計の正にはしない。
- 構成変更時は、まず本章を更新してから実装へ反映する。

---

## 3. 各ディレクトリの責務

### 3.1 app.py

**役割**：Stack の組み立てと依存関係の定義のみ

許可：
- Stack / Construct の import
- 環境情報・デプロイ定数（ドメイン名、Hosted Zone ID など）の定義
- Stack 間のリソース受け渡し（`vpc`、`api_id`、`certificate_arn` など）

禁止：
- aws_* リソースの直接定義
- if 文によるリソース分岐

追加ルール：
- `DEFAULT_ENV` は `ap-northeast-1` を基準とする。

---

### 3.2 cdk_stacks/

**重要原則：Stack はデプロイ単位、依存は app.py で明示配線する**

- 他 Stack の `Output` を内部連携の主手段として使わない
- 値共有は **同一 `app.py` 内でオブジェクト参照を受け渡す**
- CloudFormation の `Export / ImportValue` 依存を避ける

理由：
- Output は参照されると削除・変更がロックされる
- Stack 更新時の柔軟性が著しく低下する
- CDK の強み（オブジェクト参照）が失われる

**役割**：CloudFormation のデプロイ単位

ルール：
- 1 Stack = 1 責務（機能または共通基盤）
- Stack 名の追加・再編時は `app.py` の配線と本規約の両方を更新する

許可：
- Construct の組み立て

例外（必要時のみ）：
- クロスリージョン/クロスアカウントで参照が必要な場合は、SSM Parameter などの安定 ID 参照を優先する

禁止：
- 複雑なビジネスロジック
- 再利用前提の実装

---

### 3.3 cdk_constructs/

**役割**：再利用可能なインフラ部品（Stack 間共有の主軸）

- Stack 間で共有される値・リソースは Construct に閉じ込める
- Construct は Stack を跨いだ参照を前提に設計する

ルール：
- ファイル名は `*_construct.py` に統一する
- 機能配下（例：`get_cost_data/`）にサービス単位 Construct を配置する
- 共通部品は `common/` 配下へ配置し、機能別への重複実装を避ける

推奨：
- Construct ごとに README または docstring を付与

---

## 4. よくあるアンチパターン（禁止）

- 他 Stack の Output を直接参照する設計
- `Fn::ImportValue` 前提の依存関係
- Output を内部連携用に乱用する
- 巨大 Stack（1000 行超）
- L3 Construct 直書き
- app.py にリソース定義
- 環境差分を if/else で分岐
- 依存管理を `requirements.txt` と `pyproject.toml` で二重化する

---

## 5. 運用ルール

- 変更時は必ず `cdk synth` を実行し、差分を確認する
- Windows 環境では `.venv` 有効化または `uv run` で実行し、グローバル Python 依存を避ける
- PR には実施コマンドと `cdk synth` 結果（要点）を添付する
- 監査対応では CloudFormation YAML を正とする
