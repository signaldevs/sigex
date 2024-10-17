# node.js example app

This example is a simple node.js program that demonstrates how to use sigex to
retrieve secrets from AWS Secrets Manager and GCP Secret Manager. Either of
these secrets managers can be used or both.

## Running the example

1. Modify the `config/.dev.env` file (see comments in the file)

2. Build and install sigex: `go install`

3. Run the example:

```bash
$ cd examples/node
$ sigex -f config/.dev.env -f .env node app.js
```
