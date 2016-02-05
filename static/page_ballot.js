var base_url = "http://localhost:9999";

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
		var criterion_names = [];
		var criterion_descriptions = [];
		var alternative_names = [];
		var alternative_descriptions = [];


		$("body")
		.append("Decision_id: "+decision_id+"<br>"+"Ballot_id: "+ballot_id+"<br>");
		
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

				$("body")
				.append("name: "+decision_name+"<br>"+"desc: "+decision_desc+"<br>"+"stage: "+decision_stage+"<br>");

				get_text("/decision/"+decision_id+"/criterions", function (crits) {

					for(var i=0; i<crits.criterions.length; i++) {

						criterion_names[i] = crits.criterions[i].name;
						criterion_descriptions[i] = crits.criterions[i].description;

					}

					for(var j=0; j<crits.criterions.length; j++) {

						$("body")
						.append("Crit "+j+" name: "+criterion_names[j]+"<br>");

						$("body")
						.append("Crit "+j+" desc: "+criterion_descriptions[j]+"<br>");

					}
			
					get_text("/decision/"+decision_id+"/alternatives", function(alts) {

						for(var i=0; i<alts.alternatives.length; i++) {

							alternative_names[i] = alts.alternatives[i].name;
							alternative_descriptions[i] = alts.alternatives[i].description;
				
						}
	
						for(var j=0; j<alts.alternatives.length; j++) {
							
							$("body")
							.append("Alt "+j+" name: "+alternative_names[j]+"<br>");

							$("body")
							.append("Alt "+j+" desc: "+alternative_descriptions[j]+"<br>");
						}	

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
