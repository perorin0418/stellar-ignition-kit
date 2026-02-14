package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type violation struct {
	filePath   string
	lineNumber int
	importPath string
}

func main() {
	backendDir, repoRoot, err := findBackendDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	moduleName, err := readModuleName(filepath.Join(backendDir, "go.mod"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	importPattern := fmt.Sprintf("\"%s/(domain|logic|usecase)/[^\"]+\"", regexp.QuoteMeta(moduleName))
	re, err := regexp.Compile(importPattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "importパターンの作成に失敗しました")
		os.Exit(1)
	}

	violations, err := scanForViolations(backendDir, repoRoot, re)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if len(violations) == 0 {
		fmt.Println("依存禁止呼び出しは検出されませんでした。")
		return
	}

	fmt.Println("依存禁止呼び出しが検出されました:")
	for _, v := range violations {
		fmt.Printf("- %s:%d %s\n", v.filePath, v.lineNumber, v.importPath)
	}
	os.Exit(1)
}

func findBackendDir() (string, string, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("作業ディレクトリの取得に失敗しました: %w", err)
	}

	dir := start
	for {
		candidate := filepath.Join(dir, "backend")
		goModPath := filepath.Join(candidate, "go.mod")
		if _, statErr := os.Stat(goModPath); statErr == nil {
			return candidate, dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", "", fmt.Errorf("backend/go.mod が見つかりませんでした")
}

func readModuleName(goModPath string) (string, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("go.mod の読み込みに失敗しました: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("go.mod に module 宣言が見つかりませんでした")
}

func scanForViolations(backendDir, repoRoot string, re *regexp.Regexp) ([]violation, error) {
	var violations []violation

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "vendor" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(info.Name()) != ".go" {
			return nil
		}
		if strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		file, openErr := os.Open(path)
		if openErr != nil {
			return fmt.Errorf("ファイルを開けませんでした: %w", openErr)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			matches := re.FindAllString(line, -1)
			if len(matches) == 0 {
				continue
			}
			for _, match := range matches {
				importPath := strings.Trim(match, "\"")
				if strings.Contains(importPath, "/export") {
					continue
				}
				if strings.HasSuffix(info.Name(), "_export.go") {
					continue
				}
				relPath, relErr := filepath.Rel(repoRoot, path)
				if relErr != nil {
					relPath = path
				}
				violations = append(violations, violation{
					filePath:   filepath.ToSlash(relPath),
					lineNumber: lineNumber,
					importPath: importPath,
				})
			}
		}
		if scanErr := scanner.Err(); scanErr != nil {
			return fmt.Errorf("ファイルの読み込みに失敗しました: %w", scanErr)
		}
		return nil
	}

	if err := filepath.Walk(backendDir, walkFn); err != nil {
		return nil, err
	}
	return violations, nil
}
