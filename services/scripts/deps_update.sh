echo "Golang Deps Update Starting..."

cd "$(dirname "$0")/../.."

cd services/cmd && go mod tidy && go get -u

echo "Golang Deps Update Finished"