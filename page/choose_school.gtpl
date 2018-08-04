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
                    {{ range $key, $school := .Schools }}
                    <a class="schoolBox" href="/main/{{ $.UserType }}/{{ $school.Id }}">
                        <span class="name">{{ $school.FullName }}</span>
                        <span class="city">{{ $school.City }}</span>
                        <span class="street">{{ $school.Street }}</span>
                    </a>
                    {{ end }}
                    <br/>
                </div>
            </div>
        </div>
        <div id="made_by">
            By d0ku 2018
        </div>
    </body>
</html>
