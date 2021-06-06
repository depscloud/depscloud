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

Most release notes are auto-generated today.
When we introduce a major feature, we write a small markdown file to provide a summary.
Then at build time, we use a combination of hand written feature announcements and auto-detected fixes to generate notes.
To get a better idea for how this works, see the `./scripts/gen-changelog.sh` script.
One task that is currently manual is tracking down the GitHub handles for contributors to a release. 

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
