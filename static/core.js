var base_url = "http://localhost:9999";

/* This version also returns on 403! */
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

/* This version also returns on 403 */
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
