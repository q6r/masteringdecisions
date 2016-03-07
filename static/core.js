var base_url = "http://localhost:9999";

get_text = function(url, cb) {
  $.ajax({
    type: "GET",
    url: base_url + url,
    contentType: 'application/json; charset=utf-8',
    success: function(r) {
      cb(r);
    },
    error: function(r) {
      cb(JSON.parse(r.responseText));
    }
  });
}

post_text = function(url, data, cb) {
  $.ajax({
    type: "POST",
    data: data,
    url: base_url + url,
    contentType: 'application/json; charset=utf-8',
    success: function(r) {
      cb(r);
    },
    error: function(r) {
      cb(JSON.parse(r.responseText));
    }
  });
}

put_text = function(url, data, cb) {
  $.ajax({
    type: "PUT",
    data: data,
    url: base_url + url,
    contentType: 'application/json; charset=utf-8',
    success: function(r) {
      cb(r);
    },
    error: function(r) {
      cb(JSON.parse(r.responseText));
    }
  });
}

delete_text = function(url, cb) {
  $.ajax({
    type: "DELETE",
    url: base_url + url,
    contentType: 'application/json; charset=utf-8',
    success: function(r) {
      cb(r);
    },
    error: function(r) {
      cb(JSON.parse(r.responseText));
    }
  });
}

function assert(condition, message) {
  if (!condition) {
    throw message || "Assertion failed";
  }
}

function delay() {
  var now = new Date().getTime();
  while (new Date().getTime() < now + 500) { /* do nothing */ }
}

jQuery.loadCSS = function(url) {
  if (!$('link[href="' + url + '"]').length)
    $('head').append('<link rel="stylesheet" type="text/css" href="' + url + '">');
}

function getScript(src, callback) {
  var s = document.createElement('script');
  s.src = src;
  s.async = true;
  s.onreadystatechange = s.onload = function() {
    if (!callback.done && (!s.readyState || /loaded|complete/.test(s.readyState))) {
      callback.done = true;
      callback();
    }
  };
  document.querySelector('head').appendChild(s);
}

function confirmYesNo(title, msg, yesFn, noFn) {
  var $confirm = $("#modalConfirmYesNo");
  $confirm.modal('show');
  $("#lblTitleConfirmYesNo").html(title);
  $("#lblMsgConfirmYesNo").html(msg);
  $("#btnYesConfirmYesNo").off('click').click(function() {
    yesFn();
    $confirm.modal("hide");
  });
  $("#btnNoConfirmYesNo").off('click').click(function() {
    noFn();
    $confirm.modal("hide");
  });
}

function loggedIn(cb) {
  get_text("/whoami", function(result) {
    if (result['error'] == "unauthenticated") {
      window.location.replace("/login.html");
    } else {
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

function errorPage(error){
  $("body").attr("STYLE","BACKGROUND:#A4A4A4 !IMPORTANT").html('<H2><CENTER><IMG SRC="/static/images/logo.png" ALIGN="TOP" /><H2 STYLE="margin-top:50px;">'+error+'</H2>');
}
