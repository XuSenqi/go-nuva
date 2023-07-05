#! /bin/sh

# _ROOT: 工作目录
# ROOT: 总是正确指向build脚本所在目录
_ROOT="$(pwd)" && cd "$(dirname "$0")" && ROOT="$(pwd)"
PJROOT="$ROOT"

if [ -n "$1" ]; then
    APPNAME="$1"
fi
if [ -n "$2" ]; then
    if [ "$(echo "$2" | head -c1)" = "/" ]; then
        TARGET="$2"
    else
        TARGET="$_ROOT/$2"
    fi
fi

UNAME="$(uname)"
hash() {
    if [ "$UNAME" = "Darwin" ]; then
        shasum $*
    else
        sha1sum $*
    fi
}

# 检查golang环境
__check_go_version() {
    GO=go

    if ! which $GO >/dev/null ; then
        echo "[Error] go environment not found" >&2
        exit 1
    fi

    if $GO mod 2>&1 | grep -q -i 'unknown command'; then
        echo "[Error] low golang version(should be >=1.11), do not support go mod command"
        exit 1
    fi

    if [ ! -r $PJROOT/go.mod ]; then
        echo "[Error] go.mod not found or not readable"
        exit 1
    fi

    MODULE="$(cat $PJROOT/go.mod | grep ^module | head -n1 | awk '{print $NF}')"
}

__build_magic() {
    cd "$PJROOT"

    if which git >/dev/null 2>&1 && git status >/dev/null 2>&1; then
        _upstream="$(git rev-parse --abbrev-ref @{upstream} 2>/dev/null | cut -d/ -f1)"
        [ -z "$GIT_UPSTREAM" ] && _upstream="origin"
        GIT_REPO="$(git config --get remote.$_upstream.url 2>/dev/null)"
        GIT_BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null)"

        GIT_HASH="$(git log -n1 --pretty=format:%H 2>/dev/null)"
        TAG_NAME="$(git describe --tags --long --match v[0-9]* 2>/dev/null | sed -nE 's/(.*)-[0-9]+-g.{7,}/\1/p')"
        if [ -n "$TAG_NAME" ]; then
            TAG_HASH="$(git rev-list -n1 "$TAG_NAME")"
        fi

        GIT_STATUS_HASH="$(git status -s -uall 2>/dev/null | awk '{print $NF}' | xargs -I{} cat "{}" 2>/dev/null | hash | cut -d' ' -f1)"

        BUILD_MAGIC="$(echo -n -e "$GIT_REPO\x00$GIT_BRANCH\x00$GIT_HASH\x00$TAG_NAME\x00$TAG_HASH\x00$GIT_STATUS_HASH" | hash | cut -d' ' -f1)"
    fi
}

# 搜集待注入的编译环境信息
__env() {
    cd "$PJROOT"

    VERSION="0.0.x"
    RELEASE="1"
    GO_VERSION="$($GO version)"
    BUILD_ID="$(head -c 128 /dev/urandom | hash | cut -d' ' -f1)"
    BUILD_TIME="$(date +%s.%N)"

    if [ -n "$GIT_HASH" ]; then
        GIT_TIME="$(git log -n1 --pretty=format:%at 2>/dev/null)"
        GIT_NUMBER="$(git rev-list --count HEAD 2>/dev/null)"

        if [ -n "$TAG_NAME" ]; then
            TAG_TIME="$(git log -n1 $TAG_HASH --pretty=format:%at)"
            TAG_NUMBER="$(git rev-list --count $TAG_HASH)"

            TAG_MESSAGE="$(git tag -l v[0-9]* -n100 --sort=-v:refname | sed -n "/^$TAG_NAME /,\$p" | sed -E 's/^(v[0-9][^ \t]+)[ \t]{6}/\1\
/')"
            TAG_DIFF="$(git rev-list --count HEAD ^$TAG_HASH)"

            VERSION="$(echo $TAG_NAME | cut -c2-)"
            RELEASE="$((1+TAG_DIFF))"
        fi

        GIT_STATUS_NUMBER="$(git status -s -uall 2>/dev/null | wc -l | awk '{print $1}')"
    fi
}

# 编译，使用go mod做包管理
__build() {
    cd "$PJROOT"

    GCFLAGS="-c 4" #compile concurrency
    if [ -n "$DEVELOP" ]; then
        GCFLAGS="-N -l $GCFLAGS"  #-N: disable optimizations, -l: disable inlining. go tool compile --help
    else
        LDFLAGS="-w -s" #removes symbol table and debugging information. go tool link --help
    fi
    CGO_ENABLED=0 $GO build -gcflags="all=$GCFLAGS" \
        -o "$TARGET" -ldflags "$LDFLAGS \
        -X '$MODULE/app.appname=$APPNAME' \
        -X '$MODULE/app.version=$VERSION' \
        -X '$MODULE/app.release=$RELEASE' \
        -X '$MODULE/app.goVersion=$GO_VERSION' \
        -X '$MODULE/app.projectRoot=$PJROOT' \
        -X '$MODULE/app.gitRepo=$GIT_REPO' \
        -X '$MODULE/app.gitBranch=$GIT_BRANCH' \
        -X '$MODULE/app.gitHash=$GIT_HASH' \
        -X '$MODULE/app.gitTime=$GIT_TIME' \
        -X '$MODULE/app.gitNumber=$GIT_NUMBER' \
        -X '$MODULE/app.tagName=$TAG_NAME' \
        -X '$MODULE/app.tagHash=$TAG_HASH' \
        -X '$MODULE/app.tagTime=$TAG_TIME' \
        -X '$MODULE/app.tagNumber=$TAG_NUMBER' \
        -X '$MODULE/app.tagDiff=$TAG_DIFF' \
        -X '$MODULE/app.tagMessage=$(echo "$TAG_MESSAGE" | head -c 51200 | base64)' \
        -X '$MODULE/app.gitStatusNumber=$GIT_STATUS_NUMBER' \
        -X '$MODULE/app.gitStatusHash=$GIT_STATUS_HASH' \
        -X '$MODULE/app.buildID=$BUILD_ID' \
        -X '$MODULE/app.buildMagic=$BUILD_MAGIC' \
        -X '$MODULE/app.buildTime=$BUILD_TIME'" $MAINFILE

    if [ $? -ne 0 ]; then
        exit 1
    fi
}

main() {
    __check_go_version

    if [ -z "$APPNAME" -o "$APPNAME" = "help" ]; then
        echo "Usage: $0 {app name} [target file]"
        echo "e.g. $0 $(basename $MODULE)"
        echo "e.g. $0 $(basename $MODULE) ./bin/$(basename $MODULE)"
        exit 1
    fi

    if [ -z "$MAINFILE" ]; then
        MAINFILE="$PJROOT/main.go"
    fi

    [ -z "$TARGET" ] && TARGET="$PJROOT/bin/$APPNAME"

    __build_magic
    __env
    __build
}

main
