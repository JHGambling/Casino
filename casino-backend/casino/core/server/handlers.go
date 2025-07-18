package server

import (
	"errors"
	"fmt"
	"jhgambling/backend/core/utils"
	"jhgambling/protocol/models"
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

	// Hash the password before storing
	hash, err := ctx.Auth.HashPassword(packet.Password)
	if err != nil {
		if res, err := BuildPacket("auth/register:res", AuthRegisterResponsePacket{
			ResponsePacket: ResponsePacket{Success: false, Status: "failed", Message: "Error hashing password!"},
		}, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	}

	// Create the User record
	user := &models.UserModel{
		Username:     packet.Username,
		DisplayName:  packet.DisplayName,
		PasswordHash: hash,
		JoinedAt:     time.Now(),
	}

	// Create wallet with starting bonus
	wallet := &models.WalletModel{
		NetworthCents:         100000, // $1000 in cents
		ReceivedStartingBonus: true,
	}
	user.Wallet = *wallet

	err = ctx.Database.GetUserTable().Create(user)
	if err != nil {
		if res, err := BuildPacket("auth/register:res", AuthRegisterResponsePacket{
			ResponsePacket: ResponsePacket{Success: false, Status: "failed", Message: "error creating user entry: " + err.Error()},
		}, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	}
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
		utils.Log("debug", "casino::gateway", "[Auth] user ", user.ID, " with username '", user.Username, "' has registered")
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
		ctx.Client.clientType = packet.ClientType
		utils.Log("debug", "casino::gateway", "[Auth] user ", userID, " has been authenticated with type '", packet.ClientType, "'")
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
		utils.Log("debug", "casino::gateway", "[Auth] client failed authentication due to invalid token")
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
	if !ctx.Auth.CheckPasswordHash(packet.Password, user.PasswordHash) {
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
		utils.Log("debug", "casino::gateway", "[Auth] user ", user.ID, " has logged in")
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
		utils.Log("warn", "casino::gateway", "[db/op] error getting user:", err)
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

func (packet *DoesUserExistPacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	result, err := ctx.Database.GetUserTable().FindByUsername(packet.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response := DoesUserExistResponsePacket{
			ResponsePacket: ResponsePacket{Success: false, Status: "failed", Message: "internal error: " + err.Error()},
			UserExists:     false,
		}
		if res, err := BuildPacket("auth/does_user_exist:res", response, wsPacket.Nonce); err == nil {
			ctx.Client.Send(res)
		}
		return
	}

	response := DoesUserExistResponsePacket{
		ResponsePacket: ResponsePacket{Success: true, Status: "ok"},
		UserExists:     result != nil && err == nil,
	}
	if res, err := BuildPacket("auth/does_user_exist:res", response, wsPacket.Nonce); err == nil {
		ctx.Client.Send(res)
	}
}

func (packet *DatabaseSubscribePacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	if !ctx.Client.IsAuthenticated() {
		ctx.Client.SendUnauthorizedPacket(wsPacket.Nonce)
		return
	}

	utils.Log("debug", "casino::gateway", "[db/sub] user:", ctx.Client.authenticatedAs, " op:'", packet.Operation, "' table:", packet.TableID, " resource:", packet.ResourceID)

	if packet.Operation == "subscribe" {
		ctx.Client.Subscriptions = append(ctx.Client.Subscriptions, DBSubscription{
			TableID:    packet.TableID,
			ResourceID: packet.ResourceID,
		})
	} else if packet.Operation == "unsubscribe" {
		for i, sub := range ctx.Client.Subscriptions {
			// Remove matching subscriptions
			// or Remove all subscriptions (if tableID is empty)
			if sub.TableID == packet.TableID && sub.ResourceID == packet.ResourceID || len(packet.TableID) == 0 {
				ctx.Client.Subscriptions = append(ctx.Client.Subscriptions[:i], ctx.Client.Subscriptions[i+1:]...)
				break
			}
		}
	} else {
		utils.Log("warn", "casino::gateway", "[db/sub] user ", ctx.Client.authenticatedAs, " tried to perform unkown db/sub operation: ", packet.Operation)
	}
}

func (packet *SetSessionPacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	ctx.Client.SetSession(packet.SessionID)
}

func (packet *GameFinishedLoadingPacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	for _, c := range ctx.Gateway.Clients {
		if c.GetSession() == packet.SessionID {
			// Client is part of the same session, so we can send the
			// finished loading packet to the client
			if res, err := BuildPacket("game/finished_loading", packet, wsPacket.Nonce); err == nil {
				c.Send(res)
			}
		}
	}
}
