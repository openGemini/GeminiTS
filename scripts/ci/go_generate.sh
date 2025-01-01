#!/usr/bin/env bash

# It is recommended to run in a linux environment.

set -e

has_proto=$(grep "github.com/gogo/protobuf/proto" go.mod | grep -v "indirect" | wc -l)

if [[ $has_proto -ne 0 ]]; then
    echo "Don't import github.com/gogo/protobuf/proto"
    echo "Please execute make go-generate"
    exit 1
fi

go generate -x ./...

sed -i "s#github.com/gogo/protobuf/proto#github.com/openGemini/openGemini/lib/util/lifted/protobuf/proto#g" lib/util/lifted/influx/influxql/internal/internal.pb.go
sed -i "s#github.com/gogo/protobuf/proto#github.com/openGemini/openGemini/lib/util/lifted/protobuf/proto#g" lib/util/lifted/influx/meta/proto/meta.pb.go
sed -i "s#github.com/gogo/protobuf/proto#github.com/openGemini/openGemini/lib/util/lifted/protobuf/proto#g" lib/util/lifted/protobuf/proto/test_proto/test.pb.go

protoc --gogo_out=lib/netstorage/data/ lib/netstorage/data/data.proto

sed -i "7d" lib/netstorage/data/data.pb.go
sed -i "8d" lib/netstorage/data/data.pb.go
sed -i "8G"  lib/netstorage/data/data.pb.go

sed -i "9a\    proto \"github.com/openGemini/openGemini/lib/util/lifted/protobuf/proto\""  lib/netstorage/data/data.pb.go
sed -i "10a\    proto1 \"github.com/openGemini/openGemini/lib/util/lifted/influx/meta/proto\""  lib/netstorage/data/data.pb.go

# for goyacc
for i in `find . -name "*.y"`
do
  cd `dirname $i`
  if [[ `dirname $i` == *"ts-cli"* ]]; then
    goyacc -o parser.go -p QL `basename $i`
    sed -i 's/QLErrorVerbose = false/QLErrorVerbose = true/g' parser.go
  else
    goyacc `basename $i`
    sed -i 's/yyErrorVerbose = false/yyErrorVerbose = true/g' y.go
  fi
  rm -f y.output
  cd -
done

GIT_STATUS=$(git status | grep "Changes not staged for commit") || echo ""
if [[ "$GIT_STATUS" = "" ]]; then
  echo "code already generated by go generate -x ./... AND goyacc *.y"
else
  echo ""
  echo "check generated failed, please generate your code using go generate -x ./..."
  echo "git diff files:"
  git diff --stat | tee
  echo "git diff details: "
  git diff | tee
  exit 1
fi
