$(function() {

	//博客
	var blogModal = $('#blogModal');
	$('#add-blog-handle').click(function() {
		resetForm(blogModal);
		blogModal.find('.update').hide();
		blogModal.modal();
	});
	$('.del-blog-handle').click(function() {
		if (confirm('确认要删除该条博客吗？')) {
			var blogId = $(this).attr('i');
			$.post('/admin/blog/del', {
				id: blogId
			}, function(json) {
				if (json.err) {
					if (json.msg) {
						alert(json.msg);
					}
				} else {
					location.reload();
				}
			});
		}
	});
	$('.edit-blog-handle').click(function() {
		resetForm(blogModal);
		var id = $(this).attr('i');
		$.getJSON('/admin/blog/' + id, function(json) {
			if (json.err) {
				alert(json.msg)
			} else {
				blogModal.find('input[name="blog.Id"]').val(json.msg.id);
				blogModal.find('input[name="blog.Title"]').val(json.msg.title);
				blogModal.find('textarea[name="blog.Content"]').val(json.msg.content);
				if (json.msg.tags) {
					var tagE = blogModal.find('input[name=tagNames]');
					tagE.tagsinput('removeAll');
					for (var i = 0, len = json.msg.tags.length; i < len; i++) {
						tagE.tagsinput('add', json.msg.tags[i].name)
					}
				}
				blogModal.find('.update').show();
				blogModal.find('.created-at').text(showDateTime(json.msg.created));
				blogModal.find('.updated-at').text(showDateTime(json.msg.updated));
			}
			blogModal.modal();
		});
	});
	initForm(blogModal.find('form'), function(json) {
		location.reload();
	});

	$('#upload-blog-handle').click(function() {
		var t = $(this);
		uploadImg(function(json) {
			console.log(json);
			if (json.key) {
				t.parent().append('<div><img src="http://ww-blog.qiniudn.com/' + json.key + '" alt="" width="100"/></div>');
			}
		});
	});



	$.getJSON('/tagNames', function(json) {
		if (json.err) {
			alert(json.msg);
			return
		}
		$('#tagNames').tagsinput({
			typeaheadjs: {
				name: 'citynames',
				displayKey: 'name',
				valueKey: 'name',
				source: function(q, cb) {
					var matches = [];
					var substrRegex = new RegExp(q, 'i');
					$.each(json.msg, function(i, str) {
						if (substrRegex.test(str)) {
							matches.push({
								name: str
							});
						}
					});
					cb(matches);
				}
			}
		});
	});

});