$(document).ready(function () {
  function logoutRequest() {
    var token = sessionStorage.getItem("token");

    var logoutOptions = {
      method: "GET",
      credentials: "omit",
      headers: {
        Authorization: "Bearer " + token,
        "Content-Type": "text/plain",
      },
      redirect: "follow",
    };

    fetch("http://127.0.0.1:12345/logout", logoutOptions)
      .then(() =>{
        window.location.href = "../../../web/pages/login.html";
      });
  }
  
  $("#logoutButton").on("click", function () {
    logoutRequest();
  });

  $(".name").text(String(sessionStorage.getItem("username")));
});
