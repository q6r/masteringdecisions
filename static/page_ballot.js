var base_url = "http://localhost:9999";
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

		
		get_text("/decision/"+decision_id+"/info", function(dec) {

			decision_name = dec.decision.name;
			decision_desc = dec.decision.description;
			decision_stage = dec.decision.stage;

			if(decision_stage == "1") {
				$("body")
				.append("<h1>Voting has not started yet, please check back later</h1>");
				
				return;
			}
			else if (decision_stage == "3") {
				$("body")
				.append("<h1>Voting has been closed!</h1>");

				return;
			}
			else {


				get_text("/decision/"+decision_id+"/criterions", function (crits) {

					for(var i=0; i<crits.criterions.length; i++) {

						criterion_names[i] = crits.criterions[i].name;
						criterion_descriptions[i] = crits.criterions[i].description;
						criterion_ids[i] = crits.criterions[i].criterion_id;
					}

			
					get_text("/decision/"+decision_id+"/alternatives", function(alts) {

						for(var i=0; i<alts.alternatives.length; i++) {

							alternative_names[i] = alts.alternatives[i].name;
							alternative_descriptions[i] = alts.alternatives[i].description;
							alt_ids[i] = alts.alternatives[i].alternative_id;
				
						}
	
						var page="<div id=\"topbar\" class=\"navbar navbar-default navbar-fixed-top\">"
						+"<div class=\"container\">"
						+"<a class=\"navbar-brand\">"+decision_name+"</a>"
						+"</div>"
						+"</div>"
						+"<div id=\"ballotbody\" class=\"container\">"
						+"<div class=\"row\">"
						+"<div class=\"col-md-6 col-md-offset-3\" id=\"topRow\">"
						+"<p class=\"lead\">"+decision_desc+"</p>"
						+"<p>How important are each of this decision's criterion to you?</p>"
						+"<p>Please rate each one between 1 (not at all important) and 10 (vital)</p>"
						+"<p>Click on a criterion name to view its description</p>"
						+"</div>"
						+"</div>"
						+"<form class=\"form-horizontal\" role=\"form\">";

						for(var i=0; i<criterion_names.length; i++) {
						
							page+="<div class=\"row\">"
							+"<div class=\"form-group\">"
							+"<label id=crit"+i+" class=\"col-sm-4 control-label criterion\">"+criterion_names[i]+"</label>"
							+"<div class=\"col-sm-8 btn-group\" data-toggle=\"buttons\">";

							for(var j=1; j<11; j++) {

								page+="<label class=\"btn btn-default votebtn\">"+"<input type=\"radio\" name=\"crit"+parseInt(i)+"\" id=\"crit"+i+"_"+j+"\" value=\""+j+"\">"+j+"</label>";


							}

							page+="</div>"
							+"<div class=\"alert alert-success center\" id=\"crit"+i+"Desc\">"+criterion_descriptions[i]+"</div>"
							+"</div>"
							+"</div>";
			
						}

						page+="</form>"
						+"</div>"
						+"<div class=\"container\" id=\"alternative_table\">"
						+"<div class=\"row\">"
						+"<p>For each alternative below, rate how you feel each criterion will impact it.</p>"
						+"<p>Red represents least beneficial, yellow represents neutral, and green represents most beneficial</p>"
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
						+"<button class=\"btn btn-primary\" id=\"submitbtn\">Submit</button>"
						+"<button class=\"btn btn-warning\" id=\"clearbtn\">Clear</button>"
						+"<div class=\"alert alert-danger\" id=\"errordiv\"></div>"
						+"<div class=\"alert alert-success\" id=\"successdiv\">Your vote has been received, thanks!</div>"
						+"</div>"
						+"</div>"
						+"</div>"
					

						$("body").append(page);

						// Event Handlers -------------------------------------------------------

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
							for(var i=0; i<criterion_names.length; i++) {
	
								crit_votes.push( $("input[name=crit"+i+"]:checked").val());

								if(crit_votes[i] == undefined) {
									$("#errordiv").html("Please vote on all criterion<br>");
									$("#errordiv").show();
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
								
								//console.log(crit_votes);
								//console.log(alt_votes);

								for(var i=0; i<criterion_names.length; i++) {
									
									get_text("/decision/"+decision_id+"/ballot/"+ballot_id+"/criterion/"+criterion_ids[i]+"/vote/"+crit_votes[i], function(crits) {
										if(crits.error) {
											$("#errordiv").html("Your vote has already been recorded");
											$("#errordiv").show();
										}
										else {

											$('#successdiv').show();

										}
									});
								}


								for(var i=0; i<alternative_names.length; i++) {

									for(var j=0; j<criterion_names.length; j++) {

										
										get_text("/decision/"+decision_id+"/ballot/"+ballot_id+"/alternative/"+alt_ids[i]+"/criterion/"+criterion_ids[j]+"/vote/"+alt_votes[i][j], function(alts) {
											
											if(alts.error) {
												$("#errordiv").html("Your vote has already been recorded");
												$("#errordiv").show();
											}
											else {
											
												$('successdiv').show();

											}
								
										});
										

									}

								}	


							}
						
						});


						$('.criterion').click(function(event) {
							var id = "#" + this.id;
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
					


					});

				});
				


			}

		});



	}


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
		if(request.readyState == 4 && request.status == 200) {
			cb(JSON.parse(request.responseText));
		}
	}

	request.send(data);
}

