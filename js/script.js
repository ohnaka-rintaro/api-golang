// 即時関数にしてもいいかもね
(function () {
    console.log('即時関数のテスト');
    let titleId = document.getElementById('title');
    let h1Title = document.createElement('h1');
    let titleText = document.createTextNode('APIのテストページ');
    titleId.appendChild(h1Title);
    h1Title.appendChild(titleText);
}());
