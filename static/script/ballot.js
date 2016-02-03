//Hello

function main() {

	//examples....

	//$("body").append("Hello!!");

	//var test = $('<p>')
 	//	.addClass('myclass')
	//	.html('Ballot page')
	//	.appendTo('body');
	
	//todo...

	//mock up alternative voting first. (add description toggle to alt/crit vote)
	//update decision create --done
	//input test data --done

	//check cookie, split cookie and get values into vars --done
	//get user name, decision, alternative, and criterion data
	//check stage decision is in. if 1: voting not started
	//				 2: can vote
	//				 3: voting closed.
	//
	//build skeleton if stage is 2, else error
	//after submission, display thanks?				 

	var ballot_info = getBallotCookies();
	var decision_id;
	var ballot_id;


	if(ballot_info == null) {
	
		$("body").append("<h1>Error: No ballot found.  Please make sure cookies are enabled and try the link again</h1>");	

	}
	else {
		decision_id = ballot_info[0];
		ballot_id = ballot_info[1];

		$("body").append("Decision_id: "+decision_id+"<br>"+"Ballot_id: "+ballot_id);

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

	} 
	else if(cookie_1[0] == "decision_id") {
		ballot_id = cookie_2[1];
		decision_id = cookie_1[1];
	} 

	return ret = [decision_id, ballot_id];

}


//From test development page, for reference..

/*	
var my_cookies = document.cookie.split(';');

var bal_id = my_cookies[0].slice(-1);
var desc_id = my_cookies[1].slice(-1);

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


document.getElementById("get_bal").onclick = function () {

	$('#ballot_info').append("<h3>Ballot Data!</h3>");


	get_text("/decision/"+desc_id+"/info", function(x) {


		$('#ballot_info').append("<p id='dec_data'></p>");

		$('#dec_data').append("<br>dec_name "+x.decision.name);
		$('#dec_data').append("<br>dec_desc "+x.decision.description);
		$('#dec_data').append("<br>dec_stage "+x.decision.stage);
		$('#dec_data').append("<br>dec_vote_style "+x.decision.criterion_vote_style);

	});



	get_text("/decision/"+desc_id+"/criterions", function(result) {

		$('#ballot_info').append("<p id='crit_data'></p>");


		for(var i=0; i<result.criterions.length; i++) {
			$("#crit_data").append("<br>criterion"+i+" "+result.criterions[i].name);

		}

	});

	get_text("/decision/"+desc_id+"/alternatives", function(result) {

		$('#ballot_info').append("<p id='alt_data'></p>");

		for(var i=0; i<result.alternatives.length; i++) {
			$("#crit_data").append("<br>alternative"+i+" "+result.alternatives[i].name);

		}

	});

*/
//	Get the persons name, if available.

//	get_text("/decision/"+document.getElementById("dec_id").value+"/ballot/"+document.getElementById("bal_id").value+"/info", function(r) {
//	
//		console.log(r);	
//
//	});

//}


//		$('.criterion').click(function(event) {
//			var id = "#" + this.id;
//
//			$(id+"Desc").toggle();
//
//			$('.alert').not(id +"Desc").hide();

		//$('.color1').click(function(event) {
		//	var id = "#"+this.id.slice(0,-1)+"_color";
		//	$(id).css('background-color','red');
		//});

		//$('.color2').click(function(event) {
		//	var id = "#"+this.id.slice(0,-1)+"_color";
		//	$(id).css('background-color','yellow');
		//});

		//$('.color3').click(function(event) {
		//	var id = "#"+this.id.slice(0,-1)+"_color";
		//	$(id).css('background-color','green');
		//});



//		});



