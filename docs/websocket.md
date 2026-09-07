# WebSocket protocol

[Back to README](../README.md) · [HTTP API](http-api.md)

Endpoint: `/api/v1/ws` (`ws://localhost:5050/api/v1/ws` for local development).
HTTP remains the command/query transport. WebSocket delivers server notifications;
it is not an alternative endpoint for sending or editing messages.

## Connection lifecycle

1. Open the WebSocket connection. For browsers, the current handler accepts the
   same origin by default; `HTTP_ALLOWED_ORIGINS` does not configure this policy.
2. Within **5 seconds**, send the first JSON message:

   ```json
   {"type":"authenticate","access_token":"<access JWT>"}
   ```

3. Wait for the acknowledgement:

   ```json
   {"type":"authenticated"}
   ```

4. Receive events. Additional client application messages, including another
   authentication message, close the connection with a policy violation.
5. The connection ends on access-token expiry, application cancellation,
   disconnect, write failure/timeout, or outgoing queue overflow. Refresh the
   token through HTTP and open a new socket; there is no in-place reauthentication.

The initial read limit is 4096 bytes. Each connection has a queue of 64 payloads
and a 5-second write timeout. Slow clients are disconnected instead of blocking
other subscribers indefinitely. The reader handles control frames, but the
application does not schedule periodic heartbeat pings yet.

Do not rely on one close code for every failure: some paths close the connection
without a graceful close handshake. Invalid authentication type/token and extra
application messages have explicit policy-violation closes.

## Events

Unlike HTTP's response envelope, events have `type` and `data` directly.
Participants are resolved after persistence. The sender is included, so their
other tabs/devices receive updates too; all registered connections for each
recipient are offered the event.

### Message created

```json
{
  "type":"message_created",
  "data":{
    "id":"10000000-0000-4000-8000-000000000001",
    "chat_id":"10000000-0000-4000-8000-000000000002",
    "sender_id":"10000000-0000-4000-8000-000000000003",
    "content":"Hello",
    "created_at":"2026-09-07T12:00:00Z",
    "updated_at":null
  }
}
```

`message_edited` has the same fields, with the updated content and `updated_at`.
The WebSocket message DTO does **not** currently include `client_message_id`,
unlike the HTTP DTO. Use the server's message `id` to reconcile a notification
with the HTTP response and avoid displaying the same message twice.

### Message deleted

```json
{
  "type":"message_deleted",
  "data":{
    "id":"10000000-0000-4000-8000-000000000001",
    "chat_id":"10000000-0000-4000-8000-000000000002"
  }
}
```

Only these three event types exist. There are no live events for read markers,
group/membership changes, typing, presence, or account updates. A deleted-message
event does not carry the new chat preview or read markers; query HTTP for current
state when needed.

## Minimal browser example

Run this in a browser page served from the API's origin (the backend does not
itself serve a frontend). Supply a token obtained from HTTP login/register.
This example logs events; it deliberately does not implement a reconnect loop
or application state management.

```js
function connectNotifications(accessToken) {
  const url = new URL("/api/v1/ws", window.location.origin);
  url.protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const socket = new WebSocket(url);

  socket.addEventListener("open", () => {
    socket.send(JSON.stringify({
      type: "authenticate",
      access_token: accessToken,
    }));
  });

  socket.addEventListener("message", ({data}) => {
    const event = JSON.parse(data);
    switch (event.type) {
      case "authenticated":
        console.log("Ready for notifications");
        break;
      case "message_created":
      case "message_edited":
        console.log("Upsert message by id", event.data);
        break;
      case "message_deleted":
        console.log("Remove message by id", event.data.chat_id, event.data.id);
        break;
    }
  });

  socket.addEventListener("close", ({code, reason}) => {
    console.log("Disconnected", code, reason);
    // Obtain a valid access token, reconnect with backoff, and reconcile via HTTP.
  });
  return socket;
}
```

## Delivery guarantees and recovery

- Notifications are best effort. No clients online is not an error; events are
  not retained for offline users. Lookup/encoding/publishing failures are logged
  by the notifier and do not turn committed message operations into failures.
- A process crash between database commit and notification can lose the event.
  There is no broker, outbox, retry worker, acknowledgement, or replay cursor.
- The hub is process-local. Multiple application instances do not distribute
  events to each other's connections.
- Queue insertion preserves its own order, but concurrent publishers do not
  establish a database-commit order. There are no event sequence numbers or
  strict ordering guarantees across concurrent operations.
- A client can receive an event before its own HTTP response. Reconcile by
  message ID and fetch authoritative HTTP state if local state is uncertain.
- After reconnect, HTTP `after=true` can page through newer messages from a
  known position. It does not recover missed edits/deletions. Reload relevant
  history and chat previews; there is no complete offline change feed in v1.
- Logging out or changing a password revokes refresh sessions, but does not
  immediately revoke access JWTs or close sockets. Socket authentication does
  not consult session storage. Do not treat logout as immediate server-side
  disconnection of every device.
