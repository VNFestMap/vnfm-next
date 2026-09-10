# 原项目功能清单（重构对照）

本目录记录原站 **VNFest 地图**（`china-visualnovelcircle-maps`，线上 `https://www.map.vnfest.top`）的完整功能。`vnfm-next` 将用 **Nuxt + Go Fiber** 重写该平台。

**当前实现范围以 [05-v1-scope.md](./05-v1-scope.md) 为准。** 按模块划分看 [00-modules.md](./00-modules.md)。`01 / 02 / 03` 只是原站功能对照，不是这一期的待办。

| 文档 | 用途 |
|------|------|
| [05 第一期范围](./05-v1-scope.md) | 只做用户、同好会创建/加入、活动、通知；其余冻结 |
| [00 重构模块划分](./00-modules.md) | 5 个业务大模块 + 内核 + 卫星，依赖和落地顺序 |
| [01 核心发现与账号平台](./01-core-platform.md) | 地图、同好会、会籍、登录、用户中心、通知、超管壳（功能细节） |
| [02 社区运营与内容协作](./02-community-operations.md) | 活动、刊物、Wiki、企划、论坛、GalOnly、广场（功能细节） |
| [03 赛事认证与外围系统](./03-contests-and-satellites.md) | 投票、十二器、萌战、试炼、Quiz、Bot、模拟器（功能细节） |
| [04 nextmoe 接入对照](./04-nextmoe-architecture.md) | 只借账号：OIDC + 按需 catalog；图床自建 R2；本站权限自管；禁止 S2S、不改 infra |

---

## 原项目是什么

VNFest 不是单纯的社团目录。它把下面这条链路串在一个站点里：

```text
发现同好会 → 查看详情 → 申请加入 → 参与活动 → 投稿刊物 / Wiki → 参与企划赛事
```

产品自称 **中日高校 Galgame / 视觉小说同好会导航 · 社团运营 · 活动发布 · 刊物征稿 · 企划赛事 · Wiki 共建**。版本号约 **v2.1.0**。

## 原技术栈（将被替换）

| 层 | 原实现 | 本仓库目标 |
|----|--------|------------|
| 前端 | 根目录 HTML + Vanilla JS + D3.js 7；用户中心是 React 18 SPA | Nuxt |
| 后端 | PHP 8.x，`api/*.php` 零框架 | Go Fiber |
| 数据 | JSON 运行时文件 + SQLite / MySQL（PDO） | 待定（重构时需一并迁移） |
| 部署 | Docker + GitHub Actions + Watchtower | 待定 |

原站大量业务仍写在 JSON 文件里（同好会目录、活动、刊物、企划、反馈），另一部分在关系库里（账号、会籍、投票、论坛、GalOnly、试炼）。重构时必须按模块决定「继续文件存储」还是「全部入库」。

## 角色一览

| 角色 | 原系统标识 | 能做什么 |
|------|------------|----------|
| 访客 | `visitor` | 浏览地图、公开同好会、Wiki、赛事；`?guest=1` 可跳过登录墙 |
| 外交成员 | `external` | 跨社团联络身份（IEM），不算正式会籍，不能选为代表同好会 |
| 注册用户 / 成员 | `member` | 申请加入、报名、投稿、编辑 Wiki、投票、发帖 |
| 管理员 | `manager` | 协助负责人维护社团、审核成员、管理企划/刊物 |
| 负责人 | `representative` | 维护资料、发布活动/征稿、管理成员、创建企划、发起赛事 |
| 超级管理员 | `super_admin` | 全站审核、用户管理、GalOnly、赛事运营、公告、操作日志 |
| GalOnly 陪审 | `users.is_audit` | 摊位/Staff 审核投票，独立于系统角色 |

权限层级（原 `ROLE_HIERARCHY`）：`visitor(0) < external(0.5) < member(1) < manager(2) < representative(3) < super_admin(4)`。系统角色与会籍角色并存：一个人可以是全站 `visitor`，同时是某社团的 `representative`。

## 公开入口（原 Web 根）

原项目刻意把路由 HTML 放在站点根上，避免旧链接失效：

`index.html` · `login.html` · `user.html` · `club_square.html` · `club_share.html` · `vote.html` · `star_map.html` · `star_map_3d.html` · `submit.html` · `submit_event.html` · `submit_publication.html` · `feedback.html` · `achievements.html` · `verify.html`

子目录入口：`wiki/` · `Forum/` · `moe/` · `twelve/` · `Galgame_events/` · `trial/` · `admin/` · `Game/` · `club-operation-portrait/` · `tools/`

## 数据来源速查

- **JSON**：`data/clubs.json`、`data/clubs_japan.json`、`data/events.json`、`data/event_registrations.json`、`data/publications.json`、`data/projects.json` 及相关企划文件、`data/feedback.json`
- **关系库**：`users`、`sessions`、`club_memberships`、`notifications`、`galonly_*`、投票/`moe`/`twelve`、论坛、`recognition_*`、`quiz_results`、`star_unions`
- **上传**：`data/avatars/`、`data/club_avatars/`、`data/manuscripts/`、`data/publication_images/`、`uploads/galonly/`、`Forum/uploads/`、`wiki/uploads/`

详细字段与接口见三份子文档。
