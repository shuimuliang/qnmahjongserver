function iframePay(cfg) {
	var defCfg = {
		zindex: 10000,
		backTxt: "返回游戏",
	}
	if (cfg) {
		var type = typeof cfg;
		if (type == 'object') {
			if (cfg['zindex']) {
				defCfg.zindex = cfg['zindex'];
			}
			if (cfg['backTxt']) {
				defCfg.backTxt = cfg['backTxt'];
			}
		} else if (type == 'string') {
			defCfg.backTxt = cfg;
		} else if (type == 'number') {
			defCfg.zindex = cfg;
		}

	}
	this.cfg = defCfg;
	// 随机获取动态id
	this.baseid = new Date().getTime() + "_" + (Math.random() + "").replace(".", "");
	this.iframeBox = 'iframeBox' + this.baseid;
	this.maskid = 'mask' + this.baseid;
	this.maskzIndex = this.cfg.zindex + 1;
	this.iframeBoxzIndex = this.cfg.zindex + 2;
	this.createIframeBox = function (url) {
		var oDiv = document.createElement('div');
		oDiv.id = this.iframeBox;
		oDiv.style.cssText = "width:100%; height:100%;position:fixed;z-index:" + this.iframeBoxzIndex + ";overflow:hidden;left:0;top:0";
		var u = navigator.userAgent;
		if (u.indexOf('Android') > -1 || u.indexOf('Linux') > -1) {
			oDiv.style["-webkito-overflow-scrolling"] = "touch";
			oDiv.style["overflow"] = "hidden";
		} else {
			oDiv.style["-webkit-overflow-scrolling"] = "touch";
			oDiv.style["overflow-y"] = "scroll";
		}
		document.getElementsByTagName('body').item(0).appendChild(oDiv);

		var oScript = document.createElement("iframe");
		oScript.style.cssText = "background-color:#fff;width:100%; height:100%;border:0";
		oScript.src = url;
		oDiv.appendChild(oScript);
	}
	// 创建返回按钮
	this.createBack = function (onclose) {
		var myBackTxt = this.cfg.backTxt;
		var acBack = document.createElement('div');
		acBack.setAttribute('data-iframebox', this.iframeBox);
		acBack.setAttribute('data-maskid', this.maskid);
		acBack.style.cssText = "width:70px; height:30px;line-height: 30px;text-align:center;font-size:12px;color:#fff;background-color:rgba(0,0,0,0.4);position:fixed;border-radius:4px;right:20px;bottom:40px;fontFamily:Microsoft YaHei";
		acBack.onclick = function () {
			removeElement(document.getElementById(this.getAttribute('data-iframebox')));
			removeElement(document.getElementById(this.getAttribute('data-maskid')));
			if (onclose) {
				onclose();
			}
		};
		acBack.appendChild(document.createTextNode(myBackTxt));
		document.getElementById(this.iframeBox).appendChild(acBack);
	}
	// mask
	this.mask = function () {
		var mask = document.createElement('div');
		mask.id = this.maskid;
		mask.style.cssText = "z-index:" + this.maskzIndex + ";width: 100%;height: 100%;top: 0px;left: 0px;background-color: rgba(0, 0, 0, 0);opacity: 0.6;cursor: wait;position: fixed;";
		document.getElementsByTagName('body').item(0).appendChild(mask);
	}
	// 移除DOM元素
	function removeElement(ele) {
		if (ele == null || ele == undefined || ele == "undefined")
			return;
		var parentElement = ele.parentNode;
		if (parentElement) {
			parentElement.removeChild(ele);
		}
	}
	// iframe参数对象里边有backTxt（返回按钮的自定义文字）、baseZIndex（z-index的值）
	this.open = function (url, onclose) {
		this.mask();
		this.createIframeBox(url);
		this.createBack(onclose);
	};
	return this;
}
window.iframePay = new iframePay();