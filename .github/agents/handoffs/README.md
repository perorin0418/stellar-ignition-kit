# Agent Handoff YAML Files

エージェント間の契約は、このディレクトリ配下に作成する YAML ファイルで授受する。

- 配置先: `.github/agents/handoffs/<request_id>/`
- request ファイル: `<step>-<agent-slug>-run<nn>.request.yaml`
- response ファイル: `<step>-<agent-slug>-run<nn>.response.yaml`
- `run<nn>` は同一 `request_id` かつ同一 `<step>-<agent-slug>` 内の実行連番とし、初回は `run01`、ユーザー確認後の再委譲などは `run02`, `run03` ... を使う。
- `request_id` ごとにディレクトリを分け、同一ワークフローの契約をまとめる。
- 上位エージェントは request YAML を作成してから委譲する。
- 下位エージェントは response YAML を作成してから、そのファイルパスのみを返す。
- 再委譲時は既存ファイルを上書きせず、新しい `run<nn>` の request / response ペアを作成する。
- 後続エージェントへは `context.prior_output_files` で response YAML ファイルパスを引き渡す。

本文中に YAML を直接埋め込むだけの受け渡しは禁止する。