# dws Auth 契约（CLI stderr / exit code）

> **归属：** dws 仓库 skill。Agent 从 `dws auth login --sender-id <id> --device` 及业务命令输出读取。  
> 命令规范见 [dws-auth-workflow.md](./dws-auth-workflow.md)「命令规范」。

## login / auth exit code（简表）

| exit | 含义 |
|------|------|
| `0` | 成功（login 时 token 已落盘） |
| `2` | auth 类失败（含 Step4 CLI 拒绝、login 失败） |
| `4` | PAT 相关（非 device login） |
| `5` | 业务命令 fail-closed：`IDENTITY_NOT_AUTHENTICATED`（该 `--sender-id` 无 token） |

**说明：** Agent 编排靠 **stderr 原文 + JSON 错误码**（如 `IDENTITY_NOT_AUTHENTICATED`），不依赖 connector 侧解析器。login 超时/授权码过期以 stderr 中文提示为准（如「操作超时」「授权码已过期」），重新 login 即可。

## Step4 CLI 拒绝（stderr 关键词）

device login Step4（`/cli/cliAuthEnabled`）失败时，CLI 输出中文说明（exit `2`），token **未落盘**：

| stderr 关键词 / 场景 | Agent 引导 |
|----------------------|-----------|
| 不在 CLI 授权人员范围 / `user_not_allowed` | 联系管理员在「开发者设置」加入 CLI 授权名单 |
| 尚未开启 CLI 数据访问权限 / `cli_not_enabled` | 联系管理员开启组织 CLI |
| 禁止所有成员使用 CLI / `user_forbidden` | 组织策略禁止，联系管理员 |
| 其它 auth 拒绝 | 引用 stderr 原文，勿猜测 |

业务 API 返回的权限错误以当次 stderr 为准，不能反推 CLI 名单。

## Ready 判据

```bash
dws auth status --sender-id <id> --format json
# authenticated: true
```
