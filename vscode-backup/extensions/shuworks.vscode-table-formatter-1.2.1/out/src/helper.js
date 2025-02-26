'use strict';
const vscode = require("vscode");
const table_1 = require("./table");
var strWidth = require('string-width');
var trim = require('trim');
;
// Table format type
var TableFormatType;
(function (TableFormatType) {
    // separate with pipe table
    TableFormatType[TableFormatType["Normal"] = 0] = "Normal";
    // rest simple table
    TableFormatType[TableFormatType["Simple"] = 1] = "Simple";
    // CSV
    // Csv
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
// ヘルパークラス
class TableHelper {
    constructor(config) {
        // オブジェクトは参照渡し
        this.settings = config;
        // キャッシュ
        this.regExpMap = {};
        // "|"を含む
        this.pipeRegExp = /\|/;
        // +-=のみで構成されている
        this.plusSepRegExp = /^(?=.*?\+)[\-=+]+$/;
        // -のみまたは=のみで構成されている
        this.simpleSepRegExp = /^[\- ]+$|^[= ]+$/;
        // ,を含む
        this.commaRegExp = /,/;
        // Common
        this.commonMinusRegExp = /^-+$/;
        this.commonEqualRegExp = /^=+$/;
        // Markdown
        this.markdownLeftRegExp = /^:-+$/;
        this.markdownRightRegExp = /^-+:$/;
        this.markdownCenterRegExp = /^:-+:$/;
        // Textile
        this.textileHeaderRegExp = /^_\./;
        this.textileLeftRegExp = /^<\./;
        this.textileRightRegExp = /^>\./;
        this.textileCenterRegExp = /^=\./;
    }
    dispose() {
    }
    // テーブル記法の行か
    isTableLine(line, flag) {
        if (line.isEmptyOrWhitespace)
            return false;
        if (flag & TableLineFlag.HasPipe) {
            // 行が"|"を含む
            if (this.pipeRegExp.test(line.text))
                return true;
        }
        if (flag & TableLineFlag.PlusSeparator) {
            // 行が+-=のみで構成されている
            if (this.plusSepRegExp.test(line.text))
                return true;
        }
        if (flag & TableLineFlag.SimpleSeparator) {
            // 行が-のみまたは=のみで構成されている
            if (this.simpleSepRegExp.test(line.text))
                return true;
        }
        if (flag & TableLineFlag.Comma) {
            // 行が,を含む
            if (this.commaRegExp.test(line.text))
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
        }
        // 複数列にヒットしない場合はisEmptyに引っかかるように0にする
        var endChar = (startLine == endLine) ? 0 : doc.lineAt(endLine).range.end.character;
        return new vscode.Range(new vscode.Position(startLine, 0), new vscode.Position(endLine, endChar));
    }
    // 行の文字列を分割する（必ず先頭が空白で末尾は空白でないようにして返す）
    getSplittedLineText(text, formatType) {
        var cells = [];
        var delimiter = table_1.DelimiterType.Pipe;
        switch (formatType) {
            case TableFormatType.Normal:
                // | で分割
                if (text.indexOf('|') != -1) {
                    // 末尾の空白も含めるため -1
                    cells = text.split("|", -1);
                    delimiter = table_1.DelimiterType.Pipe;
                }
                else {
                    // 末尾の空白も含めるため -1
                    cells = text.split("+", -1);
                    delimiter = table_1.DelimiterType.Plus;
                }
                break;
            case TableFormatType.Simple:
                cells = this.getSplittedTextByRegExp(text, " ");
                // @TODO: CSV対応したらSimpleの動作確認もする
                // console.log(JSON.stringify(cells));
                delimiter = table_1.DelimiterType.Space;
                break;
        }
        // 先頭要素が空白でない場合、空白要素の追加
        var isAdded = false;
        if (cells.length >= 1 && trim(cells[0]) != "") {
            cells.unshift("");
            isAdded = true;
        }
        // 末尾要素が空白の場合、削除
        if (cells.length >= 1 && trim(cells[cells.length - 1]) == "") {
            cells.pop();
        }
        return { cells: cells, delimiter: delimiter, isAddedBlankHead: isAdded };
    }
    // 正規表現を使って区切った文字列を返す
    getSplittedTextByRegExp(text, delimiter) {
        if (!this.regExpMap[delimiter]) {
            // デリミタで区切る（""か''で囲まれた範囲の空白は無視する）
            // @TODO: CSV対応
            this.regExpMap[delimiter] = new RegExp("('[^']*'|\"[^\"]*\")|[^" + delimiter + "]+", "g");
        }
        return text.match(this.regExpMap[delimiter]);
    }
    // 行の解析
    getCellInfoList(line, formatType) {
        if (line.isEmptyOrWhitespace)
            return { list: [], isAddedBlankHead: false };
        var obj = this.getSplittedLineText(line.text, formatType);
        var list = [];
        var cells = obj.cells;
        // 先頭は必ず空白になっており、もともとのオフセット分数の空白文字のセルを追加
        let spaces = Array(line.firstNonWhitespaceCharacterIndex + 1).join(" ");
        list.push(new table_1.CellInfo(this.settings, spaces, obj.delimiter));
        for (var i = 0; i < cells.length; i++) {
            var trimmed = trim(cells[i]);
            // 先頭の空白は追加済みなので無視する
            if (i == 0)
                continue;
            var type = (trimmed.length == 0) ? table_1.CellType.CM_Blank : table_1.CellType.CM_Content;
            var align = table_1.CellAlign.Left;
            switch (formatType) {
                case TableFormatType.Normal:
                    // Common  ----------------
                    if (this.commonMinusRegExp.test(trimmed)) {
                        type = table_1.CellType.CM_MinusSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    else if (this.commonEqualRegExp.test(trimmed)) {
                        type = table_1.CellType.CM_EquallSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    else if (this.markdownLeftRegExp.test(trimmed)) {
                        type = table_1.CellType.MD_LeftSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    else if (this.markdownRightRegExp.test(trimmed)) {
                        type = table_1.CellType.MD_RightSeparator;
                        align = table_1.CellAlign.Right;
                    }
                    else if (this.markdownCenterRegExp.test(trimmed)) {
                        type = table_1.CellType.MD_CenterSeparator;
                        align = table_1.CellAlign.Center;
                    }
                    else if (this.textileHeaderRegExp.test(trimmed)) {
                        // ._を削除
                        trimmed = trim(trimmed.substring(2));
                        type = table_1.CellType.TT_HeaderPrefix;
                        align = table_1.CellAlign.Left;
                    }
                    else if (this.textileLeftRegExp.test(trimmed)) {
                        // <_を削除
                        trimmed = trim(trimmed.substring(2));
                        type = table_1.CellType.TT_LeftPrefix;
                        align = table_1.CellAlign.Left;
                    }
                    else if (this.textileRightRegExp.test(trimmed)) {
                        // >_を削除
                        trimmed = trim(trimmed.substring(2));
                        type = table_1.CellType.TT_RightPrefix;
                        align = table_1.CellAlign.Right;
                    }
                    else if (this.textileCenterRegExp.test(trimmed)) {
                        // =_を削除
                        trimmed = trim(trimmed.substring(2));
                        type = table_1.CellType.TT_CenterPrefix;
                        align = table_1.CellAlign.Center;
                    }
                    break;
                case TableFormatType.Simple:
                    // Common  ----------------
                    if (this.commonMinusRegExp.test(trimmed)) {
                        type = table_1.CellType.CM_MinusSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    else if (this.commonEqualRegExp.test(trimmed)) {
                        type = table_1.CellType.CM_EquallSeparator;
                        align = table_1.CellAlign.Left;
                    }
                    break;
            }
            list.push(new table_1.CellInfo(this.settings, trimmed, obj.delimiter, type, align));
        }
        return { list: list, isAddedBlankHead: obj.isAddedBlankHead };
    }
    // 表データの取得
    getTableInfo(doc, range, formatType) {
        // 各行の解析
        var grid = [];
        var info = {
            // 行頭にデリミタがあったか
            hasDelimiterAtLineHead: false
        };
        for (var i = range.start.line; i <= range.end.line; i++) {
            var obj = this.getCellInfoList(doc.lineAt(i), formatType);
            grid.push(obj.list);
            // 一行でも先頭に空白を追加していなかったら、行頭デリミタがあったと判定
            if (!obj.isAddedBlankHead)
                info.hasDelimiterAtLineHead = true;
        }
        return new table_1.TableInfo(this.settings, range, grid, info);
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