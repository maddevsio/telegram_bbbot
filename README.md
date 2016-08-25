# telegram_bbbot
Telegram Bug Bounty Bot

# History

* This bot adopted special for deploying to Heroku
* General purposes of this got - "Be helpful for infosec community!"
* Bot use ```https://github.com/maddevsio/bbcrawler``` for fetching information
* Used heroku ```https://github.com/heroku/go-getting-started``` as a template for project
* For bot used free account on ```heroku.com``` and ```firebase.com```

# Purpose

* Purposes of bot:
  * "Deliver information as fast as possible!"
  * "Be helpful for infosec community"

# Architecture

* For web server used ```GIN``` 
  * ```github.com/gin-gonic/gin```
* For Bot functionality used ```telegram-bot-api.v4```
  * ```gopkg.in/telegram-bot-api.v4```

# Bot configuration

* ```TELEGRAM_BBBOT_TOKEN``` - Telegram Api token received from @BotFather
* ```TELEGRAM_BBBOT_URL``` - Webhook url to bot public web address
* ```PORT``` -  Standard heroku ENV variable for port number
* ```TELEGRAM_BBBOT_FIREBASE_TOKEN``` - Firebase database token
* ```TELEGRAM_BBBOT_FIREBASE_URL``` -  Url to firebase project
* ```TELEGRAM_BBBOT_HO_SEARCH_URL``` - HackerOne search url (crawler)
* ```TELEGRAM_BBBOT_CHANNEL``` -  Public channel identifier, for example ```@some_channel_name```
* ```TELEGRAM_BBBOT_HOST``` - Public bot host url for ping purposes (for disabling sleeping functionality after 30 min of inactivity)
* ```TELEGRAM_BBBOT_H1_HACK_SEARCH_URL``` - HackerOne hacktivity url (crawler)
* ```TELEGRAM_BBBOT_BUGCROWD_NEW_PROG_URL``` - BugCrowd url for crawling new programs (crawler)

# Bot workflow

* Bot started
* Fetching data from firebase ```(synchronising)```
* Crawling programs from hackerone.com ```(in parallel)```
* Crawling hacktivity from hackerone.com ```(in parallel)```
* Crawling programs from bugcrowd.com ```(in parallel)```
* Determining new data from all crawled information ```(in parallel)```
* Publishing data to telegram channel from ```ENV``` variable

* **Note:** If instance of bot at heroku.com restarted all data restored from firebase storage.

# MIT License

Copyright (c) 2016 Maddevs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.