var utils = require("./utils.js");

class Data {
  constructor(dirname) {
    this.dirname = dirname;
  }

  get_raw_note(note) {
    var fs = require('fs');
    var contents = fs.readFileSync(this.dirname+note, 'utf8');
    return contents;
  }

  write_raw_note(note, contents) {
    var fs = require('fs');
    fs.writeFileSync(this.dirname+note, contents, 'utf8');
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

    if(query == "*") {
      return true;
    }

    if(entry.status != query) {
      return false;
    }

    return true;
  }

  get_ref_by_id(id) {
    var jsonfile = require('jsonfile')
    var refs_path = this.dirname+"_refs";
    var refs = jsonfile.readFileSync(refs_path);
    //console.log(refs);
    if(refs[id] == undefined || refs[id] == null) {
      //console.log(refs.length)
      var ref = 1+Object.keys(refs).length;
      refs[id] = ""+ref;
      jsonfile.writeFileSync(refs_path, refs)
    } else {
      var ref = refs[id];
    }
    
    //refs = JSON.parse(get_raw_note("_refs"));
    return ref;
  }

  complete(data, ref, callback) {
    if(data == 'tasks') {

      this.query(data, "*", function(tasks) {
        for(let task of tasks) {
          if(task.ref === ref) {
            console.log("Found task to close!! %s", task.id);
          }
        }
      })

    } else {
      throw "Do not understand data: {0}".format(data);
    }
    
  }

  query(data, query, callback) {
    if(data == 'tasks') {
      var documents = []
      var used_ids = {}
      var match_query = this.match_query;
      var sha1 = require('sha1');
      var self = this;

      this.process_notes(function(file, contents) {
        var todo_regex = /(\[\s*( |x|X|\/|-)\s*\](.+))/;

        var lines = contents.split('\n');
        var lc = 0;
        for(let line of lines) {
          lc++;
          var match = todo_regex.exec(line)
          //console.log(match)
          if(match != null) {
            var id = sha1(line+file+lc) //.substring(0, 6)
            var ref = self.get_ref_by_id(id);

            if(used_ids[id] == true) {
              throw "Duplicate Todo ID used!";
            }

            var entry = {
              'type': "task",
              'status': utils.char_to_status( match[2] ),
              'description': match[3],
              'source': file,
              'id': id,
              'ref': ref,
            }

            if(match_query(data, query, entry)) {
              used_ids[id] = true;
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
