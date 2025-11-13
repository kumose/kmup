import {sanitizeRepoName, substituteRepoOpenWithUrl} from './repo-common.ts';

test('substituteRepoOpenWithUrl', () => {
  // For example: "x-github-client://openRepo/https://github.com/go-kmup/kmup"
  expect(substituteRepoOpenWithUrl('proto://a/{url}', 'https://kmup')).toEqual('proto://a/https://kmup');
  expect(substituteRepoOpenWithUrl('proto://a?link={url}', 'https://kmup')).toEqual('proto://a?link=https%3A%2F%2Fkmup');
});

test('sanitizeRepoName', () => {
  expect(sanitizeRepoName(' a b ')).toEqual('a-b');
  expect(sanitizeRepoName('a-b_c.git ')).toEqual('a-b_c');
  expect(sanitizeRepoName('/x.git/')).toEqual('-x.git-');
  expect(sanitizeRepoName('.profile')).toEqual('.profile');
  expect(sanitizeRepoName('.profile.')).toEqual('.profile');
  expect(sanitizeRepoName('.pro..file')).toEqual('.pro.file');

  expect(sanitizeRepoName('foo.rss.atom.git.wiki')).toEqual('foo');

  expect(sanitizeRepoName('.')).toEqual('');
  expect(sanitizeRepoName('..')).toEqual('');
  expect(sanitizeRepoName('-')).toEqual('');
});
