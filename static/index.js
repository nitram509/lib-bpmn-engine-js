'use strict';

(function initApplication() {
    let oldLog = console.log;
    console.log = function (message) {
        let box = $("#log");
        box.val(box.val() + message + "\n");
        oldLog.apply(console, arguments);
    };
})();

/**
 * @typedef ActivatedJob
 * @type {object}
 * @property {function:number} GetKey
 * @property {function:number} GetProcessInstanceKey
 * @property {function:string} GetBpmnProcessId
 * @property {function:number} GetProcessDefinitionVersion
 * @property {function:number} GetProcessDefinitionKey
 * @property {function:string} GetElementId
 * @property {function:Date} GetCreatedAt
 * @property {function(string)} Fail
 * @property {function} Complete
 */

/**
 *
 * @param job {ActivatedJob}
 */
function jobHandler(job) {
    console.log("Key                      = " + job.GetKey());
    console.log("ElementId                = " + job.GetElementId());
    console.log("BpmnProcessId            = " + job.GetBpmnProcessId());
    console.log("ProcessDefinitionKey     = " + job.GetProcessDefinitionKey());
    console.log("ProcessDefinitionVersion = " + job.GetProcessDefinitionVersion());
    console.log("CreatedAt                = " + job.GetCreatedAt());
    job.Complete();
}

function doIt() {
    let e = NewBpmnEngine()
    e.NewTaskHandlerForId("hello-world", jobHandler)
    e.CreateAndRunInstance()
}
