<!DOCTYPE html>
<html lang="pl">
    <head>
        <title>E-dziennik: Przekierowanie</title>
        <meta charset="utf-8">
        <link rel="stylesheet" type="text/css" href="CSS/not_logged.css">
        <meta name="description" content= "Strona przekierowująca do logowania."> 
        <meta name="viewport"  content= "width=device-width, initial-scale=1.0"/> 
        <link href="https://fonts.googleapis.com/css?family=Lato|Merriweather|Roboto|Yellowtail" rel="stylesheet"> 
    </head>
    <body>
        <div id="container">
            <div class="frame">
                <span class="warning">
                    Wprowadzona nazwa użytkownika i hasło nie pasują do żadnego użytkownika.
                </span>
                <span class="help">
                    {{ if .HasTimeout}}
                    Z powodu wielu prób logowania następna próba może nastąpić dopiero za {{ .Timeout }} sekund.
                    {{ else }} 
                    Kliknij na przycisk, aby spróbować ponownie.
                    {{ end }}
                </span>
                <a class="button" href="/login/{{ .UserType }}">PRZEKIERUJ</a>
            </div>
       </div>
       <div id="made_by">
           By d0ku 2018
       </div>
    </body>
</html>
