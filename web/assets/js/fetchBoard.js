let queryString = window.location.search;
let urlParams = new URLSearchParams(queryString);
let projectID = urlParams.get("id");
let projectTeamID = urlParams.get("team");
let projectDetail;

sessionStorage.removeItem("board");

let kanbanBoard = [];

let getProjectDetailOptions = {
  method: "GET",
  credentials: "omit",
  headers: {
    Authorization: "Bearer " + sessionStorage.getItem("token"),
    "Content-Type": "text/plain",
  },
  redirect: "follow",
};

fetch(
  "http://127.0.0.1:12345/teams/" + projectTeamID + "/projects/" + projectID,
  getProjectDetailOptions
)
  .then((response) => response.json())
  .then((result) => {
    if (result.error.code == 200) {
      console.log(result);
      projectDetail = result;
      if (projectDetail.columns != null) {
        for (let i = 0; i < projectDetail.columns.length; i++) {
          let column = {
            id: projectDetail.columns[i].id.toString(),
            title: projectDetail.columns[i].title.toString(),
            project_id: projectDetail.columns[i].project_id,
            item: projectDetail.columns[i].cards,
            created_by: projectDetail.columns[i].created_by,
            index: projectDetail.columns[i].index,
            maxIndex: projectDetail.columns[i].maxIndex,
            projectIndex: projectDetail.columns[i].projectIndex,
            createdAt: projectDetail.columns[i].createdAt,
            updatedAt: projectDetail.columns[i].updatedAt,
            deletedAt: projectDetail.columns[i].deletedAt,
          };
          kanbanBoard[i] = column;
        }
      }
      sessionStorage.setItem("board", JSON.stringify(kanbanBoard));
    }
  });
