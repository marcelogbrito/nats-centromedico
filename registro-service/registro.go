package registro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marcelogbrito/nats-centromedico/shared"
	"github.com/nats-io/nuid"
)

const (
	Version = "0.1.0"
)

type Server struct {
	*shared.Component
}

// ListenAndServe pega o endereço de rede e porta que o servidor Http deve vincular e inicia
func (s *Server) ListenAndServe(addr string) error {
	r := mux.NewRouter()
	router := r.PathPrefix("/cmed/paciente/").Subrouter()
	//Handles referentes aos paths
	// GET /cmed/paciente
	router.HandleFunc("/", s.HandleHomeLink)

	// Handle de registro de paciente
	// POST /cmed/paciente/registro
	router.HandleFunc("registro", s.HandleRegistro).Methods("POST")
}

func (s *Server) HandleHomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Serviço de Registro v%s\n", Version))
}

// HandleRegistro processa requests de registro de pacientes
func (s *Server) HandleRegistro(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Bad Requesst", http.StatusBadRequest)
		return
	}

	var registration *shared.RegistrationRequest
	err = json.Unmarshal(body & registration)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Insere dados no banco de dados
	db := s.DB()

	insForm, err := db.Prepare("INSERT INTO paciente_detalhes(id, nome_completo, endereco, sexo, telefone, observacoes) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	insForm.Exec(registration.ID, registration.NomeCompleto, registration.Endereco, registration.Sexo, registration.Telefone,
		registration.Observacoes)

	// Marca o request com um ID para tracing nos logs
	registration.RequestID = nuid.Next()
	fmt.Println(registration)

	//Publica evento no servidor NATS
	nc := s.NATS()

	//var registration_event shared.RegistrationEvent
	tokenNo := generateTokenNumber(0)
	registration_event := shared.RegistrationEvent{registration.ID, tokenNo}
	reg_event, err := json.Marshal(registration_event)

	if err != nil {
		log.Fatal(err)
		return
	}
}
