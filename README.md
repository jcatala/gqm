# Go Quick Message (Telegram notification)


## Install

Remember to have the `$GOPATH/bin` on your `$PATH`

```bash
go get -u github.com/jcatala/gqm
mkdir ~/.config/gqm/
vim ~/.config/gqm/gqm.ini
```

The config file must be something like this:

```yaml
apikey        = YOUR TOKEN THAT BOTHFATHERS GIVES TO YOU
```

The `gqm` will automatically update the `chat_id` on the config file, that must match the `last update` from the bot.
Sometimes the bot `does not have new updates`, so you can't rely on the `updates` to get the **chat id**, that's the reason why I save it on the config file.

## Usage

```bash
$ gqm -h
  -follow
        To keep the stdin open
  -markdown
        Force markdown on the entire message, if is not, do it by yourself adding backquotes
  -verbose
        To be verbose
```

## Use cases

```bash
# To get a response without markdown
echo "123 test without markdown" | gqm -verbose

# To get a response with forced markdown
echo "123 test with markdown" | gqm -verbose -markdown

# To make the stdin open to get constant updates
tail -f "dns.log" | gqm -follow
```