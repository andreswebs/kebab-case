# kebab-case

CLI that recursively renames files and directories to kebab-case, walking the tree concurrently with goroutines.

## Layout

Single `main` package at the repo root:

- [main.go](main.go) — entrypoint, arg parsing, top-level WaitGroup/error channel
- [format.go](format.go) — `Format(string) string`, the core string transform (regex pipeline)
- [utils.go](utils.go) — `ProcessFile` recursion, rename, dir traversal helpers
- `*_test.go` — table-driven tests, including a `t.TempDir`-based integration test for `ProcessFile`

## Commands

```sh
go build ./...
go test ./...
golangci-lint run
```

CI: [.github/workflows/golangci-lint.yml](.github/workflows/golangci-lint.yml), [.github/workflows/gosec.yml](.github/workflows/gosec.yml).

## Conventions worth knowing

- `Format` is the source of truth for the transform. When adding a character class or rule, extend the regex pipeline there and add a row to the table in [main_test.go](main_test.go).
- `ProcessFile` spawns a goroutine per child entry and shares one `*sync.WaitGroup` + buffered error channel passed down the recursion. Preserve that shape — don't introduce a worker pool or per-call channel without a reason.
- Error channel in [main.go](main.go) is buffered to size 1; `run()` drains it after `wg.Wait()`. If you change buffering, also check the drain loop still terminates.
