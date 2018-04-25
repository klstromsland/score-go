$(document).ready(function () {
  $("#sw_status").empty();
  $("#sw_status").append("Ready");
  var $start = "";
  $("#sw_start").click(function(){
    $("#sw_status").empty();
    $("#sw_status").append("Started");
    $start = setInterval("delayedPost()", 10);
  });
  $("#sw_pause").click(function(){
    $("#sw_status").empty();
    $("#sw_status").append("Stopped");
    clearInterval($start);
  });
//  $("#sw_reset").click(function(){
//    $("#sw_status").empty();
//    $("#sw_status").append("Reset");
//    clearInterval($start);
//  });
});
function delayedPost() {
$.post("https:/intense-gorge-68089.herokuapp.com/gettime", "", function(data, status) {
   $("#timeread").empty();
   $("#timeread").append(data);
 });
}
