# Diary Bot

Receive the day's appointments from google calendar

## Running

This project requires two things to be set up:

1. Some env vars
1. Some google credentials

Env vars are much the same as my other bots:

* `$SASL_USER` - the user to connect with
* `$SASL_PASSWORD` - the password to connect with
* `$SERVER` - IRC connection details, as `irc://server:6667` or `ircs://server:6697` (`ircs` implies irc-over-tls)
* `$VERIFY_TLS` - Verify TLS, or sack it off. This is of interest to people, like me, running an ircd on localhost with a self-signed cert. Matches "true" as true, and anything else as false
* `$CREDENTIALS_FILE` - Google API credentials file, see below
* `$TOKEN_FILE` - Google API tokens file
* `$TZ` - Timezone to render diary in

However, because this is google, and because google APIs are just awful to work with, we need to do some extra work to get credentials:

Your best bet is to pretty much follow this: https://developers.google.com/calendar/api/quickstart/go

And to then point the env vars `$CREDENTIALS_FILE` and `$TOKEN_FILE` to the files `credentials.json` and `token.json` respectively.

It's a real faff, but it's the only sensible thing. You should keep those files safe, though quite how I don't know (I don't know what happens if either becomes public, or how desktop google apps package these files. It's all a bit of a mystery, and google is unhelpful when I try to... ahem... google for it)
