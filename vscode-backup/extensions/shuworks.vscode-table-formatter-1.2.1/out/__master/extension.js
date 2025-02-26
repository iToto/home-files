'use strict';
const vscode = require("vscode");
const helper_1 = require("./helper");
const formatter_1 = require("./formatter");
const editor_1 = require("./editor");
// Main
function activate(context) {
    let isInitilized = false;
    let configTitle = "tableformatter";
    let settings = {
        markdown: {
            oneSpacePadding: true
        },
        common: {
            explicitFullwidthChars: []
        }
    };
    let tableHelper = new helper_1.TableHelper(settings);
    let tableFormatter = new formatter_1.TableFormatter(settings, tableHelper);
    let tableEditor = new editor_1.TableEditor();
    function initialize(config) {
        settings.markdown.oneSpacePadding = config.get('markdown.oneSpacePadding', true);
        let chars = config.get('common.explicitFullwidthChars', []).filter(function (elem, i, self) {
            return ((self.indexOf(elem) === i) && (elem.length == 1));
        });
        settings.common.explicitFullwidthChars = [];
        chars.forEach((char, i) => {
            settings.common.explicitFullwidthChars.push(new RegExp(char, 'g'));
        });
        isInitilized = true;
    }
    vscode.workspace.onDidChangeConfiguration(() => {
        initialize(vscode.workspace.getConfiguration(configTitle));
    }, null, context.subscriptions);
    initialize(vscode.workspace.getConfiguration(configTitle));
    let formatCommand = vscode.commands.registerTextEditorCommand('extension.table.formatCurrent', (editor, edit) => {
        // 初期化
        if (!isInitilized)
            initialize(vscode.workspace.getConfiguration(configTitle));
        // 範囲の取得（Normal）
        var pos = editor.selection.active;
        var range = tableHelper.getTableRange(editor.document, pos.line, helper_1.TableFormatType.Normal);
        // フォーマット（Normal）
        if (!range.isEmpty) {
            var info = tableHelper.getTableInfo(editor.document, range, helper_1.TableFormatType.Normal);
            var formatted = tableFormatter.getFormatTableText(editor.document, info, helper_1.TableFormatType.Normal);
            edit.replace(info.range, formatted);
            console.log("Table: Formatting succeeded!", "start: " + info.range.start.line, "end: " + info.range.end.line, "row: " + info.size.row, "col: " + (info.size.col - 1));
        }
        else {
            // 範囲の取得（Simple）
            var range = tableHelper.getTableRange(editor.document, pos.line, helper_1.TableFormatType.Simple);
            // フォーマット（Simple）
            if (!range.isEmpty) {
                var info = tableHelper.getTableInfo(editor.document, range, helper_1.TableFormatType.Simple);
                var formatted = tableFormatter.getFormatTableText(editor.document, info, helper_1.TableFormatType.Simple);
                edit.replace(info.range, formatted);
                console.log("Table: Formatting simple succeeded!", "start: " + info.range.start.line, "end: " + info.range.end.line, "row: " + info.size.row, "col: " + (info.size.col - 1));
            }
        }
    });
    let formatAllCommand = vscode.commands.registerTextEditorCommand('extension.table.formatAll', (editor, edit) => {
        // 初期化
        if (!isInitilized)
            initialize(vscode.workspace.getConfiguration(configTitle));
        var normalNum = 0;
        var simpleNum = 0;
        var targetLines = [];
        var rangeLines = [];
        for (var i = 0; i < editor.document.lineCount; i++) {
            // 範囲の取得（Normal）
            var range = tableHelper.getTableRange(editor.document, i, helper_1.TableFormatType.Normal);
            // フォーマット（Normal）
            if (!range.isEmpty) {
                var info = tableHelper.getTableInfo(editor.document, range, helper_1.TableFormatType.Normal);
                var formatted = tableFormatter.getFormatTableText(editor.document, info, helper_1.TableFormatType.Normal);
                edit.replace(info.range, formatted);
                // フォーマット済み範囲を積んでおく（偶数がstart、奇数がend）
                rangeLines.push(range.start.line);
                rangeLines.push(range.end.line);
                normalNum++;
                i = info.range.end.line + 1;
            }
            else {
                // SimpleTableの候補行を積んでおく
                if (tableHelper.isTableLine(editor.document.lineAt(i), helper_1.TableLineFlag.SimpleSeparator)) {
                    targetLines.push(i);
                }
            }
        }
        var endLine = -1;
        var checkedIndex = 0;
        for (var i = 0; i < targetLines.length; i++) {
            var line = targetLines[i];
            if (line <= endLine)
                continue;
            // フォーマット済み範囲を走査しフォーマット対象の範囲を決める
            var min = (rangeLines.length > 0) ? rangeLines[rangeLines.length - 1] + 1 : 0;
            var cnt = editor.document.lineCount - min;
            for (var j = checkedIndex; j < rangeLines.length; j++) {
                // 行が通り越したら
                if (line < rangeLines[j]) {
                    // 偶数行（start）
                    if (j % 2 == 0) {
                        // フォーマット済み範囲間なので範囲を設定
                        min = (j == 0) ? 0 : rangeLines[j - 1] + 1;
                        cnt = rangeLines[j] - min;
                        checkedIndex = j;
                    }
                    else {
                        // フォーマット済み範囲内なので無視する
                        min = 0;
                        cnt = 0;
                    }
                    break;
                }
            }
            // 範囲の取得（Simple）
            var range = tableHelper.getTableRange(editor.document, line, helper_1.TableFormatType.Simple, min, cnt);
            // フォーマット（Simple）
            if (!range.isEmpty) {
                var info = tableHelper.getTableInfo(editor.document, range, helper_1.TableFormatType.Simple);
                var formatted = tableFormatter.getFormatTableText(editor.document, info, helper_1.TableFormatType.Simple);
                edit.replace(info.range, formatted);
                simpleNum++;
                endLine = info.range.end.line;
            }
        }
        if (normalNum + simpleNum > 0) {
            console.log("Table: Formatting succeeded!", "total: " + (normalNum + simpleNum), "(normal: " + normalNum, "simple: " + simpleNum + ")");
        }
    });
    context.subscriptions.push(tableFormatter);
    context.subscriptions.push(tableEditor);
    context.subscriptions.push(tableHelper);
    context.subscriptions.push(formatCommand);
    context.subscriptions.push(formatAllCommand);
}
exports.activate = activate;
function deactivate() {
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map