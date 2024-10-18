# sigex

[![GitHub License](https://img.shields.io/github/license/signaldevs/sigex?style=flat)](https://github.com/signaldevs/sigex/blob/master/LICENSE) ![GitHub Release](https://img.shields.io/github/v/release/signaldevs/sigex)

`sigex` is a process runner/executor with support for multiple `.env` file
configuration with automatic retrieval of secrets from supported secrets
manager platforms.

You can run any process command with `sigex`.

## Installation

### MacOS

With [Homebrew](https://brew.sh/):

```bash
brew tap signaldevs/tap
brew install sigex
```

### Windows

Coming soon...

## Usage

```bash
Usage:
  sigex [flags] command

Flags:
      --debug                    debug the resolved environment variables
  -f, --env-file strings         specify one or more .env files to use
  -e, --env-var stringToString   specify one or more environment variables to use (ex: -e FOO=bar) (default [])
  -h, --help                     help for sigex
      --skip-secrets             skip the automatic resolution of secret values
```

Example running a python app:

```bash
sigex python test.py
```

Running a node app with a `.env` file:

```bash
sigex -f .env node app.js
```

Running a node app with multiple `.env` files and specific env vars

```bash
sigex -f config/.dev.env -f .env -e FOO=BAR node app.js
```

## Secret Token Format

`sigex` resolves environment variables from common secret managers. Instead of hard coding values in your env vars, you can use the `sigex-secret-{secret_manager}://` prefix to resolve values from supported secret managers.

Supported secret managers:

- [Google Cloud Secrets Manager](#google-cloud-secrets-manager)
- [AWS Secrets Manager](#aws-secrets-manager)

Example:

```bash
SECRET_GCP_KEY=sigex-secret-gcp://projects/00000000000/secrets/mysecret/versions/latest
SECRET_AWS_KEY=sigex-secret-aws://path/to/secret
```

## Environment Files (`.env` files)

`sigex` supports using one or more `.env` files. The format for the variables in the files should be in `key=value` format like so:

```text
SECRET_KEY=sigex-secret-gcp://projects/00000000000/secrets/mysecret/versions/latest
MODE=FOO
API_URL=http://www.signaladvisors.com
```

### Google Cloud Secrets Manager

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

### Rot13 Secrets [DEPRECATED]

This is used for testing or very light obfuscation it provides zero real security.

Token Format: `sigex-secret-rot13://uryyb_jbeyq`

```bash
# format: sigex-secret-rot13://{rot13 encoded text}
MY_ROT13_SECRET=sigex-secret-rot13://uryyb_jbeyq
```

## Testing Secrets Resolution

You can run `sigex` with the `--debug` flag to see the resolved environment variables and their values.

```bash
export SECRET_KEY=sigex-secret-gcp://projects/00000000000/secrets/mysecret/versions/latest
sigex --debug | grep SECRET_KEY
```

## Running the Example

Check out the [example](examples/node) for a simple node.js program that
demonstrates how to use sigex to retrieve secrets from AWS Secrets Manager and
GCP Secret Manager.
