# Tears of Butterflies: Colors of Blood

[![Build status](https://cloud.drone.io/api/badges/szabba/tob-cob/status.svg)](https://cloud.drone.io/szabba/tob-cob)

[![Itch.io page](https://img.shields.io/badge/Itch.io-FA5C5C?style=for-the-badge&logo=itch.io&logoColor=white)](https://szabba.itch.io/tears-of-butterflies-colors-of-blood)

## Build

```bash
drone exec
```

Prerequisites:

- drone CLI
- docker

## Run

```bash
(cd ./out/linux-x64 && steam-run ./tob-cob)
```

Prerequisites:

- steam-run

## Build with deployment

```bash
echo "BUTLER_API_KEY=xxx" > SECRETS # Use actual key instead.
drone exec --secret-file SECRETS
```

Prerequisites:

- drone CLI
- docker
- itch.io API key
