# Milestones

## このディレクトリの役割

このディレクトリは、Vision を
**実行可能な到達点（Milestone）に分解する計画レイヤー**である。

Milestone は、
- 「いつまでに」
- 「どこまでできていればよいか」
を定義する。

---

## 上位・下位との関係

- 上位参照：
  - `00_vision.md`
  - `05_architecture/`
- 下位展開：
  - `20_features/`

Milestone は **設計や実装の詳細を定義しない**。
それらは Feature / Task の責務とする。

---

## Milestone に含めるべき内容

各 Milestone には、最低限以下を含める。

- Milestone 名
- 目的（この段階で達成すべき状態）
- スコープ（含む / 含まない）
- 完了条件（Definition of Done）
- 対象 Feature 一覧
- 依存関係（前提 / 後続）
- リスク・注意点

---

## Milestone に含めてはいけない内容

以下は **他レイヤーの責務**であり、ここには書かない。

- 実装手順や技術詳細 → `30_tasks/`
- 機能の細かな仕様 → `20_features/`
- システム構成の決定 → `05_architecture/`
- 設計判断の背景 → `15_design-notes/`

---

## 拘束力レベル

Milestone の記述は、以下の意味を持つ。

- **スコープ・完了条件**
  → Feature / Task は必ず満たす必要がある
- **期限・到達点**
  → 優先順位付け・取捨選択の判断基準となる
- **未記載事項**
  → Feature 側で補完してよいが、逸脱は禁止

---

## AI / opencode 向け指示

- Feature を生成・修正する際は、
  必ず対応する Milestone を参照すること
- Feature が Milestone のスコープを超える場合、
  **計画逸脱として扱い、人の判断を要する**
- Milestone は自動的に分割・統合してはならない

---

## 命名・ファイルルール

- ファイル名は以下を推奨する
```
milestone-<Vision No.>.<Milestone No>-<簡潔な変更名>.md
milestone-1.2-Foundation.md
milestone-2.3-CoreFeatures.md
milestone-3.4-Release.md
```

- Milestone ID は Feature / Task から参照される前提とする

---

## 更新ルール

- Milestone の変更は **計画変更**である
- Feature / Task 側の都合で勝手に変更しない
- 変更理由・影響範囲を必ず明記する

---

## まとめ

Milestone は、
**「何を作るか」ではなく「どこまで進めばよいか」を定義する層**である。

このディレクトリは、
プロジェクト全体の進捗とスコープを守るための
**計画上の基準点**として扱う。
