package main

import (
	"context"
	"log"
	"time"

	pb2 "michael/proto/distraccion"
	pb1 "michael/proto/negociacion"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address_lester   = "lester:50051"
	address_franklin = "franklin:50052"
	address_trevor   = "trevor:50053"
)

func main() {

	// ************** Fase 1 ******************

	//Conexion con lester

	conn, err := grpc.Dial(address_lester, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a Lester: %v", err)
	}
	defer conn.Close()

	michael := pb1.NewNegociacionClient(conn)

	//Espacio para guardar la oferta (puntero a ContratoResponse)
	var contratoAceptado *pb1.ContratoResponse
	rechazos := 0

	//bucle hasta aceptar algun contrato
	for {
		contrato, err := michael.SolicitarContrato(context.Background(), &pb1.ContratoRequest{})

		//En caso de error
		if err != nil {
			log.Println("Error al soliictar contrato:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		//Caso en el que no hay oferta
		if contrato.Botin == 0 && contrato.ProbExitoFranklin == 0 && contrato.ProbExitoTrevor == 0 && contrato.RiesgoPolicial == 0 {
			log.Println("Michael: No hay oferta disponible, esperando...")
			time.Sleep(2 * time.Second)
			continue
		}

		//Contrato recibido
		log.Printf("Michael: Oferta recibida -> Botin: %d | ProbExitoFranklin: %d | ProbExitoTrevor: %d | RiesgoPolicial: %d",
			contrato.Botin, contrato.ProbExitoFranklin, contrato.ProbExitoTrevor, contrato.RiesgoPolicial)

		//Se rechaza si algun campo falta en el contrato
		if contrato.Botin == 0 || contrato.ProbExitoFranklin == 0 || contrato.ProbExitoTrevor == 0 || contrato.RiesgoPolicial == 0 {
			log.Println("Michael: Oferta incompleta, rechazando")
			rechazos++
			_, _ = michael.RechazarContrato(context.Background(), &pb1.DecisionRequest{})

			if rechazos%3 == 0 {
				log.Println("Michael: Esperando 10s por paciencia de Lester...")
				time.Sleep(10 * time.Second)
			}
			continue
		}

		// Si la oferta cumple con las condiciones se acepta, de otra forma se rechaza
		if (contrato.ProbExitoFranklin > 50 || contrato.ProbExitoTrevor > 50) && contrato.RiesgoPolicial < 80 {
			_, _ = michael.AceptarContrato(context.Background(), &pb1.DecisionRequest{})

			log.Println("Michael: Oferta aceptada")
			contratoAceptado = contrato
			break
		} else {
			rechazos++
			_, _ = michael.RechazarContrato(context.Background(), &pb1.DecisionRequest{})
			log.Println("Michael: Oferta rechazada")
			if rechazos%3 == 0 {
				log.Println("Michael: Esperando 10s por paciencia de Lester...")
				time.Sleep(10 * time.Second)
			}
		}
	}

	// ************** Fase 2 ******************

	//Elegir personaje para fase 2 (en base a sus probabilidades de exito)
	var personaje string
	var prob int32
	var addrDist string

	if contratoAceptado.ProbExitoFranklin > contratoAceptado.ProbExitoTrevor {
		personaje = "Franklin"
		prob = contratoAceptado.ProbExitoFranklin
		addrDist = address_franklin
	} else {
		personaje = "Trevor"
		prob = contratoAceptado.ProbExitoTrevor
		addrDist = address_trevor
	}

	log.Printf("Michael: Enviando a %s a la primera misión con probabilidad %d%%", personaje, prob)

	connDist, err := grpc.Dial(addrDist, grpc.WithTransportCredentials(insecure.NewCredentials()))

	//Error al conectar
	if err != nil {
		log.Fatalf("No se pudo conectar a %s: %v", personaje, err)
	}

	defer connDist.Close()

	michael2 := pb2.NewDistraccionClient(connDist)

	resp, err := michael2.IniciarDistraccion(context.Background(), &pb2.DistraccionRequest{
		ProbExito: prob,
	})

	if err != nil {
		log.Fatalf("Error iniciando distracción: %v", err)
	}

	if resp.Confirmacion {
		log.Printf("Michael: Confirmacion de inicio de mision recibida")
	} else {
		log.Printf("Michael: Confirmacion de inicio de mision no recibida")
	}

	if resp.Exito {
		log.Printf("Michael: Misión de %s completada con éxito", personaje)
	} else {
		log.Printf("Michael: Misión de %s fallida, motivo: %s", personaje, resp.Motivo)
	}

	// ************** Fase 3 ******************

	// ************** Fase 4 ******************

}
