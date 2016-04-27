$(document).ready(function () {
  $("#sw_status").append("Waiting for system time.."); 
  var $start = "";
  $("#sw_start").click(function(){
    $start = setInterval("delayedPost()", 1000);
  });
  
  $("#sw_pause").click(function(){
    clearInterval($start);
  });
});
function delayedPost() {
 $.post("http://localhost:8080/gettime", "", function(data, status) {
   $("#timeread").empty();
   $("#timeread").append(data);
//   $("#timeread").load(data);
  //  var hello = "hello";
  //  $('#messages').append("<p>"+hello+"</p>"); 
 });
}