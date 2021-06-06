# momento

Momento is the backend for a minimalistic, privacy-oriented, open source
gallery app.

It works with PostgreSQL to store it's data and other than that it has no
runtime dependencies.

## Installing

The recommended way is to install via your package manager so it is better
integrated with your default system monitor, logging facilities, etc. In near
future a package will be available for installing via Alpine's apk.

You can however, build the binaries using make:

```
$ make
$ doas make install
```

## Contributing

Send patches to [my email] or open a pull request on Github. You can discuss
more on the [IRC channel].

## License

AGPLv3, see COPYING.

[my email]: mailto:porcellis@eletrotupi.com
[IRC channel]: ircs://irc.libera.chat/##oodnet
