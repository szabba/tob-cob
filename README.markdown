# Tears of Butterflies: Colors of Blood

[![Build status](https://cloud.drone.io/api/badges/szabba/tob-cob/status.svg)](https://cloud.drone.io/szabba/tob-cob)

[![Itch.io page](https://img.shields.io/badge/Itch.io-FA5C5C?style=for-the-badge&logo=itch.io&logoColor=white)](https://szabba.itch.io/tears-of-butterflies-colors-of-blood)

## Build

```bash
nix develop
go build
```

Prerequisites:

- nix (with flakes support)

## Run

```bash
./tob-cob
```

Prerequisites:

- steam-run

## Publish

```bash
echo "BUTLER_API_KEY=xxx" > SECRETS # Use actual key instead.
earthly --push --secret "$(cat SECRETS)" +deploy
```

Prerequisites:

- drone CLI
- docker
- itch.io API key
