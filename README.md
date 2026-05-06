# GoSX CMS

GoSX CMS is the opinionated content layer for GoSX applications.

It is meant to feel like the batteries-included side of Django: site settings,
pages, posts, media, SEO, and opinionated block catalogs built on top of
`gosx-admin`.

Core GoSX stays small; this module carries reusable CMS patterns for apps that
want them. It can grow independently as the examples harden into packages.

Current package surface:

- `blocks`: CMS content block catalogs that can be rendered by admin block
  editors or public GoSX surfaces.

```sh
go get github.com/odvcencio/gosx-cms/blocks
```
