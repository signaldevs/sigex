const vars = [
    'DEV_VAR',
    'OVERRIDE_TEST',
    'CLI_VAR',
    'GCP_SECRET',
    'AWS_SECRET',
    'ROT13_SECRET',
    'NOT_SECRET'
];
// noinspection JSUnresolvedVariable
vars.forEach(v => console.log(`${v}: ${process.env[v]}`));
