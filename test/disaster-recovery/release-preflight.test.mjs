import assert from 'node:assert/strict';
import {mkdtempSync, writeFileSync, readFileSync, rmSync} from 'node:fs';
import {tmpdir} from 'node:os';
import {join} from 'node:path';
import {spawnSync} from 'node:child_process';
import test from 'node:test';

const workflow = readFileSync(new URL('../../.github/workflows/release-cli.yml', import.meta.url), 'utf8');
const validation = workflow.split('\n  validate:\n')[1].split('\n  build:\n')[0];
const script = validation.split('        run: |\n')[1].split('\n').map(line => line.replace(/^          /, '')).join('\n');
const revision = '1'.repeat(40);
const fakeGh = `#!${process.execPath}
const fs = require('node:fs');
const args = process.argv.slice(2);
fs.appendFileSync(process.env.STORM_TEST_CALLS, JSON.stringify(args) + '\\n');
const fixture = JSON.parse(process.env.STORM_TEST_FIXTURE);
if (args[0] === '--version') console.log('gh version ' + (fixture.ghVersion || '2.103.0'));
else if (args.includes('--help')) process.exit(0);
else if (args[0] !== 'api') process.exit(90);
else if (args[1].includes('immutable-releases')) {
  console.error('Resource not accessible by integration (HTTP 403)');
  process.exit(1);
} else if (args[1].includes('/actions/runs?')) {
  if (fixture.denied) process.exit(1);
  console.log(JSON.stringify({workflow_runs: ['go-cli.yml', 'documentation.yml'].map(path => ({
    path: '.github/workflows/' + path, head_sha: fixture.foreign ? '2'.repeat(40) : process.env.GITHUB_SHA,
    head_branch: 'master', run_number: 1, status: fixture.pending ? 'in_progress' : 'completed', conclusion: 'success'
  }))}));
} else if (args[1].includes('/git/matching-refs/')) {
  console.log(JSON.stringify(fixture.usedTag ? [{ref: 'refs/tags/' + process.env.VERSION}] : []));
} else if (args.includes('--paginate')) {
  if (fixture.usedRelease) console.log(process.env.VERSION);
} else process.exit(91);
`;

for (const [name, fixture, allowed] of [
  ['repository-scoped token succeeds while administrative settings are inaccessible', {}, true],
  ['foreign repository is denied', {repo: 'other/project'}, false],
  ['development dispatch is denied', {ref: 'refs/heads/develop'}, false],
  ['outdated verification client is denied', {ghVersion: '2.80.0'}, false],
  ['pending source checks are denied', {pending: true}, false],
  ['failed metadata access is denied', {denied: true}, false],
  ['foreign source revision is denied', {foreign: true}, false],
  ['existing tag is denied', {usedTag: true}, false],
  ['existing release is denied', {usedRelease: true}, false],
  ['shell text in version is denied', {version: 'v0.1.0-alpha.1; exit 0'}, false],
]) {
  test(`release preflight: ${name}`, () => {
    const directory = mkdtempSync(join(tmpdir(), 'storm-release-preflight-'));
    try {
      writeFileSync(join(directory, 'gh'), fakeGh, {mode: 0o700});
      const calls = join(directory, 'calls.jsonl');
      const result = spawnSync('/bin/bash', ['-c', script], {
        cwd: directory,
        encoding: 'utf8',
        timeout: 10000,
        env: {
          ...process.env,
          PATH: `${directory}:${process.env.PATH}`,
          GITHUB_REPOSITORY: fixture.repo || 'Storm-Harbor/storm-harbor-control-plane',
          GITHUB_REF: fixture.ref || 'refs/heads/master',
          GITHUB_SHA: revision,
          VERSION: fixture.version || 'v0.1.0-alpha.1',
          GH_TOKEN: 'synthetic-no-network-token',
          STORM_TEST_FIXTURE: JSON.stringify(fixture),
          STORM_TEST_CALLS: calls,
        },
      });
      assert.ifError(result.error);
      assert.equal(result.status === 0, allowed, result.stderr);
      if (allowed) {
        const requests = readFileSync(calls, 'utf8');
        assert.ok(requests.includes('/actions/runs?'));
        assert.ok(requests.includes('/git/matching-refs/'));
        assert.ok(requests.includes('--paginate'));
        assert.ok(!requests.includes('immutable-releases'));
      }
    } finally {
      rmSync(directory, {recursive: true, force: true});
    }
  });
}
