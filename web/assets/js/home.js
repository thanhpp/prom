var token = sessionStorage.getItem("token");
var teams;

getTeamsRequest(token);

function getTeamsRequest(tokenString) {
  var getTeamOptions = {
    method: "GET",
    credentials: "omit",
    headers: {
      Authorization: "Bearer " + tokenString,
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
        console.log(personalProjects);
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
        console.log(projects);
        for (i = projects.length - 1; i >= 0; i--) {
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
