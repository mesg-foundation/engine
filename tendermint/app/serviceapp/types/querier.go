package types

type QueryService struct {
	Hash       string `json:"hash"`
	Definition string `json:"definition"`
}

// Query Result Payload for a resolve servcie.
type QueryServiceResolve struct {
	Service QueryService `json:"service"`
}

// Query Result Payload for a resolve servcies.
type QueryServicesResolve struct {
	Services []QueryService `json:"services"`
}
