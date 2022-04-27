#!/bin/bash
set -e

# Function to read version based on SDK
read_version() {
  local sdk=$1
  case $sdk in
    go)
      if [ -f VERSION ]; then
        cat VERSION
      else
        echo "VERSION file not found for Go"
        exit 1
      fi
      ;;
    javascript)
      if [ -f package.json ]; then
        jq -r '.version' package.json
      else
        echo "package.json not found for JavaScript"
        exit 1
      fi
      ;;
    php)
      if [ -f composer.json ]; then
        jq -r '.version' composer.json
      else
        echo "composer.json not found for PHP"
        exit 1
      fi
      ;;
    java)
      if [ -f pom.xml ]; then
        grep -m1 '<version>' pom.xml | sed -E 's/.*<version>([^<]+)<\/version>.*/\1/'
      else
        echo "pom.xml not found for Java"
        exit 1
      fi
      ;;
    python)
      if [ -f pyproject.toml ]; then
        grep '^version' pyproject.toml | head -n1 | cut -d '"' -f2
      elif [ -f setup.py ]; then
        grep 'version=' setup.py | head -n1 | sed -E "s/.*version=['\"]([^'\"]+)['\"].*/\1/"
      else
        echo "No pyproject.toml or setup.py found for Python"
        exit 1
      fi
      ;;
    ruby)
      if [ -f lib/news_api/version.rb ]; then
        grep "VERSION = " lib/news_api/version.rb | sed -E "s/.*'([^']+)'.*/\1/"
      else
        echo "version.rb not found for Ruby"
        exit 1
      fi
      ;;
    *)
      echo "Unsupported SDK: $sdk"
      exit 1
      ;;
  esac
}

# Function to deploy each SDK
deploy_sdk() {
  local sdk=$1
  echo "Preparing deployment for $sdk SDK..."

  VERSION=$(read_version "$sdk")
  echo "Deploying $sdk SDK version $VERSION"

  case $sdk in
    javascript)
      npm version $VERSION --no-git-tag-version
      npm publish
      ;;
    php)
      sed -i "s/\"version\": \".*\"/\"version\": \"$VERSION\"/" composer.json
      composer validate
      # Upload to Packagist manually or via API
      ;;
    java)
      mvn versions:set -DnewVersion=$VERSION
      mvn deploy
      ;;
    python)
      sed -i "s/version=\".*\"/version=\"$VERSION\"/" setup.py
      python -m build
      python -m twine upload dist/*
      ;;
    ruby)
      sed -i "s/VERSION = '.*'/VERSION = '$VERSION'/" lib/news_api/version.rb
      gem build news_api.gemspec
      gem push news_api-$VERSION.gem
      ;;
    go)
      echo "For Go, code is ready — tagging will handle versioning"
      ;;
  esac

  echo "$sdk SDK deployed successfully"
}

# List of SDKs to deploy
SDKS=(go)

# Deploy each SDK
for sdk in "${SDKS[@]}"; do
  (
    cd "$sdk"
    deploy_sdk "$sdk"
  )
done

# After all deployments, read the version properly
# Prefer Go's VERSION file if exists, else fallback to first SDK

if [ -f VERSION ]; then
  FINAL_VERSION=$(read_version "go")
else
  first_sdk="${SDKS[0]}"
  (
    cd "$first_sdk"
    FINAL_VERSION=$(read_version "$first_sdk")
    echo "$FINAL_VERSION" > /tmp/deploy_version.txt
  )
  FINAL_VERSION=$(cat /tmp/deploy_version.txt)
  rm /tmp/deploy_version.txt
fi

# Always tag the repository
git tag -a "v$FINAL_VERSION" -m "Release v$FINAL_VERSION"
git push origin "v$FINAL_VERSION"
echo "Repository tagged with v$FINAL_VERSION"

echo "All SDKs deployed successfully."
