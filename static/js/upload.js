function Uploader() {
    this.init()
}

Uploader.prototype = {
    fileTypes: ['image/jpeg', 'image/png'],
    fileTypeErrorMsg: "Incorrect File type. Only JPEG and PNG files",
    defaultText: "Select Image",
    maxSize: (2 * 1024 * 1024),
    maxSizeMsg: "File too large. Max size 2MB",
    updateFileInfo: function(e) {
        var t = e.value;
        var n = t.match(/([^\/\\]+)$/);
        var r;
        if (n == null) {
            r = this.defaultText;
        } else {
            r = n[1];
        }
        $('label[for^="' + e.id + '"]').text(r);
        var i = $('form#uploader input.uploader');
        var s = true;
        for (var o = 0; o < i.length; o++) {
            if (i[o].value == "") {
                s = false;
            }
        }
        if (s) {
            $('#upload').removeAttr("disabled")
        } else {
            $('#upload').attr('disabled', 'disabled')
        }
    },
    fileCheck: function(e) {
        if ($('input[id="' + e.id + '"]')[0].files.length > 0) {
            $('div[id="fileError"]').addClass("hide");
            var t = $('input[id="' + e.id + '"]')[0].files[0].size;
            var n = $('input[id="' + e.id + '"]')[0].files[0].type;
            if (t > this.maxSize) {
                $('input[id="' + e.id + '"]')[0].type = "text";
                $('input[id="' + e.id + '"]')[0].type = "file";
                $('p[id="fileMessage"]').html(this.maxSizeMsg);
                $('div[id="fileError"]').removeClass("hide");
                return
            }
            console.log(n);
        	if (this.fileTypes.indexOf(n) > -1) {
        		$('div[id="fileError"]').addClass("hide");
        		return;
        	} else {
        		$('input[id="' + e.id + '"]')[0].type = "text";
        		$('input[id="' + e.id + '"]')[0].type = "file";
        		$('p[id="fileMessage"]').html(this.fileTypeErrorMsg);
        		$('div[id="fileError"]').removeClass("hide");
        	}
        }
    },
    init: function() {
        $('button#upload').attr('disabled', 'disabled');
        b =$('button#upload');
        i = $('input.uploader');
        f = $('form#uploader');
        if (b.length == 0 || i.length == 0 || f.length == 0) {
            console.log('Upload form must have an id of "uploader"\nFile input must have a class of "uploader"\nSubmit button must hav and id of "upload"');
        } else {
            $('button#upload').click(function(e) {
                $('#importModal').modal('hide');
                $('#content').addClass("hide");
                $('div[id="uploadSpinner"]').removeClass("hide");
            });
            $('input.uploader').change(function() {
                uploader.fileCheck(this);
                uploader.updateFileInfo(this);
            });
        }
    }
}

var uploader = new Uploader();
