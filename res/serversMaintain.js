$(function () {
    function updateDatabases(content) {
        var options = ''
        for (var i = 0; i < content.Dbs; ++i) {
            if (i === content.DefaultDb) {
                options += '<option value="' + i + '" selected>' + i + '</option>'
            } else {
                options += '<option value="' + i + '">' + i + '</option>'
            }
        }

        $('#databases').html(options)
        $.refreshKeys()
    }

    $('#databases').change($.refreshKeys)
    $('#servers').change(function () {
        var redisServer  = $("#servers option:selected").text()
        $.ajax({
            type: 'GET', url: contextPath + "/changeRedisServer",
            data: {redisServer : redisServer},
            success: function (content, textStatus, request) {
                if (content.OK !== "OK") {
                    alert(content.OK)
                } else {
                    updateDatabases(content)
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("2." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    })

    $('#maintainRedisServers').click(function (e) {
        var contentHtml = '<div><h3>Redis Servers</h3></div>' +
            '<div id="redisServerDiv"><textarea id="redisServerTextArea" cols="100" rows="20"></textarea>' +
            '<br><button id="redisServerSave">Save</button>' +
            '</div>'
        $('#frame').html(contentHtml)
        var codeMirror = CodeMirror.fromTextArea(document.getElementById('redisServerTextArea'), {
            mode: 'text/x-toml', lineNumbers: true
        })

        $.ajax({
            type: 'GET', url: contextPath + "/loadRedisServerConfig",
            success: function (content, textStatus, request) {
                codeMirror.setValue(content.RedisServerConfig)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("2." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })

        $('#redisServerSave').click(function () {
            var redisServerConfig = codeMirror.getValue()
            $.ajax({
                type: 'GET', url: contextPath + "/saveRedisServerConfig",
                data: {redisServerConfig: redisServerConfig},
                success: function (content, textStatus, request) {
                    if (content.OK !== "OK") {
                        alert(content.OK)
                    } else {
                        var serverOptions = ''
                        for (var i = 0; i < content.Servers.length; ++i) {
                            serverOptions += '<option value="' + content.Servers[i] + '">' + content.Servers[i] + '</option>'
                        }

                        $('#servers').html(serverOptions)
                        updateDatabases(content)
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    alert("2." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })
        })
    })
})