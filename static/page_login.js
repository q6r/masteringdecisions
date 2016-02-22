function main(body) {
  checkLogin();
  $.loadCSS('static/css/login.css');
  buildPage();
}

function buildPage() {
  $('title').html('Please Login!');

  $('<div>').addClass('wrapper').appendTo('body');
  $('<div>').attr('id', 'header').appendTo('.wrapper');
  var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return login()').appendTo('.wrapper');
  $('<h2>').addClass('form-signin-heading').text('Please Login').appendTo(form);
  $('<div>').attr('id', 'status').addClass('alert alert-danger').appendTo(form);
  $('#status').hide();
  $('<input type="text" />').addClass('form-control')
    .attr('name', 'username')
    .attr('placeholder', 'Email Address')
    .attr('required', '')
    .attr('autofocus', '')
    .attr('id', 'username')
    .appendTo(form);
  $('<input type="password" />').addClass('form-control')
    .attr('name', 'password')
    .attr('placeholder', 'Password')
    .attr('required', '')
    .attr('id', 'password')
    .appendTo(form);
  var cbBox = $('<div>').addClass('cbBox');
  var rememberMe = $('<label class="checkbox">').addClass('rememberMeCB');
  $('<input type="checkbox">').attr('value', 'remember-me').attr('id', 'rememberMe').attr('name', 'rememberMe').appendTo(rememberMe);
  rememberMe.append('Remember Me');
  rememberMe.appendTo(cbBox);
  cbBox.appendTo(form);
  $('<a>').click(function() {
    alert("Please email the admin to request a password reset.");
  }).text('Forgot Password?').appendTo(cbBox);
  $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('type', 'submit').text('Submit').appendTo(form);
}

//If credentials are correct login and redirect, else show error
function login() {
  $('#status').hide();

  new_login = {
    "email": $('#username').val(),
    "password": $('#password').val()
  };

  post_text("/login", JSON.stringify(new_login), function(result) {
    if (result['status'] == "logged in") {
      get_text("/whoami", function(result) {
        window.location.replace("/");
      });
    } else if (result['error']) {
      $('#status').html('<b>Error:</b> ' + result['error']);
      $('#status').show();
    } else {
      $('#status').html('<b>Error:</b> Something went wrong :(');
      $('#status').show();
    }
  });

  $('#password').val('');
  //Avoid refreshing
  return false;
}

//If the user is already logged in redirect them to home
function checkLogin() {
  get_text("/whoami", function(result) {
    if (result['person_id']) {
      window.location.replace("/");
    }
  });
}