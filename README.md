# Chatter

Chatter is a way for me to experiment with learning Go by building a toy chat
server. It is **in no way** appropriate for production or any sort of public
use.

## How to use

This should work:

```bash
$ go get github.com/mnesvold/chatter
$ ./bin/chatter [--port=8000]
```

Chatter listens on all interfaces and doesn't daemonize: once you get a
message like `Listening on port 8000` on your console, point your favorite
web browser to `localhost:8000`.

## Architecture

Chatter clients (web browsers) communicate with the server over websockets.
The same Go server handles both static file requests (`index.html`, JS, CSS)
and the WS connections. The server's state consists only of the clients'
connections; no chat log is stored, nor any manner of user records.

Data is sent to and from the server as JSON; both directions use the same
schema:

```json
  {
    "nickname": "J. Doe",
    "message": "Hello, world!"
  }
```

## Potential Improvements (Known Flaws)

* User authentication
* Login/logout notifications
* Channels/rooms
* Desktop/audio notifications
* Graceful clientside handling of connection failures
* Integrate with nginx for static files and general reverse proxying (Docker!)
* Serverside validation of client messages
* WSS support
* JS/CSS preprocessors
