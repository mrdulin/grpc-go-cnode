package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"

	"google.golang.org/grpc/reflection"

	"github.com/mrdulin/grpc-go-cnode/configs"
	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/topic"
	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"
	"github.com/mrdulin/grpc-go-cnode/internal/utils/grpclogger"
	api "github.com/mrdulin/grpc-go-cnode/internal/utils/http"
	"github.com/mrdulin/grpc-go-cnode/internal/utils/interceptors"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	conf   *viper.Viper
	logger grpclog.LoggerV2
)

func init() {
	conf = configs.Read()
	logger = grpclogger.New(conf.GetString(configs.GRPC_GO_LOG_SEVERITY_LEVEL), conf.GetString(configs.GRPC_GO_LOG_VERBOSITY_LEVEL))
}

func main() {

	// Start z-Pages server.
	go func() {
		mux := http.NewServeMux()
		zpages.Handle(mux, "/debug")
		log.Fatal(http.ListenAndServe("127.0.0.1:8081", mux))
	}()

	stackdriverExporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: conf.GetString(configs.GOOGLE_CLOUD_PROJECT),
		TraceClientOptions: []option.ClientOption{
			option.WithCredentialsFile(conf.GetString(configs.TRACE_ADMIN_CREDENTIAL_FILE)),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Register stats and trace exporters to export
	// the collected data.
	trace.RegisterExporter(stackdriverExporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	view.RegisterExporter(&exporter.PrintExporter{})

	// Register the views to collect server request count.
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}

	port := conf.GetString(configs.PORT)
	baseurl := conf.GetString(configs.BASE_URL)
	if baseurl == "" {
		log.Fatal("missing api url")
	}
	//lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}

	// create TLS server
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	certFile := path.Join(dir, "./assets/server.crt")
	keyFile := path.Join(dir, "./assets/server.key")
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.Creds(creds),
		grpc.UnaryInterceptor(interceptors.NewUnaryInterceptor(logger)),
	)
	// health check
	hs := health.NewServer()
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, hs)

	// reflection
	reflection.Register(grpcServer)

	// application services
	httpClient := api.NewClient()
	userServiceImpl := user.NewUserServiceImpl(httpClient, baseurl)
	topicServiceImpl := topic.NewTopicServiceImpl(httpClient, baseurl)
	user.RegisterUserServiceServer(grpcServer, userServiceImpl)
	topic.RegisterTopicServiceServer(grpcServer, topicServiceImpl)

	// create mux for handling gRPC and HTTP request
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("ok"))
	})

	log.Printf("gRPC and HTTP server are listening on port: %s\n", port)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", port), certFile, keyFile, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//fmt.Println("ProtoMajor:", request.ProtoMajor)
		//fmt.Println("Content-Type:", request.Header.Get("Content-Type"))
		if request.ProtoMajor != 2 {
			mux.ServeHTTP(writer, request)
			return
		}
		if strings.Contains(request.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(writer, request)
			return
		}
		mux.ServeHTTP(writer, request)
	})))
	//log.Fatal(grpcServer.Serve(lis))
}
