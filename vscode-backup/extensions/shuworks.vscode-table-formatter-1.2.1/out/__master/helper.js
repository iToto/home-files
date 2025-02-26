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
})(TableFormatType = exports.TableFormatType || (exports.TableFormatType = {}));
;
// Table line
var TableLineFlag;
(function (TableLineFlag) {
    TableLineFlag[TableLineFlag["None"] = 0] = "None";
    TableLineFlag[TableLineFlag["HasPipe"] = 1] = "HasPipe";
    TableLineFlag[TableLineFlag["PlusSeparator"] = 2] = "PlusSeparator";
    TableLineFlag[TableLineFlag["SimpleSeparator"] = 4] = "SimpleSeparator";
    TableLineFlag[TableLineFlag["NotEmpty"] = 8] = "NotEmpty";
})(TableLineFlag = exports.TableLineFlag || (exports.TableLineFlag = {}));
;
// Separator type
var SeparatorType;
(function (SeparatorType) {
    SeparatorType[SeparatorType["None"] = 0] = "None";
    SeparatorType[SeparatorType["Minus"] = 1] = "Minus";
    SeparatorType[SeparatorType["Equall"] = 2] = "Equall";
})(SeparatorType || (SeparatorType = {}));
;
class TableHelper {
    constructor(config) {
        // オブジェクトは参照渡し
        this.settings = config;
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
                // 空白で区切る（""か''で囲まれた範囲の空白は無視する）
                cells = text.match(/('[^']*'|"[^"]*")|[^ ]+/g);
                delimiter = table_1.DelimiterType.Space;
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
    // 文字数の取得
    getStringLength(str) {
        if (this.settings.common.explicitFullwidthChars.length == 0) {
            return strWidth(str);
        }
        // 強制的に全角判定の文字が含まれていたら個数分加算する
        let cnt = strWidth(str);
        this.settings.common.explicitFullwidthChars.forEach((reg, i) => {
            cnt += (str.match(reg) || []).length;
        });
        return cnt;
    }
    // 行の解析
    getCellInfoList(line, formatType) {
        if (line.isEmptyOrWhitespace)
            return [];
        var splitText = this.getSplitLineText(line.text, formatType);
        var list = [];
        var cells = splitText.cells;
        var delimiter = splitText.delimiter;
        // 空白列を追加
        list.push(new table_1.CellInfo(line.firstNonWhitespaceCharacterIndex, delimiter));
        for (var i = 0; i < cells.length; i++) {
            var trimed = trim(cells[i]);
            // 先頭は空白で追加済みなので無視する
            if (i == 0)
                continue;
            var size = this.getStringLength(trimed);
            var type = (size == 0) ? table_1.CellType.CM_Blank : table_1.CellType.CM_Content;
            var align = table_1.CellAlign.Left;
            // Common  ----------------
            if (/^-+$/.test(trimed)) {
                type = table_1.CellType.CM_MinusSeparator;
                align = table_1.CellAlign.Left;
            }
            else if (/^=+$/.test(trimed)) {
                type = table_1.CellType.CM_EquallSeparator;
                align = table_1.CellAlign.Left;
            }
            else {
                switch (formatType) {
                    case TableFormatType.Normal:
                        // Markdown ----------------
                        if (/^:-+$/.test(trimed)) {
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
                            size = this.getStringLength(trim(trimed.substring(2)));
                            type = table_1.CellType.TT_HeaderPrefix;
                            align = table_1.CellAlign.Left;
                        }
                        else if (/^<\./.test(trimed)) {
                            size = this.getStringLength(trim(trimed.substring(2)));
                            type = table_1.CellType.TT_LeftPrefix;
                            align = table_1.CellAlign.Left;
                        }
                        else if (/^>\./.test(trimed)) {
                            size = this.getStringLength(trim(trimed.substring(2)));
                            type = table_1.CellType.TT_RightPrefix;
                            align = table_1.CellAlign.Right;
                        }
                        else if (/^=\./.test(trimed)) {
                            size = this.getStringLength(trim(trimed.substring(2)));
                            type = table_1.CellType.TT_CenterPrefix;
                            align = table_1.CellAlign.Center;
                        }
                        break;
                    case TableFormatType.Simple:
                        // NOP ----------------
                        break;
                }
            }
            list.push(new table_1.CellInfo(size, delimiter, type, align));
        }
        return list;
    }
    // 表データの正規化
    getNormalizedCellGrid(grid) {
        // 全行のサイズを揃える
        var max = 0;
        grid.forEach(row => {
            max = Math.max(max, row.length);
        });
        grid.forEach(row => {
            // デリミタも揃える
            var delimiter = (row.length > 0) ? row[0].delimiter : table_1.DelimiterType.Pipe;
            for (var c = row.length; c < max; c++) {
                row.push(new table_1.CellInfo(0, delimiter));
            }
        });
        // セパレータタイプを確定する
        grid.forEach(row => {
            // セパレータ行かの判定
            var rowType = SeparatorType.None;
            for (var i = 0; i < row.length; i++) {
                var cell = row[i];
                // デリミタがPlusなら初期値をMinusにする（一応この時点でセパレータ行で確定ではあるが特に何もしない）
                if (i == 0 && cell.delimiter == table_1.DelimiterType.Plus) {
                    rowType = SeparatorType.Minus;
                }
                // 文字列なら非セパレータ行で確定
                if (cell.type == table_1.CellType.CM_Content) {
                    rowType = SeparatorType.None;
                    break;
                }
                // セルのタイプで判定
                if (cell.type == table_1.CellType.CM_MinusSeparator) {
                    rowType = SeparatorType.Minus;
                }
                else if (cell.type == table_1.CellType.CM_EquallSeparator) {
                    rowType = SeparatorType.Equall;
                }
                else if (cell.type == table_1.CellType.MD_LeftSeparator || cell.type == table_1.CellType.MD_RightSeparator || cell.type == table_1.CellType.MD_CenterSeparator) {
                    rowType = SeparatorType.Minus;
                }
            }
            // セパレータタイプを補正
            row.forEach((cell, i) => {
                // セパレータ行でない場合、セルタイプを文字列に（-や:-を文字として扱う）
                if (rowType == SeparatorType.None) {
                    if (cell.type == table_1.CellType.CM_MinusSeparator || cell.type == table_1.CellType.CM_EquallSeparator ||
                        cell.type == table_1.CellType.MD_LeftSeparator || cell.type == table_1.CellType.MD_RightSeparator || cell.type == table_1.CellType.MD_CenterSeparator) {
                        cell.setType(table_1.CellType.CM_Content);
                    }
                }
                else {
                    if (i != 0 && cell.type == table_1.CellType.CM_Blank) {
                        switch (rowType) {
                            case SeparatorType.Minus:
                                cell.setType(table_1.CellType.CM_MinusSeparator);
                                break;
                            case SeparatorType.Equall:
                                cell.setType(table_1.CellType.CM_EquallSeparator);
                                break;
                        }
                    }
                }
            });
        });
        // Markdownのセパレータのサイズを設定
        // 左右のスペーサー分がフォーマット時に加算されるため-2する
        let offset = (this.settings.markdown.oneSpacePadding) ? 0 : -2;
        grid.forEach(row => {
            row.forEach(cell => {
                if (cell.type == table_1.CellType.CM_MinusSeparator || cell.type == table_1.CellType.CM_EquallSeparator) {
                    // 最小である3文字にする（---）
                    cell.setSize(3 + offset);
                }
                else if (cell.type == table_1.CellType.MD_LeftSeparator || cell.type == table_1.CellType.MD_RightSeparator) {
                    // 最小である4文字にする（:---, ---:）
                    cell.setSize(4 + offset);
                }
                else if (cell.type == table_1.CellType.MD_CenterSeparator) {
                    // 最小である5文字にする（:---:）
                    cell.setSize(5 + offset);
                }
            });
        });
        // Textileのサイズを設定
        grid.forEach(row => {
            // プレフィックス分のパディングを他の行の同列に設定する
            row.forEach((cell, i) => {
                if (cell.type == table_1.CellType.TT_HeaderPrefix || cell.type == table_1.CellType.TT_LeftPrefix || cell.type == table_1.CellType.TT_RightPrefix || cell.type == table_1.CellType.TT_CenterPrefix) {
                    grid.forEach(elem => {
                        if (i < elem.length) {
                            elem[i].setPadding(2);
                        }
                    });
                }
            });
        });
        // Markdownの位置揃え
        for (var r = grid.length - 1; r >= 0; r--) {
            var row = grid[r];
            // 各列の位置揃えを変更する
            row.forEach((cell, c) => {
                switch (cell.type) {
                    case table_1.CellType.MD_LeftSeparator:
                        grid.forEach(elem => {
                            if (c < elem.length) {
                                elem[c].setAlign(table_1.CellAlign.Left);
                            }
                        });
                        break;
                    case table_1.CellType.MD_RightSeparator:
                        grid.forEach(elem => {
                            if (c < elem.length) {
                                elem[c].setAlign(table_1.CellAlign.Right);
                            }
                        });
                        break;
                    case table_1.CellType.MD_CenterSeparator:
                        grid.forEach(elem => {
                            if (c < elem.length) {
                                elem[c].setAlign(table_1.CellAlign.Center);
                            }
                        });
                        break;
                }
            });
        }
        // Textileの位置揃え
        grid.forEach(row => {
            // 各セルの位置揃えを変更する
            row.forEach(cell => {
                switch (cell.type) {
                    case table_1.CellType.TT_LeftPrefix:
                        cell.setAlign(table_1.CellAlign.Left);
                        break;
                    case table_1.CellType.TT_RightPrefix:
                        cell.setAlign(table_1.CellAlign.Right);
                        break;
                    case table_1.CellType.TT_CenterPrefix:
                        cell.setAlign(table_1.CellAlign.Center);
                        break;
                }
            });
        });
        return grid;
    }
    // 表データの取得
    getTableInfo(doc, range, formatType) {
        // 各行の解析
        var grid = [];
        for (var i = range.start.line; i <= range.end.line; i++) {
            grid.push(this.getCellInfoList(doc.lineAt(i), formatType));
        }
        return new table_1.TableInfo(range, this.getNormalizedCellGrid(grid));
    }
}
exports.TableHelper = TableHelper;
//# sourceMappingURL=helper.js.map