var utils = require("./utils.js");

class Data {
  constructor(dirname) {
    this.dirname = dirname;
  }

  process_notes(process_callback, end_callback) {
      var finder = require('findit')(this.dirname);
      var fs = require('fs');
      finder.on('file', function (file) {
        var contents = fs.readFileSync(file, 'utf8');
        process_callback(utils.basename(file), contents)
      });

      finder.on('end', end_callback);

  }

  match_query(data, query, entry) {
    if(data.match("/"+entry.type+"/")) {
      return false;
    }

    if(entry.status != query) {
      return false;
    }

    return true;
  }

  query(data, query, callback) {
    if(data == 'tasks') {
      var documents = []
      var match_query = this.match_query;

      this.process_notes(function(file, contents) {
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

            if(match_query(data, query, entry)) {
              documents.push(entry);
            }
          }

        }
      }, function(file) {
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
