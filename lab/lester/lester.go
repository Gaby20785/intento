package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	pb1 "lester/proto/negociacion"

	"google.golang.org/grpc"
)

// Asegura compatibilidad si se añaden más RPCs al servicio.
type servidor struct {
	pb1.UnimplementedNegociacionServer
	contratos []*pb1.ContratoResponse
}

// Implementacion de SolicitarContrato
func (s *servidor) SolicitarContrato(ctx context.Context, req *pb1.ContratoRequest) (*pb1.ContratoResponse, error) {
	log.Println("Lester: Michael solicitó una oferta")

	//90% de probabilidad de que tenga una oferta
	if rand.Intn(100) < 90 && len(s.contratos) > 0 {
		// Se elige un indice cualquiera
		idx := rand.Intn(len(s.contratos))
		contrato := s.contratos[idx]

		// y se saca de la lista (para que no se repita esa oferta)
		s.contratos = append(s.contratos[:idx], s.contratos[idx+1:]...)

		log.Printf("Lester: Enviando oferta: Botin=%d, ProbExitoFranklin=%d, ProbExitoTrevor=%d, RiesgoPolicial=%d",
			contrato.Botin, contrato.ProbExitoFranklin, contrato.ProbExitoTrevor, contrato.RiesgoPolicial)
		return contrato, nil
	}

	log.Println("Lester: No hay oferta disponible, intente más tarde")
	return &pb1.ContratoResponse{}, nil
}

// Implementacion AceptarContrato
func (s *servidor) AceptarContrato(ctx context.Context, req *pb1.DecisionRequest) (*pb1.DecisionResponse, error) {
	log.Println("Lester: Michael aceptó la oferta")
	return &pb1.DecisionResponse{}, nil
}

// Implementacion RechazarContrato
func (s *servidor) RechazarContrato(ctx context.Context, req *pb1.DecisionRequest) (*pb1.DecisionResponse, error) {
	log.Println("Lester: Michael rechazó la oferta")
	return &pb1.DecisionResponse{}, nil
}

// Leer ofertas desde cvs
func leerCSV(path string) []*pb1.ContratoResponse {
	f, err := os.Open(path)

	//error al abrir
	if err != nil {
		log.Fatalf("Error abriendo CSV: %v", err)
	}
	defer f.Close()

	//leer
	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error leyendo CSV: %v", err)
	}

	var contratos []*pb1.ContratoResponse
	for _, rec := range records {
		if len(rec) != 4 {
			continue // saltar si no cumplen con el formato
		}
		botin, _ := strconv.Atoi(rec[0])
		probF, _ := strconv.Atoi(rec[1])
		probT, _ := strconv.Atoi(rec[2])
		riesgo, _ := strconv.Atoi(rec[3])
		contratos = append(contratos, &pb1.ContratoResponse{
			Botin:             int32(botin),
			ProbExitoFranklin: int32(probF),
			ProbExitoTrevor:   int32(probT),
			RiesgoPolicial:    int32(riesgo),
		})
	}
	return contratos
}

func main() {
	contratos := leerCSV("ofertas_grande.csv")

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Error escuchando: %v", err)
	}

	//se crea el servidor
	s := grpc.NewServer()
	pb1.RegisterNegociacionServer(s, &servidor{contratos: contratos})
	log.Println("Lester: Servidor gRPC escuchando en :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error sirviendo: %v", err)
	}
}
