'use strict';
const vscode = require("vscode");
const helper_1 = require("./helper");
const formatter_1 = require("./formatter");
const editor_1 = require("./editor");
// Main
function activate(context) {
    let tableFormatter = new formatter_1.TableFormatter();
    let tableEditor = new editor_1.TableEditor();
    let tableHelper = new helper_1.TableHelper();
    // --------------------------------
    // カーソル位置のテーブルのフォーマット
    // --------------------------------
    let formatCommand = vscode.commands.registerTextEditorCommand('extension.table.formatCurrent', (editor, edit) => {
        var pos = editor.selection.active;
        // 範囲の取得（Normal）
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
    // --------------------------------
    // 全テーブルのフォーマット
    // --------------------------------
    let formatAllCommand = vscode.commands.registerTextEditorCommand('extension.table.formatAll', (editor, edit) => {
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
            // フォーマット対象の範囲を取得
            var targetRange = tableHelper.getTargetRange(line, checkedIndex, rangeLines, editor.document.lineCount);
            checkedIndex = targetRange.checkedIndex;
            // 範囲の取得（Simple）
            var range = tableHelper.getTableRange(editor.document, line, helper_1.TableFormatType.Simple, targetRange.min, targetRange.count);
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
    // --------------------------------
    // 現在カーソル位置のCSVのフォーマット
    // --------------------------------
    let formatCsvCommand = vscode.commands.registerTextEditorCommand('extension.table.formatCurrentCsv', (editor, edit) => {
        var pos = editor.selection.active;
        // 範囲の取得（Normal）
        var range = tableHelper.getTableRange(editor.document, pos.line, helper_1.TableFormatType.Csv);
        // フォーマット（Normal）
        if (!range.isEmpty) {
            var info = tableHelper.getTableInfo(editor.document, range, helper_1.TableFormatType.Csv);
            var formatted = tableFormatter.getFormatTableText(editor.document, info, helper_1.TableFormatType.Csv);
            edit.replace(info.range, formatted);
            console.log("Table: Formatting succeeded!", "start: " + info.range.start.line, "end: " + info.range.end.line, "row: " + info.size.row, "col: " + (info.size.col - 1));
        }
    });
    // --------------------------------
    // 全CSVのフォーマット
    // --------------------------------
    let formatAllCsvCommand = vscode.commands.registerTextEditorCommand('extension.table.formatCurrentAllCsv', (editor, edit) => {
        var normalNum = 0;
        for (var i = 0; i < editor.document.lineCount; i++) {
            // 範囲の取得（Normal）
            var range = tableHelper.getTableRange(editor.document, i, helper_1.TableFormatType.Csv);
            // フォーマット（Normal）
            if (!range.isEmpty) {
                var info = tableHelper.getTableInfo(editor.document, range, helper_1.TableFormatType.Csv);
                var formatted = tableFormatter.getFormatTableText(editor.document, info, helper_1.TableFormatType.Csv);
                edit.replace(info.range, formatted);
                normalNum++;
                i = info.range.end.line + 1;
            }
        }
        if (normalNum > 0) {
            console.log("Table: Formatting succeeded!", "total: " + normalNum);
        }
    });
    context.subscriptions.push(tableFormatter);
    context.subscriptions.push(tableEditor);
    context.subscriptions.push(tableHelper);
    context.subscriptions.push(formatCommand);
    context.subscriptions.push(formatAllCommand);
    context.subscriptions.push(formatCsvCommand);
    context.subscriptions.push(formatAllCsvCommand);
}
exports.activate = activate;
function deactivate() {
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map