# Tasks

## このディレクトリの役割

このディレクトリは、Feature で定義された内容を、
**実装計画として具体化するレイヤー（Task）**である。

Task は、
- 「何を達成するか」
- 「どの範囲を変更するか」
- 「どの ChangeSet に分解するか」
を明示する。

**実際のコード変更指示は ChangeSet に委ねる**。

---

## 上位・下位との関係

- 上位参照：
  - `20_features/`
  - `25_interfaces/`
  - `05_architecture/`
- 下位：
  - 各 Task 配下の `changeset-XX/`

Task は **新たな要件・設計を導入しない**。  
また、Task 自体にはコードレベルの操作指示を書かない。

---

## Task と ChangeSet の責務分離（重要）

### Task の責務

- 目的・背景の明示
- 実装範囲の確定
- 変更対象の洗い出し（粒度はファイル・モジュール程度）
- ChangeSet への分解
- 完了条件・成果物の定義
- テスト観点の整理

### ChangeSet の責務

- 具体的な変更操作の指示
- ファイル単位・作業単位の変更内容
- Add / Modify / Delete の明示

Task は **ChangeSet を束ねる上位計画**である。

---

## Task に含めるべき内容

各 Task には、最低限以下を含める。

- Task 名
- 作業目的
- 作業内容
- 実装制約
- 完了条件
- 対象 ChangeSets
- チェックリスト

※ 実装手順・処理詳細は ChangeSet に記載する。

---

## Task に含めてはいけない内容

以下は **Task の責務外**であり、ここには書かない。

- ファイル内の具体的な修正手順 → ChangeSet
- 関数単位の変更内容 → ChangeSet
- 機能追加の提案 → `20_features/`
- API / スキーマ変更 → `25_interfaces/`
- 構成変更 → `05_architecture/`
- 設計判断の背景 → `15_design-notes/`

---

## 拘束力レベル（重要）

Task の記述は、以下の拘束力を持つ。

- **変更範囲・非変更範囲**
  → ChangeSet はこれを厳守する
- **ChangeSet 分割方針**
  → 実装者（人 / AI）は勝手に統合・分割しない
- **未記載の判断**
  → Task でも ChangeSet でも補完しない
- **矛盾がある場合**
  → 上位レイヤーを正とする

---

## AI / opencode 向け指示

- `/implement` は **ChangeSet を直接実行対象**とする
- Task は実装の前提条件・制約としてのみ参照する
- Task と ChangeSet に矛盾がある場合は実装を止める
- ChangeSet が存在しない Task は実装不可とする

---

## 命名・ファイルルール

### Task ディレクトリ

```
task-<Vision No.>.<Milestone No.>.<Feature No.>.<Task No.>-<簡潔な変更名>.md/
task-1.2.3.4-AddAuthMiddleware.md/
task-2.3.4.5.-CreateUserTable.md/
task-3.4.5.6-UpdateApiHandler.md/
```


### Task 配下

```
task-01-add-auth-middleware/
├─ README.md
├─ task-1.2.3.4-AddAuthMiddleware.md/
├─ task-2.3.4.5.-CreateUserTable.md/
```


- Task ID は Feature から参照される前提とする

---

## 実装粒度の基準

- 1 Task = 複数 ChangeSet を許容
- 1 ChangeSet = 1 PR / 1 commit 相当
- Task に「コード操作レベルの判断」が含まれる場合は不適切
- ChangeSet が 1 つで済む場合でも、Task は省略しない

---

## 更新ルール

- Task は **計画の履歴**として残す
- 完了後に Task / ChangeSet を書き換えない
- 問題があれば新しい Task または ChangeSet を追加する

---

## まとめ

Task は、
**「何を、どこまでやるか」を確定させる計画書**であり、
ChangeSet は、
**「どう手を動かすか」を定義する命令書**である。

このディレクトリは、
ChangeSet を安全に機能させるための
**統制レイヤー**として扱う。
