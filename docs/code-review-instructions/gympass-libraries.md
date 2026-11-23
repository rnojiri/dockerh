---
applyTo: "**"
description: Gympass libraries - none are used; the library must stay dependency-light
---

# Gympass libraries

- This repo imports no Gympass package (only `testify` in tests). New code MUST NOT add a Gympass dependency or a
  heavy third-party one (no cgo, no confluent-kafka-go): consumers import this module into their tests.
- Need not covered by the standard library? Prefer a command through `docker exec` (see `Exec`).

## Known exceptions (as of 2026-11-23)

None.
