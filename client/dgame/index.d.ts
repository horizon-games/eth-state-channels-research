import * as ethers from 'ethers';
export declare class DGame {
    private account;
    constructor(gameAddress: string, account?: ethers.Wallet | undefined);
    readonly address: string;
    deposit(value: ethers.utils.BigNumber): Promise<void>;
    readonly matchDuration: Promise<number>;
    isSecretSeedValid(address: string, secretSeed: Uint8Array): Promise<boolean>;
    createMatch(secretSeed: Uint8Array, onChange?: ChangeCallback, onCommit?: CommitCallback): Promise<Match>;
    private signer?;
    private arcadeumContract;
    private gameContract;
}
export interface ChangeCallback {
    (match: Match, previousState: State, currentState: State, aMove: Move, anotherMove?: Move): void;
}
export interface CommitCallback {
    (match: Match, previousState: State, move: Move): void;
}
export interface Match {
    readonly playerID: number;
    readonly state: Promise<State>;
    commit(move: Move): Promise<void>;
}
export interface State {
    readonly winner: Promise<Winner>;
    readonly nextPlayers: Promise<NextPlayers>;
    isMoveLegal(move: Move): Promise<{
        isLegal: boolean;
        reason: number;
    }>;
    nextState(aMove: Move, anotherMove?: Move): Promise<State>;
    nextState(moves: [Move] | [Move, Move]): Promise<State>;
    readonly encoding: any;
    readonly hash: Promise<ethers.utils.BigNumber>;
}
export declare enum Winner {
    None = 0,
    Player0 = 1,
    Player1 = 2,
}
export declare enum NextPlayers {
    None = 0,
    Player0 = 1,
    Player1 = 2,
    Both = 3,
}
export declare class Move {
    readonly move: {
        playerID: number;
        data: Uint8Array;
        signature?: any;
    };
    constructor(move: {
        playerID: number;
        data: Uint8Array;
        signature?: any;
    });
    sign(subkey: ethers.Wallet, state: State): Promise<void>;
    readonly playerID: number;
    readonly data: Uint8Array;
    private signature?;
}
