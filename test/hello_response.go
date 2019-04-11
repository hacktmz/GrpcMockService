package main

//go:generate protoc -I ./ --go_out=plugins=grpc:./  ./hello.proto
import (
	"io/ioutil"
	//"bytes"
	"log"
	pb"gitlab.rokid-inc.com/open-platform/mockservice/pbs/hello"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
)
func jsonres() {
	kv := make(map[string]string)
	kv["aaa"] = "aax"
	kv["baaa"] = "bbx"
	marshaler := jsonpb.Marshaler{EmitDefaults: true}
	msg, err := marshaler.MarshalToString(&pb.HelloResponse{
		Message: "this is a HelloResponse",
		Options: kv,
	})
	if err != nil {
		log.Printf("jsonpb.MarshalToString() = %#v", err)
		return 
	}

	//err = jsonpb.Unmarshal(bytes.NewReader([]byte(msg), out))
	log.Printf("Marshal(): %v", msg)
	if err := ioutil.WriteFile("hello.json", []byte(msg), 0644); err != nil {
		log.Printf("WriteFile(): %v", msg)
	}

}

func protores() {
	kv := make(map[string]string)
	kv["aaa"] = "aax"
	kv["baaa"] = "bbx"
	msg, err := proto.Marshal(&pb.HelloResponse{
		Message: "this is a HelloResponse",
		Options: kv,
	})
	if err != nil {
		log.Printf("jsonpb.MarshalToString() = %#v", err)
		return 
	}

	//err = jsonpb.Unmarshal(bytes.NewReader([]byte(msg), out))
	log.Printf("Marshal(): %v", msg)
	if err := ioutil.WriteFile("hello.res", []byte(msg), 0644); err != nil {
		log.Printf("WriteFile(): %v", msg)
	}

}
func main() {
	jsonres()
	protores()

}
