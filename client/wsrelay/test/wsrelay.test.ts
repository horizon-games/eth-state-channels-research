import { Relay, Message, Signature } from '../index'
import { Server } from 'mock-socket'
import { setTimeout } from 'timers'

describe('Relay', () => {
  const host = 'localhost'
  const port = 8888
  const matchID = 1 
  const message = 'foo'
  const signature = new Signature(92, '03PpsnqhjDKZxvhbPTeTnZgJI2O814WLcdKK7kTn3g==', 'ls8hxypZp31Ubbget5DM4P09x1b2FL61FnfwUkCqH9g=')
  const gameID = 1
  const playerIdx: number = 0
  const subkey = '0x5409ed021d9299bf6814279a6a1411a7e866a631'
  const seed = '[1,2]'
  let mockServer: Server

  beforeEach(() => {
    const token = new Relay(host, port, true, seed, signature, subkey, gameID).token()
    mockServer = new Server(`wss://${host}:${port}/ws?token=${token}`)
  })

  afterEach(() => {
    mockServer.stop()
  })

  test('should return UUID on initial socket connection', async (done) => {
    mockServer.on('connection', server => {
      mockServer.send(`{"meta": {"matchID": ${matchID}, "index": ${playerIdx}, "code": 1}, "payload": "{\\"signature\\":\\"signedmessage\\"}"}`)
    })
    const relay = new Relay(host, port, true, seed, signature, subkey, gameID)
    const result: Message[] = []
    relay.connect().subscribe((val: Message) => {
      result.push(val)
    })

    setTimeout(() => {
      expect(relay.isInitialized()).toBe(true)
      expect(result.length).toBe(1)
      expect(result[0].meta.matchID).toBe(matchID)
      expect(result[0].meta.index).toBe(playerIdx)
      done()
    }, 100)
  })

  test('should send payload to server', async (done) => {
    const payload = 'some-random-payload'
    mockServer.on('connection', server => {
      mockServer.send(`{"meta": {"matchID": "${matchID}", "playerIdx": ${playerIdx}, "code": 1}, "payload": "{\\"playerIndices\\": [0, 1]}"}`)
    })
    mockServer.on('message', server => {
      mockServer.send(`{"meta": {"matchID": "${matchID}", "playerIdx": ${playerIdx}, "code": 0}, "payload": "${payload}"}`)
    })
      const relay = new Relay(host, port, true, seed, signature, subkey, gameID)
    const result: Message[] = []
    relay.connect().subscribe((val: Message) => {
      result.push(val)
    })
    relay.send(payload)

    setTimeout(() => {
      expect(relay.isInitialized()).toBe(true)
      expect(result.length).toBe(2)
      expect(result[1].payload).toBe(payload)
      done()
    }, 100)
  })

  test('should return messages', async (done) => {
    const payload1 = 'something1'
    const payload2 = 'something2'
    mockServer.on('connection', server => {
      mockServer.send(`{"meta": {"matchID": "${matchID}", "playerIdx": ${playerIdx}, "code": 0}, "payload": "${payload1}"}`)
      mockServer.send(`{"meta": {"matchID": "${matchID}", "playerIdx": ${playerIdx}, "code": 0}, "payload": "${payload2}"}`)
    })
    const relay = new Relay(host, port, true, seed, signature, subkey, gameID)
    const result: Message[] = []
    relay.connect().subscribe((val: Message) => {
      result.push(val)
    })

    setTimeout(() => {
      expect(result.length).toBe(2)
      expect(result[0].payload).toBe(payload1)
      expect(result[1].payload).toBe(payload2)
      done()
    }, 100)
  })

  test('should return error if message is not parseable', async (done) => {
    mockServer.on('connection', server => {
      mockServer.send(`{lksjdfsdf}`) // unparseable JSON
    })
    const relay = new Relay(host, port, true, seed, signature, subkey, gameID)
    const result: Message[] = []
    const errs: Message[] = []
    relay.connect().subscribe((val: Message) => {
      result.push(val)
    }, (err) => {
      errs.push(err)
    })

    setTimeout(() => {
      expect(result.length).toBe(0)
      expect(errs.length).toBe(1)
      expect(errs[0].payload).toBe('Error parsing message.')
      done()
    }, 100)
  })

  test('should return error if server returns error message', async (done) => {
    mockServer.on('connection', server => {
      mockServer.send(`{"meta": {"matchID": "${matchID}", "playerIdx": ${playerIdx}, "code": -1}, "payload": "Uh oh!"}`)
    })
    const relay = new Relay(host, port, true, seed, signature, subkey, gameID)
    const result: Message[] = []
    const errs: Message[] = []
    relay.connect().subscribe((val: Message) => {
      result.push(val)
    }, (err) => {
      errs.push(err)
    })

    setTimeout(() => {
      expect(result.length).toBe(0)
      expect(errs.length).toBe(1)
      expect(errs[0].payload).toBe('Uh oh!')
      done()
    }, 100)
  })
})
