package shared

// RegsitrationRequest contem dados sobre o paciente
type RegistrationRequest struct {

	//Nome completo do paciente
	NomeCompleto string `json:"nome_completo,omitempty"`

	// endereço do paciente
	Endereco string `json:"endereco,omitempty"`

	//Numero SUS do paciente
	ID int `json:"id"`

	// Orientaçao sexual
	sexo string `json:"sexo,omitempty"`

	// Endereço de email
	Email string `json:"email,omitempty"`

	// Numero de telefone
	Telefone string `json:"telefone,omitempty"`

	// Outros detalhes
	Observacoes string `json:"observacoes,omitempty"`

	//REquestID é o ID do request
	RequestID string `json:"request_id,omitempty"`
}

// RegistrationEvent contem os detalhes de uma instancia de registration
type RegistrationEvent struct {

	// ID do paciente
	ID int `json:"id"`

	//token do paciente
	Token uint64 `json:"token"`
}
