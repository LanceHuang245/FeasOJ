#!/bin/bash

echo "Building FeasOJ Backend and JudgeCore for Linux x86_64..."

start_time=$(date +%s)

show_spinner() {
  local pid=$1
  local delay=0.1
  local spinstr='|/-\'
  local temp
  while ps -p $pid > /dev/null; do
    temp=${spinstr#?}
    printf " [%c]  " "$spinstr"
    spinstr=$temp${spinstr%"$temp"}
    sleep $delay
    printf "\b\b\b\b\b\b"
  done
  printf "    \b\b\b\b"
}

cd "$(dirname "$0")/../.."

(
  cd services/cmd/app/backend && GOOS=linux go build -o backend-linux main.go
) &
show_spinner $!

(
  cd services/cmd/app/judgecore && GOOS=linux go build -o judgecore-linux main.go
) &
show_spinner $!

end_time=$(date +%s)
duration=$((end_time - start_time))

echo "Build Finished in $duration seconds"

echo "Enter any key to exit..."
read -n 1 -s
exit 0