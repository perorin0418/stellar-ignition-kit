# Python コーディング規約（CLIツール / Windows exe配布前提）

（関数型志向・運用容易性重視・Windows互換）

## 0. 基本方針（MUST）

- 本規約は **新規開発向け** の Python 製 CLI ツールを対象とする
- Python 3.12 以上を前提とする
- Python 実行環境管理は **uv を必須** とする（グローバル環境を汚さない）
- 配布形態は **Windows で実行可能な exe** を標準とする
- ビジネスロジックは **関数ベース** で実装する
- CLI（引数解析）と業務ロジック（処理本体）を分離する
- グローバル状態の書き換えは禁止
- `if __name__ == "__main__":` 以外で副作用初期化を行わない

### 0.1 uv 利用ルール（MUST）

- Python の実行・依存解決・テスト実行は `uv` 経由で行う
- グローバル環境への `pip install` は禁止する
- 仮想環境はプロジェクトローカル（`.venv/`）に限定する
- 依存ロックは `uv.lock` を標準とし、コミット対象に含める

---

## 1. CLI 設計

### 1.1 エントリーポイント（MUST）

- エントリーポイントは `main()` に集約する
- `main()` は終了コード（`int`）を返し、`SystemExit(main())` で終了する

```python
import sys


def main(argv: list[str] | None = None) -> int:
    # 引数解析 + 実行
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
```

### 1.2 引数解析（MUST）

- `argparse` を標準採用する（特別な理由がある場合のみ代替可）
- 必須引数・選択肢・デフォルト値は明示する
- `--help` で利用者が迷わない説明を提供する

### 1.3 標準出力/標準エラー（MUST）

- 正常結果は `stdout` に出力する
- エラー・警告は `stderr` に出力する
- 人向け表示と機械処理向け表示（JSON/TSV など）を混在させない

### 1.4 終了コード（MUST）

- `0`: 正常終了
- `1`: 予期可能な業務エラー（入力不正、ファイル未存在など）
- `2`: 引数エラー
- `>=10`: システムエラー（外部I/O障害など）

---

## 2. 関数設計と責務分離

### 2.1 推奨スタイル（MUST）

- ビジネスロジックはトップレベル関数で定義する
- 外部依存（ファイル、時刻、環境変数、ネットワーク）は引数注入する

```python
def run_usecase(*, read_text, now, input_path: str) -> dict:
    text = read_text(input_path)
    return {"length": len(text), "executed_at": now().isoformat()}
```

### 2.2 禁止事項（MUST NOT）

- CLI 層で直接ファイルI/O・DBアクセス・API呼び出しを行う
- 例外を握りつぶす（`except Exception: pass`）

### 2.3 非公開関数の命名（MUST）

- 非公開用途の関数は `_internal_` プレフィックスを付与する

```python
def _internal_validate_path(path: str) -> None:
    ...
```

---

## 3. 名前付け規則

### 3.1 モジュール・ファイル

- スネークケース

```text
main_cli.py
file_exporter.py
```

### 3.2 関数・変数

- スネークケース

### 3.3 クラス・例外

- パスカルケース

### 3.4 定数

- 全大文字スネークケース

```python
DEFAULT_ENCODING = "utf-8"
MAX_RETRY_COUNT = 3
```

---

## 4. Windows 互換ルール

### 4.1 パス操作（MUST）

- パス文字列の手組みを禁止し、`pathlib.Path` を使用する
- 区切り文字（`\\` / `/`）の違いを実装で吸収する

```python
from pathlib import Path

output_path = Path(base_dir) / "output" / "result.txt"
```

### 4.2 文字コード（MUST）

- ファイルI/Oは原則 `encoding="utf-8"` を明示する
- 改行コード差分を避けるため、テキスト出力時は `newline="\n"` を検討する

### 4.3 依存ライブラリ（MUST）

- OS依存のネイティブライブラリは最小限にする
- 採用時は Windows でのビルド可否・ライセンス・保守性を確認する

---

## 5. exe パッケージング規約（PyInstaller）

### 5.1 ツール選定（MUST）

- Windows 向け exe 生成は **PyInstaller** を標準採用する
- 原則 `--onefile` で配布する（起動性能要件がある場合は `--onedir` を検討）

### 5.2 ビルド方針（MUST）

- **ビルドは配布先と同じ OS（Windows）で実施する**
- ビルド成果物は `dist/` に出力する
- 実行ファイル名は用途が分かる英小文字スネークケースを採用する

### 5.3 リソース同梱（MUST）

- 設定ファイル・テンプレート・証明書など同梱物は `--add-data` で明示する
- 相対パス依存を避け、実行時の基準パスは `sys._MEIPASS` を考慮する

### 5.4 秘密情報（MUST NOT）

- APIキー、トークン、パスワードを exe に埋め込まない
- 秘密情報は環境変数または安全な外部ストアから取得する

### 5.5 配布物（MUST）

- 最低限以下を同梱する
  - exe 本体
  - `README`（利用方法、引数、終了コード）
  - `LICENSE`（依存ライブラリ含む）
  - `CHANGELOG`（任意だが推奨）

### 5.6 標準パッケージング手順（MUST）

- Windows exe 化は以下の順序で実施する

1. 依存関係を固定した状態でセットアップする
2. テストを実行して成功を確認する
3. PyInstaller で exe を生成する
4. 生成した exe のスモークテストを実施する
5. 配布物（exe / README / LICENSE）を同一バージョンでまとめる

- 代表例（`src/app/main_cli.py` を入口とする場合）

```text
uv python install 3.12
uv sync --locked
uv run pytest
uv run pyinstaller --noconfirm --clean --onefile --name my_cli_tool src/app/main_cli.py
dist\my_cli_tool.exe --help
```

- `uv` 未導入端末では、事前に `uv` の導入を行う（導入手順はプロジェクト `README` に明記する）

### 5.7 必須ビルドオプション（MUST）

- 原則として以下オプションを指定する
    - `--noconfirm`: 既存成果物の上書きを自動化する
    - `--clean`: キャッシュ影響を低減する
    - `--onefile`: 単一 exe を生成する
    - `--name <tool_name>`: 成果物名を固定化する
- アイコンや追加データが必要な場合は、`--icon` と `--add-data` を明示する

### 5.8 必須検証（MUST）

- 生成後は最低限以下を確認する
    - `--help` が終了コード `0` で完了する
    - 正常系コマンドが終了コード `0` で完了する
    - 想定異常系で終了コード `1` または `2` を返す
    - `stderr` に機密情報が出力されない

### 5.9 配布アーカイブ規約（MUST）

- 配布名は以下形式を採用する
    - `<tool_name>_<version>_windows_amd64.zip`
- ZIP には以下を格納する
    - `<tool_name>.exe`
    - `README.md`
    - `LICENSE`
- 改ざん検知のため `SHA256` ハッシュファイルを同梱または別添する

### 5.10 再現性の担保（MUST）

- 依存バージョンは `uv.lock` により固定する
- 同一 Git リビジョンから同一手順で再ビルド可能であることを保証する
- 手順は `README` または `BUILD.md` に維持し、CI でも同等手順を実行する

---

## 6. ログ・例外・メッセージ

### 6.1 ログ出力（MUST）

- `logging` モジュールを使用する
- `print` はユーザー向け最終結果表示に限定する
- 既定ログレベルは `INFO`、デバッグ時のみ `DEBUG` を許可する

### 6.2 例外設計（MUST）

- 例外は利用者向けメッセージと内部原因を分離する
- 原因例外は `raise ... from err` で保持する

```python
class CliExecutionError(Exception):
    """CLI実行エラー。"""


def execute_task(*, run_job) -> None:
    try:
        run_job()
    except OSError as err:
        raise CliExecutionError("処理の実行に失敗しました。入力ファイルを確認してください。") from err
```

### 6.3 出力制限（MUST NOT）

- ログ・エラーメッセージに以下を出力しない
  - 認証情報（API Key / Token / Password）
  - 個人情報
  - 内部ネットワーク構成

---

## 7. テスト

### 7.1 テスト方針（MUST）

- `pytest` を標準採用する
- 主要ロジックはパラメタライズテストで記述する
- 正常系・異常系を同一テーブルで明示する

```python
import pytest


@pytest.mark.parametrize(
    "argv,want_code",
    [
        (["--input", "ok.txt"], 0),
        (["--input", "missing.txt"], 1),
    ],
)
def test_main(argv, want_code):
    assert main(argv) == want_code
```

### 7.2 テスト対象の分離（MUST）

- CLI層のテスト（引数、終了コード、標準出力/標準エラー）
- ロジック層のテスト（純粋関数ベース）
- exe スモークテスト（生成物が起動し、`--help` が成功すること）

---

## 8. 依存関係・プロジェクト構成

### 8.1 推奨ディレクトリ構成

```text
cli-tool/
    src/
        app/                # CLI入口
        usecase/            # ユースケース
        domain/             # 業務ルール
        infrastructure/     # ファイル/外部I/O
    tests/
    pyproject.toml
    README.md
```

### 8.2 依存ルール（MUST）

- `app` は `usecase` のみ依存可
- `usecase` は `domain` / `infrastructure` に依存可
- `domain` は外部I/Oに依存しない

---

## 9. レビュー観点（チェックリスト）

### 9.1 基本方針（セクション0）

- [ ] Python 3.12 以上を前提としているか
- [ ] 実行環境管理が `uv` 前提になっているか
- [ ] CLI層とロジック層が分離されているか
- [ ] グローバル状態を書き換えていないか

### 9.2 CLI設計（セクション1）

- [ ] `main()` が終了コード（int）を返しているか
- [ ] 引数定義・`--help` が明確か
- [ ] `stdout` と `stderr` の用途分離ができているか

### 9.3 Windows互換（セクション4）

- [ ] `pathlib.Path` を使用しているか
- [ ] 文字コードが明示されているか
- [ ] Windows での動作確認が実施されているか

### 9.4 exeパッケージング（セクション5）

- [ ] PyInstaller 設定が再現可能か
- [ ] `uv sync --locked` を前提にビルドしているか
- [ ] 標準パッケージング手順（5.6）に従っているか
- [ ] 必須ビルドオプション（5.7）が満たされているか
- [ ] 必須検証（5.8）が記録されているか
- [ ] 同梱リソースの解決方法が実装されているか
- [ ] 配布アーカイブ規約（5.9）に従っているか
- [ ] 秘密情報がバイナリに埋め込まれていないか

### 9.5 ログ・例外（セクション6）

- [ ] `raise ... from err` で原因を保持しているか
- [ ] 機密情報がログ/エラーメッセージに含まれていないか

### 9.6 テスト（セクション7）

- [ ] 正常系・異常系テストがあるか
- [ ] CLI層とロジック層の双方を検証しているか
- [ ] exe スモークテストがあるか

---

## 補足: レビュー優先度

| 優先度 | 観点 |
|:------:|------|
| 🔴 高 | 秘密情報埋め込み、終了コード不整合、Windows非互換、テスト未実装 |
| 🟡 中 | 依存分離違反、ログ設計不備、引数設計不足 |
| 🟢 低 | 命名規約、フォーマット、コメント不足 |
