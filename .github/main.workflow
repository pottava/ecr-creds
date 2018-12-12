workflow "Test & build" {
  on = "push"
  resolves = ["TestResult"]
}

workflow "Release a new version" {
  on = "release"
  resolves = ["ReleaseResult"]
}

action "Branch" {
  uses = "actions/bin/filter@master"
  args = "branch"
}

action "Deps" {
  uses = "pottava/github-actions/go/deps@master"
}

action "Lint" {
  needs = ["Branch", "Deps"]
  uses = "pottava/github-actions/go/lint@master"
}

action "Test" {
  needs = ["Deps"]
  uses = "pottava/github-actions/go/test@master"
}

action "Build" {
  needs = ["Deps"]
  uses = "pottava/github-actions/go/build@master"
  env = {
    BUILD_OPTIONS = "-X main.version=${version}-${GITHUB_SHA:0:7} -X main.date=$(date '+%Y-%m-%d')"
  }
}

action "TestResult" {
  needs = ["Lint", "Test", "Build"]
  uses = "actions/bin/debug@master"
}

action "Tags" {
  uses = "actions/bin/filter@master"
  args = "tag v*"
}

action "ReleaseBuild" {
  needs = ["Tags", "Deps"]
  uses = "pottava/github-actions/go/build@master"
  env = {
    BUILD_OPTIONS = "-X main.version=${version}-${GITHUB_SHA:0:7} -X main.date=$(date '+%Y-%m-%d')"
  }
}

action "Release" {
  needs = ["ReleaseBuild"]
  uses = "pottava/github-actions/github/release@master"
  secrets = ["GITHUB_TOKEN"]
}

action "ReleaseResult" {
  needs = ["Release"]
  uses = "actions/bin/debug@master"
}
