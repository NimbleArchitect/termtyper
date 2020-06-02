let prevlistitem = -1;
let currentlistitem = -1;
let currentlistlength = -1;
let listselect = undefined;
let list = undefined;

document.addEventListener('keyup', (e) => {
    if (e.keyCode === 27) {
        let txt = document.getElementById('searchbox').value;
        if (txt.length == 0) {
            closesnip();
        } else {
            document.getElementById('searchbox').value = "";
            //inp.closeAllLists();
        }
    }
});

// function getsearchlist(data) {
//     searchsnip(data).then(function(result) { 
//         return addtolist(result);
//     })
// }
// function addtolist(data) {
//     let json = JSON.parse(data);
//     if (json == null) return;
    
//     let list = [];
// 	for (var key in json) {
//         let obj = json[key];

//         itm.name = obj.name;
//         itm.hash = obj.name;
        
//         list.push(itm);
//     }
//     return list;
// }

function autocomplete(inp) {
    /*the autocomplete function takes two arguments,
    the text field element and an array of possible autocompleted values:*/
    var currentFocus;
    var arr;
    /*execute a function when someone writes in the text field:*/
    inp.addEventListener("input", function(e) {
        searchsnip(this.value).then(function(result) { 
            return inp.populateList(result);
        })
    });
    /*execute a function presses a key on the keyboard:*/
    inp.addEventListener("keydown", function(e) {
        var x = document.getElementById(this.id + "autocomplete-list");
        if (x) x = x.getElementsByTagName("div");
        if (e.keyCode == 40) {
          /*If the arrow DOWN key is pressed,
          increase the currentFocus variable:*/
          currentFocus++;
          /*and and make the current item more visible:*/
          addActive(x);
        } else if (e.keyCode == 38) { //up
          /*If the arrow UP key is pressed,
          decrease the currentFocus variable:*/
          currentFocus--;
          /*and and make the current item more visible:*/
          addActive(x);
        } else if (e.keyCode == 13) {
          /*If the ENTER key is pressed, prevent the form from being submitted,*/
          e.preventDefault();
          if (currentFocus > -1) {
            /*and simulate a click on the "active" item:*/
            if (x) x[currentFocus].click();
          }
        }
    });
    inp.populateList = function(data) {
        let list = [];
        if (data == undefined) return;
        let json = JSON.parse(data);
        if (json == null) return;
        
        for (var key in json) {
            let obj = json[key];
            list.push({
                hash: obj.hash,
                name: obj.name
              });
        }
        inp.populateSearch(list);
    }
    inp.populateSearch = function(arr) {
        var a, b, i, val = this.value;
        /*close any already open lists of autocompleted values*/
        closeAllLists();
        if (!val) { return false;}
        currentFocus = -1;
        /*create a DIV element that will contain the items (values):*/
        a = document.createElement("DIV");
        a.setAttribute("id", this.id + "autocomplete-list");
        a.setAttribute("class", "autocomplete-items");
        /*append the DIV element as a child of the autocomplete container:*/
        this.parentNode.appendChild(a);
        for (i = 0; i < arr.length; i++) {
            /*create a DIV element for each matching element:*/
            b = document.createElement("DIV");
            b.innerHTML += arr[i].name;
            /*insert a input field that will hold the current array item's value:*/
            b.innerHTML += "<input type='hidden' value='" + arr[i].name + "'>";
            b.innerHTML += "<input type='hidden' id='hash' value='" + arr[i].hash + "'>"
            /*execute a function when someone clicks on the item value (DIV element):*/
            b.addEventListener("click", function(e) {
                /*insert the value for the autocomplete text field:*/
                inp.value = this.getElementsByTagName("input")[0].value;
                /*close the list of autocompleted values,
                (or any other open lists of autocompleted values:*/
                closeAllLists();
            });
            a.appendChild(b);
        }
    }
    function addActive(x) {
      /*a function to classify an item as "active":*/
      if (!x) return false;
      /*start by removing the "active" class on all items:*/
      removeActive(x);
      if (currentFocus >= x.length) currentFocus = 0;
      if (currentFocus < 0) currentFocus = (x.length - 1);
      /*add class "autocomplete-active":*/
      x[currentFocus].classList.add("autocomplete-active");
    }
    function removeActive(x) {
      /*a function to remove the "active" class from all autocomplete items:*/
      for (var i = 0; i < x.length; i++) {
        x[i].classList.remove("autocomplete-active");
      }
    }
    function closeAllLists(elmnt) {
      /*close all autocomplete lists in the document,
      except the one passed as an argument:*/
      var x = document.getElementsByClassName("autocomplete-items");
      for (var i = 0; i < x.length; i++) {
        if (elmnt != x[i] && elmnt != inp) {
          x[i].parentNode.removeChild(x[i]);
        }
      }
    }
    /*execute a function when someone clicks in the document:*/
    document.addEventListener("click", function (e) {
        writesnip(e.target.children.namedItem("hash").value);
        //closeAllLists(e.target);

    });
  }


// function addNodes(data) {
//     if (data.length <=0) return
    
//     document.getElementById("myUL").innerHTML = "";
//     let ul = document.getElementById("myUL");
//     let json = JSON.parse(data);
//     if (json == null) return 
	
// 	for (var key in json) {
//         let obj = json[key];

//         itm = buildsearchitem(obj.hash, obj.name);
// 		ul.appendChild(itm);
//     }
    
//     listselect = document.querySelector('#myUL');
//     list = listselect.querySelectorAll('li');
//     currentlistitem = -1;
//     currentlistlength = json.length;
// }

// function buildsearchitem(hash, name) {
//     var li = document.createElement("li");

//     var divtxt = document.createElement('div');
//     divtxt.id = "data";
//     divtxt.textContent = name;
//     li.appendChild(divtxt);
    
//     var divhash = document.createElement('div');
//     divhash.id = "hash";
//     divhash.textContent = hash;
//     li.appendChild(divhash);
    
//     return li;
// }

// function movelist(direction) {
//     if (direction == 0) return;
//     if (currentlistitem <= -2) return;
    
//     prevlistitem = currentlistitem;
//     currentlistitem += direction;
    
//     boundpos = document.getElementById('searchcombo').getBoundingClientRect()
//     pos = list[currentlistitem].getBoundingClientRect();

//     //rolled to far forward, set back to start
//     if (currentlistitem >= (currentlistlength) ) {
//         currentlistitem = -1;
//         document.getElementById('searchcombo').scrollTop = boundpos.top;
//     }
//     if (currentlistitem <= -2) {
//         currentlistitem = currentlistlength - 1;
//         document.getElementById('searchcombo').scrollTop = boundpos.height;
//     }
    
//     //goes down but not up
//     if (pos != undefined) {
//         if ((boundpos.height + boundpos.top) < (pos.height + pos.top)) {
//             list[prevlistitem].scrollIntoView();
//         }
//     }

//     if (currentlistitem == -1) {
//         list[prevlistitem].className = "";
//     }

//     if (prevlistitem != -1) {
//         list[prevlistitem].className = "";
//     }
//     list[currentlistitem].className = "selected";
// }

// function searchFor() {
//     switch (event.keyCode) {

//         case 38: // up
//             //console.log("up")
//             movelist(-1);
//             //console.log(list[currentlistitem].innerText)
//             break;
        
//         case 40: // down
//             //console.log("down")
//             movelist(1);
//             //console.log(list[currentlistitem].innerText)
//             break;
            
//         case 13:
//             event.preventDefault();
//             typesnippet()
//             break;
            
//         default:
//             let txt = document.getElementById('searchbox').value;
//             searchsnip(txt).then(function(result) { 
//                 return addNodes(result);
//             })
//             break;
//     }
// }

// function typesnippet() {
//     if (currentlistitem != -1) {
//         hashid = list[currentlistitem].querySelector('#hash').innerText
//         writesnip(hashid)
//     } else {
//         //console.log("search");
//     }
// }

function saveform() {
    let txttitle = document.getElementById('title').value;
    let txtcode = document.getElementById('code').value;

    savesnip(txttitle, txtcode);
    document.getElementById('box-addnew').style.display=''
}
