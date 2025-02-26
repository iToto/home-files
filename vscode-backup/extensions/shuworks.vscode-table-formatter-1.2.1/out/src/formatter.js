'use strict';
const table_1 = require("./table");
const helper_1 = require("./helper");
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
    getFormattedTableText(doc, info, formatType, option) {
        if (!info.isValid())
            return "";
        var formatted = "";
        var maxList = info.getMaxCellSizeList();
        info.cellGrid.forEach((row, i) => {
            var line = i + info.range.start.line;
            formatted += this.getFormattedLineText(i, info.property, row, maxList, formatType);
            if (line != info.range.end.line)
                formatted += '\n';
        });
        return formatted;
    }
    // フォーマットするサイズを取得
    getFormattingSize(cellIndex, cellInfo, maxCount, cellCount, isMarkdown, isPutOneSpace, isPutEdges) {
        // utilPad()がlength基準で処理されるので、差分から長さを修正
        let size = maxCount - cellInfo.diff;
        switch (cellInfo.type) {
            case table_1.CellType.CM_MinusSeparator:
            case table_1.CellType.CM_EquallSeparator:
                size += cellInfo.padding;
                // コンフィグ：セパレータの左右をパディングする場合
                if (this.settings.markdown.oneSpacePadding) {
                    size += (cellInfo.delimiter == table_1.DelimiterType.Plus) ? 2 : 0;
                }
                else {
                    size += (cellInfo.delimiter == table_1.DelimiterType.Plus || cellInfo.delimiter == table_1.DelimiterType.Pipe) ? 2 : 0;
                }
                // コンフィグ：テーブル両端のデリミタをなくす場合
                if (!isPutEdges && !isPutOneSpace) {
                    // ２セル目か末尾セル（先頭セルは必ず空白）
                    if (cellIndex == 1) {
                        size -= 1;
                    }
                    else if (cellIndex == cellCount - 1) {
                        size -= 1;
                    }
                }
                break;
            case table_1.CellType.MD_LeftSeparator:
            case table_1.CellType.MD_RightSeparator:
            case table_1.CellType.MD_CenterSeparator:
                size += cellInfo.padding;
                // コンフィグ：セパレータの左右をパディングする場合
                if (!this.settings.markdown.oneSpacePadding) {
                    size += 2;
                }
                // コンフィグ：テーブル両端のデリミタをなくす場合
                if (!isPutEdges && !isPutOneSpace) {
                    // ２セル目か末尾セル（先頭セルは必ず空白）
                    if (cellIndex == 1) {
                        size -= 1;
                    }
                    else if (cellIndex == cellCount - 1) {
                        size -= 1;
                    }
                }
                break;
        }
        return size;
    }
    // デリミタの取得
    getDelimiter(cellIndex, cellInfo, cellCount, isMarkdown, isPutOneSpace, isPutEdges) {
        // マークダウンでないときはそのまま返す（isPutEdgesがtrueのときはMarkdown確定のはずだが念のため）
        if (!isMarkdown)
            return cellInfo.delimiter;
        // コンフィグ：通常通りテーブル両端を表示する場合そのまま返す
        if (isPutEdges)
            return cellInfo.delimiter;
        let delimiter = cellInfo.delimiter;
        switch (cellInfo.type) {
            case table_1.CellType.CM_MinusSeparator:
            case table_1.CellType.CM_EquallSeparator:
            case table_1.CellType.MD_LeftSeparator:
            case table_1.CellType.MD_RightSeparator:
            case table_1.CellType.MD_CenterSeparator:
                // コンフィグ：テーブル両端のデリミタをなくす場合
                if (cellIndex == cellCount - 1) {
                    delimiter = table_1.DelimiterType.None;
                }
                break;
            case table_1.CellType.CM_Blank:
            case table_1.CellType.CM_Content:
                // コンフィグ：テーブル両端のデリミタをなくす場合
                if (cellIndex == 0) {
                    delimiter = table_1.DelimiterType.None;
                }
                else if (cellIndex == cellCount - 1) {
                    delimiter = table_1.DelimiterType.None;
                }
                break;
        }
        return delimiter;
    }
    //
    getAlign(cellIndex, cellInfo, rowIndex, prop, formatType) {
        // コンフィグ：強制的に中心揃えにしない場合はそのまま返す
        if (!this.settings.common.centerAlignedHeader)
            return cellInfo.align;
        let align = cellInfo.align;
        switch (cellInfo.type) {
            case table_1.CellType.TT_HeaderPrefix:
                // コンフィグ：強制的に中心揃えにする
                align = table_1.CellAlign.Center;
                break;
            case table_1.CellType.CM_Blank:
            case table_1.CellType.CM_Content:
                // コンフィグ：強制的に中心揃えにする
                if (formatType == helper_1.TableFormatType.Normal) {
                    if (prop.markdownTableHeaderIndexes.has(rowIndex)) {
                        align = table_1.CellAlign.Center;
                    }
                    else if (prop.gridTableHeaderIndexes.has(rowIndex)) {
                        align = table_1.CellAlign.Center;
                    }
                }
                else if (formatType == helper_1.TableFormatType.Simple) {
                    if (prop.simpleTableHeaderIndexes.has(rowIndex)) {
                        align = table_1.CellAlign.Center;
                    }
                }
                break;
        }
        return align;
    }
    getIsPutPaddingOneSpace(cellIndex, isPutEdges, defaultOut) {
        // 通常通りテーブルの両端にセパレータを表示する場合はそのまま返す
        if (isPutEdges)
            return defaultOut;
        // コンフィグ：ボーターレスの場合で２セル目（先頭セルは空白）なら詰めるためにfalseを返す
        if (cellIndex == 1) {
            return false;
        }
        return defaultOut;
    }
    getIsPutOneSpace(prop) {
        return this.settings.markdown.oneSpacePadding;
    }
    getIsPutEdges(prop) {
        // マークダウン以外は無視
        if (!prop.isMarkdown)
            return true;
        let edgesType = this.settings.markdown.tableEdgesType;
        if (edgesType == table_1.TableEdgesType.Auto) {
            edgesType = (prop.hasDelimiterAtLineHead) ? table_1.TableEdgesType.Normal : table_1.TableEdgesType.Borderless;
        }
        switch (edgesType) {
            case table_1.TableEdgesType.Normal:
                return true;
            case table_1.TableEdgesType.Borderless:
                return false;
        }
        return false;
    }
    // 行のフォーマット済みの文字列を返す
    getFormattedLineText(rowIndex, prop, cellInfoList, maxlist, formatType, option) {
        if (cellInfoList && maxlist && cellInfoList.length != maxlist.length)
            return "";
        // 行全体のフォーマット設定
        let isPutOneSpace = this.getIsPutOneSpace(prop);
        let isPutEdges = this.getIsPutEdges(prop);
        var formatted = "";
        var cellCount = maxlist.length;
        maxlist.forEach((maxCount, i) => {
            var cellInfo = cellInfoList[i];
            // 各セルのフォーマット設定
            let formattingSize = this.getFormattingSize(i, cellInfo, maxCount, cellCount, prop.isMarkdown, isPutOneSpace, isPutEdges);
            let delimiter = this.getDelimiter(i, cellInfo, cellCount, prop.isMarkdown, isPutOneSpace, isPutEdges);
            let align = this.getAlign(i, cellInfo, rowIndex, prop, formatType);
            let isPutPaddingOneSpace = false;
            switch (cellInfo.type) {
                // Common ----------------
                case table_1.CellType.CM_MinusSeparator:
                    isPutPaddingOneSpace = this.getIsPutPaddingOneSpace(i, isPutEdges, isPutOneSpace);
                    formatted += this.getPaddingText(i, formattingSize, 0, delimiter, isPutPaddingOneSpace);
                    formatted += this.getAlignedText("", formattingSize, "-", table_1.CellAlign.Center);
                    formatted += this.getDelimiterText(i, cellInfoList.length, delimiter, isPutOneSpace);
                    break;
                case table_1.CellType.CM_EquallSeparator:
                    isPutPaddingOneSpace = this.getIsPutPaddingOneSpace(i, isPutEdges, isPutOneSpace);
                    formatted += this.getPaddingText(i, formattingSize, 0, delimiter, isPutPaddingOneSpace);
                    formatted += this.getAlignedText("", formattingSize, "=", table_1.CellAlign.Center);
                    formatted += this.getDelimiterText(i, cellInfoList.length, delimiter, isPutOneSpace);
                    break;
                // Markdown ----------------
                case table_1.CellType.MD_LeftSeparator:
                    isPutPaddingOneSpace = this.getIsPutPaddingOneSpace(i, isPutEdges, isPutOneSpace);
                    formatted += this.getPaddingText(i, formattingSize, 0, delimiter, isPutPaddingOneSpace);
                    formatted += this.getAlignedText(":---", formattingSize, "-", table_1.CellAlign.Left);
                    formatted += this.getDelimiterText(i, cellInfoList.length, delimiter, isPutOneSpace);
                    break;
                case table_1.CellType.MD_RightSeparator:
                    isPutPaddingOneSpace = this.getIsPutPaddingOneSpace(i, isPutEdges, isPutOneSpace);
                    formatted += this.getPaddingText(i, formattingSize, 0, delimiter, isPutPaddingOneSpace);
                    formatted += this.getAlignedText("---:", formattingSize, "-", table_1.CellAlign.Right);
                    formatted += this.getDelimiterText(i, cellInfoList.length, delimiter, isPutOneSpace);
                    break;
                case table_1.CellType.MD_CenterSeparator:
                    isPutPaddingOneSpace = this.getIsPutPaddingOneSpace(i, isPutEdges, isPutOneSpace);
                    formatted += this.getPaddingText(i, formattingSize, 0, delimiter, isPutPaddingOneSpace);
                    formatted += ":" + this.getAlignedText("", formattingSize - 2, "-", table_1.CellAlign.Center) + ":";
                    formatted += this.getDelimiterText(i, cellInfoList.length, delimiter, isPutOneSpace);
                    break;
                // Textile ----------------
                case table_1.CellType.TT_HeaderPrefix:
                    formatted += (i == 0) ? "" : (formattingSize == 0) ? "_." : "_. ";
                    formatted += this.getAlignedText(cellInfo.string, formattingSize, " ", align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_LeftPrefix:
                    formatted += (i == 0) ? "" : (formattingSize == 0) ? "<." : "<. ";
                    formatted += this.getAlignedText(cellInfo.string, formattingSize, " ", table_1.CellAlign.Left);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_RightPrefix:
                    formatted += (i == 0) ? "" : (formattingSize == 0) ? ">." : ">. ";
                    formatted += this.getAlignedText(cellInfo.string, formattingSize, " ", table_1.CellAlign.Right);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                case table_1.CellType.TT_CenterPrefix:
                    formatted += (i == 0) ? "" : (formattingSize == 0) ? "=." : "=. ";
                    formatted += this.getAlignedText(cellInfo.string, formattingSize, " ", table_1.CellAlign.Center);
                    formatted += this.getDelimiterText(i, cellInfoList.length, cellInfo.delimiter);
                    break;
                // 空白、文字列 ----------------
                case table_1.CellType.CM_Blank:
                case table_1.CellType.CM_Content:
                    isPutPaddingOneSpace = this.getIsPutPaddingOneSpace(i, isPutEdges, true);
                    formatted += this.getPaddingText(i, formattingSize, cellInfo.padding, delimiter, isPutPaddingOneSpace);
                    formatted += this.getAlignedText(cellInfo.string, formattingSize, " ", align);
                    formatted += this.getDelimiterText(i, cellInfoList.length, delimiter);
                    break;
            }
        });
        // コンフィグ：末尾スペースを削除する
        if (this.settings.common.trimTrailingWhitespace) {
            formatted = trim.right(formatted);
        }
        return formatted;
    }
    getPaddingText(cellIndex, size, padding, delimiter, isPutOneSpace = true) {
        var spacer = "";
        switch (delimiter) {
            case table_1.DelimiterType.None:
                // 行の先頭セルか末尾セルしか来ないはず
                if (isPutOneSpace) {
                    spacer = (cellIndex == 0 || size == 0) ? "" : " ";
                }
                else {
                    spacer = "";
                }
                break;
            case table_1.DelimiterType.Pipe:
                if (isPutOneSpace) {
                    spacer = (cellIndex == 0 || size == 0) ? "" : " ";
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
    getDelimiterText(cellIndex, rowSize, delimiter, isPutOneSpace = true) {
        switch (delimiter) {
            case table_1.DelimiterType.None:
                // 行の先頭セルか末尾セルしか来ないはず
                if (isPutOneSpace) {
                    return (cellIndex == 0) ? "" : " ";
                }
                return "";
            case table_1.DelimiterType.Pipe:
                if (isPutOneSpace) {
                    return (cellIndex == 0) ? "|" : " |";
                }
                return "|";
            case table_1.DelimiterType.Plus:
                return "+";
            case table_1.DelimiterType.Space:
                // 言語仕様的に２スペース推奨らしい
                return (cellIndex == 0 || cellIndex == rowSize - 1) ? "" : "  ";
        }
        return "";
    }
}
exports.TableFormatter = TableFormatter;
//# sourceMappingURL=formatter.js.map