package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// main はエクスポートインターフェース生成ツールのエントリポイントです。
// Goファイルからトップレベルの関数変数のシグネチャを抽出し、
// Export構造体とコンストラクタを新しいGoファイルとして出力します（DIやテストで利用）。
func main() {
	// 入力となるGoファイルのパス取得
	inputFile := flag.String("in", "", "input .go file (required)")
	flag.Parse()
	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "-in フラグは必須です")
		os.Exit(1)
	}

	// ファイルパス関連の準備（出力ファイルは入力ファイルの親ディレクトリに配置）
	absInput, _ := filepath.Abs(*inputFile)
	inDir := filepath.Dir(absInput)
	base := strings.TrimSuffix(filepath.Base(absInput), filepath.Ext(absInput))
	// 出力先は入力Goファイルの2つ上のディレクトリ(../../export)にする
	grandParentDir := filepath.Dir(filepath.Dir(inDir))
	exportDir := filepath.Join(grandParentDir, "export")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "exportディレクトリの作成に失敗しました: %v\n", err)
		os.Exit(1)
	}
	outFile := filepath.Join(exportDir, base+"_export.go")

	// 親ディレクトリからパッケージ命名規則(usecase/domain/logic)に一致するものを探索
	pkgNames := []string{"usecase", "domain", "logic"}
	_, pkgName := findParentDirForPackage(inDir, pkgNames)
	if pkgName == "" {
		fmt.Fprintln(os.Stderr, "パッケージフォルダ(usecase/domain/logic)が見つかりませんでした。")
		os.Exit(1)
	}

	// go.modのルートを見つけてモジュール名+サブパスでインポートパスを生成
	goModRoot := findGoModRoot(inDir)
	if goModRoot == "" {
		fmt.Fprintln(os.Stderr, "go.modが見つかりませんでした。")
		os.Exit(1)
	}
	// go.modからモジュール名を取得し、importパスを正規化
	modBytes, err := os.ReadFile(filepath.Join(goModRoot, "go.mod"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "go.modの読み込みに失敗しました。")
		os.Exit(1)
	}
	lines := strings.Split(string(modBytes), "\n")
	moduleName := ""
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "module ") {
			moduleName = strings.TrimSpace(strings.TrimPrefix(l, "module "))
			break
		}
	}
	if moduleName == "" {
		fmt.Fprintln(os.Stderr, "go.modにmodule宣言が見つかりませんでした。")
		os.Exit(1)
	}
	relImportPath, _ := filepath.Rel(goModRoot, inDir)
	importPath := moduleName
	if relImportPath != "." {
		importPath = toPathSlash(moduleName + "/" + relImportPath)
	}
	importAlias := filepath.Base(inDir)

	// 入力GoファイルをASTとしてパース
	fs := token.NewFileSet()
	af, err := parser.ParseFile(fs, absInput, nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "パースエラー: %v\n", err)
		os.Exit(1)
	}

	// 抽出した関数変数の情報を保持する内部構造体
	type FuncVar struct {
		Name, TypeString string
	}
	var funcVars []FuncVar

	// 宣言（Decls）からトップレベルの関数変数（無名関数）だけを抽出
	// ASTから直接関数シグネチャを抽出（型チェッカーに依存しない）
	for _, decl := range af.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range vs.Names {
				if len(vs.Values) > i {
					if flit, ok := vs.Values[i].(*ast.FuncLit); ok {
						// ASTのFuncTypeから関数シグネチャを生成
						sig := funcTypeToString(fs, flit.Type)
						funcVars = append(funcVars, FuncVar{name.Name, sig})
					}
				}
			}
		}
	}

	if len(funcVars) == 0 {
		fmt.Fprintln(os.Stderr, "トップレベルの関数変数が見つかりませんでした。")
		os.Exit(1)
	}

	// エクスポート用構造体名を作成
	structName := base + "Export"
	var b bytes.Buffer
	// パッケージ、インポート、構造体定義を出力内容に追加
	importBlock, hasImportBlock, existingAlias, hasImportPath, err := readExistingImportBlock(outFile, importPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "既存ファイルのimport解析に失敗しました: %v\n", err)
		os.Exit(1)
	}
	if existingAlias != "" {
		importAlias = existingAlias
	}
	fmt.Fprintf(&b, "package %s\n\n", pkgName)
	if hasImportBlock {
		fmt.Fprintf(&b, "%s\n\n", strings.TrimRight(importBlock, "\n"))
		if !hasImportPath {
			fmt.Fprintf(&b, "import %s \"%s\"\n\n", importAlias, importPath)
		}
	} else {
		fmt.Fprintf(&b, "import %s \"%s\"\n\n", importAlias, importPath)
	}
	fmt.Fprintf(&b, "type %s struct {\n", structName)
	for _, fv := range funcVars {
		fmt.Fprintf(&b, "\t%s %s\n", fv.Name, fv.TypeString)
	}
	fmt.Fprintf(&b, "}\n\nfunc New%s() %s {\n\treturn %s{\n", upperCamel(base), structName, structName)
	qualifier := importAlias
	if qualifier == "." {
		qualifier = ""
	}
	for _, fv := range funcVars {
		if qualifier == "" {
			fmt.Fprintf(&b, "\t\t%s: %s,\n", fv.Name, fv.Name)
		} else {
			fmt.Fprintf(&b, "\t\t%s: %s.%s,\n", fv.Name, qualifier, fv.Name)
		}
	}
	fmt.Fprintf(&b, "\t}\n}\n")

	// 出力Goコードを整形してファイル出力
	formatted, err := format.Source(b.Bytes())
	if err != nil {
		// 整形エラー時は生成途中の内容を一時ディレクトリに出力
		tmpDir := os.TempDir()
		tmpFile := filepath.Join(tmpDir, base+"_export_failed.go")
		_ = os.WriteFile(tmpFile, b.Bytes(), 0644)
		fmt.Fprintf(os.Stderr, "整形エラー: %v\n生成途中のファイルを %s に出力しました。\n", err, tmpFile)
		os.Exit(1)
	}
	if err := os.WriteFile(outFile, formatted, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "ファイル書き込みエラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("エクスポートファイルを生成しました: %s (出力先: ../../export)\n", outFile)
}

// findParentDirForPackageは、startディレクトリから上位に辿って、targetsで指定したディレクトリ名に一致するパッケージディレクトリを探します。
// 該当ディレクトリのパスと名称を返します。見つからない場合は空文字列を返します。
func findParentDirForPackage(start string, targets []string) (string, string) {
	dir := start
	for {
		base := filepath.Base(dir)
		for _, t := range targets {
			if base == t {
				return dir, t
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ""
		}
		dir = parent
	}
}

// findGoModRootはstartディレクトリから上位に移動しながら、go.modファイルのあるモジュールルートを探します。
// 見つかればそのディレクトリパス、見つからなければ空文字列を返します。
func findGoModRoot(start string) string {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// upperCamelは与えられた文字列の先頭だけ大文字に変換して返します。
func upperCamel(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// toPathSlashは、パス区切りをスラッシュに統一します（Goのimportパス用）。
func toPathSlash(s string) string {
	return filepath.ToSlash(s)
}

// funcTypeToStringは、*ast.FuncTypeをGoの関数型文字列（例: func(int, string) error）に変換します。
func funcTypeToString(fs *token.FileSet, ft *ast.FuncType) string {
	var buf bytes.Buffer
	// go/printerを使ってAST→ソースコード文字列に変換
	if err := printer.Fprint(&buf, fs, ft); err != nil {
		return "func()"
	}
	return buf.String()
}

// readExistingImportBlockは、既存ファイルのimport宣言をそのまま取得します。
// importPathが含まれている場合は対応するエイリアスも返します。
func readExistingImportBlock(pathToFile, importPath string) (string, bool, string, bool, error) {
	src, err := os.ReadFile(pathToFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, "", false, nil
		}
		return "", false, "", false, err
	}
	fs := token.NewFileSet()
	af, err := parser.ParseFile(fs, pathToFile, src, parser.ParseComments)
	if err != nil {
		return "", false, "", false, err
	}

	var (
		startOffset   = -1
		endOffset     = -1
		existingAlias string
		hasImportPath bool
	)

	for _, decl := range af.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}
		start := fs.Position(gen.Pos()).Offset
		end := fs.Position(gen.End()).Offset
		if startOffset == -1 || start < startOffset {
			startOffset = start
		}
		if endOffset == -1 || end > endOffset {
			endOffset = end
		}
		for _, spec := range gen.Specs {
			is, ok := spec.(*ast.ImportSpec)
			if !ok || is.Path == nil {
				continue
			}
			pathValue := strings.Trim(is.Path.Value, "\"")
			if pathValue == importPath {
				hasImportPath = true
				if is.Name != nil {
					existingAlias = is.Name.Name
				} else if existingAlias == "" {
					existingAlias = path.Base(pathValue)
				}
			}
		}
	}

	if startOffset == -1 || endOffset == -1 {
		return "", false, existingAlias, hasImportPath, nil
	}

	return string(src[startOffset:endOffset]), true, existingAlias, hasImportPath, nil
}
