package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/ynori7/go-irc/model"
)

type Connection struct {
	Connection       net.Conn
	ConnectionString string
	UseSSL           bool
	Nick             string
}

type MessageHandler func(connection Connection, message model.Message)

/**
 * Establish connection to the server according to the configuration.
 */
func (c *Connection) Connect() (err error) {
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

func (c *Connection) Listen(messageHandler MessageHandler) {
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
func (c *Connection) SendMessage(msg string, to string) {
	fmt.Fprintf(c.Connection, "PRIVMSG %s :%s\r\n", to, msg)
}

/**
 * Join the specified channel
 */
func (c *Connection) JoinChannel(channel string) {
	fmt.Fprintf(c.Connection, "JOIN %s\r\n", channel)
}

/**
 * Respond to server ping
 */
func (c *Connection) Pong(server string) {
	fmt.Fprintf(c.Connection, "PONG %s\r\n", server)
}
