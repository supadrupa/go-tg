# tg

[![GoDoc Widget]][GoDoc]
[![Go Report Card Widget]][Go Report Card]
[![Travis Widget]][Travis]

🤖 Golang bindings for Telegram Bot API (WIP 🚧)


## Todo

 - [x] Transport
   - [x] Request
   - [x] Response
 - [ ] Client
   - [ ] [getUpdates](https://core.telegram.org/bots/api#getUpdates)
   - [ ] [setWebhook](https://core.telegram.org/bots/api#setwebhook)
   - [ ] [deleteWebhook](https://core.telegram.org/bots/api#deletewebhook)
   - [ ] [getWebhookInfo](https://core.telegram.org/bots/api#getwebhookinfo)
   - [x] [getMe](https://core.telegram.org/bots/api#getme)
   - [ ] [sendMessage](https://core.telegram.org/bots/api#sendmessage)
   - [ ] [forwardMessage](https://core.telegram.org/bots/api#forwardmessage)
   - [ ] [sendPhoto](https://core.telegram.org/bots/api#sendphoto)
   - [ ] [sendAudio](https://core.telegram.org/bots/api#sendaudio)
   - [ ] [sendDocument](https://core.telegram.org/bots/api#senddocument)
   - [ ] [sendVideo](https://core.telegram.org/bots/api#sendvideo)
   - [ ] [sendAnimation](https://core.telegram.org/bots/api#sendanimation)
   - [ ] [sendVoice](https://core.telegram.org/bots/api#sendvoice)
   - [ ] [sendVideoNote](https://core.telegram.org/bots/api#sendvideonote)
   - [ ] [sendMediaGroup](https://core.telegram.org/bots/api#sendMediaGroup)
   - [ ] [sendLocation](https://core.telegram.org/bots/api#sendlocation)
   - [ ] [editMessageLiveLocation](https://core.telegram.org/bots/api#editmessagelivelocation)
   - [ ] [stopMessageLiveLocation](https://core.telegram.org/bots/api#stopmessagelivelocation)
   - [ ] [sendVenue](https://core.telegram.org/bots/api#sendvenue)
   - [ ] [sendContact](https://core.telegram.org/bots/api#sendcontact)
   - [ ] [sendPoll](https://core.telegram.org/bots/api#sendpoll)
   - [ ] [sendChatAction](https://core.telegram.org/bots/api#sendchataction)
   - [x] [getUserProfilesPhoto](https://core.telegram.org/bots/api#getuserprofilephotos)
   - [x] [getFile](https://core.telegram.org/bots/api#getfile)
   - [x] [kickChatMember](https://core.telegram.org/bots/api#kickchatmember)
   - [x] [unbanChatMember](https://core.telegram.org/bots/api#unbanchatmember)
   - [x] [restrictChatMember](https://core.telegram.org/bots/api#restrictchatmember)
   - [ ] [promoteChatMember](https://core.telegram.org/bots/api#promotechatmember)
   - [ ] [exportChatInviteLink](https://core.telegram.org/bots/api#exportchatinvitelink)
   - [ ] [setChatPhoto](https://core.telegram.org/bots/api#setchatphoto)
   - [ ] [deleteChatPhoto](https://core.telegram.org/bots/api#deletechatphoto)
   - [x] [setChatTitle](https://core.telegram.org/bots/api#setchattitle)
   - [x] [setChatDescription](https://core.telegram.org/bots/api#setchatdescription)
   - [ ] [pinChatMessage](https://core.telegram.org/bots/api#pinchatmessage)
   - [ ] [unpinChatMessage](https://core.telegram.org/bots/api#unpinchatmessage)
   - [ ] [leaveChat](https://core.telegram.org/bots/api#leavechat)
   - [x] [getChat](https://core.telegram.org/bots/api#getchat)
   - [x] [getChatAdministrators](https://core.telegram.org/bots/api#getchatadministrators)
   - [x] [getChatMembersCount](https://core.telegram.org/bots/api#getchatmemberscount)
   - [ ] [getChatMember](https://core.telegram.org/bots/api#getchatmember)
   - [ ] [setChatStickerSet](https://core.telegram.org/bots/api#setchatstickerset)
   - [ ] [deleteChatStickerSet](https://core.telegram.org/bots/api#deletechatstickerset)
   - [ ] [answerCallbackQuery](https://core.telegram.org/bots/api#answercallbackquery)
   - [ ] [editMessageText](https://core.telegram.org/bots/api#editmessagetext)
   - [ ] [editMessageCaption](https://core.telegram.org/bots/api#editmessagecaption)
   - [ ] [editMessageReplyMarkup](https://core.telegram.org/bots/api#editmessagereplymarkup)
   - [ ] [stopPoll](https://core.telegram.org/bots/api#stoppoll)
   - [ ] [deleteMessage](https://core.telegram.org/bots/api#deletemessage)
   - [ ] [sendSticker](https://core.telegram.org/bots/api#sendsticker)
   - [ ] [getStickerSet](https://core.telegram.org/bots/api#getstickerset)
   - [ ] [createNewStickerSet](https://core.telegram.org/bots/api#createnewstickerset)
   - [ ] [addStickerToSet](https://core.telegram.org/bots/api#addstickertoset)
   - [ ] [setStickerPositionInSet](https://core.telegram.org/bots/api#setstickerpositioninset)
   - [ ] [answerInlineQuery](https://core.telegram.org/bots/api#answerinlinequery)
   - [ ] [sendInvoice](https://core.telegram.org/bots/api#payments)
   - [ ] [answerShippingQuery](https://core.telegram.org/bots/api#answershippingquery)
   - [ ] [answerPreCheckoutQuery](https://core.telegram.org/bots/api#answerprecheckoutquery)
   - [ ] [sendGame](https://core.telegram.org/bots/api#sendgame)
   - [ ] [setGameScore](https://core.telegram.org/bots/api#setgamescore)
   - [ ] [getGameHighScore](https://core.telegram.org/bots/api#setgamescore)


[GoDoc]: https://godoc.org/github.com/mr-linch/go-tg
[GoDoc Widget]: https://godoc.org/github.com/mr-linch/go-tg?status.svg
[Go Report Card]: https://goreportcard.com/report/github.com/mr-linch/go-tg
[Go Report Card Widget]: https://goreportcard.com/badge/github.com/mr-linch/go-tg
[Travis]: https://travis-ci.org/mr-linch/go-tg
[Travis Widget]: https://travis-ci.org/mr-linch/go-tg.svg?branch=master