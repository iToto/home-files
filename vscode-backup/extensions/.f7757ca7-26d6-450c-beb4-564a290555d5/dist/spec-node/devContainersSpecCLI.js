/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Microsoft Corporation. All rights reserved.
 *  Licensed under the MIT License. See License.txt in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __esm = (fn, res) => function __init() {
  return fn && (res = (0, fn[__getOwnPropNames(fn)[0]])(fn = 0)), res;
};
var __commonJS = (cb, mod) => function __require() {
  return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
};
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
  // If the importer is in node compatibility mode or this is not an ESM
  // file that has been converted to a CommonJS file using a Babel-
  // compatible transform (i.e. "__esModule" has not been set), then set
  // "default" to the CommonJS "module.exports" for node compatibility.
  isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
  mod
));
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// node_modules/text-table/index.js
var require_text_table = __commonJS({
  "node_modules/text-table/index.js"(exports, module2) {
    module2.exports = function(rows_, opts) {
      if (!opts)
        opts = {};
      var hsep = opts.hsep === void 0 ? "  " : opts.hsep;
      var align2 = opts.align || [];
      var stringLength = opts.stringLength || function(s) {
        return String(s).length;
      };
      var dotsizes = reduce(rows_, function(acc, row) {
        forEach(row, function(c2, ix) {
          var n = dotindex(c2);
          if (!acc[ix] || n > acc[ix])
            acc[ix] = n;
        });
        return acc;
      }, []);
      var rows = map2(rows_, function(row) {
        return map2(row, function(c_, ix) {
          var c2 = String(c_);
          if (align2[ix] === ".") {
            var index = dotindex(c2);
            var size = dotsizes[ix] + (/\./.test(c2) ? 1 : 2) - (stringLength(c2) - index);
            return c2 + Array(size).join(" ");
          } else
            return c2;
        });
      });
      var sizes = reduce(rows, function(acc, row) {
        forEach(row, function(c2, ix) {
          var n = stringLength(c2);
          if (!acc[ix] || n > acc[ix])
            acc[ix] = n;
        });
        return acc;
      }, []);
      return map2(rows, function(row) {
        return map2(row, function(c2, ix) {
          var n = sizes[ix] - stringLength(c2) || 0;
          var s = Array(Math.max(n + 1, 1)).join(" ");
          if (align2[ix] === "r" || align2[ix] === ".") {
            return s + c2;
          }
          if (align2[ix] === "c") {
            return Array(Math.ceil(n / 2 + 1)).join(" ") + c2 + Array(Math.floor(n / 2 + 1)).join(" ");
          }
          return c2 + s;
        }).join(hsep).replace(/\s+$/, "");
      }).join("\n");
    };
    function dotindex(c2) {
      var m = /\.[^.]*$/.exec(c2);
      return m ? m.index + 1 : c2.length;
    }
    function reduce(xs, f, init) {
      if (xs.reduce)
        return xs.reduce(f, init);
      var i = 0;
      var acc = arguments.length >= 3 ? init : xs[i++];
      for (; i < xs.length; i++) {
        f(acc, xs[i], i);
      }
      return acc;
    }
    function forEach(xs, f) {
      if (xs.forEach)
        return xs.forEach(f);
      for (var i = 0; i < xs.length; i++) {
        f.call(xs, xs[i], i);
      }
    }
    function map2(xs, f) {
      if (xs.map)
        return xs.map(f);
      var res = [];
      for (var i = 0; i < xs.length; i++) {
        res.push(f.call(xs, xs[i], i));
      }
      return res;
    }
  }
});

// node_modules/ncp/lib/ncp.js
var require_ncp = __commonJS({
  "node_modules/ncp/lib/ncp.js"(exports, module2) {
    var fs6 = require("fs");
    var path32 = require("path");
    module2.exports = ncp3;
    ncp3.ncp = ncp3;
    function ncp3(source, dest, options, callback) {
      var cback = callback;
      if (!callback) {
        cback = options;
        options = {};
      }
      var basePath = process.cwd(), currentPath = path32.resolve(basePath, source), targetPath = path32.resolve(basePath, dest), filter = options.filter, rename2 = options.rename, transform = options.transform, clobber = options.clobber !== false, modified = options.modified, dereference = options.dereference, errs = null, started = 0, finished = 0, running = 0, limit = options.limit || ncp3.limit || 16;
      limit = limit < 1 ? 1 : limit > 512 ? 512 : limit;
      startCopy(currentPath);
      function startCopy(source2) {
        started++;
        if (filter) {
          if (filter instanceof RegExp) {
            if (!filter.test(source2)) {
              return cb(true);
            }
          } else if (typeof filter === "function") {
            if (!filter(source2)) {
              return cb(true);
            }
          }
        }
        return getStats(source2);
      }
      function getStats(source2) {
        var stat2 = dereference ? fs6.stat : fs6.lstat;
        if (running >= limit) {
          return setImmediate(function() {
            getStats(source2);
          });
        }
        running++;
        stat2(source2, function(err, stats) {
          var item = {};
          if (err) {
            return onError(err);
          }
          item.name = source2;
          item.mode = stats.mode;
          item.mtime = stats.mtime;
          item.atime = stats.atime;
          if (stats.isDirectory()) {
            return onDir(item);
          } else if (stats.isFile()) {
            return onFile(item);
          } else if (stats.isSymbolicLink()) {
            return onLink(source2);
          }
        });
      }
      function onFile(file) {
        var target = file.name.replace(currentPath, targetPath);
        if (rename2) {
          target = rename2(target);
        }
        isWritable(target, function(writable) {
          if (writable) {
            return copyFile2(file, target);
          }
          if (clobber) {
            rmFile(target, function() {
              copyFile2(file, target);
            });
          }
          if (modified) {
            var stat2 = dereference ? fs6.stat : fs6.lstat;
            stat2(target, function(err, stats) {
              if (file.mtime.getTime() > stats.mtime.getTime())
                copyFile2(file, target);
              else
                return cb();
            });
          } else {
            return cb();
          }
        });
      }
      function copyFile2(file, target) {
        var readStream = fs6.createReadStream(file.name), writeStream = fs6.createWriteStream(target, { mode: file.mode });
        readStream.on("error", onError);
        writeStream.on("error", onError);
        if (transform) {
          transform(readStream, writeStream, file);
        } else {
          writeStream.on("open", function() {
            readStream.pipe(writeStream);
          });
        }
        writeStream.once("finish", function() {
          if (modified) {
            fs6.utimesSync(target, file.atime, file.mtime);
            cb();
          } else
            cb();
        });
      }
      function rmFile(file, done) {
        fs6.unlink(file, function(err) {
          if (err) {
            return onError(err);
          }
          return done();
        });
      }
      function onDir(dir) {
        var target = dir.name.replace(currentPath, targetPath);
        isWritable(target, function(writable) {
          if (writable) {
            return mkDir(dir, target);
          }
          copyDir(dir.name);
        });
      }
      function mkDir(dir, target) {
        fs6.mkdir(target, dir.mode, function(err) {
          if (err) {
            return onError(err);
          }
          copyDir(dir.name);
        });
      }
      function copyDir(dir) {
        fs6.readdir(dir, function(err, items) {
          if (err) {
            return onError(err);
          }
          items.forEach(function(item) {
            startCopy(path32.join(dir, item));
          });
          return cb();
        });
      }
      function onLink(link) {
        var target = link.replace(currentPath, targetPath);
        fs6.readlink(link, function(err, resolvedPath) {
          if (err) {
            return onError(err);
          }
          checkLink(resolvedPath, target);
        });
      }
      function checkLink(resolvedPath, target) {
        if (dereference) {
          resolvedPath = path32.resolve(basePath, resolvedPath);
        }
        isWritable(target, function(writable) {
          if (writable) {
            return makeLink(resolvedPath, target);
          }
          fs6.readlink(target, function(err, targetDest) {
            if (err) {
              return onError(err);
            }
            if (dereference) {
              targetDest = path32.resolve(basePath, targetDest);
            }
            if (targetDest === resolvedPath) {
              return cb();
            }
            return rmFile(target, function() {
              makeLink(resolvedPath, target);
            });
          });
        });
      }
      function makeLink(linkPath, target) {
        fs6.symlink(linkPath, target, function(err) {
          if (err) {
            return onError(err);
          }
          return cb();
        });
      }
      function isWritable(path33, done) {
        fs6.lstat(path33, function(err) {
          if (err) {
            if (err.code === "ENOENT")
              return done(true);
            return done(false);
          }
          return done(false);
        });
      }
      function onError(err) {
        if (options.stopOnError) {
          return cback(err);
        } else if (!errs && options.errs) {
          errs = fs6.createWriteStream(options.errs);
        } else if (!errs) {
          errs = [];
        }
        if (typeof errs.write === "undefined") {
          errs.push(err);
        } else {
          errs.write(err.stack + "\n\n");
        }
        return cb();
      }
      function cb(skipped) {
        if (!skipped)
          running--;
        finished++;
        if (started === finished && running === 0) {
          if (cback !== void 0) {
            return errs ? cback(errs) : cback(null);
          }
        }
      }
    }
  }
});

// node_modules/pull-stream/pull.js
var require_pull = __commonJS({
  "node_modules/pull-stream/pull.js"(exports, module2) {
    "use strict";
    module2.exports = function pull(a) {
      var length = arguments.length;
      if (typeof a === "function" && a.length === 1) {
        var args = new Array(length);
        for (var i = 0; i < length; i++)
          args[i] = arguments[i];
        return function(read3) {
          if (args == null) {
            throw new TypeError("partial sink should only be called once!");
          }
          var ref = args;
          args = null;
          switch (length) {
            case 1:
              return pull(read3, ref[0]);
            case 2:
              return pull(read3, ref[0], ref[1]);
            case 3:
              return pull(read3, ref[0], ref[1], ref[2]);
            case 4:
              return pull(read3, ref[0], ref[1], ref[2], ref[3]);
            default:
              ref.unshift(read3);
              return pull.apply(null, ref);
          }
        };
      }
      var read2 = a;
      if (read2 && typeof read2.source === "function") {
        read2 = read2.source;
      }
      for (var i = 1; i < length; i++) {
        var s = arguments[i];
        if (typeof s === "function") {
          read2 = s(read2);
        } else if (s && typeof s === "object") {
          s.sink(read2);
          read2 = s.source;
        }
      }
      return read2;
    };
  }
});

// node_modules/looper/index.js
var require_looper = __commonJS({
  "node_modules/looper/index.js"(exports, module2) {
    var looper = module2.exports = function(fun) {
      (function next() {
        var loop = true, returned = false, sync = false;
        do {
          sync = true;
          loop = false;
          fun.call(this, function() {
            if (sync)
              loop = true;
            else
              next();
          });
          sync = false;
        } while (loop);
      })();
    };
  }
});

// node_modules/stream-to-pull-stream/index.js
var require_stream_to_pull_stream = __commonJS({
  "node_modules/stream-to-pull-stream/index.js"(exports, module2) {
    var pull = require_pull();
    var looper = require_looper();
    function destroy(stream) {
      if (!stream.destroy)
        console.error(
          "warning, stream-to-pull-stream: \nthe wrapped node-stream does not implement `destroy`, \nthis may cause resource leaks."
        );
      else
        stream.destroy();
    }
    function write(read4, stream, cb) {
      var ended, closed = false, did;
      function done() {
        if (did)
          return;
        did = true;
        cb && cb(ended === true ? null : ended);
      }
      function onClose() {
        if (closed)
          return;
        closed = true;
        cleanup2();
        if (!ended)
          read4(ended = true, done);
        else
          done();
      }
      function onError(err) {
        cleanup2();
        if (!ended)
          read4(ended = err, done);
      }
      function cleanup2() {
        stream.on("finish", onClose);
        stream.removeListener("close", onClose);
        stream.removeListener("error", onError);
      }
      stream.on("close", onClose);
      stream.on("finish", onClose);
      stream.on("error", onError);
      process.nextTick(function() {
        looper(function(next) {
          read4(null, function(end, data) {
            ended = ended || end;
            if (end === true)
              return stream._isStdio ? done() : stream.end();
            if (ended = ended || end) {
              destroy(stream);
              return done(ended);
            }
            if (stream._isStdio)
              stream.write(data, function() {
                next();
              });
            else {
              var pause = stream.write(data);
              if (pause === false)
                stream.once("drain", next);
              else
                next();
            }
          });
        });
      });
    }
    function read2(stream) {
      var ended = false, waiting = false;
      var _cb;
      function read4() {
        var data = stream.read();
        if (data !== null && _cb) {
          var cb = _cb;
          _cb = null;
          cb(null, data);
        }
      }
      stream.on("readable", function() {
        waiting = true;
        _cb && read4();
      }).on("end", function() {
        ended = true;
        _cb && _cb(ended);
      }).on("error", function(err) {
        ended = err;
        _cb && _cb(ended);
      });
      return function(end, cb) {
        _cb = cb;
        if (ended)
          cb(ended);
        else if (waiting)
          read4();
      };
    }
    function read1(stream) {
      var buffer = [], cbs = [], ended, paused = false;
      var draining;
      function drain() {
        while ((buffer.length || ended) && cbs.length)
          cbs.shift()(buffer.length ? null : ended, buffer.shift());
        if (!buffer.length && paused) {
          paused = false;
          stream.resume();
        }
      }
      stream.on("data", function(data) {
        buffer.push(data);
        drain();
        if (buffer.length && stream.pause) {
          paused = true;
          stream.pause();
        }
      });
      stream.on("end", function() {
        ended = true;
        drain();
      });
      stream.on("close", function() {
        ended = true;
        drain();
      });
      stream.on("error", function(err) {
        ended = err;
        drain();
      });
      return function(abort, cb) {
        if (!cb)
          throw new Error("*must* provide cb");
        if (abort) {
          let onAbort2 = function() {
            while (cbs.length)
              cbs.shift()(abort);
            cb(abort);
          };
          var onAbort = onAbort2;
          if (ended)
            return onAbort2();
          stream.once("close", onAbort2);
          destroy(stream);
        } else {
          cbs.push(cb);
          drain();
        }
      };
    }
    var read3 = read1;
    var sink = function(stream, cb) {
      return function(read4) {
        return write(read4, stream, cb);
      };
    };
    var source = function(stream) {
      return read1(stream);
    };
    exports = module2.exports = function(stream, cb) {
      return stream.writable && stream.write ? stream.readable ? function(_read) {
        write(_read, stream, cb);
        return read1(stream);
      } : sink(stream, cb) : source(stream);
    };
    exports.sink = sink;
    exports.source = source;
    exports.read = read3;
    exports.read1 = read1;
    exports.read2 = read2;
    exports.duplex = function(stream, cb) {
      return {
        source: source(stream),
        sink: sink(stream, cb)
      };
    };
    exports.transform = function(stream) {
      return function(read4) {
        var _source = source(stream);
        sink(stream)(read4);
        return _source;
      };
    };
  }
});

// node_modules/shell-quote/quote.js
var require_quote = __commonJS({
  "node_modules/shell-quote/quote.js"(exports, module2) {
    "use strict";
    module2.exports = function quote(xs) {
      return xs.map(function(s) {
        if (s && typeof s === "object") {
          return s.op.replace(/(.)/g, "\\$1");
        }
        if (/["\s]/.test(s) && !/'/.test(s)) {
          return "'" + s.replace(/(['\\])/g, "\\$1") + "'";
        }
        if (/["'\s]/.test(s)) {
          return '"' + s.replace(/(["\\$`!])/g, "\\$1") + '"';
        }
        return String(s).replace(/([A-Za-z]:)?([#!"$&'()*,:;<=>?@[\\\]^`{|}])/g, "$1\\$2");
      }).join(" ");
    };
  }
});

// node_modules/shell-quote/parse.js
var require_parse = __commonJS({
  "node_modules/shell-quote/parse.js"(exports, module2) {
    "use strict";
    var CONTROL = "(?:" + [
      "\\|\\|",
      "\\&\\&",
      ";;",
      "\\|\\&",
      "\\<\\(",
      "\\<\\<\\<",
      ">>",
      ">\\&",
      "<\\&",
      "[&;()|<>]"
    ].join("|") + ")";
    var META = "|&;()<> \\t";
    var BAREWORD = `(\\\\['"` + META + `]|[^\\s'"` + META + "])+";
    var SINGLE_QUOTE = '"((\\\\"|[^"])*?)"';
    var DOUBLE_QUOTE = "'((\\\\'|[^'])*?)'";
    var TOKEN = "";
    for (i = 0; i < 4; i++) {
      TOKEN += (Math.pow(16, 8) * Math.random()).toString(16);
    }
    var i;
    function parseInternal(s, env3, opts) {
      var chunker = new RegExp([
        "(" + CONTROL + ")",
        // control chars
        "(" + BAREWORD + "|" + SINGLE_QUOTE + "|" + DOUBLE_QUOTE + ")*"
      ].join("|"), "g");
      var match = s.match(chunker).filter(Boolean);
      if (!match) {
        return [];
      }
      if (!env3) {
        env3 = {};
      }
      if (!opts) {
        opts = {};
      }
      var commented = false;
      function getVar(_, pre, key) {
        var r = typeof env3 === "function" ? env3(key) : env3[key];
        if (r === void 0 && key != "") {
          r = "";
        } else if (r === void 0) {
          r = "$";
        }
        if (typeof r === "object") {
          return pre + TOKEN + JSON.stringify(r) + TOKEN;
        }
        return pre + r;
      }
      return match.map(function(s2, j) {
        if (commented) {
          return void 0;
        }
        if (RegExp("^" + CONTROL + "$").test(s2)) {
          return { op: s2 };
        }
        var SQ = "'";
        var DQ = '"';
        var DS = "$";
        var BS = opts.escape || "\\";
        var quote = false;
        var esc = false;
        var out = "";
        var isGlob = false;
        var i2;
        function parseEnvVar() {
          i2 += 1;
          var varend;
          var varname;
          if (s2.charAt(i2) === "{") {
            i2 += 1;
            if (s2.charAt(i2) === "}") {
              throw new Error("Bad substitution: " + s2.substr(i2 - 2, 3));
            }
            varend = s2.indexOf("}", i2);
            if (varend < 0) {
              throw new Error("Bad substitution: " + s2.substr(i2));
            }
            varname = s2.substr(i2, varend - i2);
            i2 = varend;
          } else if (/[*@#?$!_-]/.test(s2.charAt(i2))) {
            varname = s2.charAt(i2);
            i2 += 1;
          } else {
            varend = s2.substr(i2).match(/[^\w\d_]/);
            if (!varend) {
              varname = s2.substr(i2);
              i2 = s2.length;
            } else {
              varname = s2.substr(i2, varend.index);
              i2 += varend.index - 1;
            }
          }
          return getVar(null, "", varname);
        }
        for (i2 = 0; i2 < s2.length; i2++) {
          var c2 = s2.charAt(i2);
          isGlob = isGlob || !quote && (c2 === "*" || c2 === "?");
          if (esc) {
            out += c2;
            esc = false;
          } else if (quote) {
            if (c2 === quote) {
              quote = false;
            } else if (quote == SQ) {
              out += c2;
            } else {
              if (c2 === BS) {
                i2 += 1;
                c2 = s2.charAt(i2);
                if (c2 === DQ || c2 === BS || c2 === DS) {
                  out += c2;
                } else {
                  out += BS + c2;
                }
              } else if (c2 === DS) {
                out += parseEnvVar();
              } else {
                out += c2;
              }
            }
          } else if (c2 === DQ || c2 === SQ) {
            quote = c2;
          } else if (RegExp("^" + CONTROL + "$").test(c2)) {
            return { op: s2 };
          } else if (/^#$/.test(c2)) {
            commented = true;
            if (out.length) {
              return [out, { comment: s2.slice(i2 + 1) + match.slice(j + 1).join(" ") }];
            }
            return [{ comment: s2.slice(i2 + 1) + match.slice(j + 1).join(" ") }];
          } else if (c2 === BS) {
            esc = true;
          } else if (c2 === DS) {
            out += parseEnvVar();
          } else {
            out += c2;
          }
        }
        if (isGlob) {
          return { op: "glob", pattern: out };
        }
        return out;
      }).reduce(function(prev, arg) {
        if (arg === void 0) {
          return prev;
        }
        return prev.concat(arg);
      }, []);
    }
    module2.exports = function parse5(s, env3, opts) {
      var mapped = parseInternal(s, env3, opts);
      if (typeof env3 !== "function") {
        return mapped;
      }
      return mapped.reduce(function(acc, s2) {
        if (typeof s2 === "object") {
          return acc.concat(s2);
        }
        var xs = s2.split(RegExp("(" + TOKEN + ".*?" + TOKEN + ")", "g"));
        if (xs.length === 1) {
          return acc.concat(xs[0]);
        }
        return acc.concat(xs.filter(Boolean).map(function(x3) {
          if (RegExp("^" + TOKEN).test(x3)) {
            return JSON.parse(x3.split(TOKEN)[1]);
          }
          return x3;
        }));
      }, []);
    };
  }
});

// node_modules/shell-quote/index.js
var require_shell_quote = __commonJS({
  "node_modules/shell-quote/index.js"(exports) {
    "use strict";
    exports.quote = require_quote();
    exports.parse = require_parse();
  }
});

// node_modules/tar/lib/high-level-opt.js
var require_high_level_opt = __commonJS({
  "node_modules/tar/lib/high-level-opt.js"(exports, module2) {
    "use strict";
    var argmap = /* @__PURE__ */ new Map([
      ["C", "cwd"],
      ["f", "file"],
      ["z", "gzip"],
      ["P", "preservePaths"],
      ["U", "unlink"],
      ["strip-components", "strip"],
      ["stripComponents", "strip"],
      ["keep-newer", "newer"],
      ["keepNewer", "newer"],
      ["keep-newer-files", "newer"],
      ["keepNewerFiles", "newer"],
      ["k", "keep"],
      ["keep-existing", "keep"],
      ["keepExisting", "keep"],
      ["m", "noMtime"],
      ["no-mtime", "noMtime"],
      ["p", "preserveOwner"],
      ["L", "follow"],
      ["h", "follow"]
    ]);
    module2.exports = (opt) => opt ? Object.keys(opt).map((k) => [
      argmap.has(k) ? argmap.get(k) : k,
      opt[k]
    ]).reduce((set2, kv) => (set2[kv[0]] = kv[1], set2), /* @__PURE__ */ Object.create(null)) : {};
  }
});

// node_modules/minipass/index.js
var require_minipass = __commonJS({
  "node_modules/minipass/index.js"(exports, module2) {
    "use strict";
    var proc = typeof process === "object" && process ? process : {
      stdout: null,
      stderr: null
    };
    var EE = require("events");
    var Stream = require("stream");
    var stringdecoder = require("string_decoder");
    var SD = stringdecoder.StringDecoder;
    var EOF = Symbol("EOF");
    var MAYBE_EMIT_END = Symbol("maybeEmitEnd");
    var EMITTED_END = Symbol("emittedEnd");
    var EMITTING_END = Symbol("emittingEnd");
    var EMITTED_ERROR = Symbol("emittedError");
    var CLOSED = Symbol("closed");
    var READ = Symbol("read");
    var FLUSH = Symbol("flush");
    var FLUSHCHUNK = Symbol("flushChunk");
    var ENCODING = Symbol("encoding");
    var DECODER = Symbol("decoder");
    var FLOWING = Symbol("flowing");
    var PAUSED = Symbol("paused");
    var RESUME = Symbol("resume");
    var BUFFER = Symbol("buffer");
    var PIPES = Symbol("pipes");
    var BUFFERLENGTH = Symbol("bufferLength");
    var BUFFERPUSH = Symbol("bufferPush");
    var BUFFERSHIFT = Symbol("bufferShift");
    var OBJECTMODE = Symbol("objectMode");
    var DESTROYED = Symbol("destroyed");
    var ERROR = Symbol("error");
    var EMITDATA = Symbol("emitData");
    var EMITEND = Symbol("emitEnd");
    var EMITEND2 = Symbol("emitEnd2");
    var ASYNC = Symbol("async");
    var ABORT = Symbol("abort");
    var ABORTED = Symbol("aborted");
    var SIGNAL = Symbol("signal");
    var defer = (fn) => Promise.resolve().then(fn);
    var doIter = global._MP_NO_ITERATOR_SYMBOLS_ !== "1";
    var ASYNCITERATOR = doIter && Symbol.asyncIterator || Symbol("asyncIterator not implemented");
    var ITERATOR = doIter && Symbol.iterator || Symbol("iterator not implemented");
    var isEndish = (ev) => ev === "end" || ev === "finish" || ev === "prefinish";
    var isArrayBuffer = (b) => b instanceof ArrayBuffer || typeof b === "object" && b.constructor && b.constructor.name === "ArrayBuffer" && b.byteLength >= 0;
    var isArrayBufferView = (b) => !Buffer.isBuffer(b) && ArrayBuffer.isView(b);
    var Pipe = class {
      constructor(src, dest, opts) {
        this.src = src;
        this.dest = dest;
        this.opts = opts;
        this.ondrain = () => src[RESUME]();
        dest.on("drain", this.ondrain);
      }
      unpipe() {
        this.dest.removeListener("drain", this.ondrain);
      }
      // istanbul ignore next - only here for the prototype
      proxyErrors() {
      }
      end() {
        this.unpipe();
        if (this.opts.end)
          this.dest.end();
      }
    };
    var PipeProxyErrors = class extends Pipe {
      unpipe() {
        this.src.removeListener("error", this.proxyErrors);
        super.unpipe();
      }
      constructor(src, dest, opts) {
        super(src, dest, opts);
        this.proxyErrors = (er) => dest.emit("error", er);
        src.on("error", this.proxyErrors);
      }
    };
    var Minipass = class extends Stream {
      constructor(options) {
        super();
        this[FLOWING] = false;
        this[PAUSED] = false;
        this[PIPES] = [];
        this[BUFFER] = [];
        this[OBJECTMODE] = options && options.objectMode || false;
        if (this[OBJECTMODE])
          this[ENCODING] = null;
        else
          this[ENCODING] = options && options.encoding || null;
        if (this[ENCODING] === "buffer")
          this[ENCODING] = null;
        this[ASYNC] = options && !!options.async || false;
        this[DECODER] = this[ENCODING] ? new SD(this[ENCODING]) : null;
        this[EOF] = false;
        this[EMITTED_END] = false;
        this[EMITTING_END] = false;
        this[CLOSED] = false;
        this[EMITTED_ERROR] = null;
        this.writable = true;
        this.readable = true;
        this[BUFFERLENGTH] = 0;
        this[DESTROYED] = false;
        if (options && options.debugExposeBuffer === true) {
          Object.defineProperty(this, "buffer", { get: () => this[BUFFER] });
        }
        if (options && options.debugExposePipes === true) {
          Object.defineProperty(this, "pipes", { get: () => this[PIPES] });
        }
        this[SIGNAL] = options && options.signal;
        this[ABORTED] = false;
        if (this[SIGNAL]) {
          this[SIGNAL].addEventListener("abort", () => this[ABORT]());
          if (this[SIGNAL].aborted) {
            this[ABORT]();
          }
        }
      }
      get bufferLength() {
        return this[BUFFERLENGTH];
      }
      get encoding() {
        return this[ENCODING];
      }
      set encoding(enc) {
        if (this[OBJECTMODE])
          throw new Error("cannot set encoding in objectMode");
        if (this[ENCODING] && enc !== this[ENCODING] && (this[DECODER] && this[DECODER].lastNeed || this[BUFFERLENGTH]))
          throw new Error("cannot change encoding");
        if (this[ENCODING] !== enc) {
          this[DECODER] = enc ? new SD(enc) : null;
          if (this[BUFFER].length)
            this[BUFFER] = this[BUFFER].map((chunk) => this[DECODER].write(chunk));
        }
        this[ENCODING] = enc;
      }
      setEncoding(enc) {
        this.encoding = enc;
      }
      get objectMode() {
        return this[OBJECTMODE];
      }
      set objectMode(om) {
        this[OBJECTMODE] = this[OBJECTMODE] || !!om;
      }
      get ["async"]() {
        return this[ASYNC];
      }
      set ["async"](a) {
        this[ASYNC] = this[ASYNC] || !!a;
      }
      // drop everything and get out of the flow completely
      [ABORT]() {
        this[ABORTED] = true;
        this.emit("abort", this[SIGNAL].reason);
        this.destroy(this[SIGNAL].reason);
      }
      get aborted() {
        return this[ABORTED];
      }
      set aborted(_) {
      }
      write(chunk, encoding, cb) {
        if (this[ABORTED])
          return false;
        if (this[EOF])
          throw new Error("write after end");
        if (this[DESTROYED]) {
          this.emit(
            "error",
            Object.assign(
              new Error("Cannot call write after a stream was destroyed"),
              { code: "ERR_STREAM_DESTROYED" }
            )
          );
          return true;
        }
        if (typeof encoding === "function")
          cb = encoding, encoding = "utf8";
        if (!encoding)
          encoding = "utf8";
        const fn = this[ASYNC] ? defer : (f) => f();
        if (!this[OBJECTMODE] && !Buffer.isBuffer(chunk)) {
          if (isArrayBufferView(chunk))
            chunk = Buffer.from(chunk.buffer, chunk.byteOffset, chunk.byteLength);
          else if (isArrayBuffer(chunk))
            chunk = Buffer.from(chunk);
          else if (typeof chunk !== "string")
            this.objectMode = true;
        }
        if (this[OBJECTMODE]) {
          if (this.flowing && this[BUFFERLENGTH] !== 0)
            this[FLUSH](true);
          if (this.flowing)
            this.emit("data", chunk);
          else
            this[BUFFERPUSH](chunk);
          if (this[BUFFERLENGTH] !== 0)
            this.emit("readable");
          if (cb)
            fn(cb);
          return this.flowing;
        }
        if (!chunk.length) {
          if (this[BUFFERLENGTH] !== 0)
            this.emit("readable");
          if (cb)
            fn(cb);
          return this.flowing;
        }
        if (typeof chunk === "string" && // unless it is a string already ready for us to use
        !(encoding === this[ENCODING] && !this[DECODER].lastNeed)) {
          chunk = Buffer.from(chunk, encoding);
        }
        if (Buffer.isBuffer(chunk) && this[ENCODING])
          chunk = this[DECODER].write(chunk);
        if (this.flowing && this[BUFFERLENGTH] !== 0)
          this[FLUSH](true);
        if (this.flowing)
          this.emit("data", chunk);
        else
          this[BUFFERPUSH](chunk);
        if (this[BUFFERLENGTH] !== 0)
          this.emit("readable");
        if (cb)
          fn(cb);
        return this.flowing;
      }
      read(n) {
        if (this[DESTROYED])
          return null;
        if (this[BUFFERLENGTH] === 0 || n === 0 || n > this[BUFFERLENGTH]) {
          this[MAYBE_EMIT_END]();
          return null;
        }
        if (this[OBJECTMODE])
          n = null;
        if (this[BUFFER].length > 1 && !this[OBJECTMODE]) {
          if (this.encoding)
            this[BUFFER] = [this[BUFFER].join("")];
          else
            this[BUFFER] = [Buffer.concat(this[BUFFER], this[BUFFERLENGTH])];
        }
        const ret = this[READ](n || null, this[BUFFER][0]);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [READ](n, chunk) {
        if (n === chunk.length || n === null)
          this[BUFFERSHIFT]();
        else {
          this[BUFFER][0] = chunk.slice(n);
          chunk = chunk.slice(0, n);
          this[BUFFERLENGTH] -= n;
        }
        this.emit("data", chunk);
        if (!this[BUFFER].length && !this[EOF])
          this.emit("drain");
        return chunk;
      }
      end(chunk, encoding, cb) {
        if (typeof chunk === "function")
          cb = chunk, chunk = null;
        if (typeof encoding === "function")
          cb = encoding, encoding = "utf8";
        if (chunk)
          this.write(chunk, encoding);
        if (cb)
          this.once("end", cb);
        this[EOF] = true;
        this.writable = false;
        if (this.flowing || !this[PAUSED])
          this[MAYBE_EMIT_END]();
        return this;
      }
      // don't let the internal resume be overwritten
      [RESUME]() {
        if (this[DESTROYED])
          return;
        this[PAUSED] = false;
        this[FLOWING] = true;
        this.emit("resume");
        if (this[BUFFER].length)
          this[FLUSH]();
        else if (this[EOF])
          this[MAYBE_EMIT_END]();
        else
          this.emit("drain");
      }
      resume() {
        return this[RESUME]();
      }
      pause() {
        this[FLOWING] = false;
        this[PAUSED] = true;
      }
      get destroyed() {
        return this[DESTROYED];
      }
      get flowing() {
        return this[FLOWING];
      }
      get paused() {
        return this[PAUSED];
      }
      [BUFFERPUSH](chunk) {
        if (this[OBJECTMODE])
          this[BUFFERLENGTH] += 1;
        else
          this[BUFFERLENGTH] += chunk.length;
        this[BUFFER].push(chunk);
      }
      [BUFFERSHIFT]() {
        if (this[OBJECTMODE])
          this[BUFFERLENGTH] -= 1;
        else
          this[BUFFERLENGTH] -= this[BUFFER][0].length;
        return this[BUFFER].shift();
      }
      [FLUSH](noDrain) {
        do {
        } while (this[FLUSHCHUNK](this[BUFFERSHIFT]()) && this[BUFFER].length);
        if (!noDrain && !this[BUFFER].length && !this[EOF])
          this.emit("drain");
      }
      [FLUSHCHUNK](chunk) {
        this.emit("data", chunk);
        return this.flowing;
      }
      pipe(dest, opts) {
        if (this[DESTROYED])
          return;
        const ended = this[EMITTED_END];
        opts = opts || {};
        if (dest === proc.stdout || dest === proc.stderr)
          opts.end = false;
        else
          opts.end = opts.end !== false;
        opts.proxyErrors = !!opts.proxyErrors;
        if (ended) {
          if (opts.end)
            dest.end();
        } else {
          this[PIPES].push(
            !opts.proxyErrors ? new Pipe(this, dest, opts) : new PipeProxyErrors(this, dest, opts)
          );
          if (this[ASYNC])
            defer(() => this[RESUME]());
          else
            this[RESUME]();
        }
        return dest;
      }
      unpipe(dest) {
        const p = this[PIPES].find((p2) => p2.dest === dest);
        if (p) {
          this[PIPES].splice(this[PIPES].indexOf(p), 1);
          p.unpipe();
        }
      }
      addListener(ev, fn) {
        return this.on(ev, fn);
      }
      on(ev, fn) {
        const ret = super.on(ev, fn);
        if (ev === "data" && !this[PIPES].length && !this.flowing)
          this[RESUME]();
        else if (ev === "readable" && this[BUFFERLENGTH] !== 0)
          super.emit("readable");
        else if (isEndish(ev) && this[EMITTED_END]) {
          super.emit(ev);
          this.removeAllListeners(ev);
        } else if (ev === "error" && this[EMITTED_ERROR]) {
          if (this[ASYNC])
            defer(() => fn.call(this, this[EMITTED_ERROR]));
          else
            fn.call(this, this[EMITTED_ERROR]);
        }
        return ret;
      }
      get emittedEnd() {
        return this[EMITTED_END];
      }
      [MAYBE_EMIT_END]() {
        if (!this[EMITTING_END] && !this[EMITTED_END] && !this[DESTROYED] && this[BUFFER].length === 0 && this[EOF]) {
          this[EMITTING_END] = true;
          this.emit("end");
          this.emit("prefinish");
          this.emit("finish");
          if (this[CLOSED])
            this.emit("close");
          this[EMITTING_END] = false;
        }
      }
      emit(ev, data, ...extra) {
        if (ev !== "error" && ev !== "close" && ev !== DESTROYED && this[DESTROYED])
          return;
        else if (ev === "data") {
          return !this[OBJECTMODE] && !data ? false : this[ASYNC] ? defer(() => this[EMITDATA](data)) : this[EMITDATA](data);
        } else if (ev === "end") {
          return this[EMITEND]();
        } else if (ev === "close") {
          this[CLOSED] = true;
          if (!this[EMITTED_END] && !this[DESTROYED])
            return;
          const ret2 = super.emit("close");
          this.removeAllListeners("close");
          return ret2;
        } else if (ev === "error") {
          this[EMITTED_ERROR] = data;
          super.emit(ERROR, data);
          const ret2 = !this[SIGNAL] || this.listeners("error").length ? super.emit("error", data) : false;
          this[MAYBE_EMIT_END]();
          return ret2;
        } else if (ev === "resume") {
          const ret2 = super.emit("resume");
          this[MAYBE_EMIT_END]();
          return ret2;
        } else if (ev === "finish" || ev === "prefinish") {
          const ret2 = super.emit(ev);
          this.removeAllListeners(ev);
          return ret2;
        }
        const ret = super.emit(ev, data, ...extra);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [EMITDATA](data) {
        for (const p of this[PIPES]) {
          if (p.dest.write(data) === false)
            this.pause();
        }
        const ret = super.emit("data", data);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [EMITEND]() {
        if (this[EMITTED_END])
          return;
        this[EMITTED_END] = true;
        this.readable = false;
        if (this[ASYNC])
          defer(() => this[EMITEND2]());
        else
          this[EMITEND2]();
      }
      [EMITEND2]() {
        if (this[DECODER]) {
          const data = this[DECODER].end();
          if (data) {
            for (const p of this[PIPES]) {
              p.dest.write(data);
            }
            super.emit("data", data);
          }
        }
        for (const p of this[PIPES]) {
          p.end();
        }
        const ret = super.emit("end");
        this.removeAllListeners("end");
        return ret;
      }
      // const all = await stream.collect()
      collect() {
        const buf = [];
        if (!this[OBJECTMODE])
          buf.dataLength = 0;
        const p = this.promise();
        this.on("data", (c2) => {
          buf.push(c2);
          if (!this[OBJECTMODE])
            buf.dataLength += c2.length;
        });
        return p.then(() => buf);
      }
      // const data = await stream.concat()
      concat() {
        return this[OBJECTMODE] ? Promise.reject(new Error("cannot concat in objectMode")) : this.collect().then(
          (buf) => this[OBJECTMODE] ? Promise.reject(new Error("cannot concat in objectMode")) : this[ENCODING] ? buf.join("") : Buffer.concat(buf, buf.dataLength)
        );
      }
      // stream.promise().then(() => done, er => emitted error)
      promise() {
        return new Promise((resolve7, reject) => {
          this.on(DESTROYED, () => reject(new Error("stream destroyed")));
          this.on("error", (er) => reject(er));
          this.on("end", () => resolve7());
        });
      }
      // for await (let chunk of stream)
      [ASYNCITERATOR]() {
        let stopped = false;
        const stop = () => {
          this.pause();
          stopped = true;
          return Promise.resolve({ done: true });
        };
        const next = () => {
          if (stopped)
            return stop();
          const res = this.read();
          if (res !== null)
            return Promise.resolve({ done: false, value: res });
          if (this[EOF])
            return stop();
          let resolve7 = null;
          let reject = null;
          const onerr = (er) => {
            this.removeListener("data", ondata);
            this.removeListener("end", onend);
            this.removeListener(DESTROYED, ondestroy);
            stop();
            reject(er);
          };
          const ondata = (value) => {
            this.removeListener("error", onerr);
            this.removeListener("end", onend);
            this.removeListener(DESTROYED, ondestroy);
            this.pause();
            resolve7({ value, done: !!this[EOF] });
          };
          const onend = () => {
            this.removeListener("error", onerr);
            this.removeListener("data", ondata);
            this.removeListener(DESTROYED, ondestroy);
            stop();
            resolve7({ done: true });
          };
          const ondestroy = () => onerr(new Error("stream destroyed"));
          return new Promise((res2, rej) => {
            reject = rej;
            resolve7 = res2;
            this.once(DESTROYED, ondestroy);
            this.once("error", onerr);
            this.once("end", onend);
            this.once("data", ondata);
          });
        };
        return {
          next,
          throw: stop,
          return: stop,
          [ASYNCITERATOR]() {
            return this;
          }
        };
      }
      // for (let chunk of stream)
      [ITERATOR]() {
        let stopped = false;
        const stop = () => {
          this.pause();
          this.removeListener(ERROR, stop);
          this.removeListener(DESTROYED, stop);
          this.removeListener("end", stop);
          stopped = true;
          return { done: true };
        };
        const next = () => {
          if (stopped)
            return stop();
          const value = this.read();
          return value === null ? stop() : { value };
        };
        this.once("end", stop);
        this.once(ERROR, stop);
        this.once(DESTROYED, stop);
        return {
          next,
          throw: stop,
          return: stop,
          [ITERATOR]() {
            return this;
          }
        };
      }
      destroy(er) {
        if (this[DESTROYED]) {
          if (er)
            this.emit("error", er);
          else
            this.emit(DESTROYED);
          return this;
        }
        this[DESTROYED] = true;
        this[BUFFER].length = 0;
        this[BUFFERLENGTH] = 0;
        if (typeof this.close === "function" && !this[CLOSED])
          this.close();
        if (er)
          this.emit("error", er);
        else
          this.emit(DESTROYED);
        return this;
      }
      static isStream(s) {
        return !!s && (s instanceof Minipass || s instanceof Stream || s instanceof EE && // readable
        (typeof s.pipe === "function" || // writable
        typeof s.write === "function" && typeof s.end === "function"));
      }
    };
    module2.exports = Minipass;
  }
});

// node_modules/minizlib/constants.js
var require_constants = __commonJS({
  "node_modules/minizlib/constants.js"(exports, module2) {
    var realZlibConstants = require("zlib").constants || /* istanbul ignore next */
    { ZLIB_VERNUM: 4736 };
    module2.exports = Object.freeze(Object.assign(/* @__PURE__ */ Object.create(null), {
      Z_NO_FLUSH: 0,
      Z_PARTIAL_FLUSH: 1,
      Z_SYNC_FLUSH: 2,
      Z_FULL_FLUSH: 3,
      Z_FINISH: 4,
      Z_BLOCK: 5,
      Z_OK: 0,
      Z_STREAM_END: 1,
      Z_NEED_DICT: 2,
      Z_ERRNO: -1,
      Z_STREAM_ERROR: -2,
      Z_DATA_ERROR: -3,
      Z_MEM_ERROR: -4,
      Z_BUF_ERROR: -5,
      Z_VERSION_ERROR: -6,
      Z_NO_COMPRESSION: 0,
      Z_BEST_SPEED: 1,
      Z_BEST_COMPRESSION: 9,
      Z_DEFAULT_COMPRESSION: -1,
      Z_FILTERED: 1,
      Z_HUFFMAN_ONLY: 2,
      Z_RLE: 3,
      Z_FIXED: 4,
      Z_DEFAULT_STRATEGY: 0,
      DEFLATE: 1,
      INFLATE: 2,
      GZIP: 3,
      GUNZIP: 4,
      DEFLATERAW: 5,
      INFLATERAW: 6,
      UNZIP: 7,
      BROTLI_DECODE: 8,
      BROTLI_ENCODE: 9,
      Z_MIN_WINDOWBITS: 8,
      Z_MAX_WINDOWBITS: 15,
      Z_DEFAULT_WINDOWBITS: 15,
      Z_MIN_CHUNK: 64,
      Z_MAX_CHUNK: Infinity,
      Z_DEFAULT_CHUNK: 16384,
      Z_MIN_MEMLEVEL: 1,
      Z_MAX_MEMLEVEL: 9,
      Z_DEFAULT_MEMLEVEL: 8,
      Z_MIN_LEVEL: -1,
      Z_MAX_LEVEL: 9,
      Z_DEFAULT_LEVEL: -1,
      BROTLI_OPERATION_PROCESS: 0,
      BROTLI_OPERATION_FLUSH: 1,
      BROTLI_OPERATION_FINISH: 2,
      BROTLI_OPERATION_EMIT_METADATA: 3,
      BROTLI_MODE_GENERIC: 0,
      BROTLI_MODE_TEXT: 1,
      BROTLI_MODE_FONT: 2,
      BROTLI_DEFAULT_MODE: 0,
      BROTLI_MIN_QUALITY: 0,
      BROTLI_MAX_QUALITY: 11,
      BROTLI_DEFAULT_QUALITY: 11,
      BROTLI_MIN_WINDOW_BITS: 10,
      BROTLI_MAX_WINDOW_BITS: 24,
      BROTLI_LARGE_MAX_WINDOW_BITS: 30,
      BROTLI_DEFAULT_WINDOW: 22,
      BROTLI_MIN_INPUT_BLOCK_BITS: 16,
      BROTLI_MAX_INPUT_BLOCK_BITS: 24,
      BROTLI_PARAM_MODE: 0,
      BROTLI_PARAM_QUALITY: 1,
      BROTLI_PARAM_LGWIN: 2,
      BROTLI_PARAM_LGBLOCK: 3,
      BROTLI_PARAM_DISABLE_LITERAL_CONTEXT_MODELING: 4,
      BROTLI_PARAM_SIZE_HINT: 5,
      BROTLI_PARAM_LARGE_WINDOW: 6,
      BROTLI_PARAM_NPOSTFIX: 7,
      BROTLI_PARAM_NDIRECT: 8,
      BROTLI_DECODER_RESULT_ERROR: 0,
      BROTLI_DECODER_RESULT_SUCCESS: 1,
      BROTLI_DECODER_RESULT_NEEDS_MORE_INPUT: 2,
      BROTLI_DECODER_RESULT_NEEDS_MORE_OUTPUT: 3,
      BROTLI_DECODER_PARAM_DISABLE_RING_BUFFER_REALLOCATION: 0,
      BROTLI_DECODER_PARAM_LARGE_WINDOW: 1,
      BROTLI_DECODER_NO_ERROR: 0,
      BROTLI_DECODER_SUCCESS: 1,
      BROTLI_DECODER_NEEDS_MORE_INPUT: 2,
      BROTLI_DECODER_NEEDS_MORE_OUTPUT: 3,
      BROTLI_DECODER_ERROR_FORMAT_EXUBERANT_NIBBLE: -1,
      BROTLI_DECODER_ERROR_FORMAT_RESERVED: -2,
      BROTLI_DECODER_ERROR_FORMAT_EXUBERANT_META_NIBBLE: -3,
      BROTLI_DECODER_ERROR_FORMAT_SIMPLE_HUFFMAN_ALPHABET: -4,
      BROTLI_DECODER_ERROR_FORMAT_SIMPLE_HUFFMAN_SAME: -5,
      BROTLI_DECODER_ERROR_FORMAT_CL_SPACE: -6,
      BROTLI_DECODER_ERROR_FORMAT_HUFFMAN_SPACE: -7,
      BROTLI_DECODER_ERROR_FORMAT_CONTEXT_MAP_REPEAT: -8,
      BROTLI_DECODER_ERROR_FORMAT_BLOCK_LENGTH_1: -9,
      BROTLI_DECODER_ERROR_FORMAT_BLOCK_LENGTH_2: -10,
      BROTLI_DECODER_ERROR_FORMAT_TRANSFORM: -11,
      BROTLI_DECODER_ERROR_FORMAT_DICTIONARY: -12,
      BROTLI_DECODER_ERROR_FORMAT_WINDOW_BITS: -13,
      BROTLI_DECODER_ERROR_FORMAT_PADDING_1: -14,
      BROTLI_DECODER_ERROR_FORMAT_PADDING_2: -15,
      BROTLI_DECODER_ERROR_FORMAT_DISTANCE: -16,
      BROTLI_DECODER_ERROR_DICTIONARY_NOT_SET: -19,
      BROTLI_DECODER_ERROR_INVALID_ARGUMENTS: -20,
      BROTLI_DECODER_ERROR_ALLOC_CONTEXT_MODES: -21,
      BROTLI_DECODER_ERROR_ALLOC_TREE_GROUPS: -22,
      BROTLI_DECODER_ERROR_ALLOC_CONTEXT_MAP: -25,
      BROTLI_DECODER_ERROR_ALLOC_RING_BUFFER_1: -26,
      BROTLI_DECODER_ERROR_ALLOC_RING_BUFFER_2: -27,
      BROTLI_DECODER_ERROR_ALLOC_BLOCK_TYPE_TREES: -30,
      BROTLI_DECODER_ERROR_UNREACHABLE: -31
    }, realZlibConstants));
  }
});

// node_modules/minizlib/node_modules/minipass/index.js
var require_minipass2 = __commonJS({
  "node_modules/minizlib/node_modules/minipass/index.js"(exports, module2) {
    "use strict";
    var proc = typeof process === "object" && process ? process : {
      stdout: null,
      stderr: null
    };
    var EE = require("events");
    var Stream = require("stream");
    var SD = require("string_decoder").StringDecoder;
    var EOF = Symbol("EOF");
    var MAYBE_EMIT_END = Symbol("maybeEmitEnd");
    var EMITTED_END = Symbol("emittedEnd");
    var EMITTING_END = Symbol("emittingEnd");
    var EMITTED_ERROR = Symbol("emittedError");
    var CLOSED = Symbol("closed");
    var READ = Symbol("read");
    var FLUSH = Symbol("flush");
    var FLUSHCHUNK = Symbol("flushChunk");
    var ENCODING = Symbol("encoding");
    var DECODER = Symbol("decoder");
    var FLOWING = Symbol("flowing");
    var PAUSED = Symbol("paused");
    var RESUME = Symbol("resume");
    var BUFFERLENGTH = Symbol("bufferLength");
    var BUFFERPUSH = Symbol("bufferPush");
    var BUFFERSHIFT = Symbol("bufferShift");
    var OBJECTMODE = Symbol("objectMode");
    var DESTROYED = Symbol("destroyed");
    var EMITDATA = Symbol("emitData");
    var EMITEND = Symbol("emitEnd");
    var EMITEND2 = Symbol("emitEnd2");
    var ASYNC = Symbol("async");
    var defer = (fn) => Promise.resolve().then(fn);
    var doIter = global._MP_NO_ITERATOR_SYMBOLS_ !== "1";
    var ASYNCITERATOR = doIter && Symbol.asyncIterator || Symbol("asyncIterator not implemented");
    var ITERATOR = doIter && Symbol.iterator || Symbol("iterator not implemented");
    var isEndish = (ev) => ev === "end" || ev === "finish" || ev === "prefinish";
    var isArrayBuffer = (b) => b instanceof ArrayBuffer || typeof b === "object" && b.constructor && b.constructor.name === "ArrayBuffer" && b.byteLength >= 0;
    var isArrayBufferView = (b) => !Buffer.isBuffer(b) && ArrayBuffer.isView(b);
    var Pipe = class {
      constructor(src, dest, opts) {
        this.src = src;
        this.dest = dest;
        this.opts = opts;
        this.ondrain = () => src[RESUME]();
        dest.on("drain", this.ondrain);
      }
      unpipe() {
        this.dest.removeListener("drain", this.ondrain);
      }
      // istanbul ignore next - only here for the prototype
      proxyErrors() {
      }
      end() {
        this.unpipe();
        if (this.opts.end)
          this.dest.end();
      }
    };
    var PipeProxyErrors = class extends Pipe {
      unpipe() {
        this.src.removeListener("error", this.proxyErrors);
        super.unpipe();
      }
      constructor(src, dest, opts) {
        super(src, dest, opts);
        this.proxyErrors = (er) => dest.emit("error", er);
        src.on("error", this.proxyErrors);
      }
    };
    module2.exports = class Minipass extends Stream {
      constructor(options) {
        super();
        this[FLOWING] = false;
        this[PAUSED] = false;
        this.pipes = [];
        this.buffer = [];
        this[OBJECTMODE] = options && options.objectMode || false;
        if (this[OBJECTMODE])
          this[ENCODING] = null;
        else
          this[ENCODING] = options && options.encoding || null;
        if (this[ENCODING] === "buffer")
          this[ENCODING] = null;
        this[ASYNC] = options && !!options.async || false;
        this[DECODER] = this[ENCODING] ? new SD(this[ENCODING]) : null;
        this[EOF] = false;
        this[EMITTED_END] = false;
        this[EMITTING_END] = false;
        this[CLOSED] = false;
        this[EMITTED_ERROR] = null;
        this.writable = true;
        this.readable = true;
        this[BUFFERLENGTH] = 0;
        this[DESTROYED] = false;
      }
      get bufferLength() {
        return this[BUFFERLENGTH];
      }
      get encoding() {
        return this[ENCODING];
      }
      set encoding(enc) {
        if (this[OBJECTMODE])
          throw new Error("cannot set encoding in objectMode");
        if (this[ENCODING] && enc !== this[ENCODING] && (this[DECODER] && this[DECODER].lastNeed || this[BUFFERLENGTH]))
          throw new Error("cannot change encoding");
        if (this[ENCODING] !== enc) {
          this[DECODER] = enc ? new SD(enc) : null;
          if (this.buffer.length)
            this.buffer = this.buffer.map((chunk) => this[DECODER].write(chunk));
        }
        this[ENCODING] = enc;
      }
      setEncoding(enc) {
        this.encoding = enc;
      }
      get objectMode() {
        return this[OBJECTMODE];
      }
      set objectMode(om) {
        this[OBJECTMODE] = this[OBJECTMODE] || !!om;
      }
      get ["async"]() {
        return this[ASYNC];
      }
      set ["async"](a) {
        this[ASYNC] = this[ASYNC] || !!a;
      }
      write(chunk, encoding, cb) {
        if (this[EOF])
          throw new Error("write after end");
        if (this[DESTROYED]) {
          this.emit("error", Object.assign(
            new Error("Cannot call write after a stream was destroyed"),
            { code: "ERR_STREAM_DESTROYED" }
          ));
          return true;
        }
        if (typeof encoding === "function")
          cb = encoding, encoding = "utf8";
        if (!encoding)
          encoding = "utf8";
        const fn = this[ASYNC] ? defer : (f) => f();
        if (!this[OBJECTMODE] && !Buffer.isBuffer(chunk)) {
          if (isArrayBufferView(chunk))
            chunk = Buffer.from(chunk.buffer, chunk.byteOffset, chunk.byteLength);
          else if (isArrayBuffer(chunk))
            chunk = Buffer.from(chunk);
          else if (typeof chunk !== "string")
            this.objectMode = true;
        }
        if (this[OBJECTMODE]) {
          if (this.flowing && this[BUFFERLENGTH] !== 0)
            this[FLUSH](true);
          if (this.flowing)
            this.emit("data", chunk);
          else
            this[BUFFERPUSH](chunk);
          if (this[BUFFERLENGTH] !== 0)
            this.emit("readable");
          if (cb)
            fn(cb);
          return this.flowing;
        }
        if (!chunk.length) {
          if (this[BUFFERLENGTH] !== 0)
            this.emit("readable");
          if (cb)
            fn(cb);
          return this.flowing;
        }
        if (typeof chunk === "string" && // unless it is a string already ready for us to use
        !(encoding === this[ENCODING] && !this[DECODER].lastNeed)) {
          chunk = Buffer.from(chunk, encoding);
        }
        if (Buffer.isBuffer(chunk) && this[ENCODING])
          chunk = this[DECODER].write(chunk);
        if (this.flowing && this[BUFFERLENGTH] !== 0)
          this[FLUSH](true);
        if (this.flowing)
          this.emit("data", chunk);
        else
          this[BUFFERPUSH](chunk);
        if (this[BUFFERLENGTH] !== 0)
          this.emit("readable");
        if (cb)
          fn(cb);
        return this.flowing;
      }
      read(n) {
        if (this[DESTROYED])
          return null;
        if (this[BUFFERLENGTH] === 0 || n === 0 || n > this[BUFFERLENGTH]) {
          this[MAYBE_EMIT_END]();
          return null;
        }
        if (this[OBJECTMODE])
          n = null;
        if (this.buffer.length > 1 && !this[OBJECTMODE]) {
          if (this.encoding)
            this.buffer = [this.buffer.join("")];
          else
            this.buffer = [Buffer.concat(this.buffer, this[BUFFERLENGTH])];
        }
        const ret = this[READ](n || null, this.buffer[0]);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [READ](n, chunk) {
        if (n === chunk.length || n === null)
          this[BUFFERSHIFT]();
        else {
          this.buffer[0] = chunk.slice(n);
          chunk = chunk.slice(0, n);
          this[BUFFERLENGTH] -= n;
        }
        this.emit("data", chunk);
        if (!this.buffer.length && !this[EOF])
          this.emit("drain");
        return chunk;
      }
      end(chunk, encoding, cb) {
        if (typeof chunk === "function")
          cb = chunk, chunk = null;
        if (typeof encoding === "function")
          cb = encoding, encoding = "utf8";
        if (chunk)
          this.write(chunk, encoding);
        if (cb)
          this.once("end", cb);
        this[EOF] = true;
        this.writable = false;
        if (this.flowing || !this[PAUSED])
          this[MAYBE_EMIT_END]();
        return this;
      }
      // don't let the internal resume be overwritten
      [RESUME]() {
        if (this[DESTROYED])
          return;
        this[PAUSED] = false;
        this[FLOWING] = true;
        this.emit("resume");
        if (this.buffer.length)
          this[FLUSH]();
        else if (this[EOF])
          this[MAYBE_EMIT_END]();
        else
          this.emit("drain");
      }
      resume() {
        return this[RESUME]();
      }
      pause() {
        this[FLOWING] = false;
        this[PAUSED] = true;
      }
      get destroyed() {
        return this[DESTROYED];
      }
      get flowing() {
        return this[FLOWING];
      }
      get paused() {
        return this[PAUSED];
      }
      [BUFFERPUSH](chunk) {
        if (this[OBJECTMODE])
          this[BUFFERLENGTH] += 1;
        else
          this[BUFFERLENGTH] += chunk.length;
        this.buffer.push(chunk);
      }
      [BUFFERSHIFT]() {
        if (this.buffer.length) {
          if (this[OBJECTMODE])
            this[BUFFERLENGTH] -= 1;
          else
            this[BUFFERLENGTH] -= this.buffer[0].length;
        }
        return this.buffer.shift();
      }
      [FLUSH](noDrain) {
        do {
        } while (this[FLUSHCHUNK](this[BUFFERSHIFT]()));
        if (!noDrain && !this.buffer.length && !this[EOF])
          this.emit("drain");
      }
      [FLUSHCHUNK](chunk) {
        return chunk ? (this.emit("data", chunk), this.flowing) : false;
      }
      pipe(dest, opts) {
        if (this[DESTROYED])
          return;
        const ended = this[EMITTED_END];
        opts = opts || {};
        if (dest === proc.stdout || dest === proc.stderr)
          opts.end = false;
        else
          opts.end = opts.end !== false;
        opts.proxyErrors = !!opts.proxyErrors;
        if (ended) {
          if (opts.end)
            dest.end();
        } else {
          this.pipes.push(!opts.proxyErrors ? new Pipe(this, dest, opts) : new PipeProxyErrors(this, dest, opts));
          if (this[ASYNC])
            defer(() => this[RESUME]());
          else
            this[RESUME]();
        }
        return dest;
      }
      unpipe(dest) {
        const p = this.pipes.find((p2) => p2.dest === dest);
        if (p) {
          this.pipes.splice(this.pipes.indexOf(p), 1);
          p.unpipe();
        }
      }
      addListener(ev, fn) {
        return this.on(ev, fn);
      }
      on(ev, fn) {
        const ret = super.on(ev, fn);
        if (ev === "data" && !this.pipes.length && !this.flowing)
          this[RESUME]();
        else if (ev === "readable" && this[BUFFERLENGTH] !== 0)
          super.emit("readable");
        else if (isEndish(ev) && this[EMITTED_END]) {
          super.emit(ev);
          this.removeAllListeners(ev);
        } else if (ev === "error" && this[EMITTED_ERROR]) {
          if (this[ASYNC])
            defer(() => fn.call(this, this[EMITTED_ERROR]));
          else
            fn.call(this, this[EMITTED_ERROR]);
        }
        return ret;
      }
      get emittedEnd() {
        return this[EMITTED_END];
      }
      [MAYBE_EMIT_END]() {
        if (!this[EMITTING_END] && !this[EMITTED_END] && !this[DESTROYED] && this.buffer.length === 0 && this[EOF]) {
          this[EMITTING_END] = true;
          this.emit("end");
          this.emit("prefinish");
          this.emit("finish");
          if (this[CLOSED])
            this.emit("close");
          this[EMITTING_END] = false;
        }
      }
      emit(ev, data, ...extra) {
        if (ev !== "error" && ev !== "close" && ev !== DESTROYED && this[DESTROYED])
          return;
        else if (ev === "data") {
          return !data ? false : this[ASYNC] ? defer(() => this[EMITDATA](data)) : this[EMITDATA](data);
        } else if (ev === "end") {
          return this[EMITEND]();
        } else if (ev === "close") {
          this[CLOSED] = true;
          if (!this[EMITTED_END] && !this[DESTROYED])
            return;
          const ret2 = super.emit("close");
          this.removeAllListeners("close");
          return ret2;
        } else if (ev === "error") {
          this[EMITTED_ERROR] = data;
          const ret2 = super.emit("error", data);
          this[MAYBE_EMIT_END]();
          return ret2;
        } else if (ev === "resume") {
          const ret2 = super.emit("resume");
          this[MAYBE_EMIT_END]();
          return ret2;
        } else if (ev === "finish" || ev === "prefinish") {
          const ret2 = super.emit(ev);
          this.removeAllListeners(ev);
          return ret2;
        }
        const ret = super.emit(ev, data, ...extra);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [EMITDATA](data) {
        for (const p of this.pipes) {
          if (p.dest.write(data) === false)
            this.pause();
        }
        const ret = super.emit("data", data);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [EMITEND]() {
        if (this[EMITTED_END])
          return;
        this[EMITTED_END] = true;
        this.readable = false;
        if (this[ASYNC])
          defer(() => this[EMITEND2]());
        else
          this[EMITEND2]();
      }
      [EMITEND2]() {
        if (this[DECODER]) {
          const data = this[DECODER].end();
          if (data) {
            for (const p of this.pipes) {
              p.dest.write(data);
            }
            super.emit("data", data);
          }
        }
        for (const p of this.pipes) {
          p.end();
        }
        const ret = super.emit("end");
        this.removeAllListeners("end");
        return ret;
      }
      // const all = await stream.collect()
      collect() {
        const buf = [];
        if (!this[OBJECTMODE])
          buf.dataLength = 0;
        const p = this.promise();
        this.on("data", (c2) => {
          buf.push(c2);
          if (!this[OBJECTMODE])
            buf.dataLength += c2.length;
        });
        return p.then(() => buf);
      }
      // const data = await stream.concat()
      concat() {
        return this[OBJECTMODE] ? Promise.reject(new Error("cannot concat in objectMode")) : this.collect().then((buf) => this[OBJECTMODE] ? Promise.reject(new Error("cannot concat in objectMode")) : this[ENCODING] ? buf.join("") : Buffer.concat(buf, buf.dataLength));
      }
      // stream.promise().then(() => done, er => emitted error)
      promise() {
        return new Promise((resolve7, reject) => {
          this.on(DESTROYED, () => reject(new Error("stream destroyed")));
          this.on("error", (er) => reject(er));
          this.on("end", () => resolve7());
        });
      }
      // for await (let chunk of stream)
      [ASYNCITERATOR]() {
        const next = () => {
          const res = this.read();
          if (res !== null)
            return Promise.resolve({ done: false, value: res });
          if (this[EOF])
            return Promise.resolve({ done: true });
          let resolve7 = null;
          let reject = null;
          const onerr = (er) => {
            this.removeListener("data", ondata);
            this.removeListener("end", onend);
            reject(er);
          };
          const ondata = (value) => {
            this.removeListener("error", onerr);
            this.removeListener("end", onend);
            this.pause();
            resolve7({ value, done: !!this[EOF] });
          };
          const onend = () => {
            this.removeListener("error", onerr);
            this.removeListener("data", ondata);
            resolve7({ done: true });
          };
          const ondestroy = () => onerr(new Error("stream destroyed"));
          return new Promise((res2, rej) => {
            reject = rej;
            resolve7 = res2;
            this.once(DESTROYED, ondestroy);
            this.once("error", onerr);
            this.once("end", onend);
            this.once("data", ondata);
          });
        };
        return { next };
      }
      // for (let chunk of stream)
      [ITERATOR]() {
        const next = () => {
          const value = this.read();
          const done = value === null;
          return { value, done };
        };
        return { next };
      }
      destroy(er) {
        if (this[DESTROYED]) {
          if (er)
            this.emit("error", er);
          else
            this.emit(DESTROYED);
          return this;
        }
        this[DESTROYED] = true;
        this.buffer.length = 0;
        this[BUFFERLENGTH] = 0;
        if (typeof this.close === "function" && !this[CLOSED])
          this.close();
        if (er)
          this.emit("error", er);
        else
          this.emit(DESTROYED);
        return this;
      }
      static isStream(s) {
        return !!s && (s instanceof Minipass || s instanceof Stream || s instanceof EE && (typeof s.pipe === "function" || // readable
        typeof s.write === "function" && typeof s.end === "function"));
      }
    };
  }
});

// node_modules/minizlib/index.js
var require_minizlib = __commonJS({
  "node_modules/minizlib/index.js"(exports) {
    "use strict";
    var assert = require("assert");
    var Buffer2 = require("buffer").Buffer;
    var realZlib = require("zlib");
    var constants = exports.constants = require_constants();
    var Minipass = require_minipass2();
    var OriginalBufferConcat = Buffer2.concat;
    var _superWrite = Symbol("_superWrite");
    var ZlibError = class extends Error {
      constructor(err) {
        super("zlib: " + err.message);
        this.code = err.code;
        this.errno = err.errno;
        if (!this.code)
          this.code = "ZLIB_ERROR";
        this.message = "zlib: " + err.message;
        Error.captureStackTrace(this, this.constructor);
      }
      get name() {
        return "ZlibError";
      }
    };
    var _opts = Symbol("opts");
    var _flushFlag = Symbol("flushFlag");
    var _finishFlushFlag = Symbol("finishFlushFlag");
    var _fullFlushFlag = Symbol("fullFlushFlag");
    var _handle = Symbol("handle");
    var _onError = Symbol("onError");
    var _sawError = Symbol("sawError");
    var _level = Symbol("level");
    var _strategy = Symbol("strategy");
    var _ended = Symbol("ended");
    var _defaultFullFlush = Symbol("_defaultFullFlush");
    var ZlibBase = class extends Minipass {
      constructor(opts, mode) {
        if (!opts || typeof opts !== "object")
          throw new TypeError("invalid options for ZlibBase constructor");
        super(opts);
        this[_sawError] = false;
        this[_ended] = false;
        this[_opts] = opts;
        this[_flushFlag] = opts.flush;
        this[_finishFlushFlag] = opts.finishFlush;
        try {
          this[_handle] = new realZlib[mode](opts);
        } catch (er) {
          throw new ZlibError(er);
        }
        this[_onError] = (err) => {
          if (this[_sawError])
            return;
          this[_sawError] = true;
          this.close();
          this.emit("error", err);
        };
        this[_handle].on("error", (er) => this[_onError](new ZlibError(er)));
        this.once("end", () => this.close);
      }
      close() {
        if (this[_handle]) {
          this[_handle].close();
          this[_handle] = null;
          this.emit("close");
        }
      }
      reset() {
        if (!this[_sawError]) {
          assert(this[_handle], "zlib binding closed");
          return this[_handle].reset();
        }
      }
      flush(flushFlag) {
        if (this.ended)
          return;
        if (typeof flushFlag !== "number")
          flushFlag = this[_fullFlushFlag];
        this.write(Object.assign(Buffer2.alloc(0), { [_flushFlag]: flushFlag }));
      }
      end(chunk, encoding, cb) {
        if (chunk)
          this.write(chunk, encoding);
        this.flush(this[_finishFlushFlag]);
        this[_ended] = true;
        return super.end(null, null, cb);
      }
      get ended() {
        return this[_ended];
      }
      write(chunk, encoding, cb) {
        if (typeof encoding === "function")
          cb = encoding, encoding = "utf8";
        if (typeof chunk === "string")
          chunk = Buffer2.from(chunk, encoding);
        if (this[_sawError])
          return;
        assert(this[_handle], "zlib binding closed");
        const nativeHandle = this[_handle]._handle;
        const originalNativeClose = nativeHandle.close;
        nativeHandle.close = () => {
        };
        const originalClose = this[_handle].close;
        this[_handle].close = () => {
        };
        Buffer2.concat = (args) => args;
        let result;
        try {
          const flushFlag = typeof chunk[_flushFlag] === "number" ? chunk[_flushFlag] : this[_flushFlag];
          result = this[_handle]._processChunk(chunk, flushFlag);
          Buffer2.concat = OriginalBufferConcat;
        } catch (err) {
          Buffer2.concat = OriginalBufferConcat;
          this[_onError](new ZlibError(err));
        } finally {
          if (this[_handle]) {
            this[_handle]._handle = nativeHandle;
            nativeHandle.close = originalNativeClose;
            this[_handle].close = originalClose;
            this[_handle].removeAllListeners("error");
          }
        }
        if (this[_handle])
          this[_handle].on("error", (er) => this[_onError](new ZlibError(er)));
        let writeReturn;
        if (result) {
          if (Array.isArray(result) && result.length > 0) {
            writeReturn = this[_superWrite](Buffer2.from(result[0]));
            for (let i = 1; i < result.length; i++) {
              writeReturn = this[_superWrite](result[i]);
            }
          } else {
            writeReturn = this[_superWrite](Buffer2.from(result));
          }
        }
        if (cb)
          cb();
        return writeReturn;
      }
      [_superWrite](data) {
        return super.write(data);
      }
    };
    var Zlib = class extends ZlibBase {
      constructor(opts, mode) {
        opts = opts || {};
        opts.flush = opts.flush || constants.Z_NO_FLUSH;
        opts.finishFlush = opts.finishFlush || constants.Z_FINISH;
        super(opts, mode);
        this[_fullFlushFlag] = constants.Z_FULL_FLUSH;
        this[_level] = opts.level;
        this[_strategy] = opts.strategy;
      }
      params(level, strategy) {
        if (this[_sawError])
          return;
        if (!this[_handle])
          throw new Error("cannot switch params when binding is closed");
        if (!this[_handle].params)
          throw new Error("not supported in this implementation");
        if (this[_level] !== level || this[_strategy] !== strategy) {
          this.flush(constants.Z_SYNC_FLUSH);
          assert(this[_handle], "zlib binding closed");
          const origFlush = this[_handle].flush;
          this[_handle].flush = (flushFlag, cb) => {
            this.flush(flushFlag);
            cb();
          };
          try {
            this[_handle].params(level, strategy);
          } finally {
            this[_handle].flush = origFlush;
          }
          if (this[_handle]) {
            this[_level] = level;
            this[_strategy] = strategy;
          }
        }
      }
    };
    var Deflate = class extends Zlib {
      constructor(opts) {
        super(opts, "Deflate");
      }
    };
    var Inflate = class extends Zlib {
      constructor(opts) {
        super(opts, "Inflate");
      }
    };
    var _portable = Symbol("_portable");
    var Gzip = class extends Zlib {
      constructor(opts) {
        super(opts, "Gzip");
        this[_portable] = opts && !!opts.portable;
      }
      [_superWrite](data) {
        if (!this[_portable])
          return super[_superWrite](data);
        this[_portable] = false;
        data[9] = 255;
        return super[_superWrite](data);
      }
    };
    var Gunzip = class extends Zlib {
      constructor(opts) {
        super(opts, "Gunzip");
      }
    };
    var DeflateRaw = class extends Zlib {
      constructor(opts) {
        super(opts, "DeflateRaw");
      }
    };
    var InflateRaw = class extends Zlib {
      constructor(opts) {
        super(opts, "InflateRaw");
      }
    };
    var Unzip = class extends Zlib {
      constructor(opts) {
        super(opts, "Unzip");
      }
    };
    var Brotli = class extends ZlibBase {
      constructor(opts, mode) {
        opts = opts || {};
        opts.flush = opts.flush || constants.BROTLI_OPERATION_PROCESS;
        opts.finishFlush = opts.finishFlush || constants.BROTLI_OPERATION_FINISH;
        super(opts, mode);
        this[_fullFlushFlag] = constants.BROTLI_OPERATION_FLUSH;
      }
    };
    var BrotliCompress = class extends Brotli {
      constructor(opts) {
        super(opts, "BrotliCompress");
      }
    };
    var BrotliDecompress = class extends Brotli {
      constructor(opts) {
        super(opts, "BrotliDecompress");
      }
    };
    exports.Deflate = Deflate;
    exports.Inflate = Inflate;
    exports.Gzip = Gzip;
    exports.Gunzip = Gunzip;
    exports.DeflateRaw = DeflateRaw;
    exports.InflateRaw = InflateRaw;
    exports.Unzip = Unzip;
    if (typeof realZlib.BrotliCompress === "function") {
      exports.BrotliCompress = BrotliCompress;
      exports.BrotliDecompress = BrotliDecompress;
    } else {
      exports.BrotliCompress = exports.BrotliDecompress = class {
        constructor() {
          throw new Error("Brotli is not supported in this version of Node.js");
        }
      };
    }
  }
});

// node_modules/tar/lib/normalize-windows-path.js
var require_normalize_windows_path = __commonJS({
  "node_modules/tar/lib/normalize-windows-path.js"(exports, module2) {
    var platform2 = process.env.TESTING_TAR_FAKE_PLATFORM || process.platform;
    module2.exports = platform2 !== "win32" ? (p) => p : (p) => p && p.replace(/\\/g, "/");
  }
});

// node_modules/tar/lib/read-entry.js
var require_read_entry = __commonJS({
  "node_modules/tar/lib/read-entry.js"(exports, module2) {
    "use strict";
    var MiniPass = require_minipass();
    var normPath = require_normalize_windows_path();
    var SLURP = Symbol("slurp");
    module2.exports = class ReadEntry extends MiniPass {
      constructor(header, ex, gex) {
        super();
        this.pause();
        this.extended = ex;
        this.globalExtended = gex;
        this.header = header;
        this.startBlockSize = 512 * Math.ceil(header.size / 512);
        this.blockRemain = this.startBlockSize;
        this.remain = header.size;
        this.type = header.type;
        this.meta = false;
        this.ignore = false;
        switch (this.type) {
          case "File":
          case "OldFile":
          case "Link":
          case "SymbolicLink":
          case "CharacterDevice":
          case "BlockDevice":
          case "Directory":
          case "FIFO":
          case "ContiguousFile":
          case "GNUDumpDir":
            break;
          case "NextFileHasLongLinkpath":
          case "NextFileHasLongPath":
          case "OldGnuLongPath":
          case "GlobalExtendedHeader":
          case "ExtendedHeader":
          case "OldExtendedHeader":
            this.meta = true;
            break;
          default:
            this.ignore = true;
        }
        this.path = normPath(header.path);
        this.mode = header.mode;
        if (this.mode) {
          this.mode = this.mode & 4095;
        }
        this.uid = header.uid;
        this.gid = header.gid;
        this.uname = header.uname;
        this.gname = header.gname;
        this.size = header.size;
        this.mtime = header.mtime;
        this.atime = header.atime;
        this.ctime = header.ctime;
        this.linkpath = normPath(header.linkpath);
        this.uname = header.uname;
        this.gname = header.gname;
        if (ex) {
          this[SLURP](ex);
        }
        if (gex) {
          this[SLURP](gex, true);
        }
      }
      write(data) {
        const writeLen = data.length;
        if (writeLen > this.blockRemain) {
          throw new Error("writing more to entry than is appropriate");
        }
        const r = this.remain;
        const br = this.blockRemain;
        this.remain = Math.max(0, r - writeLen);
        this.blockRemain = Math.max(0, br - writeLen);
        if (this.ignore) {
          return true;
        }
        if (r >= writeLen) {
          return super.write(data);
        }
        return super.write(data.slice(0, r));
      }
      [SLURP](ex, global2) {
        for (const k in ex) {
          if (ex[k] !== null && ex[k] !== void 0 && !(global2 && k === "path")) {
            this[k] = k === "path" || k === "linkpath" ? normPath(ex[k]) : ex[k];
          }
        }
      }
    };
  }
});

// node_modules/tar/lib/types.js
var require_types = __commonJS({
  "node_modules/tar/lib/types.js"(exports) {
    "use strict";
    exports.name = /* @__PURE__ */ new Map([
      ["0", "File"],
      // same as File
      ["", "OldFile"],
      ["1", "Link"],
      ["2", "SymbolicLink"],
      // Devices and FIFOs aren't fully supported
      // they are parsed, but skipped when unpacking
      ["3", "CharacterDevice"],
      ["4", "BlockDevice"],
      ["5", "Directory"],
      ["6", "FIFO"],
      // same as File
      ["7", "ContiguousFile"],
      // pax headers
      ["g", "GlobalExtendedHeader"],
      ["x", "ExtendedHeader"],
      // vendor-specific stuff
      // skip
      ["A", "SolarisACL"],
      // like 5, but with data, which should be skipped
      ["D", "GNUDumpDir"],
      // metadata only, skip
      ["I", "Inode"],
      // data = link path of next file
      ["K", "NextFileHasLongLinkpath"],
      // data = path of next file
      ["L", "NextFileHasLongPath"],
      // skip
      ["M", "ContinuationFile"],
      // like L
      ["N", "OldGnuLongPath"],
      // skip
      ["S", "SparseFile"],
      // skip
      ["V", "TapeVolumeHeader"],
      // like x
      ["X", "OldExtendedHeader"]
    ]);
    exports.code = new Map(Array.from(exports.name).map((kv) => [kv[1], kv[0]]));
  }
});

// node_modules/tar/lib/large-numbers.js
var require_large_numbers = __commonJS({
  "node_modules/tar/lib/large-numbers.js"(exports, module2) {
    "use strict";
    var encode = (num, buf) => {
      if (!Number.isSafeInteger(num)) {
        throw Error("cannot encode number outside of javascript safe integer range");
      } else if (num < 0) {
        encodeNegative(num, buf);
      } else {
        encodePositive(num, buf);
      }
      return buf;
    };
    var encodePositive = (num, buf) => {
      buf[0] = 128;
      for (var i = buf.length; i > 1; i--) {
        buf[i - 1] = num & 255;
        num = Math.floor(num / 256);
      }
    };
    var encodeNegative = (num, buf) => {
      buf[0] = 255;
      var flipped = false;
      num = num * -1;
      for (var i = buf.length; i > 1; i--) {
        var byte = num & 255;
        num = Math.floor(num / 256);
        if (flipped) {
          buf[i - 1] = onesComp(byte);
        } else if (byte === 0) {
          buf[i - 1] = 0;
        } else {
          flipped = true;
          buf[i - 1] = twosComp(byte);
        }
      }
    };
    var parse5 = (buf) => {
      const pre = buf[0];
      const value = pre === 128 ? pos(buf.slice(1, buf.length)) : pre === 255 ? twos(buf) : null;
      if (value === null) {
        throw Error("invalid base256 encoding");
      }
      if (!Number.isSafeInteger(value)) {
        throw Error("parsed number outside of javascript safe integer range");
      }
      return value;
    };
    var twos = (buf) => {
      var len = buf.length;
      var sum = 0;
      var flipped = false;
      for (var i = len - 1; i > -1; i--) {
        var byte = buf[i];
        var f;
        if (flipped) {
          f = onesComp(byte);
        } else if (byte === 0) {
          f = byte;
        } else {
          flipped = true;
          f = twosComp(byte);
        }
        if (f !== 0) {
          sum -= f * Math.pow(256, len - i - 1);
        }
      }
      return sum;
    };
    var pos = (buf) => {
      var len = buf.length;
      var sum = 0;
      for (var i = len - 1; i > -1; i--) {
        var byte = buf[i];
        if (byte !== 0) {
          sum += byte * Math.pow(256, len - i - 1);
        }
      }
      return sum;
    };
    var onesComp = (byte) => (255 ^ byte) & 255;
    var twosComp = (byte) => (255 ^ byte) + 1 & 255;
    module2.exports = {
      encode,
      parse: parse5
    };
  }
});

// node_modules/tar/lib/header.js
var require_header = __commonJS({
  "node_modules/tar/lib/header.js"(exports, module2) {
    "use strict";
    var types = require_types();
    var pathModule = require("path").posix;
    var large = require_large_numbers();
    var SLURP = Symbol("slurp");
    var TYPE = Symbol("type");
    var Header = class {
      constructor(data, off, ex, gex) {
        this.cksumValid = false;
        this.needPax = false;
        this.nullBlock = false;
        this.block = null;
        this.path = null;
        this.mode = null;
        this.uid = null;
        this.gid = null;
        this.size = null;
        this.mtime = null;
        this.cksum = null;
        this[TYPE] = "0";
        this.linkpath = null;
        this.uname = null;
        this.gname = null;
        this.devmaj = 0;
        this.devmin = 0;
        this.atime = null;
        this.ctime = null;
        if (Buffer.isBuffer(data)) {
          this.decode(data, off || 0, ex, gex);
        } else if (data) {
          this.set(data);
        }
      }
      decode(buf, off, ex, gex) {
        if (!off) {
          off = 0;
        }
        if (!buf || !(buf.length >= off + 512)) {
          throw new Error("need 512 bytes for header");
        }
        this.path = decString(buf, off, 100);
        this.mode = decNumber(buf, off + 100, 8);
        this.uid = decNumber(buf, off + 108, 8);
        this.gid = decNumber(buf, off + 116, 8);
        this.size = decNumber(buf, off + 124, 12);
        this.mtime = decDate(buf, off + 136, 12);
        this.cksum = decNumber(buf, off + 148, 12);
        this[SLURP](ex);
        this[SLURP](gex, true);
        this[TYPE] = decString(buf, off + 156, 1);
        if (this[TYPE] === "") {
          this[TYPE] = "0";
        }
        if (this[TYPE] === "0" && this.path.slice(-1) === "/") {
          this[TYPE] = "5";
        }
        if (this[TYPE] === "5") {
          this.size = 0;
        }
        this.linkpath = decString(buf, off + 157, 100);
        if (buf.slice(off + 257, off + 265).toString() === "ustar\x0000") {
          this.uname = decString(buf, off + 265, 32);
          this.gname = decString(buf, off + 297, 32);
          this.devmaj = decNumber(buf, off + 329, 8);
          this.devmin = decNumber(buf, off + 337, 8);
          if (buf[off + 475] !== 0) {
            const prefix = decString(buf, off + 345, 155);
            this.path = prefix + "/" + this.path;
          } else {
            const prefix = decString(buf, off + 345, 130);
            if (prefix) {
              this.path = prefix + "/" + this.path;
            }
            this.atime = decDate(buf, off + 476, 12);
            this.ctime = decDate(buf, off + 488, 12);
          }
        }
        let sum = 8 * 32;
        for (let i = off; i < off + 148; i++) {
          sum += buf[i];
        }
        for (let i = off + 156; i < off + 512; i++) {
          sum += buf[i];
        }
        this.cksumValid = sum === this.cksum;
        if (this.cksum === null && sum === 8 * 32) {
          this.nullBlock = true;
        }
      }
      [SLURP](ex, global2) {
        for (const k in ex) {
          if (ex[k] !== null && ex[k] !== void 0 && !(global2 && k === "path")) {
            this[k] = ex[k];
          }
        }
      }
      encode(buf, off) {
        if (!buf) {
          buf = this.block = Buffer.alloc(512);
          off = 0;
        }
        if (!off) {
          off = 0;
        }
        if (!(buf.length >= off + 512)) {
          throw new Error("need 512 bytes for header");
        }
        const prefixSize = this.ctime || this.atime ? 130 : 155;
        const split = splitPrefix(this.path || "", prefixSize);
        const path32 = split[0];
        const prefix = split[1];
        this.needPax = split[2];
        this.needPax = encString(buf, off, 100, path32) || this.needPax;
        this.needPax = encNumber(buf, off + 100, 8, this.mode) || this.needPax;
        this.needPax = encNumber(buf, off + 108, 8, this.uid) || this.needPax;
        this.needPax = encNumber(buf, off + 116, 8, this.gid) || this.needPax;
        this.needPax = encNumber(buf, off + 124, 12, this.size) || this.needPax;
        this.needPax = encDate(buf, off + 136, 12, this.mtime) || this.needPax;
        buf[off + 156] = this[TYPE].charCodeAt(0);
        this.needPax = encString(buf, off + 157, 100, this.linkpath) || this.needPax;
        buf.write("ustar\x0000", off + 257, 8);
        this.needPax = encString(buf, off + 265, 32, this.uname) || this.needPax;
        this.needPax = encString(buf, off + 297, 32, this.gname) || this.needPax;
        this.needPax = encNumber(buf, off + 329, 8, this.devmaj) || this.needPax;
        this.needPax = encNumber(buf, off + 337, 8, this.devmin) || this.needPax;
        this.needPax = encString(buf, off + 345, prefixSize, prefix) || this.needPax;
        if (buf[off + 475] !== 0) {
          this.needPax = encString(buf, off + 345, 155, prefix) || this.needPax;
        } else {
          this.needPax = encString(buf, off + 345, 130, prefix) || this.needPax;
          this.needPax = encDate(buf, off + 476, 12, this.atime) || this.needPax;
          this.needPax = encDate(buf, off + 488, 12, this.ctime) || this.needPax;
        }
        let sum = 8 * 32;
        for (let i = off; i < off + 148; i++) {
          sum += buf[i];
        }
        for (let i = off + 156; i < off + 512; i++) {
          sum += buf[i];
        }
        this.cksum = sum;
        encNumber(buf, off + 148, 8, this.cksum);
        this.cksumValid = true;
        return this.needPax;
      }
      set(data) {
        for (const i in data) {
          if (data[i] !== null && data[i] !== void 0) {
            this[i] = data[i];
          }
        }
      }
      get type() {
        return types.name.get(this[TYPE]) || this[TYPE];
      }
      get typeKey() {
        return this[TYPE];
      }
      set type(type2) {
        if (types.code.has(type2)) {
          this[TYPE] = types.code.get(type2);
        } else {
          this[TYPE] = type2;
        }
      }
    };
    var splitPrefix = (p, prefixSize) => {
      const pathSize = 100;
      let pp = p;
      let prefix = "";
      let ret;
      const root = pathModule.parse(p).root || ".";
      if (Buffer.byteLength(pp) < pathSize) {
        ret = [pp, prefix, false];
      } else {
        prefix = pathModule.dirname(pp);
        pp = pathModule.basename(pp);
        do {
          if (Buffer.byteLength(pp) <= pathSize && Buffer.byteLength(prefix) <= prefixSize) {
            ret = [pp, prefix, false];
          } else if (Buffer.byteLength(pp) > pathSize && Buffer.byteLength(prefix) <= prefixSize) {
            ret = [pp.slice(0, pathSize - 1), prefix, true];
          } else {
            pp = pathModule.join(pathModule.basename(prefix), pp);
            prefix = pathModule.dirname(prefix);
          }
        } while (prefix !== root && !ret);
        if (!ret) {
          ret = [p.slice(0, pathSize - 1), "", true];
        }
      }
      return ret;
    };
    var decString = (buf, off, size) => buf.slice(off, off + size).toString("utf8").replace(/\0.*/, "");
    var decDate = (buf, off, size) => numToDate(decNumber(buf, off, size));
    var numToDate = (num) => num === null ? null : new Date(num * 1e3);
    var decNumber = (buf, off, size) => buf[off] & 128 ? large.parse(buf.slice(off, off + size)) : decSmallNumber(buf, off, size);
    var nanNull = (value) => isNaN(value) ? null : value;
    var decSmallNumber = (buf, off, size) => nanNull(parseInt(
      buf.slice(off, off + size).toString("utf8").replace(/\0.*$/, "").trim(),
      8
    ));
    var MAXNUM = {
      12: 8589934591,
      8: 2097151
    };
    var encNumber = (buf, off, size, number) => number === null ? false : number > MAXNUM[size] || number < 0 ? (large.encode(number, buf.slice(off, off + size)), true) : (encSmallNumber(buf, off, size, number), false);
    var encSmallNumber = (buf, off, size, number) => buf.write(octalString(number, size), off, size, "ascii");
    var octalString = (number, size) => padOctal(Math.floor(number).toString(8), size);
    var padOctal = (string, size) => (string.length === size - 1 ? string : new Array(size - string.length - 1).join("0") + string + " ") + "\0";
    var encDate = (buf, off, size, date) => date === null ? false : encNumber(buf, off, size, date.getTime() / 1e3);
    var NULLS = new Array(156).join("\0");
    var encString = (buf, off, size, string) => string === null ? false : (buf.write(string + NULLS, off, size, "utf8"), string.length !== Buffer.byteLength(string) || string.length > size);
    module2.exports = Header;
  }
});

// node_modules/tar/lib/pax.js
var require_pax = __commonJS({
  "node_modules/tar/lib/pax.js"(exports, module2) {
    "use strict";
    var Header = require_header();
    var path32 = require("path");
    var Pax = class {
      constructor(obj, global2) {
        this.atime = obj.atime || null;
        this.charset = obj.charset || null;
        this.comment = obj.comment || null;
        this.ctime = obj.ctime || null;
        this.gid = obj.gid || null;
        this.gname = obj.gname || null;
        this.linkpath = obj.linkpath || null;
        this.mtime = obj.mtime || null;
        this.path = obj.path || null;
        this.size = obj.size || null;
        this.uid = obj.uid || null;
        this.uname = obj.uname || null;
        this.dev = obj.dev || null;
        this.ino = obj.ino || null;
        this.nlink = obj.nlink || null;
        this.global = global2 || false;
      }
      encode() {
        const body = this.encodeBody();
        if (body === "") {
          return null;
        }
        const bodyLen = Buffer.byteLength(body);
        const bufLen = 512 * Math.ceil(1 + bodyLen / 512);
        const buf = Buffer.allocUnsafe(bufLen);
        for (let i = 0; i < 512; i++) {
          buf[i] = 0;
        }
        new Header({
          // XXX split the path
          // then the path should be PaxHeader + basename, but less than 99,
          // prepend with the dirname
          path: ("PaxHeader/" + path32.basename(this.path)).slice(0, 99),
          mode: this.mode || 420,
          uid: this.uid || null,
          gid: this.gid || null,
          size: bodyLen,
          mtime: this.mtime || null,
          type: this.global ? "GlobalExtendedHeader" : "ExtendedHeader",
          linkpath: "",
          uname: this.uname || "",
          gname: this.gname || "",
          devmaj: 0,
          devmin: 0,
          atime: this.atime || null,
          ctime: this.ctime || null
        }).encode(buf);
        buf.write(body, 512, bodyLen, "utf8");
        for (let i = bodyLen + 512; i < buf.length; i++) {
          buf[i] = 0;
        }
        return buf;
      }
      encodeBody() {
        return this.encodeField("path") + this.encodeField("ctime") + this.encodeField("atime") + this.encodeField("dev") + this.encodeField("ino") + this.encodeField("nlink") + this.encodeField("charset") + this.encodeField("comment") + this.encodeField("gid") + this.encodeField("gname") + this.encodeField("linkpath") + this.encodeField("mtime") + this.encodeField("size") + this.encodeField("uid") + this.encodeField("uname");
      }
      encodeField(field) {
        if (this[field] === null || this[field] === void 0) {
          return "";
        }
        const v = this[field] instanceof Date ? this[field].getTime() / 1e3 : this[field];
        const s = " " + (field === "dev" || field === "ino" || field === "nlink" ? "SCHILY." : "") + field + "=" + v + "\n";
        const byteLen = Buffer.byteLength(s);
        let digits = Math.floor(Math.log(byteLen) / Math.log(10)) + 1;
        if (byteLen + digits >= Math.pow(10, digits)) {
          digits += 1;
        }
        const len = digits + byteLen;
        return len + s;
      }
    };
    Pax.parse = (string, ex, g) => new Pax(merge2(parseKV(string), ex), g);
    var merge2 = (a, b) => b ? Object.keys(a).reduce((s, k) => (s[k] = a[k], s), b) : a;
    var parseKV = (string) => string.replace(/\n$/, "").split("\n").reduce(parseKVLine, /* @__PURE__ */ Object.create(null));
    var parseKVLine = (set2, line) => {
      const n = parseInt(line, 10);
      if (n !== Buffer.byteLength(line) + 1) {
        return set2;
      }
      line = line.slice((n + " ").length);
      const kv = line.split("=");
      const k = kv.shift().replace(/^SCHILY\.(dev|ino|nlink)/, "$1");
      if (!k) {
        return set2;
      }
      const v = kv.join("=");
      set2[k] = /^([A-Z]+\.)?([mac]|birth|creation)time$/.test(k) ? new Date(v * 1e3) : /^[0-9]+$/.test(v) ? +v : v;
      return set2;
    };
    module2.exports = Pax;
  }
});

// node_modules/tar/lib/strip-trailing-slashes.js
var require_strip_trailing_slashes = __commonJS({
  "node_modules/tar/lib/strip-trailing-slashes.js"(exports, module2) {
    module2.exports = (str2) => {
      let i = str2.length - 1;
      let slashesStart = -1;
      while (i > -1 && str2.charAt(i) === "/") {
        slashesStart = i;
        i--;
      }
      return slashesStart === -1 ? str2 : str2.slice(0, slashesStart);
    };
  }
});

// node_modules/tar/lib/warn-mixin.js
var require_warn_mixin = __commonJS({
  "node_modules/tar/lib/warn-mixin.js"(exports, module2) {
    "use strict";
    module2.exports = (Base) => class extends Base {
      warn(code, message, data = {}) {
        if (this.file) {
          data.file = this.file;
        }
        if (this.cwd) {
          data.cwd = this.cwd;
        }
        data.code = message instanceof Error && message.code || code;
        data.tarCode = code;
        if (!this.strict && data.recoverable !== false) {
          if (message instanceof Error) {
            data = Object.assign(message, data);
            message = message.message;
          }
          this.emit("warn", data.tarCode, message, data);
        } else if (message instanceof Error) {
          this.emit("error", Object.assign(message, data));
        } else {
          this.emit("error", Object.assign(new Error(`${code}: ${message}`), data));
        }
      }
    };
  }
});

// node_modules/tar/lib/winchars.js
var require_winchars = __commonJS({
  "node_modules/tar/lib/winchars.js"(exports, module2) {
    "use strict";
    var raw = [
      "|",
      "<",
      ">",
      "?",
      ":"
    ];
    var win = raw.map((char) => String.fromCharCode(61440 + char.charCodeAt(0)));
    var toWin = new Map(raw.map((char, i) => [char, win[i]]));
    var toRaw = new Map(win.map((char, i) => [char, raw[i]]));
    module2.exports = {
      encode: (s) => raw.reduce((s2, c2) => s2.split(c2).join(toWin.get(c2)), s),
      decode: (s) => win.reduce((s2, c2) => s2.split(c2).join(toRaw.get(c2)), s)
    };
  }
});

// node_modules/tar/lib/strip-absolute-path.js
var require_strip_absolute_path = __commonJS({
  "node_modules/tar/lib/strip-absolute-path.js"(exports, module2) {
    var { isAbsolute: isAbsolute3, parse: parse5 } = require("path").win32;
    module2.exports = (path32) => {
      let r = "";
      let parsed = parse5(path32);
      while (isAbsolute3(path32) || parsed.root) {
        const root = path32.charAt(0) === "/" && path32.slice(0, 4) !== "//?/" ? "/" : parsed.root;
        path32 = path32.slice(root.length);
        r += root;
        parsed = parse5(path32);
      }
      return [r, path32];
    };
  }
});

// node_modules/tar/lib/mode-fix.js
var require_mode_fix = __commonJS({
  "node_modules/tar/lib/mode-fix.js"(exports, module2) {
    "use strict";
    module2.exports = (mode, isDir, portable) => {
      mode &= 4095;
      if (portable) {
        mode = (mode | 384) & ~18;
      }
      if (isDir) {
        if (mode & 256) {
          mode |= 64;
        }
        if (mode & 32) {
          mode |= 8;
        }
        if (mode & 4) {
          mode |= 1;
        }
      }
      return mode;
    };
  }
});

// node_modules/tar/lib/write-entry.js
var require_write_entry = __commonJS({
  "node_modules/tar/lib/write-entry.js"(exports, module2) {
    "use strict";
    var MiniPass = require_minipass();
    var Pax = require_pax();
    var Header = require_header();
    var fs6 = require("fs");
    var path32 = require("path");
    var normPath = require_normalize_windows_path();
    var stripSlash = require_strip_trailing_slashes();
    var prefixPath = (path33, prefix) => {
      if (!prefix) {
        return normPath(path33);
      }
      path33 = normPath(path33).replace(/^\.(\/|$)/, "");
      return stripSlash(prefix) + "/" + path33;
    };
    var maxReadSize = 16 * 1024 * 1024;
    var PROCESS = Symbol("process");
    var FILE = Symbol("file");
    var DIRECTORY = Symbol("directory");
    var SYMLINK = Symbol("symlink");
    var HARDLINK = Symbol("hardlink");
    var HEADER = Symbol("header");
    var READ = Symbol("read");
    var LSTAT = Symbol("lstat");
    var ONLSTAT = Symbol("onlstat");
    var ONREAD = Symbol("onread");
    var ONREADLINK = Symbol("onreadlink");
    var OPENFILE = Symbol("openfile");
    var ONOPENFILE = Symbol("onopenfile");
    var CLOSE = Symbol("close");
    var MODE = Symbol("mode");
    var AWAITDRAIN = Symbol("awaitDrain");
    var ONDRAIN = Symbol("ondrain");
    var PREFIX = Symbol("prefix");
    var HAD_ERROR = Symbol("hadError");
    var warner = require_warn_mixin();
    var winchars = require_winchars();
    var stripAbsolutePath = require_strip_absolute_path();
    var modeFix = require_mode_fix();
    var WriteEntry = warner(class WriteEntry extends MiniPass {
      constructor(p, opt) {
        opt = opt || {};
        super(opt);
        if (typeof p !== "string") {
          throw new TypeError("path is required");
        }
        this.path = normPath(p);
        this.portable = !!opt.portable;
        this.myuid = process.getuid && process.getuid() || 0;
        this.myuser = process.env.USER || "";
        this.maxReadSize = opt.maxReadSize || maxReadSize;
        this.linkCache = opt.linkCache || /* @__PURE__ */ new Map();
        this.statCache = opt.statCache || /* @__PURE__ */ new Map();
        this.preservePaths = !!opt.preservePaths;
        this.cwd = normPath(opt.cwd || process.cwd());
        this.strict = !!opt.strict;
        this.noPax = !!opt.noPax;
        this.noMtime = !!opt.noMtime;
        this.mtime = opt.mtime || null;
        this.prefix = opt.prefix ? normPath(opt.prefix) : null;
        this.fd = null;
        this.blockLen = null;
        this.blockRemain = null;
        this.buf = null;
        this.offset = null;
        this.length = null;
        this.pos = null;
        this.remain = null;
        if (typeof opt.onwarn === "function") {
          this.on("warn", opt.onwarn);
        }
        let pathWarn = false;
        if (!this.preservePaths) {
          const [root, stripped] = stripAbsolutePath(this.path);
          if (root) {
            this.path = stripped;
            pathWarn = root;
          }
        }
        this.win32 = !!opt.win32 || process.platform === "win32";
        if (this.win32) {
          this.path = winchars.decode(this.path.replace(/\\/g, "/"));
          p = p.replace(/\\/g, "/");
        }
        this.absolute = normPath(opt.absolute || path32.resolve(this.cwd, p));
        if (this.path === "") {
          this.path = "./";
        }
        if (pathWarn) {
          this.warn("TAR_ENTRY_INFO", `stripping ${pathWarn} from absolute path`, {
            entry: this,
            path: pathWarn + this.path
          });
        }
        if (this.statCache.has(this.absolute)) {
          this[ONLSTAT](this.statCache.get(this.absolute));
        } else {
          this[LSTAT]();
        }
      }
      emit(ev, ...data) {
        if (ev === "error") {
          this[HAD_ERROR] = true;
        }
        return super.emit(ev, ...data);
      }
      [LSTAT]() {
        fs6.lstat(this.absolute, (er, stat2) => {
          if (er) {
            return this.emit("error", er);
          }
          this[ONLSTAT](stat2);
        });
      }
      [ONLSTAT](stat2) {
        this.statCache.set(this.absolute, stat2);
        this.stat = stat2;
        if (!stat2.isFile()) {
          stat2.size = 0;
        }
        this.type = getType(stat2);
        this.emit("stat", stat2);
        this[PROCESS]();
      }
      [PROCESS]() {
        switch (this.type) {
          case "File":
            return this[FILE]();
          case "Directory":
            return this[DIRECTORY]();
          case "SymbolicLink":
            return this[SYMLINK]();
          default:
            return this.end();
        }
      }
      [MODE](mode) {
        return modeFix(mode, this.type === "Directory", this.portable);
      }
      [PREFIX](path33) {
        return prefixPath(path33, this.prefix);
      }
      [HEADER]() {
        if (this.type === "Directory" && this.portable) {
          this.noMtime = true;
        }
        this.header = new Header({
          path: this[PREFIX](this.path),
          // only apply the prefix to hard links.
          linkpath: this.type === "Link" ? this[PREFIX](this.linkpath) : this.linkpath,
          // only the permissions and setuid/setgid/sticky bitflags
          // not the higher-order bits that specify file type
          mode: this[MODE](this.stat.mode),
          uid: this.portable ? null : this.stat.uid,
          gid: this.portable ? null : this.stat.gid,
          size: this.stat.size,
          mtime: this.noMtime ? null : this.mtime || this.stat.mtime,
          type: this.type,
          uname: this.portable ? null : this.stat.uid === this.myuid ? this.myuser : "",
          atime: this.portable ? null : this.stat.atime,
          ctime: this.portable ? null : this.stat.ctime
        });
        if (this.header.encode() && !this.noPax) {
          super.write(new Pax({
            atime: this.portable ? null : this.header.atime,
            ctime: this.portable ? null : this.header.ctime,
            gid: this.portable ? null : this.header.gid,
            mtime: this.noMtime ? null : this.mtime || this.header.mtime,
            path: this[PREFIX](this.path),
            linkpath: this.type === "Link" ? this[PREFIX](this.linkpath) : this.linkpath,
            size: this.header.size,
            uid: this.portable ? null : this.header.uid,
            uname: this.portable ? null : this.header.uname,
            dev: this.portable ? null : this.stat.dev,
            ino: this.portable ? null : this.stat.ino,
            nlink: this.portable ? null : this.stat.nlink
          }).encode());
        }
        super.write(this.header.block);
      }
      [DIRECTORY]() {
        if (this.path.slice(-1) !== "/") {
          this.path += "/";
        }
        this.stat.size = 0;
        this[HEADER]();
        this.end();
      }
      [SYMLINK]() {
        fs6.readlink(this.absolute, (er, linkpath) => {
          if (er) {
            return this.emit("error", er);
          }
          this[ONREADLINK](linkpath);
        });
      }
      [ONREADLINK](linkpath) {
        this.linkpath = normPath(linkpath);
        this[HEADER]();
        this.end();
      }
      [HARDLINK](linkpath) {
        this.type = "Link";
        this.linkpath = normPath(path32.relative(this.cwd, linkpath));
        this.stat.size = 0;
        this[HEADER]();
        this.end();
      }
      [FILE]() {
        if (this.stat.nlink > 1) {
          const linkKey = this.stat.dev + ":" + this.stat.ino;
          if (this.linkCache.has(linkKey)) {
            const linkpath = this.linkCache.get(linkKey);
            if (linkpath.indexOf(this.cwd) === 0) {
              return this[HARDLINK](linkpath);
            }
          }
          this.linkCache.set(linkKey, this.absolute);
        }
        this[HEADER]();
        if (this.stat.size === 0) {
          return this.end();
        }
        this[OPENFILE]();
      }
      [OPENFILE]() {
        fs6.open(this.absolute, "r", (er, fd) => {
          if (er) {
            return this.emit("error", er);
          }
          this[ONOPENFILE](fd);
        });
      }
      [ONOPENFILE](fd) {
        this.fd = fd;
        if (this[HAD_ERROR]) {
          return this[CLOSE]();
        }
        this.blockLen = 512 * Math.ceil(this.stat.size / 512);
        this.blockRemain = this.blockLen;
        const bufLen = Math.min(this.blockLen, this.maxReadSize);
        this.buf = Buffer.allocUnsafe(bufLen);
        this.offset = 0;
        this.pos = 0;
        this.remain = this.stat.size;
        this.length = this.buf.length;
        this[READ]();
      }
      [READ]() {
        const { fd, buf, offset, length, pos } = this;
        fs6.read(fd, buf, offset, length, pos, (er, bytesRead) => {
          if (er) {
            return this[CLOSE](() => this.emit("error", er));
          }
          this[ONREAD](bytesRead);
        });
      }
      [CLOSE](cb) {
        fs6.close(this.fd, cb);
      }
      [ONREAD](bytesRead) {
        if (bytesRead <= 0 && this.remain > 0) {
          const er = new Error("encountered unexpected EOF");
          er.path = this.absolute;
          er.syscall = "read";
          er.code = "EOF";
          return this[CLOSE](() => this.emit("error", er));
        }
        if (bytesRead > this.remain) {
          const er = new Error("did not encounter expected EOF");
          er.path = this.absolute;
          er.syscall = "read";
          er.code = "EOF";
          return this[CLOSE](() => this.emit("error", er));
        }
        if (bytesRead === this.remain) {
          for (let i = bytesRead; i < this.length && bytesRead < this.blockRemain; i++) {
            this.buf[i + this.offset] = 0;
            bytesRead++;
            this.remain++;
          }
        }
        const writeBuf = this.offset === 0 && bytesRead === this.buf.length ? this.buf : this.buf.slice(this.offset, this.offset + bytesRead);
        const flushed = this.write(writeBuf);
        if (!flushed) {
          this[AWAITDRAIN](() => this[ONDRAIN]());
        } else {
          this[ONDRAIN]();
        }
      }
      [AWAITDRAIN](cb) {
        this.once("drain", cb);
      }
      write(writeBuf) {
        if (this.blockRemain < writeBuf.length) {
          const er = new Error("writing more data than expected");
          er.path = this.absolute;
          return this.emit("error", er);
        }
        this.remain -= writeBuf.length;
        this.blockRemain -= writeBuf.length;
        this.pos += writeBuf.length;
        this.offset += writeBuf.length;
        return super.write(writeBuf);
      }
      [ONDRAIN]() {
        if (!this.remain) {
          if (this.blockRemain) {
            super.write(Buffer.alloc(this.blockRemain));
          }
          return this[CLOSE]((er) => er ? this.emit("error", er) : this.end());
        }
        if (this.offset >= this.length) {
          this.buf = Buffer.allocUnsafe(Math.min(this.blockRemain, this.buf.length));
          this.offset = 0;
        }
        this.length = this.buf.length - this.offset;
        this[READ]();
      }
    });
    var WriteEntrySync = class extends WriteEntry {
      [LSTAT]() {
        this[ONLSTAT](fs6.lstatSync(this.absolute));
      }
      [SYMLINK]() {
        this[ONREADLINK](fs6.readlinkSync(this.absolute));
      }
      [OPENFILE]() {
        this[ONOPENFILE](fs6.openSync(this.absolute, "r"));
      }
      [READ]() {
        let threw = true;
        try {
          const { fd, buf, offset, length, pos } = this;
          const bytesRead = fs6.readSync(fd, buf, offset, length, pos);
          this[ONREAD](bytesRead);
          threw = false;
        } finally {
          if (threw) {
            try {
              this[CLOSE](() => {
              });
            } catch (er) {
            }
          }
        }
      }
      [AWAITDRAIN](cb) {
        cb();
      }
      [CLOSE](cb) {
        fs6.closeSync(this.fd);
        cb();
      }
    };
    var WriteEntryTar = warner(class WriteEntryTar extends MiniPass {
      constructor(readEntry, opt) {
        opt = opt || {};
        super(opt);
        this.preservePaths = !!opt.preservePaths;
        this.portable = !!opt.portable;
        this.strict = !!opt.strict;
        this.noPax = !!opt.noPax;
        this.noMtime = !!opt.noMtime;
        this.readEntry = readEntry;
        this.type = readEntry.type;
        if (this.type === "Directory" && this.portable) {
          this.noMtime = true;
        }
        this.prefix = opt.prefix || null;
        this.path = normPath(readEntry.path);
        this.mode = this[MODE](readEntry.mode);
        this.uid = this.portable ? null : readEntry.uid;
        this.gid = this.portable ? null : readEntry.gid;
        this.uname = this.portable ? null : readEntry.uname;
        this.gname = this.portable ? null : readEntry.gname;
        this.size = readEntry.size;
        this.mtime = this.noMtime ? null : opt.mtime || readEntry.mtime;
        this.atime = this.portable ? null : readEntry.atime;
        this.ctime = this.portable ? null : readEntry.ctime;
        this.linkpath = normPath(readEntry.linkpath);
        if (typeof opt.onwarn === "function") {
          this.on("warn", opt.onwarn);
        }
        let pathWarn = false;
        if (!this.preservePaths) {
          const [root, stripped] = stripAbsolutePath(this.path);
          if (root) {
            this.path = stripped;
            pathWarn = root;
          }
        }
        this.remain = readEntry.size;
        this.blockRemain = readEntry.startBlockSize;
        this.header = new Header({
          path: this[PREFIX](this.path),
          linkpath: this.type === "Link" ? this[PREFIX](this.linkpath) : this.linkpath,
          // only the permissions and setuid/setgid/sticky bitflags
          // not the higher-order bits that specify file type
          mode: this.mode,
          uid: this.portable ? null : this.uid,
          gid: this.portable ? null : this.gid,
          size: this.size,
          mtime: this.noMtime ? null : this.mtime,
          type: this.type,
          uname: this.portable ? null : this.uname,
          atime: this.portable ? null : this.atime,
          ctime: this.portable ? null : this.ctime
        });
        if (pathWarn) {
          this.warn("TAR_ENTRY_INFO", `stripping ${pathWarn} from absolute path`, {
            entry: this,
            path: pathWarn + this.path
          });
        }
        if (this.header.encode() && !this.noPax) {
          super.write(new Pax({
            atime: this.portable ? null : this.atime,
            ctime: this.portable ? null : this.ctime,
            gid: this.portable ? null : this.gid,
            mtime: this.noMtime ? null : this.mtime,
            path: this[PREFIX](this.path),
            linkpath: this.type === "Link" ? this[PREFIX](this.linkpath) : this.linkpath,
            size: this.size,
            uid: this.portable ? null : this.uid,
            uname: this.portable ? null : this.uname,
            dev: this.portable ? null : this.readEntry.dev,
            ino: this.portable ? null : this.readEntry.ino,
            nlink: this.portable ? null : this.readEntry.nlink
          }).encode());
        }
        super.write(this.header.block);
        readEntry.pipe(this);
      }
      [PREFIX](path33) {
        return prefixPath(path33, this.prefix);
      }
      [MODE](mode) {
        return modeFix(mode, this.type === "Directory", this.portable);
      }
      write(data) {
        const writeLen = data.length;
        if (writeLen > this.blockRemain) {
          throw new Error("writing more to entry than is appropriate");
        }
        this.blockRemain -= writeLen;
        return super.write(data);
      }
      end() {
        if (this.blockRemain) {
          super.write(Buffer.alloc(this.blockRemain));
        }
        return super.end();
      }
    });
    WriteEntry.Sync = WriteEntrySync;
    WriteEntry.Tar = WriteEntryTar;
    var getType = (stat2) => stat2.isFile() ? "File" : stat2.isDirectory() ? "Directory" : stat2.isSymbolicLink() ? "SymbolicLink" : "Unsupported";
    module2.exports = WriteEntry;
  }
});

// node_modules/yallist/iterator.js
var require_iterator = __commonJS({
  "node_modules/yallist/iterator.js"(exports, module2) {
    "use strict";
    module2.exports = function(Yallist) {
      Yallist.prototype[Symbol.iterator] = function* () {
        for (let walker = this.head; walker; walker = walker.next) {
          yield walker.value;
        }
      };
    };
  }
});

// node_modules/yallist/yallist.js
var require_yallist = __commonJS({
  "node_modules/yallist/yallist.js"(exports, module2) {
    "use strict";
    module2.exports = Yallist;
    Yallist.Node = Node;
    Yallist.create = Yallist;
    function Yallist(list) {
      var self2 = this;
      if (!(self2 instanceof Yallist)) {
        self2 = new Yallist();
      }
      self2.tail = null;
      self2.head = null;
      self2.length = 0;
      if (list && typeof list.forEach === "function") {
        list.forEach(function(item) {
          self2.push(item);
        });
      } else if (arguments.length > 0) {
        for (var i = 0, l = arguments.length; i < l; i++) {
          self2.push(arguments[i]);
        }
      }
      return self2;
    }
    Yallist.prototype.removeNode = function(node) {
      if (node.list !== this) {
        throw new Error("removing node which does not belong to this list");
      }
      var next = node.next;
      var prev = node.prev;
      if (next) {
        next.prev = prev;
      }
      if (prev) {
        prev.next = next;
      }
      if (node === this.head) {
        this.head = next;
      }
      if (node === this.tail) {
        this.tail = prev;
      }
      node.list.length--;
      node.next = null;
      node.prev = null;
      node.list = null;
      return next;
    };
    Yallist.prototype.unshiftNode = function(node) {
      if (node === this.head) {
        return;
      }
      if (node.list) {
        node.list.removeNode(node);
      }
      var head = this.head;
      node.list = this;
      node.next = head;
      if (head) {
        head.prev = node;
      }
      this.head = node;
      if (!this.tail) {
        this.tail = node;
      }
      this.length++;
    };
    Yallist.prototype.pushNode = function(node) {
      if (node === this.tail) {
        return;
      }
      if (node.list) {
        node.list.removeNode(node);
      }
      var tail = this.tail;
      node.list = this;
      node.prev = tail;
      if (tail) {
        tail.next = node;
      }
      this.tail = node;
      if (!this.head) {
        this.head = node;
      }
      this.length++;
    };
    Yallist.prototype.push = function() {
      for (var i = 0, l = arguments.length; i < l; i++) {
        push(this, arguments[i]);
      }
      return this.length;
    };
    Yallist.prototype.unshift = function() {
      for (var i = 0, l = arguments.length; i < l; i++) {
        unshift(this, arguments[i]);
      }
      return this.length;
    };
    Yallist.prototype.pop = function() {
      if (!this.tail) {
        return void 0;
      }
      var res = this.tail.value;
      this.tail = this.tail.prev;
      if (this.tail) {
        this.tail.next = null;
      } else {
        this.head = null;
      }
      this.length--;
      return res;
    };
    Yallist.prototype.shift = function() {
      if (!this.head) {
        return void 0;
      }
      var res = this.head.value;
      this.head = this.head.next;
      if (this.head) {
        this.head.prev = null;
      } else {
        this.tail = null;
      }
      this.length--;
      return res;
    };
    Yallist.prototype.forEach = function(fn, thisp) {
      thisp = thisp || this;
      for (var walker = this.head, i = 0; walker !== null; i++) {
        fn.call(thisp, walker.value, i, this);
        walker = walker.next;
      }
    };
    Yallist.prototype.forEachReverse = function(fn, thisp) {
      thisp = thisp || this;
      for (var walker = this.tail, i = this.length - 1; walker !== null; i--) {
        fn.call(thisp, walker.value, i, this);
        walker = walker.prev;
      }
    };
    Yallist.prototype.get = function(n) {
      for (var i = 0, walker = this.head; walker !== null && i < n; i++) {
        walker = walker.next;
      }
      if (i === n && walker !== null) {
        return walker.value;
      }
    };
    Yallist.prototype.getReverse = function(n) {
      for (var i = 0, walker = this.tail; walker !== null && i < n; i++) {
        walker = walker.prev;
      }
      if (i === n && walker !== null) {
        return walker.value;
      }
    };
    Yallist.prototype.map = function(fn, thisp) {
      thisp = thisp || this;
      var res = new Yallist();
      for (var walker = this.head; walker !== null; ) {
        res.push(fn.call(thisp, walker.value, this));
        walker = walker.next;
      }
      return res;
    };
    Yallist.prototype.mapReverse = function(fn, thisp) {
      thisp = thisp || this;
      var res = new Yallist();
      for (var walker = this.tail; walker !== null; ) {
        res.push(fn.call(thisp, walker.value, this));
        walker = walker.prev;
      }
      return res;
    };
    Yallist.prototype.reduce = function(fn, initial) {
      var acc;
      var walker = this.head;
      if (arguments.length > 1) {
        acc = initial;
      } else if (this.head) {
        walker = this.head.next;
        acc = this.head.value;
      } else {
        throw new TypeError("Reduce of empty list with no initial value");
      }
      for (var i = 0; walker !== null; i++) {
        acc = fn(acc, walker.value, i);
        walker = walker.next;
      }
      return acc;
    };
    Yallist.prototype.reduceReverse = function(fn, initial) {
      var acc;
      var walker = this.tail;
      if (arguments.length > 1) {
        acc = initial;
      } else if (this.tail) {
        walker = this.tail.prev;
        acc = this.tail.value;
      } else {
        throw new TypeError("Reduce of empty list with no initial value");
      }
      for (var i = this.length - 1; walker !== null; i--) {
        acc = fn(acc, walker.value, i);
        walker = walker.prev;
      }
      return acc;
    };
    Yallist.prototype.toArray = function() {
      var arr = new Array(this.length);
      for (var i = 0, walker = this.head; walker !== null; i++) {
        arr[i] = walker.value;
        walker = walker.next;
      }
      return arr;
    };
    Yallist.prototype.toArrayReverse = function() {
      var arr = new Array(this.length);
      for (var i = 0, walker = this.tail; walker !== null; i++) {
        arr[i] = walker.value;
        walker = walker.prev;
      }
      return arr;
    };
    Yallist.prototype.slice = function(from, to) {
      to = to || this.length;
      if (to < 0) {
        to += this.length;
      }
      from = from || 0;
      if (from < 0) {
        from += this.length;
      }
      var ret = new Yallist();
      if (to < from || to < 0) {
        return ret;
      }
      if (from < 0) {
        from = 0;
      }
      if (to > this.length) {
        to = this.length;
      }
      for (var i = 0, walker = this.head; walker !== null && i < from; i++) {
        walker = walker.next;
      }
      for (; walker !== null && i < to; i++, walker = walker.next) {
        ret.push(walker.value);
      }
      return ret;
    };
    Yallist.prototype.sliceReverse = function(from, to) {
      to = to || this.length;
      if (to < 0) {
        to += this.length;
      }
      from = from || 0;
      if (from < 0) {
        from += this.length;
      }
      var ret = new Yallist();
      if (to < from || to < 0) {
        return ret;
      }
      if (from < 0) {
        from = 0;
      }
      if (to > this.length) {
        to = this.length;
      }
      for (var i = this.length, walker = this.tail; walker !== null && i > to; i--) {
        walker = walker.prev;
      }
      for (; walker !== null && i > from; i--, walker = walker.prev) {
        ret.push(walker.value);
      }
      return ret;
    };
    Yallist.prototype.splice = function(start, deleteCount, ...nodes) {
      if (start > this.length) {
        start = this.length - 1;
      }
      if (start < 0) {
        start = this.length + start;
      }
      for (var i = 0, walker = this.head; walker !== null && i < start; i++) {
        walker = walker.next;
      }
      var ret = [];
      for (var i = 0; walker && i < deleteCount; i++) {
        ret.push(walker.value);
        walker = this.removeNode(walker);
      }
      if (walker === null) {
        walker = this.tail;
      }
      if (walker !== this.head && walker !== this.tail) {
        walker = walker.prev;
      }
      for (var i = 0; i < nodes.length; i++) {
        walker = insert(this, walker, nodes[i]);
      }
      return ret;
    };
    Yallist.prototype.reverse = function() {
      var head = this.head;
      var tail = this.tail;
      for (var walker = head; walker !== null; walker = walker.prev) {
        var p = walker.prev;
        walker.prev = walker.next;
        walker.next = p;
      }
      this.head = tail;
      this.tail = head;
      return this;
    };
    function insert(self2, node, value) {
      var inserted = node === self2.head ? new Node(value, null, node, self2) : new Node(value, node, node.next, self2);
      if (inserted.next === null) {
        self2.tail = inserted;
      }
      if (inserted.prev === null) {
        self2.head = inserted;
      }
      self2.length++;
      return inserted;
    }
    function push(self2, item) {
      self2.tail = new Node(item, self2.tail, null, self2);
      if (!self2.head) {
        self2.head = self2.tail;
      }
      self2.length++;
    }
    function unshift(self2, item) {
      self2.head = new Node(item, null, self2.head, self2);
      if (!self2.tail) {
        self2.tail = self2.head;
      }
      self2.length++;
    }
    function Node(value, prev, next, list) {
      if (!(this instanceof Node)) {
        return new Node(value, prev, next, list);
      }
      this.list = list;
      this.value = value;
      if (prev) {
        prev.next = this;
        this.prev = prev;
      } else {
        this.prev = null;
      }
      if (next) {
        next.prev = this;
        this.next = next;
      } else {
        this.next = null;
      }
    }
    try {
      require_iterator()(Yallist);
    } catch (er) {
    }
  }
});

// node_modules/tar/lib/pack.js
var require_pack = __commonJS({
  "node_modules/tar/lib/pack.js"(exports, module2) {
    "use strict";
    var PackJob = class {
      constructor(path33, absolute) {
        this.path = path33 || "./";
        this.absolute = absolute;
        this.entry = null;
        this.stat = null;
        this.readdir = null;
        this.pending = false;
        this.ignore = false;
        this.piped = false;
      }
    };
    var MiniPass = require_minipass();
    var zlib = require_minizlib();
    var ReadEntry = require_read_entry();
    var WriteEntry = require_write_entry();
    var WriteEntrySync = WriteEntry.Sync;
    var WriteEntryTar = WriteEntry.Tar;
    var Yallist = require_yallist();
    var EOF = Buffer.alloc(1024);
    var ONSTAT = Symbol("onStat");
    var ENDED = Symbol("ended");
    var QUEUE = Symbol("queue");
    var CURRENT = Symbol("current");
    var PROCESS = Symbol("process");
    var PROCESSING = Symbol("processing");
    var PROCESSJOB = Symbol("processJob");
    var JOBS = Symbol("jobs");
    var JOBDONE = Symbol("jobDone");
    var ADDFSENTRY = Symbol("addFSEntry");
    var ADDTARENTRY = Symbol("addTarEntry");
    var STAT = Symbol("stat");
    var READDIR = Symbol("readdir");
    var ONREADDIR = Symbol("onreaddir");
    var PIPE = Symbol("pipe");
    var ENTRY = Symbol("entry");
    var ENTRYOPT = Symbol("entryOpt");
    var WRITEENTRYCLASS = Symbol("writeEntryClass");
    var WRITE = Symbol("write");
    var ONDRAIN = Symbol("ondrain");
    var fs6 = require("fs");
    var path32 = require("path");
    var warner = require_warn_mixin();
    var normPath = require_normalize_windows_path();
    var Pack = warner(class Pack extends MiniPass {
      constructor(opt) {
        super(opt);
        opt = opt || /* @__PURE__ */ Object.create(null);
        this.opt = opt;
        this.file = opt.file || "";
        this.cwd = opt.cwd || process.cwd();
        this.maxReadSize = opt.maxReadSize;
        this.preservePaths = !!opt.preservePaths;
        this.strict = !!opt.strict;
        this.noPax = !!opt.noPax;
        this.prefix = normPath(opt.prefix || "");
        this.linkCache = opt.linkCache || /* @__PURE__ */ new Map();
        this.statCache = opt.statCache || /* @__PURE__ */ new Map();
        this.readdirCache = opt.readdirCache || /* @__PURE__ */ new Map();
        this[WRITEENTRYCLASS] = WriteEntry;
        if (typeof opt.onwarn === "function") {
          this.on("warn", opt.onwarn);
        }
        this.portable = !!opt.portable;
        this.zip = null;
        if (opt.gzip) {
          if (typeof opt.gzip !== "object") {
            opt.gzip = {};
          }
          if (this.portable) {
            opt.gzip.portable = true;
          }
          this.zip = new zlib.Gzip(opt.gzip);
          this.zip.on("data", (chunk) => super.write(chunk));
          this.zip.on("end", (_) => super.end());
          this.zip.on("drain", (_) => this[ONDRAIN]());
          this.on("resume", (_) => this.zip.resume());
        } else {
          this.on("drain", this[ONDRAIN]);
        }
        this.noDirRecurse = !!opt.noDirRecurse;
        this.follow = !!opt.follow;
        this.noMtime = !!opt.noMtime;
        this.mtime = opt.mtime || null;
        this.filter = typeof opt.filter === "function" ? opt.filter : (_) => true;
        this[QUEUE] = new Yallist();
        this[JOBS] = 0;
        this.jobs = +opt.jobs || 4;
        this[PROCESSING] = false;
        this[ENDED] = false;
      }
      [WRITE](chunk) {
        return super.write(chunk);
      }
      add(path33) {
        this.write(path33);
        return this;
      }
      end(path33) {
        if (path33) {
          this.write(path33);
        }
        this[ENDED] = true;
        this[PROCESS]();
        return this;
      }
      write(path33) {
        if (this[ENDED]) {
          throw new Error("write after end");
        }
        if (path33 instanceof ReadEntry) {
          this[ADDTARENTRY](path33);
        } else {
          this[ADDFSENTRY](path33);
        }
        return this.flowing;
      }
      [ADDTARENTRY](p) {
        const absolute = normPath(path32.resolve(this.cwd, p.path));
        if (!this.filter(p.path, p)) {
          p.resume();
        } else {
          const job = new PackJob(p.path, absolute, false);
          job.entry = new WriteEntryTar(p, this[ENTRYOPT](job));
          job.entry.on("end", (_) => this[JOBDONE](job));
          this[JOBS] += 1;
          this[QUEUE].push(job);
        }
        this[PROCESS]();
      }
      [ADDFSENTRY](p) {
        const absolute = normPath(path32.resolve(this.cwd, p));
        this[QUEUE].push(new PackJob(p, absolute));
        this[PROCESS]();
      }
      [STAT](job) {
        job.pending = true;
        this[JOBS] += 1;
        const stat2 = this.follow ? "stat" : "lstat";
        fs6[stat2](job.absolute, (er, stat3) => {
          job.pending = false;
          this[JOBS] -= 1;
          if (er) {
            this.emit("error", er);
          } else {
            this[ONSTAT](job, stat3);
          }
        });
      }
      [ONSTAT](job, stat2) {
        this.statCache.set(job.absolute, stat2);
        job.stat = stat2;
        if (!this.filter(job.path, stat2)) {
          job.ignore = true;
        }
        this[PROCESS]();
      }
      [READDIR](job) {
        job.pending = true;
        this[JOBS] += 1;
        fs6.readdir(job.absolute, (er, entries) => {
          job.pending = false;
          this[JOBS] -= 1;
          if (er) {
            return this.emit("error", er);
          }
          this[ONREADDIR](job, entries);
        });
      }
      [ONREADDIR](job, entries) {
        this.readdirCache.set(job.absolute, entries);
        job.readdir = entries;
        this[PROCESS]();
      }
      [PROCESS]() {
        if (this[PROCESSING]) {
          return;
        }
        this[PROCESSING] = true;
        for (let w = this[QUEUE].head; w !== null && this[JOBS] < this.jobs; w = w.next) {
          this[PROCESSJOB](w.value);
          if (w.value.ignore) {
            const p = w.next;
            this[QUEUE].removeNode(w);
            w.next = p;
          }
        }
        this[PROCESSING] = false;
        if (this[ENDED] && !this[QUEUE].length && this[JOBS] === 0) {
          if (this.zip) {
            this.zip.end(EOF);
          } else {
            super.write(EOF);
            super.end();
          }
        }
      }
      get [CURRENT]() {
        return this[QUEUE] && this[QUEUE].head && this[QUEUE].head.value;
      }
      [JOBDONE](job) {
        this[QUEUE].shift();
        this[JOBS] -= 1;
        this[PROCESS]();
      }
      [PROCESSJOB](job) {
        if (job.pending) {
          return;
        }
        if (job.entry) {
          if (job === this[CURRENT] && !job.piped) {
            this[PIPE](job);
          }
          return;
        }
        if (!job.stat) {
          if (this.statCache.has(job.absolute)) {
            this[ONSTAT](job, this.statCache.get(job.absolute));
          } else {
            this[STAT](job);
          }
        }
        if (!job.stat) {
          return;
        }
        if (job.ignore) {
          return;
        }
        if (!this.noDirRecurse && job.stat.isDirectory() && !job.readdir) {
          if (this.readdirCache.has(job.absolute)) {
            this[ONREADDIR](job, this.readdirCache.get(job.absolute));
          } else {
            this[READDIR](job);
          }
          if (!job.readdir) {
            return;
          }
        }
        job.entry = this[ENTRY](job);
        if (!job.entry) {
          job.ignore = true;
          return;
        }
        if (job === this[CURRENT] && !job.piped) {
          this[PIPE](job);
        }
      }
      [ENTRYOPT](job) {
        return {
          onwarn: (code, msg, data) => this.warn(code, msg, data),
          noPax: this.noPax,
          cwd: this.cwd,
          absolute: job.absolute,
          preservePaths: this.preservePaths,
          maxReadSize: this.maxReadSize,
          strict: this.strict,
          portable: this.portable,
          linkCache: this.linkCache,
          statCache: this.statCache,
          noMtime: this.noMtime,
          mtime: this.mtime,
          prefix: this.prefix
        };
      }
      [ENTRY](job) {
        this[JOBS] += 1;
        try {
          return new this[WRITEENTRYCLASS](job.path, this[ENTRYOPT](job)).on("end", () => this[JOBDONE](job)).on("error", (er) => this.emit("error", er));
        } catch (er) {
          this.emit("error", er);
        }
      }
      [ONDRAIN]() {
        if (this[CURRENT] && this[CURRENT].entry) {
          this[CURRENT].entry.resume();
        }
      }
      // like .pipe() but using super, because our write() is special
      [PIPE](job) {
        job.piped = true;
        if (job.readdir) {
          job.readdir.forEach((entry) => {
            const p = job.path;
            const base = p === "./" ? "" : p.replace(/\/*$/, "/");
            this[ADDFSENTRY](base + entry);
          });
        }
        const source = job.entry;
        const zip = this.zip;
        if (zip) {
          source.on("data", (chunk) => {
            if (!zip.write(chunk)) {
              source.pause();
            }
          });
        } else {
          source.on("data", (chunk) => {
            if (!super.write(chunk)) {
              source.pause();
            }
          });
        }
      }
      pause() {
        if (this.zip) {
          this.zip.pause();
        }
        return super.pause();
      }
    });
    var PackSync = class extends Pack {
      constructor(opt) {
        super(opt);
        this[WRITEENTRYCLASS] = WriteEntrySync;
      }
      // pause/resume are no-ops in sync streams.
      pause() {
      }
      resume() {
      }
      [STAT](job) {
        const stat2 = this.follow ? "statSync" : "lstatSync";
        this[ONSTAT](job, fs6[stat2](job.absolute));
      }
      [READDIR](job, stat2) {
        this[ONREADDIR](job, fs6.readdirSync(job.absolute));
      }
      // gotta get it all in this tick
      [PIPE](job) {
        const source = job.entry;
        const zip = this.zip;
        if (job.readdir) {
          job.readdir.forEach((entry) => {
            const p = job.path;
            const base = p === "./" ? "" : p.replace(/\/*$/, "/");
            this[ADDFSENTRY](base + entry);
          });
        }
        if (zip) {
          source.on("data", (chunk) => {
            zip.write(chunk);
          });
        } else {
          source.on("data", (chunk) => {
            super[WRITE](chunk);
          });
        }
      }
    };
    Pack.Sync = PackSync;
    module2.exports = Pack;
  }
});

// node_modules/fs-minipass/node_modules/minipass/index.js
var require_minipass3 = __commonJS({
  "node_modules/fs-minipass/node_modules/minipass/index.js"(exports, module2) {
    "use strict";
    var proc = typeof process === "object" && process ? process : {
      stdout: null,
      stderr: null
    };
    var EE = require("events");
    var Stream = require("stream");
    var SD = require("string_decoder").StringDecoder;
    var EOF = Symbol("EOF");
    var MAYBE_EMIT_END = Symbol("maybeEmitEnd");
    var EMITTED_END = Symbol("emittedEnd");
    var EMITTING_END = Symbol("emittingEnd");
    var EMITTED_ERROR = Symbol("emittedError");
    var CLOSED = Symbol("closed");
    var READ = Symbol("read");
    var FLUSH = Symbol("flush");
    var FLUSHCHUNK = Symbol("flushChunk");
    var ENCODING = Symbol("encoding");
    var DECODER = Symbol("decoder");
    var FLOWING = Symbol("flowing");
    var PAUSED = Symbol("paused");
    var RESUME = Symbol("resume");
    var BUFFERLENGTH = Symbol("bufferLength");
    var BUFFERPUSH = Symbol("bufferPush");
    var BUFFERSHIFT = Symbol("bufferShift");
    var OBJECTMODE = Symbol("objectMode");
    var DESTROYED = Symbol("destroyed");
    var EMITDATA = Symbol("emitData");
    var EMITEND = Symbol("emitEnd");
    var EMITEND2 = Symbol("emitEnd2");
    var ASYNC = Symbol("async");
    var defer = (fn) => Promise.resolve().then(fn);
    var doIter = global._MP_NO_ITERATOR_SYMBOLS_ !== "1";
    var ASYNCITERATOR = doIter && Symbol.asyncIterator || Symbol("asyncIterator not implemented");
    var ITERATOR = doIter && Symbol.iterator || Symbol("iterator not implemented");
    var isEndish = (ev) => ev === "end" || ev === "finish" || ev === "prefinish";
    var isArrayBuffer = (b) => b instanceof ArrayBuffer || typeof b === "object" && b.constructor && b.constructor.name === "ArrayBuffer" && b.byteLength >= 0;
    var isArrayBufferView = (b) => !Buffer.isBuffer(b) && ArrayBuffer.isView(b);
    var Pipe = class {
      constructor(src, dest, opts) {
        this.src = src;
        this.dest = dest;
        this.opts = opts;
        this.ondrain = () => src[RESUME]();
        dest.on("drain", this.ondrain);
      }
      unpipe() {
        this.dest.removeListener("drain", this.ondrain);
      }
      // istanbul ignore next - only here for the prototype
      proxyErrors() {
      }
      end() {
        this.unpipe();
        if (this.opts.end)
          this.dest.end();
      }
    };
    var PipeProxyErrors = class extends Pipe {
      unpipe() {
        this.src.removeListener("error", this.proxyErrors);
        super.unpipe();
      }
      constructor(src, dest, opts) {
        super(src, dest, opts);
        this.proxyErrors = (er) => dest.emit("error", er);
        src.on("error", this.proxyErrors);
      }
    };
    module2.exports = class Minipass extends Stream {
      constructor(options) {
        super();
        this[FLOWING] = false;
        this[PAUSED] = false;
        this.pipes = [];
        this.buffer = [];
        this[OBJECTMODE] = options && options.objectMode || false;
        if (this[OBJECTMODE])
          this[ENCODING] = null;
        else
          this[ENCODING] = options && options.encoding || null;
        if (this[ENCODING] === "buffer")
          this[ENCODING] = null;
        this[ASYNC] = options && !!options.async || false;
        this[DECODER] = this[ENCODING] ? new SD(this[ENCODING]) : null;
        this[EOF] = false;
        this[EMITTED_END] = false;
        this[EMITTING_END] = false;
        this[CLOSED] = false;
        this[EMITTED_ERROR] = null;
        this.writable = true;
        this.readable = true;
        this[BUFFERLENGTH] = 0;
        this[DESTROYED] = false;
      }
      get bufferLength() {
        return this[BUFFERLENGTH];
      }
      get encoding() {
        return this[ENCODING];
      }
      set encoding(enc) {
        if (this[OBJECTMODE])
          throw new Error("cannot set encoding in objectMode");
        if (this[ENCODING] && enc !== this[ENCODING] && (this[DECODER] && this[DECODER].lastNeed || this[BUFFERLENGTH]))
          throw new Error("cannot change encoding");
        if (this[ENCODING] !== enc) {
          this[DECODER] = enc ? new SD(enc) : null;
          if (this.buffer.length)
            this.buffer = this.buffer.map((chunk) => this[DECODER].write(chunk));
        }
        this[ENCODING] = enc;
      }
      setEncoding(enc) {
        this.encoding = enc;
      }
      get objectMode() {
        return this[OBJECTMODE];
      }
      set objectMode(om) {
        this[OBJECTMODE] = this[OBJECTMODE] || !!om;
      }
      get ["async"]() {
        return this[ASYNC];
      }
      set ["async"](a) {
        this[ASYNC] = this[ASYNC] || !!a;
      }
      write(chunk, encoding, cb) {
        if (this[EOF])
          throw new Error("write after end");
        if (this[DESTROYED]) {
          this.emit("error", Object.assign(
            new Error("Cannot call write after a stream was destroyed"),
            { code: "ERR_STREAM_DESTROYED" }
          ));
          return true;
        }
        if (typeof encoding === "function")
          cb = encoding, encoding = "utf8";
        if (!encoding)
          encoding = "utf8";
        const fn = this[ASYNC] ? defer : (f) => f();
        if (!this[OBJECTMODE] && !Buffer.isBuffer(chunk)) {
          if (isArrayBufferView(chunk))
            chunk = Buffer.from(chunk.buffer, chunk.byteOffset, chunk.byteLength);
          else if (isArrayBuffer(chunk))
            chunk = Buffer.from(chunk);
          else if (typeof chunk !== "string")
            this.objectMode = true;
        }
        if (this[OBJECTMODE]) {
          if (this.flowing && this[BUFFERLENGTH] !== 0)
            this[FLUSH](true);
          if (this.flowing)
            this.emit("data", chunk);
          else
            this[BUFFERPUSH](chunk);
          if (this[BUFFERLENGTH] !== 0)
            this.emit("readable");
          if (cb)
            fn(cb);
          return this.flowing;
        }
        if (!chunk.length) {
          if (this[BUFFERLENGTH] !== 0)
            this.emit("readable");
          if (cb)
            fn(cb);
          return this.flowing;
        }
        if (typeof chunk === "string" && // unless it is a string already ready for us to use
        !(encoding === this[ENCODING] && !this[DECODER].lastNeed)) {
          chunk = Buffer.from(chunk, encoding);
        }
        if (Buffer.isBuffer(chunk) && this[ENCODING])
          chunk = this[DECODER].write(chunk);
        if (this.flowing && this[BUFFERLENGTH] !== 0)
          this[FLUSH](true);
        if (this.flowing)
          this.emit("data", chunk);
        else
          this[BUFFERPUSH](chunk);
        if (this[BUFFERLENGTH] !== 0)
          this.emit("readable");
        if (cb)
          fn(cb);
        return this.flowing;
      }
      read(n) {
        if (this[DESTROYED])
          return null;
        if (this[BUFFERLENGTH] === 0 || n === 0 || n > this[BUFFERLENGTH]) {
          this[MAYBE_EMIT_END]();
          return null;
        }
        if (this[OBJECTMODE])
          n = null;
        if (this.buffer.length > 1 && !this[OBJECTMODE]) {
          if (this.encoding)
            this.buffer = [this.buffer.join("")];
          else
            this.buffer = [Buffer.concat(this.buffer, this[BUFFERLENGTH])];
        }
        const ret = this[READ](n || null, this.buffer[0]);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [READ](n, chunk) {
        if (n === chunk.length || n === null)
          this[BUFFERSHIFT]();
        else {
          this.buffer[0] = chunk.slice(n);
          chunk = chunk.slice(0, n);
          this[BUFFERLENGTH] -= n;
        }
        this.emit("data", chunk);
        if (!this.buffer.length && !this[EOF])
          this.emit("drain");
        return chunk;
      }
      end(chunk, encoding, cb) {
        if (typeof chunk === "function")
          cb = chunk, chunk = null;
        if (typeof encoding === "function")
          cb = encoding, encoding = "utf8";
        if (chunk)
          this.write(chunk, encoding);
        if (cb)
          this.once("end", cb);
        this[EOF] = true;
        this.writable = false;
        if (this.flowing || !this[PAUSED])
          this[MAYBE_EMIT_END]();
        return this;
      }
      // don't let the internal resume be overwritten
      [RESUME]() {
        if (this[DESTROYED])
          return;
        this[PAUSED] = false;
        this[FLOWING] = true;
        this.emit("resume");
        if (this.buffer.length)
          this[FLUSH]();
        else if (this[EOF])
          this[MAYBE_EMIT_END]();
        else
          this.emit("drain");
      }
      resume() {
        return this[RESUME]();
      }
      pause() {
        this[FLOWING] = false;
        this[PAUSED] = true;
      }
      get destroyed() {
        return this[DESTROYED];
      }
      get flowing() {
        return this[FLOWING];
      }
      get paused() {
        return this[PAUSED];
      }
      [BUFFERPUSH](chunk) {
        if (this[OBJECTMODE])
          this[BUFFERLENGTH] += 1;
        else
          this[BUFFERLENGTH] += chunk.length;
        this.buffer.push(chunk);
      }
      [BUFFERSHIFT]() {
        if (this.buffer.length) {
          if (this[OBJECTMODE])
            this[BUFFERLENGTH] -= 1;
          else
            this[BUFFERLENGTH] -= this.buffer[0].length;
        }
        return this.buffer.shift();
      }
      [FLUSH](noDrain) {
        do {
        } while (this[FLUSHCHUNK](this[BUFFERSHIFT]()));
        if (!noDrain && !this.buffer.length && !this[EOF])
          this.emit("drain");
      }
      [FLUSHCHUNK](chunk) {
        return chunk ? (this.emit("data", chunk), this.flowing) : false;
      }
      pipe(dest, opts) {
        if (this[DESTROYED])
          return;
        const ended = this[EMITTED_END];
        opts = opts || {};
        if (dest === proc.stdout || dest === proc.stderr)
          opts.end = false;
        else
          opts.end = opts.end !== false;
        opts.proxyErrors = !!opts.proxyErrors;
        if (ended) {
          if (opts.end)
            dest.end();
        } else {
          this.pipes.push(!opts.proxyErrors ? new Pipe(this, dest, opts) : new PipeProxyErrors(this, dest, opts));
          if (this[ASYNC])
            defer(() => this[RESUME]());
          else
            this[RESUME]();
        }
        return dest;
      }
      unpipe(dest) {
        const p = this.pipes.find((p2) => p2.dest === dest);
        if (p) {
          this.pipes.splice(this.pipes.indexOf(p), 1);
          p.unpipe();
        }
      }
      addListener(ev, fn) {
        return this.on(ev, fn);
      }
      on(ev, fn) {
        const ret = super.on(ev, fn);
        if (ev === "data" && !this.pipes.length && !this.flowing)
          this[RESUME]();
        else if (ev === "readable" && this[BUFFERLENGTH] !== 0)
          super.emit("readable");
        else if (isEndish(ev) && this[EMITTED_END]) {
          super.emit(ev);
          this.removeAllListeners(ev);
        } else if (ev === "error" && this[EMITTED_ERROR]) {
          if (this[ASYNC])
            defer(() => fn.call(this, this[EMITTED_ERROR]));
          else
            fn.call(this, this[EMITTED_ERROR]);
        }
        return ret;
      }
      get emittedEnd() {
        return this[EMITTED_END];
      }
      [MAYBE_EMIT_END]() {
        if (!this[EMITTING_END] && !this[EMITTED_END] && !this[DESTROYED] && this.buffer.length === 0 && this[EOF]) {
          this[EMITTING_END] = true;
          this.emit("end");
          this.emit("prefinish");
          this.emit("finish");
          if (this[CLOSED])
            this.emit("close");
          this[EMITTING_END] = false;
        }
      }
      emit(ev, data, ...extra) {
        if (ev !== "error" && ev !== "close" && ev !== DESTROYED && this[DESTROYED])
          return;
        else if (ev === "data") {
          return !data ? false : this[ASYNC] ? defer(() => this[EMITDATA](data)) : this[EMITDATA](data);
        } else if (ev === "end") {
          return this[EMITEND]();
        } else if (ev === "close") {
          this[CLOSED] = true;
          if (!this[EMITTED_END] && !this[DESTROYED])
            return;
          const ret2 = super.emit("close");
          this.removeAllListeners("close");
          return ret2;
        } else if (ev === "error") {
          this[EMITTED_ERROR] = data;
          const ret2 = super.emit("error", data);
          this[MAYBE_EMIT_END]();
          return ret2;
        } else if (ev === "resume") {
          const ret2 = super.emit("resume");
          this[MAYBE_EMIT_END]();
          return ret2;
        } else if (ev === "finish" || ev === "prefinish") {
          const ret2 = super.emit(ev);
          this.removeAllListeners(ev);
          return ret2;
        }
        const ret = super.emit(ev, data, ...extra);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [EMITDATA](data) {
        for (const p of this.pipes) {
          if (p.dest.write(data) === false)
            this.pause();
        }
        const ret = super.emit("data", data);
        this[MAYBE_EMIT_END]();
        return ret;
      }
      [EMITEND]() {
        if (this[EMITTED_END])
          return;
        this[EMITTED_END] = true;
        this.readable = false;
        if (this[ASYNC])
          defer(() => this[EMITEND2]());
        else
          this[EMITEND2]();
      }
      [EMITEND2]() {
        if (this[DECODER]) {
          const data = this[DECODER].end();
          if (data) {
            for (const p of this.pipes) {
              p.dest.write(data);
            }
            super.emit("data", data);
          }
        }
        for (const p of this.pipes) {
          p.end();
        }
        const ret = super.emit("end");
        this.removeAllListeners("end");
        return ret;
      }
      // const all = await stream.collect()
      collect() {
        const buf = [];
        if (!this[OBJECTMODE])
          buf.dataLength = 0;
        const p = this.promise();
        this.on("data", (c2) => {
          buf.push(c2);
          if (!this[OBJECTMODE])
            buf.dataLength += c2.length;
        });
        return p.then(() => buf);
      }
      // const data = await stream.concat()
      concat() {
        return this[OBJECTMODE] ? Promise.reject(new Error("cannot concat in objectMode")) : this.collect().then((buf) => this[OBJECTMODE] ? Promise.reject(new Error("cannot concat in objectMode")) : this[ENCODING] ? buf.join("") : Buffer.concat(buf, buf.dataLength));
      }
      // stream.promise().then(() => done, er => emitted error)
      promise() {
        return new Promise((resolve7, reject) => {
          this.on(DESTROYED, () => reject(new Error("stream destroyed")));
          this.on("error", (er) => reject(er));
          this.on("end", () => resolve7());
        });
      }
      // for await (let chunk of stream)
      [ASYNCITERATOR]() {
        const next = () => {
          const res = this.read();
          if (res !== null)
            return Promise.resolve({ done: false, value: res });
          if (this[EOF])
            return Promise.resolve({ done: true });
          let resolve7 = null;
          let reject = null;
          const onerr = (er) => {
            this.removeListener("data", ondata);
            this.removeListener("end", onend);
            reject(er);
          };
          const ondata = (value) => {
            this.removeListener("error", onerr);
            this.removeListener("end", onend);
            this.pause();
            resolve7({ value, done: !!this[EOF] });
          };
          const onend = () => {
            this.removeListener("error", onerr);
            this.removeListener("data", ondata);
            resolve7({ done: true });
          };
          const ondestroy = () => onerr(new Error("stream destroyed"));
          return new Promise((res2, rej) => {
            reject = rej;
            resolve7 = res2;
            this.once(DESTROYED, ondestroy);
            this.once("error", onerr);
            this.once("end", onend);
            this.once("data", ondata);
          });
        };
        return { next };
      }
      // for (let chunk of stream)
      [ITERATOR]() {
        const next = () => {
          const value = this.read();
          const done = value === null;
          return { value, done };
        };
        return { next };
      }
      destroy(er) {
        if (this[DESTROYED]) {
          if (er)
            this.emit("error", er);
          else
            this.emit(DESTROYED);
          return this;
        }
        this[DESTROYED] = true;
        this.buffer.length = 0;
        this[BUFFERLENGTH] = 0;
        if (typeof this.close === "function" && !this[CLOSED])
          this.close();
        if (er)
          this.emit("error", er);
        else
          this.emit(DESTROYED);
        return this;
      }
      static isStream(s) {
        return !!s && (s instanceof Minipass || s instanceof Stream || s instanceof EE && (typeof s.pipe === "function" || // readable
        typeof s.write === "function" && typeof s.end === "function"));
      }
    };
  }
});

// node_modules/fs-minipass/index.js
var require_fs_minipass = __commonJS({
  "node_modules/fs-minipass/index.js"(exports) {
    "use strict";
    var MiniPass = require_minipass3();
    var EE = require("events").EventEmitter;
    var fs6 = require("fs");
    var writev = fs6.writev;
    if (!writev) {
      const binding = process.binding("fs");
      const FSReqWrap = binding.FSReqWrap || binding.FSReqCallback;
      writev = (fd, iovec, pos, cb) => {
        const done = (er, bw) => cb(er, bw, iovec);
        const req = new FSReqWrap();
        req.oncomplete = done;
        binding.writeBuffers(fd, iovec, pos, req);
      };
    }
    var _autoClose = Symbol("_autoClose");
    var _close = Symbol("_close");
    var _ended = Symbol("_ended");
    var _fd = Symbol("_fd");
    var _finished = Symbol("_finished");
    var _flags = Symbol("_flags");
    var _flush = Symbol("_flush");
    var _handleChunk = Symbol("_handleChunk");
    var _makeBuf = Symbol("_makeBuf");
    var _mode = Symbol("_mode");
    var _needDrain = Symbol("_needDrain");
    var _onerror = Symbol("_onerror");
    var _onopen = Symbol("_onopen");
    var _onread = Symbol("_onread");
    var _onwrite = Symbol("_onwrite");
    var _open = Symbol("_open");
    var _path = Symbol("_path");
    var _pos = Symbol("_pos");
    var _queue = Symbol("_queue");
    var _read = Symbol("_read");
    var _readSize = Symbol("_readSize");
    var _reading = Symbol("_reading");
    var _remain = Symbol("_remain");
    var _size = Symbol("_size");
    var _write = Symbol("_write");
    var _writing = Symbol("_writing");
    var _defaultFlag = Symbol("_defaultFlag");
    var _errored = Symbol("_errored");
    var ReadStream3 = class extends MiniPass {
      constructor(path32, opt) {
        opt = opt || {};
        super(opt);
        this.readable = true;
        this.writable = false;
        if (typeof path32 !== "string")
          throw new TypeError("path must be a string");
        this[_errored] = false;
        this[_fd] = typeof opt.fd === "number" ? opt.fd : null;
        this[_path] = path32;
        this[_readSize] = opt.readSize || 16 * 1024 * 1024;
        this[_reading] = false;
        this[_size] = typeof opt.size === "number" ? opt.size : Infinity;
        this[_remain] = this[_size];
        this[_autoClose] = typeof opt.autoClose === "boolean" ? opt.autoClose : true;
        if (typeof this[_fd] === "number")
          this[_read]();
        else
          this[_open]();
      }
      get fd() {
        return this[_fd];
      }
      get path() {
        return this[_path];
      }
      write() {
        throw new TypeError("this is a readable stream");
      }
      end() {
        throw new TypeError("this is a readable stream");
      }
      [_open]() {
        fs6.open(this[_path], "r", (er, fd) => this[_onopen](er, fd));
      }
      [_onopen](er, fd) {
        if (er)
          this[_onerror](er);
        else {
          this[_fd] = fd;
          this.emit("open", fd);
          this[_read]();
        }
      }
      [_makeBuf]() {
        return Buffer.allocUnsafe(Math.min(this[_readSize], this[_remain]));
      }
      [_read]() {
        if (!this[_reading]) {
          this[_reading] = true;
          const buf = this[_makeBuf]();
          if (buf.length === 0)
            return process.nextTick(() => this[_onread](null, 0, buf));
          fs6.read(this[_fd], buf, 0, buf.length, null, (er, br, buf2) => this[_onread](er, br, buf2));
        }
      }
      [_onread](er, br, buf) {
        this[_reading] = false;
        if (er)
          this[_onerror](er);
        else if (this[_handleChunk](br, buf))
          this[_read]();
      }
      [_close]() {
        if (this[_autoClose] && typeof this[_fd] === "number") {
          const fd = this[_fd];
          this[_fd] = null;
          fs6.close(fd, (er) => er ? this.emit("error", er) : this.emit("close"));
        }
      }
      [_onerror](er) {
        this[_reading] = true;
        this[_close]();
        this.emit("error", er);
      }
      [_handleChunk](br, buf) {
        let ret = false;
        this[_remain] -= br;
        if (br > 0)
          ret = super.write(br < buf.length ? buf.slice(0, br) : buf);
        if (br === 0 || this[_remain] <= 0) {
          ret = false;
          this[_close]();
          super.end();
        }
        return ret;
      }
      emit(ev, data) {
        switch (ev) {
          case "prefinish":
          case "finish":
            break;
          case "drain":
            if (typeof this[_fd] === "number")
              this[_read]();
            break;
          case "error":
            if (this[_errored])
              return;
            this[_errored] = true;
            return super.emit(ev, data);
          default:
            return super.emit(ev, data);
        }
      }
    };
    var ReadStreamSync = class extends ReadStream3 {
      [_open]() {
        let threw = true;
        try {
          this[_onopen](null, fs6.openSync(this[_path], "r"));
          threw = false;
        } finally {
          if (threw)
            this[_close]();
        }
      }
      [_read]() {
        let threw = true;
        try {
          if (!this[_reading]) {
            this[_reading] = true;
            do {
              const buf = this[_makeBuf]();
              const br = buf.length === 0 ? 0 : fs6.readSync(this[_fd], buf, 0, buf.length, null);
              if (!this[_handleChunk](br, buf))
                break;
            } while (true);
            this[_reading] = false;
          }
          threw = false;
        } finally {
          if (threw)
            this[_close]();
        }
      }
      [_close]() {
        if (this[_autoClose] && typeof this[_fd] === "number") {
          const fd = this[_fd];
          this[_fd] = null;
          fs6.closeSync(fd);
          this.emit("close");
        }
      }
    };
    var WriteStream = class extends EE {
      constructor(path32, opt) {
        opt = opt || {};
        super(opt);
        this.readable = false;
        this.writable = true;
        this[_errored] = false;
        this[_writing] = false;
        this[_ended] = false;
        this[_needDrain] = false;
        this[_queue] = [];
        this[_path] = path32;
        this[_fd] = typeof opt.fd === "number" ? opt.fd : null;
        this[_mode] = opt.mode === void 0 ? 438 : opt.mode;
        this[_pos] = typeof opt.start === "number" ? opt.start : null;
        this[_autoClose] = typeof opt.autoClose === "boolean" ? opt.autoClose : true;
        const defaultFlag = this[_pos] !== null ? "r+" : "w";
        this[_defaultFlag] = opt.flags === void 0;
        this[_flags] = this[_defaultFlag] ? defaultFlag : opt.flags;
        if (this[_fd] === null)
          this[_open]();
      }
      emit(ev, data) {
        if (ev === "error") {
          if (this[_errored])
            return;
          this[_errored] = true;
        }
        return super.emit(ev, data);
      }
      get fd() {
        return this[_fd];
      }
      get path() {
        return this[_path];
      }
      [_onerror](er) {
        this[_close]();
        this[_writing] = true;
        this.emit("error", er);
      }
      [_open]() {
        fs6.open(
          this[_path],
          this[_flags],
          this[_mode],
          (er, fd) => this[_onopen](er, fd)
        );
      }
      [_onopen](er, fd) {
        if (this[_defaultFlag] && this[_flags] === "r+" && er && er.code === "ENOENT") {
          this[_flags] = "w";
          this[_open]();
        } else if (er)
          this[_onerror](er);
        else {
          this[_fd] = fd;
          this.emit("open", fd);
          this[_flush]();
        }
      }
      end(buf, enc) {
        if (buf)
          this.write(buf, enc);
        this[_ended] = true;
        if (!this[_writing] && !this[_queue].length && typeof this[_fd] === "number")
          this[_onwrite](null, 0);
        return this;
      }
      write(buf, enc) {
        if (typeof buf === "string")
          buf = Buffer.from(buf, enc);
        if (this[_ended]) {
          this.emit("error", new Error("write() after end()"));
          return false;
        }
        if (this[_fd] === null || this[_writing] || this[_queue].length) {
          this[_queue].push(buf);
          this[_needDrain] = true;
          return false;
        }
        this[_writing] = true;
        this[_write](buf);
        return true;
      }
      [_write](buf) {
        fs6.write(this[_fd], buf, 0, buf.length, this[_pos], (er, bw) => this[_onwrite](er, bw));
      }
      [_onwrite](er, bw) {
        if (er)
          this[_onerror](er);
        else {
          if (this[_pos] !== null)
            this[_pos] += bw;
          if (this[_queue].length)
            this[_flush]();
          else {
            this[_writing] = false;
            if (this[_ended] && !this[_finished]) {
              this[_finished] = true;
              this[_close]();
              this.emit("finish");
            } else if (this[_needDrain]) {
              this[_needDrain] = false;
              this.emit("drain");
            }
          }
        }
      }
      [_flush]() {
        if (this[_queue].length === 0) {
          if (this[_ended])
            this[_onwrite](null, 0);
        } else if (this[_queue].length === 1)
          this[_write](this[_queue].pop());
        else {
          const iovec = this[_queue];
          this[_queue] = [];
          writev(
            this[_fd],
            iovec,
            this[_pos],
            (er, bw) => this[_onwrite](er, bw)
          );
        }
      }
      [_close]() {
        if (this[_autoClose] && typeof this[_fd] === "number") {
          const fd = this[_fd];
          this[_fd] = null;
          fs6.close(fd, (er) => er ? this.emit("error", er) : this.emit("close"));
        }
      }
    };
    var WriteStreamSync = class extends WriteStream {
      [_open]() {
        let fd;
        if (this[_defaultFlag] && this[_flags] === "r+") {
          try {
            fd = fs6.openSync(this[_path], this[_flags], this[_mode]);
          } catch (er) {
            if (er.code === "ENOENT") {
              this[_flags] = "w";
              return this[_open]();
            } else
              throw er;
          }
        } else
          fd = fs6.openSync(this[_path], this[_flags], this[_mode]);
        this[_onopen](null, fd);
      }
      [_close]() {
        if (this[_autoClose] && typeof this[_fd] === "number") {
          const fd = this[_fd];
          this[_fd] = null;
          fs6.closeSync(fd);
          this.emit("close");
        }
      }
      [_write](buf) {
        let threw = true;
        try {
          this[_onwrite](
            null,
            fs6.writeSync(this[_fd], buf, 0, buf.length, this[_pos])
          );
          threw = false;
        } finally {
          if (threw)
            try {
              this[_close]();
            } catch (_) {
            }
        }
      }
    };
    exports.ReadStream = ReadStream3;
    exports.ReadStreamSync = ReadStreamSync;
    exports.WriteStream = WriteStream;
    exports.WriteStreamSync = WriteStreamSync;
  }
});

// node_modules/tar/lib/parse.js
var require_parse2 = __commonJS({
  "node_modules/tar/lib/parse.js"(exports, module2) {
    "use strict";
    var warner = require_warn_mixin();
    var 