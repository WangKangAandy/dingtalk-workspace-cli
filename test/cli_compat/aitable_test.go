package cli_compat_test

import (
	"testing"
)

// ── base list ───────────────────────────────────────────────

func TestAitableBaseList_should_call_list_bases(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	err := execCmd(t, root, []string{"aitable", "base", "list"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertToolName(t, cap, "list_bases")
	assertCallCount(t, cap, 1)
}

// ── base search ─────────────────────────────────────────────

func TestAitableBaseSearch_should_call_search_bases(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "search"}, map[string]string{
		"query": "项目管理",
	})
	assertToolName(t, cap, "search_bases")
	assertToolArg(t, cap, "query", "项目管理")
}

// ── base get ────────────────────────────────────────────────

func TestAitableBaseGet_should_call_get_base(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "get"}, map[string]string{
		"base-id": "BASE_001",
	})
	assertToolName(t, cap, "get_base")
	assertToolArg(t, cap, "baseId", "BASE_001")
}

// ── base create ─────────────────────────────────────────────

func TestAitableBaseCreate_should_call_create_base(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "create"}, map[string]string{
		"name": "新表格",
	})
	assertToolName(t, cap, "create_base")
	assertToolArg(t, cap, "baseName", "新表格")
}

func TestAitableBaseCreate_should_pass_template_id(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "create"}, map[string]string{
		"name": "模板表格", "template-id": "TPL_001",
	})
	assertToolArg(t, cap, "templateId", "TPL_001")
}

// ── base update ─────────────────────────────────────────────

func TestAitableBaseUpdate_should_call_update_base(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "update"}, map[string]string{
		"base-id": "BASE_001", "name": "新名称",
	})
	assertToolName(t, cap, "update_base")
	assertToolArg(t, cap, "baseId", "BASE_001")
	assertToolArg(t, cap, "newBaseName", "新名称")
}

func TestAitableBaseUpdate_should_pass_description(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "update"}, map[string]string{
		"base-id": "B1", "name": "N", "desc": "备注",
	})
	assertToolArg(t, cap, "description", "备注")
}

// ── base delete ─────────────────────────────────────────────

func TestAitableBaseDelete_should_call_delete_base(t *testing.T) {
	cap := setupTestDepsAutoConfirm(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "delete"}, map[string]string{
		"base-id": "BASE_DEL",
	})
	assertToolName(t, cap, "delete_base")
	assertToolArg(t, cap, "baseId", "BASE_DEL")
}

func TestAitableBaseDelete_should_pass_reason(t *testing.T) {
	cap := setupTestDepsAutoConfirm(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "delete"}, map[string]string{
		"base-id": "B1", "reason": "不再需要",
	})
	assertToolArg(t, cap, "reason", "不再需要")
}

func TestAitableBaseDelete_should_not_call_mcp_in_dry_run(t *testing.T) {
	cap := setupTestDepsWithDryRun(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "base", "delete"}, map[string]string{
		"base-id": "B1",
	})
	assertCallCount(t, cap, 0)
}

// ── table get ───────────────────────────────────────────────

func TestAitableTableGet_should_call_get_tables(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "table", "get"}, map[string]string{
		"base-id": "BASE_001",
	})
	assertToolName(t, cap, "get_tables")
	assertToolArg(t, cap, "baseId", "BASE_001")
}

func TestAitableTableGet_should_pass_table_ids(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "table", "get"}, map[string]string{
		"base-id": "B1", "table-ids": "tbl1,tbl2",
	})
	assertToolArg(t, cap, "tableIds", []string{"tbl1", "tbl2"})
}

// ── table create ────────────────────────────────────────────

func TestAitableTableCreate_should_call_create_table(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "table", "create"}, map[string]string{
		"base-id": "B1", "name": "任务表",
		"fields": `[{"fieldName":"名称","type":"text"}]`,
	})
	assertToolName(t, cap, "create_table")
	assertToolArg(t, cap, "baseId", "B1")
	assertToolArg(t, cap, "tableName", "任务表")
}

// ── table update ────────────────────────────────────────────

func TestAitableTableUpdate_should_call_update_table(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "table", "update"}, map[string]string{
		"base-id": "B1", "table-id": "TBL_001", "name": "新表名",
	})
	assertToolName(t, cap, "update_table")
	assertToolArg(t, cap, "newTableName", "新表名")
}

// ── table delete ────────────────────────────────────────────

func TestAitableTableDelete_should_call_delete_table(t *testing.T) {
	cap := setupTestDepsAutoConfirm(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "table", "delete"}, map[string]string{
		"base-id": "B1", "table-id": "TBL_DEL",
	})
	assertToolName(t, cap, "delete_table")
	assertToolArg(t, cap, "tableId", "TBL_DEL")
}

func TestAitableTableDelete_should_pass_reason(t *testing.T) {
	cap := setupTestDepsAutoConfirm(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "table", "delete"}, map[string]string{
		"base-id": "B1", "table-id": "T1", "reason": "测试清理",
	})
	assertToolArg(t, cap, "reason", "测试清理")
}

// ── field get ───────────────────────────────────────────────

func TestAitableFieldGet_should_call_get_fields(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "field", "get"}, map[string]string{
		"base-id": "B1", "table-id": "T1",
	})
	assertToolName(t, cap, "get_fields")
	assertToolArg(t, cap, "baseId", "B1")
	assertToolArg(t, cap, "tableId", "T1")
}

func TestAitableFieldGet_should_pass_field_ids(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "field", "get"}, map[string]string{
		"base-id": "B1", "table-id": "T1", "field-ids": "fld1,fld2",
	})
	assertToolArg(t, cap, "fieldIds", []string{"fld1", "fld2"})
}

// ── field create ────────────────────────────────────────────

func TestAitableFieldCreate_should_call_create_fields(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "field", "create"}, map[string]string{
		"base-id": "B1", "table-id": "T1",
		"fields": `[{"fieldName":"状态","type":"singleSelect"}]`,
	})
	assertToolName(t, cap, "create_fields")
	assertToolArg(t, cap, "baseId", "B1")
	assertToolArg(t, cap, "tableId", "T1")
}

// ── field update ────────────────────────────────────────────

func TestAitableFieldUpdate_should_call_update_field(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "field", "update"}, map[string]string{
		"base-id": "B1", "table-id": "T1", "field-id": "FLD_001",
		"name": "新字段名",
	})
	assertToolName(t, cap, "update_field")
	assertToolArg(t, cap, "fieldId", "FLD_001")
	assertToolArg(t, cap, "newFieldName", "新字段名")
}

func TestAitableFieldUpdate_should_pass_config_json(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "field", "update"}, map[string]string{
		"base-id": "B1", "table-id": "T1", "field-id": "F1",
		"config": `{"options":[{"name":"X"}]}`,
	})
	expected := map[string]any{
		"options": []any{map[string]any{"name": "X"}},
	}
	assertToolArg(t, cap, "config", expected)
}

// ── field delete ────────────────────────────────────────────

func TestAitableFieldDelete_should_call_delete_field(t *testing.T) {
	cap := setupTestDepsAutoConfirm(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "field", "delete"}, map[string]string{
		"base-id": "B1", "table-id": "T1", "field-id": "FLD_DEL",
	})
	assertToolName(t, cap, "delete_field")
	assertToolArg(t, cap, "fieldId", "FLD_DEL")
}

// ── record query ────────────────────────────────────────────

func TestAitableRecordQuery_should_call_query_records(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "record", "query"}, map[string]string{
		"base-id": "B1", "table-id": "T1",
	})
	assertToolName(t, cap, "query_records")
	assertToolArg(t, cap, "baseId", "B1")
	assertToolArg(t, cap, "tableId", "T1")
}

func TestAitableRecordQuery_should_pass_record_ids(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "record", "query"}, map[string]string{
		"base-id": "B1", "table-id": "T1", "record-ids": "rec1,rec2",
	})
	assertToolArg(t, cap, "recordIds", []string{"rec1", "rec2"})
}

// ── record create ───────────────────────────────────────────

func TestAitableRecordCreate_should_call_create_records(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "record", "create"}, map[string]string{
		"base-id": "B1", "table-id": "T1",
		"records": `[{"cells":{"fld1":"hello"}}]`,
	})
	assertToolName(t, cap, "create_records")
	assertToolArg(t, cap, "baseId", "B1")
}

// ── record update ───────────────────────────────────────────

func TestAitableRecordUpdate_should_call_update_records(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "record", "update"}, map[string]string{
		"base-id": "B1", "table-id": "T1",
		"records": `[{"recordId":"rec1","cells":{"fld1":"updated"}}]`,
	})
	assertToolName(t, cap, "update_records")
}

// ── record delete ───────────────────────────────────────────

func TestAitableRecordDelete_should_call_delete_records(t *testing.T) {
	cap := setupTestDepsAutoConfirm(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "record", "delete"}, map[string]string{
		"base-id": "B1", "table-id": "T1", "record-ids": "rec1,rec2",
	})
	assertToolName(t, cap, "delete_records")
	assertToolArg(t, cap, "recordIds", []string{"rec1", "rec2"})
}

// ── template search ─────────────────────────────────────────

func TestAitableTemplateSearch_should_call_search_templates(t *testing.T) {
	cap := setupTestDeps(t, "aitable")
	root := buildRoot()
	_ = execCmd(t, root, []string{"aitable", "template", "search"}, map[string]string{
		"query": "项目管理",
	})
	assertToolName(t, cap, "search_templates")
	assertToolArg(t, cap, "query", "项目管理")
}
