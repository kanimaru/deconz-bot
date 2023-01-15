# Deconz Bot

Control Zigbee devices via Deconz and their REST Api.
Uses a lightweight templating engine to have a unified view for all bots.

![Example-1](/doc/bot-example-1.png)

Notice: currently only Telegram is supported other Bot APIs can be easily implemented.

## Features

- [x] Light control (brightness, color, temperature) via chat
- [-] Activate Scan for new lights / sensors (TODO: Maybe show new joined devices, also need to be tested)
- [ ] Display of sensor information
- [ ] Managing of groups
- [ ] Optional features for more features via MQTT

## Environment Variables

- DECONZ_ADDRESS - ip or domain of your deconz gateway ex: 192.168.178.21
- DECONZ_PROTO - protocol that should be used default: http
- DECONZ_API_KEY - look in the deconz [documentation](https://dresden-elektronik.github.io/deconz-rest-doc/getting_started/#acquire-an-api-key) how to get it ex: 0123456789abc36
- TELEGRAM_API_KEY - look in the telegram [documentation](https://core.telegram.org/bots#how-do-i-create-a-bot) 1234567890:ABC-DEFGHILJKLMNOPQRSTUVWXYZ0123456
- TELEGRAM_CHAT_ID - enable the commands only for this chat otherwise everyone can control your Smarthome

## Install for Raspberry PI

TODO

## Links to the products
If you want to set up your own Deconz Gateway.

- [RaspBee II](https://amzn.to/3WjZTjC)
- [ConBee II](https://amzn.to/3YvVhZg)

*As an Amazon Associate I earn from qualifying purchases.*
I'm not related to dresden-elektronik in any way.