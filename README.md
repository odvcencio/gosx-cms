# GoSX CMS

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
- `flows`: reusable server-action-backed flow definitions and Studio panels
  for contact, scheduling, enrollment, newsletter, checkout handoff, and
  app-owned handlers.
- `studio`: reusable three-pane authoring shell model and server-rendered
  structure for canvas, preview, panels, and actions.

```sh
go get github.com/odvcencio/gosx-cms
```
