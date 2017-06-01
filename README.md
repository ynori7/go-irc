# Go IRC

Go IRC is a very basic IRC client which can connect to a server (supports SSL connections), parse messages, and listen for input from the server. 

The simplest way to use it is to create a new `Client` and then call the `Listen` method, providing a custom `MessageHandler`. Here's the most basic example:

```
conn := client.NewConnection("server.whatever.net:6667", false, "myuser")

conn.Listen(func(conn client.Client, message model.Message) {
    if message.Type == "PING" {
        conn.Pong(message.Message)
    }
})
```

The provided `MessageHandler` will be called (asynchronously) any time the client receives new data from the server. In this example, the client will connect to the server and respond to any pings from the server.
