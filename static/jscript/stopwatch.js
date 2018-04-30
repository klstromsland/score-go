var start = 0;
var running = 0;
var timetxt = ""

function startPause(){
     if (running == 0){
          running = 1;
          increment();
          document.getElementById("sw_start").innerHTML = "Pause";
          start = Date.now();
     }else{
          running = 0;
          document.getElementById("sw_start").innerHTML = "Resume";
     }
};

function resetWatch(){
     running = 0;
     time = 0;
     document.getElementById("timeread").innerHTML = "00:00:00";
     document.getElementById("sw_start").innerHTML = "Start";
};

function increment(){
  if(running == 1){
    setTimeout(function(){
      var millis = Date.now() - start;
      if (millis < 1000){  // no seconds, no minutes
        var secs = 0;
        var mins = 0;
        if (millis > 99){
          hundreths = Math.floor(millis/10);
        }else{
          hundreths = millis;
        }
      }else{  // we have seconds
        var secs = Math.floor(millis/1000);
        var hundreths = millis%1000;  // milliseconds
        if (secs < 60){  // no minutes
          mins = 0;
          if (hundreths > 99){
            hundreths = Math.floor(hundreths/10);
          }
        }else{  // we have minutes
          var mins = Math.floor(secs/60);
          secs = secs + secs%60;
          if (hundreths > 99){
            hundreths = Math.floor(hundreths/10);
          }
        }
      }
      var minstr = mins.toString();
      if (mins < 10){
        minstr = "0" + minstr;
      }
      var secstr = secs.toString();
      if (secs < 10){
        secstr = "0" + secstr;
      }
      var hundrethstr = hundreths.toString();
      if(hundreths < 10){
        hundrethstr = "0" + hundrethstr;
      }
      timetxt = minstr + ":" + secstr + ":" + hundrethstr;
      document.getElementById("timeread").innerHTML = timetxt;
      increment();
    },100);
  }
};
