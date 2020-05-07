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
