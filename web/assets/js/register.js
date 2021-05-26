$(document).ready(function () {
  var token = sessionStorage.getItem("token");

  //Login request
  function registerRequest(user, pass) {
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

    fetch("http://127.0.0.1:12345/user", requestOptions)
      .then((response) => response.json())
      .then((result) => {
        console.log("regi");
        console.log(result);
        if (result.error.code == 200) {          
          loginRequest(user,pass);
        } else {
          alert("Registration Failed");
        }
      })
  }

  $("#registerButton").on("click", function () {
    var username = $("#username").val();
    var password = $("#password").val();
    if (username == "" || password == "") return;

    registerRequest(username, password);
  });

  function loginRequest(user, pass) {
    console.log("login req")
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
        console.log("login");
        console.log(result);
        if (result.error.code == 200) {
          var token = result.data.token;
          sessionStorage.setItem("token", token);
          sessionStorage.setItem("username", String(user));
          window.location.href = "../../../web/pages/login.html";
        } else {
          alert("Incorrect account credentials");
        }
      })
      .catch((error) => {
        console.log("Không kết nối được tới máy chủ", error);
        alert("Không kết nối được tới máy chủ");
      });
  }
});
