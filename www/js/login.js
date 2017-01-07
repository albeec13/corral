var SubmitLogin = function(command) {
    var xhr = new XMLHttpRequest();
    var loginJSON = {
        email: '',
        password: ''
    };

    xhr.open('POST', 'https://thewalr.us/corral/API/' + command);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function() {
        var userInfo = JSON.parse(xhr.responseText);
        
        if(userInfo["session_token"]) {
            try{
                if('localStorage' in window && window['localStorage'] !== null) {
                    localStorage.setItem("session_token",userInfo["session_token"]);
                }
            } catch (e) {
                window.name = JSON.stringify(userInfo);
            }
            document.getElementById("response").innerHTML = "";
            window.location = "user.html"
        } else if(userInfo["error"]) {
            document.getElementById("response").innerHTML = userInfo["error"];
        } else {
            document.getElementById("response").innerHTML = JSON.stringify(userInfo);
        }
    };

    loginJSON.email = document.getElementById("email").value;
    loginJSON.password = document.getElementById("password").value;
    xhr.send(JSON.stringify(loginJSON));
}
