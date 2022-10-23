'use strict';

(function initApplication() {
  let oldLog = console.log;
  console.log = function (message) {
    let box = $("#log");
    box.val(box.val() + message + "\n");
    oldLog.apply(console, arguments);
  };
})();
