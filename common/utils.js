
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
    if (e.keyCode === 27) { //27 = esc key
        snipClose();
    }
});

function saveform() {
    let txttitle = document.getElementById('title').value;
    let txtcode = document.getElementById('code').value;
    let txtcmdtyp = document.getElementById('cmdtypselect').value;

    snipSave(txttitle, txtcode, txtcmdtyp);
    $( '#box-addnew' ).hide();
    $( '#searchbox' ).focus();
}

function writeFromHash(hash) {
    //make sure we have something as a value
    if (hash != "") {
        snipWrite(''+hash);
    }
}

function runwithvars() {
    let args = [];
    let hash = $( "#searchbox" ).data('hashid');
    let nodes = document.getElementById("argument-list").childNodes
    for (let i=0; i<nodes.length; i++) {
        //get argument name
        let n = $('#var' + i).data('argname');
        //get argument value
        let v = nodes[i].getElementsByTagName("input")[0].value;
        args[i] = {"name": n, "value": v};
    }
    snipWrite(''+hash, JSON.stringify(args));
}

function populateVarsList(item) {
    let args = item.argument;
    $(document).ready(function() {
        let strautofocus = "autofocus";
        $("#argument-list").empty();
        for (var key in args) {
        n = args[key].name;
            if (n != undefined && n.length >= 1) {
                v = args[key].value;
                if (v == undefined) { v = "" }
                let txtlabel = "<label class='varList' for='var" + key + "'>" + n + ":</label><br>";
                let txtbox = "<input class='varList' type='text' id='var" + key + "' value='" + v + "' " + strautofocus + ">";
                $( "#argument-list" ).append("<div>" + txtlabel + txtbox + "</div>")
                $( '#var' + key ).data('argname', n);
            }
        }
    });
}

function getCodeFromArguments(e) {
    snipCodeFromArg().then(function(result) { 
        return function(data) {
            if (data == undefined) return;
            let json = JSON.parse(data);
            if (json == null) return;
            document.getElementById("code").value = json.code;
        } (result);
    })
}


$(function() {
$( "#searchbox" ).keypress(function(event){
    if (event.which == 13) {
        writeFromHash($( "#searchbox" ).data("hashid"));
    }
}).autocomplete({
    autoFocus: true,
    source: function( request, response ) {
        $.when(
            //snipSearchRemote(request, response),
            snipSearch(request.term)
        ).then(
            function(data) {
                let list = [];
                if (data == undefined) return;
                let json = JSON.parse(data);
                if (json == null) return;
                
                for (var key in json) {
                    let obj = json[key];
                    obj.value = obj.name;
                    list.push(
                        obj
                    );
                }
                response(list);
            }
        );
    },
    minLength: 1,
    delay: 100,
    select: function( event, ui ) {
        populateVarsList(ui.item);
        $( "#searchbox" ).data( "hashid", ''+ui.item.hash);
    },
    open: function() {
        $( "#searchbox" ).data('isopen', true);
        $( this ).removeClass( "ui-corner-all" ).addClass( "ui-corner-top" );
    },
    close: function() {
        $( "#searchbox" ).data('isopen', false);
        $( this ).removeClass( "ui-corner-top" ).addClass( "ui-corner-all" );
    }
}).data('ui-autocomplete')._renderItem = function(ul, item) {
    //if cmdtype in bash add class to searchcmdlinux and typename to bash
    //typename = "bash";
    typename = item.cmdtype 
    cmdtype = "searchcmd_" + typename;
    //the above should really be a function
    schtype = "<div class='searchcmdtype " + cmdtype + "'>" + typename + "</div>";
    schname = "<div class='searchname'>" + item.name + schtype + "</div>";
    schcmd = "<div class='searchcmd'>" + item.code + "</div>";
    schinfo = "<div class='searchinfo'>" + schname + schcmd + "</div>";
    lstitm = "<div class='listitem-div'>" + schinfo + "</div>";
    return $('<li>')
    .append(lstitm)
    .appendTo(ul);
};
});
