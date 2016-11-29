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
        document.getElementById("response").innerHTML = JSON.stringify(userInfo);
    };

    loginJSON.email = document.getElementById("email").value;
    loginJSON.password = document.getElementById("password").value;
    xhr.send(JSON.stringify(loginJSON));
}
