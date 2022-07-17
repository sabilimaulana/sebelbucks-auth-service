package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/config"
	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/db"
	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/pb"
	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/services"
	"github.com/sabilimaulana/sebelbucks-auth-service/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 30,
	}

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
