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
	"testing"
)

func TestLoadPersonaRegistry(t *testing.T) {
	t.Parallel()

	registry, err := loadPersonaRegistry()
	if err != nil {
		t.Fatalf("loadPersonaRegistry() error = %v", err)
	}
	if len(registry.Personas) == 0 {
		t.Fatal("loadPersonaRegistry() returned empty personas")
	}
}

func TestLoadRecipeRegistry(t *testing.T) {
	t.Parallel()

	registry, err := loadRecipeRegistry()
	if err != nil {
		t.Fatalf("loadRecipeRegistry() error = %v", err)
	}
	if len(registry.Recipes) == 0 {
		t.Fatal("loadRecipeRegistry() returned empty recipes")
	}
}

func TestValidateRegistryReferences(t *testing.T) {
	t.Parallel()

	personas, err := loadPersonaRegistry()
	if err != nil {
		t.Fatalf("loadPersonaRegistry() error = %v", err)
	}
	recipes, err := loadRecipeRegistry()
	if err != nil {
		t.Fatalf("loadRecipeRegistry() error = %v", err)
	}
	if err := validateRegistryReferences(personas, recipes); err != nil {
		t.Fatalf("validateRegistryReferences() error = %v", err)
	}
}

func TestValidateRegistryReferencesRejectsUnknownWorkflow(t *testing.T) {
	t.Parallel()

	personas := PersonaRegistry{
		Personas: []PersonaEntry{
			{
				Name:         "pm",
				Title:        "PM",
				Description:  "project manager",
				Services:     []string{"doc"},
				Workflows:    []string{"unknown-workflow"},
				Instructions: []string{"do work"},
			},
		},
	}
	recipes := RecipeRegistry{
		Recipes: []RecipeEntry{
			{
				Name:        "meeting-prep",
				Title:       "Meeting Prep",
				Description: "prepare",
				Category:    "collaboration",
				Services:    []string{"doc"},
				Steps:       []string{"step"},
			},
		},
	}

	if err := validateRegistryReferences(personas, recipes); err == nil {
		t.Fatal("validateRegistryReferences() should fail for unknown workflow")
	}
}

func TestUniqueSkillProducts(t *testing.T) {
	t.Parallel()

	got := uniqueSkillProducts([]string{"doc", "drive", "doc", " DRIVE ", ""})
	want := []string{"doc", "drive"}
	if len(got) != len(want) {
		t.Fatalf("uniqueSkillProducts() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("uniqueSkillProducts()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestServiceRefsSupportServicesAndProducts(t *testing.T) {
	t.Parallel()

	persona := PersonaEntry{
		Services: []string{"doc", "drive"},
		Products: []string{"doc", "chat"},
	}
	got := persona.serviceRefs()
	want := []string{"chat", "doc", "drive"}
	if len(got) != len(want) {
		t.Fatalf("serviceRefs() len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("serviceRefs()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestNormalizeRegistryToken(t *testing.T) {
	t.Parallel()

	if got := normalizeRegistryToken(" Meeting_Prep "); got != "meeting-prep" {
		t.Fatalf("normalizeRegistryToken() = %q, want %q", got, "meeting-prep")
	}
}

func TestValidateRegistryReferencesUsesNormalizedWorkflowName(t *testing.T) {
	t.Parallel()

	personas := PersonaRegistry{
		Personas: []PersonaEntry{
			{
				Name:         "pm",
				Title:        "PM",
				Description:  "project manager",
				Services:     []string{"doc"},
				Workflows:    []string{"meeting-prep"},
				Instructions: []string{"do work"},
			},
		},
	}
	recipes := RecipeRegistry{
		Recipes: []RecipeEntry{
			{
				Name:        "meeting_prep",
				Title:       "Meeting Prep",
				Description: "prepare",
				Category:    "collaboration",
				Services:    []string{"doc"},
				Steps:       []string{"step"},
			},
		},
	}

	if err := validateRegistryReferences(personas, recipes); err != nil {
		t.Fatalf("validateRegistryReferences() error = %v", err)
	}
}

func TestReadRegistryYAMLFromEnvOverride(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "personas.yaml")
	want := []byte("personas: []\n")
	if err := os.WriteFile(path, want, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	t.Setenv(PersonaRegistryPathEnv, path)

	got, err := readRegistryYAML(PersonaRegistryPathEnv, []byte("fallback"))
	if err != nil {
		t.Fatalf("readRegistryYAML() error = %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("readRegistryYAML() = %q, want %q", string(got), string(want))
	}
}
