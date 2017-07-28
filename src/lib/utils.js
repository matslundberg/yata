"use strict";

// From https://stackoverflow.com/a/18234317
String.prototype.format = String.prototype.format ||
function () {
    var str = this.toString();
    if (arguments.length) {
        var t = typeof arguments[0];
        var key;
        var args = ("string" === t || "number" === t) ?
            Array.prototype.slice.call(arguments)
            : arguments[0];

        for (key in args) {
            str = str.replace(new RegExp("\\{" + key + "\\}", "gi"), args[key]);
        }
    }

    return str;
};

module.exports = {
  char_to_status: function(chr) {
    var char_to_status = {
      "X": "completed",
      "x": "completed",
      "/": "ongoing",
      "-": "rejected",
      " ": "open"
    }

    if(chr == undefined)
      return char_to_status

    return char_to_status[chr];
  },

  status_to_char: function(status) {
    var chrstsmap = this.char_to_status();
    for(let chr in chrstsmap) {
      if(chrstsmap[chr] == status) {
        return ""+chr;
      }
    }

    return undefined;
  },

  // https://stackoverflow.com/questions/3820381/need-a-basename-function-in-javascript
  basename: function(str) {
    var base = new String(str).substring(str.lastIndexOf('/') + 1); 
    if(base.lastIndexOf(".") != -1)       
      base = base.substring(0, base.lastIndexOf("."));
    return base;
  }
}
