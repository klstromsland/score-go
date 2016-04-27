$(document).ready(function() {		
//  var timeout_m = {{ .Maxtime_m }};
//  var timeout_s = {{ .Maxtime_s }};
//  var timeout_ms = {{ .Maxtime_ms }};
//  var timeout = (timeout_m * 60000) + (timeout_s * 1000) + (timeout_ms * 10);
  var hello = {{ $maxtime_m }};
  $('#messages').append("<p>"+"hello"+"</p>");
});