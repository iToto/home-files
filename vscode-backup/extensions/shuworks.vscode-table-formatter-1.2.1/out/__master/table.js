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
    DelimiterType[DelimiterType["Pipe"] = 0] = "Pipe";
    DelimiterType[DelimiterType["Plus"] = 1] = "Plus";
    DelimiterType[DelimiterType["Space"] = 2] = "Space";
})(DelimiterType = exports.DelimiterType || (exports.DelimiterType = {}));
;
// Cell info
class CellInfo {
    constructor(size, delimiter = DelimiterType.Pipe, type = CellType.CM_Blank, align = CellAlign.Left, padding = 0) {
        this._size = size;
        this._delimiter = delimiter;
        this._type = type;
        this._align = align;
        this._padding = padding;
    }
    get size() {
        return this._size;
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
    setSize(size) {
        this._size = size;
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
}
exports.CellInfo = CellInfo;
class TableInfo {
    constructor(range, cells) {
        this._range = range;
        this._cellGrid = cells;
        this._size = this.getSize();
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
    getSize() {
        if (!this.isValid())
            return { row: 0, col: 0 };
        // 正規化済みなのでどの行も同じサイズ
        return {
            row: this._cellGrid.length,
            col: this._cellGrid[0].length
        };
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