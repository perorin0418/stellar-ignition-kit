# WPF + CommunityToolkit.Mvvm + UIA 開発規約

## 1. 目的
本規約は、WPF アプリケーションを **CommunityToolkit.Mvvm / UI Automation / 自動ビルド** を用いて開発・運用する際の
**標準的な構成・設計・運用ルール**を定め、

- テスト容易性
- 保守性
- CI/CD 安定性
- チーム開発での一貫性

を最大化することを目的とする。

---

## 2. 基本方針（原則）

### 2.1 ソリューション構成
- **1アプリケーション = 1 Solution（.sln / .slnx）**
- 実行物・テスト・ビルドスクリプトは **同一 Solution で管理**する

### 2.2 責務分離
- UI / MVVM / ドメイン / 外部依存は **必ず物理分離（別 csproj）**する
- 依存方向は **一方向のみ**とする

```
App → Presentation → Domain → Infrastructure
```

逆依存は禁止する。

---

## 3. プロジェクト構成

```
StarGazer.sln(x)
│
├─ Src/
│  ├─ StarGazer.App.Main/StarGazer.Main.csproj
│  ├─ StarGazer.App.業務コード/StarGazer.App.業務コード.csproj
│  ├─ StarGazer.Presentation.業務コード/StarGazer.Presentation.業務コード.csproj
│  ├─ StarGazer.Domain.業務コード/StarGazer.Domain.業務コード.csproj
│  └─ StarGazer.Infrastructure.業務コード/StarGazer.Infrastructure.業務コード.csproj
│
├─ Tests/
│  ├─ StarGazer.Unit.業務コード/StarGazer.Unit.業務コード.csproj
│  └─ StarGazer.Uia.業務コード/StarGazer.Uia.業務コード.csproj
│
├─ Build/
│  ├─ StarGazer.Compile/compile.bat
│  └─ StarGazer.Publish/publish.bat
│
└─ CiCd/
   ├─ StarGazer.Bootstrap/build.bat
   ├─ StarGazer.Deploy/deploy.exe
   └─ StarGazer.Update/update.exe
```

---

## 4. プロジェクト別規約

### 4.1 StarGazer.App.Main（エントリポイント、WPFアプリケーション）
| 項目 | 規約 |
|---|---|
| テンプレート | WPF アプリケーション |
| TargetFramework | net10.0-windows |
| 依存 | StarGazer.App |
| 役割 | エントリポイント、App.xaml |

**ルール**
- MainWindow.xaml は配置する。
- 業務ロジックは一切記述しない。

### 4.2 StarGazer.App（WPF クラスライブラリ）

| 項目 | 規約 |
|---|---|
| テンプレート | WPF クラスライブラリ |
| TargetFramework | net10.0-windows |
| UseWPF | true |
| 役割 | View |

**禁止事項**
- ビジネスロジックの記述
- 直接的なデータアクセス

View（XAML）はすべてこのプロジェクトに配置する。

#### ディレクトリ構成例

```
StarGazer.App.業務コード/
├─ Views/ ・・・ 画面定義（XAML + CodeBehind）
│  └─ 各種 View 定義
└─ Resources/ ・・・ アプリケーション リソース定義（画像、アイコンなど）
   └─ 各種 リソース定義
```

---

### 4.3 StarGazer.Presentation（MVVM）

| 項目 | 規約 |
|---|---|
| テンプレート | クラス ライブラリ |
| TargetFramework | net10.0 |
| 依存 | CommunityToolkit.Mvvm |

**ルール**
- XAML を置かない
- WPF 型（Window / Brush / Visibility 等）を使用しない
- ViewModel は UI 非依存とする

#### ディレクトリ構成例

```
StarGazer.Presentation.業務コード/
├─ ViewModels/ ・・・ 画面単位の状態管理・UIロジック（ObservableProperty、RelayCommand、画面状態、Domain 呼び出し、例外処理の制御など）
│  └─ 各種 ViewModel 定義
├─ Behaviors/ ・・・ XAML から使う振る舞い。
│  └─ 各種 ビヘイビア定義
└─ Validation/ ・・・ 入力値検証ロジック定義（※UI入力制約 → Presentation、業務整合性 → Domain に分離）
   └─ 各種 バリデーション定義
```

---

### 4.4 StarGazer.Domain（ドメイン）

| 項目 | 規約 |
|---|---|
| テンプレート | クラス ライブラリ |
| TargetFramework | net10.0 |
| 依存 | 原則なし |

- 業務ロジックのみを配置
- フレームワーク非依存

#### ディレクトリ構成例

```
StarGazer.Domain.業務コード/
├─ Entities/ ・・・ 業務オブジェクト（ID を持つ、状態変更メソッドを持つ、ビジネスルールを内包する）
|  └─ 各種 エンティティ定義
└─ Services/ ・・・ ドメインサービス（複数エンティティに跨るビジネスロジック、外部依存を必要としない純粋な業務ロジック）
   └─ 各種 ドメインサービス定義
```

---

### 4.5 StarGazer.Infrastructure（外部依存）

| 項目 | 規約 |
|---|---|
| テンプレート | クラス ライブラリ |
| TargetFramework | net10.0 |
| 依存 | 必要に応じて追加 |

- API / ファイル / DB / 設定取得などを担当
- Domain に定義された Interface を実装

#### ディレクトリ構成例

```
StarGazer.Infrastructure.業務コード/
├─ Api/
│  ├─ Clients/ ・・・ REST / GraphQL / gRPC クライアント
│  │  └─ 各種 API クライアント実装
│  └─ Models/ ・・・ DTO、API モデル定義
│     └─ 各種 API モデル定義
├─ FileSystem/ ・・・ ファイル読み書き、ファイル監視
│  └─ ファイル操作関連実装
├─ Configuration/ ・・・ 環境変数、設定ファイル読み込み
│  └─ 設定取得関連実装
├─ Logging/ ・・・ ロギング
│  └─ ロギング関連実装
├─ Security/ ・・・ 暗号化、トークン保存
│  └─ セキュリティ関連実装
├─ Time/ ・・・ 現在時刻取得の抽象化
│  └─ 日時関連実装
└─ Extensions/ ・・・ 共通拡張メソッド、内部ヘルパー
   └─ 各種拡張メソッド

```

---

## 5. テスト規約

### 5.1 単体テスト

- テンプレート: MSTest テスト プロジェクト
- 対象: Domain / Presentation
- TargetFramework: net10.0-windows
- App（WPF）は直接参照しない

---

### 5.2 UI Automation テスト（UIA）

- テンプレート: MSTest テスト プロジェクト
- TargetFramework: net10.0-windows
- Page Object パターンを採用

**必須ルール**
- View には AutomationId を必ず付与

---

## 6. 自動ビルド（CI/CD）

### 6.1 基本方針

- Solution 単位でビルド
- ローカルと CI のビルド手順を統一

```
build\build.bat
```

### 7.2 デプロイ

```
tools\auto_deploy\auto_deploy.exe
```

---

## 8. 禁止事項まとめ

- Presentation に WPF 依存を入れる
- Domain に UI / Framework 依存を入れる
- Tests から App を直接参照する
- slnx の GUID に意味を持たせる
