var utils = require("./utils.js");

class Data {
  constructor(dirname) {
    this.dirname = dirname;
  }

  query(data, query, callback) {
    console.log("data:", data)
    console.log("query:", query)

    if(data == 'tasks') {
      var finder = require('findit')(this.dirname);
      var fs = require('fs');
      var documents = []

      finder.on('file', function (file) {
        var contents = fs.readFileSync(file, 'utf8');
        var todo_regex = /(\[\s*( |x|X|\/|-)\s*\](.+))/;

        var lines = contents.split('\n');
        for(let line of lines) {
          var match = todo_regex.exec(line)
          //console.log(match)
          if(match != null) {
            var entry = {
              'type': "task",
              'status': utils.char_to_status( match[2] ),
              'description': match[3],
              'source': file,
            }

            documents.push(entry);
          }

        }
      });

      finder.on('end', async function (file) {
        callback(documents)
      })
    } else {
      throw "Do not understand data: {0}".format(data);
    }
  }
}

module.exports = function(dirname) {
  return new Data(dirname);
}
