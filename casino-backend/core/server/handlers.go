package server

import (
	"errors"
	"fmt"
	"jhgambling/backend/core/data/models"
	"jhgambling/backend/core/utils"
	"time"

	"gorm.io/gorm"
)

func (packet *AuthRegisterPacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	_, err := ctx.Database.GetUserTable().FindByUsername(packet.Username)

	if err == nil {
		// FindByUsername succeeded -> User already exists
		if res, err := BuildPacket("auth/register:res", AuthRegisterResponsePacket{
			ResponsePacket:    ResponsePacket{Success: false, Status: "failed", Message: "User already exists!"},
			UserAlreadyExists: true,
		}, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	} else {
		// FindByUsername failed -> Send error response if the error is not about the record not existing
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			if res, err2 := BuildPacket("auth/register:res", AuthRegisterResponsePacket{
				ResponsePacket:    ResponsePacket{Success: false, Status: "failed", Message: "internal error: " + err.Error()},
				UserAlreadyExists: true,
			}, wsPacket.Nonce); err2 == nil {
				ctx.Client.Send(res)
			}
		}
	}

	// Username is available -> continue

	user := models.UserModel{
		Username:     packet.Username,
		DisplayName:  packet.DisplayName,
		PasswordHash: packet.Password,
		JoinedAt:     time.Now(),
	}

	ctx.Database.GetUserTable().Create(user)
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

func (packet *AuthAuthenticatePacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	valid, userID, expiresAt := ctx.Auth.VerifyToken(packet.Token)

	if valid {
		ctx.Client.Authenticate(userID, expiresAt)
		// Send response
		if res, err := BuildPacket("auth/authenticate:res",
			AuthAuthenticateResponsePacket{
				ResponsePacket: ResponsePacket{Success: true, Status: "ok"},
				UserID:         userID,
				ExpiresAt:      expiresAt.UnixMilli(),
			},
			wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
	} else {
		ctx.Client.RevokeAuthentication()
		// Send response
		if res, err := BuildPacket("auth/authenticate:res",
			AuthAuthenticateResponsePacket{
				ResponsePacket: ResponsePacket{Success: false, Status: "failed", Message: "invalid token"},
			},
			wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
	}
}

func (packet *AuthLoginPacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	user, err := ctx.Database.GetUserTable().FindByUsername(packet.Username)

	if err != nil {
		if res, err := BuildPacket("auth/login:res", AuthLoginResponsePacket{
			ResponsePacket:   ResponsePacket{Success: false, Status: "failed", Message: "User not found!"},
			UserDoesNotExist: true,
		}, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	}

	// Check for correct password
	if user.PasswordHash != packet.Password {
		response := AuthLoginResponsePacket{
			ResponsePacket: ResponsePacket{Success: false, Status: "failed", Message: "Wrong password"},
			WrongPassword:  true,
		}
		if res, err := BuildPacket("auth/login:res", response, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	}

	// Generate authentication token
	token, err := ctx.Auth.CreateTokenForUser(user.ID)

	var response AuthLoginResponsePacket
	if err != nil {
		response = AuthLoginResponsePacket{
			ResponsePacket: ResponsePacket{Success: false, Status: "failed", Message: fmt.Sprintf("internal error: %v", err)},
		}
	} else {
		response = AuthLoginResponsePacket{
			ResponsePacket: ResponsePacket{Success: true, Status: "ok"},
			Token:          token,
		}
	}

	if res, err := BuildPacket("auth/login:res", response, wsPacket.Nonce); err == nil {
		ctx.Client.Send(res)
	}
}

func (packet *DatabaseOperationPacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	if !ctx.Client.IsAuthenticated() {
		ctx.Client.SendUnauthorizedPacket(wsPacket.Nonce)
		return
	}

	start := time.Now()

	userInterface, err := ctx.Database.GetUserTable().FindByID(ctx.Client.authenticatedAs)
	if err != nil {
		response := DatabaseOperationResponsePacket{
			Op:         *packet,
			Result:     nil,
			Error:      err,
			ExecTimeUs: time.Since(start).Microseconds(),
		}
		if res, err := BuildPacket("db/op:res", response, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	}

	// Convert the interface{} to models.UserModel
	userModel, ok := userInterface.(*models.UserModel)
	if !ok {
		response := DatabaseOperationResponsePacket{
			Op:         *packet,
			Result:     nil,
			Error:      "internal error: invalid user model type",
			ExecTimeUs: time.Since(start).Microseconds(),
		}
		if res, err := BuildPacket("db/op:res", response, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	}

	// Pass the concrete UserModel to PerformOperationAsUser
	result, err := ctx.Database.PerformOperationAsUser(*userModel, packet.Table, packet.Operation, packet.OpId, packet.OpData)

	response := DatabaseOperationResponsePacket{
		Op:     *packet,
		Result: result,
		Error:  err,

		ExecTimeUs: time.Since(start).Microseconds(),
	}
	if res, err := BuildPacket("db/op:res", response, wsPacket.Nonce); err == nil {
		ctx.Client.Send(res)
	}
}
