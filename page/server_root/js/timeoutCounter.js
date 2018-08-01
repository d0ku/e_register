function countdown() {
    var redirectButton = document.getElementById("try_again_button");
    var originalLink = redirectButton.href;
    redirectButton.href = "#";

    var timeObj = document.getElementById("timeoutCounter");

    if (timeObj == null) {
        //There is no need to count down.
        redirectButton.href = originalLink;
        return;
    }
    var timeLeft = parseInt(timeObj.innerHTML) + 1;
    var timeoutInfo = document.getElementById("timeoutInfo");

    var funcid = 0;
    funcid = setInterval(function() {
        timeObj.innerHTML = timeLeft - 1;
        timeLeft--;
        if (timeLeft == 0) {
            clearInterval(funcid);
            timeoutInfo.innerHTML = "";
            redirectButton.href = originalLink;
        }
    }, 1000);
}

countdown();
