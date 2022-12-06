#!/bin/bash
function getGo(){
  curl -L -o go.tgz https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go.tgz 
  rm -rf go.tgz
  echo 'export PATH=$PATH:/usr/local/go/bin'>>~/.bashrc
  source ~/.bashrc
}
go version || getGo
go build -o bin/runner src/*.go
