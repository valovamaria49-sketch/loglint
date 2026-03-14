package analysis

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestCheckStartsWithLower(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{name: "lowercase", msg: "starting server", want: ""},
		{name: "uppercase", msg: "Starting server", want: "log message should start with a lowercase English letter"},
		{name: "space then lowercase", msg: "  starting server", want: ""},
		{name: "space then uppercase", msg: "  Starting server", want: "log message should start with a lowercase English letter"},
		{name: "no letters", msg: "  1234  ", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkStartsWithLower(tt.msg)
			if got != tt.want {
				t.Fatalf("checkStartsWithLower(%q) = %q, want %q", tt.msg, got, tt.want)
			}
		})
	}
}

func TestCheckAsciiLettersAndDigits(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want string
	}{
		{name: "only letters", msg: "starting server", want: ""},
		{name: "letters and digits", msg: "server1 started2", want: ""},
		{name: "contains punctuation", msg: "server started!", want: "log message contains disallowed special characters"},
		{name: "contains emoji", msg: "server started🚀", want: "log message contains non-ASCII characters (emoji or other alphabets)"},
		{name: "non ascii letters", msg: "запуск сервера", want: "log message contains non-ASCII characters (emoji or other alphabets)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkContainsOnlyLettersDigits(tt.msg)
			if got != tt.want {
				t.Fatalf("checkAsciiLettersAndDigits(%q) = %q, want %q", tt.msg, got, tt.want)
			}
		})
	}
}

func resetConfigForTest() {
	configOnce = sync.Once{}
	activeConfig = config{}
}

func TestCheckNoSensitiveData_DefaultKeywords(t *testing.T) {
	resetConfigForTest()

	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{name: "no sensitive", msg: "user authenticated successfully", want: false},
		{name: "password", msg: "user password: 123", want: true},
		{name: "token", msg: "token: abc", want: true},
		{name: "api_key", msg: "api_key=secret", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkNoSensitiveData(tt.msg)
			if tt.want && got == "" {
				t.Fatalf("expected violation for %q, got none", tt.msg)
			}
			if !tt.want && got != "" {
				t.Fatalf("expected no violation for %q, got %q", tt.msg, got)
			}
		})
	}
}

func TestGetSensitiveKeywords_FromConfigFile(t *testing.T) {
	resetConfigForTest()

	dir, err := os.MkdirTemp("", "loglint-config-test")
	if err != nil {
		t.Fatalf("MkDirTemp: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("failed to remove temp dir: %v", err)
		}
	}()

	prevDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	defer func() {
		if err := os.Chdir(prevDir); err != nil {
			t.Fatalf("failed to chdir back: %v", err)
		}
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	configContent := `{"sensitive_keywords":["card_number","iban"]}`
	if err := os.WriteFile(filepath.Join(dir, "loglint_config.json"), []byte(configContent), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	keywords := getSensitiveKeywords()
	if len(keywords) != 2 {
		t.Fatalf("expected 2 keywords from config, got %d", len(keywords))
	}
	if keywords[0] != "card_number" || keywords[1] != "iban" {
		t.Fatalf("unexpected keywords: %#v", keywords)
	}
}
