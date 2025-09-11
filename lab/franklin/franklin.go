package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	pb2 "franklin/proto/distraccion"

	"google.golang.org/grpc"
)

// Asegura compatibilidad si se añaden más RPCs al servicio.
type servidor struct {
	pb2.UnimplementedDistraccionServer
}

// Implementacion de IniciarDistraccion
func (s *servidor) IniciarDistraccion(ctx context.Context, req *pb2.DistraccionRequest) (*pb2.DistraccionResponse, error) {
	log.Printf("Franklin: Iniciando distracción con probabilidad %d%%", req.ProbExito)

	turnos := 200 - int(req.ProbExito)
	fracaso := false
	motivo := ""

	for t := 1; t <= turnos; t++ {
		time.Sleep(50 * time.Millisecond)

		// Cada 10 turnos s muestra progreso
		if t%10 == 0 || t == turnos {
			log.Printf("Franklin: Avanzando turno %d/%d", t, turnos)
		}

		// Se revisa si hay imprevisto
		if t == turnos/2 {
			if rand.Intn(100) < 10 {
				fracaso = true
				motivo = "Chop distrajo a Franklin"
				log.Printf("Franklin: Ocurrió un imprevisto, misión fallida")
				break
			}
		}
	}

	if !fracaso {
		log.Printf("Franklin: Misión completada con éxito")
	}

	return &pb2.DistraccionResponse{
		Exito:  !fracaso,
		Motivo: motivo,
	}, nil

}

func main() {
	srv := &servidor{}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Error escuchando: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb2.RegisterDistraccionServer(grpcServer, srv)

	log.Printf("Franklin: Servidor de distracción escuchando en :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error sirviendo: %v", err)
	}
}
