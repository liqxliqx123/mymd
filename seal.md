## seal

```go
server := app.Default()
```



```go
o.AddInitHook(initGRPCServerApplicationHook)
```



```go
func initGRPCServerApplicationHook(o *Application) error {
        ctx := context.Background()
        cb := &test.Callbacks{Debug: o.Config.Service.LogLevel == "DEBUG"}

  			//初始化db
        db, err := o.DBManager.GetMysqlDatabase("mysql")
        if err != nil {
               return errors.Errorf("failed to get mysql db: %v", err)
        }

        redisDB, redis_err := o.DBManager.GetRedisDatabase("redis")
        if redis_err != nil {
               return errors.Errorf("failed to get mysql db: %v", err)
        }

        sealDB, err := initGormDatabasesApplicationHook(o, "mysql")
        if err != nil {
               return errors.Errorf("failed to get mysql db: %v", err)
        }
  
  
				//初始化服务器对象
        srv, err := servgrpc.NewServiceServer(
               ctx,
               cb,
               db,
               redisDB,
               sealDB,
        )
        if err != nil {
               return err
        }
				
  		 //创建服务端实例
        o.GRPCServer = grpcext.NewGRPCServer(o.Config.Service.Name)

        // Register basic service on gRPC server.
        registerServer(o.GRPCServer, srv)

        // serve swagger
        mux := http.NewServeMux()
        openapi.HandleOpenAPI(mux, o.Config.Service.U	RL, pb.Seal)

        // serve grpc gateway
  			//grpc网关，将http请求或响应转换成grpc请求，将grpc响应转换成http响应
        gwmux := runtime.NewServeMux(
               runtime.WithForwardResponseOption(httpResponseHeadersModifier),
        )
        addr := fmt.Sprintf("%s:%d", o.Config.Service.Host, o.Config.Service.Port)
        // TODO: should fix dial options for stage and production env
        if err := pb.RegisterServiceHandlerFromEndpoint(context.Background(), gwmux, addr, []stdgrpc.DialOption{stdgrpc.WithInsecure()}); err != nil {
               return err
        }

        gwmux.HandlePath("POST", "/v1/files:upload", servghttp.HandleBinaryFileUpload)
        gwmux.HandlePath("GET", "/v1/files", servghttp.HandleBinaryFileUpload)

        mux.Handle("/", gwmux)

        //允许跨域访问
  		  o.HTTPServer.Any("/*any", func(ctx *gin.Context) {
               allowCORS(mux).ServeHTTP(ctx.Writer, ctx.Request)
        })

        o.AddDestroyHook(func(o *Application) error {
               o.GRPCServer.GracefulStop()
               return nil
        })

        return nil
}
```

注册envoyproxy/go-control-plane的xDs服务 到grpc server

```
func registerServer(grpcServer *grpc.Server, server server.ServiceServer) {
        discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
        endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
        clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
        routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
        listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)
        secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, server)
        runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
        sealservice.RegisterServiceServer(grpcServer, server)
        envoy_service_auth_v3.RegisterAuthorizationServer(grpcServer, server)

        reflection.Register(grpcServer)
        grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
}
```