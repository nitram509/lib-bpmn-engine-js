'use strict';

(function initApplication() {
    let oldLog = console.log;
    console.log = function (message) {
        let box = $("#log");
        box.val(box.val() + message + "\n");
        oldLog.apply(console, arguments);
    };
})();

function openSaveFileDialog(data, filename, mimetype) {
    'use strict';

    if (!data) return;

    var blob = data.constructor !== Blob
        ? new Blob([data], {type: mimetype || 'application/octet-stream'})
        : data;

    if (navigator.msSaveBlob) {
        navigator.msSaveBlob(blob, filename);
        return;
    }

    var lnk = document.createElement('a'),
        url = window.URL,
        objectURL;

    if (mimetype) {
        lnk.type = mimetype;
    }

    lnk.download = filename || 'untitled';
    lnk.href = objectURL = url.createObjectURL(blob);
    lnk.dispatchEvent(new MouseEvent('click'));
    setTimeout(url.revokeObjectURL.bind(url, objectURL));
}

/**
 * Save diagram contents and print them to the console.
 */
async function exportDiagram() {
    'use strict';
    try {
        let result = await bpmnModeler.saveXML({format: true});
        let filename = "diagram-" + new Date().getTime() + ".bpmn";
        openSaveFileDialog(result.xml, filename, 'application/bpmn+xml');
    } catch (err) {
        console.error('could not save BPMN 2.0 diagram', err);
    }
}

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
    const code = flask.getCode();
    eval(code);
}


async function runWorkflow() {
    let e = NewBpmnEngine()
    let bpmn = await bpmnModeler.saveXML({format: true});
    let processKey = e.LoadFromString(bpmn.xml);
    if (typeof processKey === 'string') {
        console.log("error loading bpmn: " + processKey);
    } else {
        let box = $("#log");
        box.val("");
        let ids = getTasksIds(bpmn.xml)
        ids.forEach(id => {
            e.NewTaskHandlerForId(id, jobHandler)
        })
        e.CreateAndRunInstance(processKey)
    }
}

/**
 *
 * @param xmlString
 * @returns {[string]}
 */
function getTasksIds(xmlString) {
    let parser = new DOMParser();
    let xmlDoc = parser.parseFromString(xmlString, "text/xml");
    let ids = [];
    xmlDoc.childNodes.forEach(function (element) {
        if (element.localName === 'definitions') {
            element.childNodes.forEach(function (element) {
                if (element.localName === 'process') {
                    element.childNodes.forEach(function (element) {
                        if (element.localName === 'serviceTask') {
                            ids.push(element.getAttribute("id"))
                        }
                    });
                }
            });
        }
    });
    return ids;
}
