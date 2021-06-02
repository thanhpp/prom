let queryString = window.location.search;
let urlParams = new URLSearchParams(queryString);
let teamID = urlParams.get("id");
let persons = [];
let personID = [];
let projects;
let getProjectsOptions = {
  method: "GET",
  credentials: "omit",
  headers: {
    Authorization: "Bearer " + token,
    "Content-Type": "text/plain",
  },
  redirect: "follow",
};

fetch("http://127.0.0.1:12345/teams/" + teamID, getProjectsOptions)
  .then((response) => response.json())
  .then((result) => {
    if (result.error.code == 200) {
      projects = result.team.projects;
      persons = result.team.members;
      console.log(persons.length);

      for (let i = 0; i < persons.length; i++) {
        personID[i] = persons[i].id;
      }
      for (let i = 0; i < persons.length; i++) {
        let person = document.createElement("div");
        person.setAttribute("class", "card");

        person.innerHTML =
          '<div class="card-header"><h4>' +
          persons[i].username +
          '</h4><div class="card-header-action"><button memberID="' +
          personID[i] +
          '" class="btn btn-primary member-remove-button"><i class="fas fa-minus-circle"></i></button></div></div>';

        // Write the <div> to the HTML container
        document.getElementById("currentMembers").appendChild(person);
        document.getElementById("addNewMember").value = "";
      }
      $(".member-remove-button").on("click", function () {
        let removedPersonID = $(this).attr("memberid");
      
        for (let i = 0; i < membersID.length; i++) {
          if (membersID[i] == removedPersonID) {
            membersID.splice(i, 1);
            members.splice(i, 1);
          }
        }
      
        $(this).parent().closest(".card").remove();
      
        let myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
      
        let raw =
          '{\n  "op": "' +
          "remove" +
          '",\n  "memberID": ' +
          parseInt(removedPersonID) +
          "\n}";
        let removeOptions = {
          method: "PUT",
          credentials: "omit",
          headers: {
            Authorization: "Bearer " + token,
            "Content-Type": "text/plain",
          },
          body: raw,
          redirect: "follow",
        };
      
        fetch("http://127.0.0.1:12345/teams/" + teamID.toString(), removeOptions)
          .then((response) => response.json())
          .then((result) => {
            if (result.error.code == 200) {
              console.log(result);
            }
          });
      });
      document.getElementById("teamName").innerHTML = result.team.name;
      console.log(result);
      for (var i = projects.length - 1; i >= 0; i--) {
        let projectName = projects[i].name;
        let projectID = projects[i].id;

        // Create a temporary <div> to load into
        let li = document.createElement("article");
        li.setAttribute("class", "article col-12 col-sm-6 col-md-6");
        li.setAttribute("project-id", projectID);

        let projectURL = new URL(
          "http://127.0.0.1:5501/web/pages/project.html"
        );

        let projectURLSearchParams = projectURL.searchParams;

        projectURLSearchParams.set("id", projectID.toString());
        projectURLSearchParams.set("team", teamID);
        li.innerHTML =
          '<a href="' +
          projectURL.toString() +
          '"><div class="article-header text-white"><div class="article-image bg-primary"></div><div class="article-title"><h2>' +
          projectName.toString() +
          "</h2></div></div></a>";

        // Write the <div> to the HTML container
        document.getElementById("projects").appendChild(li);
      }
    }
  });



let currentNewPersonID;
$("#addNewMember").on("input", function (e) {
  let value = $("#addNewMember").val();
  let getUserOptions = {
    method: "GET",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + token,
      "Content-Type": "text/plain",
    },
    redirect: "follow",
  };

  fetch(
    "http://127.0.0.1:12345/user?username=" + value.toString(),
    getUserOptions
  )
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        if (result.users != null && result.users.length > 0) {
          if (result.users[0].username == value) {
            var isDuplicated = false;
            currentNewPersonID = result.users[0].id;
            for (let i = 0; i < persons.length; i++) {
              if (persons[i].username == value) {
                isDuplicated = true;
              }
            }
            if (isDuplicated) {
              $("#addNewMemberButton")[0].classList.add("disabled");
            } else {
              $("#addNewMemberButton")[0].classList.remove("disabled");
            }
          } else {
            $("#addNewMemberButton")[0].classList.add("disabled");
          }
        } else {
          $("#addNewMemberButton")[0].classList.add("disabled");
        }
      }
    });
});
$("#addNewMemberButton").on("click", function () {
  if (!$("#addNewMemberButton")[0].classList.contains("disabled")) {
    $("#addNewMemberButton")[0].classList.add("disabled");

    persons.push($("#addNewMember").val());
    personID.push(currentNewPersonID);

    let person = document.createElement("div");
    person.setAttribute("class", "card");

    person.innerHTML =
      '<div class="card-header"><h4>' +
      persons[persons.length - 1] +
      '</h4><div class="card-header-action"><button memberID="' +
      personID[personID.length - 1] +
      '" class="btn btn-primary member-remove-button"><i class="fas fa-minus-circle"></i></button></div></div>';

    // Write the <div> to the HTML container
    document.getElementById("currentMembers").appendChild(person);
    document.getElementById("addNewMember").value = "";

    let myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");

    let raw =
      '{\n  "op": "' +
      "add" +
      '",\n  "memberID": ' +
      parseInt(personID[personID.length - 1]) +
      "\n}";
    let addOptions = {
      method: "PUT",
      credentials: "omit",
      headers: {
        Authorization: "Bearer " + token,
        "Content-Type": "text/plain",
      },
      body: raw,
      redirect: "follow",
    };

    fetch("http://127.0.0.1:12345/teams/" + teamID.toString(), addOptions)
      .then((response) => response.json())
      .then((result) => {
        if (result.error.code == 200) {
        }
      });
    $(".member-remove-button").on("click", function () {
      let removedPersonID = $(this).attr("memberid");

      for (let i = 0; i < membersID.length; i++) {
        if (membersID[i] == removedPersonID) {
          membersID.splice(i, 1);
          members.splice(i, 1);
        }
      }
      $(this).parent().closest(".card").remove();

      let myHeaders = new Headers();
      myHeaders.append("Content-Type", "application/json");

      let raw =
        '{\n  "op": "' +
        "remove" +
        '",\n  "memberID": ' +
        parseInt(removedPersonID) +
        "\n}";
      let removeOptions = {
        method: "PUT",
        credentials: "omit",
        headers: {
          Authorization: "Bearer " + token,
          "Content-Type": "text/plain",
        },
        body: raw,
        redirect: "follow",
      };

      fetch("http://127.0.0.1:12345/teams/" + teamID.toString(), removeOptions)
        .then((response) => response.json())
        .then((result) => {
          if (result.error.code == 200) {
          }
        });
    });
  }
});
