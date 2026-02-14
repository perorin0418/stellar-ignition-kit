# Python コーディング規約（バックエンド）

（関数型志向・テスト容易性重視・最小限の OOP）

## 0. 基本方針（MUST）

- 本規約は **新規開発向け** の Python バックエンドを対象とする
- Python 3.12 以上を前提とする
- **純粋関数志向** と **テスト容易性** を最優先する
- ビジネスロジックでの **クラス設計は原則禁止**
  - 例外はデータ保持の `dataclass` のみ
- グローバル状態の書き換えは禁止
- `__init__.py` での副作用初期化は禁止

---

## 1. OOP 制限と関数設計

### 1.1 禁止事項（MUST NOT）

- ビジネスロジックをクラスのメソッドとして実装してはならない

```python
# ❌ 禁止
class OrderService:
    def calculate_total(self, order):
        ...
```

### 1.2 推奨スタイル（MUST）

- ビジネスロジックは **トップレベル関数** で定義する

```python
# ✅ 推奨
def calculate_total(order):
    ...
```

---

## 2. テスト容易性のための関数設計

### 2.1 外部依存の注入（MUST）

- 外部依存（DB/API/ファイル/時刻/乱数/環境変数）は **引数として受け取る**
- 結果は **戻り値として返す**

```python
def calculate_result(*, now, fetch_user, input_data):
    user = fetch_user(input_data.user_id)
    return {
        "now": now(),
        "user": user,
    }
```

### 2.2 非公開関数の命名（MUST）

- 非公開用途の関数は **`_internal_` プレフィックス** を付与する

```python
def _internal_validate_input(input_data):
    ...
```

### 2.3 許可されるクラス利用（MAY）

- `@dataclass(frozen=True)` による **不変データモデル**
- 例外クラスの定義

```python
from dataclasses import dataclass

@dataclass(frozen=True)
class Order:
    id: str
    amount: int
```

---

## 3. 名前付け規則

### 3.1 モジュール・ファイル

- **スネークケース**

```text
order_service.py
payment_gateway.py
```

### 3.2 関数・変数

- **スネークケース**

```python
retry_count = 3

def fetch_user(user_id):
    ...
```

### 3.3 クラス・例外

- **パスカルケース**

```python
class OrderNotFoundError(Exception):
    ...
```

### 3.4 定数

- **全大文字スネークケース**

```python
DEFAULT_TIMEOUT_SEC = 30
```

---

## 4. 純粋関数を目指した設計

### 4.1 基本方針

- 可能な限り **純粋関数（副作用なし）** を目指す
- 副作用は **境界層** に閉じ込める

### 4.2 外部依存の扱い（MUST）

- 依存は **関数引数で渡す**
- 依存の参照は **引数経由のみ**

```python
def create_invoice(*, now, generate_id, customer):
    return {
        "id": generate_id(),
        "customer": customer,
        "created_at": now(),
    }
```

---

## 5. エラーハンドリング

### 5.1 基本ルール（MUST）

- 例外は **明示的に補足** して意図を示す
- 例外は **`raise ... from err`** で原因を保持する

```python
def load_user(*, fetch_user, user_id):
    try:
        return fetch_user(user_id)
    except RuntimeError as err:
        raise UserLoadError("ユーザー取得に失敗しました") from err
```

### 5.2 エラーメッセージ指針（MUST）

- ユーザー向けメッセージは **意味が理解できる日本語**
- 以下の情報は **絶対に出力してはならない**
  - 認証情報（API Key / Token / Password）
  - 個人情報
  - 内部ネットワーク構成

---

## 6. コメント・ドキュメント

### 6.1 コメント必須（MUST）

- 公開関数・公開クラスには **docstring を必ず付与**
- 日本語で **引数・戻り値・例外** を説明する

```python
def create_report(*, fetch_data, report_id):
    """
    create_report はレポートを生成する。

    引数:
        fetch_data: データ取得関数
        report_id: 取得対象 ID
    戻り値:
        生成したレポート
    例外:
        ReportCreateError: 生成に失敗した場合
    """
    ...
```

---

## 7. フォーマット・静的解析

- `black` / `ruff` を使用できる場合は **必須**
- 行長は `black` の設定に従う（未設定の場合は 88 文字）

---

## 8. 依存関係・構成

### 8.1 依存関係の考え方（MUST）

- アプリケーション層は **インフラ詳細に依存しない**
- 外部サービスとの境界は **関数引数で注入**

### 8.2 依存ライブラリのルール（MUST）

- 標準ライブラリを最優先する
- 追加依存は **最小限** とし、必ず根拠を示す
- バージョンは **固定** する

---

## 9. ディレクトリ構成

### 9.1 基本方針（MUST）

- 構成は **責務単位** に分割し、依存方向は **内側（ドメイン）へ** のみ許可する
- 入口（CLI/API）とロジックを分離し、外部依存は境界層に閉じ込める

### 9.2 推奨ツリー（例）

```text
backend/
    app/                # 入口（CLI/API）
    usecase/            # ユースケース
    domain/             # ドメイン（業務ルール）
    logic/              # 汎用ロジック
    model/              # DTO/データ定義
    infrastructure/     # DB/API など外部依存
    tests/              # テスト
```

### 9.3 依存ルール（MUST）

- `app` -> `usecase` のみ依存可
- `usecase` -> `domain` / `logic` / `model` のみ依存可
- `domain` -> `logic` / `model` / `infrastructure` のみ依存可
- `logic` -> `model` のみ依存可
- `model` / `infrastructure` は他層に依存しない

---

## 10. テスト

### 9.1 テスト必須範囲（MUST）

- 主要ロジックは **pytest のパラメタライズ** で記述する
- 正常系・異常系を同一テーブルで明示する

```python
import pytest

@pytest.mark.parametrize(
    "input_data, setup, want, want_error",
    [
        ("1", lambda: None, {"ok": True}, False),
        ("2", lambda: (_ for _ in ()).throw(RuntimeError("api error")), None, True),
    ],
)
def test_process_queue(input_data, setup, want, want_error):
    def fetch():
        setup()
        return {"ok": True}

    if want_error:
        with pytest.raises(RuntimeError):
            process_queue(fetch, input_data)
        return

    assert process_queue(fetch, input_data) == want
```

### 9.2 外部依存のモック化（MUST）

- 外部依存は **関数として注入** し、直接置き換える
- `unittest.mock.patch` などの **パッチ系モックは原則禁止**

```python
def process_queue(fetch_from_api, queue_id):
    return fetch_from_api(queue_id)

# テスト側で差し替え
process_queue(lambda _id: {"ok": True}, "1")
```

---

## 11. レビュー観点（チェックリスト）

### 11.1 基本方針（セクション0）

- [ ] Python 3.12 以上を前提としているか
- [ ] ビジネスロジックがクラスメソッドで実装されていないか
- [ ] グローバル状態を書き換えていないか

### 11.2 関数設計（セクション1, 2, 4）

- [ ] ビジネスロジックはトップレベル関数か
- [ ] 外部依存が引数注入されているか
- [ ] 非公開関数に `_internal_` プレフィックスがあるか

### 11.3 名前付け規則（セクション3）

- [ ] ファイル名がスネークケースか
- [ ] 関数・変数がスネークケースか
- [ ] 定数が全大文字スネークケースか

### 11.4 エラーハンドリング（セクション5）

- [ ] `raise ... from err` で原因が保持されているか
- [ ] 機密情報がメッセージに含まれていないか

### 11.5 コメント・ドキュメント（セクション6）

- [ ] 公開関数に docstring があるか
- [ ] 引数・戻り値・例外が説明されているか

### 11.6 フォーマット（セクション7）

- [ ] `black` / `ruff` が適用されているか
- [ ] 行長が `black` 設定に一致しているか（未設定の場合は 88 文字）

### 11.7 テスト（セクション10）

### 11.8 ディレクトリ構成（セクション9）

- [ ] 責務単位で分割されているか
- [ ] 依存方向が内側（ドメイン）へ向いているか

- [ ] 主要ロジックにパラメタライズテストがあるか
- [ ] 正常系・異常系が同一テーブルで明示されているか
- [ ] 外部依存が関数注入で差し替え可能か

---

## 補足: レビュー優先度

| 優先度 | 観点 |
|:------:|------|
| 🔴 高 | 外部依存注入漏れ、機密情報漏洩、テスト未実装 |
| 🟡 中 | 命名規則違反、docstring不足、OOP逸脱 |
| 🟢 低 | フォーマット、行長超過 |
