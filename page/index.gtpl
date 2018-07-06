<!DOCTYPE html>
<html>
    <head>
    <title>Login Page</title>
    <link rel="stylesheet" type="text/css" href="CSS/reset.css">
    </head>
    <body>
        <h1 class="warning"> TEST </h1>
        {{ if .IsLogged }}
        Witamy na stronie domowej szanowny panie {{ .UserName }}!
        <a href="/logout">Log Out</a>
        {{ else }}
        Witamy na stronie domowej!
        <a href="/login">Log In</a>
        {{ end }}
    </body>
</html>
