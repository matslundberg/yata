#!/usr/bin/env node

require('./../lib/utils.js');

console.log('Lists');
var program = require('commander');
program
  .command('<data> [query]', 'List all tags')
  .parse(process.argv);

//console.log(program)
for(arg of program.args) {
  //console.log(arg);
}

var query = program.args;
var data  = program.args[0]
query.shift();
query = query.join(" ");
console.log("data:", data)
console.log("query:", query)

var __dirname = '/mnt/c/Users/mats.lundberg.INTERNAL/Dropbox/notes/'
//var __dirname = './tests/'
//var lunr = require('lunr')
var awaitEvent = require('await-event')


//const Database = require('./database');
//console.log(finder);
//finder.find
//This listens for files found


//console.log( awaitEvent(finder, 'end') );
//yield awaitEvent(finder, 'end')




try {
  var miin = require('./../lib/data.js')(__dirname)
  var view = require('./clioutput_listtasks.js')

  var tasks = miin.query(data, query, function(tasks) {
    view.render(tasks);
  });

} catch(exception) {
  console.log("Something went wrong! %s", exception)
}
