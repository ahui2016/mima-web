function randomPassword () {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', 'api/random-password');
    xhr.onerror = function () {
        window.alert('An error occurred during the transaction');
    }
    xhr.onload = function () {
        document.querySelector('#password').value = this.responseText;
        update_color_password();
    }
    xhr.send();
}

function update_color_password () {
    let password = document.getElementById('password').value;
    let pwd_html = "";
    for (let n of password) {
        if (isNaN(Number(n))) {
            pwd_html += `<span style="color: blue;">${n}</span>`
        } else {
            pwd_html += `<span style="color: red;">${n}</span>`;
        }
    }
    document.querySelector('#colorPassword').innerHTML = pwd_html;
}

function simpleDate(timestamp) {
    const date = new Date(timestamp);
    let year = '' + date.getFullYear(),
        month = '' + (date.getMonth() + 1),
        day = '' + date.getDate();
    if (month.length < 2) month = '0' + month;
    if (day.length < 2) day = '0' + day;
    return [year, month, day].join('-');
}
