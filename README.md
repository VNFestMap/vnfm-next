# vnfm-next

VNFest 同好会地图重写。Nuxt 4 前端 + Go Fiber API。账号借用 NextMoe OIDC，本站权限自管。当前范围见 [docs/05-v1-scope.md](docs/05-v1-scope.md)。

## 本地运行

需要 Node 24+（Corepack）、Go 1.26+、pnpm。API 的 `/healthz` 不依赖 Postgres；跑迁移才需要数据库。

```bash
corepack enable
pnpm install

cp apps/api/.env.example apps/api/.env
cp apps/web/.env.example apps/web/.env

# Web http://127.0.0.1:3710  API http://127.0.0.1:3711
pnpm dev
```

分开启动：`pnpm dev:web` / `pnpm dev:api`。迁移：`pnpm migrate`（需在 `apps/api/.env` 填写 `DATABASE_URL`）。
