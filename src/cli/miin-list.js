#!/usr/bin/env node

require('./../lib/utils.js');

var program = require('commander');
program
  .command('<data> [query]', 'List all tags')
  .parse(process.argv);


var query = program.args;
var data  = program.args[0]
query.shift();
query = query.join(" ");

console.log("Getting %s %s", data, query)

var __dirname = process.env.NOTES_DIR

try {
  var miin = require('./../lib/data.js')(__dirname)
  var view = require('./clioutput_listtasks.js')

  var tasks = miin.query(data, query, function(tasks) {
    view.render(tasks);
  });

} catch(exception) {
  console.log("Something went wrong! %s", exception)
}
