$(function() {
	confirmacao()
});

function confirmacao() {
	var textoOriginal = $("#texto-turnos").text()

	$('#enviar').click(function(e) {
		var turnos = $.map($("#turnos > input:checked + label"), function(x) {
			return $(x).text()
		});

		if (turnos.length == 0) {
			$("#texto-turnos").text("Você não selecionou nenhum turno!");
		} else {
			$("#texto-turnos").text(textoOriginal)
		}
		$("#confirmacao-turnos").html(turnos.map(function(turno) {return "<li>" + turno + "</li>"}));
		$("#confirmacao-zona").html($("#zona > :selected").html());

		e.preventDefault();
		$('#confirmacao').dialog({
			resizable: false,
			dialogClass: "dialog",
			height: 330,
			width: 480,
			modal: true,
			title: "Confirmação",
			buttons: {
				"Não": function() {
					$(this).dialog("close");
				},
				"Sim": function() {
					$("#formulario").submit();
					$(this).dialog("close");
				}
			}
		});
	});
}
