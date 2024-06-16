function explore(row, column, selectNode, action) {
    if (selectNode == null) {
        console.log(row, column, "selectNode is null");
        return;
    }
    console.log(row, column, selectNode.value, action);
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/instructions");
    xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8")
    xhr.send(JSON.stringify({
        row: row,
        column: column,
        action: action,
        selections: getSelectValues(selectNode)
    }));
    xhr.onload = function () {
        window.location.reload();
    }
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
    return result;
}