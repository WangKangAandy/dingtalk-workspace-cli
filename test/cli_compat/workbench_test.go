package cli_compat_test

// 暂未上线，保留测试源码但不执行
// 对应 register.go 中注释掉的 workbenchCmd

/*
import "testing"

// ── workbench app list ─────────────────────────────────────

func TestWorkbenchAppList_should_call_correct_tool(t *testing.T) {
	cap := setupTestDeps(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	err := execCmd(t, root, []string{"workbench", "app", "list"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertToolName(t, cap, "获取用户所有工作台应用")
}

func TestWorkbenchAppList_should_pass_fromCLI_input(t *testing.T) {
	cap := setupTestDeps(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	_ = execCmd(t, root, []string{"workbench", "app", "list"}, nil)
	assertToolArg(t, cap, "input", "fromCLI")
}

func TestWorkbenchAppList_should_make_exactly_one_call(t *testing.T) {
	cap := setupTestDeps(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	_ = execCmd(t, root, []string{"workbench", "app", "list"}, nil)
	assertCallCount(t, cap, 1)
}

func TestWorkbenchAppList_should_not_call_when_dry_run(t *testing.T) {
	cap := setupTestDepsWithDryRun(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	err := execCmd(t, root, []string{"workbench", "app", "list"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertCallCount(t, cap, 0)
}

func TestWorkbenchAppList_should_succeed_with_no_flags(t *testing.T) {
	_ = setupTestDeps(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	err := execCmd(t, root, []string{"workbench", "app", "list"}, nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// ── workbench app get ──────────────────────────────────────

func TestWorkbenchAppGet_should_call_correct_tool(t *testing.T) {
	cap := setupTestDeps(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	err := execCmd(t, root, []string{"workbench", "app", "get"}, map[string]string{"ids": "app1,app2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertToolName(t, cap, "批量拉取应用详情")
}

func TestWorkbenchAppGet_should_pass_ids_as_string(t *testing.T) {
	cap := setupTestDeps(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	_ = execCmd(t, root, []string{"workbench", "app", "get"}, map[string]string{"ids": "app1,app2,app3"})
	assertToolArg(t, cap, "appIds", "app1,app2,app3")
}

func TestWorkbenchAppGet_should_handle_single_id(t *testing.T) {
	cap := setupTestDeps(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	_ = execCmd(t, root, []string{"workbench", "app", "get"}, map[string]string{"ids": "single_app"})
	assertToolArg(t, cap, "appIds", "single_app")
}

func TestWorkbenchAppGet_should_handle_empty_ids(t *testing.T) {
	cap := setupTestDeps(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	_ = execCmd(t, root, []string{"workbench", "app", "get"}, map[string]string{"ids": ""})
	assertToolArg(t, cap, "appIds", "")
}

func TestWorkbenchAppGet_should_not_call_when_dry_run(t *testing.T) {
	cap := setupTestDepsWithDryRun(t, "workbench")
	root := buildRoot()
	skipIfNotRegistered(t, root, "workbench")
	_ = execCmd(t, root, []string{"workbench", "app", "get"}, map[string]string{"ids": "app1"})
	assertCallCount(t, cap, 0)
}
*/
