
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
        if (e.keyCode >= 48 && e.keyCode <= 57) {
            if ($( '#box-vars' ).css('display') != 'none') {
                argFromClip(e.keyCode)
            }
        }
    }
});
document.addEventListener('keyup', (e) => {
    if (e.keyCode === 27) { //27 = esc key
        snipClose();
    }
});


function argFromClip(keynumber) {
    number = keynumber - 49;
    if (number == -1 ) {
        number = 9;
    }

    if (number >= 0 && number <= 9) {
        if ( $('#var' + number).length ) {
            $.when(
                snipFromClip()
            ).then(
                function (data) {
                    $('#var' + number).val(data) 
                }  
            );
        }
    }
}

function clearform() {
    $('#code').val('');
    $('#title').val('')
    $('#cmdtypselect').prop('selectedIndex', 0);
}

function saveform() {
    let txttitle = document.getElementById('title').value;
    let txtcode = document.getElementById('code').value;
    let txtcmdtyp = document.getElementById('cmdtypselect').value;

    snipSave(txttitle, txtcode, txtcmdtyp);
    $( '#box-addnew' ).hide();
    $( '#searchbox' ).focus();
    clearform();
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
    var hotkey = "";

    $(document).ready(function() {
        let strautofocus = "autofocus";
        $("#argument-list").empty();
        for (var key in args) {
        n = args[key].name;
            if (n != undefined && n.length >= 1) {
                v = args[key].value;
                if (v == undefined) { v = "" }
                if (key == 10) {
                    hotkey = "<u>0</u>.) ";
                } else {
                    let k = parseInt(key, 10) + 1;
                    hotkey = "<u>" + k + "</u>.) ";
                }
                let txtlabel = "<label class='varList' for='var" + key + "'>" + hotkey + n + ":</label><br>";
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


let asyncJob = {
    deferredQueue : [],
    fullName      : function() {
        return this.firstName + " " + this.lastName;
    },
    GotData   : function(Qid, result) {
        dq = this.deferredQueue[Qid];
        dq.resolve(result);
    },
    SendJob : function(query) {
        var dfquery = $.Deferred();

        Qid = uuid();
        snipAsyncSearch(Qid, query); //sending the query and a unique id
        this.deferredQueue[Qid] = dfquery; //we save the query for our snipGotData function
        return dfquery.promise();
    }
};

function uuid() {
    return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
      (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
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
            asyncJob.SendJob(request.term)
        ).then(
            function(data) {
                let list = [];
                if (data == undefined) return;
                let json = JSON.parse(data);
                if (json == null) return;
                response(json);
            }
        );
    },
    minLength: 1,
    delay: 0,
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
    schname = "<div class='searchname'>" + item.value + schtype + "</div>";
    schcmd = "<div class='searchcmd'>" + item.code + "</div>";
    schinfo = "<div class='searchinfo'>" + schname + schcmd + "</div>";
    lstitm = "<div class='listitem-div'>" + schinfo + "</div>";
    return $('<li>')
    .append(lstitm)
    .appendTo(ul);
};
});
