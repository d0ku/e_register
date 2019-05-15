# e_register

[![Build Status](https://travis-ci.com/d0ku/e_register.svg?token=czCs7ySFgsJtHB5vZwPp&branch=master)](https://travis-ci.com/d0ku/e_register)

## Disclaimer

That's first web app I have ever written, and I really believed it will be much less work than it happened to be. That version has base for sessions, login and stuff like that implemented, it has no real use at the moment.

## Important

By default Makefile configures PostgreSQL to be run on local machine, not on remote server. If that's not the case for you, it has to be changed.

### Dependencies
[PostgreSQL driver for Go](https://github.com/lib/pq)

[Minifier for JS](https://github.com/tdewolff/minify)

[Sass](https://sass-lang.com)

### How to build

You should configure files in config directory and change variables on top of Makefile in main directory to suit your needs.

### Notes

Timeout for too many login tries won't work if server is run behind a proxy (nginx, apache, etc.)

Config files located in config/ should only be used as a base for writing your own. Certificate is self-signed, so it is not worth anything, and config.cfg file should be edited to match specific server and domain.
