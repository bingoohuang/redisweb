$(function () {
    $('#redisInfo').click(function () {
        $.ajax({
            type: 'GET',
            url: contextPath + "/redisInfo",
            data: {server: $('#servers').val(), database: $('#databases').val()},
            success: function (result, textStatus, request) {
                var contentHtml = '<div><span class="key">Redis info</span></div>' +
                    '<pre>' + result + '</pre>'

                $('#frame').html(contentHtml)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("16." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    })
})