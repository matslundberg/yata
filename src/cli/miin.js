#!/usr/bin/env node
console.log('in miin.js');
require('dotenv').config()
var program = require('commander');

program
  .command('tags [query]', 'List all tags')
  .command('list tasks [query]', 'List all tasks')
  .command('open [note]', 'List all tags')
  .parse(process.argv);
