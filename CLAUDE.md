# vnfm-next — AI Agent Project Guide

VNFest club map rewrite. `apps/api` = Go Fiber + GORM + Postgres (planned). `apps/web` = Nuxt 4 + KunUI (planned). This repo is an **independent product** that borrows NextMoe·未萌 **OIDC only** (catalog later, not in v1). It is **not** a nextmoe-infra downstream.

**Current build scope:** [docs/05-v1-scope.md](docs/05-v1-scope.md) — users, club create/join, events, notifications, China/Japan maps. Everything else is frozen.

## 铁律 (Iron Rules — non-negotiable; these override every other guideline in this file)

1. **No background gradients in any UI, ever.** Never use gradient backgrounds (`bg-gradient-*`, `from-*/via-*/to-*`, `linear-gradient()`, `radial-gradient()`, `conic-gradient()`, etc.); use solid colors from the project's palette.
2. **Prefer KunUI components; do not modify KunUI itself.** When adding or changing frontend UI, reach for a KunUI component (`@kungal/ui-*`) first — do not hand-roll a native/custom component unless there is genuinely no KunUI equivalent. If KunUI appears to have a bug or is missing a feature, **do not edit KunUI's code** (shared upstream library) — report it to the user and let them decide.
3. **Do not implement frozen features.** If it is not in the “做” section of `docs/05-v1-scope.md`, do not build it, scaffold empty modules for it, or “just add a table for later.” Unfreeze by editing that doc first.
4. **Do not treat this repo as a nextmoe downstream.** Do not modify `../nextmoe-infra`. Do not call infra S2S (`/users/batch`, Client Basic, Docker DNS `oauth:9277`, moemoepoint, trust, image service). OIDC authorization-code + PKCE + `/oauth/userinfo` on public URLs is allowed. Site authorization is 100% vnfm tables; IdP `roles` never grant access here. Images go to this site's R2, not nextmoe image.

These UI rules (KunUI-first, project palette only, no gradients) **override any global or user-level design skill** (e.g. a generic `frontend-design` skill). When they conflict, this file wins.

## Core Engineering Principles

> Shared baseline with KUN Galgame repositories. Defaults, not dogma — apply judgment.

1. All commit messages must be written entirely in English.
2. Comments are governed by the **Comments** section below — the default is none, and what survives is written in English.
3. Keep each source file under ~500 lines where practical; once a file grows past ~300 lines, consider splitting it (a guideline, not a hard rule).
4. Write every frontend function as an arrow function; compose/merge class names with `cn` wherever practical.
5. Deliberately balance elegant modularity against necessary duplication — choose per case instead of always favoring either.
6. Constantly verify that frontend and backend agree on the data: field shapes and response formats must match what each side expects.
7. After every change, watch for unintended side effects elsewhere.
8. If a change requires running a migration, tell the user explicitly at the end — which command, and against which database.
9. Always seek the most modern, elegant solution that fits the project's current state; consult the latest official docs and resources online when useful.
10. Never let the pursuit of elegance or modularity make the code complex or hard to follow, and don't write over-defensive code.
11. A Nuxt page — and any component used as a page/route root — must have a **single real root element**: never `display: contents` and never a leading comment / whitespace / sibling at the template root. Keep explanatory comments *inside* the root element.
12. Reserve the scrollbar gutter globally — `html { scrollbar-gutter: stable }`, with an `overflow-y: scroll` `@supports` fallback — so the document width is constant across routes.
13. **One task = one session, and every path has exactly one writer.** Parallel work is allowed only when the user assigns non-overlapping writable paths. Never rewrite shared Git state a peer may be standing on: on a shared checkout, no branch switch, reset, rebase, merge, cherry-pick, clean, stash or prune. Commit with explicit paths (`git commit -- <paths>`), never `add -A`.
14. DB-backed tests must use an explicit dedicated DSN from the session environment (`TEST_DATABASE_DSN` if set). Never discover or fall back to a DSN from `.env`, never print a password, never run against a live database. Unique ports for concurrent services; never stop a process whose owner is unknown.

## Comments

**Default: none.** Code that can be understood by reading it gets no comment. Most code is that code.

**A comment is earned by a mistake that already happened, not by one you predict.** Do not comment while writing. Comment when something went wrong there: an agent or a person got it wrong, a review caught it, a test went red, production broke. The comment records the wrong conclusion that was actually reached. If you cannot name the incident, there is no comment to write.

Standing exceptions:

- **Migrations** — `apps/api/migrations/**` (or equivalent). A migration is history. Say what it changes and why, including existing rows.
- **A constraint that is true but invisible from this file**: a version floor, an upstream bug, a required ordering.
- **A completeness assertion over a hand-maintained list** that a growing schema will silently outgrow.

Write the conclusion, not the mechanism. Quote real system output verbatim when reproducing a symptom — such a quote may keep its original language.

Never write: restatements of the code, section banners, `TODO` without an owner, or doc comments that only echo the identifier. If a comment explains what a name means, rename the thing and delete the comment.

**Not comments, never removed:** `//go:embed`, `//go:build`, `//go:generate`, `//nolint`, `// Code generated … DO NOT EDIT`, `eslint-disable`, `@ts-expect-error`, `prettier-ignore`, and anything inside a generated file.

English, and short. When in doubt, delete it.

Applies to `.go`, `.ts`, and `.vue`. It does **not** apply to config and onboarding files (`.env.example`, compose, CI), where the comment is often the only thing a person reads before running anything.

## Architecture (do not reinvent)

Read before coding:

| Doc | Use |
|-----|-----|
| [docs/05-v1-scope.md](docs/05-v1-scope.md) | What to build now; frozen list |
| [docs/00-modules.md](docs/00-modules.md) | Domain boundaries (Club / Ops / …) |
| [docs/04-nextmoe-architecture.md](docs/04-nextmoe-architecture.md) | OIDC borrow, no S2S, R2, local authz |
| [docs/01-core-platform.md](docs/01-core-platform.md) / [02](docs/02-community-operations.md) / [03](docs/03-contests-and-satellites.md) | Original-site inventory only |

`../kun-galgame-forum` is a **code-shape** reference (Nuxt pages, Fiber layers, OIDC BFF cookie). Do not copy its downstream topology, `/users/batch` client, image hashes, or IdP-role gates.

Planned layout:

```text
apps/web/     Nuxt 4, extends @kungal/ui-nuxt
apps/api/     Fiber; cmd/server + internal/{module}/{handler,service,repository,dto,model}
```

- Club identity is always `(club_id, country)` — never a global autoincrement that collapses China/Japan.
- Do not AutoMigrate on process start. Schema changes are numbered SQL migrations.
- User PK = OIDC `id`. Cache name/avatar from userinfo for display; account-center profile remains source of truth.
- Site permissions: club membership roles + vnfm tables. Never `roles.contains("admin")` on the IdP claim.
- Media: Cloudflare R2, object key in our tables. No content-addressed image service, no variant matrix, no refcount GC.
- Deployment is undecided and **may not be Docker**. Config uses public URLs (`https://account.nextmoe.com`), not compose service names.

## Frontend conventions

### UI

- Use `@kungal/ui-nuxt` / `@kungal/ui-vue` / `@kungal/ui-core` (`cn`). There is no in-repo `components/kun/` to edit.
- Color tokens from the KunUI palette (`text-foreground`, `text-default-500`, `border-default-200`, `primary` / `success` / `danger` / `warning` / `default` / `secondary` / `info` with 50–950). Do not use Tailwind built-in `gray` / `indigo` / `blue`. Tokens follow light/dark; no `dark:` prefix for those colors.

### Pages

- `pages/` is routes only: `definePageMeta` + one container component.
- Feature UI lives under `components/{feature}/`. Do not repeat the directory prefix in the filename (`components/users/Container.vue` → `UsersContainer`, not `UsersContainer.vue` inside `users/`).
- Constants: `app/constants/`. Types: `shared/types/` (Nuxt auto-import).
- Arrow functions only in frontend application code.

## Backend conventions

Each in-scope module (`user`/`auth`, `club`, `event`, `notification`, `media`) follows:

```text
internal/{module}/
  model/ dto/ repository/ service/ handler/
```

Handlers: parse → `MustGetUser` if needed → service → `response.OK` / `response.Error`. Services return a typed app error, not a raw `error` to the client. Repositories are GORM (or SQL) only — no HTTP, no OIDC.

After any schema change, end the task with: whether a production migration must run, which command, which database. Skipping this causes silent zero values on missing columns.

## What “frozen” means in the tree

Do not add routes, Vue pages, or migrations for: Jiangsu submap, Asia basemap, star map / star union, club comments/recommendations/moe king, publications, project hub, GalOnly, wiki, forum, square, vote/moe/twelve, recognition, quiz, bots, simulators, announcements, local password auth, nextmoe image/catalog clients.

If a user request is frozen, say so and point at `docs/05-v1-scope.md` instead of implementing it.
