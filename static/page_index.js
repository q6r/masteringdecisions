var intervalID = 0; //Used to refresh pages on the fly

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
  get_text("/whoami", function(person) {
    get_text("/person/" + person['person_id'] + "/decisions", function(decisions) {
      var inprogress = [];
      var completed = [];
      var locked = [];
      for (var i in decisions["decisions"]) {
        d = decisions["decisions"][i];
        if (d["stage"] < 3) {
          inprogress.push(d);
        } else if (d["stage"] == 3) {
          completed.push(d);
        } else {
          locked.push(d);
        }
      }
      cb(inprogress, completed, locked);
    });
  });
}

//Updates the list of inprogress and completed decisions in the leftNav
function updateLeftNav() {
  decisionListByCategory(function(inprogress, completed, locked) {

    $("#userMenu3").html("");
    $("#userMenu4").html("");
    $("#userMenu5").html("");

    for (var i in inprogress) {
      dname = inprogress[i]["name"];
      did = inprogress[i]["decision_id"];
      $("#userMenu3").append("<li><a onclick=\"buildDecisionHome(" + did + ")\"><i class=\"glyphicon glyphicon-list-alt\"></i> " + dname + " </a></li>");
    }

    for (var i in completed) {
      dname = completed[i]["name"];
      did = completed[i]["decision_id"];
      $("#userMenu4").append("<li><a onclick=\"buildDecisionHome(" + did + ")\"><i class=\"glyphicon glyphicon-list-alt\"></i> " + dname + " </a></li>");
    }

    for (var i in locked) {
      dname = locked[i]["name"];
      did = locked[i]["decision_id"];
      $("#userMenu5").append("<li><a onclick=\"buildDecisionHome(" + did + ")\"><i class=\"glyphicon glyphicon-list-alt\"></i> " + dname + " </a></li>");
    }

  });
}

//builds the main template for the page
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

  //topnav section
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
    '<li><a onclick="buildEditProfile()">Edit Profile</a></li>',
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
    '<li class="nav-header"> <a  data-toggle="collapse" data-target="#userMenu5" aria-expanded="false" class="collapsed">Locked <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
    '<ul class="nav nav-stacked collapse" id="userMenu5" aria-expanded="false" style="height: 0px;">',
    '</ul>',
    '</li>',
    '</ul>',
    '</li>',
    '<li class="nav-header"> <a  data-toggle="collapse" data-target="#Menu2" aria-expanded="false" class="collapsed">Settings <i id="arrow_change" class="glyphicon glyphicon-chevron-right"></i></a>',
    '<ul class="nav nav-stacked collapse" id="Menu2" aria-expanded="fasle">',
    '<li class="active"> <a onclick="buildHome()"><i class="glyphicon glyphicon-home"></i> Home</a></li>',
    '<li id="manageUsers"><a onclick="buildManageUsers()"><i class="glyphicon glyphicon-user"></i> Manage Users</a></li>',
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
  $('#manageUsers').hide();

  //Toggles the arrows to be correct
  $("a").attr("aria-expanded", "true");
  $("a").click(function() {
    $(this).find('i#arrow_change').toggleClass('glyphicon-chevron-right glyphicon-chevron-down');
  });
}

//clears the content section of the page
function clearContent() {
  $('#content').empty();
  if (intervalID != 0) {
    window.clearInterval(intervalID);
  }
}

//Updates the username in the topNav
function updateUserText() {
  get_text("/whoami", function(result) {
    if (result['person_id'] == 0) {
      $('#manageUsers').show();
    }
    get_text("/person/" + result['person_id'] + "/info", function(result) {
      $('#userName').text(' ' + result['person']['name_first'] + ' ' + result['person']['name_last']);
    });
  });
}

function buildHome() {
  $('title').html('Decision Home');

  clearContent();

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
  decisionListByCategory(function(inprogress, completed, locked) {

    $("#inprogress_list").html("<a class=\"list-group-item active\">Decisions In Progress</a>");
    $("#completed_list").html("<a class=\"list-group-item active\">Decisions Completed</a>");

    for (var i in inprogress) {
      dname = inprogress[i]["name"];
      did = inprogress[i]["decision_id"];
      $("#inprogress_list").append("<a  onclick=\"buildDecisionHome(" + did + ")\" class=\"list-group-item\">" + dname + "</a>");
    }

    for (var i in completed) {
      dname = completed[i]["name"];
      did = completed[i]["decision_id"];
      $("#completed_list").append("<a  onclick=\"buildDecisionHome(" + did + ")\" class=\"list-group-item\">" + dname + "</a>");
    }

    // Set progressbar lengths
    totalLength = inprogress.length + completed.length;
    inprogress_progress = ((inprogress.length / totalLength) * 100);
    completed_progress = ((completed.length / totalLength) * 100);
    $("#inprogress_length").width(inprogress_progress + "%");
    $("#completed_length").width(completed_progress + "%");

  });
}

/**** Edit Profile ****/
function buildEditProfile() {
  $('title').html('Update User!');

  clearContent();

  $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Person</strong><hr/>').appendTo('#content');

  var wrapper = $('<div>').css('max-width', '500px').appendTo('#content');
  var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return updateProfile()').appendTo(wrapper);
  $('<h3>').addClass('form-signin-heading').text('Update Details').appendTo(form);

  $('<div>').attr('id', 'success').addClass('alert alert-success').appendTo(form);
  $('#success').hide();
  $('<div>').attr('id', 'error').addClass('alert alert-danger').appendTo(form);
  $('#error').hide();

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

  $('<div>').attr('id', 'passwordError').addClass('alert alert-danger').appendTo(form);
  $('#passwordError').hide();

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

  get_text("/whoami", function(result) {
    get_text("/person/" + result['person_id'] + "/info", function(result) {
      $('#username').val(result['person']['email']);
      $('#firstname').val(result['person']['name_first']);
      $('#lastname').val(result['person']['name_last']);
    });
  });
}

function updateProfile() {
  $('#success').hide();
  $('#error').hide();
  $('#passwordError').hide();

  if ($('#password').val() != $('#password2').val()) {
    $('#passwordError').html('<b>Error:</b> Passwords do not match!');
    $('#passwordError').show();
  } else {
    if ($('#password').val() == "") {
      new_info = {
        "email": $('#username').val(),
        "name_first": $('#firstname').val(),
        "name_last": $('#lastname').val()
      }
    } else {
      new_info = {
        "email": $('#username').val(),
        "pw_hash": $('#password').val(),
        "name_first": $('#firstname').val(),
        "name_last": $('#lastname').val()
      }
    }

    get_text("/whoami", function(result) {
      put_text("/person/" + result['person_id'], JSON.stringify(new_info), function(result) {
        if (result['error']) {
          $('#error').html('<b>Error:</b> ' + result['error']);
          $('#error').show()
        } else if (result['person']) {
          $('#password').val('');
          $('#password2').val('');
          $('#success').html('<b>Update successful!</b>');
          $('#success').show();
          updateUserText();
        }
      });
    });
  }

  return false;
}

/**** Manage Users ****/
function buildManageUsers() {
  $('title').html('Manage Users');
  clearContent();

  $('<strong><i class="glyphicon glyphicon-user"></i> Manage Users</strong><hr/>').appendTo('#content');

  $('<div id="userList">').appendTo('#content');
  $('<hr/>').appendTo('#content');

  var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return false;').appendTo('#content');
  $('<div>').attr('id', 'success').addClass('alert alert-success').appendTo(form);
  $('#success').hide();
  $('<div>').attr('id', 'error').addClass('alert alert-danger').appendTo(form);
  $('#error').hide();
  $('<div>').attr('id', 'userForm').appendTo(form);

  showAddUser();
  updateUserList();
}

function showAddUser() {
  $('#userForm').html("");
  $('<h3>Add New User</h3>').appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userEmail">Email</label>'),
      $('<input type="text" />').addClass('form-control')
      .attr('name', 'userEmail')
      .attr('placeholder', 'Email')
      .attr('id', 'userEmail')
      .attr('required', ''))
    .appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userPassword">Password</label>'),
      $('<input type="password" />').addClass('form-control')
      .attr('name', 'userPassword')
      .attr('placeholder', 'Password')
      .attr('id', 'userPassword')
      .attr('required', ''))
    .appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userFirst">First Name</label>'),
      $('<input type="text" />').addClass('form-control')
      .attr('name', 'userFirst')
      .attr('placeholder', 'First Name')
      .attr('id', 'userFirst')
      .attr('required', ''))
    .appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userLast">Last Name</label>'),
      $('<input type="text" />').addClass('form-control')
      .attr('name', 'userLast')
      .attr('placeholder', 'Last Name')
      .attr('id', 'userLast')
      .attr('required', ''))
    .appendTo('#userForm');

  $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('onclick', 'addUser();').text('Add User').appendTo('#userForm');
}

function showEditUser(userID) {
  $('#userForm').html("");
  $('<h3>Add New User</h3>').appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userEmail">Email</label>'),
      $('<input type="text" />').addClass('form-control')
      .attr('name', 'userEmail')
      .attr('placeholder', 'Email')
      .attr('id', 'userEmail')
      .attr('required', ''))
    .appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userPassword">Password</label>'),
      $('<input type="password" />').addClass('form-control')
      .attr('name', 'userPassword')
      .attr('placeholder', 'Password')
      .attr('id', 'userPassword'))
    .appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userPassword2">Confirm Password</label>'),
      $('<input type="password" />').addClass('form-control')
      .attr('name', 'userPassword2')
      .attr('placeholder', 'Password')
      .attr('id', 'userPassword2'))
    .appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userFirst">First Name</label>'),
      $('<input type="text" />').addClass('form-control')
      .attr('name', 'userFirst')
      .attr('placeholder', 'First Name')
      .attr('id', 'userFirst')
      .attr('required', ''))
    .appendTo('#userForm');

  $('<div class="form-group">').append(
      $('<label for="userLast">Last Name</label>'),
      $('<input type="text" />').addClass('form-control')
      .attr('name', 'userLast')
      .attr('placeholder', 'Last Name')
      .attr('id', 'userLast')
      .attr('required', ''))
    .appendTo('#userForm');

  $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('onclick', 'editUser(' + userID + ');').text('Submit Changes').appendTo('#userForm');

  get_text("/person/" + userID + "/info", function(result) {
    $('#userEmail').val(result['person']['email']);
    $('#userFirst').val(result['person']['name_first']);
    $('#userLast').val(result['person']['name_last']);
  });
}

function addUser() {
  $("#success").hide()
  $("#error").hide()

  if ($('#userEmail').val() == '') {
    $('#error').html('<b>Error:</b> No email set!');
    $('#error').show();
  } else if ($('#userPassword').val() == '') {
    $('#error').html('<b>Error:</b> No password set!');
    $('#error').show();
  } else if ($('#userFirst').val() == '') {
    $('#error').html('<b>Error:</b> No first name set!');
    $('#error').show()
  } else if ($('#userLast').val() == '') {
    $('#error').html('<b>Error:</b> No last name set!');
    $('#error').show()
  } else if (!isEmail($('#userEmail').val())) {
    $('#error').html('<b>Error:</b> Invalid Email!');
    $('#error').show()
  } else {
    new_signup = {
      "email": $('#userEmail').val(),
      "pw_hash": $('#userPassword').val(),
      "name_first": $('#userFirst').val(),
      "name_last": $('#userLast').val()
    };

    post_text("/person", JSON.stringify(new_signup), function(person) {
      if (person['error']) {
        $('#error').html('<b>Error:</b> ' + person['error']);
        $('#error').show()
      } else if (person['person']) {
        buildManageUsers();
        $("#successful").html("Sign up successful!")
        $("#successful").show()
      } else {
        $('#error').html('<b>Error:</b> Something went wrong :(');
        $('#error').show()
      }
    });
  }
}

function editUser(userID) {
  $("#success").hide()
  $("#error").hide()

  if ($('#userEmail').val() == '') {
    $('#error').html('<b>Error:</b> No email set!');
    $('#error').show();
  } else if ($('#userFirst').val() == '') {
    $('#error').html('<b>Error:</b> No first name set!');
    $('#error').show();
  } else if ($('#userLast').val() == '') {
    $('#error').html('<b>Error:</b> No last name set!');
    $('#error').show();
  } else if (!isEmail($('#userEmail').val())) {
    $('#error').html('<b>Error:</b> Invalid Email!');
    $('#error').show();
  } else if ($('#userPassword').val() != $('#userPassword2').val()) {
    $('#error').html('<b>Error:</b> Passwords do not match!');
    $('#error').show();

    //Clear them
    $('userPassword').val('');
    $('userPassword2').val('');
  } else {
    //Don't update the Password unless they typed one in
    if ($('userPassword') == '') {
      user = {
        "email": $('#userEmail').val(),
        "name_first": $('#userFirst').val(),
        "name_last": $('#userLast').val()
      };
    } else {
      user = {
        "email": $('#userEmail').val(),
        "pw_hash": $('#userPassword').val(),
        "name_first": $('#userFirst').val(),
        "name_last": $('#userLast').val()
      };
    }

    put_text("/person/" + userID, JSON.stringify(user), function(person) {
      if (person['error']) {
        $('#error').html('<b>Error:</b> ' + person['error']);
        $('#error').show()
      } else if (person['person']) {
        buildManageUsers();
        $("#successful").html("Updated successful!")
        $("#successful").show()
      } else {
        $('#error').html('<b>Error:</b> Something went wrong :(');
        $('#error').show()
      }
    });
  }
}

function deleteUser(userID) {
  $('#success').hide();
  $('#error').hide();

  confirmYesNo(
    "Delete User",
    "Are you sure you want to delete this user?",
    function() {
      delete_text("/person/" + userID, function(result) {
        buildManageUsers();
        if (result['error']) {
          $('#signup_error').html('<b>Error:</b> ' + result['error']);
          $('#signup_error').show()
        } else if (result['result'] == "deleted") {
          $('#success').html('Deleted Successfully');
          $('#success').show();
        } else {
          $('#error').html('Soemthing went wrong :(');
          $('#error').show();
        }
      });
    },
    function() { /* Do nothing */ }
  );
}

function updateUserList() {
  //clear it to repopulate it
  $('#userList').html("");

  get_text("/persons", function(results) {
    var table = $('<table>').append($('<tbody>')).addClass('table table-striped').appendTo('#userList');
    table.append('<tr><th>Email</th><th>First Name</th><th>Last Name</th><th></th></tr>');

    for (var i in results["persons"]) {
      p = results["persons"][i];
      if (p['person_id'] == 0) { //is admin so hide delete button
        table.append('<tr><td>' + p['email'] + '</td><td>' + p['name_first'] + '</td><td>' + p['name_last'] + '</td><td>' + '<div style="width:45px; float:right;"><a onclick="showEditUser(' + p['person_id'] + ');"><span class="glyphicon glyphicon-pencil text-Primary"></span></a></div></td></tr>');
      } else {
        table.append('<tr><td>' + p['email'] + '</td><td>' + p['name_first'] + '</td><td>' + p['name_last'] + '</td><td>' + '<div style="width:45px; float:right;"><a onclick="showEditUser(' + p['person_id'] + ');"><span class="glyphicon glyphicon-pencil text-Primary"></span></a>' + '<a onclick="deleteUser(' + p['person_id'] + ');"><span class="glyphicon glyphicon-trash text-Danger" style="margin-left:10px;"></span></a></div></td></tr>');
      }
    }
  });
}

/**** Create Decision ****/
function buildCreateDecision() {
  $('title').html('Create New Decision');
  clearContent();

  $('<strong><i class="glyphicon glyphicon-cog"></i> Create New Decision</strong><hr/>').appendTo('#content');

  var wrapper = $('<div>').css('max-width', '500px').appendTo('#content');
  var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return createNewDecision()').appendTo(wrapper);
  $('<div>').attr('id', 'error').addClass('alert alert-danger').appendTo(form);
  $('#error').hide();

  $('<input type="text" />').addClass('form-control')
    .attr('name', 'name')
    .attr('placeholder', 'Decision Name')
    .attr('id', 'name')
    .attr('required', '')
    .appendTo(form);
  $('<br/>').appendTo(form);
  $('<textarea>').addClass('form-control')
    .attr('rows', '3')
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

  if ($('#name').val() == '') {
    $('#error').html('<b>Error:</b> No name set!');
    $('#error').show();
  } else if ($('#description').val() == '') {
    $('#error').html('<b>Error:</b> No description set!');
    $('#error').show();
  } else {
    get_text("/whoami", function(result) {
      if (result['person_id'] == 0) {
        $('#error').html('<b>Error:</b> The admin account cannot create decisons!');
        $('#error').show();
        return;
      }
      var new_decision = {
        "person_id": +result['person_id'],
        "name": $("#name").val(),
        "description": $("#description").val(),
        "stage": +1,
        "criterion_vote_style": "s", //sliders
        "alternative_vote_style": "3", //3-color
        "client_settings": "",
        "display_name": $("#name").val(),
        "criteria_instruction": "<h4><strong>Instructions:</strong></h4>\r\n<p style=\"margin-bottom: 50px;\">Using the sliders below, rate the importance of each evaluation factor.</p>",
        "alternative_instruction": "<h4><strong>Instructions:</strong></h4>\r\n<p>Using the boxes below, rate how well each alternative satisfies each evaluation factor.</p>\r\n<p>It is best to work on one factor at a time.</p>",
        "image": ""
      }
      post_text("/decision", JSON.stringify(new_decision), function(result) {
        if (result['error']) {
          $('#error').html('<b>Error:</b> ' + result['error']);
          $('#error').show();
        } else if (result['decision']) {
          updateLeftNav();
          buildDecisionHome(result['decision']['decision_id']);
        } else {
          $('#error').html('<b>Error:</b> Something went wrong :(');
          $('#error').show()
        }
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
  $('<li class="active"><a onclick="buildDecisionHome(' + decisionID + ')">Decision</a></li>').appendTo(ul);
  $('<li><a onclick="buildCustomizeDecision(' + decisionID + ')">Customize</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditCriteria(' + decisionID + ')">Criteria</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditAlternative(' + decisionID + ')">Alternatives</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionInvite(' + decisionID + ')">Invite</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionStatus(' + decisionID + ')">Voting Status</a></li>').appendTo(ul);

  buildEditDecision(decisionID);
  //hide the ones we don't want to show
  $('#displayNameDiv').hide();
  $('#imageDiv').hide();
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
  $('<li><a onclick="buildDecisionHome(' + decisionID + ')">Decision</a></li>').appendTo(ul);
  $('<li class="active"><a onclick="buildCustomizeDecision(' + decisionID + ')">Customize</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditCriteria(' + decisionID + ')">Criteria</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditAlternative(' + decisionID + ')">Alternatives</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionInvite(' + decisionID + ')">Invite</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionStatus(' + decisionID + ')">Voting Status</a></li>').appendTo(ul);

  buildEditDecision(decisionID);
  //hide the ones we don't want to show
  $('#nameDiv').hide();
  $('#descriptionDiv').hide();
  $('#stageDiv').hide();
  $('#deleteDecisionBtn').hide();
  $('#duplicateDecisionBtn').hide();
}

function buildEditDecision(decisionID) {
  var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');
  var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return false;').appendTo(wrapper);
  $('<div>').attr('id', 'success').addClass('alert alert-success').appendTo(form);
  $('#success').hide();
  $('<div>').attr('id', 'error').addClass('alert alert-danger').appendTo(form);
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

  $('<div id="imageDiv" class="form-group">')
    .css('padding-bottom', '20px').append(
      $('<label for="image">Custom Image</label>'),
      $('<br/><img id="decImg" src="#"/><br/>'),
      $('<input>')
      .css('float', 'left')
      .attr('type', 'file')
      .attr('id', 'image')
      .attr('accept', 'image/*')
      .change(function() {
        loadImg(this);
      }),
      $('<a>')
      .text('Remove Image')
      .addClass('text-Danger')
      .css('float', 'right')
      .click(function() {
        $('#decImg').attr('src', '');
      })
    ).appendTo(form);
  $('<div>').addClass('clearFix').appendTo(form);
  $('<div id="descriptionDiv" class="form-group">').append(
      $('<label for="description">Description</label>'),
      $('<textarea>').addClass('form-control')
      .attr('rows', '3')
      .attr('name', 'description')
      .attr('placeholder', 'Decision Description')
      .attr('id', 'description')
      .attr('required', '')) //Backend rejects code if this is null :(
    .appendTo(form);

  $('<div id="critInstructionsDiv" class="form-group">').append(
      $('<label for="critInstructions">Criteria Instructions</label>'),
      $('<textarea>').addClass('form-control')
      .attr('rows', '3')
      .attr('name', 'critInstructions')
      .attr('placeholder', 'Criteria Instructions')
      .attr('id', 'critInstructions'))
    .appendTo(form);

  $('<div id="altInstructionsDiv" class="form-group">').append(
      $('<label for="altInstructions">Alternative Instructions</label>'),
      $('<textarea>').addClass('form-control')
      .attr('rows', '3')
      .attr('name', 'altInstructions')
      .attr('placeholder', 'Alternative Instructions')
      .attr('id', 'altInstructions'))
    .appendTo(form);

  $('<div id="critStyleDiv" class="form-group">').append(
    '<label for="critStyle">Criteria Style</label>' +
    '<select id="critStyle" class="form-control">' +
    '<option value="s">Number Sliders</option>' +
    '<option value="t">Word Sliders</option>' +
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
    '<option value="4">Locked</option>' +
    '</select>').appendTo(form);

  $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('onclick', 'updateDecision(' + decisionID + ');').text('Save').appendTo(form);
  $('<button>').addClass('btn btn-lg btn-danger btn-block').attr('id', 'deleteDecisionBtn').attr('onclick', 'deleteDecision(' + decisionID + ');').text('Delete Decision').appendTo(form);
  $('<button>').addClass('btn btn-lg btn-success btn-block').attr('id', 'duplicateDecisionBtn').attr('onclick', 'duplicateDecision(' + decisionID + ');').text('Duplicate Decision').appendTo(form);
  $('<hr/>').appendTo(form);
  $('<button>').addClass('btn btn-lg btn-info').attr('id', 'decisionVotesBtn').attr('onclick', '$("<a>").attr("href","/decision/' + decisionID + '").attr("target", "_blank")[0].click();').text('View Votes').css('float', 'left').css('width', '48%').appendTo(form);
  $('<button>').addClass('btn btn-lg btn-info').attr('id', 'decisionResultsBtn').attr('onclick', '$("<a>").attr("href","/results/' + decisionID + '").attr("target", "_blank")[0].click();').text('View Results').css('float', 'right').css('width', '48%').appendTo(form);
  $('#decisionResultsBtn').hide();
  $('#decisionVotesBtn').hide();

  get_text("/decision/" + decisionID + "/info", function(result) {
    $('#name').val(result['decision']['name']);
    $('#displayName').val(result['decision']['display_name']);
    $('#decImg').attr('src', result['decision']['image']);
    $('#description').val(result['decision']['description']);
    $('#stage').val(result['decision']['stage']);
    $('#critStyle').val(result['decision']['criterion_vote_style']);
    $('#altStyle').val(result['decision']['alternative_vote_style']);
    $('#critInstructions').val(result['decision']['criteria_instruction']);
    $('#altInstructions').val(result['decision']['alternative_instruction']);

    if (result['decision']['stage'] == 3) { //if decision is completed show results Btn
      $('#decisionResultsBtn').show();
      $('#decisionVotesBtn').show();
    }
  });
}

function loadImg(input) {
  if (input.files && input.files[0]) {
    var reader = new FileReader();
    reader.onload = function(e) {
      $('#decImg').attr('src', e.target.result);
    }
    reader.readAsDataURL(input.files[0]);
  }
}

function deleteDecision(decisionID) {
  confirmYesNo(
    "Delete Decision",
    "Are you sure you want to delete this decision and all associated ballots?",
    function() {
      delete_text("/decision/" + decisionID, function(result) {
        if (result['result'] != "deleted") {
          alert("Decision Deleted!");
          $('#error').html('<b>Error:</b> Something went wrong!');
          $('#error').show();
        } else {
          updateLeftNav();
          buildHome();
        }
      });
    },
    function() { /* Do nothing */ }
  );
}

function duplicateDecision(decisionID) {
  confirmYesNo(
    "Duplicate Decision",
    "Are you sure you want to duplicate this decision?",
    function() {
      get_text("/decision/" + decisionID + "/duplicate", function(result) {
        if (result['error']) {
          $('#error').html('<b>Error:</b> ' + result['error']);
          $('#error').show();
        } else if (result['decision']) {
          updateLeftNav();
          buildDecisionHome(result['decision']['decision_id']);
          $('#success').html('Successfully duplicated decision (You are in the new decision)');
          $('#success').show();
        } else {
          $('#error').html('<b>Error:</b> Something went wrong :(');
          $('#error').show();
        }
      });
    },
    function() { /* Do nothing */ }
  );
}

function updateDecision(decisionID) {
  $('#error').hide();
  $('#success').hide();
  $('#decisionResultsBtn').hide();
  $('#decisionVotesBtn').hide();

  if ($('#name').val() == '') {
    $('#error').html('<b>Error:</b> No name set!');
    $('#error').show();
  } else if ($('#description').val() == '') {
    $('#error').html('<b>Error:</b> No description set!');
    $('#error').show();
  } else {
    get_text("/whoami", function(result) {
      var new_decision = {
        "person_id": +result['person_id'],
        "name": $("#name").val(),
        "description": $("#description").val(),
        "stage": +$("#stage").val(),
        "criterion_vote_style": $("#critStyle").val(),
        "alternative_vote_style": $("#altStyle").val(),
        "client_settings": "",
        "display_name": $('#displayName').val(),
        "criteria_instruction": $('#critInstructions').val(),
        "alternative_instruction": $('#altInstructions').val(),
        "image": $('#decImg').attr('src')
      }

      put_text("/decision/" + decisionID, JSON.stringify(new_decision), function(result) {
        if (result['error']) {
          $('#error').html('<b>Error:</b> ' + result['error']);
          $('#error').show();
        } else if (result['decision']) {
          updateLeftNav();
          $('#success').html('Updated Successfully');
          $('#success').show();

          if (new_decision['stage'] == 3) {
            $('#decisionResultsBtn').show();
            $('#decisionVotesBtn').show();
          }
        } else {
          $('#error').html('<b>Error:</b> Something went wrong :(');
          $('#error').show();
        }
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
  $('<li><a onclick="buildDecisionHome(' + decisionID + ')">Decision</a></li>').appendTo(ul);
  $('<li><a onclick="buildCustomizeDecision(' + decisionID + ')">Customize</a></li>').appendTo(ul);
  $('<li class="active"><a onclick="buildEditCriteria(' + decisionID + ')">Criteria</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditAlternative(' + decisionID + ')">Alternatives</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionInvite(' + decisionID + ')">Invite</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionStatus(' + decisionID + ')">Voting Status</a></li>').appendTo(ul);

  var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');
  //Table of existing criteria here
  $('<div id="critList">').appendTo(wrapper);
  $('#critList').hide();
  $('<hr/>').appendTo(wrapper);
  var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return false;').appendTo(wrapper);

  $('<div>').attr('id', 'success').addClass('alert alert-success').appendTo(form);
  $('#success').hide();
  $('<div>').attr('id', 'error').addClass('alert alert-danger').appendTo(form);
  $('#error').hide();

  $('<div>').attr('id', 'critForm').appendTo(form);
  showAddCriteria(decisionID);

  updateCritList(decisionID);
}

function showAddCriteria(decisionID) {
  $('#critForm').hide();
  $('#critForm').html("");
  $('<p class="text-Danger">Unable to add new criteria when not in development</p>').appendTo('#critForm');
  ifDecisionInDevelopment(decisionID, function() {
    $('#critForm').html("");
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
        .attr('rows', '3')
        .attr('name', 'critDesc')
        .attr('placeholder', 'Criterion Description')
        .attr('id', 'critDesc'))
      .appendTo('#critForm');

    $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('onclick', 'addCriteria(' + decisionID + ');').text('Add Criterion').appendTo('#critForm');
  });
  $('#critForm').show();
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
      .attr('rows', '3')
      .attr('name', 'critDesc')
      .attr('placeholder', 'Criterion Description')
      .attr('id', 'critDesc'))
    .appendTo('#critForm');

  $('<button>').addClass('btn btn-primary').attr('onclick', 'showAddCriteria(' + decisionID + ');').text('Back to Add New Criteria').attr('style', 'float: left').appendTo('#critForm');
  $('<button>').addClass('btn btn-primary').attr('onclick', 'editCriteria(' + decisionID + ', ' + criterionID + ');').text('Update Criterion').attr('style', 'float: right').appendTo('#critForm');
  $('<div>').addClass('clearfix').appendTo('#critForm'); //added to fix display issue
  get_text("/decision/" + decisionID + "/criterion/" + criterionID + "/info", function(result) {
    $('#critOrder').val(result['criterion']['order']);
    $('#critName').val(result['criterion']['name']);
    $('#critDesc').val(result['criterion']['description']);
  });
}

function addCriteria(decisionID) {
  $('#success').hide();
  $('#error').hide();

  var new_crit = {
    "name": $("#critName").val(),
    "description": $("#critDesc").val(),
    "order": +$('#critOrder').val()
  }

  post_text("/decision/" + decisionID + '/criterion', JSON.stringify(new_crit), function(result) {
    if (result['error']) {
      $('#error').html('<b>Error:</b> ' + result['error']);
      $('#error').show()
    } else if (result['criterion']) {
      updateCritList(decisionID);
      showAddCriteria(decisionID); //clears it
    } else {
      $('#error').html('<b>Error:</b> Something went wrong :(');
      $('#error').show();
    }
  });
}

function editCriteria(decisionID, criterionID) {
  $('#success').hide();
  $('#error').hide();

  var crit = {
    "name": $("#critName").val(),
    "description": $("#critDesc").val(),
    "order": +$('#critOrder').val()
  }

  put_text("/decision/" + decisionID + '/criterion/' + criterionID, JSON.stringify(crit), function(result) {
    if (result['error']) {
      $('#error').html('<b>Error:</b> ' + result['error']);
      $('#error').show()
    } else if (result['criterion']) {
      updateCritList(decisionID);
      showAddCriteria(decisionID); //clears it
    } else {
      $('#error').html('<b>Error:</b> Something went wrong :(');
      $('#error').show();
    }
  });
}

function deleteCriteria(decisionID, criterionID) {
  $('#success').hide();
  $('#error').hide();

  confirmYesNo(
    "Delete Criteria",
    "Are you sure you want to delete this criterion?",
    function() {
      delete_text("/decision/" + decisionID + "/criterion/" + criterionID, function(result) {
        if (result['error']) {
          $('#error').html('<b>Error:</b> ' + result['error']);
          $('#error').show()
        } else if (result['result'] == "deleted") {
          updateCritList(decisionID);
          showAddCriteria(decisionID);
        } else {
          $('#error').html('<b>Error:</b> Something went wrong :(');
          $('#error').show();
        }
      });
    },
    function() { /* Do nothing */ }
  );
}

function updateCritList(decisionID) {
  //clear it to repopulate it
  $('#critList').html("");

  get_text("/decision/" + decisionID + "/criterions", function(results) {
    var table = $('<table>').append($('<tbody>')).addClass('table table-striped').appendTo('#critList');
    table.append('<tr><th></th><th>Name</th><th>Description</th><th></th></tr>');

    if (results["criterions"].length < 1) $('#critList').hide();
    else $('#critList').show();

    for (var i in results["criterions"]) {
      c = results["criterions"][i];
      table.append('<tr><td>' + c['order'] + '</td><td>' + c['name'] + '</td><td>' + c['description'] + '</td><td>' + '<div style="width:45px; float:right;"><a onclick="showEditCriteria(' + decisionID + ',' + c['criterion_id'] + ');"><span class="glyphicon glyphicon-pencil text-Primary"></span></a>' + '<a onclick="deleteCriteria(' + decisionID + ',' + c['criterion_id'] + ');"><span class="glyphicon glyphicon-trash text-Danger" style="margin-left:10px;"></span></a></div></td></tr>');
    }

    //Hide dangerous icons if not in development
    $('.glyphicon.text-Danger').hide();
    ifDecisionInDevelopment(decisionID, function() {
      $('.glyphicon.text-Danger').show();
    });
  });
}

/**** Decision Alternative ****/
function buildEditAlternative(decisionID) {
  $('title').html('Edit Decision');
  clearContent();

  $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');

  var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
  $('<li><a onclick="buildDecisionHome(' + decisionID + ')">Decision</a></li>').appendTo(ul);
  $('<li><a onclick="buildCustomizeDecision(' + decisionID + ')">Customize</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditCriteria(' + decisionID + ')">Criteria</a></li>').appendTo(ul);
  $('<li class="active"><a onclick="buildEditAlternative(' + decisionID + ')">Alternatives</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionInvite(' + decisionID + ')">Invite</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionStatus(' + decisionID + ')">Voting Status</a></li>').appendTo(ul);

  var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');
  //Table of existing criteria here
  $('<div id="altList">').appendTo(wrapper);
  $('#altList').hide();
  $('<hr/>').appendTo(wrapper);
  var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return false;').appendTo(wrapper);

  $('<div>').attr('id', 'success').addClass('alert alert-success').appendTo(form);
  $('#success').hide();
  $('<div>').attr('id', 'error').addClass('alert alert-danger').appendTo(form);
  $('#error').hide();

  $('<div>').attr('id', 'altForm').appendTo(form);
  showAddAlternative(decisionID);

  updateAltList(decisionID);
}

function showAddAlternative(decisionID) {
  $('#altForm').hide();
  $('#altForm').html("");
  $('<p class="text-Danger">Unable to add new alternatives when not in development</p>').appendTo('#altForm');
  ifDecisionInDevelopment(decisionID, function() {
    $('#altForm').html("");
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
        .attr('rows', '3')
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

    $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('onclick', 'addAlternative(' + decisionID + ');').text('Add Alternative').appendTo('#altForm');
  });
  $('#altForm').show();
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
      .attr('rows', '3')
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

  $('<button>').addClass('btn btn-primary').attr('onclick', 'showAddAlternative(' + decisionID + ');').text('Back to Add New Alternative').attr('style', 'float: left').appendTo('#altForm');
  $('<button>').addClass('btn btn-primary').attr('onclick', 'editAlternative(' + decisionID + ', ' + alternativeID + ');').text('Update Alternative').attr('style', 'float: right').appendTo('#altForm');
  $('<div>').addClass('clearfix').appendTo('#altForm'); //added to fix display issue

  get_text("/decision/" + decisionID + "/alternative/" + alternativeID + "/info", function(result) {
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
    "name": $("#altName").val(),
    "description": $("#altDesc").val(),
    "cost": +$("#altCost").val(),
    "order": +$("#altOrder").val()
  }

  post_text("/decision/" + decisionID + '/alternative', JSON.stringify(new_alt), function(result) {
    if (result['error']) {
      $('#error').html('<b>Error:</b> ' + result['error']);
      $('#error').show()
    } else if (result['alternative']) {
      updateAltList(decisionID);
      showAddAlternative(decisionID); //clears it
    } else {
      $('#error').html('<b>Error:</b> Something went wrong :(');
      $('#error').show();
    }
  });
}

function editAlternative(decisionID, alternativeID) {
  $('#success').hide();
  $('#error').hide();

  var alt = {
    "name": $("#altName").val(),
    "description": $("#altDesc").val(),
    "cost": +$("#altCost").val(),
    "order": +$("#altOrder").val()
  }

  put_text("/decision/" + decisionID + '/alternative/' + alternativeID, JSON.stringify(alt), function(result) {
    if (result['error']) {
      $('#error').html('<b>Error:</b> ' + result['error']);
      $('#error').show()
    } else if (result['alternative']) {
      updateAltList(decisionID);
      showAddAlternative(decisionID); //clears it
    } else {
      $('#error').html('<b>Error:</b> Something went wrong :(');
      $('#error').show();
    }
  });
}

function deleteAlternative(decisionID, alternativeID) {
  $('#success').hide();
  $('#error').hide();

  confirmYesNo(
    "Delete Alternative",
    "Are you sure you want to delete this alternative?",
    function() {
      delete_text("/decision/" + decisionID + "/alternative/" + alternativeID, function(result) {
        if (result['error']) {
          $('#error').html('<b>Error:</b> ' + result['error']);
          $('#error').show()
        } else if (result['result'] == "deleted") {
          updateAltList(decisionID);
          showAddAlternative(decisionID);
        } else {
          $('#error').html('<b>Error:</b> Something went wrong :(');
          $('#error').show();
        }
      });
    },
    function() { /* Do nothing */ }
  );
}

function updateAltList(decisionID) {
  //clear it to repopulate it
  $('#altList').html("");

  get_text("/decision/" + decisionID + "/alternatives", function(results) {
    var table = $('<table>').append($('<tbody>')).addClass('table table-striped').appendTo('#altList');
    table.append('<tr><th></th><th>Name</th><th>Description</th><th>Cost</th><th></th></tr>');

    if (results["alternatives"].length < 1) $('#altList').hide();
    else $('#altList').show();

    for (var i in results["alternatives"]) {
      a = results["alternatives"][i];
      table.append('<tr><td>' + a['order'] + '</td><td>' + a['name'] + '</td><td>' + a['description'] + '</td><td>' + a['cost'] + '</td><td>' + '<div style="width:45px; float:right;"><a onclick="showEditAlternative(' + decisionID + ',' + a['alternative_id'] + ');"><span class="glyphicon glyphicon-pencil text-Primary"></span></a>' + '<a onclick="deleteAlternative(' + decisionID + ',' + a['alternative_id'] + ');"><span class="glyphicon glyphicon-trash text-Danger" style="margin-left:10px;"></span></a></div></td></tr>');
    }

    //Hide dangerous icons if not in development
    $('.glyphicon.text-Danger').hide();
    ifDecisionInDevelopment(decisionID, function() {
      $('.glyphicon.text-Danger').show();
    });
  });
}

/**** Decision Invite ****/
function buildDecisionInvite(decisionID) {
  $('title').html('Invite People');
  clearContent();

  $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');

  var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
  $('<li><a onclick="buildDecisionHome(' + decisionID + ')">Decision</a></li>').appendTo(ul);
  $('<li><a onclick="buildCustomizeDecision(' + decisionID + ')">Customize</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditCriteria(' + decisionID + ')">Criteria</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditAlternative(' + decisionID + ')">Alternatives</a></li>').appendTo(ul);
  $('<li class="active"><a onclick="buildDecisionInvite(' + decisionID + ')">Invite</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionStatus(' + decisionID + ')">Voting Status</a></li>').appendTo(ul);

  var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');

  $('<h2>').text('Invite Individuals').appendTo(wrapper);
  var form = $([
    '<form class ="form-signin" onsubmit = "return false" id="inviteForm">',
    '<div id= "invitation_sent" class = "alert alert-success"></div>',
    '<div id="invitation_error" class="alert alert-danger" style="display: none;"></div>',
    '<label for="bal_dec_id" >Decision Name: </label>',
    '<input type="text" class= "form-control" required="required" placeholder="Decision ID" id="dName"></input>',
    '<label for="bal_name">Name</label>',
    '<input type="text" id="i_name" class= "form-control" required="required" placeholder="Name"></input>',
    '<label for="bal_email">Email</label>',
    '<input type="email" id="i_email" class= "form-control" required="required" placeholder="Email"></input>',
    '<br/>',
    '</form>'
  ].join('\n'));

  $('<button>').addClass('btn btn-primary').attr('id', 'invite_submit').attr('onclick', 'sendInvite(' + decisionID + ');').attr('style', 'float: right;').append('<span> <i class="glyphicon glyphicon-envelope"></i>  Invite and Email</span>').appendTo(form);
  $('<button>').addClass('btn btn-primary').attr('id', 'invite_submit').attr('onclick', 'addBallot(' + decisionID + ');').attr('style', 'float: right; margin-right:10px;').append('<span> <i class="glyphicon glyphicon-user"></i>  Invite </span>').appendTo(form);

  $('<div>').addClass('clearFix').appendTo(form);

  get_text("/decision/" + decisionID + "/info", function(result) {
    $('#dName').val(result.decision.name)
  })

  $(wrapper).append(form);
  $("#invitation_sent").hide()
  $("#invitation_error").hide()

  $('<hr/>').appendTo(wrapper);
  $('<h2>').text('Invite Bulk').appendTo(wrapper);

  $('<div>').attr('id', 'bulk_sent').addClass('alert alert-success').appendTo(wrapper);
  $('#bulk_sent').hide();
  $('<div>').attr('id', 'bulk_error').addClass('alert alert-danger').appendTo(wrapper);
  $('#bulk_error').hide();

  $('<div class="form-group">').append(
      $('<label for="bulkEmails">Email List</label>'),
      $('<textarea>').addClass('form-control')
      .attr('rows', '3')
      .attr('name', 'bulkEmails')
      .attr('placeholder', 'one@email.com, two@email.com, three@email.com, etc')
      .attr('id', 'bulkEmails'))
    .appendTo(wrapper);

  $('<button>').addClass('btn btn-primary').attr('onclick', 'bulkAddBallot(' + decisionID + ');').text('Bulk Invite').appendTo(wrapper);
  $('<div>').addClass('clearfix').appendTo(wrapper);
}

function bulkAddBallot(decisionID) {
  $("#bulk_sent").hide()
  $("#bulk_error").hide()

  var emails = $('#bulkEmails').val().split(',');

  for (var i in emails) {
    if (!isEmail($.trim(emails[i]))) {
      $('#bulk_error').html('<b>Error:</b> Invalid email:' + $.trim(emails[i]));
      $('#bulk_error').show();
      break;
    }

    //Default name is just email
    new_invite = {
      "name": $.trim(emails[i].substr(0, emails[i].indexOf('@'))),
      "email": $.trim(emails[i])
    };

    //If one breaks, just stop trying
    var breakFlag = false;

    $('#bulk_sent').html('Added ballots: ');

    post_text("/decision/" + decisionID + "/ballot_silent", JSON.stringify(new_invite), function(result) {
      if (result['error']) {
        $('#bulk_error').html('<b>Error:</b> ' + result['error']);
        $('#bulk_error').show()
        breakFlag = true;
      } else if (result['ballot']) {
        $('#bulk_sent').append(result['ballot']['email'] + ' | ');
        $('#bulk_sent').show()
      } else {
        $('#bulk_error').html('<b>Error:</b> Something went wrong :(');
        $('#bulk_error').show();
        breakFlag = true;
      }
    });
    if (breakFlag) break;
  }
  $('#bulkEmails').val('');
}

function bulkInvite(decisionID) {
  $('#bulk_invite_sent').hide()
  $('#error').hide()

  get_text('/decision/' + decisionID + '/ballots', function(result) {
    ballots = result["ballots"];

    $('#bulk_invite_sent').html('Invited ballots: ');

    for (var i in ballots) {
      ballot = ballots[i];
      if (!ballot["sent"]) {
        $('#bulk_invite_sent').append(ballot["email"] + ' | ');
        get_text(ballot["url"] + "/invite", function(result) {
          if (result['error']) {
            $('#error').html('<b>Error:</b> ' + result['error']);
            $('#error').show()
          } else if (result['result'] == "invited") {
            $('#bulk_invite_sent').show();
          } else {
            $('#error').html('<b>Error:</b> Something went wrong :(');
            $('#error').show();
          }
        });

      }
    }
  });
}

//Adds Ballot and sends email with link
function sendInvite(decisionID) {
  $("#invitation_sent").hide()
  $("#invitation_error").hide()
  id = decisionID

  if ($('#i_name').val() == '') {
    $('#invitation_error').html('<b>Error:</b> No name set!');
    $('#invitation_error').show();
  } else if ($('#i_email').val() == '') {
    $('#invitation_error').html('<b>Error:</b> No email set!');
    $('#invitation_error').show();
  } else if (!isEmail($('#i_email').val())) {
    $('#invitation_error').html('<b>Error:</b> Invalid email!');
    $('#invitation_error').show()
  } else {
    new_invite = {
      "name": $('#i_name').val(),
      "email": $('#i_email').val()
    };
    post_text("/decision/" + id + "/ballot", JSON.stringify(new_invite), function(result) {
      if (result['error']) {
        $('#invitation_error').html('<b>Error:</b> ' + result['error']);
        $('#invitation_error').show()
      } else if (result['ballot']) {
        buildDecisionInvite(id);
        $('#invitation_sent').html('Invitation sent Successfully!');
        $('#invitation_sent').show()
      } else {
        $('#invitation_error').html('<b>Error:</b> Something went wrong :(');
        $('#invitation_error').show();
      }
    });
  }
}

//Adds Ballot without emailing
function addBallot(decisionID) {
  $("#invitation_sent").hide()
  $("#invitation_error").hide()
  id = decisionID

  if ($('#i_name').val() == '') {
    $('#invitation_error').html('<b>Error:</b> No name set!');
    $('#invitation_error').show();
  } else if ($('#i_email').val() == '') {
    $('#invitation_error').html('<b>Error:</b> No email set!');
    $('#invitation_error').show();
  } else if (!isEmail($('#i_email').val())) {
    $('#invitation_error').html('<b>Error:</b> Invalid email!');
    $('#invitation_error').show()
  } else {
    new_invite = {
      "name": $('#i_name').val(),
      "email": $('#i_email').val()
    };
    post_text("/decision/" + id + "/ballot_silent", JSON.stringify(new_invite), function(result) {
      if (result['error']) {
        $('#invitation_error').html('<b>Error:</b> ' + result['error']);
        $('#invitation_error').show()
      } else if (result['ballot']) {
        buildDecisionInvite(id);
        $('#invitation_sent').html('Added ballot Successfully!');
        $('#invitation_sent').show()
      } else {
        $('#invitation_error').html('<b>Error:</b> Something went wrong :(');
        $('#invitation_error').show();
      }
    });
  }
}

/**** Decision Status ****/
function buildDecisionStatus(decisionID) {
  $('title').html('Edit Decision');
  clearContent();

  $('<strong><i class="glyphicon glyphicon-cog"></i> Edit Decision</strong><hr/>').appendTo('#content');

  var ul = $('<ul>').addClass('nav nav-tabs').appendTo('#content');
  $('<li><a onclick="buildDecisionHome(' + decisionID + ')">Decision</a></li>').appendTo(ul);
  $('<li><a onclick="buildCustomizeDecision(' + decisionID + ')">Customize</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditCriteria(' + decisionID + ')">Criteria</a></li>').appendTo(ul);
  $('<li><a onclick="buildEditAlternative(' + decisionID + ')">Alternatives</a></li>').appendTo(ul);
  $('<li><a onclick="buildDecisionInvite(' + decisionID + ')">Invite</a></li>').appendTo(ul);
  $('<li class="active"><a onclick="buildDecisionStatus(' + decisionID + ')">Voting Status</a></li>').appendTo(ul);

  var wrapper = $('<div>').addClass('tabbedContent').appendTo('#content');
  var ballotsForCurrentStage = $('<div>').attr('id', 'totalBallots').appendTo(wrapper);

  $('#totalBallots').append(
    '<ul class="list-group" id="stage_ballot">' +
    '<li class="list-group-item active">Current Stage: <span id="stage_li"></span></li>' +
    '<li class="list-group-item"><span class="badge" id="voted_li"></span>Voted </li>' +
    '<li class="list-group-item"><span class="badge" id="not_voted_li"></span>Not Voted</li>' +
    '<li class="list-group-item"><span class="badge" id="total_li"></span><b>Total ballots</b> </li>' +
    '</ul>'
  );

  $('<div>').attr('id', 'success').addClass('alert alert-success').appendTo(wrapper);
  $('#success').hide();
  $('<div>').attr('id', 'error').addClass('alert alert-danger').appendTo(wrapper);
  $('#error').hide();

  $('<div>').attr('id', 'bulk_invite_sent').addClass('alert alert-success').appendTo(wrapper);
  $('#bulk_invite_sent').hide();

  $('<div>').css('margin-bottom', '20px').addClass('clearFix').append($('<button>').addClass('btn btn-primary').attr('style', 'float: right;').attr('onclick', 'bulkInvite(' + decisionID + ');').text('Send All')).appendTo(wrapper);



  $('<div>').attr('id', 'formDiv').appendTo(wrapper);

  $('<div>').attr('id', 'statusTable').appendTo(wrapper).css('margin-top', '10px');
  buildStatusTable(decisionID);

  //resets table every 5 seconds :)
  intervalID = window.setInterval(function() {
    buildStatusTable(decisionID);
  }, 5000);
}

function buildStatusTable(decisionID) {

  var status_table = $('<table class="table table-striped" id="s_table">');
  var s_body = $('<tbody id="s_body">');

  status_table.append(s_body);
  s_body.append(
    '<tr>' +
    '<th>Name</th>' +
    '<th>Email</th>' +
    '<th id="status_header">Ballot status</th>' +
    '<th>Action</th>' +
    '<th>Reset</th>' +
    '<th></th>' +
    '<th></th>' +
    '</tr>'
  );

  get_text("/decision/" + decisionID + "/ballots", function(result) {
    var vote_status = ""
    var totalBallots = 0
    var voted = 0
    var not_voted = 0
    if (result.ballots != null) {
      $('#s_table').show()

      totalBallots = result.ballots.length
      getBallotStage(decisionID)

      for (var i = 0; i < result.ballots.length; i++) {
        if (result.ballots[i].rating != null && result.ballots[i].rating.length > 0 && result.ballots[i].rating.length != "undefined") {
          vote_status = "Voted";
          voted = voted + 1;
        } else {
          vote_status = "Not Voted";
          not_voted = not_voted + 1;
        }

        var url = result.ballots[i].url

        if (result.ballots[i].sent == 0) {
          s_body.append('<tr><td>' + result.ballots[i].name + '</td><td>' + result.ballots[i].email + '</td><td>Not Emailed</td><td> <a onclick="resendEmail(' + decisionID + ', \'' + url + '\', \'' + result.ballots[i].email + '\');">Send Email</a>' + '</td><td align="center"><a onclick=\'resetVote(' + decisionID + ',' + JSON.stringify(result.ballots[i]) + ')\'><span class="glyphicon glyphicon-repeat"></span></a></td><td><a onclick=\'buildEditBallotForm(' + decisionID + ',' + JSON.stringify(result.ballots[i]) + ')\'><span class="glyphicon glyphicon-pencil text-Primary"></span></a></td><td><a onclick=deleteBallot(\'' + url + '\',' + decisionID + ')><span class="glyphicon glyphicon-trash text-Danger" style="margin-left: 10px; display: inline-block;"></span></a></td>');
        } else {
          s_body.append('<tr><td>' + result.ballots[i].name + '</td><td>' + result.ballots[i].email + '</td><td>' + vote_status + '</td><td> <a onclick="resendEmail(' + decisionID + ', \'' + url + '\', \'' + result.ballots[i].email + '\');">Resend Email</a>' + '</td><td align="center"><a onclick=\'resetVote(' + decisionID + ',' + JSON.stringify(result.ballots[i]) + ')\'><span class="glyphicon glyphicon-repeat"></span></a></td><td><a onclick=\'buildEditBallotForm(' + decisionID + ',' + JSON.stringify(result.ballots[i]) + ')\'><span class="glyphicon glyphicon-pencil text-Primary"></span></a></td><td><a onclick=deleteBallot(\'' + url + '\',' + decisionID + ')><span class="glyphicon glyphicon-trash text-Danger" style="margin-left: 10px; display: inline-block;"></span></a></td>');
        }
      }
    } else {
      getBallotStage(decisionID);
    }

    //clear old table
    $('#statusTable').html('');
    //set new one
    $('#statusTable').append(status_table);

    //Update counts
    $('#voted_li').text(voted);
    $('#not_voted_li').text(not_voted);
    $('#total_li').text(totalBallots);
  });

  sortTable();
}

function getBallotStage(decisionID) {
  get_text("/decision/" + decisionID + "/info", function(result) {

    if (result.decision.stage == 1) {
      $("#stage_li").text('In Development ');
    } else if (result.decision.stage == 2) {
      $("#stage_li").text('Voting in progress ');
    } else if (result.decision.stage == 3) {
      $("#stage_li").text('Completed ');
    } else if (result.decision.stage == 4) {
      $("#stage_li").text('Locked ');
    } else {
      $("#stage_li").text('Unknown ');
    }
  })

}

function resendEmail(decisionID, url, email) {
  $('#success').hide();
  $('#error').hide();

  get_text(url + "/invite", function(result) {
    if (result['error']) {
      $('#error').html('<b>Error:</b> ' + result['error']);
      $('#error').show()
    } else if (result['result'] == "invited") {
      buildDecisionStatus(decisionID);
      $('#success').html('Invite sent to ' + email);
      $('#success').show();
    } else {
      $('#error').html('<b>Error:</b> Something went wrong :(');
      $('#error').show();
    }
  });
}

function deleteBallot(url, decisionID) {
  $('#success').hide();
  $('#error').hide();

  confirmYesNo(
    "Delete Ballot",
    "Are you sure you want to delete this ballot?",
    function() {
      delete_text(url, function(result) {
        if (result['error']) {
          $('#error').html('<b>Error:</b> ' + result['error']);
          $('#error').show()
        } else if (result['result'] == "deleted") {
          buildStatusTable(decisionID);
        } else {
          $('#error').html('<b>Error:</b> Something went wrong :(');
          $('#error').show();
        }
      });
    },
    function() { /* Do nothing */ }
  );
}

function editBallot(decisionID, url) {
  $('#success').hide();
  $('#error').hide();

  var ballot = {
    "name": $("#ballotName").val(),
    "email": $("#ballotEmail").val(),
  }

  put_text(url, JSON.stringify(ballot), function(result) {
    if (result['error']) {
      $('#error').html('<b>Error:</b> ' + result['error']);
      $('#error').show()
    } else if (result['ballot']) {
      //Just clear the div
      $('#formDiv').html('');
      //refresh table
      buildStatusTable(decisionID);
    } else {
      $('#error').html('<b>Error:</b> Something went wrong :(');
      $('#error').show();
    }
  });
}

function cancelEdit(decisionID) {
  //Just clear the div
  $('#formDiv').html('');
}

function resetVote(decisionID, ballot) {
  $('#success').hide();
  $('#error').hide();

  console.log(ballot.votes)
  confirmYesNo(
    "Reset Ballot",
    "Are you sure you want to reset this ballot?",
    function() {
      //delete alt rating on each criteria
      for (i = 0; i < ballot.votes.length; i++) {
        delete_text(ballot.url + "/alternative/" + ballot.votes[i].alternative_id + "/criterion/" + ballot.votes[i].criterion_id + "/vote", function(result) {
          if (result['error']) {
            $('#error').html('<b>Error:</b> ' + result['error']);
            $('#error').show()
          } else if (result['result'] == "deleted") {
            buildStatusTable(decisionID);
            $('#success').html('Votes cleared!');
            $('#success').show();
          } else {
            $('#error').html('<b>Error:</b> Something went wrong :(');
            $('#error').show();
          }
        });
      }

      //delete criterion rating 
      for (i = 0; i < ballot.rating.length; i++) {
        delete_text(ballot.url + "/criterion/" + ballot.rating[i].criterion_id + "/vote", function(result) {
          if (result['error']) {
            $('#error').html('<b>Error:</b> ' + result['error']);
            $('#error').show()
          } else if (result['result'] == "deleted") {
            $('#success').html('Ratings cleared!');
            $('#success').show();
          } else {
            $('#error').html('<b>Error:</b> Something went wrong :(');
            $('#error').show();
          }
        });
      }
      buildDecisionStatus(decisionID);
    },
    function() { /* Do nothing */ }
  );
}

function buildEditBallotForm(decisionID, ballot) {
  //edit ballot form
  $('#formDiv').html('');
  $('#formDiv').append('<hr/>');
  var form = $([
    '<form class ="form-signin" onsubmit = "return false" id="editBallotForm">',

    '<label for="bal_dec_id" >Decision Name: </label>',
    '<input type="text" class= "form-control" required="required" placeholder="Decision ID" id="bName"></input>',
    '<br />',

    '<label for="bal_name">Name</label>',
    '<input type="text" id="ballotName" class= "form-control" required="required" placeholder="Name"></input>',
    '<br />',
    '<label for="bal_email">Email</label>',
    '<input type="email" id="ballotEmail" class= "form-control" required="required" placeholder="Email"></input>',
    '<br />',
    '</form>'
  ].join('\n'));
  form.appendTo('#formDiv');
  $('<button>').addClass('btn btn-primary').attr('onclick', 'editBallot(' + decisionID + ',\'' + ballot.url + '\');').attr('style', 'float: right').append('<span>  Update </span>').appendTo(form);
  $('<button>').addClass('btn btn-primary').attr('onclick', 'cancelEdit(' + decisionID + ')').append('<span> Cancel </span>').appendTo(form);
  $('#formDiv').append('<hr/>');

  get_text("/decision/" + decisionID + "/info", function(result) {
    $('#bName').val(result.decision.name);
  })
  $('#ballotName').val(ballot.name);
  $('#ballotEmail').val(ballot.email);
}

//functions to sort the table by ballot status
function sortTable() {
  $('#status_header').click(function() {  
    var table =  $(this).parents('table').eq(0);
    var rows = table.find('tr:gt(0)').toArray().sort(comparer($(this).index()));
    this.asc = !this.asc;   
    if (!this.asc) {
      rows = rows.reverse();
    }   
    for (var i = 0; i < rows.length; i++) {
      table.append(rows[i]);
    }
  })
}

function comparer(index) {  
  return function(a, b) {      
    var valA = getCellValue(a, index),
      valB = getCellValue(b, index);
    return $.isNumeric(valA) && $.isNumeric(valB) ?  valA - valB  : valA.localeCompare(valB)  ;
  }
}

function getCellValue(row, index) {
  return $(row).children('td').eq(index).html();
}

//Runs cb if true
function ifDecisionInDevelopment(decisionID, cb) {
  get_text("/decision/" + decisionID + "/info", function(result) {
    if (+result['decision']['stage'] == 1) cb();
  });
}