package main

import (

	//"time"
	"flag"
	"fmt"
	"io"
	"net"

	"os/signal"
	"syscall"

	"log"

	"github.com/kavu/go_reuseport"
	//"golang.org/x/net/context"
	//proto "github.com/golang/protobuf/proto"
	//otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	//"github.com/opentracing/opentracing-go"
	//"gitlab.rokid-inc.com/open-platform/gopkg/trace/zipkin"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/natefinch/lumberjack.v2"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/protoc-gen-go/descriptor"
	//"github.com/jhump/protoreflect/desc"
	"encoding/json"

	pb "github.com/hacktmz/GrpcMockService/pbs"
	"github.com/hacktmz/GrpcMockService/schema"
	"github.com/jhump/protoreflect/dynamic"
)

type MockServer interface {
	// 定义SayHello方法
}

// 定义helloService并实现约定的接口
type mockService struct{}

// MockService Hello服务
var MockService = mockService{}

type Info struct {
	protoFiles string
	methodName string
	servName   string
}

var res_raw []byte
var method_map = make(map[string][]interface{}, 100)
var ch = make(chan *pb.MockRequest)

//var ch_notify_start = make(chan bool)
var ch_notify_stop = make(chan bool)
var ch_stop = make(chan bool)
var ch_return_err = make(chan string)

var (
	logToFile  = flag.Bool("log_to_file", false, "Log to file?")
	logFile    = flag.String("log_file", "mockservice.log", "The log file")
	logSize    = flag.Int("log_size", 512, "The max size of log file")
	logBackups = flag.Int("log_backups", 32, "The max backups of log file")
	logAges    = flag.Int("log_max_ages", 30, "The max ages of log file")
)

// 暂时不需要用
/*
func TextHandler(method_name string) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error){
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {

		fmt.Println("get massage text in %s", method_name)
		val,ok := method_map[method_name]
		if (ok) {
			fmt.Println("text val = %d",len(val))
			for _, v := range val {
				return v, nil
			}
			return nil,nil
		} else {
			return nil,nil
		}

	}

}
*/

func StreamHandler(info *Info) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		fmt.Println("get massage stream in  %s", info.methodName)
		reqName, _ := schema.GetRequestFullName(info.protoFiles, info.servName, info.methodName)
		msg, err := schema.NewProtoMessage(reqName, info.protoFiles)
		for {
			err = stream.RecvMsg(msg)
			if err == io.EOF {
				break
			}
			fmt.Println("get RecvMsg = %v", msg)
		}

		val, ok := method_map[info.methodName]
		if ok {
			fmt.Println("stream val = %d", len(val))
			for _, v := range val {
				stream.SendMsg(v)
			}
			return nil
		} else {
			return nil
		}
	}
}

type Mock struct {
	listen net.Listener
	serv   *grpc.Server
}

func (m *Mock) StopService() {
	if m.serv != nil {
		m.serv.GracefulStop()
	}
	return
}

func (m *Mock) InitService(cfg *pb.MockRequest) (string, error) {
	port := cfg.Port
	var err error
	m.listen, err = reuseport.NewReusablePortListener("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		errStr := fmt.Sprintf("Failed to listen: %v", err)
		grpclog.Printf(errStr)
		return errStr, err
	}
	serlist := make([]*grpc.ServiceDesc, 0)
	for _, header := range cfg.Headers {
		log.Println("===========processing file : %s ===========\n", header.Filename)
		protoFileBytes := []byte(header.Protofile)
		var dmsg *dynamic.Message
		err = schema.LoadFileProto(header.Filename, protoFileBytes)
		if err != nil {
			errStr := fmt.Sprintf("Failed to open proto: %v", err)
			log.Println(errStr)
			return errStr, err
		}

		Protojson := []byte(header.Protojson)
		m_json := make(map[string]interface{})
		err = json.Unmarshal(Protojson, &m_json)
		if err != nil {
			errStr := fmt.Sprintf("Failed to Unmarshal json: %v", err)
			log.Println(errStr)
			ch_return_err <- errStr
		}
		log.Println("===========UN JSON===== %s", m_json)
		for k, v := range m_json {
			servName := k
			var textMethods []grpc.MethodDesc
			var streamMethods []grpc.StreamDesc
			for k2, v2 := range v.(map[string]interface{}) {
				methodName := k2
				Resfile, err := json.Marshal(v2)
				if err != nil {
					errStr := fmt.Sprintf("Failed to Unmarshal json get resfile: %s", err)
					log.Println(errStr)
					return errStr, err
				}
				resName, _ := schema.GetResponseFullName(header.Filename, servName, methodName)
				dmsg, err = schema.ConvertMsg(resName, Resfile, header.Filename)
				log.Println("===========get service: %s  method:%s", servName, methodName)
				if err != nil {
					errStr := fmt.Sprintf("ConvertMsg err: %s", err)
					log.Println(errStr)
					return errStr, err
				}

				var res_slice []interface{} //method -> res1 res2....
				res_slice = append(res_slice, dmsg)
				method_map[methodName] = res_slice
				info := Info{
					protoFiles: header.Filename,
					servName:   servName,
					methodName: methodName,
				}
				//textMethods = append(textMethods, grpc.MethodDesc{ MethodName: *methodName, Handler: TextHandler(*methodName),})
				streamMethods = append(streamMethods, grpc.StreamDesc{StreamName: methodName, Handler: StreamHandler(&info), ServerStreams: true, ClientStreams: true})
			}

			serviceDesc := &grpc.ServiceDesc{
				ServiceName: servName,
				HandlerType: (*MockServer)(nil),
				Methods:     textMethods,
				Streams:     streamMethods,
				//Metadata: "hello111.proto",
			}
			serlist = append(serlist, serviceDesc)
		}
	}
	if len(serlist) < 1 {
		errStr := fmt.Sprintf("not found service! \n")
		log.Printf(errStr)
		return errStr, err
	}
	m.serv = grpc.NewServer()
	for _, v := range serlist {
		// 注册HelloService
		log.Printf("RegisterService %v", v)
		m.serv.RegisterService(v, MockService)
	}
	grpclog.Println(fmt.Sprintf("Listen on %d", port))
	go m.serv.Serve(m.listen)
	return "ok", nil
}

func main() {
	//grpclog.SetLogger(log.New(ioutil.Discard, "", 0))
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
	log.SetOutput(&lumberjack.Logger{
		Filename:   *logFile,
		MaxSize:    *logSize, // megabytes
		MaxBackups: *logBackups,
		MaxAge:     *logAges, // days
		Compress:   false,    // disabled by default
		LocalTime:  true,     //
	})

	signal.Ignore(syscall.SIGPIPE)

	go Listener()

	map_port := make(map[int32]*Mock) //save prot: serv point
	for {
		select {
		case cfg := <-ch:
			log.Println("received ", cfg, " from ch\n")
			p, ok := map_port[cfg.Port]
			if p != nil && ok {
				p.StopService()
				delete(map_port, cfg.Port)
			}
			newp := new(Mock)
			res_msg, err := newp.InitService(cfg)
			if err == nil {
				map_port[cfg.Port] = newp
			} else {
				newp = nil
			}
			ch_return_err <- res_msg
		case flag := <-ch_stop:
			if flag {
				for k, v := range map_port {
					if v != nil {
						v.StopService()
						delete(map_port, k)
					}
				}
			}
			ch_notify_stop <- true
		}
	}

}
