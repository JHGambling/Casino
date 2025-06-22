import { CasinoClient } from "./client";
import { DatabaseOperation } from "./types/db";
import {
    DatabaseOperationPacket,
    DatabaseOperationResponsePacket,
} from "./types/packets";

export type DatabaseOpResult = { result: any; err: any; exec_time_us: number };

export class Database {
    constructor(private client: CasinoClient) {}

    public async performOperation(
        table: string,
        operation: DatabaseOperation,
        op_id: any,
        op_data: any,
    ): Promise<DatabaseOpResult> {
        const response = (
            await this.client.socket.request("db/op", {
                table,
                operation,
                op_id,
                op_data,
            } as DatabaseOperationPacket)
        ).payload as DatabaseOperationResponsePacket;

        // Possible problem in the future (aka i dont care right now):
        // User is not authenticated yet -> operation cannot be executed -> returns just a ResponsePacket,
        // not DatabaseOperationResponsePacket

        return {
            result: response.result,
            err: response.err,
            exec_time_us: response.exec_time_us,
        };
    }
}
