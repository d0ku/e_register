<!DOCTYPE html>
<html lang="pl">
    <head>
        <title>E-dziennik: Wybierz szkołę</title>
        <meta charset="utf-8">
        <link rel="stylesheet" type="text/css" href="/page/CSS/choose_school.css">
        <meta name="description" content= "Strona wyboru szkoły."> 
        <meta name="viewport"  content= "width=device-width, initial-scale=1.0"/> 
        <link href="https://fonts.googleapis.com/css?family=Lato|Merriweather|Roboto|Yellowtail" rel="stylesheet"> 
    </head>
    <body>
        <div id="container">
            <div class="frame">
               <div id="description">
                    <span class="about">
                    Wybierz szkołę którą chcesz zarządzać:
                    </span>
                    {{ range .Schools }}
                    <a class="schoolBox" href="/main/{{ $.UserType }}/{{ .Id }}">
                        <span class="name">{{ .FullName }}</span>
                        <span class="city">{{ .City }}</span>
                        <span class="street">{{ .Street }}</span>
                    </a>
                    {{ end }}
                    <br/>
                    <span class="help">
                    Aby wybrać szkołe kliknij na odpowiedni przycisk.
                    </span>
                </div>
            </div>
        </div>
        <div id="made_by">
            By d0ku 2018
        </div>
    </body>
</html>
