

---
# Changes

TODO (describe what your PR introduces)

---

# Type of update

Alter this statement (replace "patch" with what you want) to indicate the type of version bump:

```
+semver: patch
```

> **_NOTE:_**  This works by appending this PR description to the merge commit. Subsequently in the release pipeline it will look for the string above within the commit history in order to know how to bump the version.

Valid values to replace "patch" with are:


| keyword  | example effect    |  alias keyword  |
|----------|-------------------|-----------------|
| `major`  | `1.0.2 -> 2.0.0`  | `breaking`      |
| `minor`  | `1.0.2 -> 1.1.0`  | `feature`       |
| `patch`  | `1.0.2 -> 1.0.3`  | `fix`           |
| `none`   | `1.0.2 -> 1.0.2`  | `skip`          |


> **_NOTE:_**  It is not possible to bump from 0.x to 1.0.0. If the latest version is 0.2.3 for example, the Gitversion tool will treat instructions to perform a major bump as a minor instead, and go `0.2.3 -> 0.3.0`. When you are ready to start releasing major versions, set the line `next-version: 1.0.0` in `.azuredevops/gitversion.yaml`

---
# Work Items

Related work items: SET_WORK_ITEM
