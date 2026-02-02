package lang

import (
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
)

func AnalyzeLanguage(text string) string {
	lexer := lexers.Analyse(text)
	if lexer == nil {
		return AnalyzeLanguageInside(text)
	} else {
		return lexer.Config().Name

	}
}

func AnalyzeLanguageInside(code string) string {
	s := strings.ToLower(code)
	replacer := strings.NewReplacer(
		" ", "", "\n", "", "\t", "", "\r", "",
	)
	s = replacer.Replace(s)

	if strings.Contains(s, "<?php") {
		return "PHP"
	}
	if strings.Contains(s, "<!doctypehtml") {
		return "HTML"
	}
	if strings.Contains(s, "<?xml") {
		return "XML"
	}
	if strings.Contains(s, "@echooff") || strings.HasPrefix(s, "rem") {
		return "Batch (CMD)"
	}
	if strings.Contains(s, "usestrict;") || strings.Contains(s, "my$") {
		return "Perl"
	}
	if strings.HasPrefix(s, "from") && (strings.Contains(s, "run") || strings.Contains(s, "cmd")) {
		return "Dockerfile"
	}
	if strings.Contains(s, "resource\"") && strings.Contains(s, "provider\"") {
		return "Terraform (HCL)"
	}
	if strings.HasPrefix(s, "---") && strings.Contains(s, ":") {
		return "YAML"
	}
	if strings.Contains(s, "<-") && (strings.Contains(s, "function(") || strings.Contains(s, "c(")) {
		return "R"
	}
	if strings.Contains(s, "defmodule") && strings.Contains(s, "do") {
		return "Elixir"
	}
	if strings.Contains(s, "println!") {
		return "Rust"
	}
	if strings.Contains(s, "write-host") || strings.Contains(s, "-eq") {
		return "PowerShell"
	}
	if strings.Contains(s, "@interface") || strings.Contains(s, "#import<") {
		return "Objective-C"
	}
	if strings.Contains(s, "main=do") || strings.Contains(s, "::") && !strings.Contains(s, "std::") {
		return "Haskell"
	}
	if strings.Contains(s, "(defn") || strings.Contains(s, "(ns") {
		return "Clojure"
	}
	if strings.Contains(s, ":symbol") {
		return "Ruby"
	}
	if strings.Contains(s, ":string") || strings.Contains(s, ":number") || strings.Contains(s, "interface") {
		if strings.Contains(s, "/>") {
			return "TSX"
		}
		return "TypeScript"
	}
	if strings.Contains(s, "std::") || strings.Contains(s, "#include<iostream>") || strings.Contains(s, "cin>>") || strings.Contains(s, "cout<<") {
		return "C++"
	}
	if strings.Contains(s, "package") || strings.Contains(s, "import\"") || strings.Contains(s, ":=") {
		return "Go"
	}
	if strings.Contains(s, "importjava.") || strings.Contains(s, "publicclass") || strings.Contains(s, "system.out.println") {
		return "Java"
	}
	if strings.Contains(s, "namespace") && strings.Contains(s, "class") || strings.Contains(s, "usingsystem") {
		return "C#"
	}
	if strings.Contains(s, "import'dart:") {
		return "Dart"
	}
	if strings.Contains(s, "importscala.") || (strings.Contains(s, "def") && strings.Contains(s, ":") && strings.Contains(s, "=")) {
		return "Scala"
	}
	if (strings.Contains(s, "def") && strings.Contains(s, ":")) || (strings.Contains(s, "import") && !strings.Contains(s, "#import")) {
		return "Python"
	}
	if strings.Contains(s, "def") && strings.Contains(s, "end") {
		return "Ruby"
	}
	if strings.Contains(s, "func") && (strings.Contains(s, "->") || strings.Contains(s, "var") || strings.Contains(s, "let")) {
		return "Swift"
	}
	if strings.Contains(s, "fun") && (strings.Contains(s, "val") || strings.Contains(s, "?:")) {
		return "Kotlin"
	}
	if strings.Contains(s, "console.log(") || strings.Contains(s, "document.getelementbyid") || strings.Contains(s, "=>") {
		if strings.Contains(s, "/>") {
			return "JSX"
		}
		return "JavaScript"
	}
	if strings.Contains(s, "#include<") && strings.Contains(s, ".h>") {
		return "C"
	}
	if strings.Contains(s, "localfunction") {
		return "Lua"
	}
	if (strings.Contains(s, "select") && strings.Contains(s, "from")) || strings.Contains(s, "createtable") {
		return "SQL"
	}
	if strings.Contains(s, "background-color:") || strings.Contains(s, "font-size:") {
		return "CSS"
	}
	if strings.Contains(s, "\":\"") || strings.Contains(s, "\":{") {
		return "JSON"
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") && !strings.Contains(s, "{") {
		return "TOML"
	}
	if strings.Contains(s, "<html>") {
		return "HTML"
	}
	if strings.Contains(s, "function") && strings.Contains(s, "end") {
		return "MATLAB"
	}
	if strings.Contains(s, "echo") && strings.Contains(s, "$") {
		return "Shell"
	}
	return "Unknown"
}
