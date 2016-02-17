function main(body) {
  loggedIn(function() {
    $.loadCSS('static/css/index.css');
    
    buildTemplate();
    buildHome(); 
  });
} 

// decisionListByCategory expect a callback cb(inprogress,completed)
// contains a list of decisions in progress or completed. for the
// currently logged in user.
function decisionListByCategory(cb) {
  get_text("/whoami", function (person) {
    get_text("/person/"+person['person_id']+"/decisions", function (decisions) {
      var inprogress = [];
      var completed  = [];
      for(var i in decisions["decisions"]) {
        d = decisions["decisions"][i];
        if(d["stage"] < 3) {
          inprogress.push(d);
        } else {
          completed.push(d);
        }
      }
      cb(inprogress, completed);
    });
  });
}

function updateLeftNav() {
    decisionListByCategory(function(inprogress, completed) {

    $("#userMenu3").html("");
    $("#userMenu4").html("");

    for(var i in inprogress) {
      dname = inprogress[i]["name"];
      did   = inprogress[i]["decision_id"];
      $("#userMenu3").append("<li><a onclick=\"buildDecisionHome("+did+")\"><i class=\"glyphicon glyphicon-list-alt\"></i> "+dname+" </a></li>");
    }

    for(var i in completed) {
      dname = completed[i]["name"];
      did   = completed[i]["decision_id"];
      $("#userMenu4").append("<li><a onclick=\"buildDecisionHome("+did+")\"><i class=\"glyphicon glyphicon-list-alt\"></i> "+dname+" </a></li>");
    }

  });
}

function buildTemplate() {
  //Add alert box used in scripts
  $('body').append(
    '<div id="modalConfirmYesNo" class="modal fade">' +
        '<div class="modal-dialog">' +
            '<div class="modal-content">' +
                '<div class="modal-header">' +
                    '<button type="button" ' +
                    'class="close" data-dismiss="modal" aria-label="Close">' +
                        '<span aria-hidden="true">&times;</span>' +
                    '</button>' +
                    '<h4 id="lblTitleConfirmYesNo" class="modal-title">Confirmation</h4>' +
                '</div>' +
                '<div class="modal-body">' +
                    '<p id="lblMsgConfirmYesNo"></p>' +
                '</div>' +
                '<div class="modal-footer">' +
                    '<button id="btnYesConfirmYesNo" type="button" class="btn btn-primary">Yes</button>' +
                    '<button id="btnNoConfirmYesNo" type="button" class="btn btn-default">No</button>' +
                '</div>' +
            '</div>' +
        '</div>' +
    '</div>'
  );


  //nav section
  var nav = $('<nav>').addClass('navbar navbar-inverse navbar-fixed-top').appendTo('body')
  var div_container = $('<div class="container-fluid">').appendTo(nav)
  var div_nav_header = $('<div class="navbar-header">').appendTo(div_container)
  
  var button_nav = $([
      '<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#myNavbar">',
            '<span class="icon-bar"></span>',
            '<span class="icon-bar"></span>',
            '<span class="icon-bar"></span>',
          '</button>'
  ].join('\n'));
          
  div_nav_header.append(button_nav) 
  div_nav_header.append($('<a onclick="buildHome()"class = "navbar-brand"><img id="logo" src="../static/images/logo.png">'))
  
  var div_collapse = $('<div class="collapse navbar-collapse" id="myNavbar">')
  
  var nav_ul1 = $('<ul class="nav navbar-nav"><li> <a onclick="buildHome()">Dashboard</a></li></ul>')
  var nav_ul2 = $([
  '<ul class="nav navbar-nav navbar-right">',
    '<li class="dropdown">',
      '<a class="dropdown-toggle" role="button" data-toggle="dropdown" aria-expanded="false"><i class="glyphicon glyphicon-user"></i><span id="userName">' + '</span><span class="caret"></span></a>',
        '<ul id="g-account-menu" class="dropdown-menu" role="menu">',
          '<li><a onclick="buildEditUser()">Edit Profile</a></li>',
        '</ul>',
    '</li>',
    '<li><a href="/logout.html"><i class="glyphicon glyphicon-lock"></i> Logout</a></li>',
  '</ul>'
  ].join("\n"));
        
        
  div_collapse.append(nav_ul1)
  div_collapse.append(nav_ul2)
  div_container.append(div_collapse)
  
  //dashboard section
  
  var div_dashboard = $('<div class="container-fluid" style="margin-top:75px;">').appendTo('body')
  
  var div_row = $('<div class="row-fluid">')
  
  var nav_section = $([
      '<div class="col-sm-3" >',
          '<strong><i class="glyphicon glyphicon-wrench"></i> Tools</strong>',
        '<hr>',
        '<ul class="nav nav-stacked">',
        '<li class="nav-header"> <a  data-toggle="collapse" data-target="#userMenu" aria-expanded="false" class="collapsed">Decisions <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
          '<ul class="nav nav-stacked collapse" id="userMenu" aria-expanded="false" style="height: 0px;">',
            '<li class="active"> <a onclick="buildCreateDecision()"><i class="glyphicon glyphicon-asterisk"></i> New Decision</a></li>',
            '<li class="nav-header"> <a  data-toggle="collapse" data-target="#userMenu3" aria-expanded="false" class="collapsed">In Progress <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
              '<ul class="nav nav-stacked collapse" id="userMenu3" aria-expanded="false" style="height: 0px;">',
              '</ul>',
            '</li>',
            '<li class="nav-header"> <a  data-toggle="collapse" data-target="#userMenu4" aria-expanded="false" class="collapsed">Completed <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
              '<ul class="nav nav-stacked collapse" id="userMenu4" aria-expanded="false" style="height: 0px;">',
              '</ul>',
            '</li>',
          '</ul>',
        '</li>',
        '<li class="nav-header"> <a  data-toggle="collapse" data-target="#Menu2" aria-expanded="false" class="collapsed">Settings <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
          '<ul class="nav nav-stacked collapse" id="Menu2" aria-expanded="fasle">',
              '<li class="active"> <a onclick="buildHome()"><i class="glyphicon glyphicon-home"></i> Home</a></li>',              
              '<li><a onclick="buildAddUser()"><i class="glyphicon glyphicon-user"></i> New User</a></li>',
              //'<li><a ><i class="glyphicon glyphicon-list-alt"></i> Report</a></li>',
              '<li><a href="/logout.html"><i class="glyphicon glyphicon-off"></i> Logout</a></li>',
          '</ul>',
        '</li>',
        '</ul>',
        '<hr>',
        
      '</div>'
    ].join('\n'));

  // Update the navbar portion with decisions
  updateLeftNav();
  
  var display_section = $('<div class="col-sm-9" id="content">');
  
  div_row.append(nav_section);
  div_row.append(display_section)
  div_dashboard.append(div_row)

  updateUserText();
  
  $("a").attr("aria-expanded","true");
  $("a").click(function(){
    $(this).find('i#arrow_change').toggleClass('glyphicon-chevron-right glyphicon-chevron-down');
  });
}

function clearContent() {
  $('#content').empty();
}

function updateUserText() {
  get_text("/whoami", function (result) {
    get_text("/person/"+result['person_id']+"/info", function (result) {
      $('#userName').text(' ' + result['person']['name_first'] + ' ' + result['person']['name_last']);
    });
  });
}

function buildHome() {
  $('title').html('Decision Home');
  
  clearContent();

  // The progress bar in here are meaning less ?
  // remove them
  var display_section = $([
        '<div>',
        '<strong><i class="glyphicon glyphicon-dashboard"></i> My Dashboard</strong>',
        '<hr>',
        
        '<div class="jumbotron">',
          '<h2>Welcome to the Decision Dashboard</h2>',
          '<img src="http://blog.rameshganapathy.com/wp-content/uploads/2014/03/calvin-knowledge-is-paralyzing.gif" alt="Comic"/>',
        '</div>',
        '<hr>',
        '<div class= "panel panel-default">',
          '<div class="panel-heading">',
            '<h4>Report</h4>',
          '</div>',
          '<div class="panel-body">',
            '<small>Decisions completed</small>',
            
            '<div class="progress">',
                                '<div id="completed_length" class="progress-bar progress-bar-success" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="100" style="width: 5%">',
                                '</div>',
                        '</div>',
            '<small>Decisions in progress</small>',
            '<div class="progress">',
                                '<div id="inprogress_length" class="progress-bar progress-bar-success" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="100" style="width: 5%">',
                                '</div>',
                            '</div>',
          '</div>',
        '</div>',
        
        '<div class="list-group" id="inprogress_list">',
        '</div>',
        '<div class="list-group" id="completed_list">',
        '</div>',
        '</div>'
  ].join('\n'));
  
  display_section.appendTo('#content');

  // Update the navbar portion with decisions
  decisionListByCategory(function(inprogress, completed) {

    $("#inprogress_list").html("<a class=\"list-group-item active\">Decisions In Progress</a>");
    $("#completed_list").html("<a class=\"list-group-item active\">Decisions Completed</a>");

    for(var i in inprogress) {
      dname = inprogress[i]["name"];
      did   = inprogress[i]["decision_id"];
      $("#inprogress_list").append("<a  onclick=\"buildDecisionHome("+did+")\" class=\"list-group-item\">" + dname + "</a>");
    }

    for(var i in completed) {
      dname = completed[i]["name"];
      did   = completed[i]["decision_id"];
      $("#completed_list").append("<a  onclick=\"buildDecisionHome("+did+")\" class=\"list-group-item\">" + dname + "</a>");
    }

    // Set progressbar lengths
    totalLength = inprogress.length + completed.length;
    inprogress_progress = ((inprogress.length / totalLength) * 100); 
    completed_progress  = ((completed.length  / totalLength) * 100);
    $("#inprogress_length").width(inprogress_progress + "%");
    $("#completed_length").width(completed_progress + "%");

  });
  
          
}

/**** Edit Profile ****/
  function buildEditUser() {
    $('title').html('Update User!');

    clearContent();
    
    $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Person</strong><hr/>').appendTo('#content');
    
    var wrapper = $('<div>').css('max-width','500px').appendTo('#content');
    var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return updateUser()').appendTo(wrapper);
      $('<h3>').addClass('form-signin-heading').text('Update Details').appendTo(form);
      
      $('<div>').attr('id','success').addClass('alert alert-success').appendTo(form);
      $('#success').hide();
      
      $('<input type="text" />').addClass('form-control')
        .attr('name', 'username')
        .attr('placeholder', 'Email Address')
        .attr('id', 'username')
        .appendTo(form);
      
      $('<input type="text" />').addClass('form-control')
        .attr('name', 'firstname')
        .attr('placeholder', 'First Name')
        .attr('id', 'firstname')
        .appendTo(form);
      
      $('<input type="text" />').addClass('form-control')
        .attr('name', 'lastname')
        .attr('placeholder', 'Last Name')
        .attr('id', 'lastname')
        .appendTo(form);
      
      $('<h3>').addClass('form-signin-heading').text('Password').appendTo(form);
      
      $('<div>').attr('id','error').addClass('alert alert-danger').appendTo(form);
      $('#error').hide();
      
      $('<input type="password" />').addClass('form-control')
        .attr('name', 'password')
        .attr('placeholder', 'Password')
        .attr('id', 'password')
        .appendTo(form);
        
      $('<input type="password" />').addClass('form-control')
        .attr('name', 'password2')
        .attr('placeholder', 'Password')
        .attr('id', 'password2')
        .appendTo(form);

      $('<hr/>').appendTo(form);
      $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('type', 'submit').text('Submit').appendTo(form);
      
      get_text("/whoami", function (result) {
        get_text("/person/"+result['person_id']+"/info", function (result) {
          $('#username').val(result['person']['email']);
          $('#firstname').val(result['person']['name_first']);
          $('#lastname').val(result['person']['name_last']);
        });
      });
  }

  function updateUser() {
  $('#error').hide();
  
  if(document.getElementById('password').value != document.getElementById('password2').value) {
    $('#error').html('<b>Error:</b> Passwords do not match!');
    $('#error').show();
  }
  else {
    if(document.getElementById('password').value == "") {
      new_info = {
        "email":document.getElementById('username').value,
        "name_first":document.getElementById('firstname').value,
        "name_last":document.getElementById('lastname').value
        }
    }
    else {
      new_info = {
          "email":document.getElementById('username').value,
          "pw_hash":document.getElementById('password').value,
          "name_first":document.getElementById('firstname').value,
          "name_last":document.getElementById('lastname').value
          }
    }
    
    get_text("/whoami", function (result) {
      put_text("/person/"+result['person_id'], JSON.stringify(new_info), function (result) {
        //alert(JSON.stringify(result));
        document.getElementById('password').value = "";
        document.getElementById('password2').value = "";
        $('#success').html('<b>Update successful!</b>');
        $('#success').show();
        updateUserText();
      });
    });
  }
  
  return false;
  }

/**** Add User ****/
  function buildAddUser(){
    $('title').html('Add User!');
    clearContent();
    
    $('<strong><i class="glyphicon glyphicon-cog"></i> Add User</strong><hr/>').appendTo('#content');
    //build the sign up form
    var signupForm = $('<form id="myform" onsubmit="return false"></form>')
    $('<div id= "signup_successful" class = "alert alert-success"></div>').appendTo("#content")
    $('<div id="signup_error" class="alert alert-danger" style="display: none;"></div>').appendTo("#content")
    
    var email = $('<label> Email<input type="text" name ="email" id="emailInput" required="required" placeholder ="Email" class="form-control"/> </label></br>')
          .appendTo('body');
          
    var pwd = $('<label> Password<input type="password" name ="pw_hash" id="passwordInput" required="required" placeholder="Password" class="form-control"/> </label></br>')
          .appendTo('body');
          
    var firstname = $('<label> First Name<input type="text" name ="name_first" id="firstnameInput" required="required" placeholder="First Name" class="form-control"/> </label></br>')
          .appendTo('body');
    
    var lastname = $('<label>Last Name<input type="text" name="name_last" id="lastnameInput" required="required" placeholder="Last Name" class="form-control" /> </label></br>')
          .appendTo('body');
          
    var submit = $(' <input type="button" onclick="addUserSubmitform()" value="Sign Up" class="btn btn-sm-9 btn-primary" />')
          .appendTo('body');
          

    signupForm.append(email,pwd, firstname,lastname,submit);
    signupForm.appendTo('#content');
    $("#signup_successful").hide()
    $("#signup_error").hide()
    

  }

  function addUserSubmitform(){   
    $("#signup_successful").hide()
    $("#signup_error").hide()

    if(document.getElementById("emailInput").value =='') {
      $('#signup_error').html('<b>Error:</b> No email set!');
      $('#signup_error').show();
    }
    else if(document.getElementById("passwordInput").value =='') {
      $('#signup_error').html('<b>Error:</b> No password set!');
      $('#signup_error').show();
    }else if(document.getElementById("firstnameInput").value==''){
      $('#signup_error').html('<b>Error:</b> No first name set!');
      $('#signup_error').show()
    }else if(document.getElementById("lastnameInput").value==''){
      $('#signup_error').html('<b>Error:</b> No last name set!');
      $('#signup_error').show()
    }else if(!isEmail(document.getElementById("emailInput").value)){
      $('#signup_error').html('<b>Error:</b> Invalid Email!');
      $('#signup_error').show()
    }else{
      new_signup = {
      "email":document.getElementById("emailInput").value,
      "pw_hash":document.getElementById("passwordInput").value,
      "name_first":document.getElementById("firstnameInput").value,
      "name_last": document.getElementById("lastnameInput").value 
      };
      
      post_text("/person", JSON.stringify(new_signup), function(person){
        console.log("success")
        buildAddUser()
        $("#signup_successful").html("Sign up successful!")
        $("#signup_successful").show()
      })
    }
  }

/**** Create Decision ****/
  function buildCreateDecision() {
    $('title').html('Create New Decision');
    clearContent();
    
    $('<strong><i class="glyphicon glyphicon-cog"></i> Create New Decision</strong><hr/>').appendTo('#content');
    
    var wrapper = $('<div>').css('max-width','500px').appendTo('#content');
    var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return createNewDecision()').appendTo(wrapper);
      $('<div>').attr('id','error').addClass('alert alert-danger').appendTo(form);
      $('#error').hide();
      
      $('<input type="text" />').addClass('form-control')
        .attr('name', 'name')
        .attr('placeholder', 'Decision Name')
        .attr('id', 'name')
        .attr('required', '')
        .appendTo(form);
      $('<br/>').appendTo(form);
      $('<textarea>').addClass('form-control')
        .attr('rows','3')
        .attr('name', 'description')
        .attr('placeholder', 'Decision Description')
        .attr('id', 'description')
        .attr('required', '') //Backend rejects code if this is null :(
        .appendTo(form);

      $('<hr/>').appendTo(form);
      $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('type', 'submit').text('Submit').appendTo(form);
  }

  function createNewDecision() {
    $('#error').hide();
    
    if(document.getElementById('name') == '') {
      $('#error').html('<b>Error:</b> No name set!');
      $('#error').show();
    }
    else if(document.getElementById('description') == '') {
      $('#error').html('<b>Error:</b> No description set!');
      $('#error').show();
    }
    else {
      get_text("/whoami", function (result) {
        var new_decision = {
          "person_id":+result['person_id'],
          "name":$("#name").val(),
          "description":$("#description").val(),
          "stage":+1,
          "criterion_vote_style":"s", //sliders
          "alternative_vote_style":"3", //3-color
          "client_settings":"",
          "display_name":$("#name").val(),
          "criteria_instruction":"",
          "alternative_instruction":""
        }
        post_text("/decision", JSON.stringify(new_decision), function(result){
          updateLeftNav();
          buildEditDecision(result['decision']['decision_id']);
        });
      });
    }
    return false;
  }

/**** Edit Decision ****/
  function buildDecisionHome(decisionID) {
    $('title').html('Edit Decision');
    clearContent();
    
    $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');
    
    var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
    $('<li class="active"><a onclick="buildDecisionHome('+decisionID+')">Decision</a></li>').appendTo(ul);
    $('<li><a onclick="buildCustomizeDecision('+decisionID+')">Customize</a></li>').appendTo(ul);
    $('<li><a onclick="buildEditCriteria('+decisionID+')">Criteria</a></li>').appendTo(ul);
    $('<li><a onclick="buildEditAlternative('+decisionID+')">Alternatives</a></li>').appendTo(ul);
    $('<li><a onclick="buildDecisionStatus('+decisionID+')">Status</a></li>').appendTo(ul);
    $('<li><a onclick="buildDecisionInvite('+decisionID+')">Invite</a></li>').appendTo(ul);
    
    buildEditDecision(decisionID);
    //hide the ones we don't want to show
    $('#displayNameDiv').hide();
    $('#critStyleDiv').hide();
    $('#altStyleDiv').hide();
    $('#critInstructionsDiv').hide();
    $('#altInstructionsDiv').hide();
  }
  
  function buildCustomizeDecision(decisionID) {
    $('title').html('Edit Decision');
    clearContent();
    
    $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');
    
    var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
    $('<li><a onclick="buildDecisionHome('+decisionID+')">Decision</a></li>').appendTo(ul);
    $('<li class="active"><a onclick="buildCustomizeDecision('+decisionID+')">Customize</a></li>').appendTo(ul);
    $('<li><a onclick="buildEditCriteria('+decisionID+')">Criteria</a></li>').appendTo(ul);
    $('<li><a onclick="buildEditAlternative('+decisionID+')">Alternatives</a></li>').appendTo(ul);
    $('<li><a onclick="buildDecisionStatus('+decisionID+')">Status</a></li>').appendTo(ul);
    $('<li><a onclick="buildDecisionInvite('+decisionID+')">Invite</a></li>').appendTo(ul);
    
    buildEditDecision(decisionID);
    //hide the ones we don't want to show
    $('#nameDiv').hide();
    $('#descriptionDiv').hide();
    $('#stageDiv').hide();
    $('#deleteDecisionBtn').hide();
  }

  function buildEditDecision(decisionID) {
    var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');
    var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return false;').appendTo(wrapper);
      $('<div>').attr('id','success').addClass('alert alert-success').appendTo(form);
      $('#success').hide();
      $('<div>').attr('id','error').addClass('alert alert-danger').appendTo(form);
      $('#error').hide();
      $('<div>').addClass('clearfix').appendTo(form);
      
      $('<div id="nameDiv" class="form-group">').append(
        $('<label for="name">Name</label>'),
        $('<input type="text" />').addClass('form-control')
        .attr('name', 'name')
        .attr('placeholder', 'Decision Name')
        .attr('id', 'name')
        .attr('required', ''))
      .appendTo(form);
      
      $('<div id="displayNameDiv" class="form-group">').append(
        $('<label for="displayName">Display Name</label>'),
        $('<input type="text" />').addClass('form-control')
        .attr('name', 'displayName')
        .attr('placeholder', 'Display Name')
        .attr('id', 'displayName'))
      .appendTo(form);
      
      $('<div id="descriptionDiv" class="form-group">').append(
        $('<label for="description">Description</label>'),
        $('<textarea>').addClass('form-control')
          .attr('rows','3')
          .attr('name', 'description')
          .attr('placeholder', 'Decision Description')
          .attr('id', 'description')
          .attr('required', '')) //Backend rejects code if this is null :(
      .appendTo(form);
      
      $('<div id="critInstructionsDiv" class="form-group">').append(
        $('<label for="critInstructions">Criteria Instructions</label>'),
        $('<textarea>').addClass('form-control')
          .attr('rows','3')
          .attr('name', 'critInstructions')
          .attr('placeholder', 'Criteria Instructions')
          .attr('id', 'critInstructions'))
      .appendTo(form);
      
      $('<div id="altInstructionsDiv" class="form-group">').append(
        $('<label for="altInstructions">Alternative Instructions</label>'),
        $('<textarea>').addClass('form-control')
          .attr('rows','3')
          .attr('name', 'altInstructions')
          .attr('placeholder', 'Alternative Instructions')
          .attr('id', 'altInstructions'))
      .appendTo(form);
      
      $('<div id="critStyleDiv" class="form-group">').append(
        '<label for="critStyle">Criteria Style</label>' +
          '<select id="critStyle" class="form-control">' +
            '<option value="s">Sliders</option>' +
            '<option value="b">Buttons</option>' +
          '</select>').appendTo(form);
          
      $('<div id="altStyleDiv" class="form-group">').append(
        '<label for="altStyle">Alternative Style</label>' +
          '<select id="altStyle" class="form-control">' +
            '<option value="3">3 Color</option>' +
            '<option value="5">5 Color</option>' +
          '</select>').appendTo(form);
      
      $('<div id="stageDiv" class="form-group">').append(
        '<label for="stage">Current Stage</label>' +
          '<select id="stage" class="form-control">' +
            '<option value="1">In Development</option>' +
            '<option value="2">Voting in Progress</option>' +
            '<option value="3">Completed</option>' +
          '</select>').appendTo(form);
      
      $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('onclick', 'updateDecision('+decisionID+');').text('Submit').appendTo(form);
      $('<button>').addClass('btn btn-lg btn-danger btn-block').attr('id', 'deleteDecisionBtn').attr('onclick', 'deleteDecision('+decisionID+');').text('Delete Decision').appendTo(form);
      
      get_text("/decision/"+decisionID+"/info", function (result) {
        $('#name').val(result['decision']['name']);
        $('#displayName').val(result['decision']['display_name']);
        $('#description').val(result['decision']['description']);
        $('#stage').val(result['decision']['stage']);
        $('#critStyle').val(result['decision']['criterion_vote_style']);
        $('#altStyle').val(result['decision']['alternative_vote_style']);
        $('#critInstructions').val(result['decision']['criteria_instruction']);
        $('#altInstructions').val(result['decision']['alternative_instruction']);
      });
  }

  function deleteDecision(decisionID) {
    confirmYesNo(
        "Delete Decision",
        "Are you sure you want to delete this decision and all associated ballots?",
        function() {
          delete_text("/decision/"+decisionID, function (result) {
            if(result['result'] == "deleted")
              alert("Decision Deleted!");
            //redirect
            updateLeftNav();
            buildHome();
          });
        },
        function() { /* Do nothing */}
    );
  }

  function updateDecision(decisionID) {
    $('#error').hide();
    $('#success').hide();
    
    if(document.getElementById('name') == '') {
      $('#error').html('<b>Error:</b> No name set!');
      $('#error').show();
    }
    else if(document.getElementById('description') == '') {
      $('#error').html('<b>Error:</b> No description set!');
      $('#error').show();
    }
    else {
      get_text("/whoami", function (result) {
        var new_decision = {
          "person_id":+result['person_id'],
          "name":$("#name").val(),
          "description":$("#description").val(),
          "stage":+$("#stage").val(),
          "criterion_vote_style":$("#critStyle").val(),
          "alternative_vote_style":$("#altStyle").val(),
          "client_settings":"",
          "display_name":$('#displayName').val(),
          "criteria_instruction":$('#critInstructions').val(),
          "alternative_instruction":$('#altInstructions').val()
        }

        put_text("/decision/" + decisionID, JSON.stringify(new_decision), function(result){
          updateLeftNav();
          $('#success').html('Updated Successfully');
          $('#success').show();
        });
      });
    }
  }

  /**** Decision Criteria ****/
    function buildEditCriteria(decisionID) {
      $('title').html('Edit Decision');
      clearContent();
      
      $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');
      
      var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
      $('<li><a onclick="buildDecisionHome('+decisionID+')">Decision</a></li>').appendTo(ul);
      $('<li><a onclick="buildCustomizeDecision('+decisionID+')">Customize</a></li>').appendTo(ul);
      $('<li class="active"><a onclick="buildEditCriteria('+decisionID+')">Criteria</a></li>').appendTo(ul);
      $('<li><a onclick="buildEditAlternative('+decisionID+')">Alternatives</a></li>').appendTo(ul);
      $('<li><a onclick="buildDecisionStatus('+decisionID+')">Status</a></li>').appendTo(ul);
      $('<li><a onclick="buildDecisionInvite('+decisionID+')">Invite</a></li>').appendTo(ul);
      
      var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');
        //Table of existing criteria here
        $('<div id="critList">').appendTo(wrapper);
        $('#critList').hide();
        $('<hr/>').appendTo(wrapper);
      var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return false;').appendTo(wrapper);
        
        $('<div>').attr('id','success').addClass('alert alert-success').appendTo(form);
        $('#success').hide();
        $('<div>').attr('id','error').addClass('alert alert-danger').appendTo(form);
        $('#error').hide();
        
        $('<div>').attr('id', 'critForm').appendTo(form);    
        showAddCriteria(decisionID);
        
        updateCritList(decisionID);
    }
    
    function showAddCriteria(decisionID) {
      $('#critForm').html("");
      ifDecisionInDevelopment(decisionID, function() {
        $('<h3>Add New Criterion</h3>').appendTo('#critForm');
        
          $('<div class="form-group">').append(
            $('<label for="critOrder">Order</label>'),
            $('<input type="text" />').addClass('form-control')
            .attr('name', 'critOrder')
            .attr('placeholder', 'Order')
            .attr('id', 'critOrder'))
          .appendTo('#critForm');
          
          $('<div class="form-group">').append(
            $('<label for="critName">Name</label>'),
            $('<input type="text" />').addClass('form-control')
            .attr('name', 'critName')
            .attr('placeholder', 'Criterion Name')
            .attr('id', 'critName')
            .attr('required', ''))
          .appendTo('#critForm');
          
          $('<div class="form-group">').append(
            $('<label for="critDesc">Description</label>'),
            $('<textarea>').addClass('form-control')
              .attr('rows','3')
              .attr('name', 'critDesc')
              .attr('placeholder', 'Criterion Description')
              .attr('id', 'critDesc'))
              .appendTo('#critForm');

          $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('onclick', 'addCriteria('+decisionID+');').text('Add Criteria').appendTo('#critForm');
      });
    }
    
    function showEditCriteria(decisionID, criterionID) {
      $('#critForm').html("");
      $('<h3>Edit Criterion</h3>').appendTo('#critForm');
        $('<div class="form-group">').append(
          $('<label for="critOrder">Order</label>'),
          $('<input type="text" />').addClass('form-control')
          .attr('name', 'critOrder')
          .attr('placeholder', 'Order')
          .attr('id', 'critOrder'))
        .appendTo('#critForm');
        
        $('<div class="form-group">').append(
          $('<label for="critName">Name</label>'),
          $('<input type="text" />').addClass('form-control')
          .attr('name', 'critName')
          .attr('placeholder', 'Criterion Name')
          .attr('id', 'critName')
          .attr('required', ''))
        .appendTo('#critForm');
        
        $('<div class="form-group">').append(
          $('<label for="critDesc">Description</label>'),
          $('<textarea>').addClass('form-control')
            .attr('rows','3')
            .attr('name', 'critDesc')
            .attr('placeholder', 'Criterion Description')
            .attr('id', 'critDesc'))
        .appendTo('#critForm');
        
        $('<button>').addClass('btn btn-primary').attr('onclick', 'showAddCriteria('+decisionID+');').text('Back to Add New Criteria').attr('style','float: left').appendTo('#critForm');
        $('<button>').addClass('btn btn-primary').attr('onclick', 'editCriteria('+decisionID+', ' + criterionID + ');').text('Update Criteria').attr('style','float: right').appendTo('#critForm');
        $('<div>').addClass('clearfix').appendTo('#critForm'); //added to fix display issue
        get_text("/decision/"+decisionID+"/criterion/"+criterionID+"/info", function (result) {
          $('#critOrder').val(result['criterion']['order']);
          $('#critName').val(result['criterion']['name']);
          $('#critDesc').val(result['criterion']['description']);
        });
    }

    function addCriteria(decisionID) {
      $('#success').hide();
      $('#error').hide();
      
      var new_crit = {
        "name":$("#critName").val(),
        "description":$("#critDesc").val(),
        "order":+$('#critOrder').val()
      }

      post_text("/decision/" + decisionID + '/criterion', JSON.stringify(new_crit), function(result){
        //$('#success').html('Updated Successfully');
        //$('#success').show();
        updateCritList(decisionID);
        showAddCriteria(decisionID); //clears it
      });
    }

    function editCriteria(decisionID, criterionID) {
      $('#success').hide();
      $('#error').hide();
      
      var crit = {
        "name":$("#critName").val(),
        "description":$("#critDesc").val(),
        "order":+$('#critOrder').val()
      }

      put_text("/decision/" + decisionID + '/criterion/' + criterionID, JSON.stringify(crit), function(result){
        //$('#success').html('Updated Successfully');
        //$('#success').show();
        updateCritList(decisionID);
        showAddCriteria(decisionID); //clears it
      });
    }
    
    function deleteCriteria(decisionID, criterionID) {
      $('#success').hide();
      $('#error').hide();
      
      confirmYesNo(
          "Delete Criteria",
          "Are you sure you want to delete this criterion?",
          function() {
            delete_text("/decision/"+decisionID+"/criterion/"+criterionID, function (result) {
              if(result['result'] == "deleted") {
                //$('#success').html('Deleted Successfully');
                //$('#success').show();
              }
              updateCritList(decisionID);
              showAddCriteria(decisionID);
            });
          },
          function() { /* Do nothing */}
      );
    }

    function updateCritList(decisionID) {
      //clear it to repopulate it
      $('#critList').html("");
      
      get_text("/decision/"+decisionID+"/criterions", function (results) {
          var table = $('<table>').append($('<tbody>')).addClass('table table-striped').appendTo('#critList');
          table.append('<tr><th></th><th>Name</th><th>Description</th><th></th></tr>');
          
          if(results["criterions"].length < 1) $('#critList').hide();
          else $('#critList').show();
          
          for(var i in results["criterions"]) {
            c = results["criterions"][i];
            table.append('<tr><td>'
              + c['order'] + '</td><td>'
              + c['name'] + '</td><td>'
              + c['description'] + '</td><td>'
              + '<div style="width:45px; float:right;"><a onclick="showEditCriteria('+ decisionID + ',' + c['criterion_id'] + ');"><span class="glyphicon glyphicon-pencil text-Primary"></span></a>'
              + '<a onclick="deleteCriteria('+ decisionID + ',' + c['criterion_id'] + ');"><span class="glyphicon glyphicon-trash text-Danger" style="margin-left:10px;"></span></a></div></td></tr>');
          }
          
          //Hide dangerous icons if not in development
          $('.glyphicon.text-Danger').hide();
          ifDecisionInDevelopment(decisionID, function() {$('.glyphicon.text-Danger').show();});
        });
    }

  /**** Decision Alternative ****/
    function buildEditAlternative(decisionID) {
      $('title').html('Edit Decision');
      clearContent();
      
      $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');
      
      var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
      $('<li><a onclick="buildDecisionHome('+decisionID+')">Decision</a></li>').appendTo(ul);
      $('<li><a onclick="buildCustomizeDecision('+decisionID+')">Customize</a></li>').appendTo(ul);
      $('<li><a onclick="buildEditCriteria('+decisionID+')">Criteria</a></li>').appendTo(ul);
      $('<li class="active"><a onclick="buildEditAlternative('+decisionID+')">Alternatives</a></li>').appendTo(ul);
      $('<li><a onclick="buildDecisionStatus('+decisionID+')">Status</a></li>').appendTo(ul);
      $('<li><a onclick="buildDecisionInvite('+decisionID+')">Invite</a></li>').appendTo(ul);
      
      var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');
        //Table of existing criteria here
        $('<div id="altList">').appendTo(wrapper);
        $('#altList').hide();
        $('<hr/>').appendTo(wrapper);
      var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return false;').appendTo(wrapper);
        
        $('<div>').attr('id','success').addClass('alert alert-success').appendTo(form);
        $('#success').hide();
        $('<div>').attr('id','error').addClass('alert alert-danger').appendTo(form);
        $('#error').hide();
        
        $('<div>').attr('id', 'altForm').appendTo(form);    
        showAddAlternative(decisionID);
        
        updateAltList(decisionID);
    }
    
    function showAddAlternative(decisionID) {
      $('#altForm').html("");
      ifDecisionInDevelopment(decisionID, function() {
        $('<h3>Add New Alternative</h3>').appendTo('#altForm');
        
          $('<div class="form-group">').append(
            $('<label for="altOrder">Order</label>'),
            $('<input type="text" />').addClass('form-control')
            .attr('name', 'altOrder')
            .attr('placeholder', 'Alternative Order')
            .attr('id', 'altOrder'))
          .appendTo('#altForm');
          
          $('<div class="form-group">').append(
            $('<label for="altName">Name</label>'),
            $('<input type="text" />').addClass('form-control')
            .attr('name', 'altName')
            .attr('placeholder', 'Alternative Name')
            .attr('id', 'altName')
            .attr('required', ''))
          .appendTo('#altForm');
          
          $('<div class="form-group">').append(
            $('<label for="altDesc">Description</label>'),
            $('<textarea>').addClass('form-control')
              .attr('rows','3')
              .attr('name', 'altDesc')
              .attr('placeholder', 'Alternative Description')
              .attr('id', 'altDesc'))
          .appendTo('#altForm');
          
          $('<div class="form-group">').append(
            $('<label for="altCost">Cost</label>'),
            $('<input type="text" />').addClass('form-control')
            .attr('name', 'altCost')
            .attr('placeholder', 'Alternative Cost')
            .attr('id', 'altCost'))
          .appendTo('#altForm');
          
          $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('onclick', 'addAlternative('+decisionID+');').text('Add Alternative').appendTo('#altForm');
      });
    }
    
    function showEditAlternative(decisionID, alternativeID) {
      $('#altForm').html("");
      $('<h3>Edit Alternative</h3>').appendTo('#altForm');
      
        $('<div class="form-group">').append(
          $('<label for="altOrder">Order</label>'),
          $('<input type="text" />').addClass('form-control')
          .attr('name', 'altOrder')
          .attr('placeholder', 'Alternative Order')
          .attr('id', 'altOrder'))
        .appendTo('#altForm');
        
        $('<div class="form-group">').append(
          $('<label for="altName">Name</label>'),
          $('<input type="text" />').addClass('form-control')
          .attr('name', 'altName')
          .attr('placeholder', 'Alternative Name')
          .attr('id', 'altName')
          .attr('required', ''))
        .appendTo('#altForm');
        
        $('<div class="form-group">').append(
          $('<label for="altDesc">Description</label>'),
          $('<textarea>').addClass('form-control')
            .attr('rows','3')
            .attr('name', 'altDesc')
            .attr('placeholder', 'Alternative Description')
            .attr('id', 'altDesc'))
        .appendTo('#altForm');
        
        $('<div class="form-group">').append(
          $('<label for="altCost">Cost</label>'),
          $('<input type="text" />').addClass('form-control')
          .attr('name', 'altCost')
          .attr('placeholder', 'Alternative Cost')
          .attr('id', 'altCost'))
        .appendTo('#altForm');
        
        $('<button>').addClass('btn btn-primary').attr('onclick', 'showAddAlternative('+decisionID+');').text('Back to Add New Alternative').attr('style','float: left').appendTo('#altForm');
        $('<button>').addClass('btn btn-primary').attr('onclick', 'editAlternative('+decisionID+', ' + alternativeID + ');').text('Update Alternative').attr('style','float: right').appendTo('#altForm');
        $('<div>').addClass('clearfix').appendTo('#altForm'); //added to fix display issue
        
        get_text("/decision/"+decisionID+"/alternative/"+alternativeID+"/info", function (result) {
          $('#altOrder').val(result['alternative']['order']);
          $('#altName').val(result['alternative']['name']);
          $('#altDesc').val(result['alternative']['description']);
          $('#altCost').val(result['alternative']['cost']);
        });
    }

    function addAlternative(decisionID) {
      $('#success').hide();
      $('#error').hide();
      
      var new_alt = {
        "name":$("#altName").val(),
        "description":$("#altDesc").val(),
        "cost":+$("#altCost").val(),
        "order":+$("#altOrder").val()
      }

      post_text("/decision/" + decisionID + '/alternative', JSON.stringify(new_alt), function(result){
        //$('#success').html('Updated Successfully');
        //$('#success').show();
        updateAltList(decisionID);
        showAddAlternative(decisionID); //clears it
      });
    }

    function editAlternative(decisionID, alternativeID) {
      $('#success').hide();
      $('#error').hide();
      
      var alt = {
        "name":$("#altName").val(),
        "description":$("#altDesc").val(),
        "cost":+$("#altCost").val(),
        "order":+$("#altOrder").val()
      }

      put_text("/decision/" + decisionID + '/alternative/' + alternativeID, JSON.stringify(alt), function(result){
        //$('#success').html('Updated Successfully');
        //$('#success').show();
        updateAltList(decisionID);
        showAddAlternative(decisionID); //clears it
      });
    }
    
    function deleteAlternative(decisionID, alternativeID) {
      $('#success').hide();
      $('#error').hide();
      
      confirmYesNo(
          "Delete Alternative",
          "Are you sure you want to delete this alternative?",
          function() {
            delete_text("/decision/"+decisionID+"/alternative/"+alternativeID, function (result) {
              if(result['result'] == "deleted") {
                //$('#success').html('Deleted Successfully');
                //$('#success').show();
              }
              updateAltList(decisionID);
              showAddAlternative(decisionID);
            });
          },
          function() { /* Do nothing */}
      );
    }

    function updateAltList(decisionID) {
      //clear it to repopulate it
      $('#altList').html("");
      
      get_text("/decision/"+decisionID+"/alternatives", function (results) {
          var table = $('<table>').append($('<tbody>')).addClass('table table-striped').appendTo('#altList');
          table.append('<tr><th></th><th>Name</th><th>Description</th><th>Cost</th><th></th></tr>');
          
          if(results["alternatives"].length < 1) $('#altList').hide();
          else $('#altList').show();
          
          for(var i in results["alternatives"]) {
            a = results["alternatives"][i];
            table.append('<tr><td>'
              + a['order'] + '</td><td>'
              + a['name'] + '</td><td>'
              + a['description'] + '</td><td>'
              + a['cost'] + '</td><td>'
              + '<div style="width:45px; float:right;"><a onclick="showEditAlternative('+ decisionID + ',' + a['alternative_id'] + ');"><span class="glyphicon glyphicon-pencil text-Primary"></span></a>'
              + '<a onclick="deleteAlternative('+ decisionID + ',' + a['alternative_id'] + ');"><span class="glyphicon glyphicon-trash text-Danger" style="margin-left:10px;"></span></a></div></td></tr>');
          }
          
          //Hide dangerous icons if not in development
          $('.glyphicon.text-Danger').hide();
          ifDecisionInDevelopment(decisionID, function() {$('.glyphicon.text-Danger').show();});
        });
    }

  /**** Decision Status ****/
    function buildDecisionStatus(decisionID) {
      $('title').html('Edit Decision');
      clearContent();
      
      $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');
      
      var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
      $('<li><a onclick="buildDecisionHome('+decisionID+')">Decision</a></li>').appendTo(ul);
      $('<li><a onclick="buildCustomizeDecision('+decisionID+')">Customize</a></li>').appendTo(ul);
      $('<li><a onclick="buildEditCriteria('+decisionID+')">Criteria</a></li>').appendTo(ul);
      $('<li><a onclick="buildEditAlternative('+decisionID+')">Alternatives</a></li>').appendTo(ul);
      $('<li class="active"><a onclick="buildDecisionStatus('+decisionID+')">Status</a></li>').appendTo(ul);
      $('<li><a onclick="buildDecisionInvite('+decisionID+')">Invite</a></li>').appendTo(ul);
      
      var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');
    var ballotsForCurrentStage = $('<div>').attr('id','totalBallots').appendTo(wrapper) 
    
    var status_table = $([
    
    '<table class="table table-bordered" id="s_table">',
      '<thead class = "thead-inverse">',
      '<tr>',
      '<th>Name</th>',
      '<th>Email</th>',
      '<th>Ballot status</th>',
      '<th>Email</th>',
      '</tr>',
      '</thead>',
    '</table>'
    ].join('\n'));
  
    wrapper.append(status_table)
    $('#s_table').hide()

    getStatus(decisionID)
    }
  
    function getStatus(decisionID){
      get_text("/decision/"+ decisionID + "/ballots", function (result) {
        console.log(result)
        var vote_status= ""
        var totalBallots = 0
        if(result.ballots != null){
        $('#s_table').show()
        //get the total ballot for current stage
        totalBallots = result.ballots.length
        getTotalBallots(decisionID,totalBallots)
        
        for(var i = 0; i < result.ballots.length; i++){
          if(result.ballots[i].rating != null && result.ballots[i].rating.length > 0 && result.ballots[i].rating.length!= "undefined" ){
            vote_status = "Voted"
          }else{
            vote_status = "Note Voted"
          }
          console.log(result.ballots[i].url)
          var url = result.ballots[i].url
          
          $("#s_table").append('<tbody><tr><td>' +result.ballots[i].name + '</td><td>' + result.ballots[i].email +'</td><td>' + vote_status + '</td><td> <a onclick=resendEmail(\''+url+'\')>Resend email</a>'+'</td> </tbody>');       
    
          
        }
        }else{
          getTotalBallots(decisionID,totalBallots)
        }
      })
    }
    
    function getTotalBallots(decisionID, ballots){
      get_text("/decision/"+ decisionID + "/info", function (result) {
                          var stageBallots = $([
                                  '<ul class="list-group" id="stage_ballot">',
                                  '<li class="list-group-item active">',
                                  'Total ballots: '+ ballots ,
                                  '</li>',
                                  '</ul>'].join("\n"));
                          stageBallots.appendTo('#totalBallots');
                          if(result.decision.stage == 1) {
                                  $("#stage_ballot").append('<li class="list-group-item"><span class="badge">'+ballots+'</span>Stage: In Development </li>');
                          } else if(result.decision.stage == 2) {
                                  $("#stage_ballot").append('<li class="list-group-item"><span class="badge">'+ballots+'</span>Stage: Voting in progress </li>');
                          } else {
                                  $("#stage_ballot").append('<li class="list-group-item"><span class="badge">'+ballots+'</span>Stage: Completed </li>');
                          }
          })
      
    }
    
    function resendEmail(url){
      get_text(url + "/invite", function(result){
        console.log(result)
      })
    }

  /**** Decision Invite ****/
    function buildDecisionInvite(decisionID){

      $('title').html('Invite People');
      clearContent();
      
      $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');
      
      var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
      $('<li><a onclick="buildDecisionHome('+decisionID+')">Decision</a></li>').appendTo(ul);
      $('<li><a onclick="buildCustomizeDecision('+decisionID+')">Customize</a></li>').appendTo(ul);
      $('<li><a onclick="buildEditCriteria('+decisionID+')">Criteria</a></li>').appendTo(ul);
      $('<li><a onclick="buildEditAlternative('+decisionID+')">Alternatives</a></li>').appendTo(ul);
      $('<li><a onclick="buildDecisionStatus('+decisionID+')">Status</a></li>').appendTo(ul);
      $('<li class="active"><a onclick="buildDecisionInvite('+decisionID+')">Invite</a></li>').appendTo(ul);
      
      var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');

      var form = $([
        '<form class ="form-signin" onsubmit = "return false" id="inviteForm">',
        
        '<div id= "invitation_sent" class = "alert alert-success"></div>',
        '<div id="invitation_error" class="alert alert-danger" style="display: none;"></div>',
        '<label for="bal_dec_id" >Decision Name: </label>',
        '<input type="text" class= "form-control" required="required" placeholder="Decision ID" id="dName"></input>',
        '<br />',

        '<label for="bal_name">Name</label>',
        '<input type="text" id="i_name" class= "form-control" required="required" placeholder="Name"></input>',
        '<br />',
        '<label for="bal_email">Email</label>',
        '<input type="email" id="i_email" class= "form-control" required="required" placeholder="Email"></input>',
        '<br />',
        '</form>'
        ].join('\n'));
        
        $('<button>').addClass('btn btn-primary').attr('id','invite_submit').attr('onclick', 'buildEditDecision('+decisionID+');').append('<span> <i class="glyphicon glyphicon-arrow-left"></i>  Back to Decision </span>').appendTo(form);
        $('<button>').addClass('btn btn-primary').attr('id','invite_submit').attr('onclick', 'sendInvite('+decisionID+');').attr('style','float: right').append('<span> <i class="glyphicon glyphicon-envelope"></i>  Invite </span>').appendTo(form);
        
      
        get_text("/decision/"+decisionID+"/info", function (result) {
          $('#dName').val(result.decision.name)
        })
        
        $(wrapper).append(form)

        $("#invitation_sent").hide()
        $("#invitation_error").hide()
    }

    function sendInvite(decisionID){
      $("#invitation_sent").hide()
      $("#invitation_error").hide()
      id = decisionID
      
      if(document.getElementById("i_name").value =='') {
        $('#invitation_error').html('<b>Error:</b> No name set!');
        $('#invitation_error').show();
      }
      else if(document.getElementById("i_email").value =='') {
        $('#invitation_error').html('<b>Error:</b> No email set!');
        $('#invitation_error').show();
      }else if(!isEmail(document.getElementById("i_email").value)){
        $('#invitation_error').html('<b>Error:</b> Invalid email!');
        $('#invitation_error').show()
      }else{
        
        new_invite = {
        "name":document.getElementById("i_name").value,
        "email":document.getElementById("i_email").value 
        };
        post_text("/decision/"+id+"/ballot", JSON.stringify(new_invite), function(result) {
          console.log(result);  
          buildDecisionInvite(id)
          $('#invitation_sent').html('Invitation sent Successfully!');
          $('#invitation_sent').show()
        })
      }
    }
    
    //Runs cb if true
    function ifDecisionInDevelopment(decisionID, cb) {
      get_text("/decision/"+decisionID+"/info", function (result) {
        if(+result['decision']['stage']==1) cb();
      });
    }