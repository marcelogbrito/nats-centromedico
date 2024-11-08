package registro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

var ops uint64

// generateTokenNumber gera um numero de token para o paciente
func generateTokenNumber(start uint64) uint64 {
	if start > 0 {
		ops = start
		return ops
	}
	atomic.AddUint64(&ops, 1)
	return ops
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
	router.HandleFunc("/registro/", s.HandleRegistro).Methods("POST")

	// Handle de request de update de paciente
	// PUT /cmed/paciente/atualiza
	router.HandleFunc("atualiza", s.HandleAtualiza).Methods("PUT")

	//Handle rquest de view
	// GET /cmed/paciente/view/{id}
	router.HandleFunc("/view/{id}", s.HandleView).Methods("GET")

	// Handle request de token
	// GET /cmed/paciente/token
	router.HandleFunc("token/{id}", s.HandleToken).Methods("GET")

	//Handle reset de request de token
	// GET /cmed/paciente/token/reset/{id}
	router.HandleFunc("token/reset/{id}", s.HandleTokenReset).Methods("GET")

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go srv.Serve(l)
	return nil

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
	err = json.Unmarshal(body, &registration)
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

	log.Printf("requestID:%s - Publicando evento de registro com pacienteID %d\n", &registration.RequestID, registration.ID)
	//publicando a mensagem no servidor NATS
	nc.Publish("paciente.registro", reg_event)
	json.NewEncoder(w).Encode(registration_event)
}

func (s *Server) HandleAtualiza(w http.ResponseWriter, r *http.Request) {
	//patientID := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request *shared.RegistrationRequest
	err = json.Unmarshal(body, &request)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	db := s.DB()

	insForm, err :=
		db.Prepare("UPDATE paciente_detalhes SET nome_completo=?, endereco=?, sexo=?, telefone=?, observacoes=? WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	insForm.Exec(request.NomeCompleto, request.Endereco, request.Sexo, request.Telefone, request.Observacoes, request.ID)
	json.NewEncoder(w).Encode("Registro de Paciente Atualizado com Sucesso")
}

// HandleView processa requests para visualizaçao de dados de paciente
func (s *Server) HandleView(w http.ResponseWriter, r *http.Request) {
	pacienteID := mux.Vars(r)["ID"]
	// Inserir dados ao banco de dados
	db := s.DB()

	selDB, err := db.Query("SELECT * FROM paciente_detalhes WHERE ID=?", pacienteID)
	if err != nil {
		panic(err.Error())
	}

	registration := shared.RegistrationRequest{}
	for selDB.Next() {
		var id int
		var nome_completo, endereco, telefone, sexo, observacoes string
		err = selDB.Scan(&id, &nome_completo, &endereco, &sexo, &telefone, &observacoes)
		if err != nil {
			panic(err.Error())
		}
		registration.ID = id
		registration.NomeCompleto = nome_completo
		registration.Endereco = endereco
		registration.Sexo = sexo
		registration.Telefone = telefone
		registration.Observacoes = observacoes
	}

	fmt.Println(registration)
	json.NewEncoder(w).Encode(registration)
}

// HandleToken processa requests de geraçao de token para pacientes registrados
func (s *Server) HandleToken(w http.ResponseWriter, r *http.Request) {
	token := generateTokenNumber(0)
	pacienteID, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println("Token %d gerado para usuaário %d", token, pacienteID)
	//Publica evento no servidor NATS
	nc := s.NATS()

	registration_event := shared.RegistrationEvent{pacienteID, token}
	reg_event, err := json.Marshal(registration_event)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("TokenID: %d - Publicando evento de registro com paciente %d\n", token, pacienteID)

	// Publicando mensagem no servidor NATS
	nc.Publish("paciente.registro", reg_event)
	json.NewEncoder(w).Encode(registration_event)
}

// HandleTokenReset processa requests de reset de token
func (s *Server) HandleTokenReset(w http.ResponseWriter, r *http.Request) {
	resetID, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	generateTokenNumber(resetID)
	json.NewEncoder(w).Encode("Token resetado com sucesso")
}
