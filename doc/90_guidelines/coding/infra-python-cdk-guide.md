# Python AWS CDK ディレクトリ構成・設計規約（準デファクト）

本規約は、Python を用いた AWS CDK (v2) プロジェクトにおいて、
**理解しやすさ・保守性・監査耐性・AI支援との親和性**を最大化することを目的とする。

対象読者：
- AWS CDK を業務で利用する開発者
- ECS / EventBridge / API Gateway 等を含む中〜大規模構成
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

---

## 2. 推奨ディレクトリ構成（準デファクト）

```text
infra/
├── app.py
├── cdk.json
├── requirements.txt
├── cdk.out/
│
├── cdk_stacks/                # CloudFormation Stack 定義
│   ├── __init__.py
│   ├── network_stack.py
│   ├── security_stack.py
│   ├── ecs_stack.py
│   ├── batch_stack.py
│   └── api_stack.py
│
├── cdk_constructs/            # 再利用可能な Construct
│   ├── __init__.py
│   ├── network/
│   │   └── vpc.py
│   ├── ecs/
│   │   ├── cluster.py
│   │   └── task_definition.py
│   └── api/
│       └── api_gateway.py
│
├── cdk_envs/                  # 環境差分（アカウント・リージョン）
│   ├── __init__.py
│   ├── dev.py
│   ├── stg.py
│   └── prod.py
│
├── cdk_config/                # 定数・パラメータ
│   ├── __init__.py
│   ├── common.py
│   └── feature_flags.py
│
└── cdk_tests/                 # CDK Assertion Tests
    └── test_ecs_stack.py
```

---

## 3. 各ディレクトリの責務

### 3.1 app.py

**役割**：Stack の組み立てと依存関係の定義のみ

許可：
- Stack の import
- 環境情報の読み込み
- Stack 間の依存関係指定

禁止：
- aws_* リソースの直接定義
- if 文によるリソース分岐
- 設定値の直書き

```python
app = cdk.App()
env = load_env(app)

NetworkStack(app, "NetworkStack", env=env)
EcsStack(app, "EcsStack", env=env)

app.synth()
```

---

### 3.2 stacks/

**重要原則：Stack 間参照は Construct 共有を原則とする**

- 他 Stack の `Output` を参照する設計は **原則禁止** とする
- Stack 間で値を共有する場合は、**同一 `app.py` 内で Construct / リソースオブジェクトを直接受け渡す**
- CloudFormation の `Export / ImportValue` に依存しない構成を基本とする

理由：
- Output は参照されると削除・変更がロックされる
- Stack 更新時の柔軟性が著しく低下する
- CDK の強み（オブジェクト参照）が失われる

---


**役割**：CloudFormation のデプロイ単位

ルール：
- 1 Stack = 1 責務
- 原則 300 行以内
- Construct を組み合わせるだけにする

許可：
- Construct の組み立て

例外（必要時のみ）：
- クロスリージョン/クロスアカウントで参照が必要な場合は、SSM Parameter などの安定 ID 参照を優先する
- どうしても Export が必要な場合は、削除・置換がロックされることを前提に移行手順を用意する

禁止：
- 複雑なビジネスロジック
- 再利用前提の実装

---

### 3.3 constructs/

**役割**：再利用可能なインフラ部品（Stack 間共有の主軸）

- Stack 間で共有される値・リソースは Construct に閉じ込める
- Construct は Stack を跨いだ参照を前提に設計する



**役割**：再利用可能なインフラ部品

ルール：
- L2/L3 をそのまま使わず、必要に応じて分解
- IAM / SG / LogGroup を内部で明示定義
- 単体で意味を持つ粒度にする

推奨：
- Construct ごとに README または docstring を付与

---

### 3.4 envs/

**役割**：環境差分の集中管理

```python
DEV = cdk.Environment(account="111111111111", region="ap-northeast-1")
```

ルール：
- if 文による環境分岐は禁止
- app.py でどの env を使うか決定する

---

### 3.5 config/

**役割**：値と構成の分離

含めてよいもの：
- CIDR
- CPU / Memory
- バッチスケジュール

含めてはいけないもの：
- aws_* リソース
- Construct

---

## 4. Construct / Stack 設計ルール（AI可読向け）

### 4.1 Stack 間参照ルール（最重要）

#### 原則

- Stack 間で共有する値は **Output に出さない**
- **Construct または AWS リソースオブジェクトを直接渡す**

#### 追加方針（変動リソースの Export 回避）

- 変動しやすいリソース（短命・自動命名・頻繁に作り替えるもの）は Stack 間 Export を禁止
- 例: ECS TaskDefinition, AutoScalingGroup, LogGroup などの volatile なリソースは Export しない
- そうしたリソースは同一 Stack 内に配置し、依存は Construct 共有で解決する
- 共有が必要な場合は「同一 Stack への集約」を第一候補にする
- どうしても Stack 間共有が必要なら「安定 ID のみ」を Export（例：明示名の SSM Parameter、固定名の KMS Key、固定 ARN の EventBus）
- Export は削除・置換のロックを生むため、移行順序と切り戻し手順を必ず用意する
- イベント境界（EventBridge/キュー/通知）での疎結合を優先し、参照依存を減らす
- 変更時は `cdk diff` で Export/Import の差分を必ず確認する
- 新規 Export/Import の追加が見えた場合は理由と代替案を必ずレビューする

### 4.2 変動リソースの判定基準

- 自動命名され、デプロイごとに物理 ID が変わりやすい
- 変更のたびに置換が発生しやすい
- 運用上頻繁に更新される（タスク定義、オートスケール、ログ設定など）

### 4.3 例外（クロスリージョン/クロスアカウント）

- 同一 app 内でオブジェクト参照ができない場合は、安定 ID を SSM Parameter 等で参照する
- `cross_region_references` を使う場合は Export の増減を `cdk diff` で必ず確認する

#### 正しい例（Construct 共有）

```python
# network_stack.py
class NetworkStack(Stack):
    def __init__(self, scope, id, **kwargs):
        super().__init__(scope, id, **kwargs)
        self.vpc = ec2.Vpc(self, "Vpc")
```

```python
# app.py
network = NetworkStack(app, "NetworkStack")
EcsStack(app, "EcsStack", vpc=network.vpc)
```

#### 禁止例（Output 参照）

```python
# ❌ 非推奨
CfnOutput(self, "VpcId", value=vpc.vpc_id)
```

---



- Stack 名は CloudFormation 上の意味が分かる名前にする
- Construct は副作用を閉じ込める
- 依存関係は引数で明示する（import で参照しない）

```python
class EcsServiceConstruct(Construct):
    def __init__(self, scope, id, *, vpc, cluster, config):
        ...
```

---

## 5. よくあるアンチパターン（禁止）

- 他 Stack の Output を直接参照する設計
- `Fn::ImportValue` 前提の依存関係
- Output を内部連携用に乱用する



- 巨大 Stack（1000 行超）
- L3 Construct 直書き
- app.py にリソース定義
- 環境差分を if/else で分岐

---

## 6. 運用ルール

- 変更時は必ず `cdk synth` を確認
- PR には diff / synth 結果を添付
- 監査対応では YAML を正とする

---

## 7. 本規約の位置づけ

- AWS 公式ではない
- 実務で最も破綻しにくい「準デファクト」
- ECS / EventBridge / API Gateway を前提に最適化

---

以上

