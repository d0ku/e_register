<!DOCTYPE html>

<!-- prawy górny róg: wyloguj i licznik pozostałego czasu (js) -->
<!-- na górze opcje: dane plan lekcje(domyślnie wybrane?) klasy Uczniowie(lista całej szkoły z wyszukiwarką?)-->
<html lang="pl">
    <head>
        <title>E-dziennik: Nauczyciel</title>
        <meta charset="utf-8">
        <link rel="stylesheet" type="text/css" href="/page/CSS/main_teacher.css">
        <meta name="description" content= "Strona główna -> Nauczyciel."> 
        <meta name="viewport"  content= "width=device-width, initial-scale=1.0"/> 
        <link href="https://fonts.googleapis.com/css?family=Lato|Merriweather|Roboto|Yellowtail" rel="stylesheet"> 
    </head>
    <body>
        <div id="container">
            <div class="frame">
                <nav class="teacher_nav">
                    <ul>
                        <li><a href="/main/teacher/{{ .UserID }}/about" >Dane</a></li>
                        <li><a href="/main/teacher/{{ .UserID }}/timetable" >Plan lekcji</a></li>
                        <li><a href="/main/teacher/{{ .UserID }}/classes" >Lekcje</a></li>
                        <li><a href="/main/teacher/{{ .UserID }}/students" >Uczniowie</a></li>
                    </ul>
                </nav>
                <div id="description">
                    Witamy na stronie głównej.
                </div>
            </div>
        </div>
        <div id="made_by">
            By d0ku 2018
        </div>
    </body>
</html>
