echo "Golang Deps Update Starting..."

cd "$(dirname "$0")/../.."

cd services/cmd/server/backend && go mod tidy && go mod vendor && go get -u

cd ..

cd judgecore && go mod tidy && go mod vendor && go get -u

echo "Golang Deps Update Finished"