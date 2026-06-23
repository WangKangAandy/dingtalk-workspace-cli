# Per-senderId OAuth (P1) — 本地开发说明

对应 OpenClaw roadmap：`~/.openclaw/vendor/dingtalk-openclaw-connector-fix-Community/docs/ROADMAP-per-senderId-oauth.md`

## 构建

```bash
cd /home/mccxadmin/.openclaw/vendor/dingtalk-workspace-cli
make build   # 或 ./scripts/dev/build.sh
```

## 新增能力（本分支 `feature/per-sender-id-oauth`）

| 能力 | 用法 |
|------|------|
| 多用户 token 目录 | `~/.dws/users/<senderId>/.data` |
| 全局 flag / env | `--sender-id <id>` 或 `DWS_AUTH_IDENTITY=<id>` |
| 登录 | `dws auth login --sender-id <id> --device` |
| 状态 / 登出 | `dws auth status --sender-id <id>` / `dws auth logout --sender-id <id>` |
| Fail-closed | 设了 identity 且无 token → `IDENTITY_NOT_AUTHENTICATED`（exit 5） |
| 本人扫码校验 | **仅直连 OAuth**（自有 AppKey/Secret 换得的 userAccessToken）：落盘前 `GET /v1.0/contact/users/me`（`x-acs-dingtalk-access-token`）与 `--sender-id` 比对 → `IDENTITY_MISMATCH`。**MCP 设备流**换得的 token 不能调 `api.dingtalk.com` Raw OpenAPI（见 `dws api` 说明），故该路径下跳过 OpenAPI 比对；需依赖 connector 侧绑定（roadmap §2.3）。 |

## 验收示例

```bash
export DWS_CONFIG_DIR=/tmp/dws-test
dws auth login --sender-id userA --device
dws auth status --sender-id userA --format json

DWS_AUTH_IDENTITY=userB dws contact user get-self --format json   # 应失败
DWS_AUTH_IDENTITY=userA dws contact user get-self --format json   # 成功（已登录时）
```

## 后续

- connector P2：`getDwsSpawnEnv` + `onCommandOutput` 自动补链（见 roadmap §2.3）
- 向 `DingTalk-Real-AI/dingtalk-workspace-cli` 提 PR 前请跑全量 `make test`
