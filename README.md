# GoSX CMS — frozen, folded into gosx-studio

**This repository is frozen.** On 2026-07-06 all 13 packages were folded into
[`m31labs.dev/gosx-studio`](https://github.com/M31-Labs/gosx-studio) as the
`cms/` subtree, copied from this repo's `main` at commit `81dc8b0`. Full
provenance (path mapping, source SHA, copy method) lives in
`cms/PROVENANCE.md` in the gosx-studio repo.

No further development happens here. The code below is left intact and
`v0.2.1` (final tag) keeps resolving for existing consumers, but new work
should target the folded paths in gosx-studio.

**Import path mapping** (`m31labs.dev/gosx-cms/X` → `m31labs.dev/gosx-studio/cms/X`):

| Old import path | New import path |
|---|---|
| `m31labs.dev/gosx-cms/blocks` | `m31labs.dev/gosx-studio/cms/blocks` |
| `m31labs.dev/gosx-cms/content` | `m31labs.dev/gosx-studio/cms/content` |
| `m31labs.dev/gosx-cms/flows` | `m31labs.dev/gosx-studio/cms/flows` |
| `m31labs.dev/gosx-cms/lifecycle` | `m31labs.dev/gosx-studio/cms/lifecycle` |
| `m31labs.dev/gosx-cms/lifecycle/sqlstore` | `m31labs.dev/gosx-studio/cms/lifecycle/sqlstore` |
| `m31labs.dev/gosx-cms/media` | `m31labs.dev/gosx-studio/cms/media` |
| `m31labs.dev/gosx-cms/render` | `m31labs.dev/gosx-studio/cms/render` |
| `m31labs.dev/gosx-cms/store` | `m31labs.dev/gosx-studio/cms/store` |
| `m31labs.dev/gosx-cms/store/file` | `m31labs.dev/gosx-studio/cms/store/file` |
| `m31labs.dev/gosx-cms/store/memory` | `m31labs.dev/gosx-studio/cms/store/memory` |
| `m31labs.dev/gosx-cms/style` | `m31labs.dev/gosx-studio/cms/style` |
| `m31labs.dev/gosx-cms/studio` | `m31labs.dev/gosx-studio/cms/studio` |
| `m31labs.dev/gosx-cms/studio/collab` | `m31labs.dev/gosx-studio/cms/studio/collab` |

Migrate by a mechanical prefix sed: `m31labs.dev/gosx-cms/` → `m31labs.dev/gosx-studio/cms/`.

---

## GoSX CMS (historical description, frozen at v0.2.1)

GoSX CMS is the opinionated content layer for GoSX applications.

## Agent Skill

Agents helping someone use GoSX CMS should read the GoSX ecosystem skill: [using-gosx-ecosystem](https://github.com/odvcencio/m31labs-skills/blob/main/skills/using-gosx-ecosystem/SKILL.md).

It is meant to feel like the batteries-included side of Django: site settings,
pages, posts, media, SEO, and opinionated block catalogs built on top of
`gosx-admin`.

Core GoSX stays small; this module carries reusable CMS patterns for apps that
want them. It can grow independently as the examples harden into packages.

Current package surface:

- `blocks`: CMS content block catalogs that can be rendered by admin block
  editors or public GoSX surfaces.
- `content`: compatibility parsing and view-model helpers for structured body
  documents and the lightweight legacy block syntax.
- `lifecycle`: revision, draft, publish, preview, and rollback primitives for
  CMS stores.
- `media`: reusable asset, variant, focal point, usage, picker, and media-line
  primitives.
- `render`: generic content block rendering with hooks for app-owned product,
  flow, and custom block output.
- `store`: generic CMS store contracts and helpers for pages, posts, site
  settings, draft/publish state, and revisions.
- `store/memory`: turnkey in-memory implementation of the generic CMS store
  contracts for demos, tests, previews, and early site builds.
- `flows`: reusable server-action-backed flow definitions and Studio panels
  for contact, scheduling, enrollment, newsletter, checkout handoff, and
  app-owned handlers. The package also includes persisted authoring primitives
  for flow documents, drafts, publications, runtime instances, and lifecycle
  revision snapshots plus a memory store for demos and early Studio sites.
- `studio`: reusable three-pane authoring shell model and server-rendered
  structure for canvas, preview, panels, and actions.

```sh
go get github.com/odvcencio/gosx-cms
```

## Flow Documents

`flows.Document` is the persisted author-authored counterpart to
`flows.Definition`. It keeps flow metadata, action handler references, and each
step's `blockstudio.Document` together so Studio can draft and publish generic
flows such as contact, purchase request, checkout handoff, newsletter,
appointment, schedule tour, and enrollment.

Use `flows.StandardDocuments` to seed those generic flow shapes with app-owned
handler refs, `flows.NormalizeDocument` to clean authored step block documents,
`flows.InstanceFromDocument` to build a runtime definition, and
`flows.PublishDraft` or `flows.NewDocumentRevision` to create lifecycle
snapshots. The package exposes store interfaces only; applications decide where
documents, drafts, publications, and revisions live.
