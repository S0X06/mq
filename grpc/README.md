protoc -I ./receive ./receive/receive.proto --go_out=plugins=grpc:receive


protoc -I ./answer ./answer/answer.proto --go_out=plugins=grpc:answer