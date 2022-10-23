'use strict';

/**
 *
 * @returns {{
 *   GetName
 * }}
 * @constructor
 */
function NewBpmnEngine() {
  var idx = __newBpmnEngine();
  let retObj = {};
  retObj.GetName = __engine__getName.bind(idx)
  return retObj;
}
