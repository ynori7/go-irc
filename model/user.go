package model

import "regexp"

type User struct {
	Nick     string
	Username string
	Host     string
	Raw      string
}

/**
 * Parses the user string and returns a new IrcUser object.
 * Example string:
 * ynori7!~ynori7@unaffiliated/ynori7
 */
func NewUser(userString string) User {
	ircUser := User{Raw: userString}

	re, err := regexp.Compile(`(.*)!(.*)@(.*)`)

	if err == nil {
		res := re.FindStringSubmatch(userString)

		if len(res) == 4 {
			ircUser.Nick = res[1]
			ircUser.Username = res[2]
			ircUser.Host = res[3]
		}
	}

	return ircUser
}
