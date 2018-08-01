<!DOCTYPE html>
<!--TODO:add meta-->
<html lang=pl>
    <head>
        <title>E-dziennik: Logowanie</title>
        <meta charset="utf-8">
        <link rel="stylesheet" type="text/css" href="/page/CSS/login_form.css">
        <meta name="description" content= "Strona logowania."> 
        <meta name="viewport"  content= "width=device-width, initial-scale=1.0"/> 
        <link href="https://fonts.googleapis.com/css?family=Lato|Merriweather|Roboto|Yellowtail" rel="stylesheet">
    </head>
    <body>
        <div id="container">
            <div class="frame">
                <form action="/login" method="post">
                    <input type="hidden" name="userType" value="{{.}}">
                    Nazwa użytkownika:<input type="text" name="username">
                    Hasło:<input type="password" name="password">
                    <input type="submit" value="Login">
                </form>
            </div>
        </div>
        <div id="made_by">
            By d0ku 2018
        </div>
    </body>
</html>
