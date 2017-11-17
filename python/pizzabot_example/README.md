# Chat.io example pizza BOT integration

This is a simple example of chat.io integration in python2 that creates a BOT Agent that listens to all the messages through all the chats and reacts on `pizza` word. When anyone writes a `pizza` word the BOT joins the channel, says `Pizza is on the way!` to agents only and leaves the channel.

## What does it exactly do?

* First of all the integration must handle the installation process on a license. To do it it binds to local port to parse `HTTP GET` request that is a result of agent-sso redirect to receive a `OAuth code` and exchange it to `access_token` of an agent (see https://www.chat.io/docs/authorization/#server-side-apps).

* After getting the `access_token` it gets `access_token` iformation from agent SSO (`license` number in this case only).

* Once successfully obtaining `access_token` and `license` number it configures the necessary resources in the license:
  * it creates BOT Agent via [configuration-api](https://www.chat.io/docs/configuration-api/api-reference/v0.3/#bot-agent) that is used to receive webhooks with chat events and send messages to chats
  * while creating it initializes BOT Agent status with `not accepting chat` value to make it online but not visible to [chat router](https://www.chat.io/docs/apis-overview/#automatic-routing)

* It listens on incoming webhooks being sent via `HTTP POST`. In this case it listens on `incoming_chat_thread` and `incoming_event` webhooks only. It uses [agent-api](https://www.chat.io/docs/agent-api/) to react on incoming messages with `pizza` word with:
  * joining the chat with `pizza` message ([API method](https://www.chat.io/docs/agent-api/api-reference/v0.3/#join-chat))
  * sending `Pizza is on the way!` to all agents in the chat ([API method](https://www.chat.io/docs/agent-api/api-reference/v0.3/#send-event))
  * leaving the chat ([API method](https://www.chat.io/docs/agent-api/api-reference/v0.3/#remove-from-chat))

## I want to create this integration on my own

It's quite simple :) Just follow the steps below to run the integration on your local machine.

* Download the [integration script](./pizzabot.py)

* Create your own application in [developers console](https://console.chat.io). Remember that you have to be logged in with the same account that you created a chat.io product license with (https://chat.io) to make your private server-side integration.

* Set redirect URI to point the address where the integration script will listen on (eg `http://localhost:5000/token`)

* After creating integration in developers console prepare the [downloaded script](./pizzabot.py) replacing the following parameters at the top of the script:
  * `<CLIENT_ID>` - Client ID of newly created application in developers console
  * `<CLIENT_SECRET>` - Client Secret of newly created application in developers console
  * `<REDIRECT_URI>` - Redirect URI of newly created application in developers console (eg `http://localhost:5000/token`)
  * `<BOT_AGENT_WEBHOOKS_URL>` - URL for incoming webhooks with chat events. If you do not have possibility to expose your integration script to the internet (no external IP address) you can use eg [ngrok](https://dashboard.ngrok.com/user/signup) to forward all the webhooks to your local machine, the `<BOT_AGENT_WEBHOOKS_URL>` parameter would look like this `http://8a471b50.ngrok.io/webhook` then.

* Now you can run your script by running the following command: `./pizzabot.py 5000` to bind it on `localhost:5000`

* Now you have the setup done and you can install the integration on your chat.io product license. Currently the simplest way to do it is just open in any web browser the following link: https://accounts.chat.io/?response_type=code&client_id=**<CLIENT_ID>**&redirect_uri=**<REDIRECT_URI>** with your `<CLIENT_ID>` and `<REDIRECT_URI>` filled in (eg `https://accounts.chat.io/?response_type=code&client_id=ha778a18997bf12de0cjhd79a9s01d2sc&redirect_uri=http%3A%2F%2Flocalhost%3A5000%2Ftoken%2F`). It will automatically redirect you to your integration script with `OAuth code` and the script will do the rest. Currently it may take up to 2 minutes to make the BOT Agent ready to listen on the chat events.

* When you have installed the integration you can test your BOT now :) Just log in to [chat.io webapp](https://app.chat.io) and open customer chat link you will find [here](https://app.chat.io/settings/channel-direct-link). Then start send a message in customer window containing `pizza` word like `Hey! where is my pizza?`. Then check out the [chat.io webapp](https://app.chat.io) and you should see new chat with customer with information about joining BOT agent, writing `Pizza is on the way!` and leaving the chat.

