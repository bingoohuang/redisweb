$(function () {
    $('#redisImport').click(function () {
        var contentHtml = '<div><h3>Commands Import</h3></div>' +
            '<div id="importDiv"><textarea id="importCommands" cols="100" rows="20"></textarea>' +
            '<br><button id="redisImportBtn">Import</button>' +
            '<div id="importResult"></div>' +
            '</div>'
        $('#frame').html(contentHtml)
        var importCodeMirror = CodeMirror.fromTextArea(document.getElementById('importCommands'), {
            lineNumbers: true, matchBrackets: true, height: 500
        })

        $('#redisImportBtn').click(function () {
            var commands = importCodeMirror.getValue()
            $.ajax({
                type: 'GET', url: contextPath + "/redisImport",
                data: {server: $('#servers').val(), database: $('#databases').val(), commands: commands},
                success: function (content, textStatus, request) {
                    $('#importResult').html('<pre>' + content + '</pre>')
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    alert("2." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })
        })
    })
})