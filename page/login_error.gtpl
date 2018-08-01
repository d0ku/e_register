<!DOCTYPE html>
<html lang="pl">
    <head>
        <title>E-dziennik: Błąd logowania</title>
        <meta charset="utf-8">
        <link rel="stylesheet" type="text/css" href="/page/CSS/not_logged.css">
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
                <div class="help">
                    {{ if .HasTimeout}}
                    <div id="timeoutInfo">
                        Z powodu wielu prób logowania następna próba może nastąpić dopiero za 
                        <span id="timeoutCounter">{{ .Timeout }}</span> sek.
                    </div>
                    {{ else }} 
                    Kliknij na przycisk, aby spróbować ponownie.
                    {{ end }}
                </div>
                <a id="try_again_button" class="button" href="/login/{{ .UserType }}">Spróbuj ponownie</a>
            </div>
       </div>
       <div id="made_by">
           By d0ku 2018
       </div>
       <!-- Call js script to run timer down. -->
       <script src="/page/js/timeoutCounter.js"></script>
    </body>
</html>
