let prevlistitem = -1;
let currentlistitem = -1;
let currentlistlength = -1;
let listselect = undefined;
let list = undefined;

document.addEventListener('keyup', (e) => {
    if (e.keyCode === 27) {
        closesnip();
    }
});

function addNodes(data) {
	document.getElementById("myUL").innerHTML = "";
	var ul = document.getElementById("myUL");
	let json = JSON.parse(data);
	
	for (var key in json) {
        let obj = json[key];
        
        var divtxt = document.createElement('div');
        divtxt.id = "data";
        divtxt.textContent = obj.name;
        
        var divhash = document.createElement('div');
        divhash.id = "hash";
        divhash.textContent = obj.hash;

        var li = document.createElement("li");
        li.appendChild(divtxt);
        li.appendChild(divhash);
        
		ul.appendChild(li);
    }
    
    listselect = document.querySelector('#myUL');
    list = listselect.querySelectorAll('li');
    currentlistitem = -1;
    currentlistlength = json.length;
}

function movelist(direction) {
    if (direction == 0) return;
    if (currentlistitem <= -2) return;
    
    prevlistitem = currentlistitem;
    currentlistitem += direction;

    //rolled to far forward, set back to start
    if (currentlistitem >= (currentlistlength) ) {
        currentlistitem = -1;
    }
    if (currentlistitem <= -2) {
        currentlistitem = currentlistlength - 1;
    }

    if (currentlistitem == -1) {
        list[prevlistitem].className = ""    
    }

    if (prevlistitem != -1) {
        list[prevlistitem].className = ""
    }
    list[currentlistitem].className = "selected"
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
            let txt = document.getElementById('myInput').value;
            searchsnip(txt).then(function(result) { 
                return addNodes(result);
            })
            break;
    }
}

function typesnippet() {
    if (currentlistitem != -1) {
        hashid = list[currentlistitem].querySelector('#hash').innerText
        //console.log( hashid );
        //console.log( list[currentlistitem -1]
        writesnip(hashid)
    } else {
        //console.log("search");
    }
}

function saveform() {
    let txttitle = document.getElementById('title').value;
    let txtcode = document.getElementById('code').value;

    savesnip(txttitle, txtcode);
    window.history.back();
}
