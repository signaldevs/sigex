# sigex

## Installation

Coming soon...

## Usage

```bash
sigex is a process runner/executor with support for multiple .env file
configuration as well as automatic retrieval of secrets from
supported secrets manager platforms.

Usage:
  sigex [flags] command

Flags:
      --debug                    debug the resolved environment variables
  -f, --env-file strings         specify one or more .env files to use
  -e, --env-var stringToString   specify one or more environment variables to use (ex: -e FOO=bar) (default [])
  -h, --help                     help for sigex
      --skip-secrets             skip the automatic resolution of secret values
```

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

- Google Cloud Secrets Manager
- AWS Secrets Manager

### Google Secrets Manager

Token Format: `sigex-secret-{secret platform}://{Resource Id incl Version}`

```bash
# format: sigex-secret-gcp://{secret-resource-version-id}
MY_GCP_SECRET=sigex-secret-gcp://projects/00000000000/secrets/mysecret/versions/latest
```

### AWS Secrets Manager

Token Format: `sigex-secret-aws://{Resource Id}`

```bash
# format: sigex-secret-aws://{secret-resource-version-id}
MY_AWS_SECRET=sigex-secret-aws:///dev/sigex/test
```

### Rot13 Secrets

This is used for testing or very light obfuscation it provides zero real security.

Token Format: `sigex-secret-rot13://uryyb_jbeyq`

```bash
# format: sigex-secret-rot13://{rot13 encoded text}
MY_ROT13_SECRET=sigex-secret-rot13://uryyb_jbeyq
```

## Running the example

Check out the [example](examples/node) for a simple node.js program that
demonstrates how to use sigex to retrieve secrets from AWS Secrets Manager and
GCP Secret Manager.
