import { Observable } from 'rxjs/Rx'
import { Subscriber } from 'rxjs/Subscriber'

export class Meta {
  constructor(public matchID: number, public index: number, public code: number) { }
}

export class SessionInitMessage {
  constructor(public signature: string) { }
}

export class Message {
  constructor(public meta: Meta, public payload: string) { }
}

export class Signature {
  constructor(
    public v: number = 0, // uint8 
    public r: string = '', // base64 []byte value
    public s: string = '') {} // base64 []byte value
}

export class Token {
  // subkey, seed, signature of seed, game ID
  constructor(
    public gameID: number, 
    public subkey: string, // public address of subkey, e.g., "0x" 
    public signature: Signature, 
    public seed: string) { } // game "deck" as base64 []byte value
}

export interface ConnectHandler {
  onError: (obv: Subscriber<Message>) => (ev: Event) => {}
  onMessage: (obv: Subscriber<Message>) => (ev: MessageEvent) => {}
  onOpen: (obv: Subscriber<Message>) => (ev: Event) => {}
}

export class Relay {
  private host: string
  private port: number
  private seed: string
  private signature: Signature
  private subkey: string
  private gameID: number
  private ws: WebSocket
  private matchID: number
  private index: number // player ID
  private ssl: boolean

  constructor(
    host: string,
    port: number,
    ssl = true,
    seed = '',
    signature = new Signature(),
    subkey = '',
    gameID = -1) {
    this.host = host
    this.port = port
    this.ssl = ssl
    this.seed = seed
    this.subkey = subkey
    this.signature = signature
    this.gameID = gameID
  }

  token(): string {
    const token = JSON.stringify(new Token(this.gameID, this.subkey, this.signature, this.seed))
    return new Buffer(token).toString('base64')
  }

  connect(callbacks?: ConnectHandler): Observable<Message> {
    this.ws = new WebSocket(`${this.ssl ? 'wss' : 'ws'}://${this.host}${this.port === 80 ? '' : `:${this.port}`}/ws?token=${this.token()}`)
    return Observable.create((obv: Subscriber<Message>) => {
      this.ws.onopen    = callbacks && callbacks.onOpen != null ? callbacks.onOpen(obv) : this.onOpen(obv)
      this.ws.onmessage = callbacks && callbacks.onMessage != null ? callbacks.onMessage(obv) : this.onMessage(obv)
      this.ws.onerror   = callbacks && callbacks.onError != null ? callbacks.onError(obv) : this.onError(obv)
    })
  }

  isInitialized(): boolean {
    return typeof this.ws !== 'undefined'
  }

  send(json: string, code = 0) {
    console.log(`wsrelay: sending ${json}`)
    if (this.isInitialized()) {
      const message = new Message(new Meta(this.matchID, this.index, code), json)
      this.ws.send(JSON.stringify(message))
    }
  }

  private newError(msg: string) {
    return new Message(new Meta(this.matchID, this.index, -1), msg)
  }

  private onError(obv: Subscriber<Message>) {
    return (event: Event) => {
      console.log('error: ' + event)
      obv.error(this.newError('Error receiving message.'))
    }
  }

  private onOpen(obv: Subscriber<Message>) {
    return (event: Event) => {
      console.log('open: ' + event)
    }
  }

  private onMessage(obv: Subscriber<Message>) {
    return (msg: MessageEvent) => {
      try {
        console.log(`msg.data: ${msg.data}`)
        const data = <Message>JSON.parse(msg.data)
        if (data.meta.code === -1) {
          obv.error(data)
          return
        }
        console.log(`Relay message received: ${JSON.stringify(data)}`)
        if (data.meta.code === 1) { // cache session info
          this.matchID = data.meta.matchID
          this.index = data.meta.index
        }
        obv.next(data)
      } catch (e) {
        console.log('Error parsing message.' + e)
        obv.error(this.newError('Error parsing message.'))
      }
    }
  }
}