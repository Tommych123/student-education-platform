document.addEventListener('DOMContentLoaded', function() {
	var currentLocation = location.pathname;
	var navLinks = document.querySelectorAll('nav a');

	navLinks.forEach(function(link) {
		var linkPath = new URL(link.href).pathname;

		if (linkPath === currentLocation) {
			link.classList.add('active');
		}
	});
});

function copyToClipboard() {
	// Получаем элемент с ID zoom-id
	var zoomIdElement = document.getElementById("zoom-id");
	
	// Создаем временный элемент для выбора текста
	var tempInput = document.createElement("input");
	tempInput.value = zoomIdElement.textContent;
	document.body.appendChild(tempInput);
	tempInput.select();
	
	// Копируем текст в буфер обмена
	try {
		document.execCommand("copy");
	} catch (err) {
		alert("Не удалось скопировать Zoom ID.");
	}
	
	// Удаляем временный элемент
	document.body.removeChild(tempInput);
}