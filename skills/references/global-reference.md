# 全局参考

## 认证

> **命令规范 + OpenClaw 工作流全文：** [dws-auth-workflow.md](./dws-auth-workflow.md)（唯一编排源）。  
> 本节为速查表，不重复工作流步骤。

### 命令速查（合并，按场景）

| 命令 | 适用场景 | 状态 | token 目录 |
|------|----------|------|------------|
| `dws auth status --sender-id <id> --format json` | **OpenClaw 钉钉 IM**、多用户 | ✅ **标准** | `~/.dws/users/<id>/` |
| `dws auth login --sender-id <id> --device` | **OpenClaw 钉钉 IM**、网关无头服务器 | ✅ **标准** | `~/.dws/users/<id>/` |
| `dws auth login` | 本机开发（有浏览器，loopback） | ⚠️ **仅本机 CLI** | `~/.dws/`（default） |
| `dws auth login --device`（无 `--sender-id`） | 运维 SSH 一次性初始化 dingmbw；本机无头 | ⚠️ **非 IM 用户授权** | `~/.dws/`（default） |
| `dws auth status` / `logout` / `reset`（无 `--sender-id`） | 本机单用户 default 身份 | ⚠️ **仅本机 CLI** | `~/.dws/` |
| 裸 `dws auth login` / 省略 `--sender-id` 的 login/status | OpenClaw **Agent** 处理钉钉聊天 | ❌ **废弃禁止** | — |

**OpenClaw 钉钉 `<id>`：** prompt 中 `[DingTalk DWS Context]` 的 `DWS_AUTH_IDENTITY`。Agent 必须显式写 `--sender-id`，不要只靠 env。

**两套目录并存的原因：** `~/.dws/` 存 default（运维/本机）；`~/.dws/users/<id>/` 存各聊天用户。设了 `DWS_AUTH_IDENTITY` 时 **fail-closed**，不会混用 default token。

**Agent 默认口令（OpenClaw 钉钉）：**

```bash
dws auth status --sender-id <DWS_AUTH_IDENTITY> --format json
dws auth login --sender-id <DWS_AUTH_IDENTITY> --device
```

### Token 生命周期

登录后自动刷新，日常使用无需重复登录。

| Token | 有效期 | 说明 |
|-------|--------|------|
| Access Token | 2 小时 | 调用 API 的凭证，过期自动刷新 |
| Refresh Token | 30 天 | 换新 Access Token，使用后轮转 |

30 天内使用一次即自动续期。`refresh_token` 单设备独占，远程刷新后源设备凭证失效。

### 认证失败

| 错误 | 处理 |
|------|------|
| `IDENTITY_NOT_AUTHENTICATED` / `AUTH_TOKEN_EXPIRED`（OpenClaw 钉钉） | `dws auth login --sender-id <DWS_AUTH_IDENTITY> --device` |
| `DWS_AUTH_DENIAL reason=*` | 按 reason 引导，见 [dws-auth-contract.md](./dws-auth-contract.md) |
| `AUTH_TOKEN_EXPIRED` / `USER_TOKEN_ILLEGAL`（本机 default） | `dws auth login` 或 `dws auth login --device` |
| HTTP 403 / scope 不足 | 联系管理员开权限，**不要**反复 login |

**CI/CD：** 可 `export DWS_CLIENT_ID` / `DWS_CLIENT_SECRET` 后 `dws auth login --device`（default 身份）；与 OpenClaw per-sender 授权无关。

## Recovery

当 runtime/MCP 命令失败且 stderr 额外输出 `RECOVERY_EVENT_ID=<event_id>` 时，说明 CLI 已经持久化了失败快照，可进入 recovery 闭环：

```bash
dws recovery plan --event-id <event_id> --format json
dws recovery execute --event-id <event_id> --format json
dws recovery finalize --event-id <event_id> --outcome recovered|failed|handoff --execution-file execution.json --format json
```

- `plan` / `execute` 也支持 `--last`，但 `--last` 与 `--event-id` 互斥
- recovery 文件保存在 `DWS_CONFIG_DIR/recovery/`
- CLI 会自动清理 30 天前的 recovery 文件和事件记录
- recovery 自己发起的文档检索与只读 probe 不会再创建新的 recovery 事件

更多闭环要求见 [recovery-guide.md](./recovery-guide.md)。


## 全局标志

| 标志 | 短名 | 说明 | 默认 |
|------|:---:|------|------|
| `--format` | `-f` | 输出格式: json / table / raw | json |
| `--jq` | | jq 表达式过滤输出 (如: `.items[] \| .name`) | 无 |
| `--fields` | | 筛选输出字段 (逗号分隔, 如: name,id,status) | 无 |
| `--verbose` | `-v` | 详细日志 | false |
| `--debug` | | 调试日志 | false |
| `--yes` | `-y` | 跳过确认提示 | false |
| `--dry-run` | | 预览操作不执行 | false |
| `--timeout` | | HTTP 超时 (秒) | 30 |
| `--mock` | | Mock 数据 (开发用) | false |
| `--client-id` | | 覆盖 OAuth Client ID | 无 |
| `--client-secret` | | 覆盖 OAuth Client Secret | 无 |

## 输出格式

### --format json (机器可读, 默认)

```json
{"success": true, "body": {...}}
```

### --format table (人类可读)

```
已创建 AI 表格 "项目管理" (UUID: abc123)

下一步:
  dws aitable base get --base-id abc123
```

## 环境变量

| 变量 | 说明 |
|------|------|
| `DWS_CONFIG_DIR` | 覆盖默认配置目录 |
| `DWS_SERVERS_URL` | 自定义服务发现端点 |
| `DWS_CLIENT_ID` | 覆盖 OAuth Client ID (DingTalk AppKey) |
| `DWS_CLIENT_SECRET` | 覆盖 OAuth Client Secret (DingTalk AppSecret) |

凭证优先级: `--token` > `DWS_CLIENT_ID`/`DWS_CLIENT_SECRET` > OAuth 加密存储 (.data)
