# Release Guide

Releasing changes to the deps.cloud ecosystem is a rather easy task.
However, it requires some expertise with the system and tool chain should anything go wrong.
This document outlines the general approach used to release new versions of the project.

## Prerelease checklist

Before releasing the project, here are a few things to check before performing a release.

- [ ] Ensure all builds are passing on the main branch
  - See [Branch Checks in README.md](README.md#branch-checks)
- [ ] Permissions to push to `main`
- [ ] Permissions to push a tag

## Performing a release

To release a new version of the deps.cloud source, you simply need to push a tag.

1. Create a new version: `npm version <patch|minor|major>`
2. Push new tag: `git push --follow-tags`
3. While the builds progress, take some time to write up release notes.
   Once the tag workflows complete, you should be free to update the release.

### Writing release notes

I've found that many people ignore automated release notes.
They are often terse and difficult to read.
To alleviate some of this, I've tried to work with the following format.
This includes a small hand written summary and enumerates changes from the last version.

```markdown
## Changelog

### Summary

- Bulleted list of major fixes and feature development
- Hand written by release manager, can be highlights from commits

### Commits

- Because we squash commits going into main, we can be smart about the messages
- This section is computer generated and can be obtained using git-log

   $ start=$(git tag -l | tail -n 2 | head -n 1)
   $ end=$(git tag -l | tail -n 1)
   $ git log --format="%h: %s" ${start}...${end} | pbcopy
```

For an example, see our [v0.2.28 release notes](https://github.com/depscloud/depscloud/releases/tag/v0.2.28).

## Post-release checklist

Once a release is complete, there are some post-release tasks to complete. 

- [ ] Update deployment configuration to use newer images. ([depscloud/deploy])
- [ ] Write a blog post to advertise the new release. ([depscloud/deps.cloud])
- [ ] Write a blog post to advertise any new features. ([depscloud/deps.cloud])

[depscloud/deploy]: https://github.com/depscloud/deploy
[depscloud/deps.cloud]: https://github.com/depscloud/deps.cloud

## Becoming a release manager

In order to become a release manager, you will need to:

- Express interest
- Be an active contributor
- Participate in several releases
