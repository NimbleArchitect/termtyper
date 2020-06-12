
const autoclick = true;
const inlineStatus = true;
const statusFooter = false;

let prevlistitem = -1;
let currentlistitem = -1;
let currentlistlength = -1;
let listselect = undefined;
let list = undefined;
let searchresults = [];


document.addEventListener('keydown', (e) => {
    if (e.altKey == true) { //alt key is pressed so these are modifires
        if (e.keyCode == 65 ) { //A - run with args
            document.getElementById("btnWithArgs").click();
        }
        if (e.keyCode == 78 ) { //N - new snippet
            document.getElementById("btnNew").click();
        }
        if (e.keyCode == 69 ) { //E - edit snippet

        }
    }
});
document.addEventListener('keyup', (e) => {
    if (e.keyCode === 27) {
        // let txt = document.getElementById('searchbox').value;
        snipClose();
    }
});

function autocomplete(inp) {
    /*the autocomplete function takes two arguments,
    the text field element and an array of possible autocompleted values:*/
    var currentFocus;
    var hashid = -1;
    var isReady = false;

    /*execute a function when someone writes in the text field:*/
    inp.addEventListener("input", function(e) {
        snipSearch(this.value).then(function(result) { 
            return inp.getSearchList(result);
        })
    });
    /*execute a function presses a key on the keyboard:*/
    inp.addEventListener("keydown", function(e) {
        var x = document.getElementById(this.id + "autocomplete-list");
        if (x) x = x.getElementsByClassName("autocomplete-items-div");
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

            if (countOpenLists() == 0) {
                writeFromHash(hashid);
            }

            if (currentFocus > -1) {
                /*and simulate a click on the "active" item:*/
                if (x) x[currentFocus].click();
            } else if (currentFocus == -1) {
                if (x != null && x.length == 1) {
                    /*insert the value for the autocomplete text field:*/
                    // // /*close the list of autocompleted values,
                    // // (or any other open lists of autocompleted values:*/
                    if (autoclick == true) {
                        currentFocus++
                        x[0].click();
                    }
                }
            }
        }
    });
    inp.getSearchList = function(data) {
        let list = [];
        if (data == undefined) return;
        let json = JSON.parse(data);
        if (json == null) return;
        
        for (var key in json) {
            let obj = json[key];
            list.push(
                obj
            );
        }
        inp.populateSearch(list);

    }
    inp.populateSearch = function(arr) {
        var a, b, i, val = this.value;
        /*close any already open lists of autocompleted values*/
        closeAllLists();
        //searchresults.length = 0;
        if (!val) { return false;}
        currentFocus = -1;
        /*create a DIV element that will contain the items (values):*/
        a = document.createElement("DIV");
        a.setAttribute("id", this.id + "autocomplete-list");
        a.setAttribute("class", "autocomplete-items");
        /*append the DIV element as a child of the autocomplete container:*/
        this.parentNode.appendChild(a);
        for (i = 0; i < arr.length; i++) {
            searchresults[arr[i].hash] = arr[i];
            /*create a DIV element for each matching element:*/
            b = document.createElement("DIV");
            b.setAttribute("class", "autocomplete-items-div");
            b.innerHTML += arr[i].name;
            if (inlineStatus == true) {
                schinfo = "<div class='searchcmd'>" + arr[i].code + "</div>";
                if (arr[i].tags != undefined) {
                    schinfo += "<div class='searchtag'>tag: </div>";
                }
                b.innerHTML += "<div class='searchinfo'>" + schinfo + "</div>";
            }
            /*insert a input field that will hold the current array item's value:*/
            b.innerHTML += "<input type='hidden' value='" + arr[i].name + "'>";
            b.innerHTML += "<input type='hidden' id='hash' value='" + arr[i].hash + "'>";
            /*execute a function when someone clicks on the item value (DIV element):*/
            b.addEventListener("click", function(e) {
                /*insert the value for the autocomplete text field:*/
                inp.value = this.getElementsByTagName("input")[0].value;
                hashid = this.children.namedItem("hash").value;
                populateArgumentsList(hashid);
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
        hashid = x[currentFocus].children.namedItem("hash").value;
        if (statusFooter == true) {
            document.getElementById('cmd2run').innerText = searchresults[hashid].code;
        }
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
    function countOpenLists(elmnt) {
        let counter = 0;
        var x = document.getElementsByClassName("autocomplete-items");
        for (var i = 0; i < x.length; i++) {
            if (elmnt != x[i] && elmnt != inp) {
              counter++;
            }
        }
        return counter;
    }
  }

function saveform() {
    let txttitle = document.getElementById('title').value;
    let txtcode = document.getElementById('code').value;

    snipSave(txttitle, txtcode);
    document.getElementById('box-addnew').style.display='';
}

function writeFromHash(hash) {
    if (searchresults[hash].argument == null) {
        snipWrite(hash);
    } else { // if (searchresults[hash].argument.length <= 0) {
        snipWrite(hash);
    }
}

function populateArgumentsList(hashid) {
//let argumentList = {
    hash = -1;
    args = [];
    container = undefined;

    populateVarsList(hashid);

    function populateVarsList(hash) {
        this.args = searchresults[hash].argument;
        this.hash = hash;
        this.container = document.getElementById("variable-list");

        this.container.innerHTML = "";
        // a = document.createElement("DIV");
        // a.setAttribute("id", this.id + "");
        // a.setAttribute("class", "autocomplete-items");
        /*append the DIV element as a child of the autocomplete container:*/
        // this.container.appendChild(a);
        let strautofocus = "autofocus";
        for (var key in this.args) {
            /*create a DIV element for each matching element:*/
            n = this.args[key].name;
            if (n != undefined && n.length >= 1) {
                v = this.args[key].value;
                if (v == undefined) { v = "" }
                b = document.createElement("DIV");
                b.innerHTML += "<label class='varList' for='var" + key + "'>" + n + ":</label><br>";
                b.innerHTML += "<input class='varList' type='text' id='var" + key + "' value='" + v + "' " + strautofocus + ">";
                // a.appendChild(b);
                strautofocus = "";
                b.addEventListener("keydown", function(e) {
                    if (e.keyCode == 13) {
                        /*If the ENTER key is pressed, prevent the form from being submitted,*/
                        e.preventDefault();
                        let iid = document.activeElement.id;
                        let nodes = document.getElementById("variable-list").childNodes
                        for (let i=0; i<nodes.length; i++) {
                            let nid = nodes[i].getElementsByTagName("input").item("").id;
                            if (nid == iid) {
                                if (i + 1 >= nodes.length) {
                                    document.getElementById("btnOkVars").focus();
                                } else {
                                    nodes[i + 1].getElementsByTagName("input")[0].focus();
                                }
                            } 
                        }
                    }
                });
                this.container.appendChild(b);
            }
        }
    }

    document.getElementById("btnOkVars").addEventListener("click", function(e) {
        // let nodes = this.container.childNodes;
        let nodes = document.getElementById("variable-list").childNodes
        for (let i=0; i<nodes.length; i++) {
            args[i].value = nodes[i].getElementsByTagName("input")[0].value;
        }
        snipWrite(hash, JSON.stringify(args));
    });
}
