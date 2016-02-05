function main(body) {
  checkLogin();

  //Add CSS
  $.loadCSS('static/css/login.css');

  //Build page
  buildPage();
}

function login() {
  $('#status').hide();

  new_login = {
    "email":document.getElementById('username').value,
    "password":document.getElementById('password').value}

  post_text("/login", JSON.stringify(new_login), function (result) {
    get_text("/whoami", function (result) {
      //redirect
      window.location.replace("/");
    });
  });

  $('#status').html('<b>Error:</b> Unable to sign in');
  $('#status').show();

  //Avoid refreshing
  return false;
}

function buildPage() {
  $('title').html('Please Login!');

  $('<div>').addClass('wrapper').appendTo('body');
  $('<div>').attr('id','header').appendTo('.wrapper');
    var form = $('<form>').addClass('form-signin').attr('onsubmit', 'return login()').appendTo('.wrapper');
      $('<h2>').addClass('form-signin-heading').text('Please Login').appendTo(form);
      $('<div>').attr('id','status').addClass('alert alert-danger').appendTo(form);
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
      $('<a>').attr('href', '#').text('Forgot Password?').appendTo(cbBox);
      $('<button>').addClass('btn btn-lg btn-primary btn-block').attr('type', 'submit').text('Submit').appendTo(form);
}

function checkLogin() {
  get_text("/whoami", function (result) {
    if(result) {
      //Note: we may need to do a check and log-off to avoid an inf loop
      alert('You are already logged in, delete me and redirect after dev!');
      //window.location.replace("/");
    }
  });
}
  
/* -- Everything below here should be in core.js -- */
var base_url = "http://localhost:9999";

get_text = function(url, cb) {
  var request = new XMLHttpRequest();
  request.open('GET', base_url+url, true);
  request.setRequestHeader("Content-Type", "application/json");

  request.onreadystatechange = function() {
    if(request.readyState == 4 && request.status == 200) {
      cb(JSON.parse(request.responseText));
    }
  }

  request.send();
}

post_text = function(url, data, cb) {
  var request = new XMLHttpRequest();
  request.open('POST', base_url+url, true);
  request.setRequestHeader("Content-Type", "application/json");

  request.onreadystatechange = function() {
    if(request.readyState == 4 && request.status == 200) {
      cb(JSON.parse(request.responseText));
    }
  }

  request.send(data);
}

function assert(condition, message) {
  if (!condition) {
    throw message || "Assertion failed";
  }
}

function delay() {
  var now = new Date().getTime();
  while(new Date().getTime() < now + 500){ /* do nothing */ } 
}

jQuery.loadCSS = function(url) {
  if (!$('link[href="' + url + '"]').length)
      $('head').append('<link rel="stylesheet" type="text/css" href="' + url + '">');
}
