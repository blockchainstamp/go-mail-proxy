module github.com/blockchainstamp/go-mail-proxy

go 1.19

require (
	github.com/emersion/go-imap v1.2.1
	github.com/emersion/go-imap-id v0.0.0-20190926060100-f94a56b9ecde
	github.com/emersion/go-message v0.15.0
	github.com/emersion/go-smtp v0.16.0
	github.com/emersion/go-textwrapper v0.0.0-20200911093747-65d896831594
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.6.1
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require (
	github.com/emersion/go-sasl v0.0.0-20220912192320-0145f2c60ead // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.6.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)

replace github.com/emersion/go-imap => /Users/hyperorchid/bmail/go-imap
