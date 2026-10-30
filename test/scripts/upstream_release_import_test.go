package scripts_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPinnedUpstreamReleaseImport(t *testing.T) {
	for _, scenario := range []string{"exact import", "changelog edit", "archive edit", "new fragment", "mode change", "wrong base", "missing ancestry"} {
		t.Run(scenario, func(t *testing.T) {
			repo := newChangelogGateRepo(t)
			writeMetadata := func() {
				changelogGateWrite(t, repo.root, "CHANGELOG.md", changelogGateValidRelease, 0o644)
				changelogGateWrite(t, repo.root, ".changes/released/1.0.1-beta.1/imported.md", "### Added\n\n- Imported release.\n", 0o644)
			}
			writeMetadata()
			repo.commit(t, "upstream release")
			upstream := strings.TrimSpace(changelogGateGit(t, repo.root, "rev-parse", "HEAD"))
			if scenario == "missing ancestry" {
				changelogGateGit(t, repo.root, "switch", "-c", "without-upstream", repo.base)
				writeMetadata()
			}
			changelogGateWrite(t, repo.root, "docs/import.md", "Import the upstream baseline.\n", 0o644)
			switch scenario {
			case "changelog edit":
				changelogGateWrite(t, repo.root, "CHANGELOG.md", changelogGateValidRelease+"Extra note.\n", 0o644)
			case "archive edit":
				changelogGateWrite(t, repo.root, ".changes/released/1.0.1-beta.1/imported.md", "Tampered archive.\n", 0o644)
			case "new fragment":
				changelogGateWrite(t, repo.root, ".changes/extra.md", "### Added\n\n- Extra behavior.\n", 0o644)
			}
			repo.commit(t, "fork import")
			if scenario == "mode change" {
				changelogGateGit(t, repo.root, "update-index", "--chmod=+x", ".changes/released/1.0.1-beta.1/imported.md")
				changelogGateGit(t, repo.root, "commit", "-m", "change archive mode")
			}

			// Substitute fixture commits only in the copied helper. Production pins
			// remain literal reviewed SHAs, with no environment override.
			helper := filepath.Join(repo.root, "scripts/policy/upstream-release-import.sh")
			data, err := os.ReadFile(helper)
			if err != nil {
				t.Fatal(err)
			}
			body := strings.ReplaceAll(string(data), "dd34e2241c479bde9d70d29c4847594b2544a60a", repo.base)
			body = strings.ReplaceAll(body, "f215043e676d6889e1ddce3464828d8b1ec40c31", upstream)
			if err := os.WriteFile(helper, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
			base := repo.base
			if scenario == "wrong base" {
				base = upstream
			}
			for _, gate := range []string{"changelog", "fragments"} {
				var output string
				if gate == "changelog" {
					output, err = repo.runRefs(t, "--content-only", base, "HEAD")
				} else {
					output, err = repo.runFragmentPolicy(t, base, "HEAD")
				}
				acceptedImport := strings.Contains(output, "verified pinned upstream")
				if scenario == "exact import" {
					if err != nil || !acceptedImport {
						t.Fatalf("%s rejected exact import: %v\n%s", gate, err, output)
					}
				} else if acceptedImport || (scenario != "wrong base" && err == nil) {
					t.Fatalf("%s accepted invalid import: %v\n%s", gate, err, output)
				}
			}
			if scenario == "exact import" {
				if output, err := repo.runRefs(t, "--fast-path", base, "HEAD"); err == nil {
					t.Fatalf("import must not qualify for changelog-only fast path: %s", output)
				}
			}
		})
	}
}
