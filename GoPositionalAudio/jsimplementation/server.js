// 3D audio server 2017

var express = require('express');
var bodyParser = require('body-parser');
var morgan = require('morgan');
var cors = require('cors');

var app = express();
// var context = new AudioContext();
var Speaker = require('speaker');
// var helpers = require('./helpers');
//
// context.outStream = new Speaker({
//   channels: context.format.numberOfChannels,
//   bitDepth: context.format.bitDepth,
//   sampleRate: context.sampleRate
// });

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(morgan('dev'));
app.use(cors());

// app.post('/', function(req, res) {
//
//   console.log(req.body);
//
//   var paramsForIndex = helpers.spherical_to_hrir({alpha: Number(req.body.ada.azimuth), beta: 90-Number(req.body.ada.altitude)});
//   var ret = helpers.hrir_index(paramsForIndex);
//
//   console.log(ret);
//
//   helpers.load_hrir();
//   res.send(200);
//
// });

app.get('/sound', function(req, res) {
    res.sendFile('/home/nicolas/Documents/coding/planesafe/stratux/GoPositionalAudio/planecloned.ogg');
});


app.listen(3000, function(){
    console.log('listening on port 3000');
});
