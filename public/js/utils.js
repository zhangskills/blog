// 对Date的扩展，将 Date 转化为指定格式的String 
// 月(M)、日(d)、小时(h)、分(m)、秒(s)、季度(q) 可以用 1-2 个占位符， 
// 年(y)可以用 1-4 个占位符，毫秒(S)只能用 1 个占位符(是 1-3 位的数字) 
// 例子： 
// (new Date()).Format("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423 
// (new Date()).Format("yyyy-M-d h:m:s.S")      ==> 2006-7-2 8:9:4.18 
Date.prototype.Format = function(fmt) { //author: meizz 
	var o = {
		"M+": this.getMonth() + 1, //月份 
		"d+": this.getDate(), //日 
		"h+": this.getHours(), //小时 
		"m+": this.getMinutes(), //分 
		"s+": this.getSeconds(), //秒 
		"q+": Math.floor((this.getMonth() + 3) / 3), //季度 
		"S": this.getMilliseconds() //毫秒 
	};
	if (/(y+)/.test(fmt))
		fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
	for (var k in o) {
		if (new RegExp("(" + k + ")").test(fmt)) {
			fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
		}
	}
	return fmt;
};
var refreshDigitsImg = function(e) {
	var imgE;
	if (e) {
		imgE = e.find('.refresh-digits img');
	} else {
		imgE = $('.refresh-digits img');
	}
	imgE.attr('src', '/digits?' + Math.random());
};
var initForm = function(formE, cbFn) {
	var inputsE = formE.find(':input').not('.btn');
	inputsE.keydown(function() {
		var t = $(this);
		if (t.attr("data-original-title")) {
			t.attr('data-original-title', '').parent().removeClass('has-error');
			t.tooltip('destroy');
		}
	});
	formE.find('.refresh-digits').click(function() {
		refreshDigitsImg(formE)
	});

	var showErr = function(e, errMap) {
		if (errMap) {
			for (var k in errMap) {
				var inputE = e.find('[name="' + k + '"]').attr('title', errMap[k].Message);
				inputE.parent().addClass("has-error");
				inputE.tooltip('show');
			}
		}
	};
	formE.submit(function() {
		$(this).ajaxSubmit(function(json) {
			if (json.msg) {
				alert(json.msg);
			}
			if (json.err) {
				refreshDigitsImg(formE);
				showErr(formE, json.errMap)
			} else if (cbFn) {
				cbFn(json);
			}
		});
		return false;
	});
};

var resetForm = function(formE) {
	var inputsE = formE.find(':input').not('.btn');
	inputsE.each(function() {
		$(this).val('').tooltip('destroy').attr('data-original-title', '').parent().removeClass('has-error');
	});
};
var showDateTime = function(t) {
	return new Date(t).Format("yyyy-MM-dd hh:mm:ss");
};
var uploadImg = function(fn) {
	var formE = $('<form class="hidden" action="http://up.qiniu.com" method="post" enctype="multipart/form-data"></form>')
	formE.append('<input type="hidden" name="key"/>');
	formE.append('<input type="hidden" name="token"/>');
	formE.append('<input type="file" name="file"/>');
	formE.find("[name=file]").click().change(function() {
		var fileName = $(this).val();
		$.getJSON("/upload/token", {
			fileName: fileName
		}, function(json) {
			if (json.err) {
				alert(json.msg)
			} else {
				formE.find("[name=key]").val(json.msg.key);
				formE.find("[name=token]").val(json.msg.token);

				formE.ajaxSubmit(fn);
			}
		});
	});
};