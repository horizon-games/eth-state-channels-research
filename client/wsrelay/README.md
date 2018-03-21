## Websocket Client to the Arcadeum Server

Use the `Relay` class to interface with the Arcadeum Server via a websocket connection.

Example:

```$javascript
// Connect to the relay server. On connection, it will
// try and find a matching opponent
const relay = new Relay('localhost', 8080)
relay.connect().subscribe(
  (message: string) => {
    console.log(`Received message: ${message}`)
  },
  (err: string) => {
    console.log(`Error: ${err}`)
  })
  
// send a string to your opponent  
const json = '{"id": "abc123"}'  
relay.send(json) // fire and forget
  
```