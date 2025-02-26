'use strict';
const table_1 = require("./table");
var utilPad = require('utils-pad-string');
var strWidth = require('string-width');
var trim = require('trim');
class TableFormatter {
    constructor(config, helper) {
        // オブジェクトは参照渡し
        this.settings = config;
        this.tableHelper = helper;
    }
    dispose() {
    }
    // フォーマット済みの文字列を返す
    getFormatTableText(doc, info, formatType, option) {
        if (!info.isValid())
            return "";
        var formatted = "";
        var maxList = info.getMaxCellSizeList();
        info.cellGrid.forEach((row, i) => {
            var line = i + info.range.start.line;
            formatted += this.getFormattedLineText(this.tableHelper.getSplitLineText(doc.lineAt(line).text, formatType).cells, row, maxList, formatType);
            if (line != info.range.end.line)
                formatted += '\n';
        });
        return formatted;
    }
    // 行のフォーマット済みの文字列を返す
    getFormattedLineText(cells, cellInfoList, maxlist, formatType, option) {
        if (cellInfoList && maxlist && cellInfoList.length != maxlist.length)
            return "";
        var formatted = "";
        maxlist.forEach((elem, i) => {
            var cellInfo = cellInfoList[i];
            var trimed = (i < cells.length) ? trim(cells[i]) : "";
            // @TODO: ここでまたサイズ数えてる、paddingが普通に文字数を数えるためだと思われるので、保持するようにする
            var sub = this.tableHelper.getStringLength(trimed) - trimed.length;
            var size = (sub == 0) ? elem : elem - sub;
            let hasOneSpace = true;
            switch (cellInfo.type) {
                // Common ----------------
                case table_1.CellType.CM_MinusSeparator:
                    size += cellInfo.padding;
                    // 左右に空白を設けるか
                    if (this.settings.markdown.oneSpacePadding) {
                        size += (cellInfo.delimiter == table_1.DelimiterType.Plus) ? 2 : 0;
                        hasOneSpace = true;
                    }
                    else {
                        size += (cellInfo.delimiter == table_1.DelimiterType.Plus || cellInfo.delimiter == table_1.DelimiterType.Pipe) ? 2 : 0;
                        hasOneSpace = false;
                    }
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter, hasOneSpace);
                    formatted += this.getAlignedText("", size, "-", table_1.CellAlign.Center);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter, hasOneSpace);
                    break;
                case table_1.CellType.CM_EquallSeparator:
                    size += cellInfo.padding;
                    // 左右に空白を設けるか
                    if (this.settings.markdown.oneSpacePadding) {
                        size += (cellInfo.delimiter == table_1.DelimiterType.Plus) ? 2 : 0;
                        hasOneSpace = true;
                    }
                    else {
                        size += (cellInfo.delimiter == table_1.DelimiterType.Plus || cellInfo.delimiter == table_1.DelimiterType.Pipe) ? 2 : 0;
                        hasOneSpace = false;
                    }
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter, hasOneSpace);
                    formatted += this.getAlignedText("", size, "=", table_1.CellAlign.Center);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter, hasOneSpace);
                    break;
                // Markdown ----------------
                case table_1.CellType.MD_LeftSeparator:
                    size += cellInfo.padding;
                    // 左右に空白を設けるか
                    if (this.settings.markdown.oneSpacePadding) {
                        hasOneSpace = true;
                    }
                    else {
                        size += 2;
                        hasOneSpace = false;
                    }
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter, hasOneSpace);
                    formatted += this.getAlignedText(":---", size, "-", table_1.CellAlign.Left);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter, hasOneSpace);
                    break;
                case table_1.CellType.MD_RightSeparator:
                    size += cellInfo.padding;
                    // 左右に空白を設けるか
                    if (this.settings.markdown.oneSpacePadding) {
                        hasOneSpace = true;
                    }
                    else {
                        size += 2;
                        hasOneSpace = false;
                    }
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter, hasOneSpace);
                    formatted += this.getAlignedText("---:", size, "-", table_1.CellAlign.Right);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter, hasOneSpace);
                    break;
                case table_1.CellType.MD_CenterSeparator:
                    size += cellInfo.padding;
                    // 左右に空白を設けるか
                    if (this.settings.markdown.oneSpacePadding) {
                        hasOneSpace = true;
                    }
                    else {
                        size += 2;
                        hasOneSpace = false;
                    }
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter, hasOneSpace);
                    formatted += ":" + this.getAlignedText("", size - 2, "-", table_1.CellAlign.Center) + ":";
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter, hasOneSpace);
                    break;
                // Textile ----------------
                case table_1.CellType.TT_HeaderPrefix:
                    trimed = trim(trimed.substring(2));
                    formatted += (i == 0) ? "" : (size == 0) ? "_." : "_. ";
                    formatted += this.getAlignedText(trimed, size, " ", cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_LeftPrefix:
                    trimed = trim(trimed.substring(2));
                    formatted += (i == 0) ? "" : (size == 0) ? "<." : "<. ";
                    formatted += this.getAlignedText(trimed, size, " ", cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_RightPrefix:
                    trimed = trim(trimed.substring(2));
                    formatted += (i == 0) ? "" : (size == 0) ? ">." : ">. ";
                    formatted += this.getAlignedText(trimed, size, " ", cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_CenterPrefix:
                    trimed = trim(trimed.substring(2));
                    formatted += (i == 0) ? "" : (size == 0) ? "=." : "=. ";
                    formatted += this.getAlignedText(trimed, size, " ", cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                // Etc ----------------
                default:
                    formatted += this.getPaddingText(i, size, cellInfo.padding, cellInfo.delimiter);
                    formatted += this.getAlignedText(trimed, size, " ", cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
            }
        });
        return formatted;
    }
    getPaddingText(cell, size, padding, delimiter, hasOneSpace = true) {
        var spacer = "";
        switch (delimiter) {
            case table_1.DelimiterType.Pipe:
                if (hasOneSpace) {
                    spacer = (cell == 0 || size == 0) ? "" : " ";
                }
                else {
                    spacer = "";
                }
                break;
            case table_1.DelimiterType.Plus:
                spacer = "";
                break;
            case table_1.DelimiterType.Space:
                spacer = "";
                break;
        }
        return spacer + utilPad("", padding);
    }
    getAlignedText(text, size, pad, align) {
        var opt = {};
        switch (align) {
            case table_1.CellAlign.Left:
                opt = { rpad: pad };
                break;
            case table_1.CellAlign.Right:
                opt = { lpad: pad };
                break;
            case table_1.CellAlign.Center:
                opt = { lpad: pad, rpad: pad };
                break;
        }
        return utilPad(text, size, opt);
    }
    getDelimiterText(cell, rowSize, delimiter, hasOneSpace = true) {
        switch (delimiter) {
            case table_1.DelimiterType.Pipe:
                if (hasOneSpace) {
                    return (cell == 0) ? "|" : " |";
                }
                return "|";
            case table_1.DelimiterType.Plus:
                return "+";
            case table_1.DelimiterType.Space:
                // 2スペース推奨
                return (cell == 0 || cell == rowSize - 1) ? "" : "  ";
        }
        return "";
    }
}
exports.TableFormatter = TableFormatter;
//# sourceMappingURL=formatter.js.map