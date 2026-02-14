---
name: go-export-interface
description: |
  自動的にinterfaceファイルを生成する。usecase/*/*/*.go、domain/*/*/*.go、logic/*/*/*.go これらのパス配下にある.goファイルが更新（追加または修正）されると、scripts/go-export-interface.exeをそのファイルを引数にして実行する。生成されたinterfaceファイルにより、型の一貫性とリファクタリングの容易さを担保する。Goプロジェクトにおいてレイヤ間の依存性制御やテスト容易性の向上に貢献する。
---

# go-export-interface Skill

このスキルは、特定ディレクトリ（usecase/*/*/*.go, domain/*/*/*.go, logic/*/*/*.go）内のGoファイルが更新された際に、自動でscripts/go-export-interface.exeを実行し、対応するinterfaceファイルを生成します。

## 使い方（発動条件）
- usecase/*/*/*.go、domain/*/*/*.go、logic/*/*/*.go に該当する.goファイルが**新規作成**または**変更**された場合、自動的に本Skillがトリガーされます。

## 実行内容・フロー
1. 変更された全ての.goファイルに対し、下記コマンドを実行：

   ```bash
   .\.opencode\skills\go-export-interface\scripts\go-export-interface.exe -in <更新したGoファイルのパス>
   ```

2. コマンドの引数には**変更されたGoファイル（1ファイルずつ）**のパスが渡されます
   - 例: .\.opencode\skills\go-export-interface\scripts\go-export-interface.exe -in usecase/user/auth/create.go

3. コマンドが成功すれば、対応するinterfaceファイルが自動生成されます
   - 出力例: usecase/user/auth/create_interface.go

4. 生成/更新後、interfaceファイルを必要な箇所で利用・インポートしてください。

5. 依存関係のチェックは、下記コマンドで実行できます：

   ```bash
   .\.opencode\skills\go-export-interface\scripts\dependency-check.exe
   ```

## メリット
- 手書きによるinterface管理の煩雑さを削減
- 型の一貫性確保および依存制御の最適化
- テストコード作成やモック生成の効率向上

## 注意事項
- ルートパスやディレクトリ構成により修正が必要な場合、Skill本体またはスクリプト仕様もあわせて調整してください。

## 拡張例
- 追加オプションや複数ファイル同時処理への対応は、スクリプト側のバージョンアップで適宜拡張してください。

---
このSKILL.mdは「go-export-interface」スキルの自動発火・運用のためのガイドです。