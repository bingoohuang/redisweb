$(function () {
    function executeRedisCmd() {
        var cmd = $('#directCmd').text()
        var server = $('#servers').val()
        $.ajax({
            type: 'POST', url: contextPath + "/redisCli",
            data: {server: server, database: $('#databases').val(), cmd: cmd},
            success: function (result, textStatus, request) {
                var resultHtml = server + '&gt;&nbsp;' + cmd + '<pre>' + result + '</pre>'

                $('#directCmdResult div').append(resultHtml)
                $('#directCmd').text('')

                setTimeout(function () {
                    var scrollValue = $('#directCmdResult').height() - $(window).height()
                    $('#frame').animate({scrollTop: scrollValue + 100}, 800)
                }, 0)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("8." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    $('#redisTerminal').click(function () {
        $('#frame').html('<div id="directCmdResult"><div></div><span id="cmdPrompt"></span><span contenteditable="true" id="directCmd"></div>')
        $('#cmdPrompt').html($('#servers').val() + '&gt;&nbsp;')

        $('#directCmd').focus().keydown(function (event) {
            var keyCode = event.keyCode || event.which
            if (keyCode == 13) {
                executeRedisCmd()
            }
        })
        $('#directCmdResult').click(function () {
            $('#directCmd').focus()
        })
    })
})

