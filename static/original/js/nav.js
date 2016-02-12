function getPageName() {
    var page = window.location.pathname.split('/');
    page = page[page.length - 1];
    page = page.split('.')[0];
    if (page.indexOf('-') > 0) {
        pageList = page.split('-');
        for(i = 0; i < pageList.length; i++) {
            pageList[i] = pageList[i].charAt(0).toUpperCase() + pageList[i].substring(1,pageList[i].length);
        }
        page = pageList.join(' ');
    } else {
        page = page.charAt(0).toUpperCase() + page.substring(1,page.length);
    }
    return page;
}

var pageName = getPageName();
document.getElementById('pageName').innerHTML = (pageName === '') ? 'Home' : pageName;

function setSelected() {
    if (pageName !== '') {
        var listItems = document.querySelectorAll('ul.nav > li');
        var loc = window.location.href;

        for (var i = 0; i < listItems.length; i++) {
            var aLoc = listItems[i].children[0].href;
            if (aLoc === loc) {
                listItems[i].className += 'active';
            }
        }
    } else {
        document.querySelector('ul.nav > li a[href="/"]').parentElement.className += 'active';
    }
}
setSelected();
