$(function () {
    var pathname = window.location.pathname
    if (pathname.lastIndexOf("/", pathname.length - 1) !== -1) {
        pathname = pathname.substring(0, pathname.length - 1)
    }

    function refreshKeys(key) {
        $.ajax({
            type: 'GET', url: pathname + "/listKeys",
            data: {server: $('#servers').val(), database: $('#databases').val(), pattern: $('#serverFilterKeys').val()},
            success: function (content, textStatus, request) {
                showKeysTree(content)
                if (key) {
                    chosenKey(key)
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    $('#serverFilterKeysBtn,#refreshKeys').click(function () {
        refreshKeys()
    })

    refreshKeys()

    $('#serverFilterKeys').keydown(function (event) {
        var keyCode = event.keyCode || event.which
        if (keyCode == 13) {
            refreshKeys()
        }
    })

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
                type: 'GET', url: pathname + "/redisImport",
                data: {server: $('#servers').val(), database: $('#databases').val(), commands: commands},
                success: function (content, textStatus, request) {
                    $('#importResult').html('<pre>' + content + '</pre>')
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })
        })
    })

    $('#checkAllChk').click(function () {
        var checked = $('#checkAllChk').is(":checked")
        $('#keys ul li:visible').each(function (index, li) {
            $(li).find('input:checkbox').prop('checked', checked)
        })
    })

    $('#exportKeys').click(function () {
        var contentHtml = '<div><h3>Export</h3></div>' +
            '<div><span style="margin-right: 30px">Type:</span>' +
            '<select style="margin-right: 30px" id="exportType"><option value="Redis">Redis</option><option value="JSON">JSON</option></select>' +
            '<button id="exportBtn">Export</button><button id="downloadExportBtn">Download</button></div>' +
            '<div id="exportResult"></div>'

        $('#frame').html(contentHtml)

        $('#exportBtn,#downloadExportBtn').click(function () {
            var keys = findCheckedKeys()
            if (keys.length == 0) {
                alert("No keys chosen to be deleted!")
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
                window.open(pathname + "/exportKeys?" + $.param(data), '_blank')
                return
            }

            $.ajax({
                type: 'GET', url: pathname + "/exportKeys",
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
                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })

        })
    })

    function findCheckedKeys() {
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
        var keys = findCheckedKeys()
        if (keys.length == 0) {
            alert("No keys chosen to be deleted!")
            return
        }

        if (!confirm("Are you sure to delete " + keys.length + " keys?")) {
            return
        }

        $.ajax({
            type: 'POST', url: pathname + "/deleteMultiKeys",
            data: {server: $('#servers').val(), database: $('#databases').val(), keys: JSON.stringify(keys)},
            success: function (content, textStatus, request) {
                if (content != 'OK') {
                    alert(content)
                } else {
                    removeKeys(keys)
                    $('#checkAllChk').prop('checked', false)
                    $('#frame').html('<div><span class="key">' + keys.length + ' keys were deleted:</span></div>'
                        + '<div><br>' + joinKeysWithNo(keys) + '</div>')
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    })

    function executeRedisCmd() {
        var cmd = $('#directCmd').text()
        var server = $('#servers').val()
        $.ajax({
            type: 'POST', url: pathname + "/redisCli",
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
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
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

    var convenientConfig = null

    $('#convenientSpan').click(function () {
        $.ajax({
            type: 'POST', url: pathname + "/convenientConfig",
            success: function (result, textStatus, request) {
                if (result.Ready !== true) {
                    $('#frame').html('<div><span class="key">' + result.Error + '</span></div>')
                    return
                }

                convenientConfig = result
                var contentHtml = '<div class="itemNames">'
                for (var i = 0, len = convenientConfig.Items.length; i < len; i++) {
                    convenientConfigItem = convenientConfig.Items[i]
                    contentHtml += '<span itemIndex="' + i + '" class="convenientConfigItem">' + convenientConfigItem.Name + '</span>'
                }
                contentHtml += '<span class="convenientConfigItem">New Item</span>'
                contentHtml += '</div><div id="convenientContent"></div>'
                $('#frame').html(contentHtml)

                $('.convenientConfigItem').click(function () {
                    $('.convenientConfigItem').removeClass('convenientConfigItemSelected')
                    var $this = $(this)
                    $this.addClass('convenientConfigItemSelected')

                    var $convenientContent = $('#convenientContent');
                    if ($this.text() == 'New Item') {
                        var convenientContent =
                            '<div><span>Name:</span><span><input class="templateName"></span></div>' +
                            '<div><span>Key Template:</span><span><input class="templateValue"></span></div>' +
                            '<div><span>TTL:</span><span><input class="ttlCreated" value="-1s"></span></div>' +
                            '<div><span>Operations:</span><span><input style="width:12px" type="checkbox" id="DeleteChk" value="Delete" checked><label for="DeleteChk">Delete</label>' +
                            '<input  style="width:12px" type="checkbox" id="SaveChk" value="Save" checked><label for="SaveChk">Save</label></span></div>' +
                            '<div><span>&nbsp;</span><span><button>Save New Item</button></span></div>'
                        $convenientContent.html(convenientContent)
                        $convenientContent.find('button').click(function () {
                            var operations = ''
                            $convenientContent.find('input:checked').each(function (index, chk) {
                                if (operations.length > 0) operations += ','
                                operations += $(chk).val()
                            })
                            $.ajax({
                                type: 'GET', url: pathname + "/convenientConfigAdd",
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
                                        alert(result.Message)
                                    }
                                },
                                error: function (jqXHR, textStatus, errorThrown) {
                                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                                }
                            })
                        })

                        return
                    }

                    var item = convenientConfig.Items[+$this.attr('itemIndex')]
                    var convenientContent =
                        '<div><span>Key Template:</span><span class="templateValue">' + item.Template + ' </span></div>'

                    var variables = parseTemplateVariables(item.Template)
                    for (var i = 0, len = variables.length; i < len; i++) {
                        convenientContent += '<div class="variables"><span>' + variables[i] + ':</span><span><input placeholder="variable value"></span></div>'
                    }

                    convenientContent += '<div><span>Key:</span><span class="keyCreated"></span></div>'

                    convenientContent += '<div><span>TTL:</span><span class="ttlCreated"><input value="' + item.Ttl + '"></span></div>'
                    convenientContent += '<div><span>Value:<br/><span class="info"></span></span><span><textarea class="valueTextArea"></textarea></span></div>'

                    convenientContent += '<div><span>Operations:</span><span><button class="convenientButton">Refresh Value</button>'
                    for (var i = 0, len = item.Operations.length; i < len; i++) {
                        convenientContent += '<button class="convenientButton">' + capitalize(item.Operations[i]) + "</button>"
                    }
                    convenientContent += '<button class="convenientButton" sectionName="' + item.Section + '">Delete This Item</button>'
                    convenientContent += '</span></div>'

                    convenientContent += '<div><span>Result:</span><span class="resultSpan"></span></div>'

                    $convenientContent.html(convenientContent)

                    $(".valueTextArea").focus(function () {
                        $(this).select()
                    })

                    var instantiateTemplate = function () {
                        var templateValue = $('span.templateValue').text()
                        $('div.variables').each(function (index, div) {
                            var $div = $(div)
                            var variableName = $div.find('span:first').text()
                            variableName = variableName.substring(0, variableName.length - 1)
                            var variableValue = $div.find('input').val()
                            templateValue = templateValue.replace("{" + variableName + "}", variableValue)
                        })

                        $('.keyCreated').text(templateValue).unbind('click').click(function() {
                            showContentAjax(templateValue)
                        })
                    }
                    var refreshValue = function (resultTip) {
                        clearConvenientContentInfo()
                        var keyCreated = $('.keyCreated').text()
                        $.ajax({
                            type: 'GET', url: pathname + "/showContent",
                            data: {server: $('#servers').val(), database: $('#databases').val(), key: keyCreated},
                            success: function (result, textStatus, request) {
                                $(".valueTextArea").val(result.Exists ? result.Content : "(key does not exist)").select()
                                setConvenientContentInfo(result.Exists && result.Ttl)
                                $('.resultSpan').text(resultTip || 'Refreshed OK')
                            },
                            error: function (jqXHR, textStatus, errorThrown) {
                                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
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

                    $('div.variables input').keyup(instantiateTemplate).change(instantiateTemplate).blur(function () {
                        refreshValue(' ')
                    })

                    $('.convenientButton').click(function () {
                        var $this = $(this)
                        var keyCreated = $('.keyCreated').text()
                        if ($this.text() == 'Save') {
                            var ttl = $('.ttlCreated input').val()
                            $.ajax({
                                type: 'POST', url: pathname + "/changeContent",
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
                                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                                }
                            })
                        } else if ($this.text() == 'Delete') {
                            $.ajax({
                                type: 'POST', url: pathname + "/deleteKey",
                                data: {server: $('#servers').val(), database: $('#databases').val(), key: keyCreated},
                                success: function (content, textStatus, request) {
                                    if (content == 'OK') {
                                        refreshValue('Deleted OK')
                                    } else {
                                        $('.resultSpan').text('Deleted ' + content)
                                    }
                                },
                                error: function (jqXHR, textStatus, errorThrown) {
                                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                                }
                            })
                        } else if ($this.text() == 'Refresh Value') {
                            refreshValue()
                        } else if ($this.text() == 'Delete This Item') {
                            $.ajax({
                                type: 'POST', url: pathname + "/deleteConvenientConfigItem",
                                data: {
                                    sectionName: $this.attr('sectionName')
                                },
                                success: function (content, textStatus, request) {
                                    $('#convenientSpan').click()
                                },
                                error: function (jqXHR, textStatus, errorThrown) {
                                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                                }
                            })
                        }
                    })
                }).first().click()
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }).click()


    $('#redisInfo').click(function () {
        $.ajax({
            type: 'GET',
            url: pathname + "/redisInfo",
            data: {server: $('#servers').val(), database: $('#databases').val()},
            success: function (result, textStatus, request) {
                var contentHtml = '<div><span class="key">Redis info</span></div>' +
                    '<pre>' + result + '</pre>'

                $('#frame').html(contentHtml)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    })

    function showContentAjax(key) {
        $.ajax({
            type: 'GET', url: pathname + "/showContent",
            data: {server: $('#servers').val(), database: $('#databases').val(), key: key},
            success: function (result, textStatus, request) {
                showContent(key, result.Type, result.Content, result.Ttl, result.Size, result.Encoding, result.Error, result.Exists, result.Format)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function showKeysTree(keysArray) {
        $('#keysNum').html('(' + keysArray.length + ')')

        var keysHtml = '<ul>'
        for (var i = 0; i < keysArray.length; ++i) {
            var key = keysArray[i]
            var nodeCss = i < keysArray.length - 1 ? "sprite-tree-node" : "sprite-tree-lastnode last"
            keysHtml += '<li class="datatype-' + key.Type + ' sprite ' + nodeCss + '" data-type="' + key.Type + '">' +
                '<input type="checkbox">' +
                '<span class="sprite sprite-datatype-' + key.Type + '"></span>' +
                '<span class="keyValue">' + key.Key + '</span>'

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
            showContentAjax(key)
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

        $('#sidebar').height(window.outerHeight)
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

    function removeKeys(keys) {
        $('#keys ul li').removeClass('chosen').each(function (index, li) {
            var $span = $(li).find('.keyValue')
            if ($.inArray($span.text(), keys) > -1) {
                $(li).remove()
            }
        })
    }

    function removeKey(key) {
        $('#keys ul li').removeClass('chosen').each(function (index, li) {
            var $span = $(li).find('.keyValue')
            if ($span.text() == key) {
                $(li).remove()
                return false
            }
        })
    }


    $('#servers,#databases').change(refreshKeys)


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

    function extractValue(type) {
        var value = null
        if (type == 'string') {
            value = codeMirror != null && codeMirror.getValue() || $('#code').val()
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

        $('#format').change(function () {
            codeMirror = null
            if ($(this).val() == 'JSON' && $('#type').val() == 'string') {
                codeMirror = CodeMirror.fromTextArea(document.getElementById('code'), {
                    mode: 'application/json', lineNumbers: true, matchBrackets: true
                })
            }
        })

        $('.valueSave').click(function () {
            var type = $('#type').val()
            var key = $('#key').val()
            var ttl = $('#ttl').val()
            var format = $('#format').val()
            var jsonValue = extractValue(type)

            if (confirm("Are you sure to save save for " + key + "?")) {
                $.ajax({
                    type: 'POST', url: pathname + "/newKey",
                    data: {
                        server: $('#servers').val(), database: $('#databases').val(),
                        type: type, key: key, ttl: ttl, format: format, value: jsonValue
                    },
                    success: function (content, textStatus, request) {
                        if (content == 'OK') {
                            refreshKeys(key)
                            showContentAjax(key)
                        } else {
                            alert(content)
                        }
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                        alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                    }
                })
            }
        })
    })

    var codeMirror = null

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

        codeMirror = null
        if (format === "JSON" && size < 5000) {
            codeMirror = CodeMirror.fromTextArea(document.getElementById('code'), {
                mode: 'application/json', lineNumbers: true, matchBrackets: true
            })
        } else {
            autosize($('#code'))
        }

        $('.keyDelete').click(function () {
            if (confirm("Are you sure to delete " + key + "?")) {
                $.ajax({
                    type: 'POST', url: pathname + "/deleteKey",
                    data: {server: $('#servers').val(), database: $('#databases').val(), key: key},
                    success: function (content, textStatus, request) {
                        if (content != 'OK') {
                            alert(content)
                        } else {
                            removeKey(key)
                            $('#frame').html('<div><span class="key">' + key + ' was deleted</span></div>')
                        }
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                        alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                    }
                })
            }
        })

        $('.valueSave').click(function () {
            if (confirm("Are you sure to save changes for " + key + "?")) {
                var changedContent = extractValue(type)
                $.ajax({
                    type: 'POST', url: pathname + "/changeContent",
                    data: {
                        server: $('#servers').val(), database: $('#databases').val(),
                        key: key, type: type, ttl: $('#ttl').text(), value: changedContent
                    },
                    success: function (content, textStatus, request) {
                        if (content == 'OK') {
                            showContentAjax(key)
                        } else {
                            alert(content)
                        }
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                        alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                    }
                })
            }
        })
    }


    $('#maintainRedisServers').click(function (e) {
        $('#serversMaintain').modal()
    })

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

    $(document).on('paste', '[contenteditable]', function (e) {
        e.preventDefault();
        var text = '';
        if (e.clipboardData || e.originalEvent.clipboardData) {
            text = (e.originalEvent || e).clipboardData.getData('text/plain');
        } else if (window.clipboardData) {
            text = window.clipboardData.getData('Text');
        }
        if (document.queryCommandSupported('insertText')) {
            document.execCommand('insertText', false, text);
        } else {
            document.execCommand('paste', false, text);
        }
    });

})