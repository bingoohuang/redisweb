// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "30s", "5h" or "2h45m30s".
// Valid time units are  "s", "m", "h".
function parseDuration(s) {
    if (s.length == 0) {
        return 0
    }

    var seconds = 0
    var num = ""
    for (var i = 0; i < s.length; ++i) {
        var c = s[i]
        if (c >= '0' && c <= '9') {
            num += c
        } else if (c == 'h') {
            seconds += parseInt(num) * 60 * 60
            num = ""
        } else if (c == 'm') {
            seconds += parseInt(num) * 60
            num = ""
        } else if (c == 's') {
            seconds += parseInt(num)
            num = ""
        } else {
            return 0
        }
    }

    return seconds
}

function createDuration(seconds) {
    var h = Math.floor(seconds / 60 / 60)
    var left = seconds - h * 60 * 60
    var m = Math.floor(left / 60)
    var s = left - m * 60

    var expr = ""
    if (h > 0) {
        expr += h + "h"
    }
    if (m > 0) {
        expr += m + "m"
    }
    if (s > 0) {
        expr += s + "s"
    }

    return expr
}

function joinKeysWithNo(keys) {
    var result = []
    var length = keys.length
    var lengthSize = ('' + length).length
    for (var i = 0; i < length; ++i) {
        result.push(pad(i + 1, lengthSize) + '.&nbsp;' + escapeHtml(keys[i]) + '<br>')
    }
    return result.join('')
}

function pad(n, width, z) {
    z = z || '0'
    n = n + ''
    return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n
}

function parseTemplateVariables(template) {
    var variables = []
    var variable = ""
    var started = false
    for (var i = 0, len = template.length; i < len; i++) {
        var char = template[i]
        if (char == '{') {
            started = true
        } else if (started == true && char == '}') {
            variables.push(variable)
            started = false
            variable = ""
        } else if (started == true) {
            variable += char
        }
    }

    return variables
}

function capitalize(s) {
    return s && s.charAt(0).toUpperCase() + s.slice(1)
}

function escapeHtml(unsafe) {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

