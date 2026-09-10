# 01 核心发现与账号平台

第一批能力：让人能找到同好会、看懂资料、注册登录、加入组织、管理自己的账号。没有这部分，后面的活动、Wiki、赛事都没有主体。

对应原站主路径：`index.html` → 同好会详情 → `login.html` / `user.html` → 会籍申请 → 负责人后台。

---

## 1. 交互式地图与发现

原站的核心入口是全国 / 日本交互地图，D3.js 渲染 SVG 路径。默认未登录会跳到 `login.html`；加 `?guest=1` 进入访客浏览。

### 1.1 中国地图

- 覆盖全部省级行政区。
- 点击省份查看当地同好会列表。
- 地图模式与列表模式可切换。
- 省份索引、关键词搜索、类型筛选、多维排序。
- 同一同好会可绑定多个省份（`province` + `provinces[]`），跨地区组织能出现在多个省。
- 联系方式复制、快速查找。

筛选维度（列表态）：

| 项 | 说明 |
|----|------|
| 关键词 | 名称、学校、简介、`raw_text` |
| 类型 | `school` 学校社团 / 民间组织 / 线上社群 |
| 地区 | 中国 / 日本 |
| 排序 | 默认、按名称等 |

### 1.2 日本地图

- 精确到 47 都道府县（`js/japan.js` + `js/app-core.js` 中的别名表）。
- 独立数据源 `api/clubs_japan.php` ← `data/clubs_japan.json`。
- 中日名称互译与别名归一（「東京都 / 东京 / 東京」等）。
- 按地方分区：北海道・东北、关东、中部、关西・近畿、中国・四国、九州・冲绳。

### 1.3 江苏子地图

- 13 个地级市自包含离线 SVG，转成 `js/jiangsu.js`，运行时不请求外部地图。
- 城市：南京、无锡、徐州、常州、苏州、南通、连云港、淮安、盐城、扬州、镇江、泰州、宿迁。
- 超级管理员可做江苏城市批量维护（`clubs.php` 的 `jiangsu_city_bulk`）。
- 管理后台有「江苏专项」页签，仅超管可见。

### 1.4 亚洲 / 其他地区底图

- `js/asia.js` 提供亚洲范围路径，用于海外 / 更大视野展示。
- 同好会 `country` 目前实际落地的是 `china` 与 `japan`。

### 1.5 访客与登录墙

- 未登录访问 `index.html` 会请求 `api/auth.php?action=me`，失败则跳登录。
- `?guest=1` 跳过登录，可看地图和公开资料。
- 联系方式默认对非成员隐藏，文案变为「申请绑定后可见」。
- `protected` 保护模式：仅正式成员或负责人级以上可见联系方式。
- `visible_by_default` 可把部分信息公开给访客。

### 1.6 地图操作偏好

- 主题：亮 / 暗 / 跟随系统（`js/theme-runtime.js`，键 `themePreference`）。
- 语言：中文 / 日本語，可同步到账号（见第 6 节）。
- 地图反色控件（`js/display-preferences.js`，键 `vnfestMapInvertControls`）。
- v2.1 起这些控件从地图页拆走，统一收到用户中心「偏好设置」。

**原关键文件：** `index.html`、`js/app.js`、`js/app-core.js`、`js/china.js`、`js/japan.js`、`js/jiangsu.js`、`js/asia.js`、`js/calendar.js`（地图页内嵌活动日历）、`css/styles.css`

---

## 2. 同好会资料

每个同好会是一条目录记录，中国 / 日本分文件存储，主键是 `(id, country)`，不是全局自增唯一。

### 2.1 资料字段

| 字段 | 含义 |
|------|------|
| `id` | 国家内编号 |
| `country` | `china` / `japan` |
| `name` / `display_name` | 名称 |
| `province` / `provinces[]` | 主省份 + 多省份 |
| `prefecture` | 日本都道府县 |
| `city` | 城市（江苏子地图会用到） |
| `school` | 所属学校 |
| `type` | 组织类型，默认 `school` |
| `info` | 详细介绍 / 联系方式正文 |
| `remark` | 备注 |
| `logo_url` | 社团头像 |
| `external_links` | 外部链接 |
| `verified` | 是否已审核入库 |
| `protected` | 保护模式（联系方式仅成员可见） |
| `visible_by_default` | 是否默认可被非成员看到联系方式 |
| `project` | 项目标记，默认 `galgame` |
| `raw_text` | 搜索用拼接文本 |

接口在返回时还会附加：

- `info_hidden` / 脱敏后的 `info`
- `can_apply`：已登录且还不是该社正式成员则可申请
- `share_url`、`completeness`（资料完整度）
- `dynamic_summary`：活动数、刊物体、是否有 Wiki
- `public_contact`：按隐私规则裁剪后的公开联系方式

### 2.2 详情页能力

详情挂在地图页内，不是独立路由。能力包括：

- 名称、地区、学校、类型、介绍、头像。
- 负责人可上传并裁剪社团头像（Cropper.js + `api/club_avatar.php`）。
- **神器榜（推荐作品）**：从 Bangumi 搜作品挂到社团页，可增删排序（`api/club_recommendations.php`）。
- **萌王展示**：社团当前角色萌王（`api/club_moe_king.php`）。
- **留言板**：成员和访客留言、删除（`api/club_comments.php`，软删除）。
- 公开分享页 `club_share.html`：完整度、近期活动 / 刊物 / Wiki 摘要，给外部传播用。
- 成长摘要 `api/growth.php?action=club_summary`。

### 2.3 目录维护权限

| 操作 | 谁可以 |
|------|--------|
| 新增同好会 | 仅 `super_admin`（`POST api/clubs.php`） |
| 编辑资料 | 超管，或该社 `manager` / `representative` |
| 删除 | 超管 |
| 公开投稿新社团 | `submit.html` → 待审核队列，不直接写入正式目录 |

**原关键文件：** `api/clubs.php`、`api/clubs_japan.php`、`api/club_avatar.php`、`api/club_comments.php`、`api/club_recommendations.php`、`api/club_moe_king.php`、`club_share.html`、`includes/growth.php`、`includes/display_club.php`

---

## 3. 成员、申请与会籍

会籍存在关系库 `club_memberships`，与 JSON 目录通过 `(club_id, country)` 关联。

### 3.1 加入方式

1. **申请加入**：登录用户在详情页申请。可填 QQ / 联系方式、申请角色、是否在校生、加入方式、申请理由。状态 `pending`，等负责人审核。
2. **绑定码**：负责人生成一次性/限次码（`club_verification_codes`），用户兑换后直接加入。可设最大使用次数、过期时间、启用/吊销。
3. **外交申请（IEM）**：申请角色 `external`，走单独的「外交申请」队列，不算正式成员。
4. **外部社团名**：申请时可以填「我来自某某未入驻社团」。

申请成功会发站内通知；可选邮件通知（用户级开关 + 社团级收件人）。

### 3.2 会籍状态与角色

| 状态 | 含义 |
|------|------|
| `pending` | 待审核 |
| `active` | 有效 |
| 离开后 `left_at` | 退出或被踢，记录保留 |

会籍角色：`member` / `manager` / `representative` / `external`。

负责人可做：

- 通过 / 拒绝申请
- 查看成员列表、成员人数统计
- 改角色
- 踢人
- **转让负责人**（`transfer`）
- 设置「申请邮件接收人」
- 退出自己的会籍（`leave`）

唯一约束：`(user_id, club_id, country)`，同一人可同时加入中日两边同编号社团。

### 3.3 代表同好会

账号可从自己 **正式、活跃** 且角色为 `member` / `manager` / `representative` 的会籍里选一个公开展示，或选「不展示」。

- 出现在账号总览、论坛作者栏。
- 不公开其他会籍，也不把身份快照写进帖子。
- 退出、被踢、重新申请、会籍失效时自动清除。
- 字段：`users.display_membership_id`。

### 3.4 接口

`api/membership.php` 动作：`my`、`apply`、`set_application_email_recipient`、`members`、`club_member_counts`、`pending`、`approve`、`reject`、`leave`、`kick`、`change_role`、`transfer`。

绑定码：`api/club_codes.php` → `list` / `generate` / `revoke` / `redeem`。

---

## 4. 账号、认证与会话

### 4.1 本地账号

`api/auth.php`：

| 动作 | 说明 |
|------|------|
| `register_local` | 用户名 + 密码 + 邮箱验证码注册 |
| `send_register_code` | 发注册验证码 |
| `login_local` | 登录 |
| `logout` | 登出 |
| `me` | 当前用户 + 会籍摘要 |
| `change_password` | 改密 |
| `send_password_reset_code` / `reset_password` | 邮箱找回 |
| `send_code` / `bind_email` / `unbind_email` | 绑定/解绑邮箱 |
| `update_profile` | 昵称、简介 |
| `update_display_club` | 代表同好会 |
| `update_membership_application_email_preference` | 入会申请邮件开关 |
| `update_language_preference` | 语言 `zh` / `ja`，跨设备同步 |
| `oauth_config` | 前端查询 QQ / Discord 是否已配置 |

写操作校验同源（Origin / Referer）。登录有速率限制（`includes/rate_limit.php`）。

用户字段还包括：`username`、`nickname`、`avatar_url`、`profile_bio`、`email`、`email_verified_at`、`status`（`active` / `disabled` / `banned`）、`is_audit`。

头像：`api/avatar.php?action=upload`，存 `data/avatars/`。

### 4.2 第三方登录

- QQ OAuth：`api/auth.php?action=qq_auth` → `api/qq_callback.php`，可绑定 / 解绑。
- Discord OAuth：同样一对 `discord_auth` / `discord_callback.php`。
- 历史字段 `users.qq_openid`、`users.discord_id`；试炼模块另有 `recognition_identity_links` 作为新的外部身份表。

### 4.3 会话

- PHP Session，Cookie HttpOnly + SameSite=Lax，HTTPS 时 Secure。
- Cookie 生命周期 7 天（与 MakoQuiz `bind_token` 对齐，避免关浏览器就掉登录）。
- 另有 `sessions` 表（id、user_id、ip、ua、expires、is_valid），部分场景会用。
- 过渡期仍接受 `X-Admin-Token`（`LEGACY_AUTH_ENABLED`）。

### 4.4 登录页

`login.html`：注册、登录、找回密码、QQ / Discord 按钮、登录后 `redirect` 回原页。

---

## 5. 用户中心

v2 把原来的巨型 `user.html` 收成 React 18 SPA（源码 `user-v2-react/`，构建产物 `user-v2-assets/`，入口仍是根上的 `user.html`）。

侧栏：

| Tab | 内容 |
|-----|------|
| 总览 | 驾驶舱：资料完整度、代表同好会、未读通知、快捷入口 |
| 账户 | 昵称、简介、头像、邮箱、改密、QQ/Discord 绑定、代表同好会 |
| 偏好设置 | 语言、主题、地图反色；语言会写回账号 |
| 同好会 | 我的会籍、待审核（负责人）、退出/管理入口 |
| 通知 | 列表、已读、全部已读 |

快捷入口：活动投稿、GalOnly、同好会广场；负责人额外看到同好会管理、企划管理、刊物管理。

负责人总览会拉 `api/growth.php?action=owner_dashboard`：待审核人数、资料完整度、近期活动/刊物。

个人资料完整度按 7 项计分：昵称、头像、简介、邮箱、QQ、Discord、至少一条活跃会籍。

---

## 6. 全站偏好：语言、主题、壁纸

不是业务功能，但是所有页面的底座。

- **语言运行时** `js/language-runtime.js` + 词表 `language-catalog.js` + 静态日文 `language-static-ja.js` + 页面扫描 `page-i18n.js`。登录后以账号 `language_preference` 为准。
- **主题** `js/theme-runtime.js` + `css/theme-tokens.css`。Forum、广场、赛事页都复用。
- **壁纸** `js/page-background.js` + `api/backgrounds.php` + `image/background/`。720px 以下或粗指针设备关闭全屏壁纸。
- **顶栏归一** `js/site-header.js` + `css/site-header.css`：各页顶栏统一成「VNFest / 当前页」+ 返回地图 / 用户中心。

---

## 7. 通知与公告

### 7.1 站内通知

表 `notifications`。类型覆盖：入会申请、审核结果、踢出、角色变更、活动相关、论坛互动等。

接口 `api/notifications.php`：`list`、`count_unread`、`mark_read`、`mark_all_read`。

邮件：`includes/mailer.php`，支持 PHP `mail()` 或 SMTP。

### 7.2 全站公告

表 `announcements`。超管在超级控制台起草 / 发布 / 删除。类型 `info` 等，可设是否持久展示。公开接口 `api/announcements.php?action=active`。

---

## 8. 同好会管理后台（负责人侧）

入口 `admin/club_manager.html`。侧栏：

| 页签 | 功能 |
|------|------|
| 待审核 | 正式成员申请 |
| 外交申请 | `external` 申请 |
| 已通过 | 历史通过记录 |
| 成员 | 改角色、踢人、转让负责人 |
| 设置 | 社团资料、头像、联系方式可见性、申请邮件接收人 |
| 绑定码 | 生成 / 吊销 / 查看使用次数 |
| Bot 接入 | 社团级机器人 Token（详见 03） |
| 神器榜 | Bangumi 推荐作品（详见 02/03 交界，挂在社团页） |
| 企划枢纽 | 跳到企划管理（02） |
| 江苏专项 | 仅超管 |
| 用户管理 | 仅超管，跳超级控制台 |

顶栏可进「审核中心」`admin/reviews.html`。

---

## 9. 超级管理控制台

入口 `admin/reviews.html`，仅超管。

| 模块 | 功能 |
|------|------|
| 数据看板 | 待审核数量、近 30 天审核趋势图、同好会地区分布 Top 12 |
| 审核中心 | 同好会申请、活动申请、企划/刊物申请、成员绑定、反馈建议 |
| 用户管理 | 列表、改角色/状态、禁用 |
| 同好会总览 | 目录与成员概况 |
| 公告中心 | 起草、发布、删除公告 |
| 操作日志 | `audit_logs`（谁、对什么、何时、IP） |
| 数据导出 | 导出运营数据 |

用户 CRUD 也走 `api/users.php`：`stats` / `list` / `get` / `update` / `delete`。

公开投稿审核页还有：

- `admin/submissions.html` 同好会投稿
- `admin/submissions_event.html` 活动投稿
- `admin/events.php` 活动管理（偏运营）

---

## 10. 反馈

`feedback.html` + `api/feedback.php`。用户提交 bug / 功能建议，写入 `data/feedback.json`。超管在审核中心阅读。

---

## 11. 联合星图（发现向可视化）

把社团关系画成星空，偏探索，不是运营工具。

- `star_map.html`：Cinematic Frontend v2，径向渐变、雷达扫描、暗角；可沉浸式隐藏 HUD。
- `star_map_3d.html`：三维观测台变体。
- 数据：`star_unions`（联盟名称、描述、地区、绑定社团、星体颜色）+ `star_union_members`（联盟成员社团）。
- 接口 `api/star_unions.php`：`list` / `get` / `create` / `update` / `delete` / `add_club` / `remove_club` / `my_unions`。

---

## 12. 本批相关的横切能力

| 能力 | 说明 |
|------|------|
| 审计 | `includes/audit.php`，写 `audit_logs` |
| 限流 | `includes/rate_limit.php`，按 IP + 端点 |
| 健康检查 | `api/health.php`、`scripts/health-check.mjs` |
| 图片代理 | `api/image_proxy.php`，拉 Bangumi 等外链图 |
| 配置下发 | `api/get_config.php` |
| 可见性开关 | `api/toggle_visibility.php`、`scripts/set_contact_visibility.php` |
| 契约测试 | `scripts/test-*-contract.mjs` 四十余个，`npm run check` 全跑 |

---

## 13. 重构时建议保留的产品契约

1. `(club_id, country)` 永远成对出现，不要把中日目录合成一套自增 ID。
2. 访客可浏览；联系方式默认隐藏。
3. 系统角色 ≠ 会籍角色。
4. `external` 不是正式成员，不能当代表同好会。
5. 语言、主题、地图反色集中在账号/偏好，不要在每个页面再做一套开关。
6. 原目录在 JSON 里，成员在 SQL 里——迁移方案要先定，否则地图和会籍会对不上。
