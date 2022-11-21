#!/bin/bash

# we need git tag: "v1.2.3"
TAG=$1
# git commit message
MESSAGE=$2

if [[ $TAG != v* ]]; then
  echo "Error: TAG must start with \"v\", like v1.2.3"
  exit 1
fi

if [ -z "$MESSAGE" ]; then
  echo "Error: MESSAGE is empty!"
  exit 1
fi

echo "Tagging: $TAG"
echo "With message: $MESSAGE"

read -p "Proceed? y/n " yes

if [ $yes != "y" ]; then
  echo "cancelling"
  exit 0
fi

# update VERSION in nvm.go
sed -i "s/const VERSION = \"v.*\"/const VERSION = \"${TAG}\"/g" ./nvm/nvm.go

git add -A
git commit -am "$MESSAGE"
git tag $TAG
git push
git push --tags
