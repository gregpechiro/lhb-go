$( document ).ready( function() {
    var $container = $('.isotope').isotope({
        itemSelector: '.item',
        layoutMode: 'fitRows'
    });

    $('button.filter').click(function() {
        $('button.filter').removeClass('active');
        $(this).addClass('active');
        var filterValue = $(this).attr('data-filter');
        $container.isotope({ filter: filterValue });
    });

    setTimeout(function() {
        $container.isotope({ filter: '*' });
    }, 1000);
});
