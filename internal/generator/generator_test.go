// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DingTalk-Real-AI/dingtalk-workspace-cli/internal/ir"
)

func TestWriteArtifactsMaterializesFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := WriteArtifacts(root, []Artifact{
		{Path: "docs/generated/example.txt", Content: []byte("example\n")},
	}); err != nil {
		t.Fatalf("WriteArtifacts() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "docs/generated/example.txt")); err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
}

func TestResolveRequiredServiceSkillsSeparatesMissing(t *testing.T) {
	t.Parallel()

	available := map[string]struct{}{
		"doc":   {},
		"drive": {},
	}
	required, missing := resolveRequiredServiceSkills([]string{"doc", "calendar", "drive"}, available)

	if len(required) != 2 || required[0] != "dws-doc" || required[1] != "dws-drive" {
		t.Fatalf("required = %#v, want [dws-doc dws-drive]", required)
	}
	if len(missing) != 1 || missing[0] != "dws-calendar" {
		t.Fatalf("missing = %#v, want [dws-calendar]", missing)
	}
}

func TestFinalizeMarkdownNormalizesTrailingWhitespace(t *testing.T) {
	t.Parallel()

	input := "line one  \r\nline two\t \n\n"
	got := finalizeMarkdown(input)
	want := "line one\nline two\n"
	if got != want {
		t.Fatalf("finalizeMarkdown() = %q, want %q", got, want)
	}
}

func TestGenerateUsesFlattenedSkillsPaths(t *testing.T) {
	t.Parallel()

	artifacts, err := Generate(ir.Catalog{})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	paths := map[string]string{}
	for _, artifact := range artifacts {
		paths[artifact.Path] = string(artifact.Content)
	}

	for _, want := range []string{
		"skills/generated/apis.md",
		"skills/generated/canonical-surface/api.md",
		"skills/generated/dws-shared/api.md",
	} {
		if _, ok := paths[want]; !ok {
			t.Fatalf("Generate() missing artifact %q", want)
		}
	}

	for _, oldPath := range []string{
		"skills/dws/generated/apis.md",
		"skills/dws/generated/canonical-surface/api.md",
		"skills/dws/generated/dws-shared/api.md",
	} {
		if _, ok := paths[oldPath]; ok {
			t.Fatalf("Generate() still emitted legacy artifact %q", oldPath)
		}
	}

	readme, ok := paths["docs/generated/README.md"]
	if !ok {
		t.Fatal("Generate() missing docs/generated/README.md")
	}
	if !strings.Contains(readme, "`skills/generated/apis.md`") {
		t.Fatalf("docs/generated/README.md missing flattened skills path:\n%s", readme)
	}
	if strings.Contains(readme, "`skills/dws/generated/apis.md`") {
		t.Fatalf("docs/generated/README.md still references legacy skills path:\n%s", readme)
	}
}
