function main(body) {
  checkLogin();

  //Add CSS
  $.loadCSS('static/css/login.css');

  //Build page
  buildPage();
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

function login() {
  $('#status').hide();

  new_login = {
    "email":document.getElementById('username').value,
    "password":document.getElementById('password').value}

  post_text("/login", JSON.stringify(new_login), function (result) {
    if(result['status'] == "logged in") {
      //Authenticate and redirect
      get_text("/whoami", function (result) {
        //redirect
        window.location.replace("/");
      });
    }
    else if(result['error']) {
      $('#status').html('<b>Error:</b> ' + result['error']);
      $('#status').show();
    }
    else {
      $('#status').html('<b>Error:</b> Something went wrong :(');
      $('#status').show();
    }
  });

  //Avoid refreshing
  return false;
}

function checkLogin() {
  get_text("/whoami", function (result) {
    if(result['person_id']) {
      //Note: we may need to do a check and log-off to avoid an inf loop
      alert('You are already logged in, delete me and redirect after dev!');
      //window.location.replace("/");
    }
  });
}