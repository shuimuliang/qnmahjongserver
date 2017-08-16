#!/bin/sh
cp pf.proto pf.proto.bak
sed -i '' -e '4i \
	import "github.com/gogo/protobuf/gogoproto/gogo.proto";' pf.proto
sed -i '' -e '5i \
	option (gogoproto.marshaler_all) = true;' pf.proto
sed -i '' -e '6i \
	option (gogoproto.sizer_all) = true;' pf.proto

sed -i '' -e '7i \
	option (gogoproto.unmarshaler_all) = true;' pf.proto

protoc --proto_path=$GOPATH/src/github.com/gogo/protobuf/protobuf:../../../:. --gogo_out=.   pf.proto

sed  -i '' 's|package pf|package model|g' pf.pb.go
mv pf.pb.go ../app-logic/model/msg.pb.go
mv -f pf.proto.bak pf.proto

# # pub
#
# cp pub.proto pub.proto.bak
# sed -i '' -e '4i \
# 	import "github.com/gogo/protobuf/gogoproto/gogo.proto";' pub.proto
# sed -i '' -e '5i \
# 	option (gogoproto.marshaler_all) = true;' pub.proto
# sed -i '' -e '6i \
# 	option (gogoproto.sizer_all) = true;' pub.proto
#
# sed -i '' -e '7i \
# 	option (gogoproto.unmarshaler_all) = true;' pub.proto
#
# protoc --proto_path=$GOPATH/src/github.com/gogo/protobuf/protobuf:../../../:. --gogo_out=.  pub.proto
# mkdir -p ../pub/
# mv pub.pb.go ../pub
#
# mv -f pub.proto.bak pub.proto
