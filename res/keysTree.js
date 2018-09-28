$(function () {
    function refreshKeys(key) {
        $.ajax({
            type: 'GET', url: contextPath + "/listKeys",
            data: {server: $('#servers').val(), database: $('#databases').val(), pattern: $('#serverFilterKeys').val()},
            success: function (content, textStatus, request) {
                showKeysTree(content)
                if (typeof key === 'string' || key instanceof String) {
                    chosenKey(key)
                }

                var hashKey = $.hash().get('key')
                if (hashKey) {
                    chosenKey(hashKey)
                    $.showContentAjax(hashKey)
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("1." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    setTimeout(function () {
        // listen url hash
        $.hash().listen("key", function (hashKey) {
            chosenKey(hashKey.key)
            $.showContentAjax(hashKey.key)
        })
    }, 100)


    $('#serverFilterKeysBtn,#refreshKeys').click(refreshKeys)
    $('#serverFilterKeys').keydown(function (event) {
        var keyCode = event.keyCode || event.which
        if (keyCode == 13) {
            refreshKeys()
        }
    })

    $('#checkAllChk').click(function () {
        var checked = $('#checkAllChk').is(":checked")
        $('#keys ul li:visible').each(function (index, li) {
            $(li).find('input:checkbox').prop('checked', checked)
        })
    })

    $.refreshKeys = refreshKeys

    function showKeysTree(keysArray) {
        $('#keysNum').html('(' + keysArray.length + ')')

        var keysHtml = '<ul>'
        for (var i = 0; i < keysArray.length; ++i) {
            var key = keysArray[i]
            var nodeCss = i < keysArray.length - 1 ? "sprite-tree-node" : "sprite-tree-lastnode last"
            keysHtml += '<li class="datatype-' + key.Type + ' sprite ' + nodeCss + '" data-type="' + key.Type + '">' +
                '<input type="checkbox">' +
                '<span class="sprite sprite-datatype-' + key.Type + '"></span>' +
                '<span class="keyValue">' + escapeHtml(key.Key) + '</span>'

            var numInfo = key.Len != -1 ? '(' + key.Len + ')' : ''
            keysHtml += '<span class="info">' + numInfo + '</span></li>'
        }
        keysHtml += '</ul>'

        $('#keys').html(keysHtml)

        $('#keys ul li span.keyValue').click(function () {
            $('#keys ul li').removeClass('chosen')
            var $li = $(this).parent('li')
            $li.addClass('chosen')
            var key = $li.find('.keyValue').text()
            $.showContentAjax(key)
        })

        toggleFilterKeys()
    }

    function toggleFilterKeys() {
        var filter = $.trim($('#filterKeys').val()).toUpperCase()

        $('#keys ul li').each(function (index, li) {
            var $li = $(li)
            var text = $.trim($li.text()).toUpperCase()
            var contains = text.indexOf(filter) > -1
            $li.toggle(contains)
        })

        // $('#sidebar').height(window.outerHeight)
    }

    $('#filterKeys').keyup(toggleFilterKeys)

    function chosenKey(key) {
        $('#keys ul li').removeClass('chosen').each(function (index, li) {
            var $span = $(li).find('.keyValue')
            if ($span.text() == key) {
                $(li).addClass('chosen')
                return false
            }
        })
    }


    $.removeKey = function (key) {
        $('#keys ul li').removeClass('chosen').each(function (index, li) {
            var $span = $(li).find('.keyValue')
            if ($span.text() == key) {
                $(li).remove()
                return false
            }
        })
    }


    var keysFocused = false

    $('#keys').attr('tabindex', -1).focusin(function () {
        console.info('focusin')
        keysFocused = true
    }).focusout(function () {
        console.info('focusout')
        keysFocused = false
    })

    $(document).keydown(function (e) {
        var which = e.which
        switch (which) {
            case 37: // left
            case 38: // up
            case 39: // right
            case 40: // down
                if (keysFocused) {
                    $('#keys ul li:visible').each(function (index, li) {
                        $li = $(li)
                        if ($li.hasClass('chosen')) {
                            (which == 37 || which == 38 ? $li.prev(':visible') : $li.next(':visible')).find('span.keyValue').click()
                            e.preventDefault()
                            return false
                        }
                    })
                }
                break
            default:
                return // exit this handler for other keys
        }
    })
}())