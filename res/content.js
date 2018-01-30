$(function () {
    $.showContentAjax = function (key) {
        $.ajax({
            type: 'GET', url: contextPath + "/showContent",
            data: {server: $('#servers').val(), database: $('#databases').val(), key: key},
            success: function (result, textStatus, request) {
                showContent(key, result.Type, result.Content, result.Ttl, result.Size, result.Encoding, result.Error, result.Exists, result.Format)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("17." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function showContent(key, type, content, ttl, size, encoding, error, exists, format) {
        if (error != "") {
            contentHtml = '<div><span class="error">' + error + '</span></div>'
            $('#frame').html(contentHtml)
        }

        if (!exists) {
            contentHtml = '<div><span class="key">' + key + ' does not exits</span></div>'
            $('#frame').html(contentHtml)
            return
        }

        var contentHtml = '<div><span class="key">' + key + '</span></div>' +
            '<table class="contentTable">' +
            '<tr><td class="titleCell">Type:</td><td colspan="2">' + type + '</td></tr>' +
            '<tr><td class="titleCell">TTL:</td><td colspan="2" contenteditable="true" id="ttl">' + ttl + '</td></tr>' +
            '<tr><td class="titleCell">Encoding:</td><td colspan="2">' + encoding + '</td></tr>' +
            '<tr><td class="titleCell">Format:</td><td colspan="2">' + format + '</td></tr>' +
            '<tr><td class="titleCell">Size:</td><td colspan="2">' + size + '</td></tr>' +
            '<tr><td class="titleCell">Value:</td><td colspan="2"><span class="valueSave sprite sprite-save"></span><span class="keyDelete sprite sprite-delete"></span></td></tr>'

        switch (type) {
            case "string":
                contentHtml += '<tr class="newKeyTr string"><td colspan="3"><textarea id="code">' + content + '</textarea></td></tr>'
                break
            case "hash":
                contentHtml += '<tr class="newKeyTr hash"><td class="titleCell">Field</td><td class="titleCell" colspan="2">Value</td></tr>'
                for (var hashKey in content) {
                    contentHtml += '<tr class="newKeyTr hash"><td contenteditable="true">' + hashKey + '</td><td colspan="2" contenteditable="true">' + content[hashKey] + '</td></tr>'
                }
                break
            case "set":
            case "list":
                contentHtml += '<tr class="newKeyTr ' + type + '"><td class="titleCell">#</td><td class="titleCell" colspan="2">Value</td></tr>'
                for (var i = 0; i < content.length; ++i) {
                    contentHtml += '<tr class="newKeyTr ' + type + '"><td contenteditable="true">' + i + '</td><td colspan="2" contenteditable="true">' + content[i] + '</td></tr>'
                }
                break
            case "zset":
                contentHtml += '<tr class="newKeyTr zset"><td class="titleCell">#</td><td class="titleCell">Score</td><td class="titleCell">Value</td></tr>'
                for (var i = 0; i < content.length; ++i) {
                    contentHtml += '<tr class="newKeyTr zset"><td contenteditable="true">' + i + '</td><td contenteditable="true">' + content[i].Score + '</td><td>' + content[i].Member + '</td></tr>'
                }
                break
        }
        contentHtml += '</table><button id="addMoreRowsBtn">Add More Rows</button>'

        $('#frame').html(contentHtml)
        $('#addMoreRowsBtn').toggle(type != "string").click(function () {
            addMoreRows(type)
        })

        var showContentTtlInterval = null

        var seconds = parseDuration(ttl)
        if (seconds > 0) {
            var ttlTd = $('#ttl')
            showContentTtlInterval = setInterval(function () {
                seconds -= 1
                if (seconds > 0) {
                    ttlTd.text(createDuration(seconds))
                } else {
                    ttlTd.text(0)
                    showContentTtlInterval = null
                }
            }, 1000)
        }


        $.codeMirror = null
        if (format === "JSON" && size < 5000) {
            $.codeMirror = CodeMirror.fromTextArea(document.getElementById('code'), {
                mode: 'application/json', lineNumbers: true, matchBrackets: true
            })
        } else {
            autosize($('#code'))
        }

        $('.keyDelete').click(function () {
            if (confirm("Are you sure to delete " + key + "?")) {
                $.ajax({
                    type: 'POST', url: contextPath + "/deleteKey",
                    data: {server: $('#servers').val(), database: $('#databases').val(), key: key},
                    success: function (content, textStatus, request) {
                        if (content != 'OK') {
                            alert("20." + content)
                        } else {
                            $.removeKey(key)
                            $('#frame').html('<div><span class="key">' + key + ' was deleted</span></div>')
                        }
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                        alert("21." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                    }
                })
            }
        })

        $('.valueSave').click(function () {
            if (confirm("Are you sure to save changes for " + key + "?")) {
                var changedContent = $.extractValue(type)
                $.ajax({
                    type: 'POST', url: contextPath + "/changeContent",
                    data: {
                        server: $('#servers').val(), database: $('#databases').val(),
                        key: key, type: type, ttl: $('#ttl').text(), value: changedContent
                    },
                    success: function (content, textStatus, request) {
                        if (content == 'OK') {
                            $.showContentAjax(key)
                        } else {
                            alert("22." + content)
                        }
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                        alert("23." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                    }
                })
            }
        })
    }

    function addMoreRows(type) {
        var rows = $('tr.' + type)
        var startRowNum = rows.length - 1

        var moreRows = ''
        for (var i = 0; i < 10; ++i) {
            if (type == 'hash') {
                moreRows += '<tr class="newKeyTr hash"><td contenteditable="true"></td><td colspan="2" contenteditable="true"></td></tr>'
            } else if (type == 'list' || type == 'set') {
                moreRows += '<tr class="newKeyTr list set"><td>' + (startRowNum + i) + '</td><td colspan="2" contenteditable="true"></td></tr>'
            } else if (type == 'zset') {
                moreRows += '<tr class="newKeyTr zset"><td>' + (startRowNum + i) + '</td><td contenteditable="true"></td><td contenteditable="true"></td></tr>'
            }
        }
        $(moreRows).appendTo($('.contentTable'))
    }

    $.extractValue = function (type) {
        var value = null
        if (type == 'string') {
            value = $.codeMirror != null && $.codeMirror.getValue() || $('#code').val()
        } else if (type == 'hash') {
            value = {}
            $('tr.hash:gt(0)').each(function (index, tr) {
                var tds = $(tr).find('td')
                var key = $.trim(tds.eq(0).text())
                var val = $.trim(tds.eq(1).text())
                if (key != "" && val != "") {
                    value[key] = val
                }
            })
        } else if (type == 'list' || type == 'set') {
            value = []
            $('tr.' + type + ':gt(0)').each(function (index, tr) {
                var tds = $(tr).find('td')
                var val = $.trim(tds.eq(1).text())
                if (val != "") {
                    value.push(val)
                }
            })
        } else if (type == 'zset') {
            value = []
            $('tr.zset:gt(0)').each(function (index, tr) {
                var tds = $(tr).find('td')
                var score = $.trim(tds.eq(1).text())
                var val = $.trim(tds.eq(2).text())
                if (score != "" && val != "") {
                    value.push({Score: +score, Member: val})
                }
            })
        }

        return JSON.stringify(value)
    }

})