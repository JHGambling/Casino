package server

type ResponsePacket struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// User registration
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

// User login
type AuthLoginPacket struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type AuthLoginResponsePacket struct {
	ResponsePacket
	UserDoesNotExist bool   `json:"userDoesNotExist"`
	WrongPassword    bool   `json:"wrongPassword"`
	Token            string `json:"token,omitempty"`
}

// User authenticate
type AuthAuthenticatePacket struct {
	Token      string `json:"token"`
	ClientType string `json:"clientType"`
}
type AuthAuthenticateResponsePacket struct {
	ResponsePacket
	UserID    uint  `json:"userID"`
	ExpiresAt int64 `json:"expiresAt"`
}

// Does User exist
type DoesUserExistPacket struct {
	Username string `json:"username"`
}

type DoesUserExistResponsePacket struct {
	ResponsePacket
	UserExists bool `json:"userExists"`
}

// Database operation
type DatabaseOperationPacket struct {
	Operation string      `json:"operation"`
	Table     string      `json:"table"`
	OpId      interface{} `json:"op_id"`
	OpData    interface{} `json:"op_data"`
}

type DatabaseOperationResponsePacket struct {
	Op         DatabaseOperationPacket `json:"op"`
	Result     interface{}             `json:"result"`
	Error      interface{}             `json:"err"`
	ExecTimeUs int64                   `json:"exec_time_us"`
}

// Database subscribe
type DatabaseSubscribePacket struct {
	Operation  string `json:"operation"`
	TableID    string `json:"tableID"`
	ResourceID uint   `json:"resourceID"`
}

type DatabaseSubUpdatePacket struct {
	TableID    string      `json:"tableID"`
	ResourceID interface{} `json:"resourceID"`

	Operation string      `json:"op"`
	Data      interface{} `json:"data"`
}

// Session
type SetSessionPacket struct {
	SessionID uint `json:"sessionID"`
}

type GameFinishedLoadingPacket struct {
	SessionID uint `json:"sessionID"`
}
