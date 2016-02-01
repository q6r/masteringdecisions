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

function create_person(person, cb) {
	post_text("/person", JSON.stringify(person), cb)
}

// dev : this is used before every bunch of functions
// because we maybe running faster than the backend's
// database can commit
function delay() {
    var now = new Date().getTime();
    while(new Date().getTime() < now + 500){ /* do nothing */ } 
}

// Create a facilitator

new_person = {
	"email": "iessa@pdx.edu",
	"pw_hash":"asd",
	"name_first":"hello",
	"name_last":"world"}

//////////////////////////////////////////////////
// START : EXAMPLE ONE
// in this callback hell we do the following
// 1. Create a person
// 2. Create a decision owned by that person
// 3. Create criterions for the decision
// 4. Create ballots for the decision
// 5. Mimick how a ballot vote on a criterion
//////////////////////////////////////////////////

create_person(new_person, function (person) {

	assert(person["email"] == "iessa@pdx.edu", "Wrong email");

	// Create a decision for this facilitator
	new_decision = {
		"person_id":person.person_id,
		"name":"x",
		"description":"x",
		"owner_id":23,
		"stage":5,
		"criterion_vote_style":"x",
		"alternative_vote_style":"x",
		"client_settings":"x"}

	post_text("/decision", JSON.stringify(new_decision),
			function (decision) {
				console.log("Decision for " + person.email + " created : ");
				// Now create 10 ballots for this decision
				for(i=0;i<10;i++) {
					new_ballot = {
						"secret": 223344,
						"name": "ballot",
						"email": "useremail@a.com"
					};
					post_text("/decision/"+decision.decision_id+"/ballot",
							JSON.stringify(new_ballot),
							function (ballot) {
								console.log("ballot id #" + ballot.ballot_id + " created");
							});
				}

				// Now show the created ballots for this decision
				get_text("/decision/"+decision.decision_id+"/ballots", function (ballots) {
					for(var i in ballots) {
						ballot = ballots[i];
						console.log(ballot);
					}
				})

				// Now setup criterions for this decision
				crit1 = {"name": "crit1", "weight": 20};
				crit2 = {"name": "crit2", "weight": 20};
				crit3 = {"name": "crit3", "weight": 20};
				crit4 = {"name": "crit4", "weight": 20};
				crit5 = {"name": "crit5", "weight": 20};
				crits = [crit1, crit2, crit3, crit4, crit5];
				for(var i in crits) {
					crit = crits[i];
					post_text("/decision/"+decision.decision_id+"/criterion",
							JSON.stringify(crit),
							function (c) {
								console.log("criterion id #"+c.criterion_id+" created");
							});
				}

				// Now show the created criterions for this decision
				get_text("/decision/"+decision.decision_id+"/criterions", function (crits) {
					// Now show how a ballot vote on a criterion in a decision
					crit = crits[0]; // Taking the first crit
					get_text("/decision/"+decision.decision_id+"/ballots", function (ballots) {
						b = ballots[0]; // Taking the first ballot
						// ballot[0] vote on crit[0] on decision with weight 20
						get_text("/decision/"+decision.decision_id+"/ballot/"+b.ballot_id+"/criterion/"+crit.criterion_id+"/vote/20",
								function (v) {
									console.log("ballot "+b.ballot_id+" voted for crit " + crit.criterion_id);
									// get the votes for ballot[0]
									get_text("/decision/"+decision.decision_id+"/ballot/"+b.ballot_id+"/votes",
											function (vs) {
												console.log(JSON.stringify(vs));
											});
								});
					});
				});



			});

});
/////////////////////
// END EXAMPLE ONE //
/////////////////////
