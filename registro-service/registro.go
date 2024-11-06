package registro

import (
	"github.com/marcelogbrito/nats-centromedico/shared"
)

type Server struct {
	*shared.Component
}

// ListenAndServe pega o endereço de rede e porta que o servidor Http deve vincular e inicia
func (s *Server) ListenAndServe(addr string) error {

}
