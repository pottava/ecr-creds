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
  uses = "supinf/github-actions/go/deps@master"
}

action "Lint" {
  needs = ["Branch", "Deps"]
  uses = "supinf/github-actions/go/lint@master"
}

action "Test" {
  needs = ["Deps"]
  uses = "supinf/github-actions/go/test@master"
}

action "Build" {
  needs = ["Deps"]
  uses = "supinf/github-actions/go/build@master"
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
  needs = ["Deps"]
  uses = "supinf/github-actions/go/build@master"
  env = {
    BUILD_OPTIONS = "-X main.version=${version}-${GITHUB_SHA:0:7} -X main.date=$(date +%Y-%m-%d --utc)"
  }
}

action "Release" {
  needs = ["Tags", "ReleaseBuild"]
  uses = "supinf/github-actions/github/release@master"
  secrets = ["GITHUB_TOKEN"]
}

action "BuildImage" {
  needs = ["Tags", "ReleaseBuild"]
  uses = "supinf/github-actions/docker/build@master"
  args = "pottava/ecr-creds:1.0"
  env = {
    DOCKERFILE = "docker/1.0/Dockerfile"
    BUILD_OPTIONS = "--no-cache"
  }
}

action "TagImage" {
  needs = ["BuildImage"]
  uses = "supinf/github-actions/docker/tag@master"
  env = {
    SRC_IMAGE = "pottava/ecr-creds:1.0"
    DST_IMAGE = "pottava/ecr-creds:latest"
  }
}

action "Login" {
  needs = ["BuildImage"]
  uses = "supinf/github-actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "PushImage" {
  needs = ["TagImage", "Login"]
  uses = "supinf/github-actions/docker/push@master"
  args = "pottava/ecr-creds:1.0,pottava/ecr-creds:latest"
}

action "ReleaseResult" {
  needs = ["Release", "PushImage"]
  uses = "actions/bin/debug@master"
}
