# dws Auth 契约（CLI stderr / exit code）

> **归属：** dws 仓库 skill。Agent 从 `dws auth login --sender-id <id> --device` 输出读取。  
> 命令规范见 [dws-auth-workflow.md](./dws-auth-workflow.md)「命令规范」。

## login exit code

| exit | 含义 |
|------|------|
| `0` | token 落盘 |
| `2` | Step4 / auth 类拒绝（见 `DWS_AUTH_DENIAL`） |
| `5` | device login 超时 |
| `4` | PAT 相关（非 device login） |

## DWS_AUTH_DENIAL

device login Step4（`/cli/cliAuthEnabled`）失败时，stderr 追加：

```text
DWS_AUTH_DENIAL reason=user_not_allowed
```

| reason | 含义 |
|--------|------|
| `user_not_allowed` | 不在 CLI 个人授权名单 |
| `cli_not_enabled` | 组织未开启 CLI |
| `user_forbidden` | 组织全员禁用 |
| `auth_denied` | 其它 exit 2 兜底 |

**注意：** Step4 拒绝时 token **未落盘**。业务 API 返回的权限错误以当次 stderr 为准，不能反推 CLI 名单。

## Ready 判据

```bash
dws auth status --sender-id <id> --format json
# authenticated: true
```
