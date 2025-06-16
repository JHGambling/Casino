package server

type ResponsePacket struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AuthRegisterPacket struct {
	Username    string `json:"user"`
	DisplayName string `json:"displayName"`
	Password    string `json:"pass"`
}
type AuthRegisterResponsePacket struct {
	ResponsePacket
	Token string `json:"token,omitempty"`
}

type AuthLoginPacket struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}
type AuthLoginResponsePacket struct {
	ResponsePacket
	Token string `json:"token,omitempty"`
}

type AuthAuthenticatePacket struct {
	Token string `json:"token"`
}
type AuthAuthenticateResponsePacket struct {
	ResponsePacket
}

type AuthValidatePacket struct {
	Token string `json:"token"`
}
type AuthValidateResponsePacket struct {
	ResponsePacket
}
