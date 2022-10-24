// modeler instance
var bpmnModeler = new BpmnJS({
    container: '#canvas',
    keyboard: {
        bindTo: window
    }
});

/**
 * Open diagram in our modeler instance.
 *
 * @param {String} bpmnXML diagram to display
 */
async function openDiagram(bpmnXML) {
    // import diagram
    try {
        await bpmnModeler.importXML(bpmnXML);
        // access modeler components
        var canvas = bpmnModeler.get('canvas');
        var overlays = bpmnModeler.get('overlays');
        // zoom to fit full viewport
        canvas.zoom('fit-viewport');
    } catch (err) {
        console.error('could not import BPMN 2.0 diagram', err);
    }
}


// load external diagram file via AJAX and open it
$.get("simple_task.bpmn", openDiagram, 'text');
