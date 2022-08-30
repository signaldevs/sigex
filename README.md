# sigex

## Installation

Coming soon...

## Usage

Example running a node app:

```bash
sigex node app.js
```

Running a node app with a `.env` file:

```bash
sigex -f .env node app.js
```

Running a node app with multiple `.env` files and specific env vars

```bash
sigex -f config/.dev.env -f .env -e FOO=BAR node app.js
```

## Secrets

Secrets can be automatically resolved using supported secrets manager.

```
MY_TOKEN=secret://gcp/dev/my/secret
```

```
SIGEX_GCP_REGION=
SIGEX_GCP_PROJECTID=
```