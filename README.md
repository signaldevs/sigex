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

* Google Cloud Secrets Manager
* _[coming soon] AWS Secrets Manager_

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


# Testing:

## Testing AWS Secrets:

Admin access may be needed for secrets access, or some lesser access, coordinate with the administrator.

1. Goto Okta 
2. Click AWS 
3. Expand the account you want to test 
4. Click the "Command line or programmatic access" link
5. Copy the export command from step 1 and run those in the console you will be testing from
6. Also run this in the same console: `export AWS_REGION="us-east-2"`
7. Build sigex (from root of repo run: `go build`)
8. Go into the example app directory: `cd examples/node`
9. Run the example app: `../../sigex -f config/.dev.env -f .env node app.js`
10. Observe the output it should show `AWS_SECRET: "dev secret for sigex"` along with secrets from other sources.

