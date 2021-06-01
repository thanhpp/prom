
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
          if (result.error.code == 401) {
            alert("Please login to use Treno");
    window.location.href = "../../../web/pages/login.html";
          } 
        })
  
    $("#loginButton").on("click", function () {
      var username = $("#username").val();
      var password = $("#password").val();
      if (username == "" || password == "") return;
  
      loginRequest(username, password);
    });
  
  