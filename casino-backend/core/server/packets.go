package server

type ResponsePacket struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AuthRegisterPacket struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}
type AuthRegisterResponsePacket struct {
	ResponsePacket
	UserAlreadyExists bool   `json:"userAlreadyExists"`
	Token             string `json:"token,omitempty"`
}

type AuthLoginPacket struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type AuthLoginResponsePacket struct {
	ResponsePacket
	UserDoesNotExist bool   `json:"userDoesNotExist"`
	Token            string `json:"token,omitempty"`
}

type AuthAuthenticatePacket struct {
	Token string `json:"token"`
}
type AuthAuthenticateResponsePacket struct {
	ResponsePacket
	UserID    uint  `json:"userID"`
	ExpiresAt int64 `json:"expiresAt"`
}
