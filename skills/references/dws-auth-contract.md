# dws Auth 契约（CLI stderr / exit code）

> **归属：** dws 仓库 skill。OpenClaw 钉钉场景下 **Agent exec** `auth status` / `auth login`（per-sender）。  
> 工作流见 [dws-auth-workflow.md](./dws-auth-workflow.md)。

## login / auth exit code（简表）

| exit | 含义 |
|------|------|
| `0` | 成功（login 时 token 已落盘） |
| `2` | auth 类失败（换票失败、保存 token 失败等） |
| `4` | PAT 相关（非 device login） |
| `5` | 业务命令 fail-closed：`IDENTITY_NOT_AUTHENTICATED`（该 `--sender-id` 无 token） |

**说明：** Agent 根据 `auth status` 与 stderr/JSON 错误码编排；`auth login --device` 需等待用户扫码，exec 勿用过短超时。

## CLI 组织权限（业务层，非 login 门禁）

OAuth 扫码成功并换票后，**token 会落盘**。组织是否开启 CLI、用户是否在 CLI 授权名单，由**业务 API** 在调用时返回权限错误；Agent 根据当次 stderr 引导用户联系管理员，**不能**从 `IDENTITY_NOT_AUTHENTICATED` 反推 CLI 名单状态。

## Ready 判据

```bash
dws auth status --sender-id <id> --format json
# authenticated: true
```
