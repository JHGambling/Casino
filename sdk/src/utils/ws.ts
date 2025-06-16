import { WebsocketPacket } from "../types/packet";

/**
 * Creates a formatted packet to send to the server
 * @param type The packet type
 * @param payload The packet payload
 * @param nonce Optional nonce value (will be generated if not provided)
 * @returns The formatted packet
 */
export function createPacket(
    type: string,
    payload: any = {},
    nonce?: number,
): WebsocketPacket {
    return {
        type,
        payload,
        nonce: nonce !== undefined ? nonce : Date.now(),
    };
}

/**
 * Parse a JSON string into a WebsocketPacket
 * @param data The JSON string to parse
 * @returns The parsed packet
 */
export function parsePacket(data: string): WebsocketPacket {
    try {
        return JSON.parse(data) as WebsocketPacket;
    } catch (error) {
        throw new Error("Failed to parse packet: Invalid JSON");
    }
}
