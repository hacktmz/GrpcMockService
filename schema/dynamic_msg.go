package schema

import (
	//"time"
	//"io/ioutil"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"log"
	//"golang.org/x/net/context"
	proto "github.com/golang/protobuf/proto"
	//otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	//"github.com/opentracing/opentracing-go"
	//"gitlab.rokid-inc.com/open-platform/gopkg/trace/zipkin"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/grpclog"
	//"gopkg.in/natefinch/lumberjack.v2"
	"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	//"gitlab.rokid-inc.com/open-platform/mockservice/schema"
	"github.com/jhump/protoreflect/desc/protoparse"
)

var proto_map = make(map[string]*desc.FileDescriptor, 0)

func LoadFileProto(file string, stream []byte) error {
	//proto_stream = b
	parser := new(protoparse.Parser)
	GetBytes := func(filename string) (io.ReadCloser, error) {
		reader := ioutil.NopCloser(bytes.NewReader(stream))
		return reader, nil
	}
	parser.Accessor = GetBytes
	fds, err := parser.ParseFiles(file)
	if err != nil {
		return err
	}
	//proto.RegisterFile(file, fds[0])
	proto_map[file] = fds[0]
	return nil
}

func LoadFileDescriptor(file string) (*desc.FileDescriptor, error) {
	fd, ok := proto_map[file]
	if ok {
		return fd, nil
	} else {
		return nil, errors.New("not have FileDescriptor")
	}
}

func GetRequestFullName(file string, servName string, servMethod string) (string, error) {
	fd, err := LoadFileDescriptor(file)
	if err != nil {
		log.Println("Not Found LoadFileDescriptor %s", file)
		return "nil", err
	}
	servDes := fd.FindService(servName)
	ReqName := servDes.FindMethodByName(servMethod).GetInputType().GetFullyQualifiedName()
	log.Println("GetServices = : %v \n ReqName = : %v \n", servDes, ReqName)
	return ReqName, nil
}

func GetResponseFullName(file string, servName string, servMethod string) (string, error) {
	fd, err := LoadFileDescriptor(file)
	if err != nil {
		log.Println("Not Found LoadFileDescriptor %s", file)
		return "nil", err
	}
	servDes := fd.FindService(servName)
	ResName := servDes.FindMethodByName(servMethod).GetOutputType().GetFullyQualifiedName()
	log.Println("GetServices = : %v \n ResName = : %v \n", servDes, ResName)
	return ResName, nil
}

func IsServerStreaming(file string, servName string, servMethod string) (bool, error) {
	fd, err := LoadFileDescriptor(file)
	if err != nil {
		log.Println("Not Found LoadFileDescriptor %s", file)
		return false, err
	}
	servDes := fd.FindService(servName)
	isStream := servDes.FindMethodByName(servMethod).IsServerStreaming()
	log.Println("GetServices = : %v \n IsServerStreaming = : %v \n", servDes, isStream)
	return isStream, nil
}

func IsClientStreaming(file string, servName string, servMethod string) (bool, error) {
	fd, err := LoadFileDescriptor(file)
	if err != nil {
		log.Println("Not Found LoadFileDescriptor %s", file)
		return false, err
	}
	servDes := fd.FindService(servName)
	isStream := servDes.FindMethodByName(servMethod).IsClientStreaming()
	log.Println("GetServices = : %v \n IsClientStreaming = : %v \n", servDes, isStream)
	return isStream, nil
}

//name = pb请求全路径
func NewProtoMessage(name string, file string) (*dynamic.Message, error) {
	fd, err := LoadFileDescriptor(file)
	if err != nil {
		log.Println("Not Found LoadFileDescriptor %s", name)
		return nil, err
	}

	md := fd.FindMessage(name) //md为文件描述数组
	if md != nil {
		p := dynamic.NewMessage(md)
		log.Println("NewProtoMessage(%s) from file %v success\n ", name, p)
		return p, nil
	}
	log.Println("Not FindMessage")
	return nil, fmt.Errorf("Not Found %s", name)
}

//先用proto解码失败再用json解码
func ConvertMsg(msgName string, body []byte, protoFile string) (*dynamic.Message, error) {
	reqc, err := NewProtoMessage(msgName, protoFile)
	if err != nil {
		return nil, fmt.Errorf("NewProtoMessage(%s): %s", msgName, err)
	}
	log.Println(" open.request  = %s", msgName)

	err = proto.Unmarshal(body, reqc)
	if err != nil {
		log.Printf("Unmarshal body to %s: %s faild,will jsonpb.Unmarshal ", msgName, err)
		err = jsonpb.Unmarshal(bytes.NewReader(body), reqc)
		if err != nil {
			//log.Println(body)
			return nil, fmt.Errorf("jsonpb.Unmarshal body to %s: %s", msgName, err)
		}
	}

	return reqc, nil
}
