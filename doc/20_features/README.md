# Features

## このディレクトリの役割

このディレクトリは、Milestone で定義された到達点を、
**実装可能な機能単位（Feature）に落とし込む設計レイヤー**である。

Feature は、
- 何を作るか（What）
- どのような価値を提供するか
を定義する。

---

## 上位・下位との関係

- 上位参照：
  - `05_architecture/`
  - `10_milestones/`
  - `15_design-notes/`
- 下位展開：
  - `30_tasks/`

Feature は **Architecture を変更しない**。
Architecture に関わる判断は `05_architecture/` に集約する。

---

## Feature に含めるべき内容

各 Feature には、最低限以下を含める。

- Feature 名（IDと所属 Milestone）
- 概要
- 機能要件
- 非機能要件（必要な場合）
- インターフェース
- 設計方針
- 対象 Task 一覧
- 受け入れ条件

---

## Feature に含めてはいけない内容

以下は **他レイヤーの責務**であり、ここには書かない。

- 実装手順・コード詳細 → `30_tasks/`
- システム構成の変更 → `05_architecture/`
- 計画・スケジュールの変更 → `10_milestones/`
- 一時的な検討メモ → `15_design-notes/`

---

## 拘束力レベル

Feature の記述は、以下の意味を持つ。

- **機能要件・非機能要件**
  → Task は必ず満たす必要がある
- **受け入れ条件**
  → 実装完了判定の基準となる
- **未記載の実装詳細**
  → Task 側で決定してよい

---

## AI / opencode 向け指示

- Task を生成・修正する際は、
  必ず対応する Feature を参照すること
- Feature の要件を満たさない Task は無効とする
- Feature の範囲を超える実装は **スコープ逸脱**と判断する
- Feature を分割・統合する場合は、人の判断を要する

---

## Architecture との関係（重要）

- Feature は `05_architecture/` を **前提条件として扱う**
- Architecture に記載されていない構成要素を
  Feature 側で新規導入してはならない
- 矛盾がある場合、Feature を修正対象とする

---

## 命名・ファイルルール

- ファイル名は以下を推奨する
```
feature-<Vision No.>.<Milestone No>.<Feature No>-<簡潔な変更名>.md
feature-1.2.3-Authentication.md
feature-2.3.4-DataIngestion.md
feature-3.4.5-Dashboard.md
```

- Feature ID は Milestone / Task から参照可能であること

---

## 更新ルール

- Feature の変更は **設計変更**である
- Task の都合による逆流変更は禁止
- 変更理由・影響範囲を必ず明記する

---

## まとめ

Feature は、
**「何を作るか」を定義し、
「どう作るか」は Task に委ねる層**である。

このディレクトリは、
設計と実装の責務を分離するための
**最終防衛線**として扱う。
