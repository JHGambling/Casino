export type ResponsePacket = {
    success: boolean;
    status: string;
    message: string;
};

export type AuthRegisterPacket = {
    username: string;
    displayName: string;
    password: string;
};

export type AuthRegisterResponsePacket = ResponsePacket & {
    userAlreadyExists: boolean;
    token?: string;
};

export type AuthLoginPacket = {
    username: string;
    password: string;
};

export type AuthLoginResponsePacket = ResponsePacket & {
    userDoesNotExist: boolean;
    token?: string;
};

export type AuthAuthenticatePacket = {
    token: string;
};

export type AuthAuthenticateResponsePacket = ResponsePacket & {
    userID: number;
    expiresAt: number;
};

export type DatabaseOperationPacket = {
    operation: string;
    table: string;
    op_id: any;
    op_data: any;
};

export type DatabaseOperationResponsePacket = {
    op: DatabaseOperationPacket;
    result: any;
    err: any;
    exec_time_us: number;
};
