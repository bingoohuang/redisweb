$(function () {
    $('#addKey').click(function () {
        var contentHtml = '<div><span class="key">Add another key</span></div>' +
            '<table class="contentTable">' +
            '<tr><td class="titleCell">Type:</td><td colspan="2"><select name="type" id="type">' +
            '<option value="string">String</option><option value="hash">Hash</option><option value="list">List</option>' +
            '<option value="set">Set</option><option value="zset">Sorted Set</option>' +
            '</select></td></tr>' +
            '<tr><td class="titleCell">Key:</td><td colspan="2"><input id="key" placeholder="input the new key"></td></tr>' +
            '<tr><td class="titleCell">TTL:</td><td colspan="2"><input id="ttl" placeholder="input the expired time, like 1d/1h/10s/-1s"></td></tr>' +
            '<tr><td class="titleCell">Format:</td><td colspan="2"><select name="format" id="format">' +
            '<option value="String">String</option><option value="JSON">JSON</option><option value="Quoted">Quoted</option>' +
            '</select></td></tr>' +
            '<tr><td class="titleCell">Value:</td><td colspan="2"><span class="valueSave sprite sprite-save"></span></td></tr>'

        contentHtml += '<tr class="newKeyTr string"><td colspan="2"><textarea id="code"></textarea></td></tr>'

        contentHtml += '<tr class="newKeyTr hash"><td class="titleCell">Field</td><td colspan="2" class="titleCell">Value</td></tr>'
        for (var i = 0; i < 10; ++i) {
            contentHtml += '<tr class="newKeyTr hash"><td contenteditable="true"></td><td colspan="2" contenteditable="true"></td></tr>'
        }

        contentHtml += '<tr class="newKeyTr list set"><td class="titleCell">#</td><td colspan="2" class="titleCell">Value</td></tr>'
        for (var i = 0; i < 10; ++i) {
            contentHtml += '<tr class="newKeyTr list set"><td>' + i + '</td><td colspan="2" contenteditable="true"></td></tr>'
        }

        contentHtml += '<tr class="newKeyTr zset"><td class="titleCell">#</td><td class="titleCell">Score</td><td class="titleCell">Value</td></tr>'
        for (var i = 0; i < 10; ++i) {
            contentHtml += '<tr class="newKeyTr zset"><td>' + i + '</td><td contenteditable="true"></td><td contenteditable="true"></td></tr>'
        }

        contentHtml += '</table><button id="addMoreRowsBtn">Add More Rows</button>'

        $('#frame').html(contentHtml)

        $('tr.newKeyTr').hide()
        $('tr.string').show()
        $('#addMoreRowsBtn').hide().click(function () {
            addMoreRows($('#type').val())
        })


        $('#type').change(function () {
            var type = $('#type').val()
            $('tr.newKeyTr').hide()
            $('tr.' + type).show()
            $('#addMoreRowsBtn').toggle(type != 'string')
        })

        $.codeMirror = null
        $('#format').change(function () {
            if ($(this).val() == 'JSON' && $('#type').val() == 'string') {
                $.codeMirror = CodeMirror.fromTextArea(document.getElementById('code'), {
                    mode: 'application/json', lineNumbers: true, matchBrackets: true
                })
            } else {
                $.codeMirror = null
            }
        })
        valueSaveClickable()
    })


    function valueSaveClickable() {
        $('.valueSave').click(function () {
            var type = $('#type').val()
            var key = $('#key').val()
            var ttl = $('#ttl').val()
            var format = $('#format').val()
            var jsonValue = $.extractValue(type)

            if (confirm("Are you sure to save save for " + key + "?")) {
                $.ajax({
                    type: 'POST', url: contextPath + "/newKey",
                    data: {
                        server: $('#servers').val(), database: $('#databases').val(),
                        type: type, key: key, ttl: ttl, format: format, value: jsonValue
                    },
                    success: function (content, textStatus, request) {
                        if (content == 'OK') {
                            $.refreshKeys(key)
                            $.showContentAjax(key)
                        } else {
                            alert("18." + content)
                        }
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                        alert("19." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                    }
                })
            }
        })
    }
})