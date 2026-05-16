# LazyCosmos

LazyCosmos is a keyboard-first terminal UI for Cosmos SDK application development workflows.

It is intended to feel like `lazygit` for local appchain development: profiles, node lifecycle, keys, accounts, transactions, queries, logs, genesis, config, and state diffs from one terminal interface.

## Requirements

- Go 1.26.3 or newer

## Run From Source

```bash
go run ./cmd/lazycosmos
```

## Common Commands

```bash
just run
just test
just build
just lint
just fix
```

## What LazyCosmos Is

- A development cockpit for Cosmos SDK chains.
- A transparent wrapper around chain binaries, CometBFT RPC, REST, gRPC, config files, and local development workflows.

## What LazyCosmos Is Not

- Not a production wallet.
- Not a Mintscan replacement.
- Not a validator monitoring dashboard.
- Not a private-key export tool.
