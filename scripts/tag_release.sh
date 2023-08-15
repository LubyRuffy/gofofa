#!/usr/bin/env bash

# 检查是否有未提交的内容
if [[ $(git status) =~ "modified:" ]]; then
  echo "git is currently in a dirty state"
  exit
fi

# 从HISTORY.md中进行版本提取
if [[ ! -f HISTORY.md ]]; then
  echo "cannot found HISTORY.md file"
  exit
fi

VERSION=$(head -n1 HISTORY.md | awk '{print $2}' | tr " " _)
if [[ ! $VERSION =~ ^v[[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+$ ]]; then
  echo "version format error: $VERSION, HISTORY.md first line must be '## v0.0.1 <comment>' format "
  exit
fi

COMMENT=$(head -n1 HISTORY.md | awk '{$1="";$2="";sub("  "," ");print $0;}')
if [[ $COMMENT == "" || $COMMENT == " " ]]; then
  echo "no comment: HISTORY.md first line must be '## v0.0.1 <comment>' format "
  exit
fi

TAG=$(git tag --list "$VERSION")
if [[ $TAG == *"$VERSION"*  ]]; then
  git tag -d "$VERSION"
fi

git tag -a "$VERSION" -m "$COMMENT"
git push origin "$VERSION"
goreleaser release --rm-dist --debug