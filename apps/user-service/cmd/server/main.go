package main

import (
	"net"

	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/config"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/postgres"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/repository"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/service"
	userv1 "github.com/shah-dhwanil/grpc-chat/packages/api/gen/user/v1"
	"github.com/shah-dhwanil/grpc-chat/packages/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {
	cfg := config.GetConfig()
	pool,err := postgres.NewPool(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	repo := repository.NewRepository(pool)
	service := service.NewService(repo)
	listener, err := net.Listen("tcp", ":8001")
   	if err != nil {
      panic(err)
   }
   srv := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryAuthInterceptor))
   userv1.RegisterUserServiceServer(srv,service.UserService)
   reflection.Register(srv)
   if e := srv.Serve(listener); e != nil {
      panic(e)
   }
}