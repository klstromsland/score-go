/*

JQUERY: STOPWATCH & COUNTDOWN

This is a basic stopwatch & countdown plugin to run with jquery. Start timer, pause it, stop it or reset it. Same behaviour with the countdown besides you need to input the countdown value in seconds first. At the end of the countdown a callback function is invoked.

Any questions, suggestions? marc.fuehnen(at)gmail.com

*/

$(document).ready(function() {

    (function($){
    
        $.extend({
            
            APP : {                
                
                formatTimer : function(a) {
                    if (a < 10) {
                        a = '0' + a;
                    }                              
                    return a;
                },    
                
                startTimer : function(dir) {
                    
                    var a;
                    
                    // save type
                    $.APP.dir = dir;
                    
                    // get current date
                    $.APP.d1 = new Date();
                    
                    switch($.APP.state) {
                            
                        case 'pause' :
                            
                            // resume timer
                            // get current timestamp (for calculations) and
                            // substract time difference between pause and now
                            $.APP.t1 = $.APP.d1.getTime() - $.APP.td;                            
                            
                        break;
                            
                        default :
                            
                            // get current timestamp (for calculations)
                            $.APP.t1 = $.APP.d1.getTime(); 
                            
                            // if countdown add ms based on seconds in textfield
                            if ($.APP.dir === 'cd') {
                                $.APP.t1 += parseInt($('#cd_seconds').val())*1000;
                            }    
                        
                        break;
                            
                    }                                   
                    
                    // reset state
                    $.APP.state = 'alive';   
                    $('#' + $.APP.dir + '_status').html('Running');
                    
                    // start loop
                    $.APP.loopTimer();
                    
                },
                
                pauseTimer : function() {
				
                    var elapsed_m;
                    var elapsed_s;
                    var elapsed_ms;					
				    
                    // save timestamp of pause
                    $.APP.dp = new Date();
                    $.APP.tp = $.APP.dp.getTime();
                    
                    // save elapsed time (until pause)
                    $.APP.td = $.APP.tp - $.APP.t1;
					
                    td = $.APP.td;
                    
                    // calculate milliseconds
                    elapsed_ms = td%1000;
                    if (elapsed_ms < 1) {
                      elapsed_ms = 0;
                    } else {    
                      // calculate seconds
                      elapsed_s = (td-elapsed_ms)/1000;
                      if (elapsed_s < 1) {
                        elapsed_s = 0;
                      } else {
                        // calculate minutes   
                        elapsed_m = (elapsed_s-(elapsed_s%60))/60;
                        if (elapsed_m < 1) {
                          elapsed_m = 0;
                        }
                      }
                    }
                    
                    // substract elapsed minutes
                    elapsed_ms = Math.round(elapsed_ms/100);
                    elapsed_s  = elapsed_s-(elapsed_m*60);  
                    
                    $("#elapsed_time_m_field").val($.APP.formatTimer(elapsed_m));		            
                    $("#elapsed_time_s_field").val($.APP.formatTimer(elapsed_s));
                    $("#elapsed_time_ms_field").val($.APP.formatTimer(elapsed_ms));
                    $("#timed_out_field").val('no');			
					
                    // change button value
                    $('#' + $.APP.dir + '_start').val('Resume');
                    
                    // set state
                    $.APP.state = 'pause';
                    $('#' + $.APP.dir + '_status').html('Pause');

                },
                
                stopTimer : function() {
                    
                    // change button value
                    $('#' + $.APP.dir + '_start').val('Restart');                    
                    
                    // set state
                    $.APP.state = 'stop';
                    $('#' + $.APP.dir + '_status').html('Stopped');
                    
                },
                
                resetTimer : function() {

                    // reset display
                    $('#' + $.APP.dir + '_ms,#' + $.APP.dir + '_s,#' + $.APP.dir + '_m,#' + $.APP.dir + '_h').html('00');                 
                    
                    // change button value
                    $('#' + $.APP.dir + '_start').val('Start');                    
                    
                    // set state
                    $.APP.state = 'reset';  
                    $('#' + $.APP.dir + '_status').html('Reset & Idle');
                    
                },
                
                endTimer : function(callback) {
                   
                    // change button value
                    $('#' + $.APP.dir + '_start').val('Restart');
                    
                    // set state
                    $.APP.state = 'end';
                    
                    // invoke callback
                    if (typeof callback === 'function') {
                        callback();
                    }    
                    
                },
          
                loopTimer : function() {
                    
                    var td;
                    var d2,t2;
					
                    var ms = 0;
                    var s = 0;
                    var m = 0;
					
                    if ($.APP.state === 'alive') {
                                
                        // get current date and convert it into 
                        // timestamp for calculations
                        d2 = new Date();
                        t2 = d2.getTime();   
                        
                        // calculate time difference between
                        // initial and current timestamp
                        if ($.APP.dir === 'sw') { 
                            td = t2 - $.APP.t1;
                            //if timed out
                            if (td >= timeout) {
                                $.APP.endTimer(function(){
                                    $('#' + $.APP.dir + '_status').html('Timed out');
                                    $("#elapsed_time_m_field").val($.APP.formatTimer(timeout_m));		            
                                    $("#elapsed_time_s_field").val($.APP.formatTimer(timeout_s));
                                    $("#elapsed_time_ms_field").val($.APP.formatTimer(timeout_ms));
                                    $("#timed_out_field").val('yes');
                                });								
                            }
							
                        // reversed if countdown
                        } else {
                            td = $.APP.t1 - t2;
                            if (td <= 0) {
                                // if time difference is 0 end countdown
                                $.APP.endTimer(function(){
                                    $.APP.resetTimer();
                                    $('#' + $.APP.dir + '_status').html('Ended & Reset');
                                });
                            }    
                        }    
                        
                        // calculate milliseconds
                        ms = td%1000;
                        if (ms < 1) {
                           ms = 0;
                        } else {    
                            // calculate seconds
                            s = (td-ms)/1000;
                            if (s < 1) {
                                s = 0;
                            } else {
                                // calculate minutes   
                                var m = (s-(s%60))/60;
                                if (m < 1) {
                                    m = 0;
                                }
                            }
                        }
                      
                        // substract elapsed minutes
                        ms = Math.round(ms/100);
                        s  = s-(m*60);    
                        
                        // update display					
                        $('#' + $.APP.dir + '_ms').html($.APP.formatTimer(ms));
                        $('#' + $.APP.dir + '_s').html($.APP.formatTimer(s));
                        $('#' + $.APP.dir + '_m').html($.APP.formatTimer(m));
						
                        // loop
                        $.APP.t = setTimeout($.APP.loopTimer,1);
                    
                    } else {
                    
                        // kill loop
                        clearTimeout($.APP.t);
                        return true;
                    
                    }  
                    
                }
                    
            }    
        
        });
          
        $('#sw_start').live('click', function() {
            $.APP.startTimer('sw');
        });    

        $('#cd_start').live('click', function() {
            $.APP.startTimer('cd');
        });           
        
        $('#sw_stop,#cd_stop').live('click', function() {
            $.APP.stopTimer();
        });
        
        $('#sw_reset,#cd_reset').live('click', function() {
            $.APP.resetTimer();
        });  
        
        $('#sw_pause,#cd_pause').live('click', function() {
			$.APP.pauseTimer();			
        });                
                
    })(jQuery);        
});