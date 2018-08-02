# e_register

[![Build Status](https://travis-ci.com/d0ku/e_register.svg?token=czCs7ySFgsJtHB5vZwPp&branch=master)](https://travis-ci.com/d0ku/e_register)

### Dependencies
[PostgreSQL driver for Go](https://github.com/lib/pq)
[Minifier for JS](https://github.com/tdewolff/minify)
[Sass](https://sass-lang.com)

### Notes

Timeout for too many login tries won't work if server is run behind a proxy (nginx, apache, etc.)

Config files located in config/ should only be used as a base for writing your own. Certificate is self-signed, so it is not worth anything, and config.cfg file should be edited to match specific server and domain.
