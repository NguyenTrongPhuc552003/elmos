# Documentation site

ELMOS's documentation is published on the project site,
[haven-docs](https://github.com/NguyenTrongPhuc552003/haven-docs)
(Astro Starlight). Its `scripts/sync.mjs` copies `README.md` (the ELMOS
index), `docs/user` (the guide) and `docs/developer` (this section) at
build time and generates the CLI reference with `elmos gendocs`. The
repositories stay the source of their pages; there is no second site.

## Retiring MkDocs (design and expectations, 2026-11-13)

Owner-approved 2026-11-11 (positioning note, item 5: "once the Taskfile no
longer needs it"). Written before the change.

**Today.** `mkdocs.yml` and `docs/index.md` (its home page) are used only
by the Taskfile tasks `docs`, `docs:install`, `docs:clean` and
`docs:serve` (MkDocs Material, mkdocstrings, mkdocs_puml, installed with
pip). No workflow builds them (the Pages workflow was removed earlier) and
the `site_url` in `mkdocs.yml`, `elmos.vercel.app`, answers HTTP 404.

**Change.**

- Remove `mkdocs.yml`, `docs/index.md` and the four tasks.
- New tasks `docs` and `docs:serve` build or serve haven-docs from a
  checkout next to this one (`../haven-docs`, with `../haven` and
  `../kdal`), which is how haven-docs expects a local workspace.

**What readers lose.** Rendered PlantUML: the `puml` blocks of
[Architecture Diagrams](diagrams.md) and `docs/diagrams/*.puml` stay as
source (Starlight shows them as code); and the local MkDocs Material
preview. Nothing that is published: no MkDocs site is online.

| Check | Expected |
| - | - |
| D1: references | after the change, `git grep -i mkdocs` matches only this page |
| D2: tasks | `task --list` shows `docs` and `docs:serve` (new descriptions); `docs:install` and `docs:clean` are gone |
| D3: the remaining build | haven-docs, copied outside the owner's checkout and pointed at this repository with `SRC_ROOT`: `sync.mjs` and `astro build` succeed and every page of `docs/user` and `docs/developer` is in the built site; its link check passes |
| D4: ELMOS CI | green |
