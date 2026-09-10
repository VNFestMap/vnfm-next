# 旧站数据导入

旧站目录在 JSON 文件里，会籍 / 用户 / 绑定码 / 通知在 MySQL dump 里。两边用 `(club_id, country)` 对上。新库同好会是全局 `id`，用户主键是 NextMoe OIDC `id`，所以不能整库直灌。

脚本：`scripts/import_legacy.py`。默认 dry-run，写出 `tmp/legacy-import/report.json`、`import.sql`、`users_extracted.csv`。

## 要准备的文件

| 文件 | 来源 | 说明 |
|------|------|------|
| MySQL dump | 已有 `data/www_test_map_vnf_2026-09-09_04-00-01_mysql_data.sql` | 用户、会籍、绑定码、通知 |
| `clubs.json` | 旧站服务器 `data/clubs.json` | **不要用公开 API**。未登录的 `api/clubs.php` 会把 `info` 打成「申请绑定后可见」 |
| `clubs_japan.json` | 旧站 `data/clubs_japan.json` | 同上 |
| `events.json` | 旧站 `data/events.json` | 活动没有 `club_id`，脚本按标题/raw_text 撞同好会名 |
| `event_registrations.json` | 旧站 | 可选 |
| user map CSV | 脚本生成后你填 `oidc_id` | `legacy_id,email,username,oidc_id` |

## 命令

```bash
# 1) 先跑 schema
cd apps/api && pnpm migrate

# 2) dry-run（把 JSON 换成你从旧站拷来的未脱敏文件）
python3 scripts/import_legacy.py \
  --clubs-json /path/to/clubs.json \
  --clubs-japan-json /path/to/clubs_japan.json \
  --events-json /path/to/events.json \
  --registrations-json /path/to/event_registrations.json

# 3) 打开 tmp/legacy-import/users_extracted.csv，按邮箱填 NextMoe 用户 id
# 4) 确认 report.json 后真正写入（owner-id 必须是已经登录过本站的 OIDC 用户）
python3 scripts/import_legacy.py \
  --clubs-json /path/to/clubs.json \
  --clubs-japan-json /path/to/clubs_japan.json \
  --events-json /path/to/events.json \
  --user-map tmp/legacy-import/users_extracted.csv \
  --owner-id 2 \
  --apply
```

`--apply` 会跑 `002_legacy_import` 之后生成的 SQL。没跑过 `pnpm migrate` 会缺 `legacy_user_map` 表。

## 默认会迁 / 不会迁

迁：同好会目录、active/pending 会籍（用户已映射时）、绑定码、能匹配到同好会的活动、对应报名、会籍类通知（`join_approved` / `join_rejected` / `member_kicked` / `role_changed`）。

不迁：密码、QQ/Discord、全站公告（`notifications.type=system`，约 1.4 万条）、GalOnly/论坛/企划/赛事通知、rejected/left/kicked 会籍、远程头像文件（只记 URL 到 report）。

## 还没拍板的点

脚本按下面的默认做了。不对的话改默认再跑，不要先 `--apply`。

1. 旧 `users.id` 不是 NextMoe OIDC id，不能当新主键。会籍/码/通知只在 CSV 填了 `oidc_id` 之后才写入。
2. 公开 `api/clubs.php` 会脱敏联系方式。必须用服务器上的 `data/clubs.json` / `clubs_japan.json`。
3. 活动 JSON 没有 `club_id`。默认用标题/学校名匹配；对不上的跳过。
4. 会籍只进 `active` / `pending`。`rejected` / `left` / `kicked` 丢掉。
5. 通知默认只进会籍四类，丢掉约 1.4 万条 `system` 全站公告，以及 GalOnly/论坛/企划。
6. 头像 URL 不进 R2，只写进 report。
7. 原站 `super_admin` 不授予本站权限。
8. 多省份同好会只保留主省，其余写进 `remark`。
9. `province=海外` 的两条：目录能进，但地图没有海外层。
10. `--apply` 需要 `--owner-id`（已经在本站登录过的 OIDC 用户），给没有映射负责人的社团当 `created_by`。
