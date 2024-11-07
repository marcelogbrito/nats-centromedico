package shared

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

// component contem logica reusavel relacionada aos
// processos de conexao ao NATS e Banco de Dados no sistema
type Component struct {
	// cmu é o bloqueio do componente
	cmu sync.Mutex

	// id é o identificador unico usado para este component
	id string

	// nc é a conexão do NATS
	nc *nats.Conn

	// db é a conexão ao banco de dados
	db *sql.DB

	// kind é o tipo de component
	kind string
}

// componente retorna a atual conexao do banco de dados
func (c *Component) DB() *sql.DB {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.db
}

// NATS retorna a conexão NATS atual
func (c *Component) NATS() *nats.Conn {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.nc
}

// cria um novo Componente
func NewComponent(kind string) *Component {
	id := nuid.Next()
	return &Component{
		id:   id,
		kind: kind,
	}
}

// Name é o label utilizado para identificar a conexao NATS
func (c *Component) Name() string {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return fmt.Sprintf("%s:%s", c.kind, c.id)
}

// SetupConnectionToNATS conecta ao NATS e registra o evento
// e torna disponivel para requests de descoberta tambem
func (c *Component) SetupConnectionToNATS(servers string, options ...nats.Option) error {
	// configura a conexao com kind e id para o componente
	options = append(options, nats.Name(c.Name()))

	c.cmu.Lock()
	defer c.cmu.Unlock()

	//Conect ao NATS com opçoes customizadas
	nc, err := nats.Connect(servers, options...)
	if err != nil {
		return err
	}
	c.nc = nc

	// Configura callback de eventos NATS
	// Lida com erros de protocolo e casos de consumo lento
	nc.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		log.Printf("Erro NATS: %s\n", err)
	})
	nc.SetReconnectHandler(func(_ *nats.Conn) {
		log.Println("Reconectado ao NATS!")
	})
	nc.SetDisconnectHandler(func(_ *nats.Conn) {
		log.Println("Desconectado do NATS!")
	})
	nc.SetClosedHandler(func(_ *nats.Conn) {
		panic("Conexão ao NATS fechada!")
	})

	return err
}

// SetupConectToDB cria a conexão para o banco de dados
func (c *Component) SetupConnectionToDB(dbDriver string, connectionString string) error {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	db, err := sql.Open(dbDriver, connectionString)
	if err != nil {
		panic(err.Error())
	}
	c.db = db
	return err
}
