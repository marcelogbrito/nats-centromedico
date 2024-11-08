package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/marcelogbrito/nats-centromedico/registro-service"
	"github.com/marcelogbrito/nats-centromedico/shared"
	"github.com/nats-io/nats.go"
)

func main() {

	var (
		showHelp     bool
		showVersion  bool
		serverListen string
		natsServers  string
		dbUser       string
		dbPass       string
		dbName       string
		dbPort       string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: registro-service [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	//Setup de flags padrão
	flag.BoolVar(&showHelp, "help", false, "Mostrar ajuda")
	flag.BoolVar(&showVersion, "version", false, "Mostrar versão")
	flag.StringVar(&serverListen, "listen", "0.0.0.0:9090", "Rede host:port to listen on")
	flag.StringVar(&natsServers, "nats", nats.DefaultURL, "Lista de servidores NATS para conectar")
	flag.StringVar(&dbUser, "dbUser", "", "Username do banco de dados")
	flag.StringVar(&dbPass, "dbPassword", "", "Senha do BAnco de dados")
	flag.StringVar(&dbName, "dbName", "", "Nome do banco de dados")
	flag.StringVar(&dbPort, "dbPort", "", "Porta do banco de dados")
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "Microsserviços Centro Médico com NATS - Serviço de Registro v%s\n", registro.Version)
		os.Exit(0)
	}
	log.Printf("Iniciando Microsserviços Centro Médico com NATS - Serviço de Registro versão %s", registro.Version)

	// Registro de novo componente no sistema
	comp := shared.NewComponent("registro-service")

	// Conecta ao NATS e configura setup de descoberta de subscription
	err := comp.SetupConnectionToNATS(natsServers)
	if err != nil {
		log.Fatal(err)
	}

	// Conecta ao Banco de Dados
	err = comp.SetupConnectionToDB("mysql", "example_user:example_password@tcp(localhost:3307)/registro_db")
	if err != nil {
		log.Fatal(err)
	}

	s := registro.Server{
		Component: comp,
	}

	err = s.ListenAndServe(serverListen)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Escutando por requests HTTP em %v", serverListen)
	runtime.Goexit()
}
