import { Observable, Observer } from 'rxjs/Rx'
import { Subscriber } from 'rxjs/Subscriber'
import { first, publishReplay, refCount } from 'rxjs/operators'

export class Meta {
  constructor(public subkey: string, public index: number, public code: number) { }
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
    public s: string = '') { } // base64 []byte value
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
  private seed: string
  private signature: Signature
  private subkey: string
  private gameID: number
  private ws: WebSocket
  private index: number // player ID
  private stream: Observable<Message>

  constructor(
    private url: string,
    seed = '',
    signature = new Signature(),
    subkey = '',
    gameID = -1) {

    if (this.url.endsWith('/')) {
      this.url = this.url.slice(0, -1)
    }

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
    this.ws = new WebSocket(`${this.url}/ws?token=${this.token()}`)
    return Observable.create((obv: Subscriber<Message>) => {
      this.ws.onopen = callbacks && callbacks.onOpen != null ? callbacks.onOpen(obv) : this.onOpen(obv)
      this.ws.onmessage = callbacks && callbacks.onMessage != null ? callbacks.onMessage(obv) : this.onMessage(obv)
      this.ws.onerror = callbacks && callbacks.onError != null ? callbacks.onError(obv) : this.onError(obv)
    })
  }

  setStream(callbacks?: ConnectHandler) {
    if (!this.isInitialized()) {
      this.stream = this.connect(callbacks).publishReplay(1).refCount()
    }
  }

  connectForTimestamp(callbacks?: ConnectHandler): Promise<Message> {
    return new Promise((resolve, reject) => {
      this.setStream(callbacks)
      this.stream.pipe(first()).subscribe(msg => {
        resolve(msg)
      }, err => {
        reject(err)
      })
    })
  }

  connectForMatchVerified(callbacks?: ConnectHandler): Promise<Message> {
    return new Promise((resolve, reject) => {
      this.setStream(callbacks)
      this.stream.skip(1).take(1).subscribe(msg => {
        resolve(msg)
      }, err => {
        reject(err)
      })
    })
  }

  subscribe(observer: Observer<Message>): void {
    this.setStream()
    this.stream.skip(2).subscribe(observer)
  }

  isInitialized(): boolean {
    return typeof this.ws !== 'undefined'
  }

  send(json: string, code = 0) {
    if (this.isInitialized()) {
      const message = new Message(new Meta(this.subkey, this.index, code), json)
      this.ws.send(JSON.stringify(message))
    }
  }

  private newError(msg: string) {
    return new Message(new Meta(this.subkey, this.index, -1), msg)
  }

  private onError(obv: Subscriber<Message>) {
    return (event: Event) => {
      obv.error(this.newError('Error receiving message.'))
    }
  }

  private onOpen(obv: Subscriber<Message>) {
    return (event: Event) => {
    }
  }

  private onMessage(obv: Subscriber<Message>) {
    return (msg: MessageEvent) => {
      try {
        const data = <Message>JSON.parse(msg.data)
        if (data.meta.code === -1) {
          obv.error(data)
          return
        }
        if (data.meta.code === 1) { // cache session info
          this.index = data.meta.index
        }
        obv.next(data)
      } catch (e) {
        obv.error(this.newError('Error parsing message.'))
      }
    }
  }
}
