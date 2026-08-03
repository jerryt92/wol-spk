Ext.ns("io.github.jerryt92.spk.wol");

io.github.jerryt92.spk.wol = Ext.extend(SYNO.SDS.AppInstance, {
  appWindowName: "io.github.jerryt92.spk.wol.MainWindow",

  constructor: function(config) {
    io.github.jerryt92.spk.wol.superclass.constructor.call(this, config);
  }
});

io.github.jerryt92.spk.wol.normalizeLang = function(value) {
  var lang = String(value || "").toLowerCase().replace("_", "-");
  if (!lang) {
    return "";
  }
  if (lang === "chs" || lang === "zh-cn" || lang.indexOf("zh-hans") === 0) {
    return "zh-CN";
  }
  if (lang === "cht" || lang === "zh-tw" || lang === "zh-hk" || lang === "zh-mo" || lang.indexOf("zh-hant") === 0) {
    return "zh-TW";
  }
  if (lang === "enu" || lang.indexOf("en") === 0) {
    return "en";
  }
  if (lang === "fre" || lang === "fra" || lang.indexOf("fr") === 0) {
    return "fr";
  }
  if (lang === "jpn" || lang === "jp" || lang.indexOf("ja") === 0) {
    return "ja";
  }
  if (lang === "kor" || lang === "krn" || lang === "kr" || lang === "korean" || lang.indexOf("ko") === 0) {
    return "ko";
  }
  return "";
};

io.github.jerryt92.spk.wol.pickLang = function() {
  var session = window.SYNO && window.SYNO.SDS && window.SYNO.SDS.Session;
  var raw = session && (session.lang || session.sys_lang);
  if (!String(raw || "").trim()) {
    return "";
  }
  return io.github.jerryt92.spk.wol.normalizeLang(raw) || "en";
};

io.github.jerryt92.spk.wol.MainWindow = Ext.extend(SYNO.SDS.AppWindow, {
  constructor: function(config) {
    config = config || {};
    var src = "/webman/3rdparty/WOLManager/index.cgi";
    var lang = io.github.jerryt92.spk.wol.pickLang();
    if (lang) {
      src += "?dsm_lang=" + encodeURIComponent(lang);
    }
    io.github.jerryt92.spk.wol.MainWindow.superclass.constructor.call(this, Ext.apply({
      title: "WOL Manager",
      width: 1000,
      height: 720,
      minWidth: 760,
      minHeight: 520,
      maximizable: true,
      minimizable: true,
      resizable: true,
      layout: "fit",
      items: [{
        xtype: "box",
        autoEl: {
          tag: "iframe",
          src: src,
          frameborder: "0",
          style: "width:100%;height:100%;border:0;background:#f4f8ff;"
        }
      }]
    }, config));
  }
});
