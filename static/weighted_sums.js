//File for core functions used across multiple scripts



//Takes in a LIST of ratings and list of weights
//Returns LIST of weighted values
// Ballots should be in the form of 
function statistics(ratingList, weightList){
  //Error Check
  try{
    if(ratingList.length != weightList.length || ratingList[0].length != weightList[0].length){
      throw "Arrays of different lengths";
    }
  }
  catch(err){
    console.log("Error: 666");
    return [];
  }

  var g = 0;
  //console.log("2");
  while(ratingList[0][g]){
    ++g;
  }
  var i = 0;
  var values = [];
  var weightedSums = [];
  var h = 0;
  //console.log("3");
  for(h=0; h < g; h++){
    values.push(0);
    weightedSums.push(0);
  }
  //Adding up all the ratings
  //console.log("4");
  while(ratingList[i]){
    var ratings = ratingList[i];
    var k = 0;
    while(ratings[k]){
      values[k] = values[k] + ratings[k];
      ++k;
    }
    ++i;
  }
  var j = 0;
  //Adding up all the Weights
  //console.log("5");
  while(weightList[j]){
    var weights = weightList[j];
    var l = 0;
    while(weights[l]){
      weightedSums[l] = weightedSums[l] + weights[l];
      ++l;
    }
    ++j;
  }
  var x = 0;
  var z = 0;
  //console.log("6");
  while(weightedSums[x]){
    var temp1 = weightedSums[x] / j;
    weightedSums[x] = temp1;
    ++x;
  }
  //console.log("7");
  while(values[z]){
    var temp2 = values[z] * weightedSums[z];
    values[z] = temp2;
    ++z;
  }
  //console.log("8");
  return [values, weightedSums];
}

//Valid test
function Valid_Test(){
  console.log("FOO TEST HAS BEGUN");
  var lists = [[1,2],[2,1],[3,4],[2,5]];
  var weights = [[1,5],[2,3],[5,2],[1,2]];
  //console.log("1");
  var results = statistics(lists, weights);
  //console.log("9");
  console.log(results.toString());
  var bad = statistics([[1]],[[1,2]]);
  console.log(bad.toString());
}
