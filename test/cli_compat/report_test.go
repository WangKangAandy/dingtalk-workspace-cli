package cli_compat_test

import (
	"testing"
)

// ── report template list ───────────────────────────────────

func TestReportTemplateList_should_call_list_report_templates(t *testing.T) {
	cap := setupTestDeps(t, "report")
	root := buildRoot()
	err := execCmd(t, root, []string{"report", "template", "list"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertToolName(t, cap, "list_report_templates")
}

// ── report template detail ──────────────────────────────────

func TestReportTemplateDetail_should_call_get_report_template_detail(t *testing.T) {
	cap := setupTestDeps(t, "report")
	root := buildRoot()
	err := execCmd(t, root, []string{"report", "template", "detail"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertToolName(t, cap, "get_report_template_detail")
}

// ── report detail ───────────────────────────────────────────

func TestReportDetail_should_call_get_report_detail(t *testing.T) {
	cap := setupTestDeps(t, "report")
	root := buildRoot()
	err := execCmd(t, root, []string{"report", "detail"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertToolName(t, cap, "get_report_detail")
}

// ── report list ─────────────────────────────────────────────

// ── report sent ──────────────────────────────────────────────
