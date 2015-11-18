# tbd

A decentralised ci built on git.

## Status

Pre-Alpha. You probably don't want to use this yet.

## Why a new CI?

We're frustrated with

 * The exorbitant cost of paid solutions (*cough* bamboo)
 * Inflexible CI servers making it difficult to
  - customise your workflow (e.g. merge master into branch, then build - as the `rust` project does)
  - port build configuration between systems (or even a different server running the same software)
  - answer simple questions like `did commit <sha> 2 months ago pass or fail?`
  - answer simple questions like `What was the build configuration when <commit> failed the build?`
 * Workflows which encourage developers not to check the build status of their commits
 * Tools that *only* provide web interfaces
 * CIs claiming flexibility because you can upload a plugin
  - Plugins are separately versioned & frequently implemented in a different language
  - It's not clear which version of a plugin was running when a build ran

## How is tbd different?

Core ideas driving `tbd` which are different to traditional CI:

 * All CI configuration is maintained in the same repository as your code
 * Build results are stored in your git repository (not in the main tree)as git notes, alongside that commit's code

## Build storage

`tbd` stores build results in a git ref (specifically `.git/refs/tbd/all-build-results` by default, controlled by git config tbd.ref or TBD_REF environment).

tbd commits contain a directory for each commit/worktree which has ever been built.
Because of how git stores files, this requires very little storage.
Artifact storage is opt-in; your build configuration can include an array of artifacts which will be passed to `git add -f` and written to the build output.

Example directory structure:

```
<source worktree sha>
  <spec:coverage>
    <build timestamp and host>
      STDOUT
      STDERR
      ETC
      artifacts/
        application files
        build artifacts
<source commit sha>
  <same thing>
```

When you run a build, tbd creates a new commit on the `.git/refs/tbd/all-build-results` ref

To check a file that was modified by the build process:
`git show tbd/all-build-results:<commit>/<metadata-hash>/<target>/WORKTREE/<artifact>`
`git tbd show <ref-like> <target> [--build-number <metadata-hash, defaults to most recent>] -- artifacts/coverage.html`
`git tbd show <ref-like> <target> [--build-number <metadata-hash, defaults to most recent>] <git show args>`
`git tbd unpack <ref-like> <target> [--build-number <metadata-hash, defaults to most recent>] <dir_to_unpack_into>`

### Advantages
 * Builds can be viewed & data extracted without tbd tools installed
 * no changes required to repository config
 * git notes for each commit can link directly to the build

### Disadvantages
No garbage collection for builds of a tree that was never pushed
 * not a problem unless devs are doing local builds and pushing the branch
 * probably not a problem since even on a large repo there's very little change
 * We could write a tool to strip out builds which aren't in the history (it'd re-write the branch)

It's easy commit a large binary as part of the post-build worktree and hard to undo.
 * We'll need to make it obvious, when writing a build to git, that artifacts have been saved (name/size)

A developer could carelessly check out `tbd/all-build-results`, which would cause a *lot* of files to be written to their machine.
However, they would have to really be poking around since it's not in the `git branch` or `git tags` list.

## TODO

### Sort out a way to merge build result refs without checking them out

Currently users have to checkout the build result ref to merge their builds.
There's no good reason for this as we never want to diff/merge individual files, only trees (and we anticipate conflicts being extremely rare - will still need to figure out how to resolve them).

Options investigated:
 * Storing in a `git notes` ref
  - Only supports blobs which is a non-starter
 * Storing in a branch
  - :+1: Sets up remote syncs easily by default
  - :-1: Easy to accidentally checkout
  - :+1: Still need to checkout to merge using `git pull`
 * Storing in a ref
  - need to write a custom mergetool (but it sounds like that was going to happen anyways)
  - eg `tbd-ff-merge refs/tbd/all-build-results refs/tbd/remotes/origin`

