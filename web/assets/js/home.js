var token = sessionStorage.getItem("token");
var teams;

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
        teams = result.data.teams;
        getPersonalProjects(teams[0].id);
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
        personalProjects = result.data.projects;

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

  fetch("http://127.0.0.1:12345/teams/" + 1 + "/projects", getProjectsOptions)
    .then((response) => response.json())
    .then((result) => {
      if (result.error.code == 200) {
        projects = result.data.projects;

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
              newTeams = result.data.teams;
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
              teamProjects = result.data.projects;
              let newProjectID = teamProjects[teamProjects.length-1].id;
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
