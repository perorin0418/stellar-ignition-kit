# Nuxt.js (Nuxt 3) フロントエンド開発規約

## 0. 基本方針（MUST）

- 対象は Nuxt 3（Vue 3 + TypeScript）とする
- **最小の抽象化**と**最小のフォルダ数**を最優先する
- 1ファイル1責務を徹底し、UIは薄く保つ
- 迷ったら **標準機能（Nuxt/Vue）を使う**

---

## 1. ディレクトリ構成（MUST）

### 1.1 基本構成（最小）

```
frontend/
  app.vue
  nuxt.config.ts
  components/
  composables/
  pages/
  plugins/
  assets/
  styles/
  tests/
  types/
```

- 追加フォルダは **必要性を説明できる場合のみ** 作成する
- `composables/` と `components/` は **深い階層化を禁止**（最大 1 階層まで）

---

## 2. 命名規則（MUST）

- `components/`: **PascalCase**（例: `UserCard.vue`）
- `pages/`: **kebab-case**（例: `user-settings.vue`）
- `composables/`: **useXxx.ts**（例: `useUser.ts`）
- `types/`: **PascalCase** or **camelCase**（型名と一致）

---

## 3. CSS / Styling（MUST）

### 3.1 原則
- SFCのスタイルは **`scoped` を原則** とする
- 再利用スタイルは `styles/` に集約し、UIの分散定義を避ける
- グローバルスタイルは **最小限**（トークン/ベース/ユーティリティのみ）

### 3.2 `styles/` 構成（最小）
```
styles/
  tokens.css
  base.css
  utilities.css
```
- `components/` は **同種のスタイルが3箇所以上で再利用される場合のみ** 追加可
- グローバル読込は `app.vue` または `nuxt.config.ts` に集約する

### 3.3 命名規則
- CSSクラスは **kebab-case**
- コンポーネント用は **コンポーネント名に対応したブロック名**
- ユーティリティは **`u-` 接頭辞**（例: `u-stack`）
- 状態は **`is-` / `has-` 接頭辞**

### 3.4 トークン/変数
- デザイン値は **`styles/tokens.css` の CSS変数** を基本とする
- 色・余白・サイズは **`var(--token)`** 経由で参照する
- 直接値が必要な場合は **最小限** とし、理由をコメントで残す

### 3.5 スコープ戦略
- SFCのスタイルは **`scoped` を基本** とする
- `:global()` は **`styles/` か外部UIの上書き時のみ** 使用可
- `:deep()` は **外部UIの上書き時のみ** 最小限で使用する
- インラインスタイルは **CSS変数の受け渡し用途のみ** 許可する

### 3.6 禁止事項（MUST NOT）
- `style` 属性での通常スタイル指定（CSS変数用途を除く）
- `!important` の常用
- IDセレクタによるスタイル指定
- 目的の説明ができないグローバルセレクタ

---

## 4. Pages（画面）

### 4.1 役割
- `pages/` は **画面の組み立てのみ**
- 画面ロジックは `composables/` へ移す

### 4.2 ルール（MUST）
- `definePageMeta` は **ファイル先頭に集約**
- API 呼び出しは **composable 経由**
- `useAsyncData` / `useFetch` は **composable 内で使用**

---

## 5. Components（UI）

- UI のみを担当し、状態や副作用は持たない
- `props` は **型必須**、`any` 禁止
- ロジックが増えたら composable に切り出す

---

## 6. Composables

- 1ファイル1 composable（`useXxx`）
- API 呼び出し、状態管理、データ整形を集約する
- 副作用は **明示的に記述** する（例: `navigateTo`）

---

## 7. 状態管理

- 共有状態は **`useState` を第一選択** とする
- Pinia は **複数画面で明確に再利用される場合のみ** 採用する
- グローバル状態の乱立は禁止

---

## 8. Plugins

- **単一画面でしか使わない場合は作らない**
- Nuxt の自動登録が基本。明示登録は **順序依存や client/server 限定時のみ**
- `provide` する値は **最小限** にする

---

## 9. TypeScript（MUST）

- `strict: true` 前提
- `any` 禁止
- 境界層（API/外部入力）は `unknown` → **型ガードで絞り込み**
- 再利用する型のみ `types/` に置く

---

## 10. Lint / Format

- `eslint` と `prettier` を必須とする
- `@nuxt/eslint` を優先し、ルールは **最小限** にする
- ルール変更時は `doc/90_guidelines/` に理由を残す

---

## 11. テスト

- composable は **Vitest で単体テスト**
- 重要フローは **最小限の E2E（Playwright など）** を用意
- UI の見た目検証は必須範囲に含めない

---

## 12. 禁止事項（MUST NOT）

- `pages/` で API を直接呼ぶ
- `any` の常用、型定義の放置
- plugin の乱立、グローバル状態の無秩序な追加
- 目的の説明できないフォルダ増設
- UI と業務ロジックの同居

---

## 13. レビュー観点（チェックリスト）

- [ ] `pages/` に API 呼び出しが直接書かれていないか
- [ ] composable は `useXxx` で 1ファイル1責務か
- [ ] `any` が存在しないか（境界層は `unknown` + 型ガードか）
- [ ] `components/` と `composables/` の階層が 1 段以内か
- [ ] plugin が必要最小限で、自動登録の範囲内か
- [ ] `useState` / Pinia の使い分けが適切か
- [ ] Lint/Format/Typecheck が通るか
- [ ] `styles/` の構成が最小で、`scoped` が基本になっているか
