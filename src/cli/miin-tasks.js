#!/usr/bin/env node
console.log('Tasks');
var program = require('commander');
program
  .command('[query]', 'List all tags')
  .parse(process.argv);

//console.log(program)
for(arg of program.args) {
  //console.log(arg);
}

task_query = program.args.join(" ");
console.log(task_query);

// note we are placing WhiteSpace first as it is very common thus it will speed up the lexer.
//let allTokens = [WhiteSpace, Select, From, Where, Comma, Identifier, Integer, GreaterThan, LessThan]
/*
query_parser = require("../lib/query");

for(task_query of ['due next 7 day', 'due last 7 days', 'due tomorrow']) {
  console.log(task_query)
  ast = query_parser.toAst(task_query)
  console.log(ast)
}
*/
