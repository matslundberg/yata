var utils = require("./../lib/utils.js");

class CliOutput_ListTasks {
  render(tasks) {
    for(let task of tasks) {
      if(task.status != 'completed')
        this.render_task(task)
    }
  }

  render_task(task) {
    var colors = require('colors');
    //console.log('hello'.green)
    //console.log(documents);
    var chr = utils.status_to_char(task.status);
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

}

module.exports = new CliOutput_ListTasks();
