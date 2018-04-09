package queue

import "strings"

func (channel Channel) namespace() (res string) {
	res = strings.Join([]string{
		string(channel.Kind),
		channel.Name,
	}, ".")
	return
}
