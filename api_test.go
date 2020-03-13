package reply_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"reply"
)

func TestText(t *testing.T) {
	// Add files to be skipped.
	var skipped []string

	err := filepath.Walk("testdata/emails", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".txt" {
			t.Run(path, func(t *testing.T) {
				for _, filename := range skipped {
					if filename == filepath.Base(path) {
						t.Skipf("%s is not implemented", filename)
					}
				}

				in, err := os.Open(path)
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				expected, err := ioutil.ReadFile(fmt.Sprintf("testdata/reply/%s", filepath.Base(path)))
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				replyText, err := reply.FromReader(in)
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				if strings.TrimSpace(replyText) != strings.TrimSpace(string(expected)) {
					t.Errorf("\nexpected:\n%s\n\ngot:\n%s", string(expected), replyText)
				}
			})
		}

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
