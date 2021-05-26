$(document).ready(function () {
  var token = sessionStorage.getItem("token");

  var testTokenOptions = {
    method: "GET",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + token,
      "Content-Type": "text/plain",
    },
    redirect: "follow",
  };

  fetch("http://127.0.0.1:12345/teams", testTokenOptions)
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code != 401) {
        window.location.href = "../../../web/pages/home.html";
      }
    });

  //Login request
  function loginRequest(user, pass) {
    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "text/plain");

    var raw =
      '{\n  "username" : "' +
      user +
      '",\n  "password" : "' +
      pass +
      '"\n   \n}';

    var requestOptions = {
      method: "POST",
      headers: myHeaders,
      body: raw,
      redirect: "follow",
    };

    fetch("http://127.0.0.1:12345/login", requestOptions)
      .then((response) => response.json())
      .then((result) => {
        console.log(result);
        if (result.error.code == 200) {
          var token = result.data.token;
          sessionStorage.setItem("token", token);
          sessionStorage.setItem("username", String(user));
          window.location.href = "../../../web/pages/home.html";
        } else {
          alert("Incorrect account credentials");
        }
      })
      .catch((error) => {
      });
  }

  $("#loginButton").on("click", function () {
    var username = $("#username").val();
    var password = $("#password").val();
    if (username == "" || password == "") return;

    loginRequest(username, password);
  });
});
