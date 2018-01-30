$(function () {
    $('#exportKeys').click(function () {
        var contentHtml = '<div><h3>Export</h3></div>' +
            '<div><span style="margin-right: 30px">Type:</span>' +
            '<select style="margin-right: 30px" id="exportType"><option value="Redis">Redis</option><option value="JSON">JSON</option></select>' +
            '<button id="exportBtn">Export</button><button id="downloadExportBtn">Download</button></div>' +
            '<div id="exportResult"></div>'

        $('#frame').html(contentHtml)

        $('#exportBtn,#downloadExportBtn').click(function () {
            var keys = $.findCheckedKeys()
            if (keys.length == 0) {
                alert("3." + "No keys chosen to be deleted!")
                return
            }
            var exportType = $('#exportType').val()
            var download = $(this).prop('id') == 'downloadExportBtn'

            var data = {
                server: $('#servers').val(),
                database: $('#databases').val(),
                exportKeys: JSON.stringify(keys),
                exportType: exportType,
                download: download
            }

            if (download) {
                window.open(contextPath + "/exportKeys?" + $.param(data), '_blank')
                return
            }

            $.ajax({
                type: 'GET', url: contextPath + "/exportKeys",
                data: data,
                success: function (content, textStatus, request) {
                    if (exportType == "Redis") {
                        $('#exportResult').html('<pre>' + content.join('<br>') + '</pre>')
                    } else {
                        $('#exportResult').html('<textarea id="code">' + content + '</textarea>')
                        CodeMirror.fromTextArea(document.getElementById('code'), {
                            mode: 'application/json', lineNumbers: true, matchBrackets: true
                        })
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    alert("4." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })

        })
    })

})