var decision_name;
var criterion_names = [];
var criterion_ids = [];
var alternative_names = [];
var alternative_ids = [];
var ballots;


function main(body) {
  getScript('https://www.gstatic.com/charts/loader.js', function() {
    google.charts.load('current', {packages: ['corechart', 'table']});
    google.charts.setOnLoadCallback(function() {
      var decision_id = window.location.pathname.substr('/results/'.length);
      var decision_stage;

      $('title').html('Results');
      $('head').prepend(custom_header());

      var decision = get_decision(decision_id);

      if(decision == null || decision["error"] != null || decision["error"] != undefined) {
        alert("Unable to get decision information");
        return;
      }
        
      ballots = get_ballots(decision_id);

      if(ballots == null || ballots["error"] != null || ballots["error"] != undefined) {
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

      decision_stage = decision.stage;

      if(decision_stage != "3") {

        $("body").append("<h1>Decision stage must be complete in order to view results</h1>");
        return;

      } else {
        
        decision_name = decision.name;

        //criterion info
        for(var i=0; i<criterion.length; i++) {
          criterion_names[i] = criterion[i].name;
          criterion_ids[i] = criterion[i].criterion_id;		

        }

        //alternative info
        for(var i=0; i<alternatives.length; i++) {
      
          alternative_names[i] = alternatives[i].name;
          alternative_ids[i] = alternatives[i].alternative_id;
      
        }

        var page = '<table align="center">'
        +'<tr align="top">'
        +'<td colSpan=2>'
        +'<div id="chart_div"></div>'
        +'</td>'
        +'</tr>'
        +'<tr>'
        +'<td>'
        +'<div id="table_div"></div>'
        +'</td>'
        +'</tr>'
        +'</table>';

        $('body').append(page);
      
        drawTable();

      }
    });
  });
}


function custom_header() {
	var content='<link rel="stylesheet" type="text/css" href="/static/css/results.css">'
	return content
}

//retrieve decision 
function get_decision(decision_id) {
	var result = null;
	$.ajax({
		type: "GET",
		url: base_url+"/decision/"+decision_id+"/info",
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

function drawTable() {

        var data = new google.visualization.DataTable();
	

	var weight_row = [];
	var temp_sum = 0.0;

	var weights = [];
	var final_weights = [];
	var weight_total = 0;
	var ballot_weights = [];
	var temp_weights = [];

	var votes = [];
	var temp_votes = [];
	var ballot_votes = [];
	var final_votes = [];
	var vote_total = [];

	var tally = [];
	var final_tally = [];


	for(var i=0; i<ballots.length; i++) {

		temp_sum = 0;
		temp_weights = [];

		for(var j=0; j<ballots[i].rating.length; j++) {

			temp_weights.push(ballots[i].rating[j].rating);
			temp_sum += ballots[i].rating[j].rating;


		}
		
		weights[i] = new Array();

		for(var j=0; j<temp_weights.length; j++) {

			weights[i].push(temp_weights[j]/temp_sum);

		}
		

	}


	for(var i=0; i<weights[0].length; i++) {

		temp_sum = 0;
		temp_weights = [];

		for(var j=0; j<weights.length; j++) {

			temp_sum += weights[j][i];
		}
		
		final_weights.push(temp_sum/weights.length);

	}	
	
	for(var i=0; i<final_weights.length; i++) {
		weight_total += final_weights[i];
	}


	for(var i=0; i<ballots.length; i++) {

		temp_sum = 0;
		temp_votes = [];	


		for(var j=0; j<ballots[i].votes.length; j++) {

			temp_votes.push(ballots[i].votes[j].weight);
			temp_sum += ballots[i].votes[j].weight;

		}


		votes[i] = new Array();
	
		for(var j=0; j<temp_votes.length; j++) {
			
			votes[i].push(temp_votes[j]);	

		}

	}


	for (var i=0; i<votes[0].length; i++) {

		temp_sum=0;
		temp_votes = [];


		for(var j=0; j<votes.length; j++) {

			temp_sum += votes[j][i];

		}

		final_votes.push(temp_sum/votes.length);

	}

	var counter=0;

	for(var i=0; i<alternative_names.length; i++) {

		tally[i] = new Array();

		for (var j=0; j<criterion_names.length; j++) {

			tally[i].push(final_votes[counter]);
			counter++;
		}

	}

	var chart_junk = [];

	for(var i=0; i<tally.length; i++) {

		temp_sum = 0;

		for(var j=0; j<tally[0].length; j++) {
			chart_junk.push(tally[i][j] * final_weights[j]);
			temp_sum += (tally[i][j] * final_weights[j]);
		}

		final_tally.push(temp_sum);

	}

        data.addColumn('string', '');
	
	for(var i=0; i<criterion_names.length; i++) {

		data.addColumn('number', criterion_names[i]);

	}

        data.addColumn('number', 'Final Tally');

	weight_row = ['weights'];

	for(var i=0; i<criterion_names.length; i++) {

		weight_row.push(final_weights[i] * 100);

	}
	
	weight_row.push(weight_total * 100);

	data.addRow(weight_row);

	counter = 0;
	var chart_rows = [];

	for(var i=0; i<alternative_names.length; i++) {
		
		var row = [alternative_names[i]];
		chart_rows[i] = new Array();

		for(var j=0; j<criterion_names.length; j++) {
		
			row.push(final_votes[counter]);
			counter++;
		}

		row.push(final_tally[i]);
		data.addRow(row);

	}

        var table = new google.visualization.Table(document.getElementById('table_div'));

        table.draw(data, {showRowNumber: false, width: '100%', height: '100%'});



        var data_c = new google.visualization.DataTable();

	data_c.addColumn('string', 'criterion');

	for(var i=0; i<criterion_names.length; i++) {

		data_c.addColumn('number', criterion_names[i]);

	}

	counter = 0;

	chart_rows = [];

	for(var i=0; i<alternative_names.length; i++) {

		var row = [alternative_names[i]];
		chart_rows[i] = new Array();

		for(var j=0; j<criterion_names.length; j++) {

			row.push(chart_junk[counter]);
			counter++;

		}
				
		data_c.addRow(row);	
	}	

	new google.visualization.ColumnChart(document.getElementById('chart_div')).
		draw(data_c,
		{title:"Voting Results - "+decision_name,
		width:600, height:400,
		vAxis: {title: "Votes"}, isStacked: true}
	);

}


//retrieve all ballots
function get_ballots(decision_id) {
	var result = null;
	$.ajax({
		type: "GET",
		url: base_url+"/decision/"+decision_id+"/ballots",
		contentType: 'application/json; charset=utf-8',
		async: false,
		success: function (r) {
			result = r["ballots"];
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
		url: base_url+"/decision/"+decision["decision_id"]+"/criterions",
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
		url: base_url+"/decision/"+decision["decision_id"]+"/alternatives",
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

