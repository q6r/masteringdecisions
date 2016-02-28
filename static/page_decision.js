function main(body) {
  //grab the decisionID from the url
  var decisionID = window.location.pathname.substr('/decision/'.length);

  //Add CSS
  $.loadCSS('/static/css/decision.css');

  //Build page
  $('title').html('Decision Voting Table');

  get_text("/decision/" + decisionID + "/info", function(result) {
    if (result['decision']) {
      if (result['decision']['stage'] == 1 || result['decision']['stage'] ==2 ) {
        errorPage('Voting Results For ' + result['decision']['name'] + ' are not yet avalible');
      }
      else if (result['decision']['stage'] == 4) {
        errorPage(result['decision']['name'] + ' has been locked.');
      } else { //Voting is completed, show results!
        buildResultsPage(decisionID);
        $('<h2 id="title"></h2>').appendTo('body');
        $("#title").text('Voting Results For: ' + result['decision']['name']);
      }
    } else if (result['error']) {
      errorPage(result['error']);
    } else {
      errorPage('Could not find it :()');
    }
  });
}

function buildResultsPage(decisionID) {
  get_text("/decision/"+decisionID+"/criterions", function(results){
    var critNames = []; //Array of crit names
    var critIds = [];
    for (var i in results["criterions"]) {
      var c = results["criterions"][i];
      critNames.push(c["name"]);
      critIds.push(c["criterion_id"]);
    }

    get_text("/decision/" + decisionID + "/alternatives", function(results) {
      var altNames = []; //Array of alt names
      var altIds = [];
      
      var key = $('<table>').attr('id', 'keyTable').appendTo('body');
      key.append($('<tr>').append($('<th>').text('Alternatives Key').attr('colspan', '2')));
      
      for (i in results["alternatives"]) {
        var a = results["alternatives"][i];
        altNames.push(String.fromCharCode(65 + parseInt(i)));
        altIds.push(a["alternative_id"]);
        key.append('<tr><td>' + String.fromCharCode(65 + parseInt(i)) + '</td><td>' + a['name'] + '</td></tr>');
      }
      
      //Table and row 1
      var table = $('<table>').addClass('table table-striped').attr('id', 'votesTable').append('<tbody>').appendTo('body');
      table.append(
        $('<tr>').append(
          $('<th>').text(' '), //space to fix copy/paste bug
          $('<th>').text('Step 1').attr('colspan', critNames.length).addClass('step1'),
          $('<th>').text('Step 2').attr('colspan', critNames.length * altNames.length).addClass('step2')
        )
      );

      //Row 2 and all values
      var row2 = $('<tr>').appendTo(table);
      row2.append($('<td>')); //name
      for (j = 0; j < critNames.length; j++) {
        row2.append($('<td>'));
      }
      for (j = 0; j < critNames.length; j++) {
        row2.append($('<td>').text(critNames[j]).attr('colspan', altNames.length).addClass('critHeaders crit' + j));
      }

      //Row 3 and all values
      var row3 = $('<tr>').addClass('row3').appendTo(table);
      row3.append($('<td>').text('Name'));
      for (j = 0; j < critNames.length; j++) {
        row3.append($('<td>').text(critNames[j]));
      }
      for (j = 0; j < critNames.length; j++) {
        for (k = 0; k < altNames.length; k++) {
          row3.append($('<td>').text(altNames[k]));
        }
      }

      //Fill the table
      get_text("/decision/" + decisionID + "/ballots", function(results) {
        for (var i in results["ballots"]) {
          //Build the template <tr>
          var tr = $('<tr>').appendTo(table);
          tr.append($('<td>').attr('id', 'name' + i).text('-'));
          for (j = 0; j < critIds.length; j++) {
            tr.append($('<td>').attr('id', 'crit' + i + '-' + critIds[j]).text('-'));
          }
          for (j = 0; j < critIds.length; j++) {
            for (k = 0; k < altIds.length; k++) {
              tr.append($('<td>').attr('id', 'alt' + i + '-' + critIds[j] + '-' + altIds[k]).text('-'));
            }
          }

          //Fill in the Name
          var b = results["ballots"][i];
          $('#name' + i).text(b["name"]);
          
          //Fill in the Ratings
          for (var j in b["rating"]) {
            var r = b["rating"][j];
            $('#crit' + i + '-' + parseInt(r["criterion_id"])).text(r["rating"]);
          }

          //Fill in the Votes
          var colors = ['-', 'R', 'O', 'Y', 'GY', 'G'];
          for (var j in b["votes"]) {
            var v = b["votes"][j];
            $('#alt' + i + '-' + parseInt(v["criterion_id"]) + '-' + parseInt(v["alternative_id"])).text(colors[parseInt(v["weight"])]).addClass('alt-' + colors[parseInt(v["weight"])]);
          }
        }
      });
    });
  });
}
