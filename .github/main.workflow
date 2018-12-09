workflow "Test & build" {
  on = "push"
  resolves = ["TestAndBuild"]
}

workflow "Release a new version" {
  on = "release"
  resolves = ["BuildAndRelease"]
}

action "Branch" {
  uses = "actions/bin/filter@master"
  args = "branch"
}

action "Deps" {
  uses = "pottava/github-actions/go/deps@master"
}

action "Lint" {
  needs = ["Deps"]
  uses = "pottava/github-actions/go/lint@master"
}

action "UnitTest" {
  needs = ["Deps"]
  uses = "pottava/github-actions/go/test@master"
}

action "TestBuild" {
  needs = ["Deps"]
  uses = "pottava/github-actions/go/build@master"
}

action "TestAndBuild" {
  needs = ["Branch", "Lint", "UnitTest", "TestBuild"]
  uses = "actions/bin/debug@master"
}

action "Tags" {
  uses = "actions/bin/filter@master"
  args = "tag v*"
}

action "Release" {
  needs = ["Deps"]
  uses = "pottava/github-actions/go/build@master"
  env = {
    POST_PROCESS = "github_release"
  }
  secrets = ["GITHUB_TOKEN"]
}

action "BuildAndRelease" {
  needs = ["Tags", "Release"]
  uses = "actions/bin/debug@master"
}
