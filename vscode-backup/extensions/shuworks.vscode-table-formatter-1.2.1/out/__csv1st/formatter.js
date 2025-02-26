'use strict';
const table_1 = require("./table");
const helper_1 = require("./helper");
var utilPad = require('utils-pad-string');
var strWidth = require('string-width');
var trim = require('trim');
class TableFormatter {
    constructor() {
    }
    dispose() {
    }
    // フォーマット済みの文字列を返す
    getFormatTableText(doc, info, formatType, option) {
        if (!info.isValid())
            return "";
        var tableHelper = new helper_1.TableHelper();
        var formatted = "";
        var maxList = info.getMaxCellSizeList();
        info.cellGrid.forEach((row, i) => {
            var line = i + info.range.start.line;
            formatted += this.getFormattedLineText(tableHelper.getSplitLineText(doc.lineAt(line).text, formatType).cells, row, maxList, formatType);
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
            var sub = strWidth(trimed) - trimed.length;
            var size = (sub == 0) ? elem : elem - sub;
            switch (cellInfo.type) {
                // Common ----------------
                case table_1.CellType.CM_MinusSeparator:
                    size += cellInfo.padding;
                    size += (cellInfo.delimiter == table_1.DelimiterType.Plus) ? 2 : 0;
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter);
                    formatted += utilPad(trimed, size, { rpad: "-" });
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.CM_EquallSeparator:
                    size += cellInfo.padding;
                    size += (cellInfo.delimiter == table_1.DelimiterType.Plus) ? 2 : 0;
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter);
                    formatted += utilPad(trimed, size, { rpad: "=" });
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                // Markdown ----------------
                case table_1.CellType.MD_LeftSeparator:
                    size += cellInfo.padding;
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter);
                    formatted += utilPad(trimed, size, { rpad: "-" });
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.MD_RightSeparator:
                    size += cellInfo.padding;
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter);
                    formatted += utilPad(trimed, size, { lpad: "-" });
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.MD_CenterSeparator:
                    size += cellInfo.padding;
                    formatted += this.getPaddingText(i, size, 0, cellInfo.delimiter);
                    formatted += utilPad(":", size - 1, { rpad: "-" }) + ":";
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                // Textile ----------------
                case table_1.CellType.TT_HeaderPrefix:
                    trimed = trim(trimed.substring(2));
                    formatted += (i == 0) ? "" : (size == 0) ? "_." : "_. ";
                    formatted += this.getAlignedText(trimed, size, cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_LeftPrefix:
                    trimed = trim(trimed.substring(2));
                    formatted += (i == 0) ? "" : (size == 0) ? "<." : "<. ";
                    formatted += this.getAlignedText(trimed, size, cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_RightPrefix:
                    trimed = trim(trimed.substring(2));
                    formatted += (i == 0) ? "" : (size == 0) ? ">." : ">. ";
                    formatted += this.getAlignedText(trimed, size, cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_CenterPrefix:
                    trimed = trim(trimed.substring(2));
                    formatted += (i == 0) ? "" : (size == 0) ? "=." : "=. ";
                    formatted += this.getAlignedText(trimed, size, cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                // Etc ----------------
                default:
                    formatted += this.getPaddingText(i, size, cellInfo.padding, cellInfo.delimiter);
                    formatted += this.getAlignedText(trimed, size, cellInfo.align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
            }
        });
        return formatted;
    }
    getPaddingText(cell, size, padding, delimiter) {
        var str = "";
        switch (delimiter) {
            case table_1.DelimiterType.Pipe:
                str = (cell == 0 || size == 0) ? "" : " ";
                break;
            case table_1.DelimiterType.Plus:
                str = "";
                break;
            case table_1.DelimiterType.Space:
                str = "";
                break;
            case table_1.DelimiterType.Comma:
                str = "";
                break;
        }
        return str + utilPad("", padding);
    }
    getAlignedText(text, size, align) {
        var opt = {};
        switch (align) {
            case table_1.CellAlign.Left:
                opt = { rpad: " " };
                break;
            case table_1.CellAlign.Right:
                opt = { lpad: " " };
                break;
            case table_1.CellAlign.Center:
                opt = { lpad: " ", rpad: " " };
                break;
        }
        return utilPad(text, size, opt);
    }
    getDelimiterText(cell, rowSize, delimiter) {
        switch (delimiter) {
            case table_1.DelimiterType.Pipe:
                return (cell == 0) ? "|" : " |";
            case table_1.DelimiterType.Plus:
                return "+";
            case table_1.DelimiterType.Space:
                // 2スペース推奨
                return (cell == 0 || cell == rowSize - 1) ? "" : "  ";
            case table_1.DelimiterType.Comma:
                return (cell == 0 || cell == rowSize - 1) ? "" : ", ";
        }
        return "";
    }
}
exports.TableFormatter = TableFormatter;
//# sourceMappingURL=formatter.js.map