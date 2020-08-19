$( document ).ready(function() {

    doRequest('{"operation": "get", "value": "all"}');

    $("#btnNew").on("click", function(){
        //show new snip form
        $( '#box-addnew' ).css('display', 'block');
    });

    $("#btnWithArgs").on("click", function(){
        //show runs with args form
        $( '#box-vars' ).css('display', 'block');
    });

    $("#btnRun").on("click", function(){
        //autotype the selected snippet
        doType( $('tr.row-selected') )
        // hash = $( '#searchbox' ).data('hashid');
        // //make sure we have something as a value
        // if (hash != "") {
        //     snipTyper(''+hash); //snipTyper expects a string so we force js using single quotes
        // }
    });

    $("#btnAddNew").on("click", function(){
        //show add new box
        $( '#box-addnew' ).css('display', 'none');
        $( '#searchbox' ).focus();
        clearform();
    });

    $("#btnSaveNew").on("click", function(){
        //save new snippet
        let txttitle = $( '#title' ).val();
        let txtcode = $( '#code' ).text();
        let txtcmdtyp = $( '#cmdtypselect' ).val();
        let txtsummary = $( '#summary' ).val();

        snipSave(txttitle, txtcode, txtcmdtyp, txtsummary);
        $( '#box-addnew' ).hide();
        $( '#searchbox' ).focus();
        clearform();
    });

    $("#btnSaveEdit").on("click", function(){
        //save new snippet
        let txttitle = $( '#title' ).val();
        let txtcode = $( '#code' ).text();
        let txtcmdtyp = $( '#cmdtypselect' ).val();
        let txtsummary = $( '#summary' ).val();
        let hashid = $('#btnSaveEdit').data('updatehash');

        snipUpdate(hashid, txttitle, txtcode, txtcmdtyp, txtsummary);
        $( '#box-addnew' ).hide();
        $( '#searchbox' ).focus();
        clearform();
    });
    $("#btnDelete").on("click", function(){
        //save new snippet
        let txttitle = $( '#title' ).val();
        let hashid = $('#btnSaveEdit').data('updatehash');
        let ans = confirm("Im going to delete " + txttitle);
        if (ans == true) {
            snipDelete(hashid)
        };
        $( '#box-addnew' ).hide();
        $( '#searchbox' ).focus();
        clearform();
    });
    $("#btnCancelVars").on("click", function(){
        $( '#box-vars' ).css('display', 'none');
        $( '#searchbox' ).focus();
    });

    $("#btnOkVars").on("click", function(){
        //run with selected vars
        let args = [];
        let hash = $('tr.row-selected').data('tt-list-item').hash;
        let nodes = document.getElementById("argument-list").childNodes
        for (let i=0; i<nodes.length; i++) {
            //get argument name
            let n = $('#var' + i).data('argname');
            //get argument value
            let v = nodes[i].getElementsByTagName("input")[0].value;
            args[i] = {"name": n, "value": v};
        }
        snipTyper(''+hash, JSON.stringify(args));
    });

    function buildArg() {
        argname = $( "#argname").val();
        argvalue = $( "#argdefault" ).val();
        if (argvalue.length == 0) {
            data = "{:" + argname + ":}";
        } else { 
            data = "{:" + argname + "!" + argvalue + ":}";
        }
        return data;
    }

    $( "#argumentdrag" ).on("dragstart", function() {
        event.dataTransfer.setData("text", buildArg());
    });
    $( "#argumentdrag" ).on("drop", function(ev) {
        ev.preventDefault();
    });
    $( "#btnInsertArg" ).on("click", function() {
        $( "#code" ).append(buildArg());
    });

    document.addEventListener('keydown', (e) => {
        if (e.altKey == true) { //alt key is pressed so these are modifires
            if (e.keyCode == 65 ) { //A - run with args
                document.getElementById("btnWithArgs").click();
            }
            if (e.keyCode == 78 ) { //N - new snippet
                document.getElementById("btnNew").click();
            }
            if (e.keyCode == 69 ) { //E - edit snippet
                populateEditBox( $('tr.row-selected') );
                showEditBox();
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

    function doType(listitem) {
        let item = listitem.data('tt-list-item');

        if (item.hash != "") {
            snipTyper(''+item.hash); //snipTyper expects a string so we force js using single quotes
        }
    }

    function populateVarsList(listitem) {
        let item = listitem.data('tt-list-item');
        let args = item.argument;
        var hotkey = "";
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
    }

    function doRequest(query) {
        $.when(
            asyncJob.SendJob( query )
        ).then(
            function(data) {
                $('#resultstable').empty()
                jsonout = JSON.parse(data)
                if (jsonout.length == "") {
                    return
                } 
                Object.entries(jsonout).forEach(([k,v]) => {
                    item = v;
                    item.hash = k;
                    //item = itemarray
                    typename = item.cmdtype
                    cmdtype = "searchcmd_" + typename;
                    //the above should really be a function
                    schtype = "<div class='searchcmdtype " + cmdtype + "'>" + typename + "</div>";
                    schname = "<div class='searchname'>" + item.value + schtype + "</div>";
                    schcmd = "<div class='searchcmd'>" + item.code + "</div>";
                    schinfo = "<div class='searchinfo'>" + schname + schcmd + "</div>";
                    lstitm = "<div class='listitem-div'><div class='edititem'>E</div>" + schinfo + "</div>";
                    itmout = $( "<tr>" ).append( lstitm ).data('tt-list-item', item);
                    itmout.on("click", function () {
                        //console.log("clock");
                        $('tr.row-selected').removeClass("row-selected")
                        $(this).addClass("row-selected")  
                        populateVarsList($(this))
                    });
        
                    $('#resultstable').append( itmout )
                })
                $(".edititem").on('click', function () {
                    populateEditBox( $(this).parent().parent() );
                    showEditBox();
                });
            }
        );
    }

    function populateEditBox( listitem ) {
        let item = listitem.data('tt-list-item')
        $('#code').text(item.code);
        $('#title').val(item.value);
        $('#summary').val(item.summary);
        $('#cmdtypselect').prop('selectedIndex', 0);
        $('#btnSaveEdit').data('updatehash', item.hash);
    }

    function showEditBox() {
        $( '#btnSaveEdit' ).css('display', 'block');
        $( '#btnDelete' ).css('display', 'block');
        $( '#btnSaveNew' ).css('display', 'none');
        $( '#box-addnew' ).css('display', 'block');
    }

    $('table').keydown(function (e) {
        if (e.which == 13) {
            doType( $('tr.row-selected') )
            return
        }
        if (e.which == 38) {
            // Up Arrow
            currItm = $('tr.row-selected').prev();
            if (currItm.length == 0) {
                currItm = $('tr').last();
            }
        } else if (e.which == 40) {
            // Down Arrow
            currItm = $('tr.row-selected').next();
            if (currItm.length == 0) {
                currItm = $('tr').first();
            }
        }
        $('tr.row-selected').removeClass("row-selected")
        if (currItm.length == 0) {
            currItm = $('tr').first();
        }
        currItm.addClass("row-selected")
        //deal with scrolling
        let rowTop = currItm.position().top;
        let rowHeight = currItm.height();
        let rowBot = rowTop + currItm.height();
        
        let tblHeight = $('#resultstable').height();
        if (rowTop <= 0 ) {
            //row is out of bounds, need to move it back into view
            $('#resultstable').scrollTop(currItm[0].offsetTop-6)
        } else if (rowBot >= tblHeight) {
            $('#resultstable').scrollTop(currItm[0].offsetTop + rowHeight - tblHeight - 5);
        }
        populateVarsList($('tr.row-selected'))
        
    });

    $( '#searchbox' ).on('keyup', function (e) {
        doRequest( '{"operation": "' + $('#searchfor').val() + '", "value": "' + $('#searchbox').val() + '"}' );

        if (e.which == 40) {
            // Down Arrow
            $('tr.row-selected').removeClass("row-selected");
            $('#codelist').focus();
            currItm = $('tr').first();
            currItm.addClass("row-selected");
        }
    })
});

let asyncJob = {
    deferredQueue : [],
    GotData   : function(Qid, result) {
        dq = this.deferredQueue[Qid];
        dq.resolve(result);
    },
    SendJob : function(query) {
        var dfquery = $.Deferred();
        Qid = uuid();
        snipAsyncRequest(Qid, query); //sending the query and a unique id
        this.deferredQueue[Qid] = dfquery; //we save the query for our snipGotData function
        return dfquery.promise();
    }
};

function uuid() {
    return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}

//called from go code, part of func newfromcommand()
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

function clearform() {
    $( '#code' ).text('');
    $( '#title' ).val('')
    $( '#summary' ).val('');
    $( '#cmdtypselect' ).prop('selectedIndex', 0);
    $( '#btnSaveNew' ).css('display', 'block');
    $( '#btnSaveEdit' ).css('display', 'none');
    $( '#btnDelete' ).css('display', 'none');
}
