var base_url = "http://localhost:9999";

function main(body)
{
	$('title').html('Results');
	$('head').prepend(custom_header());

	
	var decision_id = getDecisionId();

	if(decision_id == null || decision_id == undefined) {
		$("body")
		.append("<h1>Error: No ballot found.  Please make sure cookies are enabled and try the link again</h1>");	
	} else {
                //Varible List
                //Ballot Varibles
                var bname = undefined;
                var ratings = [];
                var votes = [];
                var url = undefined;
                var email = undefined;
                //This should look like
                //  [Crit_Id, [Weight,], [Alternative_Id, [Rating,]],]
                var info_set = [];
                //This should look like
    		// [[Crit_id, Weighted Avg, [Alternative_id, weighted_sum, total]]]
		var weighted_set = [];

		//$("body")
		
		get_text("/decision/"+decision_id+"/ballots", function(dec) {

			//Pull out info from Ballots
                        bname = dec.name;
                        votes = dec.votes;
                        ratings = dec.rating;
                        email = dec.email;
                        url = dec.url;
                        
			info_set = initPairedList(ratings, votes);
                        //Match up all the votes and ratings with criteron
			weighted_set = buildWeightedSums(info_set);
                        
//------------------------
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
			else if (decision_stage == "4") {
				$("body")
				.append("<h1>This decision has been archived.</h1>");

				return;
			}
			else {

				//$("body")
				//.append("name: "+decision_name+"<br>"+"desc: "+decision_desc+"<br>"+"stage: "+decision_stage+"<br>");

				get_text("/decision/"+decision_id+"/criterions", function (ballot_list) {

					for(var i=0; i<crits.criterions.length; i++) {

						criterion_names[i] = crits.criterions[i].name;
						criterion_descriptions[i] = crits.criterions[i].description;

					}
					//PAGE

			
						//PAGE REFERENCE	
						var page="<div id=\"topbar\" class=\"navbar navbar-default navbar-fixed-top\">"
						+"<div class=\"container\">"
						+"<a class=\"navbar-brand\">"+decision_name+"</a>"
						+"</div>"
						+"</div>"
						+"<div id=\"ballotbody\" class=\"container\">"
						+"<div class=\"row\">"
						+"<div class=\"col-md-6 col-md-offset-3\" id=\"topRow\">"
						+"<p class=\"lead\">"+decision_desc+"</p>"
						+"<p>Instructions for criterion</p>"
						+"<p>Click on a criterion name for description</p>"
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
						+"</div>";

						$("body").append(page);

						
					});

				});
				


			}
a
		});




	}
}
//---------------------------------------
//Set up the Paired_List
//Note: Rating and Weight have been named backwards on Backend so it looks backwards within this function
function initPairedList(ratings, votes){
	//Local Varibles
        var ret = [];

        var flag1 = 0;
	var flag2 = 0;
	var flag3 = 0;

	//Initialize all
        //First: votes array
	for (i = 0; i < votes.length; ++i){
		for (k = 0; k < ret.length; ++k){
			if (ret[k][0] == votes[i].criterion_id){
				flag1 = 1;
				for (j = 0; j < ret[k][2].length){
					if (ret[k][2][j][1] == votes[i].alternative_id){
						flag2 = 1;
						ret[k][2][j][1].push(votes[i].weight);
					}
				}
				if (flag2 == 0){
					ret[k][2].push([votes[i].alternative_id, [votes[i].weight]]);
				}
				else{
					flag2 = 0;
				}
			}
		}
                if (flag1 == 0){
			ret.push([votes[i].criterion_id, [], [votes[i].alternative_id, [votes[i].weight]]]);
                }
                else{
			flag1 = 0;
        	}
	}
	//Second: rating array
	for (h = 0; h < ratings.length; ++h){
		for (g = 0; g < ret.length; ++g){
			if (ret[g][0] == ratings[h].criterion_id){
				flag3 = 1;
				ret[g][1].put(rating[h].rating);
			}
		}
                if (flag3 == 0){
			ret.push([ratings[h].criterion_id, [rating[h].rating], []]);
                }
                else{
			flag3 = 0;
        	}

	}
	
	return ret;

}

//Input
//[[Crit_Id, [Weight,], [Alternative_Id, [Rating,],]],]
//Output
//[Crit_id, Weighted Avg, Weighted Total, [Alternative_id, weighted_sum, total_sum]]]
function buildWeightedSums(info_set){
	//Local Variables
	var ret = [];
	var weight = undefined;
	var weightedSums = [];
        var z = 0;
        var weigth_total = 0;

	//Build Weights
	for (z = 0; z < info_set.length; ++z){
                weight = mathWeight(into_set[z][1], weight_total);
                weight_total = weight_total + weight[1];
		weightedSums = mathNormalized(into_set[z][2], weight[0]);
		ret.push([into_set[z][0], weight[0], weight[1], weightedSums]);
	}
}

function mathWeight(weight_list, total){
	//Local Variables
  	var ret = [];
        var avg  = 0;
	var i = 0;

	//Build Average Weight
  	for(i = 0; i < weight_list.length; ++i){
		avg = avg + weigth_list;
  	}
        if( avg == 0 ){
		ret = [0, total];
		return ret;
 	}
        total = total + avg
	avg = avg / i;
	ret = [avg, total];
	return ret;
}

function mathNormalized(alt_list, weight ){
        if (alt_list == undefined){
		return undefined;
        }
	//Local Variables
	var ret = [];
        var alt_id = alt_list[0];
        var value = 0;
        var avg = 0;
	var i = 0;

	//Build Normalized Ratings
	for (i = 0; i < alt_list[1].length; ++i){
		value = value + alt_list[1][i];
	}
        avg = value / i;
        ret = [alt_id, avg, value];
	return ret;
}

function getDecisionId() {

	var my_cookies = document.cookie.split(';')

	if(my_cookies == "") {
		return null;
	}

	var decision_id = my_cookies[0].split('=');

	return decision_id;
}

function custom_header() {
	var content='<link rel="stylesheet" type="text/css" href="/static/css/ballot.css">'
	return content;
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
