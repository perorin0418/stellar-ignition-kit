# Electron E2E / Playwright コーディング規約

（実認証・外部ブラウザ連携・失敗時証跡重視）

## 0. 基本方針（MUST）

- 本規約は **Electron デスクトップアプリに対する E2E テスト実装規約** とする
- E2E テストフレームワークは **Playwright を標準採用** とする
- 認証を含むシナリオは **外部認証 UI を実際に通過する方式** を前提とし、認証成功をトークン注入や API モックだけで代替してはならない
- テスト資産は **リポジトリ直下の `/e2e`** に集約し、アプリ本体配下へ分散配置してはならない
- テストは **シナリオ単位でフォルダー / ファイルを分離** し、1 ファイルに複数機能の関心事を混在させてはならない
- Electron アプリ操作と外部認証ブラウザ操作は **責務を分離** し、シナリオコードからウィンドウ取得・環境変数読込・資格情報解決を直接書いてはならない
- セレクタは **`data-testid` を第一選択** とし、CSS セレクタや表示文言依存を最小化する
- 待機は **Playwright の自動待機と期待条件** を基本とし、固定スリープへ依存してはならない
- 実認証を含むため、**秘密情報保護・テストアカウント分離・失敗時証跡・タイムアウト設計・CI での検知性** を規約で明示する
- ローカル実行と CI 実行で設定差分を吸収しつつ、**同一シナリオを同一規約で再実行できる構成** を維持する
- 暫定的な回避策を入れる場合でも、**恒久対応の方向性と除去条件をコメントまたはレビュー記録へ残す**
- Electron アプリ構造、プロセス境界、`.env`、設定定数、main / preload / renderer の責務、画面 / IPC / 依存関係ルールは **`electron-typescript-coding-guide.md` を正本** とし、本規約では E2E に必要な接点のみ扱う
- 秘密情報のログ出力禁止、明示的エラーハンドリング、外部依存の扱いは **`backend-golang-coding-guide.md` の共通原則に依拠** し、本規約では E2E 実装時の適用方法のみ定める
- 設定の正本一元化、責務分離、暗黙より明示、ローカル / CI 設定分離は **`infra-python-cdk-guide.md` の原則に依拠** し、本規約では E2E の設定設計へ読み替えて適用する

---

## 1. 対象範囲と責務分離

### 1.1 対象範囲（MUST）

本規約の対象は次とする。

- Electron デスクトップアプリの主要利用者操作
- 外部認証 UI を経由するログイン / ログアウト
- 権限別メニュー表示、主要画面遷移、主要業務フローの正常系 / 権限制御 / 代表的異常系
- Electron 固有のウィンドウ制御、外部ブラウザ遷移、認証コールバック待機

本規約の対象外は次とする。

- 単体テスト、結合テスト、API 契約テストの詳細実装
- バックエンド内部ロジック単体の検証
- Electron アプリ本体のディレクトリ構成・画面構造・IPC 契約の詳細定義
- 手動試験専用の探索的確認手順

### 1.2 依拠先規約との責務分担（MUST）

- **Electron アプリ本体の設計**
  - `electron-typescript-coding-guide.md` を参照する
  - main / preload / renderer の責務、設定定数の置き場所、`.env` の扱い、画面構成、IPC、依存関係はアプリ規約に従う
- **共通実装原則**
  - `backend-golang-coding-guide.md` の「秘密情報を出力しない」「エラーを握りつぶさない」「外部依存を明示する」を E2E 実装へ適用する
- **設定設計・構成原則**
  - `infra-python-cdk-guide.md` の「設定の正本を 1 箇所に寄せる」「暗黙より明示」「責務分離」「ローカル / CI を分離する」を E2E 構成へ適用する

### 1.3 E2E 層ごとの責務（MUST）

- `e2e/scenarios/`: 利用者視点のシナリオ意図、前提、期待結果のみを書く
- `e2e/fixtures/`: Electron 起動、外部認証ブラウザ取得、テストアカウント解決、後始末を担う
- `e2e/helpers/`: 待機、再試行、証跡保存、ログ秘匿化などの横断処理を担う
- `e2e/selectors/`: `data-testid` とアクセシブルロールの参照定義を担う
- `e2e/config/`: アプリ側設定と E2E 専用設定の読込・検証を担う
- `e2e/reports/`, `e2e/artifacts/`: 実行結果の保存先のみを担い、テストロジックを置いてはならない

### 1.4 シナリオコードの禁止責務（MUST NOT）

- シナリオファイルに資格情報の値を直書きしてはならない
- シナリオファイルに `process.env.*` 参照を散在させてはならない
- シナリオファイルに XPath / 長い CSS セレクタを直書きしてはならない
- シナリオファイルに `waitForTimeout(...)` を常設してはならない
- シナリオファイルに Electron 起動コマンド組み立てやログ保存パス生成を直書きしてはならない

---

## 2. ディレクトリ構成

### 2.1 推奨ディレクトリツリー

```text
/ (リポジトリルート)
├── e2e/
│   ├── README.md
│   ├── .env.example
│   ├── playwright.shared.config.ts
│   ├── playwright.local.config.ts
│   ├── playwright.ci.config.ts
│   ├── config/
│   │   ├── env.ts
│   │   ├── projects.ts
│   │   └── timeouts.ts
│   ├── fixtures/
│   │   ├── electronApp.fixture.ts
│   │   ├── authBrowser.fixture.ts
│   │   ├── session.fixture.ts
│   │   └── testAccounts.fixture.ts
│   ├── helpers/
│   │   ├── artifacts.ts
│   │   ├── logRedaction.ts
│   │   ├── retry.ts
│   │   ├── wait.ts
│   │   └── windows.ts
│   ├── selectors/
│   │   ├── common/
│   │   │   └── shell.selectors.ts
│   │   ├── auth/
│   │   │   └── login.selectors.ts
│   │   └── <domain>/
│   │       └── <feature>.selectors.ts
│   ├── scenarios/
│   │   ├── auth/
│   │   │   ├── login/
│   │   │   │   ├── login-success.spec.ts
│   │   │   │   ├── login-timeout.spec.ts
│   │   │   │   └── logout.spec.ts
│   │   │   └── authorization/
│   │   │       └── admin-menu-visibility.spec.ts
│   │   └── <domain>/
│   │       └── <feature>/
│   │           └── <scenario>.spec.ts
│   ├── .auth/
│   │   ├── .gitignore
│   │   └── README.md
│   ├── artifacts/
│   │   └── .gitignore
│   └── reports/
│       └── .gitignore
├── <app-root>/
│   ├── .env
│   ├── .env.example
│   └── package.json
└── doc/
```

### 2.2 構成ルール（MUST）

- ドメイン分類ディレクトリは **アプリ本体の機能分類と一対一で対応する命名** に揃える
- シナリオは `e2e/scenarios/<domain>/<feature>/` へ配置する
- 画面横断の共通 UI は `e2e/selectors/common/` と `e2e/helpers/` へ集約する
- Playwright 関連依存関係の正本は **原則として `/e2e/package.json` と対応する lockfile** に置き、E2E を `/e2e` 配下で独立管理する
- 移行期間などでアプリ本体側の `package.json` を正本とする場合でも、**依存関係の正本は 1 箇所に限定** し、E2E 側へ重複定義してはならない
- 認証済み状態の一時ファイルは `e2e/.auth/` 配下とし、**Git 管理対象へ含めてはならない**
- 実行証跡は `e2e/artifacts/` と `e2e/reports/` に分離し、成果物種別を混在させてはならない

### 2.3 `README.md` の責務（SHOULD）

`e2e/README.md` には少なくとも次を記載する。

- 実行前提（必要なアプリ側 `.env` と E2E 用秘密情報）
- Playwright 依存関係の正本の配置場所、インストール起点、更新責務
- ローカル実行手順
- CI 実行時の差分
- テストアカウントの利用ルール
- 失敗時証跡の保存場所

---

## 3. 命名規則

### 3.1 ディレクトリ・ファイル

- シナリオ分類ディレクトリ: **ケバブケース**
  - 例: `data-export`, `user-settings`, `login`
- テストファイル: `<scenario-name>.spec.ts`
  - 例: `login-success.spec.ts`, `create-item.spec.ts`
- fixture ファイル: `<purpose>.fixture.ts`
  - 例: `electronApp.fixture.ts`, `session.fixture.ts`
- helper ファイル: **キャメルケース**
  - 例: `logRedaction.ts`, `artifacts.ts`
- selector 定義ファイル: `<screen-or-feature>.selectors.ts`
  - 例: `login.selectors.ts`, `shell.selectors.ts`
- Playwright 設定ファイル: `playwright.<mode>.config.ts`
  - 例: `playwright.local.config.ts`, `playwright.ci.config.ts`

### 3.2 テスト名・describe 名（MUST）

- `test(...)` 名は **利用者視点の結果** を表す日本語または明確な英語とする
- `describe(...)` は業務目的単位で分ける
- 実装詳細だけを表す曖昧な名前を避ける

```ts
// ✅ 推奨
test('管理者がログインするとユーザー管理メニューが表示される', async () => {
  // ...
});
```

### 3.3 定数・タグ（SHOULD）

- タイムアウト、再試行回数、アカウント種別は定数化する
- タグは `@smoke`, `@auth`, `@serial`, `@nightly` のように **目的が分かる単語** を用いる
- CI 条件分岐で使うタグは `e2e/config/projects.ts` に集約する
- 定数配置や命名方針は、対象アプリの `electron-typescript-coding-guide.md` に定義された定数管理方針へ合わせる

---

## 4. シナリオ分割規約

### 4.1 1 シナリオ = 1 業務目的（MUST）

- 1 つの `*.spec.ts` は **1 つの業務目的または 1 つの異常条件** を主題とする
- 「ログイン + データ作成 + 一覧確認 + 出力」のような **多目的な長大シナリオ** を 1 ファイルへ詰め込んではならない
- 正常系、権限制御、タイムアウト、入力異常は原則として別ファイルへ分ける

### 4.2 分割単位の推奨（SHOULD）

- 認証系: `login-success.spec.ts`, `login-timeout.spec.ts`, `logout.spec.ts`
- 権限制御系: `admin-menu-visibility.spec.ts`, `viewer-access-denied.spec.ts`
- CRUD 系: 一覧取得、作成、更新、削除、権限変更を必要に応じて分離する
- 再現性の低い長時間シナリオは `@nightly` を付与し、PR の常時実行対象から分離する

### 4.3 共有前提の扱い（MUST）

- 認証済み状態の生成は fixture へ寄せる
- 画面到達までの共通導線は helper または fixture 化する
- ただし、**期待結果の検証自体** はシナリオ側へ残し、helper に業務判定を書いてはならない

### 4.4 シリアル実行の基準（MUST）

次に該当するシナリオは直列実行を基本とする。

- 同一テストアカウントを共有するシナリオ
- 同一業務データを更新し、ロールバック不能または競合しやすいシナリオ
- ログイン失敗回数制限やアカウントロックへ影響しうるシナリオ
- 外部認証ブラウザの状態を専有するシナリオ

---

## 5. セレクタ規約

### 5.1 優先順位（MUST）

セレクタの優先順位は次とする。

1. `getByTestId(...)`
2. `getByRole(...)` + `name`
3. `getByLabel(...)`
4. `getByText(...)`（表示文言自体が仕様の一部である場合のみ）

### 5.2 `data-testid` 命名規則（MUST）

- 形式は `<screen>-<region>-<element>-<purpose>` を推奨する
- 省略しすぎず、役割が一意に分かる名前にする
- DOM 構造変更で壊れにくい **意味ベース** の命名にする

```text
login-screen-start-button
shell-header-current-user-label
settings-screen-save-button
search-form-submit-button
```

### 5.3 禁止事項（MUST NOT）

- `.MuiButton-root > span` のようなスタイル依存セレクタを使ってはならない
- XPath を常用してはならない
- nth-child 依存の脆い指定を常設してはならない
- 1 回しか使わないからという理由で、シナリオ中へ生文字列セレクタを散在させてはならない

### 5.4 セレクタ定義の分離（SHOULD）

- セレクタ文字列は `e2e/selectors/` へ寄せる
- 共通シェル領域と画面固有領域を分ける
- 画面改修時は、まず selector 定義ファイルの差分有無を確認する

---

## 6. 待機・リトライ規約

### 6.1 基本ルール（MUST）

- 待機は `await expect(locator).toBeVisible()` などの **状態ベース待機** を使う
- API 応答待ちや非同期反映確認には `expect.poll(...)` や `waitForResponse(...)` を使う
- 画面初期化完了は、`data-testid` を付与したルート要素やローディング消失で判定する
- `waitForTimeout(...)` はデバッグ中の一時利用を除き、コミットしてはならない

### 6.2 タイムアウト設計（MUST）

- 通常 UI 操作タイムアウトと、実認証 / 外部 API を含むタイムアウトは分離する
- ログイン待機の既定値は、**アプリ側の認証タイムアウト定義を下回ってはならない**
- ログアウト待機も、**アプリ側の認証関連タイムアウトと整合する値** を採用する
- タイムアウト値は `e2e/config/timeouts.ts` に集約する

### 6.3 リトライ設計（SHOULD）

- ローカル実行は `retries = 0` を基本とする
- CI 実行は **外部認証や一時的通信揺らぎ** に限って `retries = 1` または `2` を許容する
- リトライは flaky さの隠蔽に使ってはならず、継続発生する失敗は根本原因を修正する
- 失敗時に再試行した事実はレポートへ残す
- 再試行で成功した失敗、同一シナリオの断続的失敗、worker 依存の失敗は **flaky 候補** として識別し、シナリオ名 / project 名 / worker 番号 / 再試行有無を追跡可能にする
- flaky 監視の閾値、担当者、定期棚卸し方法は別文書で定義してよいが、少なくとも **「いつ issue 化するか」「いつ `@nightly` へ退避するか」「いつ PR ゲートから外せるか」** を定義しなければならない

### 6.4 再試行の禁止対象（MUST NOT）

- 権限不足を確認する異常系
- 明示的なバリデーションエラー確認
- データ競合や設計不備が疑われる失敗
- 秘密情報漏えい検知などのセキュリティ観点の失敗

---

## 7. Electron 固有のウィンドウ遷移規約

### 7.1 基本ルール（MUST）

- Electron アプリウィンドウ取得は `electronApp.fixture.ts` へ集約する
- 外部認証ブラウザ取得は `authBrowser.fixture.ts` へ集約する
- シナリオから `BrowserWindow` 相当の取得方法や再フォーカス処理を直接扱ってはならない
- ログイン後は **認証ブラウザ側の完了確認 → Electron 側の再前面化 → セッション反映確認** の順で検証する

### 7.2 アプリ本体規約との関係（MUST）

- Electron の main / preload / renderer の責務境界は `electron-typescript-coding-guide.md` を参照する
- E2E 基盤は、**アプリ本体のプロセス境界を壊さずに観測・操作できる範囲で実装** しなければならない
- renderer が本来触れてはならない API を E2E 都合で追加公開してはならない
- ブラウザ起動責務や認証コールバックの処理が自動化困難な場合でも、**本物の認証 URL を使いながらテスト制御可能な責務分離を行う** 方針を優先する
- OS 既定ブラウザへ完全委譲したままの方式は、自動化の再現性・検知性・CI 互換性が低いため、恒久運用の標準にしてはならない

### 7.3 ウィンドウ遷移の確認項目（SHOULD）

- 認証開始前に Electron 主ウィンドウが 1 つに保たれていること
- 認証開始後に外部ブラウザコンテキストが取得できること
- コールバック後に Electron 側へ制御が戻ること
- 不要な補助ウィンドウが残留していないこと
- 認証コールバック URI の待受先が想定どおり機能すること

---

## 8. 実認証テストの扱い

### 8.1 テストアカウント方針（MUST）

- 実認証用アカウントは **個人利用アカウントと分離した専用テストアカウント** を使う
- 少なくとも次の権限群に対応するアカウントを用意する
  - 一般利用者
  - 管理者
  - 特権利用者またはそれに相当する上位権限
- 同一アカウントを高並列で使い回してはならない
- 実認証を伴う worker 数は、**同時利用可能な専用テストアカウント数と排他制御可能なデータ範囲の小さい方** を上限とする
- 原則として **1 worker = 1 専用テストアカウント** を割り当て、同一アカウントを複数 worker で共有する場合は `@serial` 指定と共有理由を必須とする
- 権限や初期データは再実行可能な形で保守し、手作業前提の状態依存を増やしてはならない
- アカウント割当、解放、ロック検知、ローテーションの詳細運用は別文書で定義してよいが、少なくとも **worker 上限、権限別アカウント本数、共有禁止条件、障害時の退避運用** を明記しなければならない

### 8.2 資格情報・秘密情報（MUST）

- ユーザー名、パスワード、クライアントシークレット相当情報、トークン、Cookie、Authorization ヘッダーを Git 管理対象へ保存してはならない
- テスト資格情報は CI の Secret Store または Git 管理外の `e2e/.env.local` 等から読み込む
- `storageState` を利用する場合も、永続保存先は `e2e/.auth/` 配下とし、コミットしてはならない
- Playwright trace / screenshot / video / console 出力へ秘密情報が残らないよう、`helpers/logRedaction.ts` 相当の共通秘匿化を通す

### 8.3 ログ出力ルール（MUST）

- 認証コード、ID トークン、アクセストークン、リフレッシュトークン、Authorization、Cookie を平文出力してはならない
- 失敗ログには **シナリオ名、ドメイン分類、画面名、待機対象、タイムアウト値** を残す
- ユーザー識別子は必要最小限に留め、メールアドレス全文や個人情報の無制限出力を避ける
- エラーの握りつぶしは禁止とし、失敗箇所が追跡できる文脈を明示する

### 8.4 タイムアウトと失敗時対応（MUST）

- 認証ページ表示待機、認証完了待機、Electron 側のセッション反映待機を別々に計測する
- 外部認証を含むため、**失敗時は必ず trace / screenshot / video / Playwright HTML report を残す**
- 可能であれば Electron main process のログと renderer console の抜粋も同時に保存する
- 認証タイムアウト失敗時は、次回実行へ影響する残存ブラウザや残存セッションを後始末する

### 8.5 実認証テストの適用範囲（SHOULD）

- すべての画面テストで毎回フルログインし直すのではなく、**権限別に 1 回は本物の認証を通したうえで、その実行単位内で安全に再利用可能なセッションを使う**
- ただし、ログイン導線自体を検証するシナリオでは必ず初回認証から確認する
- 認証連携そのものの健全性を監視するため、`@auth` タグ付きスモークを PR または定期実行へ残す

### 8.6 `@auth` スモーク失敗時の扱い（MUST）

- `@auth` タグ付きスモークの失敗は、**認証基盤障害の可能性を含むシグナル** として扱い、単純に flaky 扱いで握りつぶしてはならない
- 初動切り分けでは少なくとも次を確認する
  - 外部認証 UI への到達可否、HTTP ステータス、証明書 / DNS / ネットワーク異常の有無
  - 認証コールバック URI への到達可否、Electron 側セッション反映可否、同一時刻帯の複数シナリオ同時失敗有無
  - 直近差分がセレクタ / fixture / 待機条件などテスト実装起因か、認証設定 / 外部依存起因か
- 認証基盤障害が疑われる場合は **機能不具合と別に記録** し、再実行結果と証跡を紐付けて保管する
- テスト不備が疑われる場合は、再試行で通過しても修正対象として扱い、`retries` 増加だけで恒久対応にしてはならない
- 詳細な連絡系統、再実行回数、障害宣言基準は別文書で定義してよいが、少なくとも **エスカレーション先、暫定的な PR 判定ルール、障害復旧後の再実行条件** を明記しなければならない

---

## 9. 環境変数と設定ファイル

### 9.1 設定の責務分離（MUST）

- アプリ実行時設定の正本は **対象 Electron アプリの `.env` / `.env.example` と設定定数層** とする
- E2E 実行時の秘密情報・上書き設定は `/e2e/.env.example`、Git 管理外の `/e2e/.env.local`、CI 変数で扱う
- `e2e/config/env.ts` はアプリ側設定と E2E 用設定をまとめて検証し、シナリオ側へ生の `process.env` を渡してはならない
- Playwright の依存関係、実行スクリプト、lockfile の正本は 1 系統に固定し、アプリ本体と E2E の双方で同一依存関係を二重管理してはならない
- 設定値の正本一元化とローカル / CI 分離は `infra-python-cdk-guide.md` の原則に従い、暗黙のフォールバックや複数正本を作ってはならない

### 9.2 アプリ側環境変数の利用（MUST）

E2E では次のカテゴリの設定を **アプリ側設定の正本** として利用する。

| 設定カテゴリ | 用途 |
| --- | --- |
| 外部認証 UI の接続先 | ログイン画面 / トークン交換先などの認証基盤接続 |
| 認証関連環境変数 | 認証クライアント ID、テナント識別子、必要なスコープなど |
| 認証コールバック URI | ログイン / ログアウト後にアプリへ戻す URI |
| アプリ API 接続先 | 画面から利用する API / BFF / Gateway などの接続先 |
| 通信タイムアウト設定 | API 呼び出しや認証処理の上限時間 |

> 具体的な変数名・定数名・配置場所は `electron-typescript-coding-guide.md` の規約に従うこと。

### 9.3 E2E 専用環境変数の推奨名（SHOULD）

E2E 専用変数名は次を推奨する。

| 環境変数名 | 用途 |
| --- | --- |
| `E2E_TEST_USERNAME_STANDARD` | 一般利用者テストアカウント |
| `E2E_TEST_PASSWORD_STANDARD` | 一般利用者パスワード |
| `E2E_TEST_USERNAME_ADMIN` | 管理者テストアカウント |
| `E2E_TEST_PASSWORD_ADMIN` | 管理者パスワード |
| `E2E_TEST_USERNAME_PRIVILEGED` | 上位権限テストアカウント |
| `E2E_TEST_PASSWORD_PRIVILEGED` | 上位権限パスワード |
| `E2E_HEADLESS` | headless 実行可否 |
| `E2E_LOGIN_TIMEOUT_MS` | 認証待機タイムアウト |
| `E2E_ARTIFACT_ROOT` | 証跡出力先ルート |
| `E2E_WORKERS` | 実行 worker 数 |

### 9.4 設定ファイルの標準責務（SHOULD）

- `e2e/playwright.shared.config.ts`: 共通 reporter / trace / video / screenshot / timeout 設定
- `e2e/playwright.local.config.ts`: ローカル向け headed 実行、開発者確認用設定
- `e2e/playwright.ci.config.ts`: CI 向け headless 実行、worker 制限、artifact 保存強化
- `e2e/config/projects.ts`: 権限別 / 実行モード別の Playwright project 定義
- `e2e/config/timeouts.ts`: UI / 認証 / API / 後始末タイムアウト定数

---

## 10. 実行モード

### 10.1 ローカル実行（SHOULD）

- 開発者はリポジトリルートで E2E を起動できる形を維持する
- ローカルでは headed 実行を基本とし、デバッグしやすさを優先する
- Playwright 依存関係の正本を `/e2e/package.json` とする場合、インストールと実行は `/e2e` を起点に統一する
- 代表例として、次のようなコマンド構成を推奨する

```bash
npm --prefix ./e2e ci
npm --prefix ./e2e exec playwright test -c ./playwright.local.config.ts
```

- 単一シナリオの再実行例:

```bash
npm --prefix ./e2e exec playwright test ./scenarios/auth/login/login-success.spec.ts -c ./playwright.local.config.ts
```

### 10.2 CI 実行（MUST）

- CI では headless 実行を基本とする
- 秘密情報は CI の Secret Store から注入する
- worker 数は **利用可能なテストアカウント数と外部認証の安定性** を超えてはならない
- worker 数の決定根拠は、少なくとも **権限別アカウント本数、共有可否、`@serial` 対象数** を説明できる状態にする
- 失敗時は HTML report、trace、screenshot、video を必ず保管する
- `/e2e/package.json` を正本とする場合、CI でも Playwright の実行起点は `/e2e` に統一する
- 代表例として、次のようなコマンド構成を推奨する

```bash
npm --prefix ./e2e exec playwright test -c ./playwright.ci.config.ts --reporter=line,html
```

### 10.3 実行対象の分離（SHOULD）

- PR ゲートでは `@smoke` と `@auth` の最小重要フローを優先する
- 全量回帰や長時間シナリオは `@nightly` で定期実行へ分離する
- リリース前確認では権限差分、主要 CRUD、ログイン / ログアウト、重要処理を含める

---

## 11. ログ出力と失敗時証跡

### 11.1 保存対象（MUST）

失敗時は少なくとも次を保存する。

- Playwright trace
- 失敗時 screenshot
- 失敗時 video
- HTML report
- シナリオ名、project 名、worker 番号、再試行回数
- 必要に応じて Electron main process / renderer console の抜粋

### 11.2 保存先ルール（MUST）

- 保存先は `e2e/artifacts/<run-id>/<project>/<scenario>/` を推奨する
- HTML report は `e2e/reports/` 配下へ分離する
- 成果物の上書きを避けるため、`run-id` には時刻または CI 実行 ID を含める
- ローカル生成物は **作業者本人のみが参照可能な保存先** を原則とし、共有が必要な場合は秘匿化後に限定共有する
- CI 生成物の公開範囲は、**当該リポジトリの保守担当者・レビュー担当者・運用担当者など必要最小限** に限定し、一般公開リンクを恒常運用してはならない

### 11.3 秘匿化ルール（MUST）

- trace / log / screenshot 名へ資格情報や個人情報を含めてはならない
- URL のクエリ文字列に認証コードが含まれる場合は出力前にマスクする
- 共有チャットやチケットへ貼る前提のログは、既定で秘匿化済み文字列を使う
- 秘匿化の判断基準は `backend-golang-coding-guide.md` の情報漏えい防止原則に準拠する

### 11.4 保存期間・削除責務（MUST）

- artifact / report の保存期間は必要最小限とし、**通常運用の既定上限は 30 日以内** を目安とする
- インシデント調査などで延長保管する場合は、保管理由、保管先、削除期限、閲覧権限を別文書またはチケットで明示する
- CI 上の artifact / report 削除責務は、原則として **CI 設定の保守担当者** が負う
- ローカルに取得した artifact / report の削除責務は、原則として **取得した開発者本人** が負う
- 詳細な保存先、保持日数の例外、外部共有承認手順は別文書で定義してよいが、少なくとも **保存期間、閲覧権限、延長保管条件、自動削除方法** を定義しなければならない

---

## 12. 禁止事項

### 12.1 実装上の禁止事項（MUST NOT）

- 認証成功状態を localStorage / Cookie / token ファイルの注入だけで再現してはならない
- 外部認証 UI の遷移先を偽 URL やスタブ HTML へ差し替えてはならない
- 個人用アカウントや本番相当の共有運用アカウントを自動試験へ使ってはならない
- 生のパスワードを `*.spec.ts`、`README.md`、コミットメッセージへ記載してはならない
- `waitForTimeout(...)`、過剰な `force: true`、無制限リトライで不安定さを隠してはならない
- 1 つの spec で複数機能の変更系操作を連鎖させ、失敗時復旧を困難にしてはならない
- CSS クラス名や DOM 階層へ過度依存してはならない
- 失敗証跡を無効化したまま CI 常設運用してはならない
- E2E 都合で Electron アプリ本体の main / preload / renderer 境界を崩してはならない

### 12.2 運用上の禁止事項（MUST NOT）

- アカウントロックや MFA 変更など、運用側へ影響する設定変更を E2E から恒常的に行ってはならない
- テストデータ掃除を手作業前提にしてはならない
- flaky なシナリオを「たまに落ちるが問題ない」として放置してはならない
- artifact / report を無期限保管したり、公開範囲と削除責務が不明なまま共有してはならない
- worker 数を実認証テストアカウント数より大きく設定し、排他前提を崩したまま常設運用してはならない

---

## 13. レビュー観点（チェックリスト）

### 13.1 基本方針（セクション0, 1）

- [ ] Playwright を標準採用しているか
- [ ] テスト資産がリポジトリ直下の `/e2e` に集約されているか
- [ ] 実認証をモックやトークン注入で置き換えていないか
- [ ] Electron 操作と外部認証ブラウザ操作の責務が分離されているか
- [ ] シナリオファイルが資格情報や設定読込の責務を持っていないか
- [ ] アプリ本体の詳細規約を本規約へ重複記載せず、`electron-typescript-coding-guide.md` へ適切に委譲しているか

### 13.2 ディレクトリ構成（セクション2, 3）

- [ ] `scenarios/`, `fixtures/`, `helpers/`, `selectors/`, `config/` が役割どおり分離されているか
- [ ] ドメイン分類がアプリ本体の機能分類と対応しているか
- [ ] ファイル名が `<scenario>.spec.ts`、`<purpose>.fixture.ts`、`<feature>.selectors.ts` に従っているか
- [ ] `.auth/`, `artifacts/`, `reports/` が Git 管理対象や役割の面で適切に扱われているか
- [ ] Playwright 依存関係の正本が 1 箇所に定まり、`README.md` に配置責務が明記されているか

### 13.3 シナリオ分割（セクション4）

- [ ] 1 spec = 1 業務目的を守っているか
- [ ] 正常系 / 異常系 / 権限制御が無理なく分離されているか
- [ ] シリアル実行が必要なケースを正しく限定しているか
- [ ] helper に業務判定を押し込みすぎていないか

### 13.4 セレクタ（セクション5）

- [ ] `data-testid` を第一選択にしているか
- [ ] 表示文言依存や CSS 依存を最小化しているか
- [ ] セレクタ文字列が `selectors/` へ集約されているか
- [ ] `data-testid` 名が意味ベースで保守可能か

### 13.5 待機・リトライ（セクション6）

- [ ] 固定スリープへ依存していないか
- [ ] 認証待機タイムアウトがアプリ側の認証タイムアウト定義より短くなっていないか
- [ ] CI のリトライ回数が過剰でないか
- [ ] 再試行すべきでない失敗をリトライで隠していないか
- [ ] flaky 候補を識別・記録でき、監視ルールが別文書を含めて定義されているか

### 13.6 Electron / 実認証（セクション7, 8）

- [ ] `electron-typescript-coding-guide.md` のプロセス境界を壊していないか
- [ ] renderer が本来不要な API を E2E のために公開していないか
- [ ] 認証 URL 自体は本物を使っているか
- [ ] テストアカウントが個人利用アカウントと分離されているか
- [ ] 秘密情報がログ / trace / screenshot / レポートへ残らないか
- [ ] セッション残骸やブラウザ残骸の後始末があるか
- [ ] worker 数と実認証テストアカウント数の対応関係、共有禁止条件、`@serial` 条件が定義されているか
- [ ] `@auth` スモーク失敗時に、認証基盤障害とテスト不備を切り分ける観点があるか

### 13.7 設定・実行モード（セクション9, 10）

- [ ] アプリ設定の正本が 1 箇所に定義され、E2E がそれを参照しているか
- [ ] E2E 専用環境変数が `e2e/config/env.ts` 等で集中管理されているか
- [ ] ローカル / CI の差分が設定ファイルで吸収されているか
- [ ] `infra-python-cdk-guide.md` の原則どおり、暗黙の設定や複数正本を作っていないか
- [ ] PR 用と定期実行用の実行対象分離が考慮されているか

### 13.8 共通原則・証跡・禁止事項（セクション11, 12）

- [ ] `backend-golang-coding-guide.md` の共通原則どおり、秘密情報をログやエラーメッセージへ含めていないか
- [ ] 失敗時に trace / screenshot / video / report を保存するか
- [ ] 保存パスが衝突しにくい構成か
- [ ] 資格情報や認証コードが証跡に残らないか
- [ ] 禁止事項に該当する近道実装を採っていないか
- [ ] artifact / report の保存期間、公開範囲、削除責務が最低限定義されているか

---

## 14. 補足: レビュー優先度

| 優先度 | 観点 |
|:------:|------|
| 🔴 高 | 実認証の迂回、秘密情報漏えい、プロセス境界違反、固定スリープ依存、失敗証跡不足 |
| 🟡 中 | シナリオ分割不備、セレクタ命名不統一、ローカル / CI 設定分離不足、不要なシリアル実行 |
| 🟢 低 | タグ命名の表記揺れ、README の表現差、軽微な構成コメント不足 |
