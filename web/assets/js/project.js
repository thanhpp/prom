var projectName = "";
let KanbanTest;

var token = sessionStorage.getItem("token");

fetch(
  "http://127.0.0.1:12345/teams/" + projectTeamID + "/projects/" + projectID,
  getProjectDetailOptions
)
  .then((response) => response.json())
  .then((result) => {
    if (result.error.code == 200) {
      projectDetail = result;
      projectName = projectDetail.project.name;

      document.getElementById("projectName").innerHTML = projectName;
    }
  });

initKanban(JSON.parse(sessionStorage.getItem("board")));

function initKanban(finalBoard) {
  KanbanTest = new jKanban({
    element: "#myKanban",
    gutter: "10px",
    widthBoard: "300px",
    responsivePercentage: false,
    itemHandleOptions: {
      enabled: false,
    },
    buttonClick: function (el, boardId) {
      // create a form to enter element
      var formItem = document.createElement("form");
      formItem.setAttribute("class", "itemform");
      formItem.innerHTML =
        '<div class="new-item form-group"><input type="text" class="form-control new-item-text" rows="2" autofocus><input style="display: none; visibility: hidden; position: absolute;" type="submit" value></input></div><div class="form-group text-right"><button type="submit" id="submit" class="btn btn-success new-item-button">Add</button><button type="button" id="CancelBtn" class="btn btn-outline-danger pull-right new-item-button">Cancel</button></div>';

      KanbanTest.addForm(boardId, formItem);
      formItem.addEventListener("submit", function (e) {
        e.preventDefault();
        var text = e.target[0].value;
        KanbanTest.addElement(boardId, {
          title: text,
          description: "",
          column_id: boardId,
          created_by: sessionStorage.getItem("userID"),
          index: KanbanTest.getBoardElements(boardId).length + 1,
        });

        let raw =
          '{"card":{"assignedTo":' +
          sessionStorage.getItem("userID") +
          ',"description":"","duedate":0,"title":"' +
          text +
          '"},"columnID":' +
          parseInt(boardId) +
          "}";

        let newCardOptions = {
          method: "POST",
          headers: {
            Authorization: "Bearer " + sessionStorage.getItem("token"),
            "Content-Type": "text/plain",
          },
          body: raw,
          redirect: "follow",
        };

        fetch(
          "http://127.0.0.1:12345/teams/" +
            projectTeamID +
            "/projects/" +
            projectID +
            "/cards",
          newCardOptions
        )
          .then((response) => response.json())
          .then((result) => {
            console.log(result);
            if (result.error.code == 200) {
            }
          });

        formItem.parentNode.removeChild(formItem);
        addClassToNewBoard();
      });
      document.getElementById("CancelBtn").onclick = function () {
        formItem.parentNode.removeChild(formItem);
      };
    },
    itemAddOptions: {
      enabled: true,
      content: "+ New",
      class: "new-card btn btn-outline-primary",
      footer: false,
    },
    click: function (el) {
      console.log("aa");
    },
    dragendBoard: function (el) {
      console.log(el);
      let columnID = parseInt(el.getAttribute("data-id"));
      let nextOfIndex = parseInt(el.getAttribute("data-order"));
      let raw =
        '{\n  "columnID" : ' +
        columnID +
        ',\n  "nextOfIndex" : ' +
        nextOfIndex +
        "\n}";

      let reorderColumnOptions = {
        method: "POST",
        headers: {
          Authorization: "Bearer " + sessionStorage.getItem("token"),
          "Content-Type": "text/plain",
        },
        body: raw,
        redirect: "follow",
      };

      fetch(
        "http://127.0.0.1:12345/teams/" +
          projectTeamID +
          "/projects/" +
          projectID +
          "/columns/reorder",
        reorderColumnOptions
      )
        .then((response) => response.json())
        .then((result) => {
          console.log(result);
          if (result.error.code == 200) {
          }
        });
    },
    dropEl: function (el, target, source, sibling) {
      if (target.parentElement == source.parentElement) {
        let raw =
          '{\n  "cardID" : ' +
          parseInt(el.getAttribute("data-eid")) +
          ',\n  "aboveOfIdx" : ' +
          parseInt(Array.from(el.parentNode.children).indexOf(el) + 2) +
          ',\n  "columnID" : 0' +
          "\n}";

        let cardReorderSameColumnOptions = {
          method: "POST",
          headers: {
            Authorization: "Bearer " + sessionStorage.getItem("token"),
            "Content-Type": "text/plain",
          },
          body: raw,
          redirect: "follow",
        };

        fetch(
          "http://127.0.0.1:12345/teams/" +
            projectTeamID +
            "/projects/" +
            projectID +
            "/cards/reorder",
          cardReorderSameColumnOptions
        )
          .then((response) => response.json())
          .then((result) => {
            console.log(result);
            if (result.error.code == 200) {
            }
          });
      } else {
        let targetColumnID = target.parentElement.getAttribute("data-order");
        let raw =
          '{\n  "cardID" : ' +
          parseInt(el.getAttribute("data-eid")) +
          ',\n  "aboveOfIdx" : ' +
          parseInt(Array.from(el.parentNode.children).indexOf(el) + 1) +
          ',\n  "columnID" : ' +
          targetColumnID +
          "\n}";
        console.log(raw);
        let cardReorderSameColumnOptions = {
          method: "POST",
          headers: {
            Authorization: "Bearer " + sessionStorage.getItem("token"),
            "Content-Type": "text/plain",
          },
          body: raw,
          redirect: "follow",
        };

        fetch(
          "http://127.0.0.1:12345/teams/" +
            projectTeamID +
            "/projects/" +
            projectID +
            "/cards/reorder",
          cardReorderSameColumnOptions
        )
          .then((response) => response.json())
          .then((result) => {
            console.log(result);
            if (result.error.code == 200) {
            }
          });
      }
    },
    boards: finalBoard,
  });
}

$(document).ready(function () {
  console.log("ready!");
  addClassToNewBoard();
});

function addClassToNewBoard() {
  // var containerElements = document.getElementsByClassName("kanban-container");
  // for (i = 0; i < containerElements.length; i++) {
  //   containerElements[i].classList.add("row");

  // }

  var boardElements = document.getElementsByClassName("kanban-board");
  for (i = 0; i < boardElements.length; i++) {
    boardElements[i].classList.add("col", "card", "card-primary");
  }

  var boardHeaderElements = document.getElementsByClassName(
    "kanban-board-header"
  );
  for (i = 0; i < boardHeaderElements.length; i++) {
    boardHeaderElements[i].classList.add("card-header");
  }

  var boardBodyElements = document.getElementsByClassName("kanban-drag");
  for (i = 0; i < boardBodyElements.length; i++) {
    boardBodyElements[i].classList.add("card-body", "row");
  }

  var kanbanItems = document.getElementsByClassName("kanban-item");
  for (i = 0; i < kanbanItems.length; i++) {
    kanbanItems[i].classList.add("btn", "btn-primary");
  }
}

var addBoard = document.getElementById("createNewColumnButton");
addBoard.addEventListener("click", function () {
  let boardid = (KanbanTest.boardContainer.length + 1).toString();
  KanbanTest.addBoards([
    {
      id: boardid,
      title: $("#newColumnName").val().toString(),
      class: "card-header",
      project_id: projectID,
      item: [],
      created_by: sessionStorage.getItem("userID"),
      projectIndex: projectID,
    },
  ]);
  let container = KanbanTest.boardContainer;
  container[container.length - 1].classList.add("card-body", "row");
  KanbanTest.findBoard(boardid).classList.add("col", "card", "card-primary");

  let raw =
    '{\n  "columnName" : "' + $("#newColumnName").val().toString() + '"\n}';

  let requestOptions = {
    method: "POST",
    headers: {
      Authorization: "Bearer " + sessionStorage.getItem("token"),
      "Content-Type": "text/plain",
    },
    body: raw,
    redirect: "follow",
  };

  fetch(
    "http://127.0.0.1:12345/teams/" +
      projectTeamID +
      "/projects/" +
      projectID +
      "/columns",
    requestOptions
  )
    .then((response) => response.json())
    .then((result) => {
      console.log(result);
      if (result.error.code == 200) {
        $("#newColumnName").value = "";
        addEventListenerToHeader();
      }
    });
});
addEventListenerToHeader();
function addEventListenerToHeader() {
  let editCol = document.getElementsByClassName("kanban-board-header");
  console.log(editCol);
  for (var i = 0; i < editCol.length; i++) {
    editCol[i].children[0].setAttribute("data-toggle", "modal");
    editCol[i].children[0].setAttribute("data-target", "#editColumn");
    editCol[i].addEventListener("click", function () {
      let curentCol = this;
      console.log(this.parentElement.getAttribute("data-id"));
      document
        .getElementById("newColumnNameChange")
        .setAttribute("placeholder", this.children[0].outerText);
      $("#deleteColumn").off();
      $("#deleteColumn").on("click", function () {
        let colID = curentCol.parentElement.getAttribute("data-id");
        KanbanTest.removeBoard(colID);

        let myHeaders = new Headers();
        myHeaders.append("Content-Type", "text/plain");
      
        let raw = '{\n  "columnID" : ' + parseInt(colID) + '\n}';
      
        let deleteColOptions = {
          method: "POST",
          credentials: "omit",
          headers: {
            Authorization: "Bearer " + token,
            "Content-Type": "text/plain",
          },
          body: raw,
          redirect: "follow",
        };
      
        fetch("http://127.0.0.1:12345/teams", deleteColOptions)
          .then((response) => response.json())
          .then((result) => {
            if (result.error.code == 200) {
         
            }
          });

      });
    });
  }
}
