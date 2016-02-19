var base_url = "http://localhost:9999";

get_text = function(url, cb) {
  var request = new XMLHttpRequest();
  request.open('GET', base_url+url, true);
  request.setRequestHeader("Content-Type", "application/json");

  request.onreadystatechange = function() {
    if(request.readyState == 4 && (request.status == 200 || request.status == 403)) {
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
    if(request.readyState == 4 && (request.status == 200 || request.status == 403)) {
      cb(JSON.parse(request.responseText));
    }
  }

  request.send(data);
}

put_text = function(url, data, cb) {
	var request = new XMLHttpRequest();
	request.open('PUT', base_url+url, true);
	request.setRequestHeader("Content-Type", "application/json");

	request.onreadystatechange = function() {
		if(request.readyState == 4 && (request.status == 200 || request.status == 403)) {
			cb(JSON.parse(request.responseText));
		}
	}

	request.send(data);
}

delete_text = function(url, cb) {
  var request = new XMLHttpRequest();
  request.open('DELETE', base_url+url, true);
  request.setRequestHeader("Content-Type", "application/json");

  request.onreadystatechange = function() {
    if(request.readyState == 4 && (request.status == 200 || request.status == 403)) {
      cb(JSON.parse(request.responseText));
    }
  }

  request.send();
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

function confirmYesNo(title, msg, yesFn, noFn) {
    var $confirm = $("#modalConfirmYesNo");
    $confirm.modal('show');
    $("#lblTitleConfirmYesNo").html(title);
    $("#lblMsgConfirmYesNo").html(msg);
    $("#btnYesConfirmYesNo").off('click').click(function () {
        yesFn();
        $confirm.modal("hide");
    });
    $("#btnNoConfirmYesNo").off('click').click(function () {
        noFn();
        $confirm.modal("hide");
    });
}

function loggedIn(cb) {
  get_text("/whoami", function (result) {
    if(result['error'] == "unauthenticated") {
      alert('Please Login!');
      window.location.replace("/login.html");
    }
    else {
      cb();
    }
  });
}

function isEmail(email) {
  var regex = /^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$/;
  return regex.test(email);
}

function getParameterByName(name, url) {
  if (!url) url = window.location.href;
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
    results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, " "));
}