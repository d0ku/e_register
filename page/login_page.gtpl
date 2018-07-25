<!DOCTYPE html>
<!--TODO:add meta-->
<html lang=pl>
    <head>
    </head>
    <body>
        <form action="/login" method="post">
            <input type="hidden" name="userType" value="{{.}}">
            Nazwa użytkownika:<input type="text" name="username">
            Hasło:<input type="password" name="password">
            <input type="submit" value="Login">
        </form>
    </body>
</html>
