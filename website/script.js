var toggle = document.querySelector('.mobile-toggle');
var actions = document.querySelector('.header-actions');

toggle.addEventListener('click', function () {
    actions.classList.toggle('open');
});
