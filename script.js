function explore(row, column, selectNode, action) {
    if (selectNode == null) {
        console.log(row, column, "selectNode is null");
        return;
    }
    console.log(row, column, selectNode.value, action);
    const xhr = new XMLHttpRequest();
    xhr.open("POST", window.location.href);
    xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8")
    xhr.send(JSON.stringify({
        row: row,
        column: column,
        action: action,
        selections: getSelectValues(selectNode)
    }));
    xhr.onload = function() { window.location.reload(); }
}

function resume() {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", window.location.href);
    xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8")
    xhr.send(JSON.stringify({
        action: "resume" 
    }));
}

// Return an array of the selected option values in the control.
// Select is an HTML select element.
function getSelectValues(select) {
    var result = [];
    var options = select && select.options;
    var opt;

    for (var i = 0, iLen = options.length; i < iLen; i++) {
        opt = options[i];

        if (opt.selected) {
            result.push(opt.value || opt.text);
        }
    }
    if (result.length == 0 && options.length == 1) {
        opt = options[0];
        result.push(opt.value || opt.text);
    }
    return result;
} 