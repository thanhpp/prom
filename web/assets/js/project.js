var KanbanTest = new jKanban({
  element: "#myKanban",
  gutter: "10px",
  widthBoard: "300px",
  responsivePercentage: false,
  itemHandleOptions: {
    enabled: false,
  },
  click: function (el) {
    console.log("Trigger on all items click!");
  },
  context: function (el, e) {
    console.log("Trigger on all items right-click!");
  },
  dropEl: function (el, target, source, sibling) {
    console.log(target.parentElement.getAttribute("data-id"));
    console.log(el, target, source, sibling);
  },
  buttonClick: function (el, boardId) {
    console.log(el);
    console.log(boardId);
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
    content: "+ New Item",
    class: "new-card btn btn-outline-primary",
    footer: false,
  },
  boards: [
    {
      id: "_todo",
      title: "To Do ",
      class: "info,good",
      dragTo: ["_working"],
      item: [
        {
          id: "_test_delete",
          title: "Try drag this (Look the console)",
          drag: function (el, source) {
            console.log("START DRAG: " + el.dataset.eid);
          },
          dragend: function (el) {
            console.log("END DRAG: " + el.dataset.eid);
          },
          drop: function (el) {
            console.log("DROPPED: " + el.dataset.eid);
          },
        },
        {
          title: "Try Click This!",
          click: function (el) {
            alert("click");
          },
          context: function (el, e) {
            alert("right-click at (" + `${e.pageX}` + "," + `${e.pageX}` + ")");
          },
          class: ["peppe", "bello"],
        },
      ],
    },
    {
      id: "_working",
      title: "Working ",
      class: "warning",
      item: [
        {
          title: "Do Something!",
        },
        {
          title: "Run?",
        },
      ],
    },
    {
      id: "_done",
      title: "Done",
      class: "success",
      dragTo: ["_working"],
      item: [
        {
          title: "All right",
        },
        {
          title: "Ok!",
        },
        {
          title: "Ok!",
        },
        {
          title: "Ok!",
        },
      ],
    },
  ],
});

var toDoButton = document.getElementById("addToDo");
toDoButton.addEventListener("click", function () {
  KanbanTest.addElement("_todo", {
    title: "Test Add",
  });
  addClassToNewBoard();
});

var addBoardDefault = document.getElementById("addDefault");
addBoardDefault.addEventListener("click", function () {
  KanbanTest.addBoards([
    {
      id: "_default",
      title: "New Board",
      item: [],
    },
  ]);
  addClassToNewBoard();
});

var removeBoard = document.getElementById("removeBoard");
removeBoard.addEventListener("click", function () {
  KanbanTest.removeBoard("_done");
});

var removeElement = document.getElementById("removeElement");
removeElement.addEventListener("click", function () {
  KanbanTest.removeElement("_test_delete");
});

var allEle = KanbanTest.getBoardElements("_todo");
allEle.forEach(function (item, index) {
  //console.log(item);
});

var test = document.getElementById("test");
test.addEventListener("click", function () {
var board = KanbanTest.findBoard("_done");
console.log(board);
});

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
    boardElements[i].classList.add("col","card");
  }
  
  var boardHeaderElements = document.getElementsByClassName("kanban-board-header");
  for (i = 0; i < boardHeaderElements.length; i++) {
    boardHeaderElements[i].classList.add("card-header");
  }

   
  var boardBodyElements = document.getElementsByClassName("kanban-drag");
  for (i = 0; i < boardBodyElements.length; i++) {
    boardBodyElements[i].classList.add("card-body", "row");
  }

  var kanbanItems = document.getElementsByClassName("kanban-item");
  for (i = 0; i < kanbanItems.length; i++) {
    kanbanItems[i].classList.add("btn", "btn-primary",);
  }
  
}
