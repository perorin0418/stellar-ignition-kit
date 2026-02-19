## ディレクトリ構造

```
doc/
│  # プロジェクトルート
│  # SDDおよび全工程成果物の管理単位
│
├─ /requirements/
│   # 要件定義工程の成果物
│   # 「何を作るか」を定義するフェーズ
│   │
│   ├─ system_overview.md
│   │   # システム全体概要（目的・背景・スコープ・関係者）
│   │
│   ├─ nonfunctional_requirements.md
│   │   # 非機能要件（性能・可用性・セキュリティ・運用・制約条件）
│   │
│   ├─ glossary.md
│   │   # 用語定義集（全工程共通参照）
│   │
│   ├─ /business_flows/
│   │   # 業務フロー図群（複数業務想定）
│   │   │
│   │   ├─ BF-01_order_flow.drawio
│   │   │   # 受注業務フロー
│   │   │
│   │   ├─ BF-02_payment_flow.drawio
│   │   │   # 支払業務フロー
│   │   │
│   │   └─ BF-03_admin_flow.drawio
│   │       # 管理業務フロー
│   │
│   └─ /functional_requirements/
│       # 機能単位の要件定義書
│       # 機能IDで管理（トレーサビリティ確保）
│       │
│       ├─ F-01_user_management.md
│       │   # ユーザー管理機能の要件
│       │
│       ├─ F-02_order_management.md
│       │   # 受注管理機能の要件
│       │
│       └─ F-03_reporting.md
│           # レポート出力機能の要件
│
├─ /basic_design/
│   # 基本設計（外部設計）
│   # ユーザー視点での仕様を定義
│   │
│   ├─ architecture_overview.md
│   │   # システム構成図・レイヤ構成・技術スタック概要
│   │
│   ├─ system_context.drawio
│   │   # 外部システムとの関係図
│   │
│   ├─ data_model_logical.er
│   │   # 論理データモデル（ER図）
│   │
│   └─ /features/
│       # 機能単位の基本設計
│       │
│       ├─ /F-01_user_management/
│       │   │
│       │   ├─ overview.md
│       │   │   # 機能概要・前提条件・制約
│       │   │
│       │   ├─ /screens/
│       │   │   # 画面設計（UI仕様）
│       │   │   │
│       │   │   ├─ SCR-01_user_list.fig
│       │   │   │   # ユーザー一覧画面
│       │   │   │
│       │   │   ├─ SCR-02_user_edit.fig
│       │   │   │   # ユーザー編集画面
│       │   │   │
│       │   │   └─ screen_transition.drawio
│       │   │       # 画面遷移図
│       │   │
│       │   ├─ /interfaces/
│       │   │   └─ external_interface.md
│       │   │       # 外部IF仕様（API概要・連携方式）
│       │   │
│       │   └─ business_rules.md
│       │       # 業務ルール定義
│       │
│       ├─ /F-02_order_management/
│       │   # 受注管理機能の基本設計一式
│       │
│       └─ /F-03_reporting/
│           # レポート機能の基本設計一式
│
├─ /detailed_design/
│   # 詳細設計（内部設計）
│   # 開発者向け仕様
│   │
│   ├─ coding_standards.md
│   │   # コーディング規約
│   │
│   ├─ error_handling_policy.md
│   │   # 例外処理方針・ログ方針
│   │
│   └─ /features/
│       # 機能単位の詳細設計
│       │
│       ├─ /F-01_user_management/
│       │   │
│       │   ├─ class_diagram.drawio
│       │   │   # クラス構造図
│       │   │
│       │   ├─ sequence_login.drawio
│       │   │   # 処理シーケンス図
│       │   │
│       │   ├─ api_spec.yaml
│       │   │   # API詳細仕様（OpenAPI）
│       │   │
│       │   ├─ db_physical_design.sql
│       │   │   # テーブル定義（物理設計）
│       │   │
│       │   └─ processing_spec.md
│       │       # 処理仕様書（アルゴリズム・ロジック）
│       │
│       ├─ /F-02_order_management/
│       │   # 受注管理の詳細設計
│       │
│       └─ /F-03_reporting/
│           # レポート機能の詳細設計
│
├─ /implementation/
│   # 実装工程
│   │
│   ├─ /src/
│   │   # ソースコード
│   │   │
│   │   ├─ /user_management/
│   │   │   # F-01実装コード
│   │   │
│   │   ├─ /order_management/
│   │   │   # F-02実装コード
│   │   │
│   │   └─ /reporting/
│   │       # F-03実装コード
│   │
│   └─ build_and_deploy.md
│       # ビルド方法・環境構築手順
│
├─ /testing/
│   # テスト工程
│   │
│   ├─ test_strategy.md
│   │   # テスト方針（観点・カバレッジ基準）
│   │
│   ├─ /unit/
│   │   # 単体テスト
│   │
│   ├─ /integration/
│   │   # 結合テスト
│   │
│   ├─ /system/
│   │   # システムテスト
│   │
│   └─ /acceptance/
│       # 受入テスト（UAT）
│
└─ /release_and_maintenance/
    # リリース・保守工程
    │
    ├─ deployment_plan.md
    │   # 本番移行計画
    │
    ├─ migration_plan.md
    │   # データ移行手順
    │
    └─ operation_manual.md
        # 運用手順書

```



