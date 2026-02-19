## ディレクトリ構造

```
doc/
│  # プロジェクトルート
│  # SDDおよび全工程成果物の管理単位
│
├─ requirements/
│   # 要件定義工程の成果物
│   # 「何を作るか」を定義するフェーズ
│   │
│   ├─ system_overviews/
│   │   └─ system_overview.md
│   │       # システム全体概要（目的・背景・スコープ・関係者）
│   │
│   ├─ nonfunctional_requirements
│   │   └─ nonfunctional_requirements.md
│   │       # 非機能要件（性能・可用性・セキュリティ・運用・制約条件）
│   │
│   ├─ glossaries/
│   │   └─ glossary.md
│   │       # 用語定義集（全工程共通参照）
│   │
│   ├─ business_codes/
│   │   └─ business_codes.md
│   │       # 業務コード定義（全工程共通参照）
│   │
│   ├─ business_flows/
│   │   # 業務フロー図群（複数業務想定）
│   │   │
│   │   ├─ XX/
│   │   │   # 業務コード
│   │   │   │
│   │   │   ├─ XX0001BF_order-flow/
│   │   │   │   └─ order-flow.mmd
│   │   │   │   # 受注業務フロー
│   │   │   │
│   │   │   ├─ XX0002BF_payment-flow/
│   │   │   │   └─ payment-flow.mmd
│   │   │   │       # 支払業務フロー
│   │   │   │
│   │   │   └─ XX0003BF_admin-flow/
│   │   │       └─ admin-flow.mmd
│   │   │           # 管理業務フロー
│   │   ├─ YY/
│   │   └─ ZZ/
│   │
│   └─ functional_requirements/
│       # 機能単位の要件定義書
│       │
│       ├─ XX/
│       │   # 業務コード
│       │   │
│       │   ├─ XX0001FR_user-management/
│       │   │   └─ user-management.md
│       │   │       # ユーザー管理機能の要件
│       │   │
│       │   ├─ XX0002FR_order-management/
│       │   │   └─ order_management.md
│       │   │       # 受注管理機能の要件
│       │   │
│       │   └─ XX0003FR_reporting/
│       │       └─ reporting.md
│       │           # レポート出力機能の要件
│       ├─ YY/
│       └─ ZZ/
│
└─ /design/
    # 設計
    │
    ├─ coding_standards/
    │   # コーディング規約
    │   │
    │   ├─ python/
    │   │   # Pythonコーディング規約
    │   │   └─ python_coding_standards.md
    │   └─ go/
    │       # Goコーディング規約
    │       └─ go_coding_standards.md
    │
    ├─ log_policies/
    │   # ログ方針
    │   │
    │   ├─ log_format.md
    │   │   # ログフォーマット定義
    │   │
    │   └─ log_levels.md
    │       # ログレベル定義
    │ 
    ├─ error_policies/
    │   # エラー処理方針
    │   │
    │   ├─ error_handling.md
    │   │   # エラー処理の基本方針
    │   │
    │   └─ error_codes.md
    │       # エラーコード定義
    │
    └─ features/
        # 機能単位の設計
        │
        ├─ XX/
        │   # 業務コード
        │   │
        │   ├─ XX0001DD_user-management/
        │   │   │
        │   │   ├─ overview.md
        │   │   │   # 機能概要・前提条件・制約
        │   │   │
        │   │   ├─ system_architecture.mmd
        │   │   │   # システム構成図
        │   │   │
        │   │   ├─ layered_architecture.mmd
        │   │   │   # レイヤ構成図
        │   │   │
        │   │   ├─ tech_stack.md
        │   │   │   # 技術スタック詳細（フレームワーク・ライブラリ・ミドルウェア・クラウドサービス）
        │   │   │
        │   │   ├─ screens/
        │   │   │   # 画面設計（UI仕様）
        │   │   │   │
        │   │   │   ├─ XX0001SC-01_user-list.fig
        │   │   │   │   # ユーザー一覧画面
        │   │   │   │
        │   │   │   ├─ XX0001SC-02_user-edit.fig
        │   │   │   │   # ユーザー編集画面
        │   │   │   │
        │   │   │   └─ screen-transition.mmd
        │   │   │       # 画面遷移図
        │   │   │
        │   │   ├─ interfaces/
        │   │   │   └─ external_interface.md
        │   │   │       # 外部IF仕様（API概要・連携方式）
        │   │   │
        │   │   ├─ business_rules.md
        │   │   │   # 業務ルール定義
        │   │   │
        │   │   ├─ sequence.mmd
        │   │   │   # 処理シーケンス図
        │   │   │
        │   │   ├─ api_spec.yaml
        │   │   │   # API詳細仕様（OpenAPI）
        │   │   │
        │   │   ├─ db_physical_design.sql
        │   │   │   # テーブル定義（物理設計）
        │   │   │
        │   │   └─ processing_spec.md
        │   │       # 処理仕様書（アルゴリズム・ロジック）
        │   │
        │   ├─ XX0002DD_order-management/
        │   │   # 受注管理の詳細設計
        │   │
        │   └─ XX0003DD_reporting/
        │       # レポート機能の詳細設計
        ├─ YY/
        └─ ZZ/

```



