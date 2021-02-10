window.downloads = {
    android: 'https://wwi.lanzous.com/iZxdzleszpi',
    windows: 'https://wwe.lanzous.com/i4ZDUfwxs3g',
    linux: 'https://wwe.lanzous.com/iom0Bfwxt9i',
    mac: 'https://wwe.lanzous.com/iKnR4fwxs0d',
}

function changeLanguage(me) {
    switch (me.value) {
        case 'en':
            location.href = '/';
            break;
        case 'zh':
            location.href = '/zh/index.html';
            break;
    }
    localStorage.setItem('language','select')
}
