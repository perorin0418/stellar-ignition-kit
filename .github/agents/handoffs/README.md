# Agent Handoff YAML Files

エージェント間の契約は、このディレクトリ配下に作成する YAML ファイルで授受する。

- 配置先: `.github/agents/handoffs/<request_id>/`
- `request_id` の正式形式: `REQ-YYYYMMDD-HHMMSS-NNN`
- `YYYYMMDD-HHMMSS` は採番時刻（ローカル時刻、秒まで）を表し、`NNN` は同一秒内での 3 桁連番とする。
- request ファイル: `h<handoff-nnn>-<step>-<agent-slug>-run<nn>.request.yaml`
- response ファイル: `h<handoff-nnn>-<step>-<agent-slug>-run<nn>.response.yaml`
- `h<handoff-nnn>` は同一 `request_id` 内での handoff 発行順を表す通し番号とし、初回は `h001`、以降は `h002`, `h003` ... を使う。step が戻った場合でも、この番号で時系列順を追える。
- `run<nn>` は同一 `request_id` かつ同一 `<step>-<agent-slug>` 内の実行連番とし、初回は `run01`、ユーザー確認後の再委譲などは `run02`, `run03` ... を使う。
- `request_id` ごとにディレクトリを分け、同一ワークフローの契約をまとめる。
- `request_id` 自体が日時ベースの識別子を持つため、ファイル名には日時を重複して入れず、追跡性向上には `h<handoff-nnn>` を優先する。
- 新規 `request_id` の採番と handoff ディレクトリ作成には `.github/agents/scripts/new-request-id.ps1` を使用する。
- `.github/agents/scripts/new-request-id.ps1` は未使用の `request_id` を払い出し、対応する `.github/agents/handoffs/<request_id>/` を作成してから結果を返す。
- `.github/agents/scripts/new-request-id.ps1` の標準出力は既定で `request_id` の 1 行テキストとし、`-AsJson` 指定時は `request_id` と `handoff_directory` を含む JSON を返す。
- 上位エージェントは request YAML を作成してから委譲する。
- 下位エージェントは response YAML を作成してから、そのファイルパスのみを返す。
- 再委譲時は既存ファイルを上書きせず、新しい `h<handoff-nnn>` と必要に応じて新しい `run<nn>` の request / response ペアを作成する。
- 後続エージェントへは `context.prior_output_files` で response YAML ファイルパスを引き渡す。

本文中に YAML を直接埋め込むだけの受け渡しは禁止する。