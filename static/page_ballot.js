var criterion_names = [];
var alternative_names = [];
var criterion_ids = [];
var alt_ids = [];

function main(body) {
  $('title').html('Ballot');

  $('head').append('<script src="/static/nouislider.min.js"></script>');
  $.loadCSS('/static/css/ballot.css');
  $.loadCSS('/static/css/nouislider.min.css');
  $.loadCSS('/static/css/nouislider.pips.css');
  $.loadCSS('/static/css/nouislider.tooltips.css');

  var ballot_info = getBallotCookies();

  if (ballot_info == null || ballot_info == undefined) {
    $("body")
      .append("<h1>Error: No ballot found.  Please make sure cookies are enabled and try the link again</h1>");
  } else {

    var decision_id = ballot_info[0];
    var ballot_id = ballot_info[1];
    var decision_name;
    var decision_desc;
    var decision_stage;
    var criterion_descriptions = [];
    var alternative_descriptions = [];
    var voter_name;
    var page;
    var crit_vote_style;
    var crit_instructions;
    var alt_instructions;
    var dec_img;
    var decision = get_decision(decision_id);

    if (decision == null || decision["error"] != null || decision["error"] != undefined) {
      alert("Unable to get decision information");
      return;
    }

    var ballot = get_ballot(decision_id, ballot_id);

    if (ballot == null || ballot["error"] != null || ballot["error"] != undefined) {
      alert("Unable to get ballot information");
      return;
    }

    var criterion = get_criterion(decision);

    if (criterion == null || criterion["error"] != null || criterion["error"] != undefined) {
      alert("Unable to get criterion information");
      return;
    }

    var alternatives = get_alternatives(decision);

    if (alternatives == null || alternatives["error"] != null || alternatives["error"] != undefined) {
      alert("Unable to get alternative information");
      return;
    }

    var votes = get_votes(decision_id, ballot_id);

    if (votes == null || votes["error"] != null || votes["error"] != undefined) {
      alert("Unable to retrieve votes");
      return;
    }

    decision_name = decision.name;
    decision_desc = decision.description;
    decision_stage = decision.stage;
    crit_vote_style = decision.criterion_vote_style;
    alt_vote_style = decision.alternative_vote_style;
    voter_name = ballot.name;
    crit_instructions = decision.criteria_instruction;
    alt_instructions = decision.alternative_instruction;
    dec_img = decision.image;

    for (var i = 0; i < criterion.length; i++) {
      criterion_names[i] = criterion[i].name;
      criterion_descriptions[i] = criterion[i].description;
      criterion_ids[i] = criterion[i].criterion_id;
    }

    for (var i = 0; i < alternatives.length; i++) {
      alternative_names[i] = alternatives[i].name;
      alternative_descriptions[i] = alternatives[i].description;
      alt_ids[i] = alternatives[i].alternative_id;
    }

    //build page		
    if (decision_stage == "1") {
      errorPage("Voting has not started yet, please check back later.");
      return;
    } else if (decision_stage == "3") {
      errorPage("Voting for this decision is now closed!");
      return;
    } else if (decision_stage == "4") {
      errorPage("This decision has been locked.");
      return;
    } else if (votes.length != 0) {
      errorPage("You have already voted in this decision, contact the facilitator if you wish to change your vote.");
    } else {
      var page = "<div id=\"topbar\" class=\"navbar navbar-default navbar-fixed-top\">" + "<div class=\"container\">" + "<a class=\"navbar-brand\">" + decision_name + "</a>" + "</div>" + "</div>" + "<div id=\"ballotbody\" class=\"container\">" + "<div class=\"row\">" + "<div class=\"col-md-6 col-md-offset-3\" id=\"topRow\">" + "<img id=\"decision_image\" src=\"" + dec_img + "\"></img>" + "<h3>Welcome, " + voter_name + "! </h3>" + "<p class=\"lead\">" + decision_desc + "</p>" + "<div id='crit_inst' class='partone'>" + crit_instructions + "</div>" + "</div>" + "</div>" + "<form class=\"form-horizontal center partone\" role=\"form\">";

      //check to see if they want sliders or buttons
      if (crit_vote_style == "b") {
        for (var i = 0; i < criterion_names.length; i++) {
          page += "<div class=\"row\">" + "<div class=\"form-group\">" + "<label id=\"crit" + i + "lab\" class=\"col-sm-4 control-label criterion\">" + criterion_names[i] + "</label>" + "<div class=\"col-sm-8 btn-group\" data-toggle=\"buttons\">";
          for (var j = 1; j < 11; j++) {
            page += "<label class=\"btn btn-default votebtn\">" + "<input type=\"radio\" name=\"crit" + parseInt(i) + "\" id=\"crit" + i + "_" + j + "\" value=\"" + j + "\">" + j + "</label>";
          }
          page += "</div>" + "<div class=\"alert alert-success center\" id=\"crit" + i + "Desc\">" + criterion_descriptions[i] + "</div>" + "</div>" + "</div>";
        }
      } else {
        page += "<div class=\"container center\">";
        for (var i = 0; i < criterion_names.length; i++) {
          page += "<div class=\"row\">" + "<label for=\"crit" + i + "\" id=\"crit" + i + "lab\" class=\"criterion\">" + criterion_names[i] + "</label>" + "<span id=\"crit" + i + "slider-val\" class=\"slider-val\"></span>" + "<div class=\"alert alert-success center\" id=\"crit" + i + "Desc\">" + criterion_descriptions[i] + "</div>" + "<div id=\"crit" + i + "\" class=\"slider\"></div>" + "</div>";
        }
      }

      page += "<div class=\"alert alert-danger\" id=\"criterrordiv\"></div>" + "<button class=\"btn btn-primary\" id=\"contbtn\" type=\"button\">Continue</button>" + "</form>" + "</div>" + "<div class=\"container parttwo\" id=\"alternative_table\">" + "<div id='alt_inst'>" + alt_instructions + "</div>" + "<div class=\"row\">" + "<table class=\"table\">" + "<thead>" + "<tr>" + "<th></th>"

      for (var i = 0; i < criterion_names.length; i++) {
        page += "<th id=\"critT" + i + "\">" + criterion_names[i] + "<div class=\"alert alert-success center\" id=\"critT" + i + "Desc\">" + criterion_descriptions[i] + "</div>" + "</th>";
      }
      page += "</thead>" + "<tbody>";

      for (var i = 0; i < alternative_names.length; i++) {
        page += "<tr>" + "<td id=\"Alt" + i + "\" class=\"alternative\">" + alternative_names[i] + "<div class=\"alert alert-success center\" id=\"Alt" + i + "Desc\">" + alternative_descriptions[i] + "</div>" + "</td>";

        if (alt_vote_style == "5") {
          for (var j = 0; j < criterion_names.length; j++) {
            page += "<td>" + "<div class=\"dropdown\">" + "<button class=\"btn btn-default dropdown-toggle center\" type=\"button\" id=\"dropdownMenu" + i + j + "\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">" + "<div id=\"alt" + i + "crit" + j + "_color\" class=\"color_pick\"></div>" + "</button>" + "<ul class=\"dropdown-menu\" aria-labelledby=\"dropdownMenu" + i + j + "\">" + "<li id=\"alt" + i + "crit" + j + "g\" class=\"vote-green\">Very Well-Aligned</li>" + "<li id=\"alt" + i + "crit" + j + "y\" class=\"vote-greenyellow\">Well Aligned</li>" + "<li id=\"alt" + i + "crit" + j + "y\" class=\"vote-yellow\">Neutral</li>" + "<li id=\"alt" + i + "crit" + j + "r\" class=\"vote-orange\">Poorly Aligned</li>" + "<li id=\"alt" + i + "crit" + j + "r\" class=\"vote-red\">Very Poorly-Aligned</li>" + "</ul>" + "</div>" + "</td>";
          }
        } else { //3
          for (var j = 0; j < criterion_names.length; j++) {
            page += "<td>" + "<div class=\"dropdown\">" + "<button class=\"btn btn-default dropdown-toggle center\" type=\"button\" id=\"dropdownMenu" + i + j + "\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">" + "<div id=\"alt" + i + "crit" + j + "_color\" class=\"color_pick\"></div>" + "</button>" + "<ul class=\"dropdown-menu\" aria-labelledby=\"dropdownMenu" + i + j + "\">" + "<li id=\"alt" + i + "crit" + j + "g\" class=\"vote-green\">Very Well-Aligned</li>" + "<li id=\"alt" + i + "crit" + j + "y\" class=\"vote-yellow\">Neutral</li>" + "<li id=\"alt" + i + "crit" + j + "r\" class=\"vote-red\">Very Poorly-Aligned</li>" + "</ul>" + "</div>" + "</td>";
          }
        }
        page += "</tr>";
      }
      page += "</tbody>" + "</table>" + "</div>" + "</div>" + "</div>" + "<div id=\"ballotbottom\" class=\"container parttwo\">" + "<div class=\"col-md-6 col-md-offset-3\" id=\"bottomRow\">" + "<div class=\"alert alert-danger\" id=\"errordiv\"></div>" + "<button class=\"btn btn-primary\" id=\"submitbtn\">Submit</button>" + "<button class=\"btn btn-warning\" id=\"clearbtn\">Clear</button>" + "</div>" + "</div>" + "</div>" + "<div class=\"alert alert-success\" id=\"successdiv\">Your vote has been received, thanks!</div>"

      $("body").append(page);
      $(".parttwo").hide();

      //build slider bars
      if (crit_vote_style != 'b') {
        var sliders = $('.slider');
        var spans = $('.slider-val');

        for (var i = 0; i < sliders.length; i++) {
          noUiSlider.create(sliders[i], {
            start: [1],
            behaviour: 'tap',
            connect: 'lower',
            step: 1,
            range: {
              'min': [1],
              'max': [10]
            },
            format: {
              to: function(value) {
                return value;
              },
              from: function(value) {
                return value;
              }
            }
          });
          //Text Descriptions
          if (crit_vote_style == 't') {
            bindValuesText(sliders[i], spans[i]);
          } else { //Number Descriptions
            bindValuesNumbers(sliders[i], spans[i]);
          }
        }
      }

      //event handlers
      $('#clearbtn').click(function(event) {
        location.reload();
      });

      $('#contbtn').click(function(event) {
        var crit_votes = [];

        $("#criterrordiv").hide();
        $("#criterrordiv").empty();

        if (crit_vote_style == 'b') {
          for (var i = 0; i < criterion_names.length; i++) {
            crit_votes.push($("input[name=crit" + i + "]:checked").val());

            if (crit_votes[i] == undefined) {
              $("#criterrordiv").html("Please vote on all criterion<br>");
              $("#criterrordiv").show();
            }
          }
        } else {
          for (var i = 0; i < criterion_names.length; i++) {
            crit_votes.push($("#crit" + i + "slider-val").text());

            if (crit_votes[i] == undefined) {
              $("#criterrordiv").html("Please vote on all criterion<br>");
              $("#criterrordiv").show();
            }
          }
        }

        if ($("#criterrordiv").is(":empty")) {
          $(".partone").hide();
          $(".parttwo").show();
        }
      });

      $('#submitbtn').click(function(event) {
        var crit_votes = [];
        var alt_votes = [];
        var row_votes = [];
        var alt_missing = false;

        $("#successdiv").hide();
        $("#errordiv").empty();
        $("#errordiv").hide();

        //collect criterion votes
        if (crit_vote_style == 'b') {
          for (var i = 0; i < criterion_names.length; i++) {
            crit_votes.push($("input[name=crit" + i + "]:checked").val());
            if (crit_votes[i] == undefined) {
              $("#errordiv").html("Please vote on all criterion<br>");
              $("#errordiv").show();
            }
          }
        } else { //Type of slider
          for (var i = 0; i < criterion_names.length; i++) {
            crit_votes.push($("#crit" + i + "slider-val").data('val'));
            if (crit_votes[i] == undefined) {
              $("#errordiv").html("Please vote on all criterion<br>");
              $("#errordiv").show();
            }
          }
        }

        //collect alternative votes
        for (var i = 0; i < alternative_names.length; i++) {
          alt_votes[i] = [];

          for (var j = 0; j < criterion_names.length; j++) {
            if ($("#alt" + i + "crit" + j + "_color").css("background-color") == "rgb(255, 0, 0)") { //RED
              alt_votes[i].push("1");
            }
            if ($("#alt" + i + "crit" + j + "_color").css("background-color") == "rgb(255, 165, 0)") { //ORANGE
              alt_votes[i].push("2");
            }
            if ($("#alt" + i + "crit" + j + "_color").css("background-color") == "rgb(255, 255, 0)") { //YELLOW
              alt_votes[i].push("3");
            }
            if ($("#alt" + i + "crit" + j + "_color").css("background-color") == "rgb(173, 255, 47)") { //GREENYELLOW
              alt_votes[i].push("4");
            }
            if ($("#alt" + i + "crit" + j + "_color").css("background-color") == "rgb(0, 128, 0)") { //GREEN
              alt_votes[i].push("5");
            }
            if ($("#alt" + i + "crit" + j + "_color").css("background-color") == "rgb(128, 128, 128)") { //GREY
              alt_missing = true;
            }
          }
        }

        if (alt_missing == true) {
          $("#errordiv").append("Please assign a rating to all boxes");
          $("#errordiv").show();
        }

        //if ballot is complete, send votes
        if ($("#errordiv").is(":empty")) {
          for (var i = 0; i < criterion_names.length; i++) {
            var vote1 = vote_criterion(decision_id, ballot_id, criterion_ids[i], crit_votes[i]);

            if (vote1["error"] != null) {
              $("#errordiv").html("Your vote has already been recorded");
              $("#errordiv").show();
            }
          }

          for (var i = 0; i < alternative_names.length; i++) {
            for (var j = 0; j < criterion_names.length; j++) {
              var vote2 = rate_alternative(decision_id, ballot_id, alt_ids[i], criterion_ids[j], alt_votes[i][j]);
              if (vote2["error"] != null) {
                $("#errordiv").html("Your vote has already been recorded");
                $("#errordiv").show();
              }
            }
          }
          if ($("#errordiv").is(":empty")) {
            $(".parttwo").hide();
            $("#successdiv").show();
          }
        }
      });

      $('.criterion').click(function(event) {
        var id = "#" + this.id.slice(0, -3);
        $(id + "Desc").toggle();
        $('.alert').not(id + "Desc").hide();
      });

      $('.alternative').click(function(event) {
        var id = "#" + this.id;
        $(id + "Desc").toggle();
        $('.alert').not(id + "Desc").hide();
      });

      $('.vote-red').click(function(event) {
        var id = "#" + this.id.slice(0, -1) + "_color";
        $(id).css('background-color', 'red');
      });

      $('.vote-orange').click(function(event) {
        var id = "#" + this.id.slice(0, -1) + "_color";
        $(id).css('background-color', 'orange');
      });

      $('.vote-yellow').click(function(event) {
        var id = "#" + this.id.slice(0, -1) + "_color";
        $(id).css('background-color', 'yellow');
      });

      $('.vote-greenyellow').click(function(event) {
        var id = "#" + this.id.slice(0, -1) + "_color";
        $(id).css('background-color', 'greenyellow');
      });

      $('.vote-green').click(function(event) {
        var id = "#" + this.id.slice(0, -1) + "_color";
        $(id).css('background-color', 'green');
      });
    }
  }
}

function bindValuesText(slider, span) {
  slider.noUiSlider.on('update', function(values, handle) {
    var test = [
      '0 -', //0 not used
      '1 - Not at all important',
      '2 - Slightly important',
      '3 - Slightly important',
      '4 - Moderately important',
      '5 - Moderately important',
      '6 - Important',
      '7 - Important',
      '8 - Very important',
      '9 - Very important',
      '10 - Extremely important'
    ];
    span.innerHTML = test[values[handle]];
    $('#' + span.id).data('val', values[handle]);
  });
}

function bindValuesNumbers(slider, span) {
  slider.noUiSlider.on('update', function(values, handle) {
    span.innerHTML = values[handle];
    $('#' + span.id).data('val', values[handle]);
  });
}


function getBallotCookies() {
  var my_cookies = document.cookie.split('; ')
  var decision_id;
  var ballot_id;
  var ret = null;

  if (my_cookies != "") {
    for (var i = 0; i < my_cookies.length; i++) {
      var temp = my_cookies[i].split('=');
      if (temp[0] == "ballot_id") {
        ballot_id = temp[1];
      } else if (temp[0] == "decision_id") {
        decision_id = temp[1];
      }
    }
    ret = [decision_id, ballot_id];
  }
  return ret;
}

//retrieve decision 
function get_decision(decision_id) {
  var result = null;
  $.ajax({
    type: "GET",
    url: base_url + "/decision/" + decision_id + "/info",
    contentType: 'application/json; charset=utf-8',
    async: false,
    success: function(r) {
      result = r["decision"];
    },
    error: function(r) {
      errmsg = JSON.parse(r.responseText);
      result = errmsg;
    }
  });
  return result;
}

//retrieve ballot
function get_ballot(decision_id, ballot_id) {
  var result = null;
  $.ajax({
    type: "GET",
    url: base_url + "/decision/" + decision_id + "/ballot/" + ballot_id + "/info",
    contentType: 'application/json; charset=utf-8',
    async: false,
    success: function(r) {
      result = r["ballot"];
    },
    error: function(r) {
      errmsg = JSON.parse(r.responseText);
      result = errmsg;
    }
  });
  return result;
}

//get votes from a ballot
function get_votes(decision_id, ballot_id) {
  var result = null;
  $.ajax({
    type: "GET",
    url: base_url + "/decision/" + decision_id + "/ballot/" + ballot_id + "/votes",
    contentType: 'application/json; charset=utf-8',
    async: false,
    success: function(r) {
      result = r["votes"];
    },
    error: function(r) {
      errmsg = JSON.parse(r.responseText);
      result = errmsg;
    }
  });
  return result;
}

//retrieve criterion
function get_criterion(decision) {
  var result = null;
  $.ajax({
    type: "GET",
    url: base_url + "/decision/" + decision["decision_id"] + "/criterions",
    contentType: 'application/json; charset=utf-8',
    async: false,
    success: function(r) {
      result = r["criterions"];
    },
    error: function(r) {
      errmsg = JSON.parse(r.responseText);
      result = errmsg;
    }
  });
  return result;
}

//retrieve alternatives
function get_alternatives(decision) {
  var result = null;
  $.ajax({
    type: "GET",
    url: base_url + "/decision/" + decision["decision_id"] + "/alternatives",
    contentType: 'application/json; charset=utf-8',
    async: false,
    success: function(r) {
      result = r["alternatives"];
    },
    error: function(r) {
      errmsg = JSON.parse(r.responseText);
      result = errmsg;
    }
  });
  return result;
}

//vote on criterion
function vote_criterion(decision_id, ballot_id, criterion_id, vote) {
  var result = null;
  $.ajax({
    type: "GET",
    url: base_url + "/decision/" + decision_id + "/ballot/" + ballot_id + "/criterion/" + criterion_id + "/vote/" + vote,
    contentType: 'application/json; charset=utf-8',
    async: false,
    success: function(r) {
      result = r["rating"];
    },
    error: function(r) {
      errmsg = JSON.parse(r.responseText);
      result = errmsg;
    }
  });
  return result;
}

//rate alternative
function rate_alternative(decision_id, ballot_id, alt_id, crit_id, vote) {
  var result = null;
  $.ajax({
    type: "GET",
    url: base_url + "/decision/" + decision_id + "/ballot/" + ballot_id + "/alternative/" + alt_id + "/criterion/" + crit_id + "/vote/" + vote,
    contentType: 'application/json; charset=utf-8',
    async: false,
    success: function(r) {
      result = r["vote"];
    },
    error: function(r) {
      errmsg = JSON.parse(r.responseText);
      result = errmsg;
    }
  });
  return result;
}