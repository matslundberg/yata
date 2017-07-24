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
var lunr = require('lunr')
var finder = require('findit')(__dirname);
var fs = require('fs');

var char_to_status = function(chr) {
  var char_to_status = {
    "x": "completed",
    "/": "ongoing",
    "-": "rejected",
    " ": "open"
  }

  if(chr == undefined)
    return char_to_status

  return char_to_status[chr];
}

var status_to_char = function(status) {
  var chrstsmap = char_to_status();
  for(chr in chrstsmap) {
    if(chrstsmap[chr] == status) {
      return ""+chr;
    }
  }

  return undefined;
}

var documents = []

//console.log(finder);
//finder.find
//This listens for files found
finder.on('file', function (file) {
  console.log('File: ' + file);
  var contents = fs.readFileSync(file, 'utf8');
  //documents.push({'name':file, text:contents});
  var todo_regex = /(\[\s*( |x|X|\/|-)\s*\](.+))/;

  var lines = contents.split('\n');
  for(line of lines) {
    match = todo_regex.exec(line)
    //console.log(match)
    if(match != null) {
      var entry = {
        'type': "task",
        'status': char_to_status( match[2] ),
        'description': match[3],
        'source': file,
      }

      documents.push(entry);
    }

  }
  //var todos = contents.split(/ /)
  //console.log(documents);
});

var render_task = function(task) {
  var colors = require('colors');
  //console.log('hello'.green)
  //console.log(documents);
  var chr = status_to_char(task.status);
  if(chr == undefined) {
    var out = "      {0}".format(task.description);
  } else {
    var out = "[{0}] {1} ({2})".format(chr, task.description, task.source);
  }
  if(task.status == "completed") {
    console.log(out.green)
  } else if(task.status == "ongoing") {
    console.log(out.yellow)
  } else if(task.status == "rejected") {
    console.log(out.strikethrough.grey)
  } else if(task.status == "overdue") {
    console.log(out.red)
  } else {
    console.log(out)
  }

}

finder.on('end', function (file) {
  /*
  var idx = lunr(function () {
    this.ref('name')
    this.field('text')

    documents.forEach(function (doc) {
      this.add(doc)
    }, this)
  })

  console.log(documents, idx)

  var test = idx.search('*');
  console.log(test);
  */
  var tasks = documents;
  for(task of tasks) {
    //console.log(task)
    if(task.status != 'completed')
      render_task(task)
  }
})



class Muistiin {
  constructor() {

  }

  query(data, query) {
    console.log("data:", data)
    console.log("query:", query)
    //return [{"field":"1"}, {"field":"2"}]
    if(data == 'tasks') {

    } else {
      throw "Do not understand data: {0}".format(data);
    }
  }
}

class CliOutput_List {
  render(entries) {
    for(let entry in entries) {
      console.log(entries[entry]);
    }
  }
}

try {
  var miin = new Muistiin();
  var view = new CliOutput_List();

  var tasks = miin.query(data, query);

  view.render(tasks);
} catch(exception) {
  console.log("Something went wrong! {0}".format(exception))
}
