package model

import "strings"

type Message struct {
	Type     string
	Sender   User
	Location string
	Message  string
	Raw      string
}

/**
 * Parses the message received from the server and returns a new IrcMessage object.
 * Samples Messages:
:ynori7!~ynori7@unaffiliated/ynori7 KICK #ynori7 blorgleflorps :blorgleflorps
:blorgleflorps!~blorglefl@2001:4c50:29e:2c00:9084:4b28:8dbd:791 JOIN #ynori7
:wolfe.freenode.net 353 blorgleflorps @ #ynori7 :blorgleflorps @ynori7
:wolfe.freenode.net 366 blorgleflorps #ynori7 :End of /NAMES list.
:ynori7!~ynori7@unaffiliated/ynori7 PRIVMSG #ynori7 :hello blorgleflorps
*/
func NewMessage(msg string) Message {
	ircMsg := Message{Raw: msg}

	if strings.HasPrefix(msg, "PING") {
		ircMsg.Type = "PING"
		ircMsg.Message = strings.Fields(msg)[1]
	} else {
		if strings.HasPrefix(msg, ":") {
			msg = msg[1:]
		}

		tmp := strings.Fields(msg)
		ircMsg.Sender = NewUser(tmp[0])
		ircMsg.Type = tmp[1]

		//For JOIN there's a : in front
		if strings.HasPrefix(tmp[2], ":") {
			tmp[2] = tmp[2][1:]
		}
		ircMsg.Location = tmp[2]

		if ircMsg.Type == "KICK" && len(tmp) >= 3 { //for KICK it ends with "username :"
			ircMsg.Message = tmp[3]
		} else if len(tmp) >= 3 && strings.Contains(msg, ":") {
			ircMsg.Message = strings.TrimSpace(strings.SplitAfterN(msg, ":", 2)[1])
		}
	}

	return ircMsg
}
