# Chat.io example pizza BOT integration

This is a simple example of chat.io integration that creates a BOT Agent that listens to all the messages through all the chats and reacts on `pizza` word. When anyone writes a `pizza` word the BOT joins the channel, says `Pizza is on the way!` to agents only and leaves the channel.

## What does it exactly do?

First of all the integration must handle the installation process on a license. To do it it binds to local port to parse `HTTP GET` request that is a result of agent-sso redirect to receive a `code` and exchange it to `access_token` of an agent (see https://www.chat.io/docs/authorization/#public-server-side-apps).
