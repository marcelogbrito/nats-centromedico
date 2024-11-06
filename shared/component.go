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
