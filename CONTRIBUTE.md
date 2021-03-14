# Contributing to Shore

This repository contains only Shore core, which includes the command line interface and the currently supported `Renderers` & `Backend`.

## Definitions

* **Maintainers** - The group that sets the road map, code quality standards and priorities (may also be called `Shore/Core Team`)\
  The Maintainers will create Github Issues, answer questions and provide guidance to contributors.
* **Contributors** - Individuals or Groups that want to take part in the development and growth of the project.\
  Contributors are expected to communicate with `Maintainers` to get guidance and clarity during development.

## Contributing Fixes

It can be tempting to want to dive into an open source project and help _build the thing_ you believe you're missing.

It's a wonderful and helpful intention.

For new contributors we've labeled a few issues with `Good First Issue` as a nod to issues which will help get you familiar with Shore development, while also providing an onramp to the codebase itself.

If the fix you are working on appears to be larger scope fix than was initially thought of, get in touch with the `maintainers` through the Github-Issue!

## Proposing a Change

In order to be respectful of the time of community contributors, we aim to discuss potential changes in GitHub issues prior to implementation.

That will allow us to give design feedback up front and set expectations about the scope of the change, and, for larger changes, how best to approach the work such that the Maintainers can review it and merge it along with other concurrent work.

If the bug you wish to fix or enhancement you wish to implement isn't already covered by a GitHub issue that contains feedback from the Maintainers (Core) team, please do start a discussion in a new Github Issue or an existing one, as appropriate, before you invest significant development time.

If you mention your intent to implement the change described in your issue, the Shore team can, as best as possible, prioritize including implementation-related feedback in the subsequent discussion.

For large proposals that could entail a significant design phase, we wish to be up front with potential contributors that, unfortunately, we are unlikely to be able to give prompt feedback. We are still interested to hear about your use-cases so that we can consider ways to meet them as part of other larger projects.

Most changes will involve updates to the test suite, and changes to Shore's documentation.

The Shore team can advise on different testing strategies for specific scenarios, and may ask you to revise the specific phrasing of your proposed documentation prose to match better with the standard "voice" of Shore's documentation.

This repository is primarily maintained by a small team at Autodesk along with their other responsibilities, so unfortunately we cannot always respond promptly to pull requests, particularly if they do not relate to an existing GitHub issue where the Shore team has already participated and indicated willingness to work on the issue or accept PRs for the proposal. We *are* grateful for all contributions however, and will give feedback on pull requests as soon as we're able.

## Getting Your Pull Requests Merged

It is much easier to review pull requests that are:

1. **Well-documented:** Try to explain in the pull request comments what your change does, why you have made the change, and provide instructions for how to produce the new behavior introduced in the pull request. If you can, provide screen captures or terminal output to show what the changes look like. This helps the reviewers understand and test the change.
2. **Small:** Try to only make one change per pull request. If you found two bugs and want to fix them both, that's *awesome*, but it's still best to submit the fixes as separate pull requests. This makes it much easier for reviewers to keep in their heads all of the implications of individual code changes, and that means the PR takes less effort and energy to merge. In general, the smaller the pull request, the sooner reviewers will be able to make time to review it.
3. **Passing Tests:** Based on how much time we have, we may not review pull requests which aren't passing our tests (look below for advice on how to run unit tests). If you need help figuring out why tests are failing, please feel free to ask, but while we're happy to give guidance it is generally your responsibility to make sure that tests are passing. If your pull request changes an interface or invalidates an assumption that causes a bunch of tests to fail, then you need to fix those tests before we can merge your PR.

If we request changes, try to make those changes in a timely manner. Otherwise, PRs can go stale and be a lot more work for all of us to merge in the future.

Even with everyone making their best effort to be responsive, it can be time-consuming to get a PR merged. It can be frustrating to deal with the back-and-forth as we make sure that we understand the changes fully. Please bear with us, and please know that we appreciate the time and energy you put into the project.
