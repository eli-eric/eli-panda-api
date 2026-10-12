package db

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4/database/multistmt"
	migrateneo4j "github.com/golang-migrate/migrate/v4/database/neo4j"
)

// MigrateNeo4jMainInstance runs with x-multi-statement=true, so golang-migrate
// splits every file on ";" without understanding Cypher comments and sends
// each non-blank chunk to Neo4j. A chunk holding nothing but "//" comments —
// what a ";" inside a comment produces — fails with "Unexpected end of input"
// and leaves SchemaMigration dirty, which blocks every later API boot. This
// runs each migration through the driver's own parser and rejects such chunks.
func TestMigrationChunksAreNotCommentOnly(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("neo4j", "migrations", "*.cypher"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no migration files found")
	}

	delimiter := migrateneo4j.StatementSeparator
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if line := embeddedSeparatorLine(content); line != 0 {
			t.Errorf("%s:%d: semicolon inside a comment or literal is split into a separate Cypher statement", filepath.Base(file), line)
		}
		err = multistmt.Parse(bytes.NewReader(content), delimiter, migrateneo4j.DefaultMultiStatementMaxSize, func(chunk []byte) bool {
			// Mirror the driver: blank chunks and a bare delimiter are skipped.
			statement := bytes.TrimSuffix(bytes.TrimSpace(chunk), delimiter)
			if len(statement) > 0 && isCommentOnly(string(statement)) {
				t.Errorf("%s: statement %q contains only comments; Neo4j rejects it (is there a \";\" inside a comment?)",
					filepath.Base(file), string(statement))
			}
			return true
		})
		if err != nil {
			t.Errorf("%s: %v", filepath.Base(file), err)
		}
	}
}

// The migration driver splits on every semicolon, including ones inside
// comments, strings and escaped identifiers. Such a split can dirty the DB.
func embeddedSeparatorLine(content []byte) int {
	const (
		plain = iota
		lineComment
		blockComment
		quoted
	)
	state := plain
	var quote byte
	line := 1
	for i := 0; i < len(content); i++ {
		c := content[i]
		if c == '\n' {
			line++
			if state == lineComment {
				state = plain
			}
			continue
		}
		if state != plain && c == ';' {
			return line
		}
		switch state {
		case plain:
			if c == '/' && i+1 < len(content) && content[i+1] == '/' {
				state = lineComment
				i++
			} else if c == '/' && i+1 < len(content) && content[i+1] == '*' {
				state = blockComment
				i++
			} else if c == '\'' || c == '"' || c == '`' {
				state, quote = quoted, c
			}
		case blockComment:
			if c == '*' && i+1 < len(content) && content[i+1] == '/' {
				state = plain
				i++
			}
		case quoted:
			if c == '\\' && i+1 < len(content) {
				i++
			} else if c == quote {
				state = plain
			}
		}
	}
	return 0
}

func TestEmbeddedSeparatorLine(t *testing.T) {
	for _, example := range []struct {
		name string
		text string
		line int
	}{
		{"normal", "RETURN 1;\nRETURN 2", 0},
		{"comment", "// note;\nRETURN 1", 1},
		{"string", "RETURN 'a;b'", 1},
		{"escaped identifier", "MATCH (n:`a;b`) RETURN n", 1},
		{"block comment", "RETURN 1 /* ; */", 1},
		{"quoted slash", "RETURN '//';", 0},
	} {
		t.Run(example.name, func(t *testing.T) {
			if got := embeddedSeparatorLine([]byte(example.text)); got != example.line {
				t.Errorf("got line %d, want %d", got, example.line)
			}
		})
	}
}

func isCommentOnly(statement string) bool {
	for _, line := range strings.Split(statement, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") {
			return false
		}
	}
	return true
}
