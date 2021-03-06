package main

import (
	"github.com/0mili/mili"
	"github.com/0mili/redis-memory"
	"github.com/pkg/errors"
)

type ExampleBot struct {
	*mili.Bot
}

func main() {
	b := &ExampleBot{
		Bot: mili.New("example",
			redis.Memory("localhost:6379"),
		),
	}

	b.Respond("remember (.+) is (.+)", b.Remember)
	b.Respond("what is (.+)", b.WhatIs)
	b.Respond("show keys", b.ShowKeys)

	err := b.Run()
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
}

func (b *ExampleBot) Remember(msg mili.Message) error {
	key, value := msg.Matches[0], msg.Matches[1]
	msg.Respond("OK, I'll remember %s is %s", key, value)
	return b.Store.Set(key, value)
}

func (b *ExampleBot) WhatIs(msg mili.Message) error {
	key := msg.Matches[0]
	var value string
	ok, err := b.Store.Get(key, &value)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve key %q from brain", key)
	}

	if ok {
		msg.Respond("%s is %s", key, value)
	} else {
		msg.Respond("I do not remember %q", key)
	}

	return nil
}

func (b *ExampleBot) ShowKeys(msg mili.Message) error {
	keys, err := b.Store.Keys()
	if err != nil {
		return err
	}

	msg.Respond("I got %d keys:", len(keys))
	for i, k := range keys {
		msg.Respond("%d) %q", i+1, k)
	}
	return nil
}
