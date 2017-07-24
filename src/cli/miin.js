#!/usr/bin/env node
console.log('in miin.js');

var program = require('commander');

program
  .command('tags [query]', 'List all tags')
  .command('list tasks [query]', 'List all tasks')
  .parse(process.argv);
