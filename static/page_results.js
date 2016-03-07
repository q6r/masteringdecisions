var decision_name;
var criterion_names = [];
var criterion_ids = [];
var alternative_names = [];
var alternative_ids = [];
var ballots;


function main(body) {
  getScript('https://www.gstatic.com/charts/loader.js', function() {
    google.charts.load('current', {
      packages: ['corechart', 'table']
    });
    google.charts.setOnLoadCallback(function() {
      var decision_id = window.location.pathname.substr('/results/'.length);
      var decision_stage;

      $('title').html('Results');
      $.loadCSS('/static/css/results.css');

      var decision = get_decision(decision_id);
      if (decision == null || decision["error"] != null || decision["error"] != undefined) {
        errorPage("Unable to get decision information");
        return;
      }

      ballots = get_ballots(decision_id);
      if (ballots == null || ballots["error"] != null || ballots["error"] != undefined) {
        errorPage("Unable to get ballot information");
        return;
      }

      var criterion = get_criterion(decision);
      if (criterion == null || criterion["error"] != null || criterion["error"] != undefined) {
        errorPage("Unable to get criterion information");
        return;
      }

      var alternatives = get_alternatives(decision);
      if (alternatives == null || alternatives["error"] != null || alternatives["error"] != undefined) {
        errorPage("Unable to get alternative information");
        return;
      }

      decision_stage = decision.stage;
      if (decision_stage != "3") {
        errorPage("Decision stage must be complete in order to view results");
        return;
      } else {
        decision_name = decision.name;

        //criterion info
        for (var i = 0; i < criterion.length; i++) {
          criterion_names[i] = criterion[i].name;
          criterion_ids[i] = criterion[i].criterion_id;
        }

        //alternative info
        for (var i = 0; i < alternatives.length; i++) {
          alternative_names[i] = alternatives[i].name;
          alternative_ids[i] = alternatives[i].alternative_id;
        }

        var page = '<h1 id="chart_title">Voting Results - ' + decision_name + '</h1>' + '<div id="pie_div"></div>' + '<div id="chart_div"></div>' + '<div id="table_div"></div>' + '<input class="btn btn-warning" id="clearbtn" type="button" value="Refresh" onClick="window.location.reload()">' + '<form>' + '<input type="checkbox" name="pause" id="paused"> Pause automatic weight balancing' + '</form>';

        $('<div id="content">').appendTo('body');
        $('#content').append('<div id="header" class="navbar-inverse"><CENTER><IMG SRC="/static/images/logo.png" ALIGN="TOP" /></CENTER></div>');
        $('#content').append(page);

        var data_table = $('<table>').addClass('table table-striped').attr('id', 'dataTable').append('<tbody>').appendTo('#content');
        var temp_row = $('<tr>');

        $(temp_row).append($('<th>').text(' '));

        for (var i = 0; i < criterion_names.length; i++) {
          $(temp_row).append($('<th>').text(criterion_names[i]).addClass("crit" + i));
        }

        $(temp_row).append($('<th>').text('Final Tally'));
        $(temp_row).appendTo(data_table);

        temp_row = $('<tr>');

        $(temp_row).append($('<td>').text('Weights'));

        for (var i = 0; i < criterion_names.length; i++) {
          $(temp_row).append($('<td>').append($('<input type="text">').attr('id', 'weight_' + i).val('-').attr('onChange', 'rebalance_weights(' + i + ')')));
        }

        $(temp_row).append($('<td>').attr('id', 'weight_total').text('-'));
        $(temp_row).appendTo(data_table);

        for (var i = 0; i < alternative_names.length; i++) {
          temp_row = $('<tr>');
          $(temp_row).append($('<td>').attr('id', 'alt_names_' + i).text(alternative_names[i]));

          for (var j = 0; j < criterion_names.length; j++) {
            $(temp_row).append($('<td>').append($('<input type="text">').attr('id', 'value_' + i + '_' + j).val('-').attr('onChange', 'recalculate_tally()')));
          }

          $(temp_row).append($('<td>').attr('id', 'totals_' + i).text('-'));
          $(temp_row).appendTo(data_table);
        }

        calculate_table();
        draw_graphs();
      }
    });
  });
}

function recalculate_tally() {
  var new_tally;
  for (var i = 0; i < alternative_names.length; i++) {
    new_tally = 0;
    for (var j = 0; j < criterion_names.length; j++) {
      new_tally += parseFloat($('#value_' + i + '_' + j).val() * $('#weight_' + j).val() / 100);
    }
    $('#totals_' + i).text(new_tally.toFixed(3));
  }
  draw_graphs();
}

function draw_graphs() {
  var data_c = new google.visualization.DataTable();
  var colors = ["#3366cc", "#dc3912", "#ff9900", "#109618", "#990099", "#0099c6", "#dd4477", "#66aa00", "#b82e2e", "#316395", "#994499", "#22aa99", "#aaaa11", "#6633cc", "#e67300", "#8b0707", "#651067", "#329262", "#5574a6", "#3b3eac", "#b77322", "#16d620", "#b91383", "#f4359e", "#9c5935", "#a9c413", "#2a778d", "#668d1c", "#bea413", "#0c5922", "#743411"]
  var chartColors = new Array();
  data_c.addColumn('string', 'criterion');

  for (var i = criterion_names.length-1; i >= 0; i--) {
    data_c.addColumn('number', criterion_names[i]);
    chartColors.push(colors[i]);
  }

  chart_rows = [];
  for (var i = alternative_names.length-1; i >= 0; i--) {
    var row = [alternative_names[i]];
    chart_rows[i] = new Array();
    for (var j = criterion_names.length-1; j >= 0; j--) {
      row.push(parseFloat($('#value_' + i + '_' + j).val()) * $('#weight_' + j).val() / 100);
    }
    data_c.addRow(row);
  }

  var col_chart = new google.visualization.ColumnChart(document.getElementById('chart_div'));

  col_chart.draw(data_c, {
    width: 600,
    height: 400,
    vAxis: {
      title: "Votes"
    },
    isStacked: true,
    colors: chartColors, //reverse of Googles colors
    chartArea: {  width: "50%" }
  });

  var data_p = new google.visualization.DataTable();

  var pie_row;
  data_p.addColumn('string', 'criterion');
  data_p.addColumn('number', 'weight');

  for (var i = 0; i < criterion_names.length; i++) {
    pie_row = [];
    pie_row.push(criterion_names[i]);
    pie_row.push(parseFloat($('#weight_' + i).val()));
    data_p.addRow(pie_row);
  }

  var pie = new google.visualization.PieChart(document.getElementById('pie_div'));
  pie.draw(data_p, {
    title: "Criterion weights",
    width: 600,
    height: 400
  });
}

function rebalance_weights(x) {


  if ($("#paused").prop("checked") == false) {
    //Check overflow
    if (parseFloat($('#weight_' + x).val()) > 100) {
      $('#weight_' + x).val('100')
    }
    if (parseFloat($('#weight_' + x).val()) < 0) {
      $('#weight_' + x).val('0')
    }

    var total_weight = 0;
    for (var i = 0; i < criterion_names.length; i++) {
      total_weight += parseFloat($('#weight_' + i).val());
    }


    for (var i = 0; i < criterion_names.length; i++) {
      if (i != x) {
        var old_value = parseFloat($('#weight_' + i).val());
        var changed_value = parseFloat($('#weight_' + x).val());
        var sub_total = parseFloat(total_weight - changed_value);

        var new_value = (old_value - (old_value / sub_total * (total_weight - 100))).toFixed(3);
        var new_tally;

        $('#weight_' + i).val(new_value);
      }
    }
  }
  for (var i = 0; i < alternative_names.length; i++) {
    new_tally = 0;
    for (var j = 0; j < criterion_names.length; j++) {
      new_tally += parseFloat($('#value_' + i + '_' + j).val() * $('#weight_' + j).val() / 100);
    }
    $('#totals_' + i).text(new_tally.toFixed(3));
  }


  var new_total = 0;
  for (var i = 0; i < criterion_names.length; i++) {

    new_total += parseFloat($('#weight_' + i).val());

  }

  $("#weight_total").text(Math.round(new_total));


  draw_graphs();

}

function calculate_table() {
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

  for (var i = 0; i < ballots.length; i++) {
    temp_sum = 0;
    temp_weights = [];

    if (ballots[i].rating.length != 0) {
      for (var j = 0; j < ballots[i].rating.length; j++) {
        temp_weights.push(ballots[i].rating[j].rating);
        temp_sum += ballots[i].rating[j].rating;
      }

      weights[weights.length] = new Array();
      for (var j = 0; j < temp_weights.length; j++) {
        weights[weights.length - 1].push(temp_weights[j] / temp_sum);
      }
    }
  }

  if (weights.length == 0) {
    return;
  } //No vote checking, might expand later

  for (var i = 0; i < weights[0].length; i++) {
    temp_sum = 0;
    temp_weights = [];
    for (var j = 0; j < weights.length; j++) {
      temp_sum += weights[j][i];
    }
    final_weights.push(temp_sum / weights.length);
  }

  for (var i = 0; i < final_weights.length; i++) {
    weight_total += final_weights[i];
  }

  for (var i = 0; i < ballots.length; i++) {
    temp_sum = 0;
    temp_votes = [];

    if (ballots[i].votes.length != 0) {
      for (var j = 0; j < ballots[i].votes.length; j++) {
        temp_votes.push(ballots[i].votes[j].weight);
        temp_sum += ballots[i].votes[j].weight;
      }

      votes[votes.length] = new Array();
      for (var j = 0; j < temp_votes.length; j++) {
        votes[votes.length - 1].push(temp_votes[j]);
      }
    }
  }

  for (var i = 0; i < votes[0].length; i++) {
    temp_sum = 0;
    temp_votes = [];
    for (var j = 0; j < votes.length; j++) {
      temp_sum += votes[j][i];
    }
    final_votes.push(temp_sum / votes.length);
  }

  var counter = 0;
  for (var i = 0; i < alternative_names.length; i++) {
    tally[i] = new Array();
    for (var j = 0; j < criterion_names.length; j++) {
      tally[i].push(final_votes[counter]);
      counter++;
    }
  }

  var chart_junk = [];
  for (var i = 0; i < tally.length; i++) {
    temp_sum = 0;
    for (var j = 0; j < tally[0].length; j++) {
      chart_junk.push(tally[i][j] * final_weights[j]);
      temp_sum += (tally[i][j] * final_weights[j]);
    }
    final_tally.push(temp_sum);
  }

  for (var i = 0; i < criterion_names.length; i++) {
    $('#weight_' + i).val((final_weights[i] * 100).toFixed(3));
  }

  $('#weight_total').text(Math.round(weight_total * 100));

  counter = 0;
  var chart_rows = [];
  for (var i = 0; i < alternative_names.length; i++) {
    for (var j = 0; j < criterion_names.length; j++) {
      $('#value_' + i + '_' + j).val(final_votes[counter].toFixed(3));
      counter++;
    }
    $('#totals_' + i).text(final_tally[i].toFixed(3));
  }
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

//retrieve all ballots
function get_ballots(decision_id) {
  var result = null;
  $.ajax({
    type: "GET",
    url: base_url + "/decision/" + decision_id + "/ballots",
    contentType: 'application/json; charset=utf-8',
    async: false,
    success: function(r) {
      result = r["ballots"];
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