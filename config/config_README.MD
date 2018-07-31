It does not make much sense (if any) to run that website without TLS, so there is no option to disable it.
host -> hostname on which you want to serve (defaults to localhost)
server_cert -> path to tls server certificate
server_key -> path to tls server private key
redirect_http_to_https -> 1 enables redirect, other values do nothing
http_port -> port on which http request should be catched and redirected to https_port (matters only if redirect_http_to_https is set to 1). You will probably want to set that value to 80.
https_port -> port on which application will be accessible (all catched http requests will be redirected to that port). You will probably want to set that value to 443.
db_username -> username used to access database
db_name -> name of database that should be accessed
db_sslmode -> SSL mode to be used when connecting with database
web_assets_path -> path to folder which contains all HTML, CSS etc. specific files (relative will work, absolute are much better option)
cookie_life_time -> describes time (in seconds) after which user will be automatically logged out if he does not do something. If not specified, defaults to 900 seconds.
