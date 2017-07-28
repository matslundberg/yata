#!/usr/bin/env node

require('./../lib/utils.js');

var program = require('commander');
program
  .command('<data> [query]', 'List all tags')
  .parse(process.argv);


var ref = program.args[1];
var data  = program.args[0]

console.log("Closing %s %s", data, ref)

var __dirname = process.env.NOTES_DIR

try {
  var miin = require('./../lib/data.js')(__dirname)
  var view = require('./clioutput_listtasks.js')

  var tasks = miin.complete(data, ref, function(tasks) {
    console.log("Closed %s %s", data, ref)
  });

} catch(exception) {
  console.log("Something went wrong! %s", exception)
}
