# cloud-mail-ru-cli
Implementation of [mail.ru cloud](https://cloud.mail.ru/) API written in golang.

A trivial command line client (**cloud-cli**) based on the api also supplied.

# Installation & Usage
    go get github.com/PhantomXCool/cloud-mail-ru-cli
    cd $HOME/go/src/github.com/PhantomXCool/cloud-mail-ru-cli
    go build -o cloud-cli

    export MAILRU_USER=<your mail.ru username>
    export MAILRU_PASSWORD=<your mail.ru password>
    ./cloud-cli -help

# Documentation
Most of API documented using godoc. To view the inline documentation use:

    godoc github.com/PhantomXCool/cloud-mail-ru-cli/Api
