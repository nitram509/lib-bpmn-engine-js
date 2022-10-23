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
 * @typedef NewBpmnEngine
 * @type {function}
 * @returns {BpmnEngine}
 */

/**
 * @typedef BpmnEngine
 * @type {object}
 * @property {function(id:string, callback:function(ActivatedJob))} NewTaskHandlerForId
 * @property {function(bpmn:string):number} LoadFromString loads BPMN returns processKey
 * @property {function(processKey:number)} CreateAndRunInstance
 */

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
 * @property {function(reason:string)} Fail
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


async function runWorkflow() {
    let e = NewBpmnEngine()
    let bpmn = await bpmnModeler.saveXML({format: true});
    let processKey = e.LoadFromString(bpmn.xml);
    if (typeof processKey === 'string') {
        console.log("error loading bpmn: " + processKey);
    } else {
        e.NewTaskHandlerForId("id", jobHandler)
        e.CreateAndRunInstance(processKey)
    }
}
