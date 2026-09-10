# nextmoe 接入对照（vnfm-next 怎么借账号）

`vnfm-next` **不是** nextmoe 生态的下游产品。论坛（`kun-galgame-forum`）和补丁站才是。本仓只是 **借用 NextMoe·未萌 的账号**：用它的 OIDC 做登录，顺便拿便宜的用户资料；**本站权限完全由 vnfm 自己的表和规则决定**。

论坛仓可以当「Nuxt + Fiber + KunUI + OIDC 怎么写」的参考实现，不要当成部署拓扑或权限模型的模板。

权威源：

- 枢纽代码与契约：`../nextmoe-infra`（只读；本仓 **不改 infra**）
- 产品仓写法参考：`../kun-galgame-forum`
- OIDC 公开协议：infra `docs/integration/oauth/`、`docs/auth/03-oidc-standardization-design.md`

---

## 0. 已锁定的关系

| # | 裁定 |
|---|------|
| 1 | vnfm 独立产品，不是 kungal / moyu 那种下游。不进 nextmoe 的 compose 网络，不共用它的 Postgres / Redis / MinIO。 |
| 2 | 只借用账号。登录、注册、改密、资料主档在 `account.nextmoe.com`。 |
| 3 | **本站授权 100% 由 vnfm 控制。** nextmoe 的 `roles` / `site_roles` / `ren` / `admin` 最多当展示或调试信息，**不得**用来放行本站任何管理操作。 |
| 4 | 向 nextmoe 只借两件事：**OIDC**（必须）和 **catalog**（按需）。图床本站自建（Cloudflare R2）。不用 image、artifact、community、trust、AI、萌萌点账本、开放 API 开发者平台。 |
| 5 | **禁止与 infra 做 S2S。** 不调 `/users/batch`、不调萌萌点 s2s、不用 OAuth Client Basic 打业务接口、不用 Docker 服务名 `oauth:9277`。只走公网（或本机）上的 OIDC / catalog **用户面或公开面**。 |
| 6 | 本仓不改 infra。要登记 OAuth client、回调域、scope，走账户管理台，不提 infra PR 当默认路径。 |
| 7 | 部署形态未定，**有可能不用 Docker**。所有对接按独立站点写：配置里是 `https://account.nextmoe.com` 这类 URL，不写死 compose 服务名或共享 volume。 |

「S2S」在这里特指枢纽给第一方下游的 **服务账号通道**（Client Basic、内部 DNS、共享库、调账、批量用户资料）。标准 OIDC 作为 Relying Party 换票、带用户 access token 调 `/oauth/userinfo`、读 JWKS，**不算 S2S**，那是协议本身。

---

## 1. 和论坛差在哪

| | 论坛 kungal | 本仓 vnfm |
|--|-------------|-----------|
| 生态位置 | nextmoe 下游，compose 同一网络 | 外部站点，只借账号 |
| 调 infra | 服务名 + 部分 s2s | 仅公开 HTTPS（或开发时的 localhost 端口） |
| 用户资料 | `/users/batch` s2s + userinfo | **只有**当前用户的 `/oauth/userinfo` |
| 权限 | 会把 IdP `roles` 映射进站内能力 | 站内权限表自己说了算；IdP 角色不授权 |
| 图床 | nextmoe image service | 本站 R2，不接他们的图床 |
| catalog | 下游客户端，常带服务凭证 | 公开读或用户 Bearer，按需 |
| 部署 | Dokploy + 与 infra 同网 | 独立，Docker 非必须 |
| UI | KunUI | 同样用 KunUI（这是组件库，不是生态隶属） |

---

## 2. 借用范围：OIDC + 按需 catalog

### 2.1 OIDC（必须）

自建 OP，issuer 目标域 `https://account.nextmoe.com`。Authorization Code + PKCE S256。

本站作为 **Relying Party**：

```text
用户点登录
  → 前端生成 PKCE + state
  → 跳到 https://account.nextmoe.com 的 /oauth/authorize
     scope=openid profile email
  → 未登录则在账户中心登录/注册
  → 带回 code
  → 本站后端用 code（+ PKCE，confidential 时再加 client_secret）调 /oauth/token
  → 用 access_token 调 /oauth/userinfo，得到 id / sub / name / picture …
  → 本站自己发会话（cookie 或本站 session 存储），token 不要进浏览器
```

协议端点（裸 JSON，不是 `{code,message,data}` 信封）：`/oauth/authorize`、`/oauth/token`、`/oauth/userinfo`、`/oauth/revoke`、`/oauth/jwks`、`/.well-known/openid-configuration`。

用户身份主键用 userinfo 里的整数 `id`（与 `sub` 对应）。本站用户表以这个 id 为外键，**不要另发一套登录号**。这只是「用 IdP 的 subject」，不是加入他们的库。

注册跳账户中心即可，不要在地图站做注册表单。改密、换绑邮箱、改头像也回账户中心。

会话怎么存（本机 Redis、数据库、还是别的）由 vnfm 自己定，**不要**接 infra 的 Redis。

### 2.2 catalog（按需，非第一期阻塞）

作品 / 厂商 / 角色元数据。给神器榜、十二器、萌战提名当数据源，替代或并行原 Bangumi/VNDB 代理。

只调 catalog 的 **公开读或用户 Bearer 写**。不要用第一方 admin 面、claim-events feed、s2s `works/claim`。

### 2.3 明确不用

nextmoe **image service**、artifact、community、trust、ai、萌萌点 s2s 调账、开发者平台 MCP、infra 的 Postgres 多库、OpenSearch、与 kungal 共用的 session 前缀。

---

## 2a. 本站图床（R2，保持简单）

原站图片散落在 `data/avatars/`、`data/club_avatars/`、`uploads/`、`wiki/uploads/`、`Forum/uploads/`。重构后统一进 **本站对象存储**，后端只当签发者和元数据记录，字节放 Cloudflare R2。

不要抄 nextmoe image service。那套是给多站共用的：内容寻址 hash、libwebp/cgo、preset 变体、引用计数、独立 Postgres、GC worker。本站只有一个产品，不需要。

### 做

- 一个 R2 bucket，S3 兼容 API（endpoint / account id / access key 配在本仓环境变量）
- 公开读：自定义域或 r2.dev，业务表只存 **object key**（或拼好的公开 URL），不存 hash 当主键
- 上传：登录用户走本站 API；小文件可由 Fiber 代传到 R2，或签发短时 **presigned PUT** 让浏览器直传。二选一即可，先代传更少移动部分
- key 按用途分前缀，例如 `club/{id}/logo`、`event/{id}/{uuid}`、`wiki/{club_key}/{uuid}`、`user/{id}/upload/{uuid}`，覆盖写入或 uuid 都行，不要做跨用户内容去重
- 限制 MIME（jpeg/png/webp/gif）和大小；可选：服务端压一次再存，失败就拒收，不要做变体矩阵
- 删除：有权限的人删业务记录时顺手 `DeleteObject`；不做引用计数、不做延迟 GC
- PDF / 稿件 / 大文件若以后要存，另开前缀或另开 bucket，不要塞进「图床预处理」

### 不做

- 独立 image 微服务、独立图片库
- `{aa}/{bb}/{sha256}_avatar-100.webp` 那种寻址和 preset
- 跨站去重、refping、软删回收站
- 把用户头像主档写进 R2——账户中心的头像仍以 OIDC `picture` 为准；本站只缓存 URL。社团头像、活动图、Wiki 图才是本站对象

配置示例：

```text
R2_ACCOUNT_ID=...
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET=vnfm
R2_ENDPOINT=https://<accountid>.r2.cloudflarestorage.com
R2_PUBLIC_BASE=https://media.example.com
```

---

## 3. 本站权限（vnfm 自己管）

nextmoe 五角色（`creator` / `moderator` / `admin` / `ren`，普通用户 claim 为空数组）是 **账户中心和第一方站点** 的契约。vnfm 可以显示「此账号在未萌是管理员」，但：

- 放行本站超管台、审核、改会籍、发公告，只查 **vnfm 自己的角色/权限表**
- 原 `super_admin`、`representative`、`manager`、`member`、`external`、GalOnly 陪审，全部落在本站（Club 会籍或本站 RBAC）
- 不要实现 `ren ⊇ admin ⊇ moderator` 当本站规则
- 不要用 infra 的 `site_roles` 给地图站授权（那是下游站点域角色，本站不是下游）
- 封禁本站功能（踢出社团、禁止报名）是本站状态；账户中心封禁只会让对方登不进 OIDC，两套不要混成一张表

建议本站权限模型：

1. **身份**：OIDC subject → `users.id`
2. **本站角色**：例如站点管理员、陪审，存在 vnfm 表
3. **会籍角色**：`member` / `manager` / `representative` / `external`，存在 `club_memberships`
4. 代码检查本站权限串或会籍，**永不** `if roles.contains("admin")` 拿 IdP 角色开门

---

## 4. 仓形状（可参考论坛，不绑定其运维）

pnpm workspace + `apps/web`（Nuxt 4 + KunUI）+ `apps/api`（Fiber）仍然适合本仓，因为这是技术选型，不是生态隶属。

Go 模块五层（handler / service / repository / dto / model）、编号 SQL 迁移、启动不做 AutoMigrate，都可以照抄。

前端约定同样适用：KunUI 优先、禁止背景渐变、页面单一根节点、`cn` + 项目色板。KunUI 是发布在 npm 的组件库，跟是不是下游无关。

对接 URL 全部来自本仓环境变量，例如：

```text
OIDC_ISSUER=https://account.nextmoe.com
OIDC_CLIENT_ID=...
OIDC_CLIENT_SECRET=...          # 仅后端；若选 public client 则不持有 secret，只靠 PKCE
OIDC_REDIRECT_URI=https://<vnfm>/auth/callback
CATALOG_BASE=https://<catalog-public>          # 按需
R2_ENDPOINT=https://<accountid>.r2.cloudflarestorage.com
R2_BUCKET=vnfm
R2_PUBLIC_BASE=https://media.example.com
```

开发机可以指向 `http://127.0.0.1:9277`，那也是本机端口，不是 Docker 网络别名。生产按独立站点反代即可，不假设和 infra 同机、同 compose、同 Dokploy 项目。

---

## 5. 登录实现注意（和论坛的取舍）

论坛 BFF 把 OAuth token 放进 Redis、cookie 只带会话 id——这个 **模式** 可以留（token 不进浏览器）。但：

- cookie 名、Redis（若用）实例、key 前缀都是 vnfm 自己的
- 换票只打公开 `/oauth/token`，不要再写一份 `userclient` 去打 `/users/batch`
- 展示他人资料（论坛作者栏那种）不能靠 s2s 批量拉主档；用本站缓存的公开字段（登录时从 userinfo 写下的 name/avatar），或仅在对方访问本站时刷新
- client 用 confidential 还是 public：独立站点、后端可保管 secret 时 confidential 更稳；若部署环境不方便藏 secret，用 public + PKCE。两者都是 OIDC，都不是 S2S

---

## 6. KunUI

| 包 | 用途 |
|----|------|
| `@kungal/ui-nuxt` | Nuxt layer，`extends` |
| `@kungal/ui-vue` | 组件 |
| `@kungal/ui-core` | `cn`、色 token |
| `@kungal/ui-tokens` | 设计 token |
| `@kungal/editor-nuxt` | 富文本；Wiki 需要时再加 |

铁律：先用 KunUI；缺功能或 bug 报告给用户，不改 KunUI 源码；不用背景渐变。

---

## 7. 和 `00-modules.md` 的衔接

领域五模块不变。Platform 内核变成：

```text
account.nextmoe.com     登录 / 注册 / 资料主档     （OIDC，公网）
catalog（按需）          作品元数据                  （公开读或用户 Bearer）
        │
        ▼
vnfm Platform           本站会话 · 本站权限 · R2 图床 · 通知 · layout
        │
        ▼
Club / Ops / Wiki / Social / Program
```

原 01 里的本地注册、QQ/Discord、改密、绑邮箱不再在本站实现。用户中心只留本站业务：会籍、代表同好会、本站通知、本站主题/语言偏好。

---

## 8. 建议落地顺序（技术）

1. `apps/web` + `apps/api` 骨架，KunUI 接上
2. 在账户管理台登记 **外部** OAuth client（回调指向 vnfm 自己的域名），本仓只保存 client 配置
3. OIDC 登录：PKCE、token、userinfo、本站会话；**没有** s2s userclient
4. 本站 `users` + 权限表；IdP `roles` 不参与鉴权
5. Club 模块（地图 + 目录 + 会籍）；头像等上传走本站 R2
6. catalog 在神器榜 / 提名真正需要时再接

不要做：改 infra、接 nextmoe 图床、Docker 服务名、S2S、把 nextmoe `admin` 当成本站超管、自签一套登录 JWT、把 token 放 localStorage、为「和论坛一致」去接萌萌点或 trust。
