$('a[id^="brochure-"]').on('click', function() {
    var title = document.getElementById('brochure-modal-title');
    title.innerText = this.dataset.title;
    var body = document.getElementById('modal-body');
    body.data = this.dataset.body;
    $('#brochure-modal').modal('toggle');
});
