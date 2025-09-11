echo "Go Lint Checker Starting..."

cd "$(dirname "$0")/../../"

if [ -z "$(which golangci-lint)" ]; then
  echo "golangci-lint not installed"
  exit 1
fi

echo "Now Linting Backend..."
echo "====================================================================================="

cd services/cmd/app/backend

golangci-lint run

echo "====================================================================================="

echo "Now Linting Judgecore..."

echo "====================================================================================="

cd ../judgecore

golangci-lint run

echo "====================================================================================="

echo "Go Lint Checker Finished"