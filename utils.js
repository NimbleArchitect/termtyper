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
        }
    }
});

function addNodes(data) {
    if (data.length <=0) return
    
    document.getElementById("myUL").innerHTML = "";
    let ul = document.getElementById("myUL");
    let json = JSON.parse(data);
    if (json == null) return 
	
	for (var key in json) {
        let obj = json[key];

        itm = buildsearchitem(obj.hash, obj.name);
		ul.appendChild(itm);
    }
    
    listselect = document.querySelector('#myUL');
    list = listselect.querySelectorAll('li');
    currentlistitem = -1;
    currentlistlength = json.length;
}

function buildsearchitem(hash, name) {
    var li = document.createElement("li");

    var divtxt = document.createElement('div');
    divtxt.id = "data";
    divtxt.textContent = name;
    li.appendChild(divtxt);
    
    var divhash = document.createElement('div');
    divhash.id = "hash";
    divhash.textContent = hash;
    li.appendChild(divhash);
    
    return li;
}

function movelist(direction) {
    if (direction == 0) return;
    if (currentlistitem <= -2) return;
    
    prevlistitem = currentlistitem;
    currentlistitem += direction;
    
    boundpos = document.getElementById('searchcombo').getBoundingClientRect()
    pos = list[currentlistitem].getBoundingClientRect();

    //rolled to far forward, set back to start
    if (currentlistitem >= (currentlistlength) ) {
        currentlistitem = -1;
        document.getElementById('searchcombo').scrollTop = boundpos.top;
    }
    if (currentlistitem <= -2) {
        currentlistitem = currentlistlength - 1;
        document.getElementById('searchcombo').scrollTop = boundpos.height;
    }
    
    //goes down but not up
    if (pos != undefined) {
        if ((boundpos.height + boundpos.top) < (pos.height + pos.top)) {
            list[prevlistitem].scrollIntoView();
        }
    }

    if (currentlistitem == -1) {
        list[prevlistitem].className = "";
    }

    if (prevlistitem != -1) {
        list[prevlistitem].className = "";
    }
    list[currentlistitem].className = "selected";
}

function searchFor() {
    switch (event.keyCode) {

        case 38: // up
            //console.log("up")
            movelist(-1);
            //console.log(list[currentlistitem].innerText)
            break;
        
        case 40: // down
            //console.log("down")
            movelist(1);
            //console.log(list[currentlistitem].innerText)
            break;
            
        case 13:
            event.preventDefault();
            typesnippet()
            break;
            
        default:
            let txt = document.getElementById('searchbox').value;
            searchsnip(txt).then(function(result) { 
                return addNodes(result);
            })
            break;
    }
}

function typesnippet() {
    if (currentlistitem != -1) {
        hashid = list[currentlistitem].querySelector('#hash').innerText
        writesnip(hashid)
    } else {
        //console.log("search");
    }
}

function saveform() {
    let txttitle = document.getElementById('title').value;
    let txtcode = document.getElementById('code').value;

    savesnip(txttitle, txtcode);
    document.getElementById('box-addnew').style.display=''
}
