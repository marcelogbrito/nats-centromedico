package shared

import (
	"database/sql"
	"sync"

	"github.com/nats-io/nats.go"
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
