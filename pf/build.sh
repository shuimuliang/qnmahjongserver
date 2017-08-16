# https://github.com/gogo/protobuf

# Install the protoc-gen-gofast binary
# go get -u github.com/gogo/protobuf/protoc-gen-gofast

# Use it to generate faster marshaling and unmarshaling go code for your protocol buffers.
# protoc --gofast_out=. pf.proto


# Other binaries are also included:
# protoc-gen-gogofast (same as gofast, but imports gogoprotobuf)
# protoc-gen-gogofaster (same as gogofast, without XXX_unrecognized, less pointer fields)
# protoc-gen-gogoslick (same as gogofaster, but with generated string, gostring and equal methods)

# Installing any of these binaries is easy. Simply run:
# go get github.com/gogo/protobuf/proto
# go get github.com/gogo/protobuf/protoc-gen-gogofast
# go get github.com/gogo/protobuf/protoc-gen-gogofaster
# go get github.com/gogo/protobuf/protoc-gen-gogoslick
# go get github.com/gogo/protobuf/gogoproto

# protoc --gogofast_out=. pf.proto
# protoc --gogofaster_out=. pf.proto
cd ../proto
svn update
cd ../pf
sed '3s/mj/pf/' ../proto/pf.proto > pf.proto
protoc --gogoslick_out=. pf.proto
mv pf.pb.go pf.go