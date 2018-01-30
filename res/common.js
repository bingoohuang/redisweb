var contextPath = window.location.pathname
if (contextPath.lastIndexOf("/", contextPath.length - 1) !== -1) {
    contextPath = contextPath.substring(0, contextPath.length - 1)
}

var windowAlert = window.alert

window.alert = function(msg) {
    if (msg) windowAlert(msg)
}

$(document).on('paste', '[contenteditable]', function (e) {
    e.preventDefault()
    var text = ''
    if (e.clipboardData || e.originalEvent.clipboardData) {
        text = (e.originalEvent || e).clipboardData.getData('text/plain')
    } else if (window.clipboardData) {
        text = window.clipboardData.getData('Text')
    }
    if (document.queryCommandSupported('insertText')) {
        document.execCommand('insertText', false, text)
    } else {
        document.execCommand('paste', false, text)
    }


    $.codeMirror = null

})
