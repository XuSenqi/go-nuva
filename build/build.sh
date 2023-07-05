#! /bin/sh

app=$1
target=$2

if [ -z "$app" ]; then
    echo "ERROR: app is empty!"
    exit 1
fi

if [ -z "$target" ]; then
    echo "ERROR: target is empty!"
    exit 1
fi

osName=$(uname)

hash() {
  if [ $osName = "Darwin" ]; then
    shasum $*
  else
    sha1sum $*
  fi
}

# environment
git_tag="$(git describe --tags --long --match v[0-9]* 2>/dev/null | sed -nE 's/(.*)-[0-9]+-g.{7,}/\1/p')"
version="$(echo "$git_tag" | cut -c2-)"
go_version="$(go version)"
build_rand="$(head -c 128 /dev/urandom | hash | cut -d' ' -f1)"
build_time="$(date +%s)"

if which git 2>/dev/null >/dev/null && git status 2>/dev/null >/dev/null; then
  git_repo="$(git config --get remote.origin.url 2>/dev/null)"
  git_branch="$(git rev-parse --abbrev-ref HEAD 2>/dev/null)"

  tag_hash="$(git rev-list -n 1 "$git_tag")"
  tag_diff="$(git rev-list --count HEAD ^"$tag_hash")"
  git_hash="$(git log -n1 --pretty=format:%H-%at 2>/dev/null)"
  git_number="$((1 + tag_diff))"
  git_status_number="$(git status -s -uall 2>/dev/null | wc -l | awk '{print $1}')"
  git_status_hash="$(git status -s -uall 2>/dev/null | awk '{print $NF}' | xargs -I{} cat "{}" 2>/dev/null | hash | cut -d' ' -f1)"

  git_indicator="$(echo -n -e "$git_repo\x00$git_branch\x00$git_hash\x00$git_number\x00$git_status_number\x00$git_status_hash" | hash | cut -d' ' -f1)"
fi

# build
_module="$(cat go.mod | grep ^module | head -n1 | awk '{print $NF}')"
CGO_ENABLED=0 go build -gcflags="all=-c 4" -o $target -ldflags "\
        -X '$_module/cmd.appname=$app' \
        -X '$_module/cmd.version=$version' \
        -X '$_module/cmd.goVersion=$go_version' \
        -X '$_module/cmd.codeRoot=$CODEROOT' \
        -X '$_module/cmd.gitRepo=$git_repo' \
        -X '$_module/cmd.gitBranch=$git_branch' \
        -X '$_module/cmd.gitHash=$git_hash' \
        -X '$_module/cmd.gitNumber=$git_number' \
        -X '$_module/cmd.gitStatusNumber=$git_status_number' \
        -X '$_module/cmd.gitStatusHash=$git_status_hash' \
        -X '$_module/cmd.buildRand=$build_rand' \
        -X '$_module/cmd.buildIndicator=$git_indicator' \
        -X '$_module/cmd.buildTime=$build_time'" main.go

if [ $? -ne 0 ]; then
    echo "build failed Ooooooh!!!"
    exit 1
else
    echo "build succeed!"
fi
