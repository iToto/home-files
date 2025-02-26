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
    DelimiterType[DelimiterType["Comma"] = 3] = "Comma";
})(DelimiterType = exports.DelimiterType || (exports.DelimiterType = {}));
;
// Separator type
var SeparatorType;
(function (SeparatorType) {
    SeparatorType[SeparatorType["None"] = 0] = "None";
    SeparatorType[SeparatorType["Minus"] = 1] = "Minus";
    SeparatorType[SeparatorType["Equall"] = 2] = "Equall";
})(SeparatorType || (SeparatorType = {}));
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
    constructor(range, grid) {
        this._range = range;
        this.setupCellGrid(grid);
        this.setupSize();
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
        this._cellGrid.forEach(row => {
            // デリミタも揃える
            var delimiter = (row.length > 0) ? row[0].delimiter : DelimiterType.Pipe;
            for (var c = row.length; c < max; c++) {
                row.push(new CellInfo(0, delimiter));
            }
        });
    }
    // セパレータタイプの確定
    setupSeparatorType() {
        this._cellGrid.forEach(row => {
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
        });
    }
    // 全セルのサイズを確定
    setupCellSize() {
        // Markdownのセパレータのサイズを設定
        this._cellGrid.forEach(row => {
            row.forEach(cell => {
                if (cell.type == CellType.CM_MinusSeparator ||
                    cell.type == CellType.MD_LeftSeparator || cell.type == CellType.MD_RightSeparator || cell.type == CellType.MD_CenterSeparator) {
                    // 最小である3文字にする
                    cell.setSize(3);
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
    // 表データを正規化して設定
    setupCellGrid(grid) {
        this._cellGrid = grid;
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