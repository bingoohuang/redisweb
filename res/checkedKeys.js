$(function () {
    $.findCheckedKeys = function() {
        var keys = []
        $('#keys ul li:visible').each(function (index, li) {
            var $li = $(li)
            if ($li.find('input:checkbox').is(":checked")) {
                var key = $li.find('.keyValue').text()
                keys.push(key)
            }
        })
        return keys
    }

    $('#deleteCheckedKeys').click(function () {
        var keys = $.findCheckedKeys()
        if (keys.length == 0) {
            alert("5." + "No keys chosen to be deleted!")
            return
        }

        if (!confirm("Are you sure to delete " + keys.length + " keys?")) {
            return
        }

        $.ajax({
            type: 'POST', url: contextPath + "/deleteMultiKeys",
            data: {server: $('#servers').val(), database: $('#databases').val(), keys: JSON.stringify(keys)},
            success: function (content, textStatus, request) {
                if (content != 'OK') {
                    alert("6." + content)
                } else {
                    removeKeys(keys)
                    $('#checkAllChk').prop('checked', false)
                    $('#frame').html('<div><span class="key">' + keys.length + ' keys were deleted:</span></div>'
                        + '<div><br>' + joinKeysWithNo(keys) + '</div>')
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("7." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    })


    function removeKeys(keys) {
        $('#keys ul li').removeClass('chosen').each(function (index, li) {
            var $span = $(li).find('.keyValue')
            if ($.inArray($span.text(), keys) > -1) {
                $(li).remove()
            }
        })
    }
})