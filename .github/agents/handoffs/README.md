# Agent Handoff YAML Files

エージェント間の契約は、このディレクトリ配下に作成する YAML ファイルで授受する。

- 配置先: `.github/agents/handoffs/<request_id>/`
- `request_id` の正式形式: `REQ-YYYYMMDD-HHMMSS-NNN`
- `YYYYMMDD-HHMMSS` は採番時刻（ローカル時刻、秒まで）を表し、`NNN` は同一秒内での 3 桁連番とする。
- `request_id` ごとにディレクトリを分け、同一ワークフローの契約をまとめる。
- handoff ファイルの具体的なファイル名規則、slug 生成、衝突回避は `.github/agents/scripts/new-handoff-files.ps1` の実装を正本とする。本文から人や AI が手動でファイル名を組み立ててはならない。
- 新規 `request_id` の採番と handoff ディレクトリ作成には `.github/agents/scripts/new-request-id.ps1` を使用する。
- `.github/agents/scripts/new-request-id.ps1` は未使用の `request_id` を払い出し、対応する `.github/agents/handoffs/<request_id>/` を作成してから結果を返す。
- `.github/agents/scripts/new-request-id.ps1` の標準出力は JSON とし、`request_id` と `handoff_directory` を返す。
- `request_file` / `response_file` の採番と空ファイル生成には `.github/agents/scripts/new-handoff-files.ps1` を使用する。
- `.github/agents/scripts/new-handoff-files.ps1` は既存の `request_id` ディレクトリ配下に request / response の空 YAML ファイルを生成し、`request_file` / `response_file` / `request_path` / `response_path` などを JSON として返す。
- 上位エージェントは `.github/agents/scripts/new-handoff-files.ps1` が生成した request YAML ファイルへ内容を書き込んでから委譲する。
- 下位エージェントは `.github/agents/scripts/new-handoff-files.ps1` が生成した response YAML ファイルへ内容を書き込んでから、そのファイルパスのみを返す。
- 再委譲時は既存ファイルを上書きせず、`.github/agents/scripts/new-handoff-files.ps1` を再実行して新しい request / response ペアを作成する。
- 後続エージェントへは `context.prior_output_files` で response YAML ファイルパスを引き渡す。

本文中に YAML を直接埋め込むだけの受け渡しは禁止する。