from golang:1.3

run apt-get update && apt-get install -y git-core cmake
run apt-get install -y pkg-config

run useradd gobuddy
user gobuddy

# add vendored deps
add ./Godeps/_workspace/src /go/src

# add src and setup for work on the project
add . /go/src/github.com/nathanleclaire/tarzan
workdir /go/src/github.com/nathanleclaire/tarzan
run go build

entrypoint ["./tarzan"]
