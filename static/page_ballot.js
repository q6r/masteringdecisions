var criterion_names = [];
var alternative_names = [];
var criterion_ids = [];
var alt_ids = [];

function main(body)
{
	$('title').html('Ballot');
	$('head').prepend(custom_header());

	var ballot_info = getBallotCookies();

	if(ballot_info == null || ballot_info == undefined) {
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

		var decision = get_decision(decision_id);

		//Check for undefined ids here.

		//gather information

		if(decision == null || decision["error"] != null || decision["error"] != undefined) {
			alert("Unable to get decision information");
			return;
		}
		
		var ballot = get_ballot(decision_id, ballot_id);


		if(ballot == null || ballot["error"] != null || ballot["error"] != undefined) {
			alert("Unable to get ballot information");
			return;
		}

		var criterion = get_criterion(decision);

		if(criterion == null || criterion["error"] != null || criterion["error"] != undefined) {
			alert("Unable to get criterion information");
			return;
		}

		var alternatives = get_alternatives(decision);
		
		if(alternatives == null || alternatives["error"] != null || alternatives["error"] != undefined) {
			alert("Unable to get alternative information");
			return;
		}

		decision_name = decision.name;
		decision_desc = decision.description;
		decision_stage = decision.stage;
		crit_vote_style = decision.criterion_vote_style;
		voter_name = ballot.name;

		for(var i=0; i<criterion.length; i++) {
			criterion_names[i] = criterion[i].name;
			criterion_descriptions[i] = criterion[i].description;
			criterion_ids[i] = criterion[i].criterion_id;
		}				

		for(var i=0; i<alternatives.length; i++) {
			alternative_names[i] = alternatives[i].name;
			alternative_descriptions[i] = alternatives[i].description;
			alt_ids[i] = alternatives[i].alternative_id;

		}

		//build page		
		
		if(decision_stage == "1") {
			
			$("body").append("<h1>Voting has not started yet, please check back later</h1>");
			return;
		}	
		else if (decision_stage == "3") {
		
			$("body").append("<h1>Voting for this decision is now closed!</h1>");
			return;

		}
		else {

			var page="<div id=\"topbar\" class=\"navbar navbar-default navbar-fixed-top\">"
			+"<div class=\"container\">"
			+"<a class=\"navbar-brand\">"+decision_name+"</a>"
			+"</div>"
			+"</div>"
			+"<div id=\"ballotbody\" class=\"container\">"
			+"<div class=\"row\">"
			+"<div class=\"col-md-6 col-md-offset-3\" id=\"topRow\">"
			+"<h3>Welcome, "+voter_name+"! </h3>"
			+"<p class=\"lead\">"+decision_desc+"</p>"
			+"<p>How important are each of this decision's criterion to you?</p>"
			+"<p>Please rate each one between 1 (not at all important) and 10 (vital)</p>"
			+"<p>Keep in mind you can use the same rating as many times as you like.</p>"
			+"<p>Click on a criterion name to view its description</p>"
			+"</div>"
			+"</div>"
			+"<form class=\"form-horizontal\" role=\"form\">";

			//check to see if they want sliders or buttons

			if(crit_vote_style == "b") {
			

				for(var i=0; i<criterion_names.length; i++) {
				
					page+="<div class=\"row\">"
					+"<div class=\"form-group\">"
					+"<label id=\"crit"+i+"lab\" class=\"col-sm-4 control-label criterion\">"+criterion_names[i]+"</label>"
					+"<div class=\"col-sm-8 btn-group\" data-toggle=\"buttons\">";
	
					for(var j=1; j<11; j++) {
	
						page+="<label class=\"btn btn-default votebtn\">"+"<input type=\"radio\" name=\"crit"+parseInt(i)+"\" id=\"crit"+i+"_"+j+"\" value=\""+j+"\">"+j+"</label>";
	
					}
	
					page+="</div>"
					+"<div class=\"alert alert-success center\" id=\"crit"+i+"Desc\">"+criterion_descriptions[i]+"</div>"
					+"</div>"
					+"</div>";
			
				}
			} 
			else {

				page+="<div class=\"container center\">";

				for(var i=0; i<criterion_names.length; i++) {

					page+="<div class=\"row\">"
					+"<label for=\"crit"+i+"\" id=\"crit"+i+"lab\" class=\"criterion\">"+criterion_names[i]+"</label>"	
					+"<span id=\"crit"+i+"slider-val\" class=\"slider-val\"></span>"
					+"<div class=\"alert alert-success center\" id=\"crit"+i+"Desc\">"+criterion_descriptions[i]+"</div>"
					+"<div id=\"crit"+i+"\" class=\"slider\"></div>"
					+"</div>";
				}
		

			}
				

			page+="</form>"
			+"</div>"
			+"<div class=\"container\" id=\"alternative_table\">"
			+"<div class=\"row\">"
			+"<p>For each alternative below, rate how you feel each criterion will impact it.</p>"
			+"<p>Red represents least beneficial, yellow represents neutral, and green represents most beneficial</p>"
			+"<p>Feel free to use the same color as many times (or as few times) as you like!</p>"
			+"<p>Click on an alternative name to view its description</p>"
			+"<table class=\"table\">"
			+"<thead>"
			+"<tr>"
			+"<th></th>"

			for(var i=0; i<criterion_names.length; i++) {
			
				page+="<th id=\"critT"+i+"\">"
				+criterion_names[i]
				+"<div class=\"alert alert-success center\" id=\"critT"+i+"Desc\">"
				+criterion_descriptions[i]
				+"</div>"
				+"</th>";	

			}

			page+="</thead>"
			+"<tbody>";

			for(var i=0; i<alternative_names.length; i++) {
	
				page+="<tr>"
				+"<td id=\"Alt"+i+"\" class=\"alternative\">"
				+alternative_names[i]
				+"<div class=\"alert alert-success center\" id=\"Alt"+i+"Desc\">"
				+alternative_descriptions[i]
				+"</div>"
				+"</td>";

				for(var j=0; j<criterion_names.length; j++) {

					page+="<td>"
					+"<div class=\"dropup\">"
					+"<button class=\"btn btn-default dropdown-toggle center\" type=\"button\" id=\"dropdownMenu"+i+j+"\" data-toggle=\"dropdown\" aria-haspopup=\"true\" aria-expanded=\"false\">"
					+"<div id=\"alt"+i+"crit"+j+"_color\" class=\"color_pick\"></div>"
					+"</button>"
					+"<ul class=\"dropdown-menu\" aria-labelledby=\"dropdownMenu"+i+j+"\">"
					+"<li id=\"alt"+i+"crit"+j+"r\" class=\"color1\">Red</li>"
					+"<li id=\"alt"+i+"crit"+j+"y\" class=\"color2\">Yellow</li>"
					+"<li id=\"alt"+i+"crit"+j+"g\" class=\"color3\">Green</li>"
					+"</ul>"
					+"</div>"
					+"</td>";
				

				}
								
				page+="</tr>";

			}
						
			page+="</tbody>"
			+"</table>"
			+"</div>"
			+"</div>"
			+"<div id=\"ballotbottom\" class=\"container\">"
			+"<div class=\"col-md-6 col-md-offset-3\" id=\"bottomRow\">"
			+"<div class=\"alert alert-danger\" id=\"errordiv\"></div>"
			+"<div class=\"alert alert-success\" id=\"successdiv\">Your vote has been received, thanks!</div>"
			+"<button class=\"btn btn-primary\" id=\"submitbtn\">Submit</button>"
			+"<button class=\"btn btn-warning\" id=\"clearbtn\">Clear</button>"
			+"</div>"
			+"</div>"
			+"</div>"
			

			$("body").append(page);
			
			//build slider bars

			if(crit_vote_style != 'b') {

				var sliders = $('.slider');
				var spans = $('.slider-val');
	
				for(var i=0; i<sliders.length; i++) {
					
					noUiSlider.create(sliders[i], {
	
						start:[1],
						behaviour:'tap',
						connect:'lower',
						step:1,
						range: {
							'min':[1],
							'max':[10]
		
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
					bindValues(sliders[i], spans[i]);
	
				}
			}


			//event handlers

			$('#clearbtn').click(function(event) {
				location.reload();
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
				
				if(crit_vote_style == 'b') {

					for(var i=0; i<criterion_names.length; i++) {
	
						crit_votes.push( $("input[name=crit"+i+"]:checked").val());
	
						if(crit_votes[i] == undefined) {
							$("#errordiv").html("Please vote on all criterion<br>");
							$("#errordiv").show();
						}

					}
				}
				else {

					for(var i=0; i<criterion_names.length; i++) {

						crit_votes.push( $("#crit"+i+"slider-val").text());

						if(crit_votes[i] == undefined) {
							$("#errordiv").html("Please vote on all criterion<br>");
							$("#errordiv").show();
						}
					}


				}

				//collect alternative votes
				for(var i=0; i<alternative_names.length; i++) {

					alt_votes[i] = [];	

					for(var j=0; j<criterion_names.length; j++) {


						if( $("#alt"+i+"crit"+j+"_color").css("background-color") == "rgb(255, 0, 0)") {
							alt_votes[i].push("1");
						}
						if( $("#alt"+i+"crit"+j+"_color").css("background-color") == "rgb(255, 255, 0)") {
							alt_votes[i].push("3");
						}
						if( $("#alt"+i+"crit"+j+"_color").css("background-color") == "rgb(0, 128, 0)") {
							alt_votes[i].push("5");
						}
						if( $("#alt"+i+"crit"+j+"_color").css("background-color") == "rgb(128, 128, 128)") {
							alt_missing = true;
						}

					}

				}
							
				if(alt_missing == true) {
					$("#errordiv").append("Please vote on all alternatives");
					$("#errordiv").show();
				}

				//if ballot is complete, send votes
				if($("#errordiv").is(":empty")) {
								
					for(var i=0; i<criterion_names.length; i++) {
							
						var vote1 = vote_criterion(decision_id, ballot_id, criterion_ids[i], crit_votes[i]);


						if(vote1["error"] != null ) {
								$("#errordiv").html("Your vote has already been recorded");
								$("#errordiv").show();
						}
					
					}


					for(var i=0; i<alternative_names.length; i++) {

						for(var j=0; j<criterion_names.length; j++) {

							var vote2 = rate_alternative(decision_id, ballot_id, alt_ids[i], criterion_ids[j], alt_votes[i][j]);

							if(vote2["error"] != null) {
								$("#errordiv").html("Your vote has already been recorded");
								$("#errordiv").show();
							}
										

						}

					}

					if($("#errordiv").is(":empty")) {
						$("#successdiv").show();
					}	


				}
						
			});


			$('.criterion').click(function(event) {
				var id = "#" + this.id.slice(0,-3);
				$(id+"Desc").toggle();
				$('.alert').not(id +"Desc").hide();
			});

			$('.alternative').click(function(event) {
				var id = "#" + this.id;
				$(id+"Desc").toggle();
				$('.alert').not(id +"Desc").hide();
			});


			$('.color1').click(function(event) {
				var id = "#"+this.id.slice(0,-1)+"_color";
				$(id).css('background-color','red');
			});

			$('.color2').click(function(event) {
				var id = "#"+this.id.slice(0,-1)+"_color";
				$(id).css('background-color','yellow');
			});	

			$('.color3').click(function(event) {
				var id = "#"+this.id.slice(0,-1)+"_color";
				$(id).css('background-color','green');
			});

		}

	}


}

function bindValues(slider, span) {

	slider.noUiSlider.on('update', function(values, handle) {

		span.innerHTML = values[handle];

	});
	
}


function getBallotCookies() {

	var my_cookies = document.cookie.split(';')

	if(my_cookies == "") {
		return null;
	}

	var cookie_1 = my_cookies[0].split('=');
	var cookie_2 = my_cookies[1].split('=');

	var ballot_id;
	var decision_id;
	var ret;

	if(cookie_1[0] == "ballot_id") {
		ballot_id = cookie_1[1];
		decision_id = cookie_2[1];		

	} else if(cookie_1[0] == "decision_id") {
		ballot_id = cookie_2[1];
		decision_id = cookie_1[1];
	} 

	return ret = [decision_id, ballot_id];
}

function custom_header() {
	var content='<link rel="stylesheet" type="text/css" href="/static/css/ballot.css">'
	return content
}

//retrieve decision 
function get_decision(decision_id) {
	var result = null;
	$.ajax({
		type: "GET",
		url: "http://localhost:9999/decision/"+decision_id+"/info",
		contentType: 'application/json; charset=utf-8',
		async: false,
		success: function (r) {
			result = r["decision"];
		},
		error: function (r) {
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
		url: "http://localhost:9999/decision/"+decision_id+"/ballot/"+ballot_id+"/info",
		contentType: 'application/json; charset=utf-8',
		async: false,
		success: function (r) {
			result = r["ballot"];
		},
		error: function (r) {
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
		url: "http://localhost:9999/decision/"+decision["decision_id"]+"/criterions",
		contentType: 'application/json; charset=utf-8',
		async: false,
		success: function (r) {
			result = r["criterions"];
		},
		error: function (r) {
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
		url: "http://localhost:9999/decision/"+decision["decision_id"]+"/alternatives",
		contentType: 'application/json; charset=utf-8',
		async: false,
		success: function (r) {
			result = r["alternatives"];
		},
		error: function (r) {
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
		url: "http://localhost:9999/decision/"+decision_id+"/ballot/"+ballot_id+"/criterion/"+criterion_id+"/vote/"+vote,
		contentType: 'application/json; charset=utf-8',
		async: false,
		success: function (r) {
			result = r["rating"];
		},
		error: function (r) {
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
		url: "http://localhost:9999/decision/"+decision_id+"/ballot/"+ballot_id+"/alternative/"+alt_id+"/criterion/"+crit_id+"/vote/"+vote,
		contentType: 'application/json; charset=utf-8',
		async: false,
		success: function (r) {
			result = r["vote"];
		},
		error: function (r) {
			errmsg = JSON.parse(r.responseText);
			result = errmsg;
		}
	});
	return result;
}

