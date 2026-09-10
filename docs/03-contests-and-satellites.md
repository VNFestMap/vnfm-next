# 03 赛事认证与外围系统

第三批：建立在 01/02 之上的赛事引擎、徽章认证、外部连携，以及可以独立运行的小产品和打包形态。重构主站时不必第一期就全部带走，但产品上它们已经和 VNFest 品牌绑在一起。

---

## 1. 统一投票底座

十二器和萌战不是两套无关系统，而是同一套「年度投票企划」核心（`includes/vote_projects.php`），萌战/十二器 PHP 只是别名层。

### 1.1 企划（Vote Project）

类型：`twelve` | `moe`

状态：`draft` → `published` → `running` → `ended` / `archived` / `suspended`

可见性：`public` | `unlisted` | `club_only`

投票资格：`club_member` | `public` | `invite_code` | `whitelist`

接口 `api/vote_projects.php`：`list`、`my_manageable`、`get`、`create`、`update`、`publish`、`suspend`、`archive`、`delete`

### 1.2 阶段（Stage）

类型：`nomination` 提名 → `qualifier` 海选 → `group_vote` 分组投票 → `bracket` 对阵 → `final` 决赛

阶段状态：`pending` | `open` | `locked` | `reviewing` | `settled`

投票模式：`nomination` | `multi_select` | `score` | `match_single`

结果可见性：`live_votes` | `live_rank_only` | `after_stage` | `after_event` | `hidden`

`api/vote_stages.php` 动作很多，覆盖运营全流程：

- CRUD、排序、`flow_status`
- 从提名重建并开启：`rebuild_from_nomination_and_open` / `rebuild_from_nomination`
- 海选池：`open_pool` / `settle_pool` / `generate_next_pool`
- 种子：`seed_entries` / `reseed_stage`
- 晋级：`advance_from_stage` / `advance_from_nomination`
- 平票：`resolve_flow_tie` / `resolve_tie`
- 开关：`open` / `lock` / `close` / `settle`
- 查看：`stage_entries`

### 1.3 提名与作品源

`api/vote_nominations.php`：用户提交/撤回、我的提名、汇总、管理员通过/拒绝、手动创建、导入、移除、恢复。

外部源：

- Bangumi 作品 / 角色（`api/bangumi_v0_search.php`、`api/vote_sources.php`）
- VNDB 作品 / 角色（`api/vndb_search.php`）
- 手动提名

条目状态（十二器作品）：`pending` / `approved` / `rejected` / `removed`

### 1.4 投票与对阵

`api/vote_votes.php`：`eligibility`、`cast`、`my_votes`、各类 `*_results`

`api/vote_matches.php`：列出/生成对阵、开启/锁定/结算、按票数结算。萌战淘汰赛用。

平票规则（十二器）：`same_rank` / `created_order` / `manual_review`

前端公共层：`js/vote-common.js`、`js/vote-activity.js`、`vote.html`（总入口）。

定时结算：`scripts/settle-moe-contests.php`。

---

## 2. 十二器（Twelve）—— 作品评选

目标：提名 → 海选 → 分组评分/投票 → 年度 Top 12 视觉小说。

页面：

| 路径 | 用途 |
|------|------|
| `twelve/index.html` | Hub：进行中 / 即将开始 / 已结束 |
| `twelve/contest.html` | 赛事详情、阶段、提名、投票 |
| `twelve/vote.html` | 投票页 |

样式/脚本：`twelve-hub.*`、`twelve-detail.*`、`js/twelve-contest.js`

兼容 API 仍在：`api/twelve_contests.php`、`twelve_rounds.php`、`twelve_votes.php`、`twelve_works.php`，内部转到投票底座。

---

## 3. 萌战（Moe）—— 角色对决

目标：提名 → 海选 → 2 的幂人数 1v1 淘汰 → 萌王。

页面：

| 路径 | 用途 |
|------|------|
| `moe/index.html` | Hub |
| `moe/contest.html` | 详情 |
| `moe/bracket.html` | 淘汰赛对阵图 |

对阵图：可缩放、平移、拖拽，实时显示每场投票状态和票数（`moe-bracket.js` / `moe-bracket.css`）。

兼容 API：`api/moe_contests.php`、`moe_stages.php`、`moe_votes.php`、`moe_matches.php`、`moe_candidates.php`

社团详情上的「萌王」是结果展示位，不是整场赛事。

---

## 4. 同好会试炼（Recognition）

让同好会自己发「试炼」，通过后签发可公开验证的徽章。架构按职责拆在 `includes/recognition/`：capability、roles、events、rules、credential、outbox、pipeline、signature。**只有 `credential.php` 可以创建或改凭证状态。**

### 4.1 用户侧

| 页面 | 用途 |
|------|------|
| `trial/index.html` | 试炼广场、详情、答题、社团管理入口 |
| `achievements.html` | 认可图鉴（我的徽章，可滤有效/历史） |
| `verify.html` | 公开核验，输入 `VNF-CRED-...`，不展示证据原文 |

试炼类型：

- `assessment` 知识试炼（答题）
- `activity` 活动签到
- `submission` 作品提交
- `award` 人工授予

广场可按类型筛选。答题有次数上限、冷却。作品提交走审核。兑换码支持「现场先参与、回去再领」。

参与接口 `api/recognition_participate.php`：`start`、`submit`（答题）、`submit_work`、`redeem`、`status`

凭证接口 `api/recognition_credentials.php`：`my`、`verify`、`set_visibility`、`revoke`、`grant`、`club_list`

### 4.2 同好会管理

`api/recognition_programs.php`：`list` / `detail` / `create` / `update` / `publish` / `set_status` / `manage` / `badge_create` / `badge_list`

规则以 **版本快照 JSON**（`recognition_program_versions.content_snapshot`）发布，已发布版本不可变。

可配：可见性、开放/关闭时间、最大尝试、冷却、最大签发量、凭证 TTL、安全级别、参与规则、能力列表。

徽章 `recognition_badges`：名称、分类（如 `participation`）、图片、版本。

管理接口 `api/recognition_admin.php`：审作品、生成/列出兑换码、导入参与者。

### 4.3 事件、连接器、异步

外部系统通过 Connector 上报事件，**禁止直接写凭证**。

- Connector：webhook、token 前缀+哈希、HMAC、scope、可吊销
- Event：`event_id` + `idempotency_key` 唯一，防重放
- Outbox：签发后的通知/过期扫描，由 `scripts/recognition_worker.php` 消费
- 身份绑定 `recognition_identity_links`（新外部身份走这张表）
- 同好会认可角色 `recognition_club_roles`（比会籍更细的 8 角色，MVP 多从会籍映射）

HMAC 密钥：`RECOGNITION_HMAC_SECRET`；凭证号前缀：`RECOGNITION_CRED_PREFIX`（默认 `VNF-CRED-`）。

公开验证只返回签发社团、时间、状态，不泄露证据。

---

## 5. MakoQuiz 答题连携

MakoQuiz 是独立抢答产品，和主站用 HMAC 绑账号。

- 配置：`QUIZ_LINK_SECRET`（bind_token 签名，与对方 `VNFEST_LINK_SECRET` 一致）、`QUIZ_API_KEY`（回传 Bearer）
- `api/quiz_auth.php`：签发绑定令牌，有效期与 7 天会话对齐
- `api/quiz.php`：`submit_results`（房间码、标题、分数、名次、人数、结束时间；`(room_code, user, ended_at)` 去重）、`my_results`
- 表 `quiz_results`
- 试炼模块也能吃答题事件：`api/recognition_events.php?action=quiz_sync`

使用文档在 Wiki Guide 的「连接产品 / 主办一场 Quiz」。

---

## 6. AstrBot 机器人

下载包：`downloads/astrbot_plugin_galgamemap_beta0.5.zip`。

主站接口 `api/bot.php`。两种鉴权：全站 `BOT_API_KEY`，或社团级 Token（`club_bot_tokens`，负责人在后台「Bot 接入」里签发/吊销）。

可查询（公开范围）：同好会列表/搜索/详情、分享摘要、活动、刊物、Wiki、星图联盟、萌战、公告、GalOnly 活动。

可写（需权限）：入会申请列表、自动通过、通过、拒绝；GalOnly 摊位/Staff 申请查看。

`full=1` 才能看隐私字段，社团 Token 默认被拦住。

动作列表见 `action=help`。契约测试：`scripts/test-astrbot-sync-contract.mjs`。

---

## 7. 同好会运行画像

独立工具 `club-operation-portrait/`。**不是管理系统，也不是排行榜。**

填问卷 → 六维雷达图（0–100）+ 规模概览 + 优劣 + 建议。

六维：组织稳定度、活动执行力、成员规模与参与度、内容沉淀力、外部连接力、传承持续力。

评分 = 数字基础分 + 状态修正。等级：成熟稳定 / 良好 / 发展中 / 待补足 / 待补充。

可选 LLM 顾问（DeepSeek / OpenAI / Claude）：

1. `llm-correction/start` 只出 3–5 个追问
2. `llm-correction/complete` 再出修正分析；越界修正会被丢掉
3. LLM 不可用时前端仍展示本地画像，不造假分数

登录用户可 `prefill` 自己的社团信息。另有 `rank-demo.html` / `Rank/Prototype.html` 原型，正式产品强调不做排名。

---

## 8. Galgame 同好会模拟器

独立 Vite 游戏 `Game/galgame_club_sim/`，广场有入口。

玩法概要（v7.2）：

- 48 周学年，春/暑/秋/寒，学期报告 + 学年总结
- 每周在 17 个行动里选择（6 个分类）
- 6 名核心角色卡 + 技能冷却 + 普通成员池
- 大型企划：社刊 / 出摊 / 跨校 / VN / 交接
- 16 个成就、6 条角色个人剧情、每周突发事件
- 6 种结局（含「维持之年」中性结局）+ 继承到下一周目
- 经费、疲劳、六维能力；免费行动仍耗时间
- 无障碍：`aria-live`、焦点循环、`prefers-reduced-motion`

附带角色卡生成器 `card-creator/`，数据可从 zine JSON 来。数值文档 `GAME_DATA.md`、策略 `STRATEGY_GUIDE.md`。

这是完整前端游戏，和主站账号几乎无耦合，重构主站时可以继续当静态子应用挂载。

---

## 9. 其他专题与工具页

| 路径 | 说明 |
|------|------|
| `JUYOU/HAIGUITANG.html` | 「第一届橘柚海龟汤推理大赛」活动海报/专题页 |
| `tools/pdf-reader/` | 刊物 3D 翻页阅读器 |
| `feedback.html` | 反馈（能力在 01，入口全站都能挂） |
| `api/extract.php` | 文本/信息抽取类辅助 |
| `api/narrative_auth.php` | 叙事/专题向鉴权（与主登录分离的活动页） |
| `api/test.php` | 开发探测，生产不应暴露 |

---

## 10. 多端打包（历史能力）

原 `package.json` 仍保留桌面 / Android 打包，但 Web 才是主产品：

- Electron 桌面：`npm run dev` / `build:exe`，应用名「Galgame同好会地图」
- Capacitor Android：`build:apk`、`cap:sync`
- 静态站点收集：`build:www`

重构后若还要客户端，应重新评估，不必原样搬 Electron 壳。

---

## 11. 超级控制台里和本批相关的部分

`admin/reviews.html` 的看板会统计审核趋势、地区分布；赛事本身的运营主要在各 Hub + `vote_stages` 管理动作，而不是单独再做一套「萌战后台页」。

GalOnly 审核见 02。试炼管理在 `trial/index.html` 的「管理我的同好会试炼」视图。

---

## 12. 原项目外围基础设施（迁移时要知道）

这些不是用户功能，但重写会碰到：

| 项 | 原做法 |
|----|--------|
| 迁移 | `scripts/migrate.php`，MySQL / SQLite 双方言；Forum schema 在 `Forum/includes/forum_schema.php` |
| 超管种子 | `scripts/seed_superadmin.php` |
| Docker | `Dockerfile` + `docker-compose.yml`，PHP 容器 + 反向代理 |
| CI | GitHub Actions 构建镜像 → GHCR → Watchtower 热更新 |
| 部署文档 | 原 `DEPLOY.md`（宝塔 + Docker + PHP 8.4 + MySQL） |
| 配置 | `config.example.php`：库、OAuth、邮件、Quiz、Recognition、LLM、Bot Key |
| 契约测试 | 40+ Node 脚本，覆盖投票流、萌战/十二器、隐私、i18n、Bot、Wiki、GalOnly Staff 等 |

---

## 13. 重构时建议保留的产品契约

1. 十二器与萌战共用阶段引擎，不要再做成两套状态机；差异只在条目类型（作品 vs 角色）和是否生成 1v1 bracket。
2. 提名优先 Bangumi，VNDB 和手动作补充。
3. 对阵图是萌战的体验核心，不能只做表格。
4. 试炼凭证：同一 `(program_version, user, badge)` 只保留一份 active；撤销不物理删除。
5. Connector 不能写凭证，只报送事件；事件必须幂等。
6. 公开验证页不得返回证据原文或隐私联系方式。
7. Bot 的 `full` 隐私字段与 Web 端 `info_hidden` 规则必须一致。
8. 模拟器、运行画像、PDF 阅读器可以继续当独立前端，不必强行塞进 Nuxt 业务路由，用子路径托管即可。
