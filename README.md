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

## Env Files

`sigex` supports using one or more `.env` files. The format for the variables in the files should be in `key=value` format like so:

```text
FOO=Bar
BIN=Baz
URL=http://www.signaladvisors.com
```

## Secrets

Secrets can be automatically resolved using a supported secrets manager. In your environment variables, insert a token for your secret and `sigex` will request it for you.

Current secrets managers supported:

* Google Cloud Secrets Manager
* _[coming soon] AWS Secrets Manager_

### Google Secrets Manager

Token Format: `sigex-secret-gcp://{Resource Id incl Version}`

```bash
# format: sigex-secret-gcp://{secret-resource-version-id}
MY_GCP_SECRET=sigex-secret-gcp://projects/00000000000/secrets/mysecret/versions/latest
```
