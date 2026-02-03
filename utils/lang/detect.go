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

	switch {
	case strings.Contains(s, "<?php"):
		return "PHP"
	case strings.Contains(s, "<!doctypehtml"):
		return "HTML"
	case strings.Contains(s, "<?xml"):
		return "XML"
	case strings.Contains(s, "@echooff") || strings.HasPrefix(s, "rem"):
		return "Batch (CMD)"
	case strings.Contains(s, "usestrict;") || strings.Contains(s, "my$"):
		return "Perl"
	case strings.HasPrefix(s, "from") && (strings.Contains(s, "run") || strings.Contains(s, "cmd")):
		return "Dockerfile"
	case strings.Contains(s, "resource\"") && strings.Contains(s, "provider\""):
		return "Terraform (HCL)"
	case strings.HasPrefix(s, "---") && strings.Contains(s, ":"):
		return "YAML"
	case strings.Contains(s, "<-") && (strings.Contains(s, "function(") || strings.Contains(s, "c(")):
		return "R"
	case strings.Contains(s, "defmodule") && strings.Contains(s, "do"):
		return "Elixir"
	case strings.Contains(s, "println!"):
		return "Rust"
	case strings.Contains(s, "write-host") || strings.Contains(s, "-eq"):
		return "PowerShell"
	case strings.Contains(s, "@interface") || strings.Contains(s, "#import<"):
		return "Objective-C"
	case strings.Contains(s, "main=do") || (strings.Contains(s, "::") && !strings.Contains(s, "std::")):
		return "Haskell"
	case strings.Contains(s, "(defn") || strings.Contains(s, "(ns"):
		return "Clojure"
	case strings.Contains(s, ":symbol"):
		return "Ruby"
	case strings.Contains(s, ":string") || strings.Contains(s, ":number") || strings.Contains(s, "interface"):
		if strings.Contains(s, "/>") {
			return "TSX"
		}
		return "TypeScript"
	case strings.Contains(s, "std::") || strings.Contains(s, "#include<iostream>") || strings.Contains(s, "cin>>") || strings.Contains(s, "cout<<"):
		return "C++"
	case strings.Contains(s, "package") || strings.Contains(s, "import\"") || strings.Contains(s, ":="):
		return "Go"
	case strings.Contains(s, "importjava.") || strings.Contains(s, "publicclass") || strings.Contains(s, "system.out.println"):
		return "Java"
	case (strings.Contains(s, "namespace") && strings.Contains(s, "class")) || strings.Contains(s, "usingsystem"):
		return "C#"
	case strings.Contains(s, "import'dart:"):
		return "Dart"
	case strings.Contains(s, "importscala.") || (strings.Contains(s, "def") && strings.Contains(s, ":") && strings.Contains(s, "=")):
		return "Scala"
	case (strings.Contains(s, "def") && strings.Contains(s, ":")) || (strings.Contains(s, "import") && !strings.Contains(s, "#import")):
		return "Python"
	case strings.Contains(s, "def") && strings.Contains(s, "end"):
		return "Ruby"
	case strings.Contains(s, "func") && (strings.Contains(s, "->") || strings.Contains(s, "var") || strings.Contains(s, "let")):
		return "Swift"
	case strings.Contains(s, "fun") && (strings.Contains(s, "val") || strings.Contains(s, "?:")):
		return "Kotlin"
	case strings.Contains(s, "console.log(") || strings.Contains(s, "document.getelementbyid") || strings.Contains(s, "=>"):
		if strings.Contains(s, "/>") {
			return "JSX"
		}
		return "JavaScript"
	case strings.Contains(s, "#include<") && strings.Contains(s, ".h>"):
		return "C"
	case strings.Contains(s, "localfunction"):
		return "Lua"
	case (strings.Contains(s, "select") && strings.Contains(s, "from")) || strings.Contains(s, "createtable"):
		return "SQL"
	case strings.Contains(s, "background-color:") || strings.Contains(s, "font-size:"):
		return "CSS"
	case strings.Contains(s, "\":\"") || strings.Contains(s, "\":{"):
		return "JSON"
	case strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") && !strings.Contains(s, "{"):
		return "TOML"
	case strings.Contains(s, "<html>"):
		return "HTML"
	case strings.Contains(s, "function") && strings.Contains(s, "end"):
		return "MATLAB"
	case strings.Contains(s, "echo") && strings.Contains(s, "$"):
		return "Shell"
	default:
		return "Unknown"
	}
}
