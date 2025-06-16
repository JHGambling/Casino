package server

import "fmt"

func (packet *AuthRegisterPacket) Handle(wsPacket WebsocketPacket, ctx *HandlerContext) {
	fmt.Println(packet.DisplayName, packet.Password, packet.Username)
}
