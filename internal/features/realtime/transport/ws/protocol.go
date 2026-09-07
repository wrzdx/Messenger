package realtime_transport_ws

const (
	authReqType = "authenticate"
	authResType = "authenticated"
)

type AuthRequest struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type AuthResponse struct {
	Type string `json:"type"`
}

// Event is the server-to-client envelope. Data holds a transport DTO and is
// serialized once by Publish before being shared across connection queues.
type Event struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}
