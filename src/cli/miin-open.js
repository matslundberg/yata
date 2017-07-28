#!/usr/bin/env node

require('./../lib/utils.js');

var program = require('commander');
program
  .command('<data> [query]', 'test')
  .parse(process.argv);


var note = program.args[0]
var __dirname = process.env.NOTES_DIR
var ExternalEditor = require('external-editor');

try {
  if(!note)  {
      throw "No note given!";
  } 

  var miin = require('./../lib/data.js')(__dirname)
  var note_contents = miin.get_raw_note(note);

  var editor = new ExternalEditor(note_contents);
  var new_note_contents = editor.run()

  miin.write_raw_note(note, new_note_contents);

  editor.cleanup();   

} catch (exception) {
    console.log("Something went wrong! %s", exception)
}
