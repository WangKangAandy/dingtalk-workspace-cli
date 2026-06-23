package cli_compat_test

import (
	"testing"
)

// ── message send ────────────────────────────────────────────

func TestDingMessageSend_should_call_send_ding_message(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	err := execCmd(t, root, []string{"ding", "message", "send"}, map[string]string{
		"robot-code": "RC001",
		"users":      "userId1,userId2",
		"content":    "请查看",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertToolName(t, cap, "send_ding_message")
	assertToolArg(t, cap, "robotCode", "RC001")
	assertToolArg(t, cap, "content", "请查看")
}

func TestDingMessageSend_should_parse_receiver_list(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "send"}, map[string]string{
		"robot-code": "R", "users": "u1,u2,u3", "content": "C",
	})
	assertToolArg(t, cap, "receiverUserIdList", []string{"u1", "u2", "u3"})
}

func TestDingMessageSend_should_use_sms_type(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "send"}, map[string]string{
		"robot-code": "R", "type": "sms", "users": "u1", "content": "C",
	})
	assertToolArg(t, cap, "remindType", 2)
}

func TestDingMessageSend_should_use_call_type(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "send"}, map[string]string{
		"robot-code": "R", "type": "call", "users": "u1", "content": "紧急告警",
	})
	assertToolArg(t, cap, "remindType", 3)
	assertToolArg(t, cap, "content", "紧急告警")
}

func TestDingMessageSend_should_fallback_to_app_for_unknown_type(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "send"}, map[string]string{
		"robot-code": "R", "type": "unknown", "users": "u1", "content": "C",
	})
	assertToolArg(t, cap, "remindType", 1)
}

func TestDingMessageSend_should_trim_receiver_spaces(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "send"}, map[string]string{
		"robot-code": "R", "users": " u1 , u2 ", "content": "C",
	})
	assertToolArg(t, cap, "receiverUserIdList", []string{"u1", "u2"})
}

func TestDingMessageSend_should_make_single_call(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "send"}, map[string]string{
		"robot-code": "R", "users": "u1", "content": "C",
	})
	assertCallCount(t, cap, 1)
}

func TestDingMessageSend_should_not_call_mcp_in_dry_run(t *testing.T) {
	cap := setupTestDepsWithDryRun(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "send"}, map[string]string{
		"robot-code": "R", "users": "u1", "content": "C",
	})
	assertCallCount(t, cap, 0)
}

// ── message recall ──────────────────────────────────────────

func TestDingMessageRecall_should_call_recall_ding_message(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	err := execCmd(t, root, []string{"ding", "message", "recall"}, map[string]string{
		"robot-code": "RC001",
		"id":         "DING_MSG_001",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertToolName(t, cap, "recall_ding_message")
	assertToolArg(t, cap, "robotCode", "RC001")
	assertToolArg(t, cap, "openDingId", "DING_MSG_001")
}

func TestDingMessageRecall_should_handle_different_ids(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "recall"}, map[string]string{
		"robot-code": "R", "id": "MSG_XYZ",
	})
	assertToolArg(t, cap, "openDingId", "MSG_XYZ")
}

func TestDingMessageRecall_should_make_single_call(t *testing.T) {
	cap := setupTestDeps(t, "ding")
	root := buildRoot()
	_ = execCmd(t, root, []string{"ding", "message", "recall"}, map[string]string{
		"robot-code": "R", "id": "M1",
	})
	assertCallCount(t, cap, 1)
}

func TestDingMessageRecall_should_return_no_error(t *testing.T) {
	_ = setupTestDeps(t, "ding")
	root := buildRoot()
	err := execCmd(t, root, []string{"ding", "message", "recall"}, map[string]string{
		"robot-code": "R", "id": "M1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
