
// var AudioContext = require('web-audio-api').AudioContext;
var fs = require('fs');
var path = require('path');

var possibleThetas =  [-80, -65, -55, -45, -40, -35, -30, -25, -20, -15, -10, -5, 0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 55, 65, 80];

var hrir_l;
var hrir_r;
var ctx;
var gain;
var hrir_length;

var helperFunctions = {

  spherical_to_hrir: function(s)
  {

      console.log(s);
      var x = Math.sin(s.alpha) * Math.cos(s.beta);
      var rx = (Math.sqrt(1 - Math.pow(x/1,2.0)));
      var gamma = Math.acos(Math.cos(s.alpha) / rx);
      gamma = Math.sign(s.beta) * gamma;

      //Boundaries
      if(gamma < -Math.PI * 3/4.0) gamma = -Math.PI * 3/4.0;
      if(gamma > Math.PI * (3/4.0+0.03125)) gamma = Math.PI * (3/4.0+0.03125);

      return {
          x: x,
          rx: rx,
          gamma: gamma
      };
  },
  hrir_index: function(hrir)
  {
      function closest(arr, target) {
          for (var i=1; i < arr.length; i++) {
              if (arr[i] > target) {
                  var p = arr[i-1];
                  var c = arr[i];
                  return Math.abs(p-target) < Math.abs(c-target) ? i-1 : i;
              }
          }
          return arr.length-1;
      }

      // convert degrees to radians

      var thetaRads = [];
      var i = 0;
      for( ; i < possibleThetas.length; i++)
      {
          thetaRads.push(possibleThetas[i]/180.0*Math.PI);
      }

      return {
          theta: closest(thetaRads, Math.asin(hrir.x)),
          phi: Math.round(hrir.gamma/((6/4.0+0.03125)*Math.PI)*49)+24
      };
  },
  load_hrir_and_update_convolver: function(name) {
      hrir_loaded = false;
      hrir_l = [];
      hrir_r = [];
      for(var i = 0; i < 25; i++)
      {
          hrir_l.push([]);
          hrir_r.push([]);
          for(var j = 0; j < 50; j++)
          {
              hrir_l[i].push([]);
              hrir_r[i].push([]);
              for(var k = 0; k < 200; k++)
              {
                  hrir_l[i][j].push(0);
                  hrir_r[i][j].push(0);
              }
          }
      }

      // load right and left
      var client_l = {};
      var client_r = {};
      fs.readFile(path.join(__dirname,'./hrir/3_l.dat'), 'utf8', function (err, leftText) {

        console.log(err);
        client_l.responseText = leftText;

        fs.readFile(path.join(__dirname, './hrir/3_r.dat'), 'utf8', function(err, rightText) {
          client_r.responseText = rightText;


          // left
          var all_data = client_l.responseText.split('\n');
          for( i = 0; i < 25; i++)
          {
              var data = all_data[i].split(',');
              for( var c = 0; c < 50; c++)
              {
                  for(var q = 0; q < 200; q++)
                  {
                      hrir_l[i][c][q] = data[c+q*50];
                  }
              }
          }

          // right
          all_data = client_r.responseText.split('\n');
          for( i = 0; i < 25; i++)
          {
              data2 = all_data[i].split(',');
              for(var l = 0; l < 50; l++)
              {
                  for(var m = 0; m < 200; m++)
                  {
                      hrir_r[i][l][m] = data2[l+m*50];
                  }
              }
          }

          // hpl.update_convolver();

          ctx = new AudioContext();
          gain = ctx.createGain();

          hrir_length = Math.ceil(200.0 * ctx.sampleRate / 44100);



        });


      });
  },





};



module.exports = helperFunctions;
