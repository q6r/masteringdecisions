function main(body) {
  //grab the decisionID from the url
  var decisionID = window.location.pathname.substr('/decision/'.length);
  
  //Add CSS
  $.loadCSS('/static/css/decision.css');

  //Build page
  $('title').html('Decision Voting Table');
  $('<h2 id="title"></h2>').appendTo('body');
  
  get_text("/decision/"+decisionID+"/info", function (result) {
    if(result['decision']) {
      if(result['decision']['stage'] != 3) {
        $('#title').text('Voting Results For ' + result['decision']['name'] + ' are not yet avalible');
      }
      else { //Voting is completed, show results!
        buildResultsPage(decisionID);
        $('#title').text('Voting Results For: ' + result['decision']['name']);
      }
    }
    else if(result['error']) {
      $('#title').text(result['error']);
    }
    else {
      $('#title').text('Could not find it :()');
    }
  });
}

function buildResultsPage(decisionID) {
  get_text("/decision/"+decisionID+"/criterions", function (results) {
    var critNames = []; //Array of crit names
    var critIds = [];
    for(var i in results["criterions"]) {
      var c = results["criterions"][i];
      critNames.push(c["name"]);
      critIds.push(c["criterion_id"]);
    }
    
    get_text("/decision/"+decisionID+"/alternatives", function (results) {
      
      var altNames = []; //Array of alt names
      var altIds = [];
      for(i in results["alternatives"]) {
        var a = results["alternatives"][i];
        altNames.push(String.fromCharCode(65 + parseInt(i)));
        altIds.push(a["alternative_id"]);
      }
      
      //Table and row 1
      var table = $('<table>').addClass('table table-striped').append('<tbody>').appendTo('body');
      table.append(
        $('<tr>').append(
          $('<th>'), //blank as it is the names
          $('<th>').text('Step 1').attr('colspan', critNames.length).addClass('step1'),
          $('<th>').text('Step 2').attr('colspan', critNames.length * altNames.length).addClass('step2'),
          $('<th>').text('Step 4').attr('colspan', altNames.length).addClass('step4')
        )
      );
      
      //Row 2 and all values
      var row2 = $('<tr>').appendTo(table);
      row2.append($('<td>')); //name
      for(j=0; j<critNames.length; j++) {
        row2.append($('<td>'));
      }
      for(j=0; j<critNames.length; j++) {
        row2.append($('<td>').text(critNames[j]).attr('colspan', altNames.length).addClass('critHeaders crit'+j));
      }
      
      //Row 3 and all values
      var row3 = $('<tr>').addClass('row3').appendTo(table);
      row3.append($('<td>').text('Name'));
      for(j=0; j<critNames.length; j++) {
        row3.append($('<td>').text(critNames[j]));
      }
      for(j=0; j<critNames.length+1; j++) { //+1 to show the totals at the end!
        for(k=0; k<altNames.length; k++) {
          row3.append($('<td>').text(altNames[k]));
        }
      }
      
      get_text("/decision/"+decisionID+"/ballots", function (results) {
        for(var i in results["ballots"]) {
          //Build the template <tr>
          var tr = $('<tr>').appendTo(table);
          tr.append($('<td>').attr('id', 'name'+i).text('-'));
          for(j=0; j<critIds.length; j++) {
            tr.append($('<td>').attr('id', 'crit'+i+'-'+critIds[j]).text('-'));
          }
          for(j=0; j<critIds.length; j++) {
            for(k=0; k<altIds.length; k++) {
              tr.append($('<td>').attr('id', 'alt'+i+'-'+critIds[j]+'-'+altIds[k]).text('-'));
            }
          }
          
          var b = results["ballots"][i];
          $('#name'+i).text(b["name"]);
          
          //set crit values in template
          for(var j in b["rating"]) {
            var r = b["rating"][j];
            $('#crit'+i+'-'+parseInt(r["criterion_id"])).text(r["rating"]);
          }
          
          var weights = [];
          var weightTotal = 0;
          
          for(var j in altIds) {
            tr.append($('<td>').attr('id', 'weights'+i+'-'+altIds[j]).text('-'));
            
            weights.push(altIds[j]);
            weights[altIds[j]] = 0;
          }
          
          var colors = ['-', 'R', 'O', 'Y', 'GY', 'G'];          
          for(var j in b["votes"]) {
            var v = b["votes"][j];
            $('#alt'+i+'-'+parseInt(v["criterion_id"])+'-'+parseInt(v["alternative_id"])).text(colors[parseInt(v["weight"])]).addClass('alt-'+colors[parseInt(v["weight"])]);
            
            weights[parseInt(v["alternative_id"])] += parseInt(v["weight"]);
            weightTotal += parseInt(v["weight"]);
          }
          
          for(var j in altIds) {
            if(weights[altIds[j]]) {
              $('#weights'+i+'-'+altIds[j]).text((weights[altIds[j]]*100/weightTotal).toFixed(2));
            }
          }
        }
      });
    });
  });
}