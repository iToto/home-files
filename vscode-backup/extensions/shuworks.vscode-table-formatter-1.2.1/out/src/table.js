'use strict';
var strWidth = require('string-width');
var trim = require('trim');
// Cell type
var CellType;
(function (CellType) {
    CellType[CellType["CM_Blank"] = 0] = "CM_Blank";
    CellType[CellType["CM_Content"] = 1] = "CM_Content";
    CellType[CellType["CM_MinusSeparator"] = 2] = "CM_MinusSeparator";
    CellType[CellType["CM_EquallSeparator"] = 3] = "CM_EquallSeparator";
    CellType[CellType["MD_LeftSeparator"] = 4] = "MD_LeftSeparator";
    CellType[CellType["MD_RightSeparator"] = 5] = "MD_RightSeparator";
    CellType[CellType["MD_CenterSeparator"] = 6] = "MD_CenterSeparator";
    CellType[CellType["TT_HeaderPrefix"] = 7] = "TT_HeaderPrefix";
    CellType[CellType["TT_LeftPrefix"] = 8] = "TT_LeftPrefix";
    CellType[CellType["TT_RightPrefix"] = 9] = "TT_RightPrefix";
    CellType[CellType["TT_CenterPrefix"] = 10] = "TT_CenterPrefix";
})(CellType = exports.CellType || (exports.CellType = {}));
;
// Cell align
var CellAlign;
(function (CellAlign) {
    CellAlign[CellAlign["Left"] = 0] = "Left";
    CellAlign[CellAlign["Center"] = 1] = "Center";
    CellAlign[CellAlign["Right"] = 2] = "Right";
})(CellAlign = exports.CellAlign || (exports.CellAlign = {}));
;
// Delimiter type
var DelimiterType;
(function (DelimiterType) {
    DelimiterType[DelimiterType["None"] = 0] = "None";
    DelimiterType[DelimiterType["Pipe"] = 1] = "Pipe";
    DelimiterType[DelimiterType["Plus"] = 2] = "Plus";
    DelimiterType[DelimiterType["Space"] = 3] = "Space";
    // Comma
})(DelimiterType = exports.DelimiterType || (exports.DelimiterType = {}));
;
// Separator type
var SeparatorType;
(function (SeparatorType) {
    SeparatorType[SeparatorType["None"] = 0] = "None";
    SeparatorType[SeparatorType["Minus"] = 1] = "Minus";
    SeparatorType[SeparatorType["Equall"] = 2] = "Equall";
})(SeparatorType = exports.SeparatorType || (exports.SeparatorType = {}));
;
// Markdown table edges type
var TableEdgesType;
(function (TableEdgesType) {
    TableEdgesType[TableEdgesType["Auto"] = 0] = "Auto";
    TableEdgesType[TableEdgesType["Normal"] = 1] = "Normal";
    TableEdgesType[TableEdgesType["Borderless"] = 2] = "Borderless";
})(TableEdgesType = exports.TableEdgesType || (exports.TableEdgesType = {}));
;
;
;
// Cell info
class CellInfo {
    constructor(settings, trimmed, delimiter = DelimiterType.Pipe, type = CellType.CM_Blank, align = CellAlign.Left, padding = 0) {
        this._settings = settings;
        this._string = trimmed;
        this._size = this.getStringLength(this._string);
        this._diff = this._size - trimmed.length;
        this._delimiter = delimiter;
        this._type = type;
        this._align = align;
        this._padding = padding;
    }
    get string() {
        return this._string;
    }
    get size() {
        return this._size;
    }
    get diff() {
        return this._diff;
    }
    get delimiter() {
        return this._delimiter;
    }
    get type() {
        return this._type;
    }
    get align() {
        return this._align;
    }
    get padding() {
        return this._padding;
    }
    isValid() {
        if (this._size < 0)
            return false;
        return true;
    }
    setString(trimmed) {
        this._string = trimmed;
    }
    setSize(size) {
        this._size = size;
    }
    setDiff(length) {
        this._diff = length;
    }
    setDelimiter(delimiter) {
        this._delimiter = delimiter;
    }
    setType(type) {
        this._type = type;
    }
    setAlign(align) {
        this._align = align;
    }
    setPadding(padding) {
        this._padding = padding;
    }
    // 文字数の取得
    getStringLength(str) {
        if (this._settings.common.explicitFullwidthChars.length == 0) {
            return strWidth(str);
        }
        // コンフィグ：強制的に全角判定の文字が含まれていたら個数分加算する
        let cnt = strWidth(str);
        this._settings.common.explicitFullwidthChars.forEach((reg, i) => {
            cnt += (str.match(reg) || []).length;
        });
        return cnt;
    }
}
exports.CellInfo = CellInfo;
class TableInfo {
    constructor(settings, range, grid, info) {
        this._settings = settings;
        this._rowInfos = [];
        this._property = {
            isMarkdown: false,
            hasDelimiterAtLineHead: false,
            markdownTableHeaderIndexes: new Set(),
            gridTableHeaderIndexes: new Set(),
            simpleTableHeaderIndexes: new Set()
        };
        this._range = range;
        this._cellGrid = [];
        this._size = { row: 0, col: 0 };
        this.setupCellGrid(grid);
        this.setupSize();
        this.setupProperty(info);
    }
    get property() {
        return this._property;
    }
    get range() {
        return this._range;
    }
    get cellGrid() {
        return this._cellGrid;
    }
    get size() {
        return this._size;
    }
    // サイズの設定
    setupSize() {
        if (!this.isValid())
            return;
        // 正規化済みなのでどの行も同じサイズ
        this._size.row = this._cellGrid.length;
        this._size.col = this._cellGrid[0].length;
    }
    ;
    // 全行のサイズを最大のものに揃える
    setupRowSize() {
        var max = 0;
        this._cellGrid.forEach(row => {
            max = Math.max(max, row.length);
        });
        this._cellGrid.forEach((row, index) => {
            // デリミタも揃える
            var delimiter = (row.length > 0) ? row[0].delimiter : DelimiterType.Pipe;
            for (var c = row.length; c < max; c++) {
                row.push(new CellInfo(this._settings, "", delimiter));
            }
            // 保持
            if (this._rowInfos.length > index)
                this._rowInfos[index].delimiterTYpe = delimiter;
        });
    }
    // セパレータタイプの確定
    setupSeparatorType() {
        this._cellGrid.forEach((row, index) => {
            // セパレータ行かの判定
            var rowType = SeparatorType.None;
            for (var i = 0; i < row.length; i++) {
                var cell = row[i];
                // デリミタがPlusなら初期値をMinusにする（一応この時点でセパレータ行で確定ではあるが特に何もしない）
                if (i == 0 && cell.delimiter == DelimiterType.Plus) {
                    rowType = SeparatorType.Minus;
                }
                // 文字列なら非セパレータ行で確定
                if (cell.type == CellType.CM_Content) {
                    rowType = SeparatorType.None;
                    break;
                }
                // セルのタイプで判定
                if (cell.type == CellType.CM_MinusSeparator) {
                    rowType = SeparatorType.Minus;
                }
                else if (cell.type == CellType.CM_EquallSeparator) {
                    rowType = SeparatorType.Equall;
                }
                else if (cell.type == CellType.MD_LeftSeparator || cell.type == CellType.MD_RightSeparator || cell.type == CellType.MD_CenterSeparator) {
                    rowType = SeparatorType.Minus;
                }
            }
            // セパレータタイプを補正
            row.forEach((cell, i) => {
                // セパレータ行でない場合、セルタイプを文字列に（-や:-を文字として扱う）
                if (rowType == SeparatorType.None) {
                    if (cell.type == CellType.CM_MinusSeparator || cell.type == CellType.CM_EquallSeparator ||
                        cell.type == CellType.MD_LeftSeparator || cell.type == CellType.MD_RightSeparator || cell.type == CellType.MD_CenterSeparator) {
                        cell.setType(CellType.CM_Content);
                    }
                }
                else {
                    if (i != 0 && cell.type == CellType.CM_Blank) {
                        switch (rowType) {
                            case SeparatorType.Minus:
                                cell.setType(CellType.CM_MinusSeparator);
                                break;
                            case SeparatorType.Equall:
                                cell.setType(CellType.CM_EquallSeparator);
                                break;
                        }
                    }
                }
            });
            // 保持
            if (this._rowInfos.length > index)
                this._rowInfos[index].separatorType = rowType;
        });
    }
    // 全セルのサイズを確定
    setupCellSize() {
        // Markdownのセパレータのサイズを設定
        // コンフィグ：左右のスペーサー分がフォーマット時に加算されるため-2する
        let offset = (this._settings.markdown.oneSpacePadding) ? 0 : -2;
        this._cellGrid.forEach(row => {
            row.forEach(cell => {
                if (cell.type == CellType.CM_MinusSeparator || cell.type == CellType.CM_EquallSeparator) {
                    // 最小である3文字にする（---）
                    cell.setSize(3 + offset);
                }
                else if (cell.type == CellType.MD_LeftSeparator || cell.type == CellType.MD_RightSeparator) {
                    // 最小である4文字にする（:---, ---:）
                    cell.setSize(4 + offset);
                }
                else if (cell.type == CellType.MD_CenterSeparator) {
                    // 最小である5文字にする（:---:）
                    cell.setSize(5 + offset);
                }
            });
        });
        // Textileのサイズを設定
        this._cellGrid.forEach(row => {
            // プレフィックス分のパディングを他の行の同列に設定する
            row.forEach((cell, i) => {
                if (cell.type == CellType.TT_HeaderPrefix || cell.type == CellType.TT_LeftPrefix || cell.type == CellType.TT_RightPrefix || cell.type == CellType.TT_CenterPrefix) {
                    this._cellGrid.forEach(elem => {
                        if (i < elem.length) {
                            elem[i].setPadding(2);
                        }
                    });
                }
            });
        });
    }
    // 全セルの位置揃えを確定
    setupCellAlign() {
        // Markdownの位置揃え
        for (var r = this._cellGrid.length - 1; r >= 0; r--) {
            var row = this._cellGrid[r];
            // 各列の位置揃えを変更する
            row.forEach((cell, c) => {
                switch (cell.type) {
                    case CellType.MD_LeftSeparator:
                        this._cellGrid.forEach(elem => {
                            if (c < elem.length) {
                                elem[c].setAlign(CellAlign.Left);
                            }
                        });
                        break;
                    case CellType.MD_RightSeparator:
                        this._cellGrid.forEach(elem => {
                            if (c < elem.length) {
                                elem[c].setAlign(CellAlign.Right);
                            }
                        });
                        break;
                    case CellType.MD_CenterSeparator:
                        this._cellGrid.forEach(elem => {
                            if (c < elem.length) {
                                elem[c].setAlign(CellAlign.Center);
                            }
                        });
                        break;
                }
            });
        }
        // Textileの位置揃え
        this._cellGrid.forEach(row => {
            // 各セルの位置揃えを変更する
            row.forEach(cell => {
                switch (cell.type) {
                    case CellType.TT_LeftPrefix:
                        cell.setAlign(CellAlign.Left);
                        break;
                    case CellType.TT_RightPrefix:
                        cell.setAlign(CellAlign.Right);
                        break;
                    case CellType.TT_CenterPrefix:
                        cell.setAlign(CellAlign.Center);
                        break;
                }
            });
        });
    }
    // セルの内容からプロパティの更新
    setupProperty(info) {
        this._property.isMarkdown = this.isMarkdownTable(this._cellGrid, this._rowInfos);
        this._property.hasDelimiterAtLineHead = info.hasDelimiterAtLineHead;
        this.getMarkdownTableHeaderIndexes(this._property.markdownTableHeaderIndexes, this._cellGrid, this._rowInfos);
        this.getGridTableHeaderIndexes(this._property.gridTableHeaderIndexes, this._cellGrid, this._rowInfos);
        this.getSimpleTableHeaderIndexes(this._property.simpleTableHeaderIndexes, this._cellGrid, this._rowInfos);
    }
    // Markdownかどうか
    isMarkdownTable(grid, rowInfos) {
        if (grid.length < 3 || rowInfos.length < 3)
            return false;
        // 非セパレート行、セパレート行、非セパレート行の順になっているか
        if (rowInfos[0].separatorType != SeparatorType.None || rowInfos[0].delimiterTYpe != DelimiterType.Pipe)
            return false;
        if (rowInfos[1].separatorType != SeparatorType.Minus || rowInfos[1].delimiterTYpe != DelimiterType.Pipe)
            return false;
        if (rowInfos[2].separatorType != SeparatorType.None || rowInfos[2].delimiterTYpe != DelimiterType.Pipe)
            return false;
        return true;
    }
    // Markdownテーブルのヘッダー行取得
    getMarkdownTableHeaderIndexes(outIndexes, grid, rowInfos) {
        outIndexes.clear();
        // １行目以降でセパレータ行までのコンテンツ行をヘッダーとして判定
        for (var i = 0; i < rowInfos.length; i++) {
            var elem = rowInfos[i];
            if (elem.separatorType == SeparatorType.None) {
                // もし最終行ならすべて破棄して終了
                if (i == rowInfos.length - 1) {
                    outIndexes.clear();
                    break;
                }
                outIndexes.add(i);
            }
            else {
                break;
            }
        }
    }
    // Gridテーブルのヘッダー行取得
    getGridTableHeaderIndexes(outIndexes, grid, rowInfos) {
        outIndexes.clear();
        // １行目がMinusセパレータでなければ無視
        if (rowInfos[0].separatorType != SeparatorType.Minus)
            return;
        // ２行目以降でEquallセパレータ行までのコンテンツ行をヘッダーとして判定
        for (var i = 1; i < rowInfos.length; i++) {
            var elem = rowInfos[i];
            if (elem.separatorType == SeparatorType.None) {
                outIndexes.add(i);
            }
            else if (elem.separatorType == SeparatorType.Equall) {
                break;
            }
            else if (elem.separatorType == SeparatorType.Minus) {
                // すべて破棄して終了
                outIndexes.clear();
                break;
            }
        }
    }
    // Simpleテーブルのヘッダー行取得
    getSimpleTableHeaderIndexes(outIndexes, grid, rowInfos) {
        outIndexes.clear();
        // １行目がEquallセパレータでなければ無視
        if (rowInfos[0].separatorType != SeparatorType.Equall)
            return;
        // ２行目以降でEquallセパレータ行までのコンテンツ行をヘッダーとして判定
        for (var i = 1; i < rowInfos.length; i++) {
            var elem = rowInfos[i];
            if (elem.separatorType == SeparatorType.None) {
                outIndexes.add(i);
            }
            else if (elem.separatorType == SeparatorType.Equall) {
                // もし最終行ならすべて破棄して終了
                if (i == rowInfos.length - 1) {
                    outIndexes.clear();
                    break;
                }
                break;
            }
            else if (elem.separatorType == SeparatorType.Minus) {
                continue;
            }
        }
    }
    // 表データを正規化して設定
    setupCellGrid(grid) {
        this._cellGrid = grid;
        this._rowInfos = [];
        for (var i = 0; i < this._cellGrid.length; i++) {
            this._rowInfos.push({
                separatorType: SeparatorType.None,
                delimiterTYpe: DelimiterType.None
            });
        }
        this.setupRowSize();
        this.setupSeparatorType();
        this.setupCellSize();
        this.setupCellAlign();
    }
    ;
    isValid() {
        if (!this._range || this._range.isEmpty)
            return false;
        if (!this._cellGrid || this._cellGrid.length == 0)
            return false;
        var size = this._cellGrid[0].length;
        this._cellGrid.forEach(row => {
            if (!row || row.length == 0 || row.length != size)
                return false;
            row.forEach(cell => {
                if (!cell.isValid())
                    return false;
            });
        });
        return true;
    }
    ;
    getMaxCellSizeList() {
        if (!this.isValid())
            return [];
        var list = [];
        for (var c = 0; c < this._size.col; c++) {
            var max = 0;
            this._cellGrid.forEach((row, r) => {
                if (r == 0)
                    max = row[c].size;
                // １列目は空白列なので小さい方を取る（最も空白が少ない位置に合わせる）
                if (c == 0) {
                    max = Math.min(max, row[c].size);
                }
                else {
                    max = Math.max(max, row[c].size);
                }
            });
            list.push(max);
        }
        return list;
    }
    ;
}
exports.TableInfo = TableInfo;
//# sourceMappingURL=table.js.map