var SaveSessionCreds = function (userInfo) {
    try{
        if('localStorage' in window && window['localStorage'] !== null) {
            localStorage.setItem("session_token", userInfo["session_token"]);
            localStorage.setItem("session_user",userInfo["session_user"]);
        }
    } catch (e) {
        window.name = JSON.stringify(userInfo);
    }
};

var PurgeSessionCreds = function() {
    try {
        if('localStorage' in window && window['localStorage'] !== null) {
            localStorage.removeItem("session_token");
            localStorage.removeItem("session_user");
        }
    } catch (e) {
        window.name = "";
    }
};

window.onload = function () {
    if (((window.location.pathname === "/login.html") || (window.location.pathname === "/")) &&
        (localStorage.session_token) && (localStorage.session_token.length > 0) &&
        (localStorage.session_user) && (localStorage.session_user.length > 0)) {
            window.location = "/user.html";
    }
    else if (window.location.pathname === "/user.html") {
        var xhr = new XMLHttpRequest();
        
        xhr.open('GET', 'https://thewalr.us/corral/API/data/viewall');
        xhr.setRequestHeader('X-Request-ID', localStorage["session_user"]);
        xhr.setRequestHeader('X-Request-Token', localStorage["session_token"]);
        xhr.onload = function () {
            var resp = JSON.parse(xhr.responseText);

            if(resp["session_token"]) {
                SaveSessionCreds(resp);
                alert(JSON.stringify(resp));
            } else {
                PurgeSessionCreds();
                window.location = "/login.html";
                alert(JSON.stringify(resp)); }
        }
        xhr.send();
    }
};

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
            SaveSessionCreds(userInfo);
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
};
