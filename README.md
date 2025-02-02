# Kube Prow Style Bot 🌊

test

<div align="center">
    <img src="./docs/img/logo.png" width=30%>
</div>

The kube prow style bot helps manage your issues and PRs via GitHub Action, which is easy-to-use and lightweight.

## Installation

Take this as an example to use this bot: [examples](./examples/kube-prow-bot.yaml). You should edit it firstly to set up maintainers, approvers and reviewers of your project.

We support read roles from OWNERS file, see https://github.com/Xunzhuo/kube-prow-bot/blob/main/OWNERS:

``` yaml
# See the OWNERS docs at https://go.k8s.io/owners

maintainers:
  - Xunzhuo

approvers:
  - Xunzhuo

reviewers:
  - Xunzhuo

```

We also support define Roles in env (We will prioritize reading roles from OWNERS):

```yaml
env:
  REVIEWERS: |-
    Reviewer1
    Reviewer2
  APPROVERS: |-
    Approver1
    Approver2
  MAINTAINERS: |-
    Maintainer1
    Maintainer2
```

And then you can manage the commands which different roles can access and use in env:

```yaml
  # This commands is for anyone who can use it
  COMMON_PLUGINS: |-
    assign
    unassign
    kind
    remove-kind
    cc
    uncc
  # This commands is for author of issue or PR
  AUTHOR_PLUGINS: |-
    retest
  # This commands is for organization member or repository member
  MEMBERS_PLUGINS: |-
    good-first-issue
    help-wanted
    close
    reopen
  # This commands is for in the REVIEWERS environment variable
  REVIEWERS_PLUGINS: |-
    area
    remove-area
    lgtm
    hold
    retitle
  # This commands is for in the APPROVERS environment variable
  APPROVERS_PLUGINS: |-
    merge
    approve
    rebase
  # This commands is for in the MAINTAINERS environment variable
  MAINTAINERS_PLUGINS: |-
    milestone
    remove-milestone
    priority
    remove-priority
```

Then you can just copy it into `.github/workflows/` and it should work now.

## Usage

List commands which kube prow style bot supports now:

|     Command      |                      Usage                      | Default Role |
| :--------------: | :---------------------------------------------: | :----------: |
|      assign      |        `/assign`, `/assign @[GitHub ID]`        |    Anyone    |
|     unassign     |      `/unassign`, `/unassign @[GitHub ID]`      |    Anyone    |
|       kind       |                  `/kind label`                  |    Anyone    |
|   remove-kind    |              `/remove-kind label`               |    Anyone    |
|        cc        |               `/cc @[GitHub ID]`                |    Anyone    |
|       uncc       |              `/uncc @[GitHub ID]`               |    Anyone    |
|      retest      |                    `/retest`                    |    Author    |
| good-first-issue | `/good-first-issue`, `/good-first-issue cancel` |    Member    |
|   help-wanted    |      `/help-wanted`, `/help-wanted cancel`      |    Member    |
|      close       |                    `/close`                     |    Member    |
|      reopen      |                    `/reopen`                    |    Member    |
|       area       |                  `/area label`                  |   Reviewer   |
|   remove-area    |              `/remove-area label`               |   Reviewer   |
|       lgtm       |              `/lgtm`, `/lgtm cancel`            |   Reviewer   |
|       hold       |              `/hold`, `/hold cancel`            |   Reviewer   |
|      retitle     |                `/retitle xx`                    |   Reviewer   |
|     approve      |              `/approve`, `/approve cancel`      |   Approver   |
|      rebase      |                    `/rebase`                    |   Approver   |
|      merge       |   `/merge` , `/merge rebase`, `/merge squash`   |   Approver   |
|    milestone     |                `/milestone name`                |  Maintainer  |
| remove-milestone |            `/remove-milestone name`             |  Maintainer  |
|     priority     |                `/priority name`                 |  Maintainer  |
| remove-priority  |             `/remove-priority name`             |  Maintainer  |

## Automation

Pull Request will be merged automatically when it has been `/approve` and `/lgtm`.

But it will be stopped when it has `/hold`, use `/hold cancel` to process to merge it.

## Contact

Xunzhuo Liu – [@Xunzhuo](https://github.com/Xunzhuo) – mixdeers@gmail.com

## License

Distributed under the Apache license. See ``LICENSE`` for more information.

## Contributing

1. Fork it.
2. Create your feature branch (`git checkout -b feature/fooBar`).
3. Commit your changes (`git commit -am 'Add some fooBar'`).
4. Push to the branch (`git push origin feature/fooBar`).
5. Create a new Pull Request.
