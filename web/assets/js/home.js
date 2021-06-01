var token = sessionStorage.getItem("token");
var teams;
var personalTeamID;
getTeamsRequest();

function getTeamsRequest() {
  let getTeamOptions = {
    method: "GET",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + token,
      "Content-Type": "text/plain",
    },
    redirect: "follow",
  };

  fetch("http://127.0.0.1:12345/teams", getTeamOptions)
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        teams = result.teams;
        personalTeamID = teams[0].id;
        getPersonalProjects(personalTeamID);
        assignTeamsHTML();
        assignRecentProjects();
      }
    });
}

function getPersonalProjects(teamID) {
  let personalProjects;
  let getProjectsOptions = {
    method: "GET",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + token,
      "Content-Type": "text/plain",
    },
    redirect: "follow",
  };

  fetch(
    "http://127.0.0.1:12345/teams/" + teamID + "/projects",
    getProjectsOptions
  )
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        personalProjects = result.projects;

        for (i = personalProjects.length - 1; i >= 0; i--) {
          let projectName = personalProjects[i].name;
          let projectID = personalProjects[i].id;

          // Create a temporary <div> to load into
          let li = document.createElement("li");
          li.setAttribute("class", "nav-item");
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
            '" class="nav-link"><i class="fas fa-briefcase"></i> <span>' +
            projectName.toString() +
            "</span></a>";

          // Write the <div> to the HTML container
          document.getElementById("personalProjects").after(li);
        }
      }
    });
}

function assignTeamsHTML() {
  for (i = teams.length - 1; i > 0; i--) {
    let teamName = teams[i].name;
    let teamID = teams[i].id;

    // Create a temporary <div> to load into
    let li = document.createElement("li");
    li.setAttribute("class", "nav-item");
    li.setAttribute("team-id", teamID);
    let teamURL = new URL("http://127.0.0.1:5501/web/pages/team.html");

    let teamURLSearchParams = teamURL.searchParams;

    teamURLSearchParams.set("id", teamID.toString());
    

    li.innerHTML =
      '<a href="' +
      teamURL.toString() +
      '" class="nav-link"><i class="fas fa-users"></i> <span>' +
      teamName.toString() +
      "</span></a>";

    // Write the <div> to the HTML container
    document.getElementById("teams").after(li);
  }
}

function assignRecentProjects() {
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

  fetch(
    "http://127.0.0.1:12345/teams/" + personalTeamID + "/projects",
    getProjectsOptions
  )
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        projects = result.projects;
        console.log(projects);
        for (var i = projects.length - 1; i >= 0; i--) {
          let projectName = projects[i].name;
          let projectID = projects[i].id;

          // Create a temporary <div> to load into
          let li = document.createElement("article");
          li.setAttribute("class", "article col-12 col-sm-6 col-md-3");
          li.setAttribute("project-id", projectID);

          let projectURL = new URL(
            "http://127.0.0.1:5501/web/pages/project.html"
          );

          let projectURLSearchParams = projectURL.searchParams;

          projectURLSearchParams.set("id", projectID.toString());
          projectURLSearchParams.set("team", personalTeamID);
          li.innerHTML =
            '<a href="' +
            projectURL.toString() +
            '"><div class="article-header text-white"><div class="article-image bg-primary"></div><div class="article-title"><h2>' +
            projectName.toString() +
            "</h2></div></div></a>";

          // Write the <div> to the HTML container
          document.getElementById("recentProjects").appendChild(li);
        }
      }
    });
}
$("#navBarNewProject").on("click", function () {
  updateChooseTeam(teams);
});

$("#heroNewProject").on("click", function () {
  updateChooseTeam(teams);
});

$("#createNewTeamButton").on("click", function () {
  let createNewTeamName = $("#createNewTeamName").val();

  makeNewTeamQuick(createNewTeamName);
});

function makeNewTeamQuick(teamName) {
  let myHeaders = new Headers();
  myHeaders.append("Content-Type", "text/plain");

  let raw = '{\n  "teamName" : "' + teamName.toString() + '"\n}';

  let newTeamOptions = {
    method: "POST",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + token,
      "Content-Type": "text/plain",
    },
    body: raw,
    redirect: "follow",
  };

  fetch("http://127.0.0.1:12345/teams", newTeamOptions)
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        let newTeams;

        let getTeamOptions = {
          method: "GET",
          credentials: "omit",
          headers: {
            Authorization: "Bearer " + token,
            "Content-Type": "text/plain",
          },
          redirect: "follow",
        };

        fetch("http://127.0.0.1:12345/teams", getTeamOptions)
          .then((response) => response.json())
          .then((result) => {
            if (result.error.code == 200) {
              newTeams = result.teams;
              updateChooseTeam(newTeams);
              document.getElementById("createNewTeamName").value = "";
              alert("Add team " + teamName + " successfully!");
            }
          });
      }
    });
}

function updateChooseTeam(teams) {
  let teamSelectParent = document.getElementById("createSelectTeam");
  teamSelectParent.innerHTML = "";

  for (let i = 0; i < teams.length; i++) {
    let teamID = teams[i].id;
    let teamName = teams[i].name;

    let option = document.createElement("option");
    option.setAttribute("value", teamID);
    option.text = teamName;
    teamSelectParent.append(option);
  }
}

$("#createNewProjectButton").on("click", function () {
  let projectName = $("#createProjectName").val();
  let assignedTeam = parseInt($("#createSelectTeam").val());
  console.log(assignedTeam);

  let myHeaders = new Headers();
  myHeaders.append("Content-Type", "text/plain");

  let raw = '{\n  "projectName" : "' + projectName.toString() + '"\n}';

  let newProjectOptions = {
    method: "POST",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + token,
      "Content-Type": "text/plain",
    },
    body: raw,
    redirect: "follow",
  };

  fetch(
    "http://127.0.0.1:12345/teams/" + assignedTeam + "/projects",
    newProjectOptions
  )
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        let teamProjects;
        let teamProjectsOptions = {
          method: "GET",
          credentials: "omit",
          headers: {
            Authorization: "Bearer " + token,
            "Content-Type": "text/plain",
          },
          redirect: "follow",
        };

        fetch(
          "http://127.0.0.1:12345/teams/" + assignedTeam + "/projects",
          teamProjectsOptions
        )
          .then((response) => response.json())
          .then((result) => {
            if (result.error.code == 200) {
              teamProjects = result.projects;
              let newProjectID = teamProjects[teamProjects.length - 1].id;
              let projectURL = new URL(
                "http://127.0.0.1:5501/web/pages/project.html"
              );

              let projectURLSearchParams = projectURL.searchParams;

              projectURLSearchParams.set("id", newProjectID.toString());
              alert(projectURL);
              // window.location.href = projectURL;
            }
          });
      }
    });
});

$("#addNewPerson").on("input", function (e) {
  let value = $("#addNewPerson").val();
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
            currentNewMemberID = result.users[0].id;
            for (let i = 0; i < members; i++) {
              if (members[i] == value) {
                isDuplicated = true;
              }
            }
            if (isDuplicated) {
              $("#addNewPersonButton")[0].classList.add("disabled");
            } else {
              $("#addNewPersonButton")[0].classList.remove("disabled");
            }
          } else {
            $("#addNewPersonButton")[0].classList.add("disabled");
          }
        } else {
          $("#addNewPersonButton")[0].classList.add("disabled");
        }
      }
    });
});
var currentNewMemberID;
var members = [];
var membersID = [];

$("#navBarNewTeam").on("click", function () {
  members = [];
});

$("#addNewPersonButton").on("click", function () {
  if (!$("#addNewPersonButton")[0].classList.contains("disabled")) {
    $("#addNewPersonButton")[0].classList.add("disabled");

    members.push($("#addNewPerson").val());
    membersID.push(currentNewMemberID);

    let member = document.createElement("div");
    member.setAttribute("class", "card");

    member.innerHTML =
      '<div class="card-header"><h4>' +
      members[members.length - 1] +
      '</h4><div class="card-header-action"><button memberID="' +
      membersID[membersID.length - 1] +
      '" class="btn btn-primary member-remove-button"><i class="fas fa-minus-circle"></i></button></div></div>';

    // Write the <div> to the HTML container
    document.getElementById("members").appendChild(member);
    document.getElementById("addNewPerson").value = "";

    $(".member-remove-button").on("click", function () {
      let removedMemberID = $(this).attr("memberid");

      for (let i = 0; i < membersID.length; i++) {
        if (membersID[i] == removedMemberID) {
          membersID.splice(i, 1);
          members.splice(i, 1);
        }
      }
      $(this).parent().closest(".card").remove();
    });
  }
});

$("#createNewTeamModalButton").on("click", function () {
  let teamName = $("#createTeamName").val();

  let myHeaders = new Headers();
  myHeaders.append("Content-Type", "text/plain");

  let raw = '{\n  "teamName" : "' + teamName.toString() + '"\n}';

  let newTeamOptions = {
    method: "POST",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + token,
      "Content-Type": "text/plain",
    },
    body: raw,
    redirect: "follow",
  };

  fetch("http://127.0.0.1:12345/teams", newTeamOptions)
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        let newTeamsID;
        let getTeamOptions = {
          method: "GET",
          credentials: "omit",
          headers: {
            Authorization: "Bearer " + token,
            "Content-Type": "text/plain",
          },
          redirect: "follow",
        };

        fetch("http://127.0.0.1:12345/teams", getTeamOptions)
          .then((response) => response.json())
          .then((result) => {
            if (result.error.code == 200) {
              newTeamsID = result.teams[result.teams.length - 1].id;
              console.log("start");
              for (let i = 0; i < membersID.length; i++) {
                addToTeam(membersID[i], newTeamsID);
              }

              setTimeout(function () {
                console.log("done");
                window.location.href =
                  "http://127.0.0.1:5501/web/pages/team.html?id=" + newTeamsID;
              }, 2000);
            }
          });
      }
    });
});

function addToTeam(memberID, teamID) {
  let myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");

  let raw =
    '{\n  "op": "' + "add" + '",\n  "memberID": ' + parseInt(memberID) + "\n}";
  let newMemberOptions = {
    method: "PUT",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + token,
      "Content-Type": "text/plain",
    },
    body: raw,
    redirect: "follow",
  };

  fetch("http://127.0.0.1:12345/teams/" + teamID.toString(), newMemberOptions)
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        console.log("added " + memberID);
      }
    });
}
