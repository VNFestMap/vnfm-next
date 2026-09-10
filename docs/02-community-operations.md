# 02 社区运营与内容协作

第二批能力：同好会「活起来」之后要办活动、出刊物、写百科、做企划、在论坛说话，以及办线下摊位（GalOnly）。这些功能依赖 01 的账号与会籍，但彼此相对独立，可以按模块迁移。

---

## 1. 活动日历

地图页内嵌双栏日历（`js/calendar.js`），也有独立投稿入口。

### 1.1 浏览

- 日历视图 / 列表视图。
- 按日期点选、按类型过滤。
- 公开列表：`GET api/events.php?action=list`。
- 数据文件 `data/events.json`。

### 1.2 发布与管理

- 负责人 / 超管：`add` / `update` / `delete` / `replace`（整表替换，需权限）。
- 独立投稿页 `submit_event.html` → `api/submit_event.php`，进入待审核，后台 `admin/submissions_event.html`。
- 活动可关联发起同好会。

### 1.3 报名

- 登录用户 `register` / `unregister`。
- 报名记录 `data/event_registrations.json`。
- 用户中心展示「我报名的活动」。
- 报名有锁定逻辑（契约测试 `test-event-registrations-lock.mjs`）：活动开始或截止后不能再改。

日历和企划枢纽会合并展示部分活动（`test-events-merge.mjs`）。

**原关键文件：** `js/calendar.js`、`api/events.php`、`submit_event.html`、`admin/events.php`、`admin/submissions_event.html`

---

## 2. 刊物征稿与稿件

### 2.1 征稿启事

- 数据 `data/publications.json`。
- 接口 `api/publications.php`：公开 GET；登录后 POST/PUT/DELETE。
- 权限：超管，或刊物关联社团的 `manager` / `representative`。
- 刊物可绑多个 `club_ids`，旧数据用 `clubName` 反查。
- 投稿状态：待审核 / 已录用 / 已发布。
- 公开投稿页 `submit_publication.html` → `api/submit_publication.php`。
- 管理页：`admin/publication_manager.html`、`wiki/publication-manage.html`。

### 2.2 稿件文件

`api/manuscripts.php`：

- `list_by_publication` / `list_by_club`
- `upload` / `delete` / `download`
- 文件落在 `data/manuscripts/`

### 2.3 资料公开库（数字存档）

面向已经出版的同好会刊物，不是征稿流程。

- 浏览页 `wiki/publications.html`
- 上传页 `wiki/publication-upload.html`
- 管理页 `wiki/publication-manage.html`
- 支持 PDF 和页面图片。
- 元数据：标题、作者、发布日期、关联同好会。
- 预览接口 `api/publication_previews.php`：`list` / `book_data` / `create` / `upload_page` / `publish`
- 在线 3D 翻页阅读：`tools/pdf-reader/?preview_id=`（上游 pdfReader）
- 图片存 `data/publication_images/`
- 另有一次迁移脚本 `api/migrate_publications.php`

---

## 3. Wiki 百科

百科式布局：左侧导航 + 右侧正文。内容按同好会一篇，生成静态 HTML。

### 3.1 阅读

- 首页 `wiki/index.html`：最近更新、中国 / 日本、全部页面、文档库、资料公开库、维护中心。
- 搜索同时过滤地区索引、页面和文档库。
- 索引来源：`wiki/index.json`、`wiki/library/index.json`、`wiki/feature-slots.json`；读取失败有静态兜底。
- 单页：`wiki/pages/{country}-{id}.html`，正文 JSON 在 `wiki/content/`。
- 中日双语内容（编辑器里分「中文内容 / 日本語内容」两套字段）。

### 3.2 编辑器

`admin/wiki_editor.html?club_key=china-12`（或 `japan-1`）。

可编辑：

- 标题、摘要
- 信息框（学校、地区、成立时间、状态等 KV）
- 章节（每行一段）
- 图片：上传 / 外链 / 站内路径，单张 ≤ 10MB，可设宽度、对齐、裁切
- 参考资料
- 基础模板插入、JSON 导入导出
- 本地草稿（`localStorage`）
- 实时预览

保存走 `api/wiki.php?action=save`，同时写出 content JSON 和静态 HTML。图片 `action=upload`，读取 `action=get`。

权限：超管，或该社团 `manager` / `representative`。生成脚本也可离线跑 `scripts/generate-wiki-pages.mjs`。

### 3.3 使用文档（Guide）

独立文档站 `wiki/guide/`，中日分组导航 + 站内检索。

分组：

1. 开始使用（入门、账号与角色）
2. VNFest 主站（发现社团、社区流程、Wiki 与资料库、投票与 GalOnly）
3. 运营与管理（负责人、全站管理）
4. 连接产品（MakoQuiz、AstrBot）
5. 技术与运维
6. 历史更新记录

种子数据 `wiki/guide/seed/{zh-CN,ja-JP}/documents.json`。管理页 `admin/wiki_guide_editor.html`。

Guide API（`api/wiki.php`）：`guide_catalog`、`guide_article`、`guide_admin_catalog`、`guide_diff`、`guide_upload`、`guide_save_draft`、`guide_publish`、`guide_unpublish`、`guide_reset_seed`。

资料库还有写作说明 `wiki/library/wiki-writing-guide.html`。

---

## 4. 企划枢纽

同好会协作中心，数据在 JSON：`data/projects.json`、`project_items.json`、`project_participations.json`，文件 `api/project_files.php`。

### 4.1 企划本体

类型：`publication` 刊物 / `activity` 活动 / `content` 内容征集 / `recruit` 协力招募 / `other`。

状态：`draft` → `collecting` → `ongoing` → `completed` / `archived`。软删除。

字段要点：标题、类型、状态、发起社团 `organizer_club {id, country}`、介绍、时间。

接口 `api/projects.php`：列表可按 type/status 过滤；登录后创建/更新。权限：发起社团的管理者或超管。

### 4.2 企划项（Item）

一个企划下可挂多种参与入口：

| item 类型 | 含义 |
|-----------|------|
| `submission` | 稿件投稿 |
| `registration` | 活动报名 |
| `collaboration` | 申请协力 |
| `survey` | 问卷 |
| `voting` | 投票 |
| `other` | 其他 |

### 4.3 参与记录

状态：`submitted` / `reviewing` / `accepted` / `rejected` / `withdrawn`。

角色分配在产品叙述里是策划 / 美术 / 文案 / 技术；实现上主要通过参与记录和成员角色体现。

### 4.4 前端

- 地图页内嵌枢纽 UI：`js/project-hub.js`
- 负责人管理：`admin/club_project_manager.html` + `js/club-project-manager.js`
- 广场会公开展示进行中的企划

共享 PHP：`includes/project_hub.php`（JSON 读写带文件锁）。

---

## 5. 论坛

自包含模块，全部在 `Forum/`。当前产品 **只开放统一广场**，历史「同好会分区」数据保留但全部 404。

### 5.1 页面

| 页面 | 用途 |
|------|------|
| `forum-plaza.html` | 广场、最新回复、我的发帖、我的收藏、搜索、分类、排序、分页 |
| `forum-post.html?id=` | 主帖 + 连续楼层 |
| `forum-create.html` | 发帖；`?edit=` 编辑自己的帖 |

访客可浏览搜索；登录可发帖、回复、点赞、收藏、举报。广场管理只给 `super_admin`。`manager` / `representative` 不再有论坛管理权。

### 5.2 分类与内容规则

固定五类：综合讨论、资源分享、活动发布、作品交流、求助答疑。

- 标题 ≤ 100 字，正文 ≤ 50,000 字
- 每帖最多 5 个标签，每个 ≤ 20 字
- 正文只存 Markdown；原始 HTML 转义
- 链接仅 HTTP(S)；图片仅本用户上传并绑定的 `Forum/uploads/...`
- Markdown：H2/H3、粗斜体、列表、引用、行内/围栏代码、链接、站内图、`@提及`
- 图片 JPEG/PNG/GIF/WebP，单张 ≤ 10MB，每篇最多 20 张
- 编辑时按正文引用同步附件；去掉的图归档而不是物理删

### 5.3 互动与管理

- 点赞（`forum_reactions`）、收藏（`forum_favorites`）
- 回复可挂 `parent_reply_id`（同帖、未删除）
- 举报（`forum_reports`）带内容快照，超管处理
- 修订历史 `forum_revisions`
- 置顶、精华、浏览/回复/点赞/收藏计数
- 作者对象动态带可空 `display_club`（代表同好会 · 会籍角色）
- 列表有纯文本摘要（约 180 字）和首图 `preview_image`
- 页脚深链接：`#rules` / `#about` / `#privacy`

### 5.4 写作画布

受控 `contenteditable`，序列化成 Markdown；可切源码模式。工具栏、粘贴清洗、中文 IME 保护、图片即时 blob 缩略图、草稿 `localStorage`。快捷回复是精简工具栏的同一套渲染器。

### 5.5 休眠的分区功能

旧 `scope=club` 的分类、帖子、回复、附件、收藏、通知仍在库里，API 对所有角色返回 404。URL 上的 `scope=club&club=&country=` 会 `replaceState` 清掉并给一次中性提示。将来若恢复分区，必须重做权限和国家隔离，不能只把入口打开。

**表：** `forum_categories`、`forum_posts`、`forum_replies`、`forum_attachments`、`forum_reactions`、`forum_favorites`、`forum_reports`、`forum_revisions`

---

## 6. GalOnly（高校专属摊位与 Staff）

线下展会通道，独立专题，不是普通活动日历。已落地上海、北京等场次。

入口：`Galgame_events/galgameonly_list.html`（高校专属通道列表）。

### 6.1 摊位申请

公开页按城市拆：

- 上海：`Shanghai_Galonly_submit.html`
- 北京第二届：`Beijing_Galonly_submit.html`
- 通用列表 / 海报：`galgameonly_list.html`、`galo_poster.html`

申请字段（随版本叠加）：联合摊位、摊位名、联系方式、QQ、手机、备注、主图、附件、摊位类型、预计人数、布局说明、是否用电、参展经历、商品清单及附件。

**两阶段审核 v2：**

- `phase` 1 / 2
- 每阶段陪审投票 + 意见（`galonly_votes`，唯一键含 phase）
- 反馈 `phase1_feedback` / `phase2_feedback`
- 状态 `pending` 等；可驳回、要求修改、标记 `resubmitted` / `has_update`
- 公开侧可按 IP 投人气票（`galonly_public_votes`，每活动每 IP 每摊位一条）

资格检查 `check_eligibility`。图片上传 `upload_image`，文件在 `uploads/galonly/`。

### 6.2 Staff 招募

独立表单与须知：

- `galonly_staff_submit.html` / `galonly_staff_guidelines.html`
- 上海、北京各有一份 Staff 页

申请字段：姓名、QQ、手机、邮箱、所属社团、意向岗位（多选 JSON）、能否保证档期、是否 Cosplay、是否三天都在、自我介绍、性别、有无 Staff 经验、技能。

活动级 Staff 配置：报名开关、截止时间、人数上限、需求人数、名单是否锁定。

审核：陪审投票、撤回票、改状态、锁定/解锁排班名单。同一活动同一用户在 `pending`/`pooled` 时不能重复有效申请（生成列 `active_key`）。

### 6.3 审核后台

- `admin/Galonly_audit.html`：摊位审核
- `api/galonly.php`：活动 CRUD、摊位全流程、Staff 全流程、公开投票
- `api/galonly_staff.php`：Staff 专用动作（`get_my`、`submit`、`list_applications`、`vote`、`withdraw_vote`、`update_status`）
- 陪审名单 `galonly_reviewers`（按活动绑定 `user_id` + `role`）
- 用户 `is_audit=1` 才能投审核票

活动表 `galonly_events`：名称、地点、日期、报名开关、Staff 开关、活动代码、描述、海报 `image_url`。

### 6.4 前端资源

海报、城市 Logo（北京 / 上海 / 凌夏 / VNF）在 `Galgame_events/assets/`。

---

## 7. 同好会广场

`club_square.html` + `js/club-square.js` + `css/club-square.css`。

信息汇聚页，深空主题 + 五色活动卡片：

| 色/类型 | 内容 |
|---------|------|
| 萌 | 萌战 |
| 12 | 十二器 |
| 青 | 雷达 / 其他企划 |
| 暖 | 同好会模拟器 |
| 紫 | MakoQuiz |

筛选、关键词、状态（草稿/已发布/进行中/结束）。点进具体赛事走 03 的 Hub。中日双语、亮暗主题。

广场也承担「看正在进行的社区项目」的入口，和企划枢纽、投票底座相连。

---

## 8. 通用投稿入口

| 页面 | 投什么 | 后端 |
|------|--------|------|
| `submit.html` | 新同好会信息 | `api/submit.php` → 审核队列 |
| `submit_event.html` | 活动 | `api/submit_event.php` |
| `submit_publication.html` | 刊物征稿 | `api/submit_publication.php` |

未入驻组织主要通过 `submit.html` 进目录，不直接写 `clubs.json`。

---

## 9. 社团页上的社区插件

这些挂在 01 的详情页上，但属于社区互动：

- **留言板** `club_comments`：列表 / 发表 / 软删
- **神器榜** Bangumi 作品推荐：`list` / `add` / `remove` / `reorder`，代理 `api/bangumi_proxy.php`、`api/bangumi_v0_search.php`
- **萌王** 当前角色展示
- **公开分享卡** `club_share.html`：给 QQ 群、海报、机器人用的摘要页

VNDB 也有代理：`api/vndb_proxy.php`、`api/vndb_search.php`（更多被 03 的提名用）。

---

## 10. 本批管理入口汇总

| 后台页 | 用途 |
|--------|------|
| `admin/club_project_manager.html` | 企划 CRUD、参与审核、文件 |
| `admin/publication_manager.html` | 刊物 |
| `admin/wiki_editor.html` | 社团 Wiki |
| `admin/wiki_guide_editor.html` | 使用文档 |
| `admin/Galonly_audit.html` | GalOnly 摊位/Staff |
| `admin/submissions.html` | 同好会投稿 |
| `admin/submissions_event.html` | 活动投稿 |
| `admin/events.php` | 活动运营 |
| `admin/reviews.html` 的审核中心 Tab | 把上述队列收口 |

---

## 11. 重构时建议保留的产品契约

1. 活动报名与征稿状态机要可审计（谁在何时通过/拒绝）。
2. Wiki 保存 = 结构化 JSON + 可读静态页，便于备份和离线。
3. Forum 图片必须走站内上传白名单，禁止外链图和原始 HTML。
4. GalOnly 两阶段审核 + 陪审分阶段投票是硬需求，北京/上海表单字段不同，最好做成「活动配置驱动表单」而不是再复制 HTML。
5. Forum 的 `scope=club` 历史数据不要在迁移时物理删除，除非明确做归档方案。
6. 企划枢纽目前是 JSON 文件锁；Go 侧重写时建议直接入库，避免多实例文件锁。
