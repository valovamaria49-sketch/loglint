// Package analysis implements a linter for Go log messages (log/slog, go.uber.org/zap).
// It checks that log messages start with a lowercase letter, use only English and ASCII,
// avoid special characters/emojis, and do not contain sensitive data (configurable).
package analysis

import (
	"encoding/json"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the loglint analyzer for golangci-lint and go/analysis.
// It validates log messages against configurable rules. Configuration
// is read from loglint_config.json in the working directory.
var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages for style and content rules",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			_, _, msgIndex, ok := detectLogCall(pass, call)
			if !ok {
				return true
			}

			if msgIndex >= len(call.Args) {
				return true
			}

			basicLit, ok := call.Args[msgIndex].(*ast.BasicLit)
			if !ok || basicLit.Kind != token.STRING {
				return true
			}

			msg, err := strconv.Unquote(basicLit.Value)
			if err != nil {
				return true
			}

			checkMessage(pass, basicLit, msg)

			return true
		})
	}

	return nil, nil
}

func detectLogCall(pass *analysis.Pass, call *ast.CallExpr) (logger, level string, msgIndex int, ok bool) {
	msgIndex = 0

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", "", 0, false
	}

	level = sel.Sel.Name

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return "", "", 0, false
	}

	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return "", "", 0, false
	}

	if pkgName, ok := obj.(*types.PkgName); ok {
		switch pkgName.Imported().Path() {
		case "log/slog":
			return "slog", level, msgIndex, true
		case "go.uber.org/zap":
			return "zap", level, msgIndex, true
		default:
			return "", "", 0, false
		}
	}

	if tv := pass.TypesInfo.TypeOf(ident); tv != nil {
		typeString := tv.String()
		switch {
		case strings.Contains(typeString, "go.uber.org/zap"):
			return "zap", level, msgIndex, true
		case strings.Contains(typeString, "log/slog"):
			return "slog", level, msgIndex, true
		}
	}

	return "", "", 0, false
}

type config struct {
	SensitiveKeywords   []string `json:"sensitive_keywords"`
	EnableLowercaseRule bool     `json:"enable_lowercase_rule"`
	EnableAsciiRule     bool     `json:"enable_ascii_rule"`
	EnableEnglishRule   bool     `json:"enable_english_rule"`
	EnableSensitiveRule bool     `json:"enable_sensitive_rule"`
}

var (
	defaultSensitiveKeywords = []string{
		"password",
		"passwd",
		"pwd",
		"token",
		"api_key",
		"apikey",
		"secret",
		"jwt",
	}

	configOnce   sync.Once
	activeConfig config
)

func loadConfig() config {
	configOnce.Do(func() {
		activeConfig = config{
			SensitiveKeywords:   defaultSensitiveKeywords,
			EnableLowercaseRule: true,
			EnableAsciiRule:     true,
			EnableEnglishRule:   true,
			EnableSensitiveRule: true,
		}

		data, err := os.ReadFile("loglint_config.json")
		if err != nil {
			return
		}

		var cfg config
		if err := json.Unmarshal(data, &cfg); err != nil {
			return
		}

		if len(cfg.SensitiveKeywords) > 0 {
			activeConfig.SensitiveKeywords = cfg.SensitiveKeywords
		}

		activeConfig.EnableLowercaseRule = cfg.EnableLowercaseRule
		activeConfig.EnableAsciiRule = cfg.EnableAsciiRule
		activeConfig.EnableEnglishRule = cfg.EnableEnglishRule
		activeConfig.EnableSensitiveRule = cfg.EnableSensitiveRule
	})

	return activeConfig
}

func checkMessage(pass *analysis.Pass, lit *ast.BasicLit, msg string) {
	cfg := loadConfig()
	pos := lit.Pos()

	if cfg.EnableLowercaseRule {
		if reason := checkStartsWithLower(msg); reason != "" {
			d := analysis.Diagnostic{
				Pos:     pos,
				Message: reason,
			}
			if fix := suggestLowercaseFix(lit, msg); fix != nil {
				d.SuggestedFixes = []analysis.SuggestedFix{*fix}
			}
			pass.Report(d)
		}
	}

	if cfg.EnableAsciiRule {
		if reason := checkContainsOnlyLettersDigits(msg); reason != "" {
			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				Message: reason,
			})
		}
	}

	if cfg.EnableEnglishRule {
		if reason := checkEnglishOnly(msg); reason != "" {
			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				Message: reason,
			})
		}
	}

	if cfg.EnableSensitiveRule {
		if reason := checkNoSensitiveData(msg); reason != "" {
			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				Message: reason,
			})
		}
	}
}

func suggestLowercaseFix(lit *ast.BasicLit, msg string) *analysis.SuggestedFix {
	trimmed := strings.TrimSpace(msg)
	if len(trimmed) == 0 {
		return nil
	}
	first := rune(trimmed[0])
	if first < 'A' || first > 'Z' {
		return nil
	}
	lower := string(unicode.ToLower(first))
	offset := strings.Index(msg, trimmed)
	if offset < 0 {
		return nil
	}
	if lit.Kind != token.STRING || len(lit.Value) < 2 {
		return nil
	}
	start := int(lit.Pos()) + 1 + offset
	firstLen := len(string(first))
	return &analysis.SuggestedFix{
		Message: "use lowercase first letter",
		TextEdits: []analysis.TextEdit{{
			Pos:     token.Pos(start),
			End:     token.Pos(start + firstLen),
			NewText: []byte(lower),
		}},
	}
}

func checkContainsOnlyLettersDigits(msg string) string {
	for _, r := range msg {
		if r > unicode.MaxASCII {
			return "log message contains non-ASCII characters (emoji or other alphabets)"
		}

		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == ' ' {
			continue
		}

		return "log message contains disallowed special characters"
	}

	return ""
}

func checkStartsWithLower(msg string) string {
	cfg := loadConfig()

	if !cfg.EnableLowercaseRule {
		return ""
	}

	trimmed := strings.TrimSpace(msg)
	for _, r := range trimmed {
		if r >= 'A' && r <= 'Z' {
			return "log message should start with a lowercase English letter"
		}
		if r >= 'a' && r <= 'z' {
			return ""
		}
	}

	return ""
}

func checkEnglishOnly(msg string) string {
	cfg := loadConfig()

	if !cfg.EnableEnglishRule {
		return ""
	}

	for _, r := range msg {
		if !unicode.IsLetter(r) {
			continue
		}
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return "log message should contain only English letters"
		}
	}

	return ""
}

func getSensitiveKeywords() []string {
	cfg := loadConfig()
	return cfg.SensitiveKeywords
}

func checkNoSensitiveData(msg string) string {
	cfg := loadConfig()

	if !cfg.EnableSensitiveRule {
		return ""
	}

	lower := strings.ToLower(msg)

	for _, kw := range getSensitiveKeywords() {
		if strings.Contains(lower, kw) {
			return "log message contains potentially sensitive data: " + kw
		}
	}

	return ""
}
