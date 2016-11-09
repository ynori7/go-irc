package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/ynori7/go-irc/model"
)

/*
   Here is basic IRC client. The simplest way to use it is to create a new client and then call Listen, providing
   a message handler. Here's the most basic example:

   conn, err := client.NewConnection("server.whatever.net:6667", false, "myuser")
   if err != nil {
       log.Fatal(err)
   }

   conn.Listen(func(conn client.Client, message model.Message) {
	if message.Type == "PING" {
		conn.Pong(message.Message)
	}
   })

   In this example, the client will connect to the server and respond to any pings from the server.
 */

type MessageHandler func(connection Client, message model.Message)

type Client struct {
	Connection       net.Conn
	ConnectionString string
	UseSSL           bool
	Nick             string
}

func NewConnection(connectionString string, useSSL bool, nick string) (Client, error) {
	conn := Client{
		Nick:             nick,
		ConnectionString: connectionString,
		UseSSL:           useSSL,
	}

	return conn, conn.Connect()
}

/**
 * Establish connection to the server according to the configuration.
 */
func (c *Client) Connect() (err error) {
	if c.UseSSL {
		c.Connection, err = tls.Dial("tcp", c.ConnectionString, &tls.Config{InsecureSkipVerify: true})
	} else {
		c.Connection, err = net.Dial("tcp", c.ConnectionString)
	}

	if err == nil {
		fmt.Fprintf(c.Connection, "USER %s %s %s :%s\r\n", c.Nick, c.Nick, c.Nick, c.Nick)
		fmt.Fprintf(c.Connection, "NICK %s\r\n", c.Nick)
	}

	return err
}

/**
 * Listen to the connection and call the callback function, messageHandler, on any received data
 */
func (c Client) Listen(messageHandler MessageHandler) {
	//Start reading from the connection
	connbuf := bufio.NewReader(c.Connection)
	for {
		str, err := connbuf.ReadString('\n')
		if len(str) > 0 {
			log.Println(str)
			go messageHandler(c, model.NewMessage(str)) //handle message asynchronously so we can go back to listening
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

/**
 * Send the specified message to the specified recipient or channel
 */
func (c *Client) SendMessage(msg string, to string) {
	fmt.Fprintf(c.Connection, "PRIVMSG %s :%s\r\n", to, msg)
}

/**
 * Join the specified channel
 */
func (c *Client) JoinChannel(channel string) {
	fmt.Fprintf(c.Connection, "JOIN %s\r\n", channel)
}

/**
 * Respond to server ping
 */
func (c *Client) Pong(server string) {
	fmt.Fprintf(c.Connection, "PONG %s\r\n", server)
}
