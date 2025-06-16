package server

import (
	"fmt"
	"jhgambling/backend/core/data/models"
	"jhgambling/backend/core/utils"
	"time"
)

func (packet *AuthRegisterPacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	_, exists, _ := ctx.Database.GetUserByUsername(packet.Username)

	if !exists {
		if res, err := BuildPacket("auth/register:res", AuthRegisterResponsePacket{
			ResponsePacket:    ResponsePacket{Success: false, Status: "failed", Message: "User already exists!"},
			UserAlreadyExists: true,
		}, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	}

	// Username is available -> continue

	user := models.UserModel{
		Username:     packet.Username,
		DisplayName:  packet.DisplayName,
		PasswordHash: packet.Password,
		JoinedAt:     time.Now(),
	}

	ctx.Database.CreateUser(&user)
	utils.Log("ok", "casino::gateway",
		"new user registered: ",
		fmt.Sprintf(
			"User<id: %v | username: %s | display: \"%s\">",
			user.ID,
			user.Username,
			user.DisplayName,
		),
	)

	// Generate authentication token
	token, err := ctx.Auth.CreateTokenForUser(user.ID)

	var response AuthRegisterResponsePacket
	if err != nil {
		response = AuthRegisterResponsePacket{
			ResponsePacket: ResponsePacket{Success: false, Status: "failed", Message: fmt.Sprintf("internal error: %v", err)},
		}
	} else {
		response = AuthRegisterResponsePacket{
			ResponsePacket: ResponsePacket{Success: true, Status: "ok"},
			Token:          token,
		}
	}

	if res, err := BuildPacket("auth/register:res", response, wsPacket.Nonce); err == nil {
		ctx.Client.Send(res)
	}
}
