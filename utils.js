
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
        divtxt.textContent = obj.txt;
        
        var divhash = document.createElement('div');
        divhash.id = "hash";
        divhash.textContent = obj.id;

        var li = document.createElement("li");
        li.appendChild(divtxt);
        li.appendChild(divhash);
        
		ul.appendChild(li);
    }
    
    listselect = document.querySelector('#myUL');
    list = listselect.querySelectorAll('li');
    currentlistitem = 0;
    currentlistlength = json.length;
}

function movelist(direction) {
    if (currentlistitem == -1) return;

    if (currentlistitem == currentlistlength) {
        if (direction == 1) currentlistitem = -1;
    } else if (currentlistitem == 0) {
        if (direction == -1) currentlistitem = currentlistlength + 1;
    }

    currentlistitem += direction;
}

function searchFor() {
    switch (event.keyCode) {

        case 38: // up
            console.log("up")
            movelist(-1);
            break;
        
        case 40: // down
            console.log("down")
            movelist(1);
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
    if (currentlistitem != 0) {
        hashid = list[currentlistitem -1].querySelector('#hash').innerText
        console.log( hashid );
        //console.log( list[currentlistitem -1]
        writesnip("")
    } else {
        //console.log("search");
    }
}