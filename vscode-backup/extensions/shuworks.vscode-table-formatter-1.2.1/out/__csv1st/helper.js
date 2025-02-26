'use strict';
const vscode = require("vscode");
const table_1 = require("./table");
var strWidth = require('string-width');
var trim = require('trim');
// Table format type
var TableFormatType;
(function (TableFormatType) {
    // separate with pipe table
    TableFormatType[TableFormatType["Normal"] = 0] = "Normal";
    // rest simple table
    TableFormatType[TableFormatType["Simple"] = 1] = "Simple";
    // CSV
    TableFormatType[TableFormatType["Csv"] = 2] = "Csv";
})(TableFormatType = exports.TableFormatType || (exports.TableFormatType = {}));
;
// Table line
var TableLineFlag;
(function (TableLineFlag) {
    TableLineFlag[TableLineFlag["None"] = 0] = "None";
    TableLineFlag[TableLineFlag["HasPipe"] = 1] = "HasPipe";
    TableLineFlag[TableLineFlag["PlusSeparator"] = 2] = "PlusSeparator";
    TableLineFlag[TableLineFlag["SimpleSeparator"] = 4] = "SimpleSeparator";
    TableLineFlag[TableLineFlag["Comma"] = 8] = "Comma";
    TableLineFlag[TableLineFlag["NotEmpty"] = 16] = "NotEmpty";
})(TableLineFlag = exports.TableLineFlag || (exports.TableLineFlag = {}));
;
class TableHelper {
    constructor() {
    }
    dispose() {
    }
    // テーブル記法の行か
    isTableLine(line, flag) {
        if (line.isEmptyOrWhitespace)
            return false;
        if (flag & TableLineFlag.HasPipe) {
            // 行が"|"を含む
            if (/\|/.test(line.text))
                return true;
        }
        if (flag & TableLineFlag.PlusSeparator) {
            // 行が+-=のみで構成されている
            if (/^(?=.*?\+)[\-=+]+$/.test(line.text))
                return true;
        }
        if (flag & TableLineFlag.SimpleSeparator) {
            // 行が-のみまたは=のみで構成されている
            if (/^[\- ]+$|^[= ]+$/.test(line.text))
                return true;
        }
        if (flag & TableLineFlag.Comma) {
            // 行が,を含む
            if (/,/.test(line.text))
                return true;
        }
        if (flag & TableLineFlag.NotEmpty) {
            // 行が空でない
            if (!line.isEmptyOrWhitespace)
                return true;
        }
        return false;
    }
    // Tableの範囲取得
    getTableRange(doc, line, formatType, minLine = 0, maxCount = -1) {
        var startLine = line;
        var endLine = line;
        if (maxCount < 0)
            maxCount = doc.lineCount - minLine;
        switch (formatType) {
            // ----------------
            case TableFormatType.Normal:
                // 現在の行を判定
                if (this.isTableLine(doc.lineAt(line), TableLineFlag.HasPipe | TableLineFlag.PlusSeparator)) {
                    // 後方に操作し開始行を取得
                    for (var i = line - 1; i >= minLine; i--) {
                        if (!this.isTableLine(doc.lineAt(i), TableLineFlag.HasPipe | TableLineFlag.PlusSeparator))
                            break;
                        startLine = i;
                    }
                    // 前方に操作し終了行を取得
                    for (var i = line + 1; i < minLine + maxCount; i++) {
                        if (!this.isTableLine(doc.lineAt(i), TableLineFlag.HasPipe | TableLineFlag.PlusSeparator))
                            break;
                        endLine = i;
                    }
                }
                break;
            // ----------------
            case TableFormatType.Simple:
                var hasSeparator = false;
                // 現在の行を判定
                if (this.isTableLine(doc.lineAt(line), TableLineFlag.NotEmpty)) {
                    // 後方に操作し開始行を取得（空白行またドキュメント始まり、セパレーター行の構成を探す）
                    for (var i = line - 1; i >= minLine - 1; i--) {
                        if (i == minLine - 1 || !this.isTableLine(doc.lineAt(i), TableLineFlag.NotEmpty)) {
                            if (i + 1 < minLine + maxCount && this.isTableLine(doc.lineAt(i + 1), TableLineFlag.SimpleSeparator)) {
                                startLine = i + 1;
                                hasSeparator = true;
                            }
                            break;
                        }
                    }
                    // 前方に操作し終了行を取得（セパレーター行、空白行またはドキュメント終わりの構成を探す）
                    if (hasSeparator) {
                        hasSeparator = false;
                        for (var i = line + 1; i < minLine + maxCount + 1; i++) {
                            if (i == minLine + maxCount || !this.isTableLine(doc.lineAt(i), TableLineFlag.NotEmpty)) {
                                if (i - 1 >= minLine && this.isTableLine(doc.lineAt(i - 1), TableLineFlag.SimpleSeparator)) {
                                    endLine = i - 1;
                                    hasSeparator = true;
                                }
                                break;
                            }
                        }
                    }
                }
                // 正しいセパレーター行がない場合は初期化
                if (!hasSeparator) {
                    startLine = line;
                    endLine = line;
                }
                break;
            // ----------------
            case TableFormatType.Csv:
                // 現在の行を判定
                if (this.isTableLine(doc.lineAt(line), TableLineFlag.Comma)) {
                    // 後方に操作し開始行を取得
                    for (var i = line - 1; i >= minLine; i--) {
                        if (!this.isTableLine(doc.lineAt(i), TableLineFlag.Comma))
                            break;
                        startLine = i;
                    }
                    // 前方に操作し終了行を取得
                    for (var i = line + 1; i < minLine + maxCount; i++) {
                        if (!this.isTableLine(doc.lineAt(i), TableLineFlag.Comma))
                            break;
                        endLine = i;
                    }
                }
                break;
        }
        // 複数列にヒットしない場合はisEmptyに引っかかるように0にする
        var endChar = (startLine == endLine) ? 0 : doc.lineAt(endLine).range.end.character;
        return new vscode.Range(new vscode.Position(startLine, 0), new vscode.Position(endLine, endChar));
    }
    // 行の文字列を分割する（必ず先頭が空白で末尾は空白でないようにして返す）
    getSplitLineText(text, formatType) {
        var cells = [];
        var delimiter = table_1.DelimiterType.Pipe;
        switch (formatType) {
            case TableFormatType.Normal:
                // |がないときのみ+で分ける（末尾の空白も含めるため-1）
                if (text.indexOf('|') != -1) {
                    cells = text.split("|", -1);
                    delimiter = table_1.DelimiterType.Pipe;
                }
                else {
                    cells = text.split("+", -1);
                    delimiter = table_1.DelimiterType.Plus;
                }
                break;
            case TableFormatType.Simple:
                // @TODO: csv-parserを使って区切る
                // 空白で区切る（""か''で囲まれた範囲の空白は無視する）
                cells = text.match(/('[^']*'|"[^"]*")|[^ ]+/g);
                delimiter = table_1.DelimiterType.Space;
                break;
            case TableFormatType.Csv:
                // @TODO: csv-parserを使って区切る
                // text = text.replace()
                // ,で区切る（""か''で囲まれた範囲の空白は無視する）
                cells = text.match(/('[^']*'|"[^"]*")|[^,]+/g);
                delimiter = table_1.DelimiterType.Camma;
                break;
        }
        // 先頭に空白の追加
        if (cells.length >= 1 && trim(cells[0]) != "") {
            cells.unshift("");
        }
        // 末尾の空白の削除
        if (cells.length >= 1 && trim(cells[cells.length - 1]) == "") {
            cells.pop();
        }
        return { cells: cells, delimiter: delimiter };
    }
    // 行の解析
    getCellInfoList(line, formatType) {
        if (line.isEmptyOrWhitespace)
            return [];
        var splitText = this.getSplitLineText(line.text, formatType);
        var list = [];
        var cells = splitText.cells;
        var delimiter = splitText.delimiter;
        // 先頭に空白列を追加
        list.push(new table_1.CellInfo(line.firstNonWhitespaceCharacterIndex, delimiter));
        for (var i = 0; i < cells.length; i++) {
            var trimed = trim(cells[i]);
            // 先頭は空白で追加済みなので無視する
            if (i == 0)
                continue;
            var size = strWidth(trimed);
            var type = (size == 0) ? table_1.CellType.CM_Blank : table_1.CellType.CM_Content;
            var align = table_1.CellAlign.Left;
            switch (formatType) {
                case TableFormatType.Normal:
                    // Common  ----------------
                    if (/^-+$/.test(trimed)) {
                        type = table_1.CellType.CM_MinusSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    else if (/^=+$/.test(trimed)) {
                        type = table_1.CellType.CM_EquallSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    else if (/^:-+$/.test(trimed)) {
                        type = table_1.CellType.MD_LeftSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    else if (/^-+:$/.test(trimed)) {
                        type = table_1.CellType.MD_RightSeparator;
                        align = table_1.CellAlign.Right;
                    }
                    else if (/^:-+:$/.test(trimed)) {
                        type = table_1.CellType.MD_CenterSeparator;
                        align = table_1.CellAlign.Center;
                    }
                    else if (/^_\./.test(trimed)) {
                        size = strWidth(trim(trimed.substring(2)));
                        type = table_1.CellType.TT_HeaderPrefix;
                        align = table_1.CellAlign.Left;
                    }
                    else if (/^<\./.test(trimed)) {
                        size = strWidth(trim(trimed.substring(2)));
                        type = table_1.CellType.TT_LeftPrefix;
                        align = table_1.CellAlign.Left;
                    }
                    else if (/^>\./.test(trimed)) {
                        size = strWidth(trim(trimed.substring(2)));
                        type = table_1.CellType.TT_RightPrefix;
                        align = table_1.CellAlign.Right;
                    }
                    else if (/^=\./.test(trimed)) {
                        size = strWidth(trim(trimed.substring(2)));
                        type = table_1.CellType.TT_CenterPrefix;
                        align = table_1.CellAlign.Center;
                    }
                    break;
                case TableFormatType.Simple:
                    // Common  ----------------
                    if (/^-+$/.test(trimed)) {
                        type = table_1.CellType.CM_MinusSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    else if (/^=+$/.test(trimed)) {
                        type = table_1.CellType.CM_EquallSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    break;
                case TableFormatType.Csv:
                    // NOP ----------------
                    break;
            }
            list.push(new table_1.CellInfo(size, delimiter, type, align));
        }
        return list;
    }
    // 表データの取得
    getTableInfo(doc, range, formatType) {
        // 各行の解析
        var grid = [];
        for (var i = range.start.line; i <= range.end.line; i++) {
            grid.push(this.getCellInfoList(doc.lineAt(i), formatType));
        }
        return new table_1.TableInfo(range, grid);
    }
    // フォーマット済み範囲を走査しフォーマット対象の範囲を決める
    getTargetRange(line, checkedIndex, ignoreRangeLines, maxLineCount) {
        // 判定に引っかからなかった場合のために初期値として最終区画の範囲を設定しておく
        var min = (ignoreRangeLines.length > 0) ? ignoreRangeLines[ignoreRangeLines.length - 1] + 1 : 0;
        var cnt = maxLineCount - min;
        for (var j = checkedIndex; j < ignoreRangeLines.length; j++) {
            // 行が通り越したら
            if (line < ignoreRangeLines[j]) {
                // 偶数行（start）
                if (j % 2 == 0) {
                    // フォーマット済み範囲間なので範囲を設定
                    min = (j == 0) ? 0 : ignoreRangeLines[j - 1] + 1;
                    cnt = ignoreRangeLines[j] - min;
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
        return { min: min, count: cnt, checkedIndex: checkedIndex };
    }
}
exports.TableHelper = TableHelper;
//# sourceMappingURL=helper.js.map