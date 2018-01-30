$(function () {
    var convenientConfig = null

    $('#convenientSpan').click(function () {
        $.ajax({
            type: 'POST', url: contextPath + "/convenientConfig",
            success: function (result, textStatus, request) {
                if (result.Ready !== true) {
                    $('#frame').html('<div><span class="key">' + result.Error + '</span></div>')
                    return
                }

                convenientConfig = result
                var contentHtml = '<div class="itemNames">'
                for (var i = 0, len = convenientConfig.Items.length; i < len; i++) {
                    var convenientConfigItem = convenientConfig.Items[i]
                    contentHtml += '<span itemIndex="' + i + '" class="convenientConfigItem">' + convenientConfigItem.Name + '</span>'
                }
                contentHtml += '<span class="convenientConfigItem">New Item</span>'
                contentHtml += '</div><div id="convenientContent"></div>'
                $('#frame').html(contentHtml)
                itemClick()
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("15." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }).click()

    function newItem() {
        var convenientContent =
            '<div><span>Name:</span><span><input class="templateName"></span></div>' +
            '<div><span>Key Template:</span><span><input class="templateValue"></span></div>' +
            '<div><span>TTL:</span><span><input class="ttlCreated" value="-1s"></span></div>' +
            '<div><span>Operations:</span><span><input style="width:12px" type="checkbox" id="DeleteChk" value="Delete" checked><label for="DeleteChk">Delete</label>' +
            '<input  style="width:12px" type="checkbox" id="SaveChk" value="Save" checked><label for="SaveChk">Save</label></span></div>' +
            '<div><span>&nbsp;</span><span><button>Save New Item</button></span></div>'
        var $convenientContent = $('#convenientContent')
        $convenientContent.html(convenientContent)
        $convenientContent.find('button').click(function () {
            var operations = ''
            $convenientContent.find('input:checked').each(function (index, chk) {
                if (operations.length > 0) operations += ','
                operations += $(chk).val()
            })
            $.ajax({
                type: 'GET', url: contextPath + "/convenientConfigAdd",
                data: {
                    name: $convenientContent.find('input.templateName').val(),
                    template: $convenientContent.find('input.templateValue').val(),
                    operations: operations,
                    ttl: $convenientContent.find('input.ttlCreated').val()
                },
                success: function (result, textStatus, request) {
                    if (result.Message == 'OK') {
                        $('#convenientSpan').click()
                    } else {
                        alert("9." + result.Message)
                    }
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    alert("10." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })
        })
    }

    function saveItemContent(keyCreated) {
        var ttl = $('.ttlCreated input').val()
        $.ajax({
            type: 'POST', url: contextPath + "/changeContent",
            data: {
                server: $('#servers').val(), database: $('#databases').val(),
                key: keyCreated, type: 'string', ttl: ttl,
                value: JSON.stringify($(".valueTextArea").val())
            },
            success: function (content, textStatus, request) {
                setConvenientContentInfo(content == 'OK' && ttl)
                $('.resultSpan').text('Saved ' + content)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("12." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function deleteItemContent(keyCreated) {
        $.ajax({
            type: 'POST', url: contextPath + "/deleteKey",
            data: {server: $('#servers').val(), database: $('#databases').val(), key: keyCreated},
            success: function (content, textStatus, request) {
                if (content == 'OK') {
                    refreshValue('Deleted OK')
                } else {
                    $('.resultSpan').text('Deleted ' + content)
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("13." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function deleteItemConfig($this) {
        $.ajax({
            type: 'POST', url: contextPath + "/deleteConvenientConfigItem",
            data: {
                sectionName: $this.attr('sectionName')
            },
            success: function (content, textStatus, request) {
                $('#convenientSpan').click()
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("14." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function itemButtonClickable() {
        $('.convenientButton').click(function () {
            var $this = $(this)
            var keyCreated = $('.keyCreated').text()
            if ($this.text() == 'Save') {
                saveItemContent(keyCreated)
            } else if ($this.text() == 'Delete') {
                deleteItemContent(keyCreated)
            } else if ($this.text() == 'Refresh Value') {
                refreshValue()
            } else if ($this.text() == 'Delete This Item') {
                deleteItemConfig($this)
            }
        })
    }

    function showExistedItem($item) {
        var item = convenientConfig.Items[+$item.attr('itemIndex')]
        var convenientContent =
            '<div><span>Key Template:</span><span class="templateValue">' + item.Template + ' </span></div>'

        var variables = parseTemplateVariables(item.Template)
        for (var i = 0, len = variables.length; i < len; i++) {
            convenientContent += '<div class="variables"><span>' + variables[i] + ':</span><span><input placeholder="variable value"></span></div>'
        }

        convenientContent += '<div><span>Key:</span><span class="keyCreated"></span></div>'
            + '<div><span>TTL:</span><span class="ttlCreated"><input value="' + item.Ttl + '"></span></div>'
            + '<div><span>Value:<br/><span class="info"></span></span><span><textarea class="valueTextArea"></textarea></span></div>'
            + '<div><span>Operations:</span><span><button class="convenientButton">Refresh Value</button>'

        for (var i = 0, len = item.Operations.length; i < len; i++) {
            convenientContent += '<button class="convenientButton">' + capitalize(item.Operations[i]) + "</button>"
        }
        convenientContent += '<button class="convenientButton" sectionName="' + item.Section + '">Delete This Item</button>'
         + '</span></div>'
         + '<div><span>Result:</span><span class="resultSpan"></span></div>'

        $('#convenientContent').html(convenientContent)

        $(".valueTextArea").focus(function () {
            $(this).select()
        })

        $('div.variables input').keyup(instantiateTemplate).change(instantiateTemplate).blur(function () {
            refreshValue(' ')
        })
        itemButtonClickable()
    }

    function itemClick() {
        $('.convenientConfigItem').click(function () {
            $('.convenientConfigItem').removeClass('convenientConfigItemSelected')
            var $this = $(this)
            $this.addClass('convenientConfigItemSelected')

            if ($this.text() == 'New Item') {
                newItem()
            } else {
                showExistedItem($this)
            }
        }).first().click()
    }


    var instantiateTemplate = function () {
        var templateValue = $('span.templateValue').text()
        $('div.variables').each(function (index, div) {
            var $div = $(div)
            var variableName = $div.find('span:first').text()
            variableName = variableName.substring(0, variableName.length - 1)
            var variableValue = $div.find('input').val()
            templateValue = templateValue.replace("{" + variableName + "}", variableValue)
        })

        $('.keyCreated').text(templateValue).unbind('click').click(function () {
            $.showContentAjax(templateValue)
        })
    }
    var refreshValue = function (resultTip) {
        clearConvenientContentInfo()
        var keyCreated = $('.keyCreated').text()
        $.ajax({
            type: 'GET', url: contextPath + "/showContent",
            data: {server: $('#servers').val(), database: $('#databases').val(), key: keyCreated},
            success: function (result, textStatus, request) {
                $(".valueTextArea").val(result.Exists ? result.Content : "(key does not exist)").select()
                setConvenientContentInfo(result.Exists && result.Ttl)
                $('.resultSpan').text(resultTip || 'Refreshed OK')
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert("11." + jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    var ttlInterval = null
    var clearConvenientContentInfo = function () {
        clearInterval(ttlInterval)
        ttlInterval = null
        $('#convenientContent').find('.info').text('')
        $('#convenientContent').find('.valueTextArea').val('(key does not exist)').select()
        $('.resultSpan').text('')
    }

    var resetConvenientContentInfo = function () {
        clearInterval(ttlInterval)
        ttlInterval = null
    }

    var setConvenientContentInfo = function (info) {
        resetConvenientContentInfo()
        var infoSpan = $('#convenientContent').find('.info')
        infoSpan.text(info ? '(' + info + ')' : '')
        if (info) {
            var seconds = parseDuration(info)
            if (seconds > 0) {
                ttlInterval = setInterval(function () {
                    seconds -= 1
                    if (seconds > 0) {
                        infoSpan.text('(' + createDuration(seconds) + ')')
                    } else {
                        clearConvenientContentInfo()
                    }
                }, 1000)
            }
        }
    }
})