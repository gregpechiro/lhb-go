$('a.floorplan').on('click', function() {
    var a = $(this);
    var type = a.attr('data-type');
    var height = "";
    if (type == 'png') {
        type = 'image/png';
    } else if (type == 'jpg') {
        type = 'image/jpg';
    } else if (type == 'pdf') {
        type = 'application/pdf';
        height = '500px';
    } else {
        return
    }
    $('#floorplan-modal-title').text(a.attr('data-title'));
    $('#floorplan-modal-body').attr('data', a.attr('data-body'));
    $('#floorplan-modal-body').attr('type', a.attr('type'));
    $('#floorplan-modal-body').attr('height', height);
    $('#floorplan-modal').modal('toggle');
});
